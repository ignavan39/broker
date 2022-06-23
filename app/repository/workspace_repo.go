package repository

import "broker/app/models"

type WorkspaceRepository interface {
	Create(email string, name string, isPrivate bool) (*models.Workspace, error)
	Delete(workspaceID string) error
	Update(workspaceID string, name *string, isPrivate *bool) (*models.Workspace, error)
	GetManyByUserId(id string) ([]models.Workspace, error)
	GetWorkspaceByUserId(userID string, workspaceID string) (*models.Workspace, error)
	GetAccessByUserId(userID string, workspaceID string) (*string, error)
	GetWorkspaceUsersCount(userID string, workspaceID string) (int, error)
}
