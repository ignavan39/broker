package invitation

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"broker/core/service/invitation/publisher"
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
	userRepository       repository.UserRepository
	connectionService    service.ConnectionService
	mailer               mailer.Mailer
	publisher            *publisher.Publisher
	schedulerErrors      chan error
}

func NewInvitationService(
	workspaceRepository repository.WorkspaceRepository,
	invitationRepository repository.InvitationRepository,
	userRepository repository.UserRepository,
	connectionService service.ConnectionService,
	mailer mailer.Mailer,
	publisher *publisher.Publisher,
) *InvitationService {
	return &InvitationService{
		workspaceRepository:  workspaceRepository,
		invitationRepository: invitationRepository,
		userRepository:       userRepository,
		connectionService:    connectionService,
		mailer:               mailer,
		publisher:            publisher,
		schedulerErrors:      make(chan error),
	}
}

func (s *InvitationService) StartScheduler(ctx context.Context) {
	duration := config.GetConfig().Invitation.InvitationExpireDuration

	invitationScheduler := scheduler.NewScheduler(duration, func(ctx context.Context) error {
		err := s.clearExpiredInvitations(duration)

		if err != nil {
			logger.Logger.Errorf("%s", err)
			return err
		}

		logger.Logger.Info("Invitations cleared successfully")

		return nil
	})

	go invitationScheduler.Start(ctx)
	go func() {
		err := <-invitationScheduler.Error()
		s.schedulerErrors <- err
	}()

	deleteExpiredQueuesScheduler := scheduler.NewScheduler(time.Duration(time.Second*5), func(ctx context.Context) error {
		keys, err := s.publisher.RemoveDeadQueues(time.Now())
		if err != nil {
			return err
		}
		for _, key := range keys {
			s.connectionService.Remove(key)
		}
		return nil
	})

	go deleteExpiredQueuesScheduler.Start(ctx)
	go func() {
		err := <-deleteExpiredQueuesScheduler.Error()
		logger.Logger.Error("[InvitationPublisher] error %s", err)
	}()
}

func (s *InvitationService) ReadError() error {
	return <-s.schedulerErrors
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

	var access string

	if workspace.IsPrivate {
		access = "private"
	} else {
		access = "public"
	}

	_, _, err = s.mailer.SendMail(ctx,
		fmt.Sprintf("You have been invited to %s workspace '%s'. Follow this link to accept invitation:\n https://%s/invitations/%s",
			access, workspace.Name, config.GetConfig().Frontend.Host, code),
		"Invitation to workspace", payload.RecipientEmail)

	if err != nil {
		return nil, err
	}

	recipient, err := s.userRepository.GetOneByEmail(payload.RecipientEmail)

	if err == nil {
		if err := s.publisher.Publish(recipient.ID, *invitation); err != nil {
			return nil, err
		}
	}

	return &dto.SendInvitationResponse{
		ID:             invitation.ID,
		CreatedAt:      invitation.CreatedAt,
		Sender:         invitation.Sender,
		RecipientEmail: invitation.RecipientEmail,
		Workspace:    invitation.Workspace,
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
		Workspace:    invitation.Workspace,
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

func (s *InvitationService) RejectInvitation(payload dto.RejectInvitationPayload, userID string) error {
	err := s.invitationRepository.RejectInvitation(userID, payload.Code)

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

func (s *InvitationService) Connect(ctx context.Context, userID string) (*dto.ConnectInvitationResponse, error) {
	queue, err := s.publisher.CreateConnection(ctx, userID)

	if err != nil {
		return nil, err
	}

	ch := make(chan int)
	s.connectionService.Add(userID, queue.Consume.QueueName, ch)

	go func() {
		_, ok := <-ch
		if ok {
			s.publisher.SetLastUpdateTimeByUserId(userID, time.Now())
		} else {
			return
		}
	}()

	return queue, nil
}
