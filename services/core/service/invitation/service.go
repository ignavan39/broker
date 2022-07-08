package invitation

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"broker/pkg/mailer"
	"broker/pkg/scheduler"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	blogger "github.com/sirupsen/logrus"
)

type InvitationService struct {
	workspaceRepository  repository.WorkspaceRepository
	invitationRepository repository.InvitationRepository
	scheduler            scheduler.Scheduler
	mailer               mailer.Mailer
	runes                []rune
}

func NewInvitationService(
	workspaceRepository repository.WorkspaceRepository,
	invitationRepository repository.InvitationRepository,
	mailer mailer.Mailer,
) *InvitationService {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-")

	return &InvitationService{
		workspaceRepository:  workspaceRepository,
		invitationRepository: invitationRepository,
		mailer:               mailer,
		runes:                runes,
	}
}

func (s *InvitationService) StartScheduler(ctx context.Context) {
	duration := config.GetConfig().Invitation.InvitationExpireDuration

	scheduler := scheduler.NewScheduler(duration, func(ctx context.Context) error {
		err := s.clearExpiredInvitations(duration)

		if err != nil {
			blogger.Errorf("%s", err)
			return err
		}

		blogger.Info("Invitations cleared successfully")

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

	_, _, err = s.mailer.SendMail(ctx,
		fmt.Sprintf("You have been invited to workspace '%s'. Follow this link to accept invitation: https://kind-of-link-i-guess.com/invitation/%s",
			workspace.Name, s.generateBigString()),
		"Invitation to workspace", payload.RecipientEmail)

	if err != nil {
		return nil, err
	}

	invitation, err := s.invitationRepository.CreateInvitation(senderID, workspaceID, payload.RecipientEmail)

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

func (s *InvitationService) generateBigString() string {
	rand.Seed(time.Now().UnixNano())

	link := make([]rune, 100)

	for i := range link {
		link[i] = s.runes[rand.Intn(len(s.runes))]
	}

	return string(link)
}
