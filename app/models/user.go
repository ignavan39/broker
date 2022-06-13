package models

type User struct {
	Id                string            `json:"id"`
	Password          string            `json:"password,omitempty"`
	Email             string            `json:"email"`
	Nickname          string            `json:"nickname"`
	LastName          string            `json:"lastName"`
	FirstName         string            `json:"firstName"`
	WorkspaceAccesses []WorkspaceAccess `json:"workspaceAccesses,omitempty"`
	UserPeers         []UserPeers       `json:"userPeers,omitempty"`
	AvatarURL         string            `json:"avatarURL"`
}
