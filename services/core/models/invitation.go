package models

type InvitationStatus = string

var (
	PENDING  InvitationStatus = "PENDING"
	ACCEPTED InvitationStatus = "ACCEPTED"
	CANCELED InvitationStatus = "CANCELED"
)

type Invitation struct {
	ID             string           `json:"id"`
	Sender         User             `json:"sender"`
	RicipientEmail string           `json:"ricipientEmail"`
	WorkspaceID    string           `json:"workspaceID"`
	Status         InvitationStatus `json:"status"`
}
