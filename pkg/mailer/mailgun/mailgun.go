package mailer

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunApi struct {
	mg         *mailgun.MailgunImpl
	privateKey string
	publicKey  string
	domain     string
	sender     string
}

func NewMailgunApi(
	privateKey string,
	publicKey string,
	domain string,
	sender string,
) *MailgunApi {
	mg := mailgun.NewMailgun(domain, privateKey)
	return &MailgunApi{
		mg:         mg,
		privateKey: privateKey,
		publicKey:  publicKey,
		domain:     domain,
		sender:     sender,
	}
}

// return message,message id,error
func (m *MailgunApi) SendMail(ctx context.Context, msg string, subject string, recipient string) (string, string, error) {
	message := m.mg.NewMessage(m.sender, subject, msg, recipient)
	return m.mg.Send(ctx, message)
}
