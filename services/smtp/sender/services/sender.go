package services

import (
	"broker/smtp/config"
	"broker/smtp/sender/dto"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Send(payload dto.SendMailPayload) error {
	host := config.GetConfig().Host
	addr := config.GetConfig().Address
	from := config.GetConfig().Email

	auth := smtp.PlainAuth("", from, config.GetConfig().Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from); err != nil {
		return err
	}

	if err = c.Rcpt(payload.Recipient); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	message := s.buildMessage(payload, from)

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}

func (s *Sender) buildMessage(payload dto.SendMailPayload, sender string) string {
	delimiter := "beer=happiness"

	msg := fmt.Sprintf("From: %s\r\n", sender)
	msg += fmt.Sprintf("To: %s\r\n", payload.Recipient)
	msg += fmt.Sprintf("Subject: %s\r\n", payload.Subject)
	msg += fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/mixed; charset=\"UTF-8\"; boundary=%s\r\n", delimiter)
	msg += fmt.Sprintf("\r\n%s\r\n", payload.Message)

	//message
	msg += fmt.Sprintf("\r\n--%s\r\n", delimiter)
	msg += fmt.Sprintf("Content-Type: text/html; charset=\"utf-8\"\r\n")
	msg += fmt.Sprintf("\r\n%s\r\n", payload.Message)

	//attachments
	if payload.Attachments != nil {
		msg += fmt.Sprintf("\r\n--%s\r\n", delimiter)
		for attachmentName, attachment := range *payload.Attachments {
			msg += fmt.Sprintf("\r\n--%s\r\n", delimiter)
			msg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
			msg += "Content-Transfer-Encoding: base64\r\n"
			msg += "Content-Disposition: attachment;filename=\"" + attachmentName + "\"\r\n"
			msg += "\r\n" + attachment
		}
	}

	return msg
}
