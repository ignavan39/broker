package models

import "time"

type UserPeers struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Peer      Peer      `json:"peer"`
	IsBlocked bool      `json:"isBlocked"`
	CreatedAt time.Time `json:"createdAt"`
}

type Peer struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	WorkspaceId string    `json:"workspaceId"`
}
