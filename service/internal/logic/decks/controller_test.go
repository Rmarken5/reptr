package decks

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

func TestNew(t *testing.T) {
	l := New(zerolog.Nop(), nil)
	assert.NotNil(t, l)
}
func TestLogic_CreateDeck(t *testing.T) {
	var (
		ctx          = context.Background()
		haveDeckID   = uuid.NewString()
		haveDeckName = uuid.NewString()
	)
	testCases := map[string]struct {
		haveDeckName           string
		wantDeckID             string
		wantErr                error
		mockRepositoryResponse func(mockRepo *database.MockRepository)
	}{
		"should return deck ID after successful insert": {
			haveDeckName: haveDeckName,
			wantDeckID:   haveDeckID,
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertDeck(gomock.Any(), gomock.Any()).Return(haveDeckID, nil)
			},
		},
		"should return error returned from repo": {
			haveDeckName: haveDeckName,
			wantDeckID:   "",
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().InsertDeck(gomock.Any(), gomock.Any()).Return(haveDeckID, errors.Join(errors.New("error inserting"), dbErrors.ErrInsert))
			},
			wantErr: dbErrors.ErrInsert,
		},
		"should return ErrEmptyDeckName when deck name is empty": {
			haveDeckName: "",
			wantDeckID:   "",
			wantErr:      ErrEmptyDeckName,
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

			gotDeckID, gotErr := logic.CreateDeck(ctx, tc.haveDeckName)
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
		beforeFrom = from.Add(-1 * time.Second)
		limit      = 10
		offset     = 0
		haveErr    = errors.New("db error")
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
		"should return error from database": {
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().GetWithCards(gomock.Any(), from, gomock.Any(), limit, offset).Return(nil, haveErr)
			},
			toTime:  &to,
			wantErr: haveErr,
		},
		"should return cards using default to time": {
			wantDecks: []models.DeckWithCards{},
			mockRepositoryResponse: func(mockRepo *database.MockRepository) {
				mockRepo.EXPECT().GetWithCards(gomock.Any(), from, nil, limit, offset).Return([]models.DeckWithCards{}, nil)
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
	timeNow := time.Now().UTC()
	deckID := "your_deck_id"
	testCard := models.Card{
		ID:        "your_card_id",
		Front:     "Front content",
		Back:      "Back content",
		Kind:      1, // Adjust according to your model
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
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
	timeNow := time.Now().UTC()
	testCard := models.Card{
		ID:        "your_card_id",
		Front:     "Updated front content",
		Back:      "Updated back content",
		Kind:      2, // Updated kind value, adjust as per your model
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
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

func TestLogic_CreateGroup(t *testing.T) {
	var (
		haveErr     = errors.New("db error")
		haveGroupID = uuid.NewString()
	)

	testCases := map[string]struct {
		mockStore    func(mock *database.MockRepository)
		haveGroupID  string
		haveUsername string
		wantErr      error
	}{
		"should return groupID when group is inserted": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().InsertGroup(gomock.Any(), gomock.Any()).Return(haveGroupID, nil)
				mock.EXPECT().AddUserAsMemberOfGroup(gomock.Any(), gomock.Any(), haveGroupID).Return(nil)
			},
			haveGroupID:  haveGroupID,
			haveUsername: uuid.NewString(),
		},
		"should return error when adding group membership": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().InsertGroup(gomock.Any(), gomock.Any()).Return(haveGroupID, nil)
				mock.EXPECT().AddUserAsMemberOfGroup(gomock.Any(), gomock.Any(), haveGroupID).Return(haveErr)
			},
			haveGroupID:  haveGroupID,
			haveUsername: uuid.NewString(),
			wantErr:      haveErr,
		},
		"should return ErrEmptyUserName when username is not on context": {
			haveUsername: "",
			wantErr:      ErrEmptyUsername,
		},
		"should return ErrInvalidGroupName when groupName empty string": {
			haveUsername: uuid.NewString(),
			haveGroupID:  "",
			wantErr:      ErrInvalidGroupName,
		},
		"should return err when database layer returns err": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().InsertGroup(gomock.Any(), gomock.Any()).Return("", haveErr)
			},
			haveGroupID:  uuid.NewString(),
			haveUsername: uuid.NewString(),
			wantErr:      haveErr,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDB := database.NewMockRepository(ctrl)

			if tc.mockStore != nil {
				tc.mockStore(mockDB)
			}

			logic := Logic{repo: mockDB, logger: zerolog.Nop()}

			gotGroupID, err := logic.CreateGroup(context.Background(), tc.haveUsername, tc.haveGroupID)

			assert.ErrorIs(t, err, tc.wantErr)
			if err != nil {
				assert.Empty(t, gotGroupID)
			} else {
				assert.NotEmpty(t, gotGroupID)
			}
		})
	}
}

