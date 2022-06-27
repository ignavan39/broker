package dto

import (
	"broker/core/models"
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
