package workspace

import (
	"broker/app/delivery/http/middleware"
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
	"errors"
	"net/http"
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
	var payload dto.CreateWorkspacePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	ctx := r.Context()

	if err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
		return
	}

	err = payload.Validate()
	if err != nil {
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
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}
