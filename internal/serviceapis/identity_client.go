package serviceapis

import (
	"context"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	identityv1 "github.com/mcorrigan89/messaging/gen/serviceapis/identity/v1"
	identityv1connect "github.com/mcorrigan89/messaging/gen/serviceapis/identity/v1/identityv1connect"
	"github.com/mcorrigan89/messaging/internal/entities"
	"github.com/mcorrigan89/messaging/internal/usercontext"

	"github.com/mcorrigan89/messaging/internal/config"

	"github.com/rs/zerolog"
)

type IdentityClientV1 struct {
	config *config.Config
	wg     *sync.WaitGroup
	logger *zerolog.Logger
	client identityv1connect.IdentityServiceClient
}

func NewIdentityClientV1(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup) *IdentityClientV1 {
	client := identityv1connect.NewIdentityServiceClient(
		http.DefaultClient,
		"http://localhost:4000",
	)
	return &IdentityClientV1{
		config: cfg,
		wg:     wg,
		logger: logger,
		client: client,
	}
}

func (c *IdentityClientV1) GetUserBySessionToken(ctx context.Context, token string) (*entities.User, error) {
	sessionToken := usercontext.ContextGetSession(ctx)
	req := connect.NewRequest(&identityv1.GetUserBySessionTokenRequest{
		Token: token,
	})

	req.Header().Set("x-session-token", sessionToken)

	response, err := c.client.GetUserBySessionToken(ctx, req)
	if err != nil {
		c.logger.Err(err).Msg("Error getting user by session token")
		return nil, err
	}
	uuid, err := uuid.Parse(response.Msg.User.Id)
	if err != nil {
		c.logger.Err(err).Msg("Error parsing user ID")
		return nil, err
	}
	userEntity := entities.NewUserEntity(entities.NewUserEntityArgs{
		ID:         uuid,
		GivenName:  response.Msg.User.GivenName,
		FamilyName: response.Msg.User.FamilyName,
		Email:      response.Msg.User.Email,
	})

	return userEntity, nil
}

func (c *IdentityClientV1) GetUserByID(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	sessionToken := usercontext.ContextGetSession(ctx)
	req := connect.NewRequest(&identityv1.GetUserByIdRequest{
		Id: userID.String(),
	})
	req.Header().Set("x-session-token", sessionToken)

	response, err := c.client.GetUserById(ctx, req)
	if err != nil {
		c.logger.Err(err).Msg("Error getting user by ID")
		return nil, err
	}
	uuid, err := uuid.Parse(response.Msg.User.Id)
	if err != nil {
		c.logger.Err(err).Msg("Error parsing user ID")
		return nil, err
	}
	userEntity := entities.NewUserEntity(entities.NewUserEntityArgs{
		ID:         uuid,
		GivenName:  response.Msg.User.GivenName,
		FamilyName: response.Msg.User.FamilyName,
		Email:      response.Msg.User.Email,
	})

	return userEntity, nil
}
