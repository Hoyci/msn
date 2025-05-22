package auth

import (
	"crypto/rsa"
	token "msn/internal/infra/http/token"
	"msn/pkg/common/dto"
)

type jwtTokenProvider struct {
	accessKey  *rsa.PrivateKey
	refreshKey *rsa.PrivateKey
}

func NewJWTTokenProvider(accessKey, refreshKey *rsa.PrivateKey) TokenProvider {
	return &jwtTokenProvider{
		accessKey:  accessKey,
		refreshKey: refreshKey,
	}
}

func (j *jwtTokenProvider) GenerateAccessToken(user dto.UserResponse) (string, *token.Claims, error) {
	return token.GenerateToken(j.accessKey, user, AccessTokenDuration)
}

func (j *jwtTokenProvider) GenerateRefreshToken(user dto.UserResponse) (string, *token.Claims, error) {
	return token.GenerateToken(j.refreshKey, user, RefreshTokenDuration)
}

func (j *jwtTokenProvider) VerifyRefreshToken(tokenStr string) (*token.Claims, error) {
	return token.Verify(j.refreshKey, tokenStr)
}
