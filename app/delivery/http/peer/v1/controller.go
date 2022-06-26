package peer

import (
	"broker/app/delivery/http/middleware"
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
	"fmt"
	"net/http"
)

type Controller struct {
	peerService service.PeerService
}

func NewController(peerService service.PeerService) *Controller {
	return &Controller{
		peerService: peerService,
	}
}

func (c *Controller) CreateConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload dto.CreateWorkspaceConnectionPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
	}

	userId := middleware.GetUserIdFromContext(ctx)

	res, err := c.peerService.CreateConnection(ctx, userId, payload)

	if err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusCreated)
}
