package consumer

import (
	"broker/core/dto"
	"context"
)

type Consumer interface {
	Consume(handler func(payload dto.PeerEnvelope))
	CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error)
}
