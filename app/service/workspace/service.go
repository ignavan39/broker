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

func (s *WorkspaceService) Delete(userID string, workspaceID string) error {
	_, err := s.workspaceRepository.GetWorkspaceByUserId(userID, workspaceID)

	if err != nil {
		return err
	}

	err = s.workspaceRepository.Delete(workspaceID)

	if err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) Update(payload dto.UpdateWorkspacePayload, workspaceID string, userID string) (*dto.UpdateWorkspaceResponse, error) {
	_, err := s.workspaceRepository.GetWorkspaceByUserId(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	workspace, err := s.workspaceRepository.Update(workspaceID, payload.Name, payload.IsPrivate)

	if err != nil {
		return nil, err
	}

	return &dto.UpdateWorkspaceResponse{
		ID:        workspace.ID,
		Name:      workspace.Name,
		CreatedAt: workspace.CreatedAt,
		IsPrivate: workspace.IsPrivate,
	}, nil
}

func (s *WorkspaceService) GetManyByUserID(userId string) (*dto.GetManyByUserResponse, error) {
	workspaces, err := s.workspaceRepository.GetManyByUserId(userId)

	if err != nil {
		return nil, err
	}

	return &dto.GetManyByUserResponse{
		Workspaces: workspaces,
	}, nil
}
