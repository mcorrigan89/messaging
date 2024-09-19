package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mcorrigan89/messaging/internal/repositories"
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

	message := service.mailgun.NewMessage(args.FromEmail, args.Subject, args.Body, args.ToEmail)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	resp, id, err := service.mailgun.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return &resp, nil
}
