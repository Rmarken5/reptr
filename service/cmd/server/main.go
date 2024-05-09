package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	exAPI "github.com/rmarken/reptr/api"
	"github.com/rmarken/reptr/service/cmd"
	"github.com/rmarken/reptr/service/internal/api"
	"github.com/rmarken/reptr/service/internal/api/middlewares"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"os"
)

var (
	config cmd.Config
	log    zerolog.Logger
)

func init() {
	log = zerolog.New(os.Stdout).With().Str("program", "reptr server").Logger()
	env := os.Getenv("ENV")
	var path string
	switch env {
	case "local":
		path = "./config/local.yaml"
	case "dev":
		path = "./config/dev.yaml"
	default:
		log.Info().Msg("loading from env")
		config = cmd.LoadConfFromEnv(log)
		return
	}
	log.Info().Msgf("loading from file: %s", path)
	config = cmd.LoadConfigFromFile(log, path)
}

func main() {
	ctx := context.Background()

	db := cmd.MustConnectMongo(ctx, log, config)
	defer db.Client().Disconnect(ctx)
	repo := cmd.MustLoadRepo(log, db)
	l := cmd.MustLoadLogic(log, repo)

	sessionController := cmd.MustLoadSessionLogic(log, l, repo)
	authenticator := cmd.MustLoadAuth(ctx, log, config, repo)
	p := cmd.MustLoadProvider(log, repo)
	store := sessions.NewCookieStore([]byte(config.SessionKey))
	serverImpl := api.New(log, l, p, authenticator, sessionController, store)

	router := mux.NewRouter()

	wrapper := exAPI.ServerInterfaceWrapper{
		Handler: serverImpl,
	}

	router.HandleFunc("/login", wrapper.LoginPage).Methods(http.MethodGet)
	router.HandleFunc("/login", wrapper.Login).Methods(http.MethodPost)
	router.HandleFunc("/register", wrapper.Register).Methods(http.MethodPost)
	router.HandleFunc("/register", wrapper.RegistrationPage).Methods(http.MethodGet)

	styleRouter := router.PathPrefix("/styles").Subrouter()
	styleRouter.HandleFunc("/{path}/{style_name}", wrapper.ServeStyles).Methods(http.MethodGet)

	pageRoute := router.PathPrefix("/page").Subrouter()
	pageRoute.HandleFunc("/home", wrapper.HomePage)
	pageRoute.HandleFunc("/create-group", wrapper.CreateGroupPage).Methods(http.MethodGet)
	pageRoute.HandleFunc("/create-group", wrapper.CreateGroup).Methods(http.MethodPost)
	pageRoute.HandleFunc("/group/{groupID}", wrapper.GroupPage).Methods(http.MethodGet)
	pageRoute.HandleFunc("/create-deck/{group_id}", wrapper.CreateDeckPage).Methods(http.MethodGet)
	pageRoute.HandleFunc("/create-deck", wrapper.CreateDeckPage).Methods(http.MethodGet)
	pageRoute.HandleFunc("/create-deck/{group_id}", wrapper.CreateDeck).Methods(http.MethodPost)
	pageRoute.HandleFunc("/create-deck", wrapper.CreateDeck).Methods(http.MethodPost)
	pageRoute.HandleFunc("/create-cards/{deck_id}", wrapper.CreateCardForDeck).Methods(http.MethodPost)
	pageRoute.HandleFunc("/create-cards/{deck_id}", wrapper.GetCreateCardsForDeckPage).Methods(http.MethodGet)
	pageRoute.HandleFunc("/create-cards-content/{deck_id}", wrapper.GetCreateCardsForDeckContent).Methods(http.MethodGet)
	pageRoute.HandleFunc("/add-card/{deck_id}", wrapper.GetCardsForDeck).Methods(http.MethodGet)
	pageRoute.HandleFunc("/front-of-card/{deck_id}/{card_id}", wrapper.FrontOfCard).Methods(http.MethodGet)
	pageRoute.HandleFunc("/back-of-card/{deck_id}/{card_id}", wrapper.BackOfCard).Methods(http.MethodGet)
	pageRoute.HandleFunc("/view-deck/{deck_id}", wrapper.ViewDeck).Methods(http.MethodGet)
	pageRoute.HandleFunc("/upvote-card/{card_id}/{direction}", wrapper.VoteCard).Methods(http.MethodPost)
	pageRoute.HandleFunc("/answered-correct/{session_id}", wrapper.UpdateCardCorrect).Methods(http.MethodPost)
	pageRoute.HandleFunc("/answered-incorrect/{session_id}", wrapper.UpdateCardIncorrect).Methods(http.MethodPost)

	pageRoute.Use(
		middlewares.Session(log, store),
		middlewares.Authenticate(log, authenticator),
		middlewares.ExchangeSubjectForUser(log, p))

	secureRoute := router.PathPrefix("/secure").Subrouter()
	secureRoute.HandleFunc("/api/v1/deck", wrapper.AddDeck).Methods(http.MethodPost)
	secureRoute.HandleFunc("/api/v1/group", wrapper.AddGroup).Methods(http.MethodPost)
	secureRoute.HandleFunc("/api/v1/group/{group_id}/deck/{deck_id}", wrapper.AddDeckToGroup).Methods("PUT")
	secureRoute.HandleFunc("/api/v1/groups", wrapper.GetGroups).Methods(http.MethodGet)
	secureRoute.HandleFunc("/api/v1/card-input/{card-num}", wrapper.GetCardInput).Methods(http.MethodGet)

	secureRoute.Use(
		middlewares.Authenticate(log, authenticator),
		middlewares.ExchangeSubjectForUser(log, p))

	s := &http.Server{
		Handler: router,
		Addr:    net.JoinHostPort("0.0.0.0", config.PORT),
	}

	log.Fatal().Err(s.ListenAndServe())

}
