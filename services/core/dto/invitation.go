package dto

import (
	"broker/core/models"
	"errors"
)

type SendInvitationPayload struct {
	RecipientEmail string `json:"recipientEmail"`
}

func (p *SendInvitationPayload) Validate() error {
	if !isCorrectEmail(p.RecipientEmail) {
		return errors.New("email must be not empty string")
	}
	return nil
}

type SendInvitationResponse = models.Invitation

type GetInvitationsByWorkspaceResponse struct {
	Invitations []models.Invitation `json:"invitations"`
}

type CancelInvitationResponse = models.Invitation

type AcceptInvitationPayload struct {
	Code string `json:"code"`
}

func (p *AcceptInvitationPayload) Validate() error {
	if len(p.Code) == 0 {
		return errors.New("code must be not empty string")
	}
	return nil
}
