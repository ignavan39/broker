package service

import (
	"broker/app/dto"
	"context"
)
type Consumer interface{
	Start(ctx context.Context) error
	Consume(func (ctx context.Context, payload dto.PeerEnvelope)) error
	CreateConnection(ctx context.Context,recipientID string,payload dto.CreatePeerConnectionPayload) (*dto.CreatePeerConnectionResponse,error)
}