package dto

import (
	"broker/core/models"
	"context"
)

type PeerEnvelope struct {
	ctx context.Context

	Meta    Meta            `json:"meta"`
	Event   PeerEventType   `json:"event"`
	Payload *models.Message `json:"payload,omitempty"`
	FromId  string          `json:"fromId"`
	PeerId  string          `json:"peerId"`
}
