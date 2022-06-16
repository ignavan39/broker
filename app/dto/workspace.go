package dto

import (
	"broker/app/models"
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
