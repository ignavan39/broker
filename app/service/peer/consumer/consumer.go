package consumer

import (
	"broker/app/dto"
	"context"
)

type Consumer interface {
	Consume(handler func(payload dto.PeerEnvelope))
	CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionBase, error)
}
