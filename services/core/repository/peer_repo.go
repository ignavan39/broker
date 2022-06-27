package repository

import "broker/core/models"

type PeerRepository interface {
	//FindOneOrCreate()
	GetMany(userID string, workspaceID string) ([]models.Peer, error)
}
