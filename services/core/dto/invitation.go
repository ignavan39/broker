package dto

import (
	"broker/core/models"
	"errors"
)

type SendInvitationPayload struct {
	RicipientEmail string `json:"ricipientEmail"`
}

func (p *SendInvitationPayload) Validate() error {
	if len(p.RicipientEmail) == 0 {
		return errors.New("Choose an email to send invitation")
	}
	return nil
}

type SendInvitationResponse = models.Invitation

type GetInvitationsByWorkspaceResponse struct {
	Invitations []models.Invitation `json:"invitations"`
}

type CancelInvitationResponse = models.Invitation
