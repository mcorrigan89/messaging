package main

import (
	"os"
	"sync"

	"github.com/mcorrigan89/messaging/internal/api"
	"github.com/mcorrigan89/messaging/internal/config"
	"github.com/mcorrigan89/messaging/internal/repositories"
	"github.com/mcorrigan89/messaging/internal/services"
	"github.com/rs/zerolog"
)

type application struct {
	config         config.Config
	wg             *sync.WaitGroup
	logger         *zerolog.Logger
	services       *services.Services
	protoServer    *api.ProtoServer
	identityClient *api.IdentityClientV1
}

func main() {

	cfg := config.Config{}
	config.LoadConfig(&cfg)

	logger := getLogger()

	db, err := openDBPool(&cfg)
	if err != nil {
		logger.Err(err).Msg("Failed to open database connection")
		os.Exit(1)
	}
	defer db.Close()

	wg := sync.WaitGroup{}

	repositories := repositories.NewRepositories(db, &logger, &wg)
	services := services.NewServices(&repositories, &cfg, &logger, &wg)
	protoServer := api.NewProtoServer(&cfg, &logger, &wg, &services)
	identityClient := api.NewIdentityClientV1(&cfg, &logger, &wg)

	app := &application{
		wg:             &wg,
		config:         cfg,
		logger:         &logger,
		services:       &services,
		protoServer:    protoServer,
		identityClient: identityClient,
	}

	err = app.serve()
	if err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
