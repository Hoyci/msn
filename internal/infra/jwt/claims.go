package jwt

import (
	"crypto/rsa"
	"fmt"
	"msn/pkg/common/dto"
	"msn/pkg/utils/uid"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User *dto.EnrichedUserResponse `json:"user"`
	jwt.RegisteredClaims
}

func NewClaims(user *dto.EnrichedUserResponse, duration time.Duration) (*Claims, error) {
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

func Verify(secretKey *rsa.PrivateKey, v string) (*Claims, error) {
	if strings.TrimSpace(v) == "" {
		return nil, fmt.Errorf("invalid token")
	}

	v = strings.TrimPrefix(v, "Bearer ")

	keyFunc := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return &secretKey.PublicKey, nil
	}

	token, err := jwt.ParseWithClaims(v, &Claims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
