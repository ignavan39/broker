package service

import "broker/app/dto"

type WorkspaceService interface {
	Create(payload dto.CreateWorkspacePayload, userId string) (*dto.CreateWorkspaceResponse, error)
}
