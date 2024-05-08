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

const (
	hxTriggerHeaderKey = "HX-Trigger"

	stylesDir = "/styles/pages/"

	//styles
	//base
	baseStyle = stylesDir + "base.css"
	pageStyle = stylesDir + "page.css"
	// component
	formStyle  = stylesDir + "form.css"
	tableStyle = stylesDir + "table.css"

	//page level
	loginStyle        = stylesDir + "login.css"
	registrationStyle = stylesDir + "registration.css"
	homeStyle         = stylesDir + "home.css"
	groupStyle        = stylesDir + "group.css"
	deckViewStyle     = stylesDir + "deck_viewer.css"
	createDeckStyle   = stylesDir + "create_deck.css"
	errorStyle        = stylesDir + "error.css"
)

var cssFileArr = []string{baseStyle, pageStyle}

func (rc ReprtClient) ServeStyles(w http.ResponseWriter, r *http.Request, path string, styleName string) {
	log := rc.logger.With().Str("method", "ServeStyles").Logger()
	log.Info().Msgf("serving %s %s", path, styleName)

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
	err := pages.Page(pages.Register(nil), append(cssFileArr, formStyle, registrationStyle)).Render(r.Context(), w)
	if err != nil {
		log.Error().Err(err).Msg("while trying to serve registration page")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to serve registration page",
			Msg:        "Something went wrong while serving registration page",
		})
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
		err := pages.Page(pages.Register(pages.Banner("Must provide password")), cssFileArr).Render(r.Context(), w)
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
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to registering user",
			Msg:        "Something went wrong while registering user",
		})
		return
	}

	if !registrationError.IsZero() {
		w.WriteHeader(registrationError.StatusCode)
		err := pages.Page(pages.Register(pages.Banner(registrationError.Description)), cssFileArr).Render(r.Context(), w)
		if err != nil {
			rc.serveError(w, r, pages.ErrorPageData{
				StatusCode: strconv.Itoa(http.StatusInternalServerError),
				Status:     http.StatusText(http.StatusInternalServerError),
				Error:      "unable to registering user",
				Msg:        "Something went wrong while registering user",
			})
		}
		return
	}

	log.Info().Msgf("user is registered: %+v", user)
	w.WriteHeader(http.StatusCreated)
	err = pages.Page(pages.Form(pages.Banner("Registration Successful"), pages.Login()), cssFileArr).Render(r.Context(), w)
	if err != nil {
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to serve registration page",
			Msg:        "Something went wrong while serving registration page",
		})
		return
	}
}

func (rc ReprtClient) LoginPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "LoginPage").Logger()
	log.Info().Msgf("serving login page")
	err := pages.Page(pages.Form(nil, pages.Login()), append(cssFileArr, loginStyle, formStyle)).Render(r.Context(), w)
	if err != nil {
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to serve login page",
			Msg:        "Something went wrong while serving to login page",
		})
		return
	}
}

