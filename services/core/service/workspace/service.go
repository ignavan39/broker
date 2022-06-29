package workspace

import (
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"strings"
)

type WorkspaceService struct {
	workspaceRepository repository.WorkspaceRepository
	userRepository      repository.UserRepository
	peerRepository      repository.PeerRepository
}

func NewWorkspaceService(
	workspaceRepository repository.WorkspaceRepository,
	userRepository repository.UserRepository,
	peerRepository repository.PeerRepository,
) *WorkspaceService {
	return &WorkspaceService{
		workspaceRepository: workspaceRepository,
		userRepository:      userRepository,
		peerRepository:      peerRepository,
	}
}

func (s *WorkspaceService) Create(payload dto.CreateWorkspacePayload, userID string) (*dto.CreateWorkspaceResponse, error) {
	workspace, err := s.workspaceRepository.Create(userID, payload.Name, payload.IsPrivate)

	if err != nil {
		return nil, err
	}

	res := dto.CreateWorkspaceResponse(*workspace)
	return &res, nil
}

func (s *WorkspaceService) Delete(userID string, workspaceID string) error {
	accessType, err := s.workspaceRepository.GetAccessByUserId(userID, workspaceID)

	if err != nil {
		return err
	}

	if strings.Compare(*accessType, models.ADMIN) != 0 {
		return service.WorkspaceAccessDeniedErr
	}

	err = s.workspaceRepository.Delete(workspaceID)

	if err != nil {
		return err
	}

	return nil
}

func (s *WorkspaceService) Update(payload dto.UpdateWorkspacePayload, workspaceID string, userID string) (*dto.UpdateWorkspaceResponse, error) {
	accessType, err := s.workspaceRepository.GetAccessByUserId(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	if strings.Compare(*accessType, models.ADMIN) != 0 {
		return nil, service.WorkspaceAccessDeniedErr
	}

	_, err = s.workspaceRepository.GetWorkspaceByUserId(userID, workspaceID)

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

func (s *WorkspaceService) GetWorkspaceInfo(userID string, workspaceID string) (*dto.GetWorkspaceInfoResponse, error) {
	_, err := s.workspaceRepository.GetAccessByUserId(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	workspace, err := s.workspaceRepository.GetWorkspaceByUserId(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	peers, err := s.peerRepository.GetMany(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	usersCount, err := s.workspaceRepository.GetWorkspaceUsersCount(workspaceID)

	if err != nil {
		return nil, err
	}

	dtoPeers := make([]dto.PeerResponse, 0)

	for _, peer := range peers {
		dtoPeer := dto.PeerResponse{
			Id:   peer.ID,
			Name: peer.Name,
		}

		dtoPeers = append(dtoPeers, dtoPeer)
	}

	return &dto.GetWorkspaceInfoResponse{
		Name:       workspace.Name,
		Peers:      dtoPeers,
		UsersCount: usersCount,
	}, nil
}
