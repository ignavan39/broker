package repository

import "broker/app/models"

type WorkspaceRepository interface {
	Create(email string, name string, isPrivate bool) (*models.Workspace, error)
}
