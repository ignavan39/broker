package invitation

import (
	"broker/app/dto"
	"broker/app/models"
	"broker/app/repository"
	"broker/app/service"
	"strings"
)

type InvitationService struct {
	userRepository      repository.UserRepository
	workspaceRepository repository.WorkspaceRepository
}

func NewInvitationService(
	userRepository repository.UserRepository,
	workspaceRepository repository.WorkspaceRepository,
) *InvitationService {
	return &InvitationService{
		userRepository:      userRepository,
		workspaceRepository: workspaceRepository,
	}
}

func (s *InvitationService) SendInvitation(payload dto.SendInvitationPayload,
	senderID string,
	workspaceID string) (*dto.SendInvitationResponse, error) {
	accessType, err := s.workspaceRepository.GetAccessByUserId(senderID, workspaceID)

	if err != nil {
		return nil, err
	}

	if strings.Compare(*accessType, models.USER) == 0 {
		return nil, service.WorkspaceAccessDeniedErr
	}

	invitation, err := s.userRepository.SendInvitation(senderID, workspaceID, *payload.RicipientEmail)

	if err != nil {
		return nil, err
	}

	return &dto.SendInvitationResponse{
		ID:             invitation.ID,
		SenderID:       invitation.SenderID,
		RicipientEmail: invitation.RicipientEmail,
		WorkspaceID:    invitation.WorkspaceID,
		Status:         invitation.Status,
	}, nil
}

func (s *InvitationService) GetInvitationsByWorkspaceID(userID string, workspaceID string) (*dto.GetInvitationsByWorkspaceResponse, error) {
	accessType, err := s.workspaceRepository.GetAccessByUserId(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	if strings.Compare(*accessType, models.USER) == 0 {
		return nil, service.WorkspaceAccessDeniedErr
	}

	invitations, err := s.userRepository.GetInvitationsByWorkspaceID(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	return &dto.GetInvitationsByWorkspaceResponse{
		Invitations: invitations,
	}, nil
}