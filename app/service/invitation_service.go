package service

import "broker/app/dto"

type InvitationService interface {
	SendInvitation(payload dto.SendInvitationPayload, userID string, workspaceID string) (*dto.SendInvitationResponse, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) (*dto.GetInvitationsByWorkspaceResponse, error)
	CancelInvitation(payload dto.CancelInvitationPayload, senderID string, wrokspaceID string) (*dto.CancelInvitationResponse, error)
}
