package mailer

import "context"

type Mailer interface {
	SendMail(ctx context.Context, msg string, subject string, recipient string) (string, string, error)
}
