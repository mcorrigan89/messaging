package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mcorrigan89/messaging/internal/repositories"
	"github.com/mcorrigan89/messaging/internal/templates"
)

type EmailService struct {
	utils           ServicesUtils
	mailgun         *mailgun.MailgunImpl
	emailRepository *repositories.EmailRepository
}

func NewEmailService(utils ServicesUtils, emailRepo *repositories.EmailRepository) *EmailService {
	mailgun := mailgun.NewMailgun(utils.config.Mailgun.Domain, utils.config.Mailgun.APIKey)
	return &EmailService{
		utils:           utils,
		mailgun:         mailgun,
		emailRepository: emailRepo,
	}
}

type EmailSendVerificationEmailArgs struct {
	ToEmail string
	Link    string
}

func (service *EmailService) SendVerificationEmail(ctx context.Context, args EmailSendVerificationEmailArgs) (*string, error) {
	service.utils.logger.Info().Ctx(ctx).Str("toEmail", args.ToEmail).Msg("Sending Verification Email")
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	verificationEmail := templates.Base(templates.VerificationEmail(args.Link))
	htmlString := templates.RenderToString(ctx, verificationEmail)

	message := service.mailgun.NewMessage(service.utils.config.Mailgun.Email, "Verify your email", "", args.ToEmail)
	message.SetHtml(htmlString)

	resp, id, err := service.mailgun.Send(ctx, message)

	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Failed to send verification email")
		return nil, err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return &resp, nil
}

type EmailPasswordResetEmailArgs struct {
	ToEmail string
	Link    string
}

func (service *EmailService) SendPasswordResetEmail(ctx context.Context, args EmailPasswordResetEmailArgs) (*string, error) {
	service.utils.logger.Info().Ctx(ctx).Str("toEmail", args.ToEmail).Msg("Sending Password Reset Email")
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	passwordResetEmail := templates.Base(templates.PasswordResetEmail(args.Link))
	htmlString := templates.RenderToString(ctx, passwordResetEmail)

	message := service.mailgun.NewMessage(service.utils.config.Mailgun.Email, "Password Reset", "", args.ToEmail)
	message.SetHtml(htmlString)

	resp, id, err := service.mailgun.Send(ctx, message)

	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Failed to send password reset email")
		return nil, err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return &resp, nil
}
