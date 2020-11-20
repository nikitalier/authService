package main

import (
	golog "log"
	"os"
	"time"

	"github.com/nikitalier/authService/config"
	"github.com/nikitalier/authService/pkg/application"
	"github.com/nikitalier/authService/pkg/provider"
	"github.com/nikitalier/authService/pkg/repository"
	"github.com/nikitalier/authService/pkg/service"

	"github.com/rs/zerolog"
)

func main() {
	var (
		appConfig config.Config
		logger    zerolog.Logger
	)

	err := appConfig.Load()
	if err != nil {
		golog.Fatalf("%v", err)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "timestamp"
	zerolog.LevelFieldName = "logLevel"

	logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().CallerWithSkipFrameCount(2).Timestamp().Logger()

	provider := provider.New(&appConfig.DataBase[0], logger)

	provider.Open()

	rep := repository.New(provider.GetCon(), logger)

	err = rep.PingDB()
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	logger.Info().Msg("DB connected")

	serv := service.New(rep, &logger)

	app := application.New(&application.Options{Serv: appConfig.ServerOpt, Svc: serv, Logger: logger})

	app.Start()
}
