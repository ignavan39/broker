package dto

import (
	"errors"
	"time"
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

type CreateWorkspaceResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"isPrivate"`
	CreatedAt time.Time `json:"createdAt"`
}
