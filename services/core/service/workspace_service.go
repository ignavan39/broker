package service

import (
	"broker/core/dto"
)

type WorkspaceService interface {
	Create(payload dto.CreateWorkspacePayload, userID string) (*dto.CreateWorkspaceResponse, error)
	Delete(usedID string, workspaceID string) error
	Update(payload dto.UpdateWorkspacePayload, workspaceID string, userID string) (*dto.UpdateWorkspaceResponse, error)
	GetManyByUserID(userID string) (*dto.GetManyByUserResponse, error)
	GetWorkspaceInfo(userID string, workspaceID string) (*dto.GetWorkspaceInfoResponse, error)
	ChangeUserRole(payload dto.ChangeUserRoleWorkspacePayload, adminID string, userID string, workspaceID string) error
	BanUser(userID string, bannedUserID string, workspaceID string) error
}
