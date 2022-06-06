package api

import (
	"broker/app/types"
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"

	blogger "github.com/sirupsen/logrus"
)

type AuthGuard struct {
	signingKey string
}

func NewAuthGuard(signingKey string) *AuthGuard {
	return &AuthGuard{
		signingKey: signingKey,
	}
}

func (ag *AuthGuard) Next() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
				return
			} else {
				jwtToken := authHeader[1]
				customClaims := &types.Claims{}

				token, err := jwt.ParseWithClaims(jwtToken, customClaims, func(token *jwt.Token) (interface{}, error) {
					return []byte(ag.signingKey), nil
				})
				if err != nil || !token.Valid {
					blogger.Error("[AuthGuard] Error :%s", err.Error())
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}

				ctx := context.WithValue(r.Context(), types.ContextUserKey, customClaims.Id)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
