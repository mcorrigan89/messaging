package api

import (
	"context"
	"sync"

	"connectrpc.com/connect"

	messagingv1 "github.com/mcorrigan89/messaging/internal/api/serviceapis/messaging/v1"
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

func (s *MessagingServerV1) SendEmail(ctx context.Context, req *connect.Request[messagingv1.SendEmailRequest]) (*connect.Response[messagingv1.SendEmailResponse], error) {
	fromEmail := req.Msg.FromEmail
	toEmail := req.Msg.ToEmail
	subject := req.Msg.Subject
	body := req.Msg.Body

	emailStatus, err := s.services.EmailService.SendEmail(ctx, services.SendEmailArgs{
		FromEmail: fromEmail,
		ToEmail:   toEmail,
		Subject:   subject,
		Body:      body,
	})
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error sending email")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&messagingv1.SendEmailResponse{
		Status: *emailStatus,
	})
	res.Header().Set("Messaging-Version", "v1")
	return res, nil
}
