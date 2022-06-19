package publisher

import (
	"broker/app/dto"
	"context"
)

type Publisher interface {
	Publish(payload dto.PeerEnvelope)
	CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error)
}
