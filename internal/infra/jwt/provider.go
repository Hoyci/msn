package jwt

import (
	"crypto/rsa"
	"msn/pkg/common/dto"
	"time"
)

const (
	AccessTokenDuration  = 2 * time.Minute
	RefreshTokenDuration = 720 * time.Hour
)

type JWTProvider struct {
	accessKey  *rsa.PrivateKey
	refreshKey *rsa.PrivateKey
}

func NewProvider(accessKey, refreshKey *rsa.PrivateKey) *JWTProvider {
	return &JWTProvider{
		accessKey:  accessKey,
		refreshKey: refreshKey,
	}
}

func (j *JWTProvider) GenerateAccessToken(user *dto.EnrichedUserResponse) (string, *Claims, error) {
	return GenerateToken(j.accessKey, user, AccessTokenDuration)
}

func (j *JWTProvider) GenerateRefreshToken(user *dto.EnrichedUserResponse) (string, *Claims, error) {
	return GenerateToken(j.refreshKey, user, RefreshTokenDuration)
}

func (j *JWTProvider) VerifyRefreshToken(tokenStr string) (*Claims, error) {
	return Verify(j.refreshKey, tokenStr)
}
