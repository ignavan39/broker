package invitation

import (
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"strings"
)

type InvitationService struct {
	workspaceRepository  repository.WorkspaceRepository
	invitationRepository repository.InvitationRepository
}

func NewInvitationService(
	workspaceRepository repository.WorkspaceRepository,
	invitationRepository repository.InvitationRepository,
) *InvitationService {
	return &InvitationService{
		workspaceRepository:  workspaceRepository,
		invitationRepository: invitationRepository,
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

	invitation, err := s.invitationRepository.SendInvitation(senderID, workspaceID, payload.RicipientEmail)

	if err != nil {
		return nil, err
	}

	return &dto.SendInvitationResponse{
		ID:             invitation.ID,
		CreatedAt:      invitation.CreatedAt,
		Sender:         invitation.Sender,
		RicipientEmail: invitation.RicipientEmail,
		WorkspaceID:    invitation.WorkspaceID,
		Status:         invitation.Status,
		SystemStatus:   invitation.SystemStatus,
		Code:           invitation.Code,
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

	invitations, err := s.invitationRepository.GetInvitationsByWorkspaceID(userID, workspaceID)

	if err != nil {
		return nil, err
	}

	return &dto.GetInvitationsByWorkspaceResponse{
		Invitations: invitations,
	}, nil
}

func (s *InvitationService) CancelInvitation(senderID string, invitationID string) (*dto.CancelInvitationResponse, error) {

	invitation, err := s.invitationRepository.CancelInvitation(invitationID)

	if err != nil {
		return nil, err
	}

	return &dto.CancelInvitationResponse{
		ID:             invitation.ID,
		CreatedAt:      invitation.CreatedAt,
		RicipientEmail: invitation.RicipientEmail,
		Sender:         invitation.Sender,
		WorkspaceID:    invitation.WorkspaceID,
		Status:         invitation.Status,
		SystemStatus:   invitation.SystemStatus,
		Code:           invitation.Code,
	}, nil
}

func (s *InvitationService) AcceptInvitation(userID string, code string) error {
	err := s.invitationRepository.AcceptInvitation(userID, code)

	if err != nil {
		return err
	}

	return nil
}
