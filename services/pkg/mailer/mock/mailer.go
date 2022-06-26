package mock

import (
	"context"

	blogger "github.com/sirupsen/logrus"
)

type MockMailer struct{}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (mm *MockMailer) SendMail(ctx context.Context, msg string, subject string, recipient string) (string, string, error) {
	blogger.Printf("Email %s message: %s", recipient, msg)
	return msg, recipient, nil
}
