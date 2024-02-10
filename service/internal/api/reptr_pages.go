package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rmarken/reptr/api"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rmarken/reptr/service/internal/web/components/dumb"
	"github.com/rmarken/reptr/service/internal/web/components/pages"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const hxTriggerHeaderKey = "HX-Trigger"

var tailwindArr = []string{"/styles/pages/tailwind-output.css"}

func (rc ReprtClient) ServeStyles(w http.ResponseWriter, r *http.Request, path string, styleName string) {
	log := rc.logger.With().Str("method", "ServeStyles").Logger()
	log.Info().Msgf("serving  %s %s", path, styleName)

	absolutePath, err := filepath.Abs(fmt.Sprintf("./service/internal/web/styles/%s/%s", path, styleName))
	file, err := os.ReadFile(absolutePath)
	if err != nil {
		log.Error().Err(err).Msgf("while reading: %s", absolutePath)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/css")
	w.Write(file)
}

func (rc ReprtClient) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "RegistrationPage").Logger()
	log.Info().Msgf("serving registration page")
	err := pages.Page(pages.Register(nil), tailwindArr).Render(r.Context(), w)
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
		err := pages.Register(pages.Banner("Must provide email")).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := pages.Page(pages.Register(pages.Banner("Must provide password")), tailwindArr).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	//TODO: check password strength

	if password != repassword {
		w.WriteHeader(http.StatusBadRequest)
		err := pages.Register(pages.Banner("Passwords do not match")).Render(r.Context(), w)
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
		err := pages.Page(pages.Register(pages.Banner(registrationError.Description)), tailwindArr).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	log.Info().Msgf("user is registered: %+v", user)
	w.WriteHeader(http.StatusCreated)
	err = pages.Page(pages.Form(pages.Banner("Registration Successful"), pages.Login()), tailwindArr).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rc ReprtClient) LoginPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "LoginPage").Logger()
	log.Info().Msgf("serving login page")
	err := pages.Page(pages.Form(nil, pages.Login()), tailwindArr).Render(r.Context(), w)
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

	homeGroups := make([]pages.Group, len(groups))
	for i, group := range groups {
		homeGroups[i] = homeGroupFromModel(group)
	}

	pages.Page(pages.Home(pages.HomeData{Username: userName, Groups: homeGroups}), tailwindArr).Render(r.Context(), w)
}

