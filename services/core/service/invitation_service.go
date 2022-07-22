package service

import (
	"broker/core/dto"
	"context"
)

type InvitationService interface {
	AcceptInvitation(payload dto.AcceptInvitationPayload, userID string) error
	RejectInvitation(payload dto.RejectInvitationPayload, userID string) error
	SendInvitation(ctx context.Context, payload dto.SendInvitationPayload, userID string, workspaceID string) (*dto.SendInvitationResponse, error)
	GetInvitationsByWorkspaceID(userID string, workspaceID string) (*dto.GetInvitationsByWorkspaceResponse, error)
	CancelInvitation(senderID string, invitationID string) (*dto.CancelInvitationResponse, error)
	Connect(ctx context.Context, userID string) (*dto.ConnectInvitationResponse, error)
}
