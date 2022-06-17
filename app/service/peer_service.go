package service

import "broker/app/dto"

type PeerService interface {
	Create(userId string, payload dto.CreatePeerPayload) (*dto.CreatePeerResponse, error)
}
