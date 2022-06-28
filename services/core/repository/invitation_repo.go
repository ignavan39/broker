package repository

import "broker/core/models"

type InvitationRepository interface {
	CheckInvites(userID string, email string) error
	SendInvitation(senderID string, workspaceID string, ricipientEmail string) (*models.Invitation, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) ([]models.Invitation, error)
	CancelInvitation(invitationID string) (*models.Invitation, error)
}
