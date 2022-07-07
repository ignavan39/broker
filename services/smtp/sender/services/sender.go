package services

import (
	"broker/smtp/config"
	"broker/smtp/sender/dto"
	"net/smtp"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Send(payload dto.SendMailPayload) error {
	host := config.GetConfig().Host
	from := config.GetConfig().Email

	auth := smtp.PlainAuth("", from, config.GetConfig().Password, host)

	err := smtp.SendMail(config.GetConfig().Address, auth, from, []string{payload.Recipient}, []byte(payload.Message))

	if err != nil {
		return err
	}

	return nil
}