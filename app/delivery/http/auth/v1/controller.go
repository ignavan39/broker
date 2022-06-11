package auth

import (
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
	"errors"
	"net/http"

	blogger "github.com/sirupsen/logrus"
)

type Controller struct {
	authService service.AuthService
}

func NewController(
	authService service.AuthService,
) *Controller {
	return &Controller{
		authService: authService,
	}
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignUpPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	ctx := r.Context()

	if err != nil {
		httpext.AbortJSON(w, "failed decode payload", http.StatusBadRequest)
	}

	err = payload.Validate()
	if err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
	}

	res, err := c.authService.SignUp(payload)

	if err != nil {
		if errors.Is(err, service.DuplicateUserErr) {
			httpext.AbortJSON(w, "user already exists", http.StatusBadRequest)
			return
		} else {
			blogger.Errorf("[user/sign-up] CTX:[%v], ERROR:[%s]", ctx, err.Error())
			httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	httpext.JSON(w, res, http.StatusCreated)
}

func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignInPayload
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

	res, err := c.authService.SignIn(payload)

	if err != nil {
		var code int
		if errors.Is(err, service.UserNotFoundErr) {
			code = http.StatusNotFound
		} else if errors.Is(err, service.PasswordNotMatch) {
			code = http.StatusBadRequest
		} else {
			blogger.Errorf("[user/sign-up] CTX:[%v], ERROR:[%s]", ctx, err.Error())
			code = http.StatusInternalServerError
		}
		httpext.AbortJSON(w, err.Error(), code)
		return
	}

	httpext.JSON(w, res, http.StatusOK)
}
