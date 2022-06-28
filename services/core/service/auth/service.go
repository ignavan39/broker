package auth

import (
	"broker/core/config"
	"broker/core/dto"
	"broker/core/models"
	"broker/core/repository"
	"broker/core/service"
	"broker/pkg/cache"
	"broker/pkg/mailer"
	"broker/pkg/utils"
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
)

type AuthService struct {
	signingKey           []byte
	expireDuration       time.Duration
	userRepository       repository.UserRepository
	invitationRepository repository.InvitationRepository
	mailer               mailer.Mailer
	cache                cache.Cache[int]
}

func NewAuthService(
	signingKey []byte,
	expireDuration time.Duration,
	userRepository repository.UserRepository,
	invitationRepository repository.InvitationRepository,
	cache cache.Cache[int],
	mailer mailer.Mailer,
) *AuthService {
	return &AuthService{
		signingKey:           signingKey,
		expireDuration:       expireDuration,
		userRepository:       userRepository,
		invitationRepository: invitationRepository,
		cache:                cache,
		mailer:               mailer,
	}
}

func (a *AuthService) SignUp(ctx context.Context, payload dto.SignUpPayload) (*dto.SignResponse, error) {
	if err := a.verifyCode(ctx, *payload.Email, payload.Code); err != nil {
		return nil, err
	}

	user, err := a.userRepository.Create(*payload.Nickname, *payload.Email, utils.CryptString(payload.Password, config.GetConfig().JWT.HashSalt), payload.LastName, payload.FirstName)
	if err != nil {
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

	if err := a.invitationRepository.CheckInvites(user.ID, user.Email); err != nil {
		return nil, err
	}

	return &res, nil
}

func (a *AuthService) SignIn(payload dto.SignInPayload) (*dto.SignResponse, error) {
	var user *models.User
	var err error

	if payload.Email != nil {
		user, err = a.userRepository.GetOneByEmail(*payload.Email)
	} else {
		user, err = a.userRepository.GetOneByNickname(*payload.Nickname)
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
	code := utils.GenerateRandomNumber(10000, 99999)

	if err := a.cache.Set(ctx, fmt.Sprintf("%s_%s", getUserPrefix(), email), code); err != nil {
		return err
	}

	_, _, err := a.mailer.SendMail(ctx, fmt.Sprintf("Your verify code :%d", code), "Verify code", email)
	if err != nil {
		return err
	}

	return nil
}
func (a *AuthService) verifyCode(ctx context.Context, email string, code int) error {
	codeFromCache, err := a.cache.Get(ctx, fmt.Sprintf("%s_%s", getUserPrefix(), email))
	if err != nil {
		return service.VerifyCodeExpireErr
	}

	if code != *codeFromCache {
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

	_, err = a.userRepository.GetEmailById(customClaims.Id)
	if err != nil {
		return customClaims, false
	}

	return customClaims, true
}

func getUserPrefix() string {
	return "user"
}
