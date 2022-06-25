package mailer

type MailLetter struct {
	Subject  string `yaml:"subject"`
	Template string `yaml:"template"`
}

type MailConfig struct {
	Sender  string                `yaml:"sender"`
	Letters map[string]MailLetter `yaml:"letters"`
}
