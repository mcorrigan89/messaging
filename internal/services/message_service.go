package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/mcorrigan89/messaging/internal/serviceapis"
)

type MessageService struct {
	utils             ServicesUtils
	mailgun           *mailgun.MailgunImpl
	serviceApiClients *serviceapis.ServiceApiClients
	emailService      *EmailService
}

func NewMessageService(utils ServicesUtils, emailService *EmailService, serviceApiClients *serviceapis.ServiceApiClients) *MessageService {
	mailgun := mailgun.NewMailgun(utils.config.Mailgun.Domain, utils.config.Mailgun.APIKey)
	return &MessageService{
		utils:             utils,
		mailgun:           mailgun,
		serviceApiClients: serviceApiClients,
		emailService:      emailService,
	}
}

type SendVerificationEmailArgs struct {
	UserID uuid.UUID
	Link   string
}

func (service *MessageService) SendVerificationEmail(ctx context.Context, args SendVerificationEmailArgs) (*string, error) {
	service.utils.logger.Info().Ctx(ctx).Str("user ID", args.UserID.String()).Msg("Sending Verification Email")
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	user, err := service.serviceApiClients.Identity.GetUserByID(ctx, args.UserID)
	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Error getting user by ID")
		return nil, err
	}
	res, err := service.emailService.SendVerificationEmail(ctx, EmailSendVerificationEmailArgs{
		ToEmail: user.Email,
		Link:    args.Link,
	})
	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Error sending verification email")
		return nil, err
	}

	return res, nil
}

type SendPasswordResetEmailArgs struct {
	UserID uuid.UUID
	Link   string
}

func (service *MessageService) SendPasswordResetEmail(ctx context.Context, args SendPasswordResetEmailArgs) (*string, error) {
	service.utils.logger.Info().Ctx(ctx).Str("user ID", args.UserID.String()).Msg("Sending Password Reset Email")
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	user, err := service.serviceApiClients.Identity.GetUserByID(ctx, args.UserID)
	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Error getting user by ID")
		return nil, err
	}
	res, err := service.emailService.SendPasswordResetEmail(ctx, EmailPasswordResetEmailArgs{
		ToEmail: user.Email,
		Link:    args.Link,
	})
	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Error sending password reset email")
		return nil, err
	}

	return res, nil
}
