package serviceapis

import (
	"sync"

	"github.com/mcorrigan89/messaging/internal/config"
	"github.com/rs/zerolog"
)

type ServiceApiClients struct {
	Identity *IdentityClientV1
}

func NewServiceApiClients(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup) *ServiceApiClients {
	identityClient := NewIdentityClientV1(cfg, logger, wg)

	return &ServiceApiClients{
		Identity: identityClient,
	}
}
