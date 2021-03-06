package smtp

import (
	"broker/smtp/sender/dto"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type SmtpMailer struct{}

func NewSmtpMailer() *SmtpMailer {
	return &SmtpMailer{}
}

func (m *SmtpMailer) SendMail(ctx context.Context, msg string, subject string, recipient string) (string, string, error) {
	payload := dto.SendMailPayload{
		Message:   msg,
		Subject:   subject,
		Recipient: recipient,
	}

	json, err := json.Marshal(payload)

	if err != nil {
		return "", "", err
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, "http://smtp.broker.loc/api/v1/mail/send", bytes.NewBuffer(json))

	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		return "", "", err
	}

	return msg, recipient, nil
}
