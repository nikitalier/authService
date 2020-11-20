package application

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikitalier/authService/config"
	"github.com/nikitalier/authService/pkg/service"

	"github.com/rs/zerolog"
)

type Application struct {
	serv   *http.Server
	svc    *service.Service
	logger zerolog.Logger
}

type Options struct {
	Svc    *service.Service
	Serv   config.ServerOpt
	Logger zerolog.Logger
}

func New(opt *Options) *Application {
	app := &Application{
		svc: opt.Svc,
		serv: &http.Server{
			Addr: opt.Serv.Port,
		},
		logger: opt.Logger,
	}

	app.serv.Handler = app.setupRoutes()

	app.logger.Info().Msg("App started on port" + opt.Serv.Port)

	return app
}

func (app *Application) Start() {
	app.serv.ListenAndServe()
}

func (app *Application) setupRoutes() *mux.Router {
	r := &mux.Router{}

	r.HandleFunc("/gettokens", app.getPairTokens).Methods("GET")
	r.HandleFunc("/refreshtoken", app.refreshToken).Methods("POST")

	return r
}
