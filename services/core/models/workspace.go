package models

import "time"

type WorkspaceAccessesType = string

var (
	ADMIN     WorkspaceAccessesType = "ADMIN"
	USER      WorkspaceAccessesType = "USER"
	MODERATOR WorkspaceAccessesType = "MODERATOR"
)

var Roles = []WorkspaceAccessesType{ADMIN, USER, MODERATOR}

type WorkspaceAccess struct {
	ID     string                `json:"id"`
	UserId string                `json:"userId"`
	Type   WorkspaceAccessesType `json:"type"`
}

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	IsPrivate bool      `json:"isPrivate"`
	Users     []User    `json:"users,omitempty"`
}
