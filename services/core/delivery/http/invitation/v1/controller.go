package invitation

import (
	"broker/core/delivery/http/middleware"
	"broker/core/dto"
	"broker/core/service"
	"broker/pkg/httpext"
	"broker/pkg/logger"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
)

type Controller struct {
	invitationService service.InvitationService
}

func NewController(invitationService service.InvitationService) *Controller {
	return &Controller{
		invitationService: invitationService,
	}
}

func (c *Controller) SendInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload dto.SendInvitationPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspaceID := chi.URLParam(r, "workspaceID")

	if len(workspaceID) == 0 {
		httpext.AbortJSON(w, service.EmptyUrlParamsErr.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	res, err := c.invitationService.SendInvitation(ctx, payload, userID, workspaceID)

	if err != nil {
		if errors.Is(err, service.DuplicateInvitationErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Errorf("[InvitationController][SendInvitation] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) GetInvitations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workspaceID := chi.URLParam(r, "workspaceID")

	if len(workspaceID) == 0 {
		httpext.AbortJSON(w, service.EmptyUrlParamsErr.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	res, err := c.invitationService.GetInvitationsByWorkspaceID(userID, workspaceID)

	if err != nil {
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		logger.Logger.Errorf("[InvitationController][GetInvitations] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) CancelInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	invitationID := chi.URLParam(r, "invitationID")

	if len(invitationID) == 0 {
		httpext.AbortJSON(w, service.EmptyUrlParamsErr.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	res, err := c.invitationService.CancelInvitation(userID, invitationID)

	if err != nil {
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		if errors.Is(err, service.InvitationNotFoundErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Errorf("[InvitationController][CancelInvitation] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload dto.AcceptInvitationPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	err := c.invitationService.AcceptInvitation(payload, userID)

	if err != nil {
		if errors.Is(err, service.InviteAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		if errors.Is(err, service.InvitationNotFoundErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Errorf("[InvitationController][AcceptInvitation] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.EmptyResponse(w, http.StatusOK)
}

func (c *Controller) Connect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := middleware.GetUserIdFromContext(ctx)

	res, err := c.invitationService.Connect(ctx, userID)

	if err != nil {
		logger.Logger.Errorf("[InvitationController][Connect] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) RejectInvitation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload dto.RejectInvitationPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	err := c.invitationService.RejectInvitation(payload, userID)

	if err != nil {
		if errors.Is(err, service.InvitationNotFoundErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Logger.Errorf("[InvitationController][RejectInvitation] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.EmptyResponse(w, http.StatusOK)
}
