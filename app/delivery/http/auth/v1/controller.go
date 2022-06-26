package auth

import (
	"broker/app/dto"
	"broker/app/service"
	"broker/pkg/httpext"
	"encoding/json"
	"errors"
	"fmt"
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
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = payload.Validate()
	if err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := c.authService.SignUp(ctx, payload)
	if err != nil {
		if errors.Is(err, service.EmailCodeNotMatchErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.VerifyCodeExpireErr) {
			httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	ctx := r.Context()
	var payload dto.SignInPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
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

func (c *Controller) SendVerifyCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload dto.SendCodePayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		httpext.AbortJSON(w, fmt.Sprintf("failed decode payload %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := payload.Validate(); err != nil {
		httpext.AbortJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.authService.SendVerifyCode(ctx, payload.Email); err != nil {
		blogger.Errorf("[user/verifyCode] CTX:[%v], ERROR:[%s]", ctx, err.Error())
		httpext.AbortJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpext.EmptyResponse(w, http.StatusOK)
}
