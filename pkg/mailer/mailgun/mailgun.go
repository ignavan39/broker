package mailer

import (
	"broker/app/config"
	"context"

	"github.com/mailgun/mailgun-go/v4"
)

type MailgunApi struct {
	mg     *mailgun.MailgunImpl
	config config.MailgunConfig
}

func NewMailgunApi(config config.Config) *MailgunApi {
	mg := mailgun.NewMailgun(config.MailgunConfig.Domain, config.MailgunConfig.PrivateKey)
	return &MailgunApi{
		mg:     mg,
		config: config.MailgunConfig,
	}
}

// func (m *MailgunApi) SendMail(msg string, subject string, recipient string) (string, string, error) {
// 	ctx := context.Background()
// 	message := m.mg.NewMessage(m.config.Sender, subject, msg, recipient)
// 	return m.mg.Send(ctx, message)
// }
