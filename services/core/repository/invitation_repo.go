package repository

import (
	"broker/core/models"
	"time"
)

type InvitationRepository interface {
	AcceptInvitation(userID string, code string) error
	CreateInvitation(senderID string, workspaceID string, recipientEmail string) (*models.Invitation, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) ([]models.Invitation, error)
	CancelInvitation(userID string, invitationID string) (*models.Invitation, error)
	DeleteExpiredInvitations(duration time.Duration) error
}