func (rc ReprtClient) CreateGroup(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateGroup").Logger()
	logger.Debug().Msg("create group called")

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Info().Msg("username is not on context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	groupName := r.PostForm.Get("group-name")
	if groupName == "" {
		logger.Info().Msgf("create group attempt without groupName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = rc.deckController.CreateGroup(r.Context(), username, groupName)
	if err != nil {
		logger.Error().Err(err).Msgf("while calling CreateGroup")
		w.WriteHeader(toStatus(err))
	}

	pages.Page(pages.Form(pages.Banner("Group Successfully Created"), pages.CreateGroupForm()), tailwindArr).Render(r.Context(), w)
}

func (rc ReprtClient) GroupPage(w http.ResponseWriter, r *http.Request, groupID string) {
	logger := rc.logger.With().Str("method", "GroupPage").Logger()
	logger.Info().Msgf("serving group page for: %s", groupID)

	group, err := rc.deckController.GetGroupByID(r.Context(), groupID)
	if err != nil {
		http.Error(w, "while getting groups for user", toStatus(err))
		return
	}
	pages.Page(pages.Form(nil, pages.GroupPage(groupPageFromModel(group))), tailwindArr).Render(r.Context(), w)
}

func homeGroupFromModel(group models.Group) pages.Group {
	return pages.Group{
		ID:        group.ID,
		GroupName: group.Name,
		NumDecks:  len(group.DeckIDs),
		NumUsers:  0,
	}
}

func groupPageFromModel(group models.GroupWithDecks) pages.GroupData {
	return pages.GroupData{
		ID:        group.ID,
		GroupName: group.Name,
		Decks:     groupDecksFromDecks(group.Decks),
		NumUsers:  "0",
	}
}

func groupDecksFromDecks(fromService []models.GetDeckResults) []pages.Deck {
	apiDecks := make([]pages.Deck, len(fromService))
	for i, deck := range fromService {
		apiDecks[i] = pages.Deck{
			ID:           deck.ID,
			DeckName:     deck.Name,
			NumUpvotes:   deck.Upvotes,
			NumDownvotes: deck.Downvotes,
			NumCards:     0,
			CreatedAt:    deck.CreatedAt,
			UpdatedAt:    deck.UpdatedAt,
		}
	}
	return apiDecks
}

func (rc ReprtClient) CreateGroupPage(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateGroupPage").Logger()
	logger.Info().Msg("serving create group page")
	pages.Page(pages.Form(nil, pages.CreateGroupForm()), tailwindArr).Render(r.Context(), w)
}

func (rc ReprtClient) CreateDeckPage(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateDeckPage").Logger()
	logger.Info().Msg("serving create deck page")
	pages.Page(pages.Form(nil, pages.CreateDeckPage(pages.CreateDeckPageData{
		UsersGroups: nil,
	})), tailwindArr).Render(r.Context(), w)

}

func (rc ReprtClient) CreateDeck(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateDeck").Logger()
	logger.Info().Msg("serving creating deck")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deckName := r.PostForm.Get("deck-name")
	if deckName == "" {
		logger.Info().Msgf("create deck attempt without deckName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deckID, err := rc.deckController.CreateDeck(r.Context(), deckName)
	if err != nil {
		http.Error(w, "while creating deck", toStatus(err))
	}

	http.Redirect(w, r, fmt.Sprintf("/page/create-cards/%s", deckID), http.StatusSeeOther)
}

func (rc ReprtClient) GetCreateCardsForDeckPage(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "GetCreateCardsForDeckPage").Logger()
	logger.Info().Msg("serving creating deck")

	if deckID == "" {
		logger.Info().Msgf("get cards without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deck, err := rc.deckController.GetCardsByDeckID(r.Context(), deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while cards for deck %s", deckID)
		http.Error(w, "while cards for deck", toStatus(err))
		return
	}
	viewCards := make([]dumb.CardDisplay, len(deck.Cards))
	for i, card := range deck.Cards {
		viewCards[i] = dumb.CardDisplay{
			Front: card.Front,
			Back:  card.Back,
		}
	}
	pages.Page(pages.Form(nil, pages.DeckCreateCardForm(pages.DeckCreateCardData{
		DeckID:   deck.ID,
		DeckName: deck.Name,
		Cards:    viewCards,
	})), tailwindArr).Render(r.Context(), w)

}
func (rc ReprtClient) GetCardsForDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "GetCardsForDeck").Logger()
	logger.Info().Msgf("getting cards for deck: %s", deckID)

	if deckID == "" {
		logger.Info().Msgf("get cards without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deck, err := rc.deckController.GetCardsByDeckID(r.Context(), deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while cards for deck %s", deckID)
		http.Error(w, "while cards for deck", toStatus(err))
		return
	}
	viewCards := make([]dumb.CardDisplay, len(deck.Cards))
	for i, card := range deck.Cards {
		viewCards[i] = dumb.CardDisplay{
			Front: card.Front,
			Back:  card.Back,
		}
	}
	dumb.GroupCardDisplay(dumb.GroupCardDisplayPageData{Cards: viewCards}).Render(r.Context(), w)
}
func (rc ReprtClient) ViewDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "ViewDeck").Logger()
	logger.Info().Msg("creating card")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deckID == "" {
		logger.Info().Msgf("create deck attempt without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := rc.getCardViewerContent(r.Context(), username, deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting card content %s", deckID)
		http.Error(w, "while cards for deck", toStatus(err))
		return
	}

	pages.Page(pages.DeckViewerPage(content), tailwindArr).Render(r.Context(), w)
}

