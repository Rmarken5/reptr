package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/api"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rmarken/reptr/service/internal/web/components"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

//go:generate mockgen -destination ./mocks/sessions_mock.go -package api  github.com/gorilla/sessions Store

var _ api.ServerInterface = ReprtClient{}

const (
	IDToken         = "id_token"
	SessionTokenKey = "token"
	CookieSessionID = "reptr-session-id"
)

type ReprtClient struct {
	logger             zerolog.Logger
	deckController     decks.Controller
	providerController provider.Controller
	authenticator      auth.Authentication
	store              sessions.Store
}

func New(logger zerolog.Logger, deckController decks.Controller, providerController provider.Controller, authentication auth.Authentication, store sessions.Store) *ReprtClient {
	logger = logger.With().Str("module", "server").Logger()
	return &ReprtClient{
		logger:             logger,
		deckController:     deckController,
		providerController: providerController,
		authenticator:      authentication,
		store:              store,
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

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		log.Error().Msg("username not on context while calling AddGroup")
		http.Error(w, "username not on context", http.StatusBadRequest)
	}

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

	group, err := rc.deckController.CreateGroup(r.Context(), username, groupName.GroupName)
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
func (rc ReprtClient) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "RegistrationPage").Logger()
	log.Info().Msgf("serving registration page")
	err := components.Register(nil).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	repassword := r.PostForm.Get("repassword")

	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := components.Register(components.Banner("Must provide email")).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := components.Register(components.Banner("Must provide password")).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	//TODO: check password strength

	if password != repassword {
		w.WriteHeader(http.StatusBadRequest)
		err := components.Register(components.Banner("Passwords do not match")).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	user, registrationError, err := rc.authenticator.RegisterUser(r.Context(), email, password)
	if err != nil {
		log.Error().Err(err).Msg("while registering")
		http.Error(w, "while registering", http.StatusInternalServerError)
		return
	}

	if !registrationError.IsZero() {
		w.WriteHeader(registrationError.StatusCode)
		err := components.Register(components.Banner(registrationError.Description)).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	log.Info().Msgf("user is registered: %+v", user)
	w.WriteHeader(http.StatusCreated)
	err = components.Login(components.Banner("Registration Successful")).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rc ReprtClient) LoginPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "LoginPage").Logger()
	log.Info().Msgf("serving login page")
	err := components.Login(nil).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rc ReprtClient) Login(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "Login").Logger()
	logger.Debug().Msg("login called")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	email := r.PostForm.Get("email")
	if email == "" {
		logger.Info().Msgf("login attempt without email")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password := r.PostForm.Get("password")
	if password == "" {
		logger.Info().Msgf("login attempt without password - email: %s", email)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := rc.authenticator.PasswordCredentialsToken(r.Context(), email, password)
	if err != nil {
		logger.Error().Err(err).Msgf("error authenticating %s: %v", email, err)
		http.Error(w, "Bad request - Invalid username or password", http.StatusUnauthorized)
		return
	}

	if tokenString, ok := token.Extra(IDToken).(string); ok {
		_, err := rc.authenticator.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "unable to verify token", http.StatusInternalServerError)
			return
		}
		session, _ := rc.store.Get(r, CookieSessionID)
		session.Values[SessionTokenKey] = "Bearer " + tokenString
		session.Options.Secure = true
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "unable to save session", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization", "Bearer "+tokenString)
		http.Redirect(w, r, "/page/home", http.StatusSeeOther)
		return
	}

	logger.Error().Msgf("Bad request - token doesn't contain %s", IDToken)
	http.Error(w, fmt.Sprintf("Bad request - token doesn't contain %s", IDToken), http.StatusInternalServerError)
	return
}

func (rc ReprtClient) GetDecksForUser(w http.ResponseWriter, r *http.Request, params api.GetDecksForUserParams) {
	_ = rc.logger.With().Str("method", "GetDecksForUser")
	panic("implement me")
}

func (rc ReprtClient) HomePage(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "HomePage").Logger()
	logger.Debug().Msgf("HomePage called.")

	var userName string
	var ok bool
	if userName, ok = reptrCtx.Username(r.Context()); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong with getting username"))
		return
	}

	groups, err := rc.deckController.GetGroupsForUser(r.Context(), userName, time.Time{}, nil, 10, 0)
	if err != nil {
		http.Error(w, "while getting groups for user", toStatus(err))
		return
	}
	homeGroups := make([]components.Group, len(groups))
	for i, group := range groups {
		homeGroups[i] = homeGroupFromModel(group)
	}

	components.Home(components.HomeData{Username: userName, Groups: homeGroups}).Render(r.Context(), w)
}

func homeGroupFromModel(group models.Group) components.Group {
	return components.Group{
		GroupName: group.Name,
		NumDecks:  len(group.DeckIDs),
		NumUsers:  0,
	}
}

func (rc ReprtClient) CreateGroup(w http.ResponseWriter, r *http.Request) {
	components.CreateGroup().Render(r.Context(), w)
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
