package peer

import (
	"broker/app/delivery/http/middleware"
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
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
	var payload dto.CreateWorkspaceConnectionPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	ctx := r.Context()

	if err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
	}

	userId := middleware.GetUserIdFromContext(ctx)

	res, err := c.peerService.CreateConnection(ctx, userId, payload)

	if err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.JSON(w, res, http.StatusCreated)
}