func (rc ReprtClient) Login(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "Login").Logger()

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to parse form",
			Msg:        "Something went wrong while attempting to login",
		})

		return
	}

	email := r.PostForm.Get("email")
	if email == "" {
		logger.Info().Msgf("login attempt without email")
		w.WriteHeader(http.StatusBadRequest)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "login attempt without email",
			Msg:        "login attempt without email",
		})

		return
	}
	password := r.PostForm.Get("password")
	if password == "" {
		logger.Info().Msgf("login attempt without password - email: %s", email)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "login attempt without password - email",
			Msg:        "login attempt without password - emaild",
		})
		return
	}

	token, err := rc.authenticator.PasswordCredentialsToken(r.Context(), email, password)
	if err != nil {
		logger.Error().Err(err).Msgf("error authenticating %s: %v", email, err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusUnauthorized),
			Status:     http.StatusText(http.StatusUnauthorized),
			Error:      "Bad request - Invalid username or password",
			Msg:        "Bad request - Invalid username or password",
		})
		return
	}

	if tokenString, ok := token.Extra(IDToken).(string); ok {
		_, err := rc.authenticator.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			rc.serveError(w, r, pages.ErrorPageData{
				StatusCode: strconv.Itoa(http.StatusInternalServerError),
				Status:     http.StatusText(http.StatusInternalServerError),
				Error:      "unable to verify token",
				Msg:        "Something went wrong while logging in",
			})
			return
		}
		session, err := rc.store.Get(r, CookieSessionID)
		if err != nil {
			rc.serveError(w, r, pages.ErrorPageData{
				StatusCode: strconv.Itoa(http.StatusInternalServerError),
				Status:     http.StatusText(http.StatusInternalServerError),
				Error:      "unable to get session",
				Msg:        "Something went wrong while logging in",
			})
			return
		}
		session.Values[SessionTokenKey] = "Bearer " + tokenString
		session.Options.Secure = true
		err = session.Save(r, w)
		if err != nil {
			logger.Error().Err(err).Msgf("unable to save session")
			rc.serveError(w, r, pages.ErrorPageData{
				StatusCode: strconv.Itoa(http.StatusInternalServerError),
				Status:     http.StatusText(http.StatusInternalServerError),
				Error:      "unable to save session",
				Msg:        "Something went wrong while logging in",
			})
			return
		}
		w.Header().Set("Authorization", "Bearer "+tokenString)
		http.Redirect(w, r, "/page/home", http.StatusSeeOther)
		return
	}

	logger.Error().Msgf("Bad request - token doesn't contain %s", IDToken)
	rc.serveError(w, r, pages.ErrorPageData{
		StatusCode: strconv.Itoa(http.StatusInternalServerError),
		Status:     http.StatusText(http.StatusInternalServerError),
		Error:      fmt.Sprintf("Bad request - token doesn't contain %s", IDToken),
		Msg:        "Something went wrong while logging in",
	})

	return
}

func (rc ReprtClient) GetDecksForUser(w http.ResponseWriter, r *http.Request, params api.GetDecksForUserParams) {
	_ = rc.logger.With().Str("method", "GetDecksForUser")
	panic("implement me")
}

func (rc ReprtClient) HomePage(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "HomePage").Logger()

	var userName string
	var ok bool
	if userName, ok = reptrCtx.Username(r.Context()); !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     "Internal Server Error",
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	homepageData, err := rc.deckController.GetHomepageData(r.Context(), userName, time.Time{}, nil, 10, 0)
	if err != nil {
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting homepageData for use",
			Msg:        "Unable to load homepage.",
		})
		return
	}

	homeGroups := make([]pages.HomeGroupData, len(homepageData.Groups))
	for i, group := range homepageData.Groups {
		homeGroups[i] = homeGroupFromModel(group)
	}

	homeDecks := make([]dumb.Deck, len(homepageData.Decks))
	for i, deck := range homepageData.Decks {
		homeDecks[i] = webDeckFromModel(deck)
	}
	pages.Page(pages.Home(pages.HomeData{Username: userName, Groups: homeGroups, Decks: homeDecks}), append(cssFileArr, tableStyle, homeStyle, groupStyle)).Render(r.Context(), w)
}

func (rc ReprtClient) CreateGroup(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateGroup").Logger()

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     "Internal Server Error",
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to parse form",
			Msg:        "unable to create group",
		})
		return
	}

	groupName := r.PostForm.Get("group-name")
	if groupName == "" {
		logger.Error().Msgf("create group attempt without groupName")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "create group attempt without groupName",
			Msg:        "unable to create group",
		})
		return
	}

	_, err = rc.deckController.CreateGroup(r.Context(), username, groupName)
	if err != nil {
		logger.Error().Err(err).Msgf("while calling CreateGroup")
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while calling CreateGroup",
			Msg:        "unable to create group",
		})
		return
	}

	pages.Page(pages.Form(pages.Banner("Group Successfully Created"), pages.CreateGroupForm()), append(cssFileArr, formStyle, groupStyle)).Render(r.Context(), w)
}

