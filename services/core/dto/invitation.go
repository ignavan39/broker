package dto

import (
	"broker/core/models"
	"errors"
)

type SendInvitationPayload struct {
	RicipientEmail string `json:"ricipientEmail"`
}

func (p *SendInvitationPayload) Validate() error {
	if !isCorrectEmail(p.RicipientEmail) {
		return errors.New("email must be not empty string")
	}
	return nil
}

type SendInvitationResponse = models.Invitation

type GetInvitationsByWorkspaceResponse struct {
	Invitations []models.Invitation `json:"invitations"`
}

type CancelInvitationResponse = models.Invitation
