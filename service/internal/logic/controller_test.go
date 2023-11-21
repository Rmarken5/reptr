package logic

import (
	"context"
	"errors"
	"github.com/google/uuid"
	dbErrors "github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/database/mocks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestLogic_InsertDeck(t *testing.T) {
	var (
		ctx        = context.Background()
		haveDeckID = uuid.NewString()
	)
	testCases := map[string]struct {
		wantDeckID             string
		wantErr                error
		mockRepositoryResponse func(mockRepo *database.MockRepository)
	}{
		"should return deck ID after successful insert": {
			wantDeckID: haveDeckID,
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertDeck(gomock.Any(), gomock.Any()).Return(haveDeckID, nil)
			},
		},
		"should return error returned from repo": {
			wantDeckID: "",
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertDeck(gomock.Any(), gomock.Any()).Return(haveDeckID, errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}
			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}

			gotDeckID, gotErr := logic.InsertDeck(ctx, models.Deck{})
			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.wantDeckID, gotDeckID)
		})
	}
}

func TestLogic_GetDecks(t *testing.T) {
	// Set up your mock repository response and other necessary variables for testing
	var (
		ctx        = context.Background()
		from       = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		to         = time.Date(2023, time.January, 5, 0, 0, 0, 0, time.UTC)
		defaultTo  = time.Date(from.Year(), from.Month(), from.Day(), 23, 59, 59, 0, from.Location())
		beforeFrom = from.Add(-1 * time.Second)
		limit      = 10
		offset     = 0
	)

	testCases := map[string]struct {
		wantDecks              []models.DeckWithCards
		wantErr                error
		toTime                 *time.Time
		mockRepositoryResponse func(mockRepo *database.MockRepository)
	}{
		"should return decks within the specified time range": {
			wantDecks: []models.DeckWithCards{
				// Define your expected DeckWithCards here
			},
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().GetWithCards(gomock.Any(), from, gomock.Any(), limit, offset).Return([]models.DeckWithCards{}, nil)
			},
			toTime: &to,
		},
		"should return cards using default to time": {
			wantDecks: []models.DeckWithCards{},
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().GetWithCards(gomock.Any(), from, &defaultTo, limit, offset).Return([]models.DeckWithCards{}, nil)
			},
			toTime: nil,
		},
		"should return error if 'to' is before 'from'": {
			wantDecks: nil,
			wantErr:   ErrInvalidToBeforeFrom,
			toTime:    &beforeFrom,
		},
		// Add more test cases as needed to cover other scenarios
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}
			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}

			gotDecks, gotErr := logic.GetDecks(ctx, from, tc.toTime, limit, offset)
			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, tc.wantDecks, gotDecks)
		})
	}
}

func TestLogic_AddCardToDeck(t *testing.T) {
	ctx := context.Background()
	deckID := "your_deck_id"
	testCard := models.Card{
		ID:        "your_card_id",
		Front:     "Front content",
		Back:      "Back content",
		Kind:      1, // Adjust according to your model
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should add card to deck successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertCards(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if card insertion fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertCards(gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
		// Add more test cases to cover other scenarios if needed
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.AddCardToDeck(ctx, deckID, testCard)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestLogic_UpdateCard(t *testing.T) {
	ctx := context.Background()
	testCard := models.Card{
		ID:        "your_card_id",
		Front:     "Updated front content",
		Back:      "Updated back content",
		Kind:      2, // Updated kind value, adjust as per your model
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should update card successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().UpdateCard(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if card update fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().UpdateCard(gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrUpdate))
			},
			wantErr: dbErrors.ErrUpdate,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.UpdateCard(ctx, testCard)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestLogic_UpvoteDeck(t *testing.T) {
	ctx := context.Background()
	deckID := "your_deck_id"
	userID := "your_user_id"

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should upvote deck successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().AddUserToUpvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if upvoting deck fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().AddUserToUpvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.UpvoteDeck(ctx, deckID, userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestLogic_RemoveUpvoteDeck(t *testing.T) {
	ctx := context.Background()
	deckID := "your_deck_id"
	userID := "your_user_id"

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should remove upvote from deck successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().RemoveUserFromUpvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if removing upvote fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().RemoveUserFromUpvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.RemoveUpvoteDeck(ctx, deckID, userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestLogic_DownvoteDeck(t *testing.T) {
	ctx := context.Background()
	deckID := "your_deck_id"
	userID := "your_user_id"

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should add downvote to deck successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().AddUserToDownvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if adding downvote fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().AddUserToDownvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
		// Add more test cases to cover other scenarios if needed
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.DownvoteDeck(ctx, deckID, userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func TestLogic_RemoveDownvoteDeck(t *testing.T) {
	ctx := context.Background()
	deckID := "your_deck_id"
	userID := "your_user_id"

	testCases := map[string]struct {
		mockRepositoryResponse func(mockRepo *database.MockRepository)
		wantErr                error
	}{
		"should remove downvote from deck successfully": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().RemoveUserFromDownvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		"should return error if removing downvote fails": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().RemoveUserFromDownvote(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
		// Add more test cases to cover other scenarios if needed
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := database.NewMockRepository(ctrl)
			if tc.mockRepositoryResponse != nil {
				tc.mockRepositoryResponse(mockRepo)
			}

			logic := Logic{repo: mockRepo, logger: zerolog.Nop()}
			gotErr := logic.RemoveDownvoteDeck(ctx, deckID, userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}