func (rc ReprtClient) GroupPage(w http.ResponseWriter, r *http.Request, groupID string) {
	logger := rc.logger.With().Str("method", "GroupPage").Logger()
	logger.Info().Msgf("serving group page for: %s", groupID)

	group, err := rc.deckController.GetGroupByID(r.Context(), groupID)
	if err != nil {
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "unable to parse form",
			Msg:        "Problem with creating deck.",
		})
		return
	}
	pages.Page(pages.Form(nil, pages.GroupPage(groupPageFromModel(group))), append(cssFileArr, tableStyle, groupStyle)).Render(r.Context(), w)
}

func homeGroupFromModel(group models.Group) pages.HomeGroupData {
	return pages.HomeGroupData{
		ID:        group.ID,
		GroupName: group.Name,
		NumDecks:  len(group.DeckIDs),
		NumUsers:  0,
	}
}

func webDeckFromModel(deck models.GetDeckResults) dumb.Deck {
	return dumb.Deck{
		ID:           deck.ID,
		DeckName:     deck.Name,
		NumCards:     deck.NumCards,
		NumUpvotes:   deck.Upvotes,
		NumDownvotes: deck.Downvotes,
		CreatedAt:    deck.CreatedAt,
		UpdatedAt:    deck.UpdatedAt,
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

func groupDecksFromDecks(fromService []models.GetDeckResults) []dumb.Deck {
	apiDecks := make([]dumb.Deck, len(fromService))
	for i, deck := range fromService {
		apiDecks[i] = dumb.Deck{
			ID:           deck.ID,
			DeckName:     deck.Name,
			NumUpvotes:   deck.Upvotes,
			NumDownvotes: deck.Downvotes,
			NumCards:     deck.NumCards,
			CreatedAt:    deck.CreatedAt,
			UpdatedAt:    deck.UpdatedAt,
		}
	}
	return apiDecks
}

func (rc ReprtClient) CreateGroupPage(w http.ResponseWriter, r *http.Request) {
	logger := rc.logger.With().Str("method", "CreateGroupPage").Logger()
	logger.Info().Msg("serving create group page")
	pages.Page(pages.Form(nil, pages.CreateGroupForm()), append(cssFileArr, formStyle)).Render(r.Context(), w)
}

func (rc ReprtClient) CreateDeckPage(w http.ResponseWriter, r *http.Request, groupID string) {
	logger := rc.logger.With().Str("method", "CreateDeckPage").Logger()
	logger.Info().Msg("serving create deck page")
	path := "/page/create-deck"
	if groupID != "" {
		path = path + "/" + groupID
	}

	pages.Page(pages.Form(nil, pages.CreateDeckPage(path)), append(cssFileArr, formStyle, createDeckStyle)).Render(r.Context(), w)
}

func (rc ReprtClient) CreateDeck(w http.ResponseWriter, r *http.Request, groupID string) {
	logger := rc.logger.With().Str("method", "CreateDeck").Logger()
	logger.Info().Msg("serving creating deck")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		w.WriteHeader(http.StatusInternalServerError)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to parse form",
			Msg:        "Problem with creating deck.",
		})
		return
	}

	deckName := r.PostForm.Get("deck-name")
	if deckName == "" {
		logger.Error().Msgf("create deck attempt without deckName")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "create deck attempt without deckName",
			Msg:        "Problem with creating deck.",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	deckID, err := rc.deckController.CreateDeck(r.Context(), deckName, username)
	if err != nil {
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      err.Error(),
			Msg:        "Problem with creating deck.",
		})
		return
	}
	if groupID != "" {
		// TODO: bundle these in deck controller so that they can be done in a tx.
		err = rc.deckController.AddDeckToGroup(r.Context(), groupID, deckID)
		if err != nil {
			status := toStatus(err)
			rc.serveError(w, r, pages.ErrorPageData{
				StatusCode: strconv.Itoa(status),
				Status:     http.StatusText(status),
				Error:      err.Error(),
				Msg:        "Problem with creating deck.",
			})
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/page/create-cards-content/%s", deckID), http.StatusSeeOther)
}

func (rc ReprtClient) GetCreateCardsForDeckPage(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "GetCreateCardsForDeckPage").Logger()
	logger.Info().Msg("serving create cards for deck")

	if deckID == "" {
		logger.Error().Msgf("get cards without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get cards without deckID",
			Msg:        "Problem getting cards.",
		})
		return
	}

	deck, err := rc.deckController.GetCardsByDeckID(r.Context(), deckID)
	if err != nil && !errors.Is(err, database.ErrNoResults) {
		logger.Error().Err(err).Msgf("while cards for deck %s", deckID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting cards for deck",
			Msg:        "Problem getting cards.",
		})
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
	})), append(cssFileArr, formStyle, createDeckStyle)).Render(r.Context(), w)

}

