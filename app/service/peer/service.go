package peer

import (
	"broker/app/dto"
	"broker/app/service/peer/consumer"
	"broker/app/service/peer/publisher"
	"context"
)

type PeerService struct {
	consumer  consumer.Consumer
	publisher publisher.Publisher
}

func NewPeerService(
	consumer consumer.Consumer,
	// publisher publisher.Publisher,
) *PeerService {
	return &PeerService{
		consumer: consumer,
		// publisher: publisher,
	}
}

func (ps *PeerService) CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionResponse, error) {
	consume, err := ps.consumer.CreateConnection(ctx, senderID, payload)

	if err != nil {
		return nil, err
	}

	publish, err := ps.consumer.CreateConnection(ctx, senderID, payload)

	if err != nil {
		return nil, err
	}

	return &dto.CreateWorkspaceConnectionResponse{
		Consume: *consume,
		Publish: *publish,
	}, nil
}
