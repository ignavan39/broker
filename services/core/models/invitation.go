package models

type InvitationStatus = string

var (
	PENDING  InvitationStatus = "PENDING"
	ACCEPTED InvitationStatus = "ACCEPTED"
	CANCELED InvitationStatus = "CANCELED"
)

type Invitation struct {
	ID             string           `json:"id"`
	SenderID       string           `json:"senderID"`
	RicipientEmail string           `json:"ricipientEmail"`
	WorkspaceID    string           `json:"workspaceID"`
	Status         InvitationStatus `json:"status"`
}
