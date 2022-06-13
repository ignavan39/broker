package models

import "time"

type UserPeers struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	Peer      Peer      `json:"peer"`
	IsBlocked bool      `json:"isBlocked"`
	CreatedAt time.Time `json:"createdAt"`
}

type Peer struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	WorkspaceId string    `json:"workspaceId"`
}
