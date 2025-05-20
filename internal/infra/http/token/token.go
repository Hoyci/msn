package token

import (
	"crypto/rsa"
	"fmt"
	"msn/pkg/common/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secretKey *rsa.PrivateKey, user dto.UserResponse, duration time.Duration) (string, *Claims, error) {
	claims, err := NewClaims(user, duration)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create claims: %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", claims, fmt.Errorf("failed to sign token: %w", err)
	}

	return token, claims, nil
}