func (rc ReprtClient) GetCreateCardsForDeckContent(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "GetCreateCardsForDeckContent").Logger()
	logger.Info().Msg("serving create cards content for deck")

	if deckID == "" {
		logger.Error().Msgf("get cards without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get cards without deckID",
			Msg:        "Problem getting cards.",
		})
		return
	}

	deck, err := rc.deckController.GetCardsByDeckID(r.Context(), deckID)
	if err != nil && !errors.Is(err, database.ErrNoResults) {
		logger.Error().Err(err).Msgf("while cards for deck %s", deckID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting cards for deck",
			Msg:        "Problem getting cards.",
		})
		return
	}

	viewCards := make([]dumb.CardDisplay, len(deck.Cards))
	for i, card := range deck.Cards {
		viewCards[i] = dumb.CardDisplay{
			Front: card.Front,
			Back:  card.Back,
		}
	}

	pages.CreateDeckContent(pages.DeckCreateCardData{
		DeckID:   deck.ID,
		DeckName: deck.Name,
		Cards:    viewCards,
	}).Render(r.Context(), w)

}

func (rc ReprtClient) GetCardsForDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "GetCardsForDeck").Logger()
	logger.Error().Msgf("getting cards for deck: %s", deckID)

	if deckID == "" {
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get cards without deckID",
			Msg:        "Problem getting cards.",
		})
		return
	}

	deck, err := rc.deckController.GetCardsByDeckID(r.Context(), deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while cards for deck %s", deckID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while cards for deck",
			Msg:        "Problem getting cards.",
		})
		return
	}
	viewCards := make([]dumb.CardDisplay, len(deck.Cards))
	for i, card := range deck.Cards {
		viewCards[i] = dumb.CardDisplay{
			Front: card.Front,
			Back:  card.Back,
		}
	}
	dumb.GroupCardDisplay(viewCards).Render(r.Context(), w)
}
func (rc ReprtClient) ViewDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "ViewDeck").Logger()
	logger.Info().Msg("view deck")

	if deckID == "" {
		logger.Info().Msgf("view deck attempt without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "view deck attempt without deckID",
			Msg:        "Problem getting deck content.",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	content, err := rc.getCardViewerContent(r.Context(), username, deckID)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting card content %s", deckID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting card content",
			Msg:        "Problem getting deck content.",
		})
		return
	}

	pages.Page(pages.DeckViewerPage(content), append(cssFileArr, deckViewStyle)).Render(r.Context(), w)
}

