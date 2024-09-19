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

type SendEmailArgs struct {
	FromEmail string
	ToEmail   string
	Subject   string
	Body      string
}

func (service *EmailService) SendEmail(ctx context.Context, args SendEmailArgs) (*string, error) {
	service.utils.logger.Info().Ctx(ctx).Str("fromEmail", args.FromEmail).Str("toEmail", args.ToEmail).Str("subject", args.Subject).Msg("Sending Email")
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	welcomeEmail := templates.Base(templates.Welcome("Michael", "Welcome to the site"))
	htmlString := templates.RenderToString(ctx, welcomeEmail)

	message := service.mailgun.NewMessage(args.FromEmail, args.Subject, args.Body, args.ToEmail)
	message.SetHtml(htmlString)

	resp, id, err := service.mailgun.Send(ctx, message)

	if err != nil {
		service.utils.logger.Err(err).Ctx(ctx).Msg("Failed to send email")
		return nil, err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return &resp, nil
}
