package middlewares

import (
	"context"
	"crypto/rsa"
	"msn/internal/infra/jwt"
	"msn/pkg/common/fault"
	"net/http"
	"strings"
)

type AuthKey struct{}

type middleware struct {
	accessKey *rsa.PrivateKey
}

func NewWithAuth(accessKey *rsa.PrivateKey) *middleware {
	return &middleware{
		accessKey: accessKey,
	}
}

func (m *middleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")

		if len(accessToken) == 0 {
			fault.NewHTTPError(w, fault.NewUnauthorized("access token not provided"))
			return
		}

		claims, err := jwt.Verify(m.accessKey, accessToken)
		if err != nil {
			if strings.Contains(err.Error(), "token has expired") {
				fault.NewHTTPError(w, fault.NewUnauthorized("token has expired"))
				return
			}
			fault.NewHTTPError(w, fault.NewUnauthorized("invalid access token"))
			return
		}

		ctx := context.WithValue(r.Context(), AuthKey{}, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