func (rc ReprtClient) getCardViewerContent(ctx context.Context, username, deckID string) (pages.DeckViewPageData, error) {
	s, err := rc.sessionController.GetSessionForUserAndDeckID(ctx, username, deckID)
	if err != nil {
		return pages.DeckViewPageData{}, err
	}
	if s.IsFront {
		f, err := rc.deckController.GetFrontOfCardByID(ctx, deckID, s.CurrentCardID, username)
		if err != nil {
			return pages.DeckViewPageData{}, err
		}
		return pages.DeckViewPageData{
			DeckName: s.DeckName,
			DeckID:   deckID,
			Content: dumb.FrontCardDisplay(dumb.CardFront{
				DeckID:         deckID,
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

	b, err := rc.deckController.GetBackOfCardByID(ctx, deckID, s.CurrentCardID, username)
	if err != nil {
		return pages.DeckViewPageData{}, err
	}
	return pages.DeckViewPageData{
		DeckName: s.DeckName,
		DeckID:   deckID,
		Content: dumb.BackOfCardDisplay(dumb.CardBack{
			DeckID:         deckID,
			CardID:         s.CurrentCardID,
			BackContent:    b.Answer,
			NextCardID:     b.NextCard,
			PreviousCardID: b.PreviousCard,
			IsUpvoted:      bool(b.IsUpvotedByUser),
			IsDownvoted:    bool(b.IsDownvotedByUser),
		}),
	}, err

}
func (rc ReprtClient) CreateCardForDeck(w http.ResponseWriter, r *http.Request, deckID string) {
	logger := rc.logger.With().Str("method", "CreateCardForDeck").Logger()
	logger.Info().Msg("creating card")

	err := r.ParseForm()
	if err != nil {
		logger.Error().Err(err).Msg("unable to parse form")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "unable to parse form",
			Msg:        "Problem creating card",
		})
		return
	}

	if deckID == "" {
		logger.Error().Msgf("create deck attempt without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "create deck attempt without deckID",
			Msg:        "Problem creating card",
		})
		return
	}

	cardFront := r.PostForm.Get("card-front")
	if cardFront == "" {
		logger.Info().Msgf("create deck attempt without card front")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "create deck attempt without card front",
			Msg:        "Problem creating card",
		})
		return
	}

	cardBack := r.PostForm.Get("card-back")
	if cardBack == "" {
		logger.Error().Msgf("create deck attempt without card back")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "create deck attempt without card back",
			Msg:        "Problem creating card",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     "Internal Server Error",
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
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
		CreatedBy: username,
	})
	if err != nil {
		logger.Error().Err(err).Msgf("while creating a card for deck %s", deckID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while creating a card for deck",
			Msg:        "Problem creating card",
		})
		return
	}

	w.Header().Set(hxTriggerHeaderKey, "newCard")
	w.WriteHeader(http.StatusCreated)
}

func (rc ReprtClient) BackOfCard(w http.ResponseWriter, r *http.Request, deckID, cardID string) {
	logger := rc.logger.With().Str("method", "BackOfCard").Logger()
	logger.Info().Msgf("getting back of card with ID: %s", cardID)

	if deckID == "" {
		logger.Error().Msgf("get back of card without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get back of card without deckID",
			Msg:        "Problem processing getting back of card.",
		})
		return
	}

	if cardID == "" {
		logger.Error().Msgf("get back of card without cardID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get back of card without cardID",
			Msg:        "Problem processing getting back of card.",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     http.StatusText(http.StatusInternalServerError),
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	backOfCard, err := rc.deckController.GetBackOfCardByID(r.Context(), deckID, cardID, username)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting back of card for cardID: %s", cardID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting back of card",
			Msg:        "Problem processing getting back of card.",
		})
		return
	}

	dumb.BackOfCardDisplay(dumb.CardBack{
		DeckID:         deckID,
		CardID:         backOfCard.CardID,
		BackContent:    backOfCard.Answer,
		NextCardID:     backOfCard.NextCard,
		PreviousCardID: backOfCard.PreviousCard,
		IsUpvoted:      bool(backOfCard.IsUpvotedByUser),
		IsDownvoted:    bool(backOfCard.IsDownvotedByUser),
		VoteButtonData: dumb.VoteButtonsData{
			CardID:            backOfCard.CardID,
			UpvoteClass:       backOfCard.IsUpvotedByUser.UpvotedClass(),
			DownvoteClass:     backOfCard.IsDownvotedByUser.DownvotedClass(),
			UpvoteDirection:   backOfCard.IsUpvotedByUser.NextUpvoteDirection(),
			DownvoteDirection: backOfCard.IsDownvotedByUser.DownvotedClass()},
	}).Render(r.Context(), w)
}

