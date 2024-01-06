package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

var _ api.ServerInterface = ReprtClient{}

const IDToken = "id_token"

type ReprtClient struct {
	logger             zerolog.Logger
	deckController     decks.Controller
	providerController provider.Controller
	authenticator      auth.Authentication
}

func New(logger zerolog.Logger, deckController decks.Controller, providerController provider.Controller, authentication auth.Authentication) *ReprtClient {
	logger = logger.With().Str("module", "server").Logger()
	return &ReprtClient{
		logger:             logger,
		deckController:     deckController,
		providerController: providerController,
		authenticator:      authentication,
	}
}

func (rc ReprtClient) GetGroups(w http.ResponseWriter, r *http.Request, params api.GetGroupsParams) {
	log := rc.logger.With().Str("method", "GetGroups").Logger()
	w.Header().Set("Content-Type", "application/json")
	groups, err := rc.deckController.GetGroups(r.Context(), params.From, params.To, params.Limit, params.Offset)
	if err != nil {
		log.Error().Err(err).Msgf("while getting groups with: %+v", params)
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("error in request with %+v", params),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}
	g := make(api.GetGroups, len(groups))
	for i, group := range groups {
		g[i] = api.GroupWithDecks{
			CreatedAt: group.CreatedAt,
			Id:        group.ID,
			Name:      group.Name,
			Decks:     decksFromDecks(group.Decks),
			UpdatedAt: group.UpdatedAt,
		}
	}
	json.NewEncoder(w).Encode(g)
}

func decksFromDecks(fromService []models.Deck) []api.Deck {
	apiDecks := make([]api.Deck, len(fromService))
	for i, deck := range fromService {
		apiDecks[i] = api.Deck{
			CreatedAt: deck.CreatedAt,
			Id:        deck.ID,
			Name:      deck.Name,
			UpdatedAt: deck.UpdatedAt,
		}
	}
	return apiDecks
}

func (rc ReprtClient) AddGroup(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "AddGroup").Logger()

	w.Header().Set("Content-Type", "application/json")

	var groupName api.GroupName
	err := json.NewDecoder(r.Body).Decode(&groupName)
	if err != nil {
		log.Error().Err(err).Msg("while trying to read request body")
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("error in reading request body"),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}
	defer r.Body.Close()

	group, err := rc.deckController.CreateGroup(r.Context(), groupName.GroupName)
	if err != nil {
		log.Error().Err(err).Msg("while trying create group")
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("while trying create group"),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(group))
}

func (rc ReprtClient) AddDeck(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "AddDeck").Logger()

	w.Header().Set("Content-Type", "application/json")

	var deckName api.DeckName
	err := json.NewDecoder(r.Body).Decode(&deckName)
	if err != nil {
		log.Error().Err(err).Msg("while trying to read request body")
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("error in reading request body"),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}
	defer r.Body.Close()

	deck, err := rc.deckController.CreateDeck(r.Context(), deckName.DeckName)
	if err != nil {
		log.Error().Err(err).Msg("while trying create deck")
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("while trying create deck"),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}

	w.Header().Set("Content-Type", "plain/text")

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(deck))
}

func (rc ReprtClient) AddDeckToGroup(w http.ResponseWriter, r *http.Request, groupId string, deckId string) {
	log := rc.logger.With().Str("method", "AddDeckToGroup").Logger()

	w.Header().Set("Content-Type", "application/json")

	err := rc.deckController.AddDeckToGroup(r.Context(), groupId, deckId)
	if err != nil {
		log.Error().Err(err).Msg("while trying create deck")
		status := toStatus(err)
		w.WriteHeader(status)
		errObj := api.ErrorObject{
			Error:      err.Error(),
			Message:    fmt.Sprintf("while trying create deck"),
			StatusCode: status,
		}
		json.NewEncoder(w).Encode(errObj)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(groupId))
}

func (rc ReprtClient) Register(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "register").Logger()
	log.Info().Msgf("calling register")

	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	repassword := r.PostForm.Get("repassword")

	if username == "" {
		http.Error(w, "must provide a username", http.StatusBadRequest)
	}

	if password == "" {
		http.Error(w, "must provide a password", http.StatusBadRequest)
	}

	// check password strength

	if password != repassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
	}

	user, err := rc.authenticator.RegisterUser(r.Context(), username, password)
	if err != nil {
		log.Error().Err(err).Msg("while registering")
		http.Error(w, "while registering", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(user))

}

func (rc ReprtClient) Login(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "Login").Logger()
	logger.Debug().Msg("login called")

	username, password, err := parseBasicAuth(r.Header.Get("Authorization"))
	if err != nil {
		logger.Error().Err(err).Msgf("Bad request - Authorization header")
		http.Error(w, "Bad request - Authorization header", http.StatusBadRequest)
		return
	}

	token, err := rc.authenticator.PasswordCredentialsToken(r.Context(), username, password)
	if err != nil {
		logger.Error().Err(err).Msgf("error authenticating %s: %v", username, err)
		http.Error(w, "Bad request - Invalid username or password", http.StatusUnauthorized)
		return
	}

	if tokenString, ok := token.Extra(IDToken).(string); ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
		return
	}

	logger.Error().Msgf("Bad request - token doesn't contain %s", IDToken)
	http.Error(w, fmt.Sprintf("Bad request - token doesn't contain %s", IDToken), http.StatusInternalServerError)
	return

}

func toStatus(err error) int {
	switch {
	case errors.Is(err, decks.ErrInvalidToBeforeFrom),
		errors.Is(err, decks.ErrInvalidGroupName),
		errors.Is(err, decks.ErrInvalidDeckName),
		errors.Is(err, decks.ErrEmptyGroupID),
		errors.Is(err, decks.ErrEmptyDeckID):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func parseBasicAuth(basicAuthHeader string) (string, string, error) {
	authParts := strings.SplitN(basicAuthHeader, " ", 2)

	if len(authParts) != 2 || authParts[0] != "Basic" {
		return "", "", errors.New("invalid basic auth header")
	}

	// Decode the base64-encoded credentials
	credentials, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		fmt.Println("Error decoding credentials:", err)
		return "", "", errors.New("error decoding credentials")

	}

	// Split the decoded credentials into username and password
	credentialsSplit := strings.SplitN(string(credentials), ":", 2)
	if len(credentialsSplit) != 2 || credentialsSplit[0] == "" || credentialsSplit[1] == "" {
		fmt.Println("Invalid credentials format")
		return "", "", errors.New("invalid credentials format")
	}
	return credentialsSplit[0], credentialsSplit[1], nil
}
