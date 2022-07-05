package models

import "time"

type InvitationStatus = string
type SystemStatus = string

var (
	PENDING  InvitationStatus = "PENDING"
	ACCEPTED InvitationStatus = "ACCEPTED"
	CANCELED InvitationStatus = "CANCELED"
	EXPIRED InvitationStatus = "EXPIRED"

	CREATED   SystemStatus = "CREATED"
	SEND      SystemStatus = "SEND"
	DELIVERED SystemStatus = "DELIVERED"
	REJECT    SystemStatus = "REJECT"
)

type Invitation struct {
	ID             string           `json:"id"`
	CreatedAt      time.Time        `json:"createdAt"`
	Sender         User             `json:"sender"`
	RecipientEmail string           `json:"recipientEmail"`
	WorkspaceID    string           `json:"workspaceID"`
	Status         InvitationStatus `json:"status"`
	SystemStatus   SystemStatus     `json:"systemStatus"`
	Code           string           `json:"code"`
}
