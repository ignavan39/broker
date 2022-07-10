package peer

import (
	"broker/core/dto"
	"broker/core/service/peer/consumer"
	"broker/core/service/peer/publisher"
	"context"
)

type PeerService struct {
	consumer  *consumer.Consumer
	publisher *publisher.Publisher
}

func NewPeerService(
	consumer *consumer.Consumer,
	publisher *publisher.Publisher,
) *PeerService {
	return &PeerService{
		consumer:  consumer,
		publisher: publisher,
	}
}

func (ps *PeerService) CreateConnection(ctx context.Context, senderID string, payload dto.CreateWorkspaceConnectionPayload) (*dto.CreateWorkspaceConnectionResponse, error) {
	consume, err := ps.consumer.CreateConnection(ctx, senderID, payload)

	if err != nil {
		return nil, err
	}

	publish, err := ps.publisher.CreateConnection(ctx, senderID, payload)

	if err != nil {
		return nil, err
	}

	return &dto.CreateWorkspaceConnectionResponse{
		Consume: *consume,
		Publish: *publish,
	}, nil
}

func (ps *PeerService) Run() {
	// TODO save to database and move to member PeerService

	go ps.consumer.Consume(func(payload dto.PeerEnvelope) {
		// ps.publisher.Publish(payload)
	})
}
