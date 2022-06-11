package middleware

import (
	"broker/app/service"
	"context"
	"net/http"
	"strings"
)

type ContextUseType = string

const ContextUserKey ContextUseType = "user"

func GetUserIdFromContext(ctx context.Context) string {
	userId, ok := ctx.Value(ContextUserKey).(string)
	if !ok {
		return ""
	}

	return userId
}

type AuthGuard struct {
	authService service.AuthService
}

func NewAuthGuard(authService service.AuthService) *AuthGuard {
	return &AuthGuard{
		authService: authService,
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
				claims, ok := ag.authService.Validate(jwtToken)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					return
				}

				ctx := context.WithValue(r.Context(), ContextUserKey, claims.Id)
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		})
	}
}
