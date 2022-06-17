package service

import (
	"broker/app/dto"
	"context"
)

type Consumer interface{
	Start(ctx context.Context) error
	Consume(func (ctx context.Context, payload any)) error
	CreateConnection(fromId string,payload dto.CreatePeerConnectionPayload) (*dto.CreatePeerConnectionResponse,error)
}