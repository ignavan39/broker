package auth

import (
	"broker/app/config"
	"broker/app/dto"
	"broker/app/models"
	"broker/app/repository"
	"broker/app/service"
	"broker/pkg/cache"
	"broker/pkg/mailer"
	"broker/pkg/utils"
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService struct {
	signingKey     []byte
	expireDuration time.Duration
	userRepo       repository.UserRepository
	mailer         mailer.Mailer
	cache          cache.Cache[string]
}

func NewAuthService(
	signingKey []byte,
	expireDuration time.Duration,
	userRepo repository.UserRepository,
	cache cache.Cache[string],
	mailer mailer.Mailer,
) *AuthService {
	return &AuthService{
		signingKey:     signingKey,
		expireDuration: expireDuration,
		userRepo:       userRepo,
		cache:          cache,
		mailer:         mailer,
	}
}

func (a *AuthService) SignUp(ctx context.Context, payload dto.SignUpPayload) (*dto.SignResponse, error) {
	user, err := a.userRepo.Create(*payload.Nickname, *payload.Email, utils.CryptString(payload.Password, config.GetConfig().JWT.HashSalt), payload.LastName, payload.FirstName)
	if err != nil {
		return nil, err
	}

	if err := a.SendVerifyCode(ctx, *payload.Email); err != nil {
		return nil, err
	}

	payloadBuilder := dto.NewSignPayloadResponseBuilder().WithUser(*user)

	accessToken, err := a.createToken(user.ID, a.expireDuration)
	if err != nil {
		return nil, err
	}

	payloadBuilder.WithAccessToken(accessToken)

	refreshToken, err := a.createToken(user.ID, time.Duration(24*30))
	if err != nil {
		return nil, err
	}

	payloadBuilder.WithRefreshToken(refreshToken)

	res := payloadBuilder.Exec()

	return &res, nil
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

	payloadBuilder := dto.NewSignPayloadResponseBuilder().WithUser(*user)

	accessToken, err := a.createToken(user.ID, a.expireDuration)
	if err != nil {
		return nil, err
	}

	payloadBuilder.WithAccessToken(accessToken)

	refreshToken, err := a.createToken(user.ID, time.Duration(24*30))
	if err != nil {
		return nil, err
	}

	payloadBuilder.WithRefreshToken(refreshToken)

	res := payloadBuilder.Exec()
	return &res, err
}

func (a *AuthService) SendVerifyCode(ctx context.Context, email string) error {
	code := fmt.Sprintf("%d", utils.GenerateRandomNumber(1000, 9999))

	if err := a.cache.Set(ctx, fmt.Sprintf("%s_%s", getUserPrefix(), email), code); err != nil {
		return err
	}

	_, _, err := a.mailer.SendMail(ctx, fmt.Sprintf("Your verify code :%s", code), "Verify code", email)
	if err != nil {
		return err
	}

	return nil
}
func (a *AuthService) VerifyCode(ctx context.Context, userId string, payload dto.VerifyCodePayload) error {
	email, err := a.userRepo.GetEmailById(userId)
	if err != nil {
		return err
	}

	code, err := a.cache.Get(ctx, fmt.Sprintf("%s_%s", getUserPrefix(), email))
	if err != nil {
		return service.VerifyCodeExpireErr
	}

	if code != &payload.Code {
		return service.EmailCodeNotMatchErr
	}

	return nil
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

	_, err = a.userRepo.GetEmailById(customClaims.ID)
	if err != nil {
		return customClaims, false
	}

	return customClaims, true
}

func getUserPrefix() string {
	return "user"
}