func (rc ReprtClient) FrontOfCard(w http.ResponseWriter, r *http.Request, deckID, cardID string) {
	logger := rc.logger.With().Str("method", "FrontOfCard").Logger()
	logger.Info().Msgf("getting front of card with ID: %s", cardID)

	if deckID == "" {
		logger.Error().Msgf("get front of card without deckID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "get front of card without deckID",
			Msg:        "Problem processing getting front of card.",
		})
		return
	}

	if cardID == "" {
		logger.Error().Msgf("get front of card without cardID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     http.StatusText(http.StatusBadRequest),
			Error:      "getting front of card without cardID",
			Msg:        "Problem processing getting front of card.",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     "Internal Server Error",
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	frontOfCard, err := rc.deckController.GetFrontOfCardByID(r.Context(), deckID, cardID, username)
	if err != nil {
		logger.Error().Err(err).Msgf("while getting front of card for cardID: %s", cardID)
		status := toStatus(err)
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(status),
			Status:     http.StatusText(status),
			Error:      "while getting front of card",
			Msg:        "Something went wrong while getting front of card.",
		})
		http.Error(w, "while getting front of card", toStatus(err))
		return
	}

	dumb.FrontCardDisplay(dumb.CardFront{
		DeckID:         deckID,
		CardID:         frontOfCard.CardID,
		Front:          frontOfCard.Content,
		NextCardID:     frontOfCard.NextCard,
		PreviousCardID: frontOfCard.PreviousCard,
		Downvotes:      strconv.Itoa(frontOfCard.Downvotes),
		Upvotes:        strconv.Itoa(frontOfCard.Upvotes),
		CardType:       "",
	}).Render(r.Context(), w)
}

func (rc ReprtClient) VoteCard(w http.ResponseWriter, r *http.Request, cardID string, direction string) {
	logger := rc.logger.With().Str("method", "VoteCard").Logger()
	logger.Info().Msg("voting card")

	if cardID == "" {
		logger.Info().Msgf("card vote without cardID")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     "Bad Request",
			Error:      "card vote without cardIn",
			Msg:        "Something went wrong while processing vote.",
		})
		return
	}

	if direction == "" {
		logger.Error().Msgf("card vote without direction")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     "Bad Request",
			Error:      "card vote without direction",
			Msg:        "Something went wrong while processing vote.",
		})
		return
	}

	vote := models.VoteFromString(direction)

	if vote == models.Unknown {
		logger.Error().Msgf("unknown vote type")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusBadRequest),
			Status:     "Bad Request",
			Error:      "unknown vote type",
			Msg:        "Something went wrong while processing vote.",
		})
		return
	}

	username, ok := reptrCtx.Username(r.Context())
	if !ok {
		logger.Error().Msgf("username is not on context")
		rc.serveError(w, r, pages.ErrorPageData{
			StatusCode: strconv.Itoa(http.StatusInternalServerError),
			Status:     "Internal Server Error",
			Error:      "username not on context",
			Msg:        "Try logging back in.",
		})
		return
	}

	err := rc.deckController.VoteCard(r.Context(), vote, cardID, username)
	if err != nil {
		logger.Error().Err(err).Msgf("voting for card from user with vote: %s %s %s", cardID, username, vote.String())
		http.Error(w, "voting for card", toStatus(err))
	}

	dumb.VoteButtons(dumb.VoteButtonsData{
		CardID:            cardID,
		UpvoteClass:       vote.UpvoteClass(),
		DownvoteClass:     vote.DownvoteClass(),
		UpvoteDirection:   vote.NextUpvote().String(),
		DownvoteDirection: vote.NextDownvote().String(),
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

func (rc ReprtClient) serveError(w http.ResponseWriter, r *http.Request, data pages.ErrorPageData) {
	code, err := strconv.Atoi(data.StatusCode)
	if err != nil {
		rc.logger.Error().Err(err).Msgf("not able to convert data.StatusCode to int: %s", data.StatusCode)
	}
	w.WriteHeader(code)
	pages.Page(pages.Error(data), append(cssFileArr, errorStyle)).Render(r.Context(), w)
}
