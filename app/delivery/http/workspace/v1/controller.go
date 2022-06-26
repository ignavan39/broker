package workspace

import (
	"broker/app/delivery/http/middleware"
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	blogger "github.com/sirupsen/logrus"
)

type Controller struct {
	workspaceService service.WorkspaceService
}

func NewController(
	workspaceService service.WorkspaceService,
) *Controller {
	return &Controller{
		workspaceService: workspaceService,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload dto.CreateWorkspacePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := middleware.GetUserIdFromContext(ctx)

	res, err := c.workspaceService.Create(payload, userId)

	if err != nil {
		if errors.Is(err, service.DuplicateWorkspaceEmailErr) ||
			errors.Is(err, service.DuplicateWorkspaceErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.UserNotFoundErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusNotFound)
			return
		}
		blogger.Errorf("[Create] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) GetManyByUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := middleware.GetUserIdFromContext(ctx)

	res, err := c.workspaceService.GetManyByUserID(userID)
	if err != nil {
		blogger.Errorf("[WorkspaceController][GetManyByUser] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workspaceID := chi.URLParam(r, "workspaceID")

	if len(workspaceID) == 0 {
		httpext.AbortJSON(w, service.EmptyUrlParamsErr.Error(), http.StatusBadRequest)
	}

	userID := middleware.GetUserIdFromContext(ctx)

	err := c.workspaceService.Delete(userID, workspaceID)

	if err != nil {
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		blogger.Errorf("[WorkspaceController][Delete] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.EmptyResponse(w, http.StatusOK)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload dto.UpdateWorkspacePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
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

	res, err := c.workspaceService.Update(payload, workspaceID, userID)

	if err != nil {
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		blogger.Errorf("[WorkspaceController][Update] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}

func (c *Controller) GetWorkspaceInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workspaceID := chi.URLParam(r, "workspaceID")

	if len(workspaceID) == 0 {
		httpext.AbortJSON(w, service.EmptyUrlParamsErr.Error(), http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserIdFromContext(ctx)

	response, err := c.workspaceService.GetWorkspaceInfo(userID, workspaceID)

	if err != nil {
		if errors.Is(err, service.WorkspaceAccessDeniedErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusForbidden)
			return
		}
		blogger.Errorf("[GetWorkspaceInfo] Error: %s", err)
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, response, http.StatusOK)
}
