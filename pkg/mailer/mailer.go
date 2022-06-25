package mailer

type Mailer interface {
	SendMail(msg string, subject string, recipient string) (string, string, error)
}
