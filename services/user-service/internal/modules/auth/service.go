package auth

import (
	"context"
	"crypto/rsa"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/services/user-service/internal/infra/http/token"
	"msn/services/user-service/internal/modules/session"
	"msn/services/user-service/internal/modules/user"
	"net/http"
	"time"
)

const (
	accessTokenDuration  = time.Minute * 15
	refreshTokenDuration = time.Hour * 24 * 30
)

type ServiceConfig struct {
	AccessKey  *rsa.PrivateKey
	RefreshKey *rsa.PrivateKey

	UserRepo    user.Repository
	SessionRepo session.Repository

	SessionService session.Service
}

type service struct {
	userRepo       user.Repository
	sessionRepo    session.Repository
	sessionService session.Service

	AccessKey  *rsa.PrivateKey
	RefreshKey *rsa.PrivateKey
}

func (s *service) Logout(ctx context.Context) error {
	panic("unimplemented")
}

func NewService(c ServiceConfig) Service {
	return &service{
		userRepo:       c.UserRepo,
		sessionRepo:    c.SessionRepo,
		sessionService: c.SessionService,
		AccessKey:      c.AccessKey,
		RefreshKey:     c.RefreshKey,
	}
}

func (s service) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get user by email")
	}
	if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	if !crypto.PasswordMatches(password, userRecord.Password) {
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	if userRecord.DeletedAt != nil {
		return nil, fault.New(
			"user must be active to login",
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithTag(fault.DISABLED_USER),
		)
	}

	user := dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		CreatedAt: userRecord.CreatedAt,
		UpdatedAt: userRecord.UpdatedAt,
	}

	err = s.sessionRepo.DeactivateAll(ctx, userRecord.ID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	accessToken, _, err := token.Generate(s.AccessKey, user, accessTokenDuration)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to login")
	}

	refreshToken, refreshTokenClaims, err := token.Generate(s.RefreshKey, user, refreshTokenDuration)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to login")
	}

	session, err := s.sessionService.CreateSession(
		ctx,
		dto.CreateSession{
			UserID: userRecord.ID,
			JTI:    refreshTokenClaims.ID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		SessionID:    session.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
