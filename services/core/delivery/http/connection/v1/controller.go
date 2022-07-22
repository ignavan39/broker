package connection

import (
	"broker/core/delivery/http/middleware"
	"broker/core/service"
	"broker/pkg/httpext"
	"net/http"
)

type Controller struct {
	connectionService service.ConnectionService
}

func NewController(
	connectionService service.ConnectionService,
) *Controller {
	return &Controller{
		connectionService: connectionService,
	}
}

func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := middleware.GetUserIdFromContext(ctx)

	c.connectionService.Ping(userID)

	httpext.EmptyResponse(w, http.StatusOK)
}
