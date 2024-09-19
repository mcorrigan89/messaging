package services

import (
	"sync"

	"github.com/mcorrigan89/messaging/internal/config"
	"github.com/mcorrigan89/messaging/internal/repositories"
	"github.com/rs/zerolog"
)

type ServicesUtils struct {
	logger *zerolog.Logger
	wg     *sync.WaitGroup
	config *config.Config
}

type Services struct {
	utils        ServicesUtils
	EmailService *EmailService
}

func NewServices(repositories *repositories.Repositories, cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup) Services {
	utils := ServicesUtils{
		logger: logger,
		wg:     wg,
		config: cfg,
	}

	emailService := NewEmailService(utils, repositories.EmailRepository)

	return Services{
		utils:        utils,
		EmailService: emailService,
	}
}