func TestLogic_AddDeckToGroup(t *testing.T) {
	var (
		haveErr = errors.New("db error")
	)

	testCases := map[string]struct {
		mockStore   func(mock *database.MockRepository)
		haveGroupID string
		haveDeckID  string
		wantErr     error
	}{
		"should return nil when deck is added to group": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().AddDeckToGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			haveGroupID: uuid.NewString(),
			haveDeckID:  uuid.NewString(),
		},
		"should return ErrEmptyGroupID when group ID is empty": {
			haveGroupID: "",
			wantErr:     ErrEmptyGroupID,
		},
		"should return ErrEmptyDeckID when deckID is empty string": {
			haveGroupID: uuid.NewString(),
			haveDeckID:  "",
			wantErr:     ErrEmptyDeckID,
		},
		"should return err when database layer returns err": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().AddDeckToGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(haveErr)
			},
			haveGroupID: uuid.NewString(),
			haveDeckID:  uuid.NewString(),
			wantErr:     haveErr,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDB := database.NewMockRepository(ctrl)

			if tc.mockStore != nil {
				tc.mockStore(mockDB)
			}

			logic := Logic{repo: mockDB, logger: zerolog.Nop()}

			err := logic.AddDeckToGroup(context.Background(), tc.haveGroupID, tc.haveDeckID)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestLogic_GetGroups(t *testing.T) {
	var (
		timeNow    = time.Now().UTC().Truncate(time.Millisecond)
		invalidTo  = timeNow.Add(-1 * time.Second)
		username   = uuid.NewString()
		deckOneID  = uuid.NewString()
		haveErr    = errors.New("db error")
		haveGroups = []models.GroupWithDecks{
			{
				Group: models.Group{
					ID:         uuid.NewString(),
					Name:       uuid.NewString(),
					CreatedBy:  username,
					Moderators: []string{username},
					DeckIDs:    []string{deckOneID},
					CreatedAt:  timeNow,
					UpdatedAt:  timeNow,
					DeletedAt:  nil,
				},
				Decks: []models.GetDeckResults{
					{
						ID:        deckOneID,
						Name:      uuid.NewString(),
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					},
				},
			},
		}
	)

	testCases := map[string]struct {
		haveFrom   time.Time
		haveTo     *time.Time
		mockRepo   func(mock *database.MockRepository)
		wantGroups []models.GroupWithDecks
		wantErr    error
	}{
		"should return groups when database returns result": {
			haveFrom:   time.Time{},
			haveTo:     &timeNow,
			wantGroups: haveGroups,
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsWithDecks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(haveGroups, nil)
			},
		},
		"should return error when database returns error": {
			haveFrom:   time.Time{},
			haveTo:     nil,
			wantGroups: []models.GroupWithDecks(nil),
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsWithDecks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.GroupWithDecks(nil), haveErr)
			},
			wantErr: haveErr,
		},
		"should return ErrInvalidToBeforeFrom when to is before from": {
			haveFrom:   timeNow,
			haveTo:     &invalidTo,
			wantGroups: []models.GroupWithDecks(nil),
			wantErr:    ErrInvalidToBeforeFrom,
		},
		"should return empty slice no results come from database": {
			haveFrom: timeNow,
			haveTo:   nil,
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsWithDecks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.GroupWithDecks(nil), dbErrors.ErrNoResults)
			},
			wantGroups: []models.GroupWithDecks(nil),
			wantErr:    nil,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDB := database.NewMockRepository(ctrl)

			if tc.mockRepo != nil {
				tc.mockRepo(mockDB)
			}

			logic := Logic{repo: mockDB, logger: zerolog.Nop()}

			gotGroups, err := logic.GetGroups(context.Background(), tc.haveFrom, tc.haveTo, 0, 0)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantGroups, gotGroups)

		})
	}
}

