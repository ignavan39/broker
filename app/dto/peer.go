package dto

import (
	"broker/app/models"
	"context"
	"errors"
)

type CreatePeerPayload struct {
	WorkspaceId string `json:"workspaceId"`
}

func (p *CreatePeerPayload) Validate() error {
	if len(p.WorkspaceId) == 0 {
		return errors.New("workspace id must be not empty")
	}
	return nil
}

type CreatePeerResponse = models.UserPeers

type PeerEventType = string

var (
	NewMessageEvent    PeerEventType = "new_message"
	ReadMessageEvent   PeerEventType = "read_message"
	TypingEvent        PeerEventType = "typing"
	DeleteMessageEvent PeerEventType = "delete_message"
	EditMessageEvent   PeerEventType = "edit_message"
	BlockUserEvent     PeerEventType = "block_user"
	OnlineUserEvent    PeerEventType = "online_user"
)

type Meta struct {
	ReportKey    string `json:"report"`
	QueueName    string `json:"queueName"`
	ExchangeName string `json:"exchange"`
}

type CreatePeerConnectionPayload struct {
	PeerId string `json:"peerId"`
}

type CreatePeerConnectionResponse struct {
	Meta
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Vhost    string `json:"vhost"`
	Password string `json:"password"`
}

type PeerEnvelope struct {
	ctx context.Context

	Meta    Meta            `json:"meta"`
	Event   PeerEventType   `json:"event"`
	Payload *models.Message `json:"payload,omitempty"`
	FromId  string          `json:"fromId"`
	PeerId  string          `json:"peerId"`
}
