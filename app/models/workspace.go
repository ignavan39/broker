package models

import "time"

type WorkspaceAccessesType = string

var (
	ADMIN WorkspaceAccessesType = "ADMIN"
	USER  WorkspaceAccessesType = "USER"
)

type WorkspaceAccess struct {
	Id     string                `json:"id"`
	UserId string                `json:"userId"`
	Type   WorkspaceAccessesType `json:"type"`
}

type Workspace struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	IsPrivate bool      `json:"isPrivate"`
	Users     []User    `json:"users,omitempty"`
}
