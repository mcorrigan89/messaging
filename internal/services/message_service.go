package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/mcorrigan89/messaging/internal/api"
	"github.com/mcorrigan89/messaging/internal/repositories"
	"github.com/mcorrigan89/messaging/internal/templates"
)

type MessageService struct {
	utils         ServicesUtils
	mailgun       *mailgun.MailgunImpl
	idenityClient *api.IdentityClientV1
	emailService  *EmailService
}

func NewMessageService(utils ServicesUtils, emailRepo *repositories.EmailRepository) *MessageService {
	mailgun := mailgun.NewMailgun(utils.config.Mailgun.Domain, utils.config.Mailgun.APIKey)
	return &MessageService{
		utils:           utils,
		mailgun:         mailgun,
		emailRepository: emailRepo,
	}
}

type SendVerificationEmailArgs struct {
	UserID uuid.UUID
	Link   string
}

func (service *MessageService) SendVerificationEmail(ctx context.Context, args SendVerificationEmailArgs) (*string, error) {
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
