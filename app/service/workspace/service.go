package workspace

import (
	"broker/app/dto"
	"broker/app/repository"
)

type WorkspaceService struct {
	workspaceRepository repository.WorkspaceRepository
	userRepository      repository.UserRepository
}

func NewWorkspaceService(
	workspaceRepository repository.WorkspaceRepository,
	userRepository repository.UserRepository,
) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository: workspaceRepository,
		userRepository:      userRepository,
	}
}

func (s *WorkspaceService) Create(payload dto.CreateWorkspacePayload, userId string) (*dto.CreateWorkspaceResponse, error) {
	userEmail, err := s.userRepository.GetEmailById(userId)

	if err != nil {
		return nil, err
	}

	workspace, err := s.workspaceRepository.Create(userEmail, payload.Name, payload.IsPrivate)

	if err != nil {
		return nil, err
	}

	return &dto.CreateWorkspaceResponse{
		Id:        workspace.Id,
		Name:      workspace.Name,
		IsPrivate: workspace.IsPrivate,
		CreatedAt: workspace.CreatedAt,
	}, nil
}