func (rc ReprtClient) getCardViewerContent(ctx context.Context, username, deckID string) (pages.DeckViewPageData, error) {
	s, err := rc.sessionController.GetSessionForUserAndDeckID(ctx, username, deckID)
	if err != nil {
		return pages.DeckViewPageData{}, err
	}
	if s.IsFront {
		f, err := rc.deckController.GetFrontOfCardByID(ctx, deckID, s.CurrentCardID)
		if err != nil {
			return pages.DeckViewPageData{}, err
		}
		return pages.DeckViewPageData{
			DeckName: s.DeckName,
			DeckID:   deckID,
			Content: dumb.FrontCardDisplay(dumb.CardFront{
				CardType:       "",
				CardID:         s.CurrentCardID,
				Front:          f.Content,
				Upvotes:        strconv.Itoa(f.Upvotes),
				Downvotes:      strconv.Itoa(f.Downvotes),
				PreviousCardID: f.PreviousCard,
				NextCardID:     f.NextCard,
			}),
		}, err
	}

	b, err := rc.deckController.GetBackOfCardByID(ctx, deckID, s.CurrentCardID)
	if err != nil {
		return pages.DeckViewPageData{}, err
	}
	return pages.DeckViewPageData{
		DeckName: s.DeckName,
		DeckID:   deckID,
		Content: dumb.BackOfCardDisplay(dumb.CardBack{
			CardID:      s.CurrentCardID,
			BackContent: b.Answer,
			NextCardID:  b.NextCard,
			IsUpvoted:   b.IsUpvotedByUser,
			IsDownvoted: b.IsDownvotedByUser,
		}),
	}, err

}
func (rc ReprtClient) CreateCardForDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "CreateCardForDeck").Logger()
	logger.Info().Msg("creating card")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if deckID == "" {
		logger.Info().Msgf("create deck attempt without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cardFront := r.PostForm.Get("card-front")
	if cardFront == "" {
		logger.Info().Msgf("create deck attempt without card front")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cardBack := r.PostForm.Get("card-back")
	if cardBack == "" {
		logger.Info().Msgf("create deck attempt without card back")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timeNow := time.Now().UTC()
	err = rc.deckController.AddCardToDeck(r.Context(), deckID, models.Card{
		ID:        uuid.NewString(),
		Front:     cardFront,
		Back:      cardBack,
		Kind:      models.BasicCard,
		DeckID:    deckID,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("while creating a card for deck %s", deckID)
		http.Error(w, "while creating card", toStatus(err))
		return
	}

	w.Header().Set(hxTriggerHeaderKey, "newCard")
	w.WriteHeader(http.StatusCreated)
}

func (rc ReprtClient) BackOfCard(w http.ResponseWriter, r *http.Request, deckID, cardID string) {
	logger := rc.logger.With().Str("method", "BackOfCard").Logger()
	logger.Info().Msgf("getting back of card with ID: %s", cardID)

	if deckID == "" {
		logger.Info().Msgf("get back of card without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cardID == "" {
		logger.Info().Msgf("get back of card without cardID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	backOfCard, err := rc.deckController.GetBackOfCardByID(r.Context(), deckID, cardID)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting back of card for cardID: %s", cardID)
		http.Error(w, "while getting back of card", toStatus(err))
		return
	}

	dumb.BackOfCardDisplay(dumb.CardBack{
		CardID:      backOfCard.CardID,
		BackContent: backOfCard.Answer,
		NextCardID:  backOfCard.NextCard,
		IsUpvoted:   backOfCard.IsUpvotedByUser,
		IsDownvoted: backOfCard.IsDownvotedByUser,
	}).Render(r.Context(), w)
}

func (rc ReprtClient) FrontOfCard(w http.ResponseWriter, r *http.Request, deckID, cardID string) {
	logger := rc.logger.With().Str("method", "FrontOfCard").Logger()
	logger.Info().Msgf("getting front of card with ID: %s", cardID)

	if deckID == "" {
		logger.Info().Msgf("get front of card without deckID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cardID == "" {
		logger.Info().Msgf("get front of card without cardID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	frontOfCard, err := rc.deckController.GetFrontOfCardByID(r.Context(), deckID, cardID)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting front of card for cardID: %s", cardID)
		http.Error(w, "while getting back of card", toStatus(err))
		return
	}

	dumb.FrontCardDisplay(dumb.CardFront{
		CardID:     frontOfCard.CardID,
		Front:      frontOfCard.Content,
		NextCardID: frontOfCard.NextCard,
		Downvotes:  strconv.Itoa(frontOfCard.Downvotes),
		Upvotes:    strconv.Itoa(frontOfCard.Upvotes),
		CardType:   "",
	}).Render(r.Context(), w)
}

func toStatus(err error) int {
	switch {
	case errors.Is(err, decks.ErrInvalidToBeforeFrom),
		errors.Is(err, decks.ErrInvalidGroupName),
		errors.Is(err, decks.ErrInvalidDeckName),
		errors.Is(err, decks.ErrEmptyGroupID),
		errors.Is(err, decks.ErrEmptyDeckID):
		return http.StatusBadRequest
	case errors.Is(err, database.ErrNoResults):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
