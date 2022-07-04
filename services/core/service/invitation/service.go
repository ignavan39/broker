package invitation

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"broker/pkg/scheduler"
	"context"
	"strings"
	"time"

	blogger "github.com/sirupsen/logrus"
)

type InvitationService struct {
	workspaceRepository  repository.WorkspaceRepository
	invitationRepository repository.InvitationRepository
	scheduler            scheduler.Scheduler
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

func (s *InvitationService) StartScheduler(ctx context.Context) error {
	duration := config.GetConfig().Invitation.InvitationExpireDuration

	scheduler := scheduler.NewScheduler(duration, ctx, func(ctx context.Context) error {
		err := s.clearExpiredInvitations(duration)

		if err != nil {
			blogger.Panic(err)
			return err
		}

		blogger.Info("Invitations cleared successfully")

		return nil
	})

	go scheduler.Start()

	err := <-scheduler.Error()

	if err != nil {
		return err
	}

	return nil
}

func (s *InvitationService) CreateInvitation(payload dto.SendInvitationPayload,
	senderID string,
	workspaceID string) (*dto.SendInvitationResponse, error) {
	accessType, err := s.workspaceRepository.GetAccessByUserId(senderID, workspaceID)

	if err != nil {
		return nil, err
	}

	if strings.Compare(*accessType, models.USER) == 0 {
		return nil, service.WorkspaceAccessDeniedErr
	}

	invitation, err := s.invitationRepository.SendInvitation(senderID, workspaceID, payload.RecipientEmail)

	if err != nil {
		return nil, err
	}

	return &dto.SendInvitationResponse{
		ID:             invitation.ID,
		CreatedAt:      invitation.CreatedAt,
		Sender:         invitation.Sender,
		RecipientEmail: invitation.RecipientEmail,
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
	invitation, err := s.invitationRepository.CancelInvitation(senderID, invitationID)

	if err != nil {
		return nil, err
	}

	return &dto.CancelInvitationResponse{
		ID:             invitation.ID,
		CreatedAt:      invitation.CreatedAt,
		RecipientEmail: invitation.RecipientEmail,
		Sender:         invitation.Sender,
		WorkspaceID:    invitation.WorkspaceID,
		Status:         invitation.Status,
		SystemStatus:   invitation.SystemStatus,
		Code:           invitation.Code,
	}, nil
}

func (s *InvitationService) AcceptInvitation(payload dto.AcceptInvitationPayload, userID string) error {
	err := s.invitationRepository.AcceptInvitation(userID, payload.Code)

	if err != nil {
		return err
	}

	return nil
}

func (s *InvitationService) clearExpiredInvitations(duration time.Duration) error {
	err := s.invitationRepository.DeleteExpiredInvitations(duration)

	if err != nil {
		return err
	}

	return nil
}
