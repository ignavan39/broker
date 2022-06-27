package repository

import (
	"broker/core/models"
)

type UserRepository interface {
	Create(nickname string, email string, password string, lastName string, firstName string) (*models.User, error)
	GetOneByEmail(email string) (*models.User, error)
	GetOneByNickname(nickname string) (*models.User, error)
	GetEmailById(userID string) (string, error)
	CheckInvites(userID string, email string) error
	SendInvitation(senderID string, workspaceID string, ricipientEmail string) (*models.Invitation, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) ([]models.Invitation, error)
	CancelInvitation(senderID string, invitationID string) (*models.Invitation, error)
}
