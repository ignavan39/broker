package auth

import (
	"broker/app/config"
	"broker/app/dto"
	"broker/app/models"
	"broker/app/repository"
	"broker/app/service"
	"broker/pkg/utils"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService struct {
	signingKey     []byte
	expireDuration time.Duration
	userRepo       repository.UserRepository
}

func NewAuthService(
	signingKey []byte,
	expireDuration time.Duration,
	userRepo repository.UserRepository,
) *AuthService {
	return &AuthService{
		signingKey:     signingKey,
		expireDuration: expireDuration,
		userRepo:       userRepo,
	}
}

func (a *AuthService) SignUp(payload dto.SignUpPayload) (*dto.SignResponse, error) {
	user, err := a.userRepo.Create(*payload.Nickname, *payload.Email, utils.CryptString(payload.Password, config.GetConfig().JWT.HashSalt), payload.LastName, payload.FirstName)
	if err != nil {
		return nil, err
	}

	auth, err := a.refresh(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.SignResponse{
		Auth: auth,
		User: *user,
	}, nil
}

func (a *AuthService) SignIn(payload dto.SignInPayload) (*dto.SignResponse, error) {
	var user *models.User
	var err error

	if payload.Email != nil {
		user, err = a.userRepo.GetOneByEmail(*payload.Email)
	} else {
		user, err = a.userRepo.GetOneByNickname(*payload.Nickname)
	}

	if err != nil {
		return nil, err
	}

	if utils.CryptString(payload.Password, config.GetConfig().JWT.HashSalt) != user.Password {
		return nil, service.PasswordNotMatch
	}

	auth, err := a.refresh(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.SignResponse{
		Auth: auth,
		User: *user,
	}, nil
}

func (a *AuthService) refresh(id string) (map[string]string, error) {
	auth := map[string]string{}
	accessToken, err := a.createToken(id, a.expireDuration)
	if err != nil {
		return auth, err
	}
	auth["accessToken"] = accessToken
	refreshToken, err := a.createToken(id, time.Duration(24*30))
	if err != nil {
		return auth, err
	}
	auth["refreshToken"] = refreshToken
	return auth, nil
}

func (a *AuthService) createToken(id string, expireAt time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &service.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(expireAt)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Id: id,
	})

	return token.SignedString(a.signingKey)
}

func (a *AuthService) Validate(jwtToken string) (*service.Claims, bool) {
	customClaims := &service.Claims{}

	token, err := jwt.ParseWithClaims(jwtToken, customClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.signingKey), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}

	return customClaims, true
}
