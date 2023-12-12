package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/logic/provider"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"net/http"
)

var _ api.ServerInterface = ReprtClient{}

type ReprtClient struct {
	logger             zerolog.Logger
	deckController     decks.Controller
	providerController provider.Controller
}

func New(logger zerolog.Logger, controller decks.Controller) *ReprtClient {
	logger = logger.With().Str("module", "server").Logger()
	return &ReprtClient{
		logger:         logger,
		deckController: controller,
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

func (rc ReprtClient) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request - Invalid form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	passowrd := r.FormValue("password")

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