func TestLogic_GetGroupsForUser(t *testing.T) {
	var (
		timeNow    = time.Now().UTC().Truncate(time.Millisecond)
		invalidTo  = timeNow.Add(-1 * time.Second)
		username   = uuid.NewString()
		deckOneID  = uuid.NewString()
		deckTwoID  = uuid.NewString()
		haveErr    = errors.New("db error")
		haveGroups = []models.Group{
			{
				ID:         uuid.NewString(),
				Name:       uuid.NewString(),
				CreatedBy:  username,
				Moderators: []string{username},
				DeckIDs:    []string{deckOneID},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
			{
				ID:         uuid.NewString(),
				Name:       uuid.NewString(),
				CreatedBy:  username,
				Moderators: []string{username},
				DeckIDs:    []string{deckTwoID},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
		}
	)

	testCases := map[string]struct {
		haveUser   string
		haveFrom   time.Time
		haveTo     *time.Time
		mockRepo   func(mock *database.MockRepository)
		wantGroups []models.Group
		wantErr    error
	}{
		"should return groups when database returns result": {
			haveUser:   username,
			haveFrom:   time.Time{},
			haveTo:     &timeNow,
			wantGroups: haveGroups,
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsForUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(haveGroups, nil)
			},
		},
		"should return error when database returns error": {
			haveUser:   username,
			haveFrom:   time.Time{},
			haveTo:     nil,
			wantGroups: []models.Group(nil),
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsForUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Group(nil), haveErr)
			},
			wantErr: haveErr,
		},
		"should return ErrInvalidToBeforeFrom when to is before from": {
			haveUser:   username,
			haveFrom:   timeNow,
			haveTo:     &invalidTo,
			wantGroups: []models.Group(nil),
			wantErr:    ErrInvalidToBeforeFrom,
		},
		"should return empty slice no results come from database": {
			haveUser: username,
			haveFrom: timeNow,
			haveTo:   nil,
			mockRepo: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupsForUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.Group(nil), dbErrors.ErrNoResults)
			},
			wantGroups: []models.Group(nil),
			wantErr:    nil,
		},
		"should return ErrEmptyUsername when username is empty": {
			haveUser:   "",
			haveFrom:   timeNow,
			haveTo:     nil,
			wantGroups: []models.Group(nil),
			wantErr:    ErrEmptyUsername,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDB := database.NewMockRepository(ctrl)

			if tc.mockRepo != nil {
				tc.mockRepo(mockDB)
			}

			logic := Logic{repo: mockDB, logger: zerolog.Nop()}

			gotGroups, err := logic.GetGroupsForUser(context.Background(), tc.haveUser, tc.haveFrom, tc.haveTo, 0, 0)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantGroups, gotGroups)

		})
	}
}

func TestLogic_GetGroupByID(t *testing.T) {
	var (
		haveErr   = errors.New("db error")
		timeNow   = time.Now().UTC()
		haveGroup = models.GroupWithDecks{
			Group: models.Group{
				ID:         uuid.NewString(),
				Name:       uuid.NewString(),
				CreatedBy:  uuid.NewString(),
				Moderators: []string{},
				DeckIDs:    []string{},
				CreatedAt:  timeNow,
				UpdatedAt:  timeNow,
				DeletedAt:  nil,
			},
			Decks: []models.GetDeckResults{
				{
					ID:        uuid.NewString(),
					Name:      uuid.NewString(),
					Upvotes:   1,
					Downvotes: 3,
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
			},
		}
	)

	testCases := map[string]struct {
		mockStore   func(mock *database.MockRepository)
		haveGroupID string
		wantGroup   models.GroupWithDecks
		wantErr     error
	}{
		"should return group when database returns group": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).Return(haveGroup, nil)
			},
			haveGroupID: uuid.NewString(),
			wantGroup:   haveGroup,
		},
		"should return the error the db returns": {
			mockStore: func(mock *database.MockRepository) {
				mock.EXPECT().GetGroupByID(gomock.Any(), gomock.Any()).Return(models.GroupWithDecks{}, haveErr)
			},
			haveGroupID: uuid.NewString(),
			wantGroup:   models.GroupWithDecks{},
			wantErr:     haveErr,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDB := database.NewMockRepository(ctrl)

			if tc.mockStore != nil {
				tc.mockStore(mockDB)
			}

			logic := Logic{repo: mockDB, logger: zerolog.Nop()}

			gotGroup, err := logic.GetGroupByID(context.Background(), tc.haveGroupID)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantGroup, gotGroup)
		})
	}
}
