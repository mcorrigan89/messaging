package api

import (
	"context"
	"errors"
	"sync"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	messagingv1 "github.com/mcorrigan89/messaging/gen/serviceapis/messaging/v1"
	"github.com/mcorrigan89/messaging/internal/config"
	"github.com/mcorrigan89/messaging/internal/services"

	"github.com/rs/zerolog"
)

type MessagingServerV1 struct {
	config   *config.Config
	wg       *sync.WaitGroup
	logger   *zerolog.Logger
	services *services.Services
}

func newMessagingProtoUrlServer(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup, services *services.Services) *MessagingServerV1 {
	return &MessagingServerV1{
		config:   cfg,
		wg:       wg,
		logger:   logger,
		services: services,
	}
}

func (s *MessagingServerV1) SendVerificationEmail(ctx context.Context, req *connect.Request[messagingv1.SendVerificationEmailRequest]) (*connect.Response[messagingv1.SendVerificationEmailResponse], error) {
	userId := req.Msg.UserId
	link := req.Msg.VerificationLink

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error parsing user ID")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if link == "" {
		err := errors.New("link is empty")
		s.logger.Err(err).Ctx(ctx).Msg("Link is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	emailStatus, err := s.services.MessageService.SendPasswordResetEmail(ctx, services.SendPasswordResetEmailArgs{
		UserID: userUUID,
		Link:   link,
	})
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error sending verification email")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&messagingv1.SendVerificationEmailResponse{
		Message: *emailStatus,
	})
	res.Header().Set("Messaging-Version", "v1")
	return res, nil
}

func (s *MessagingServerV1) SendPasswordResetEmail(ctx context.Context, req *connect.Request[messagingv1.SendPasswordResetEmailRequest]) (*connect.Response[messagingv1.SendPasswordResetEmailResponse], error) {
	userId := req.Msg.UserId
	link := req.Msg.PasswordResetLink

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error parsing user ID")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if link == "" {
		err := errors.New("link is empty")
		s.logger.Err(err).Ctx(ctx).Msg("Link is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	emailStatus, err := s.services.MessageService.SendPasswordResetEmail(ctx, services.SendPasswordResetEmailArgs{
		UserID: userUUID,
		Link:   link,
	})
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error sending password reset email")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&messagingv1.SendPasswordResetEmailResponse{
		Message: *emailStatus,
	})
	res.Header().Set("Messaging-Version", "v1")
	return res, nil
}
