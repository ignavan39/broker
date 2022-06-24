package dto

import (
	"broker/app/models"
	"errors"
)

type SendInvitationPayload struct {
	RicipientEmail *string `json:"ricipientEmail,omitempty"`
}

func (p *SendInvitationPayload) Validate() error {
	if p.RicipientEmail == nil || len(*p.RicipientEmail) == 0 {
		return errors.New("Choose an email to send invitation")
	}
	return nil
}

type SendInvitationResponse = models.Invitation

type GetInvitationsByWorkspaceResponse struct {
	Invitations []models.Invitation `json:"invitations"`
}
