package invitation

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"broker/pkg/logger"
	"broker/pkg/mailer"
	"broker/pkg/scheduler"
	"broker/pkg/utils"
	"context"
	"fmt"
	"strings"
	"time"
)

type InvitationService struct {
	workspaceRepository  repository.WorkspaceRepository
	invitationRepository repository.InvitationRepository
	scheduler            scheduler.Scheduler
	mailer               mailer.Mailer
}

func NewInvitationService(
	workspaceRepository repository.WorkspaceRepository,
	invitationRepository repository.InvitationRepository,
	mailer mailer.Mailer,
) *InvitationService {
	return &InvitationService{
		workspaceRepository:  workspaceRepository,
		invitationRepository: invitationRepository,
		mailer:               mailer,
	}
}

func (s *InvitationService) StartScheduler(ctx context.Context) {
	duration := config.GetConfig().Invitation.InvitationExpireDuration

	scheduler := scheduler.NewScheduler(duration, func(ctx context.Context) error {
		err := s.clearExpiredInvitations(duration)

		if err != nil {
			logger.Logger.Errorf("%s", err)
			return err
		}

		logger.Logger.Info("Invitations cleared successfully")

		return nil
	})

	go scheduler.Start(ctx)

	s.scheduler = *scheduler
}

func (s *InvitationService) ReadError() error {
	return <-s.scheduler.Error()
}

func (s *InvitationService) SendInvitation(ctx context.Context, payload dto.SendInvitationPayload,
	senderID string,
	workspaceID string) (*dto.SendInvitationResponse, error) {
	accessType, err := s.workspaceRepository.GetAccessByUserId(senderID, workspaceID)

	if err != nil {
		return nil, err
	}

	if strings.Compare(*accessType, models.USER) == 0 {
		return nil, service.WorkspaceAccessDeniedErr
	}

	workspace, err := s.workspaceRepository.GetWorkspaceByUserId(senderID, workspaceID)

	if err != nil {
		return nil, err
	}

	code := utils.GenerateBigString(100)

	invitation, err := s.invitationRepository.CreateInvitation(senderID, workspaceID, payload.RecipientEmail, code)

	if err != nil {
		return nil, err
	}

	_, _, err = s.mailer.SendMail(ctx,
		fmt.Sprintf("You have been invited to workspace '%s'. Follow this link to accept invitation: https://localhost:3000/invitations/%s",
			workspace.Name, code),
		"Invitation to workspace", payload.RecipientEmail)

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
