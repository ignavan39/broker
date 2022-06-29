package service

import "broker/core/dto"

type InvitationService interface {
	AcceptInvitation(userID string, code string) error
	SendInvitation(payload dto.SendInvitationPayload, userID string, workspaceID string) (*dto.SendInvitationResponse, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) (*dto.GetInvitationsByWorkspaceResponse, error)
	CancelInvitation(senderID string, invitationID string) (*dto.CancelInvitationResponse, error)
}
