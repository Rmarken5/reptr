package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/internal/logic"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"net/http"
)

var _ api.ServerInterface = ReprtClient{}

type ReprtClient struct {
	logger     zerolog.Logger
	controller logic.Controller
}

func New(logger zerolog.Logger, controller logic.Controller) *ReprtClient {
	logger = logger.With().Str("module", "server").Logger()
	return &ReprtClient{
		logger:     logger,
		controller: controller,
	}
}

func (rc ReprtClient) GetGroups(w http.ResponseWriter, r *http.Request, params api.GetGroupsParams) {
	log := rc.logger.With().Str("method", "GetGroups").Logger()

	groups, err := rc.controller.GetGroups(r.Context(), params.From, params.To, params.Limit, params.Offset)
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
	w.Header().Set("Content-Type", "application/json")
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

	var groupName api.GroupName
	err := json.NewDecoder(r.Body).Decode(&groupName)
	if err != nil {
		log.Error().Err(err).Msg("while trying to read request")
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

	group, err := rc.controller.CreateGroup(r.Context(), groupName.GroupName)
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
	json.NewEncoder(w).Encode(group)
}

func toStatus(err error) int {
	switch {
	case errors.Is(err, logic.ErrInvalidToBeforeFrom),
		errors.Is(err, logic.ErrInvalidGroupName):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
