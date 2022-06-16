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

	res := dto.CreateWorkspaceResponse(*workspace)
	return &res, nil
}

func (s *WorkspaceService) GetManyByUserId(userId string) (*dto.GetManyByUserResponse, error) {
	email, err := s.userRepository.GetEmailById(userId)

	if err != nil {
		return nil, err
	}

	workspaces, err := s.workspaceRepository.GetManyByUserEmail(email)

	if err != nil {
		return nil, err
	}

	return &dto.GetManyByUserResponse{
		Workspaces: workspaces,
	}, nil
}
