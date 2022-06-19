package dto

import (
	"broker/app/models"
	"context"
	"errors"
)

type CreateWorkspacePayload struct {
	Name      string `json:"name"`
	IsPrivate bool   `json:"isPrivate"`
}

func (p *CreateWorkspacePayload) Validate() error {
	if len(p.Name) == 0 {
		return errors.New("workspace name must be not empty string")
	}
	return nil
}

type CreateWorkspaceResponse = models.Workspace

type GetManyByUserResponse struct {
	Workspaces []models.Workspace `json:"workspaces"`
}

type Meta struct {
	QueueName    string `json:"queueName"`
	ExchangeName string `json:"exchange"`
}

type CreateWorkspaceConnectionPayload struct {
	WorkspaceID string `json:"workspaceId"`
}

type CreateWorkspaceConnectionBase struct {
	Meta
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Vhost    string `json:"vhost"`
	Password string `json:"password"`
}

type CreateWorkspaceConnectionResponse struct {
	Consume CreateWorkspaceConnectionBase `json:"consume"`
	Publish CreateWorkspaceConnectionBase `json:"publish"`
}

type PeerEnvelope struct {
	ctx context.Context

	Meta    Meta            `json:"meta"`
	Event   PeerEventType   `json:"event"`
	Payload *models.Message `json:"payload,omitempty"`
	FromId  string          `json:"fromId"`
	PeerId  string          `json:"peerId"`
}
