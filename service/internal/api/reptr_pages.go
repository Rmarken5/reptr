package api

import (
	"errors"
	"fmt"
	"github.com/rmarken/reptr/api"
	reptrCtx "github.com/rmarken/reptr/service/internal/context"
	"github.com/rmarken/reptr/service/internal/database"
	"github.com/rmarken/reptr/service/internal/logic/decks"
	"github.com/rmarken/reptr/service/internal/models"
	"github.com/rmarken/reptr/service/internal/web/components/pages"
	"net/http"
	"time"
)

func (rc ReprtClient) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "RegistrationPage").Logger()
	log.Info().Msgf("serving registration page")
	err := pages.Page(pages.Register(nil)).Render(r.Context(), w)
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
		err := pages.Page(pages.Register(pages.Banner("Must provide password"))).Render(r.Context(), w)
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
		err := pages.Page(pages.Register(pages.Banner(registrationError.Description))).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	log.Info().Msgf("user is registered: %+v", user)
	w.WriteHeader(http.StatusCreated)
	err = pages.Login(pages.Banner("Registration Successful")).Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rc ReprtClient) LoginPage(w http.ResponseWriter, r *http.Request) {
	log := rc.logger.With().Str("method", "LoginPage").Logger()
	log.Info().Msgf("serving login page")
	err := pages.Page(pages.Login(nil)).Render(r.Context(), w)
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

	pages.Page(pages.Home(pages.HomeData{Username: userName, Groups: homeGroups})).Render(r.Context(), w)
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

	pages.Page(pages.Form(pages.Banner("Group Successfully Created"), pages.CreateGroupForm())).Render(r.Context(), w)
}

func (rc ReprtClient) GroupPage(w http.ResponseWriter, r *http.Request, groupID string) {
	logger := rc.logger.With().Str("method", "GroupPage").Logger()
	logger.Info().Msgf("serving group page for: %s", groupID)

	group, err := rc.deckController.GetGroupByID(r.Context(), groupID)
	if err != nil {
		http.Error(w, "while getting groups for user", toStatus(err))
		return
	}

	pages.Form(nil, pages.Page(pages.GroupPage(groupPageFromModel(group)))).Render(r.Context(), w)
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
	pages.Form(nil, pages.Page(pages.CreateGroupForm())).Render(r.Context(), w)
}

func (rc ReprtClient) CreateDeckPage(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (rc ReprtClient) CreateDeck(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
