package token

import (
	"msn/pkg/common/dto"
	"msn/pkg/utils/uid"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User dto.UserResponse `json:"user"`
	jwt.RegisteredClaims
}

func NewClaims(user dto.UserResponse, duration time.Duration) (*Claims, error) {
	jti := uid.New("jti")

	return &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "user-service",
			ID:        jti,
		},
	}, nil
}
