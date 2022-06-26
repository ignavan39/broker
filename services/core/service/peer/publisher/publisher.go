package publisher

import (
	"broker/core/dto"
	"context"
)

type Publisher interface {
	Publish(workspaceID string, payload dto.PeerEnvelope) error
	CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error)
}
