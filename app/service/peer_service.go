package service

import (
	"broker/app/dto"
	"context"
)

type PeerService interface {
	CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionResponse, error)
	// Create(userId string, payload dto.CreatePeerPayload) (*dto.CreatePeerResponse, error)
}
