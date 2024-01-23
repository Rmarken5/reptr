package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/rmarken/reptr/api"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/logic/auth"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rmarken/reptr/service/internal/web/components/dumb"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"
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

func decksFromDecks(fromService []models.GetDeckResults) []api.Deck {
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

func (rc ReprtClient) GetCardInput(w http.ResponseWriter, r *http.Request, cardNum int) {
	logger := rc.logger.With().Str("method", "GetCardInput").Logger()
	logger.Info().Msgf("serving card input section")
	logger.Debug().Msgf("vals: ", r.URL.)
	dumb.CardInput(strconv.Itoa(0)).Render(r.Context(), w)
}
