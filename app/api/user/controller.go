package user

import (
	"broker/app/config"
	"broker/app/repository"
	"broker/app/services"
	"broker/pkg/httpext"
	"broker/pkg/utils"
	"encoding/json"
	"net/http"

	blogger "github.com/sirupsen/logrus"
)

type Controller struct {
	authService services.AuthService
	repo        repository.UserRepository
}

func NewController(
	authService services.AuthService,
	repo repository.UserRepository) *Controller {
	return &Controller{
		authService: authService,
		repo:        repo,
	}
}

func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var payload SignUpPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	ctx := r.Context()

	if err != nil {
		httpext.JSON(w, httpext.CommonError{
			Error: "failed decode payload",
			Code:  http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}

	err = payload.Validate()
	if err != nil {
		httpext.JSON(w, httpext.CommonError{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}

	user, err := c.repo.Create(payload.Email, utils.CryptString(payload.Password, config.GetConfig().JWT.HashSalt), payload.LastName, payload.FirstName)
	if err != nil {
		blogger.Errorf("[user/sign-up] CTX:[%v], ERROR:[%s]", ctx, err.Error())
		httpext.JSON(w, httpext.CommonError{
			Error: "user already exists",
			Code:  http.StatusBadRequest,
		}, http.StatusBadRequest)
		return
	}

	accessToken, err := c.authService.CreateToken(ctx, user.Id)
	if err != nil {
		httpext.JSON(w, httpext.CommonError{
			Error: "failed created access token",
			Code:  http.StatusInternalServerError,
		}, http.StatusInternalServerError)
		return
	}

	response := SignResponse{*user, accessToken}
	httpext.JSON(w, response, http.StatusCreated)
}
