package repository

import "broker/app/models"

type PeerRepository interface {
	//FindOneOrCreate()
	GetMany(userID string, workspaceID string) ([]models.Peer, error)
}
