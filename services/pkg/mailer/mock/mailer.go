package mock

import (
	"broker/pkg/logger"
	"context"
)

type MockMailer struct{}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (mm *MockMailer) SendMail(ctx context.Context, msg string, subject string, recipient string) (string, string, error) {
	logger.Logger.Printf("Email %s message: %s", recipient, msg)
	return msg, recipient, nil
}
