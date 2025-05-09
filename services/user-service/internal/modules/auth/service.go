package auth

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/services/user-service/internal/infra/http/middleware"
	"msn/services/user-service/internal/infra/http/token"
	"msn/services/user-service/internal/modules/session"
	"msn/services/user-service/internal/modules/user"
	"net/http"
	"time"
)

const (
	AccessTokenDuration  = time.Minute * 15
	RefreshTokenDuration = time.Hour * 24 * 30
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

	accessToken, _, err := token.Generate(s.AccessKey, user, AccessTokenDuration)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to login")
	}

	refreshToken, refreshTokenClaims, err := token.Generate(s.RefreshKey, user, RefreshTokenDuration)
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

func (s service) Logout(ctx context.Context) error {
	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		slog.Error("context does not contain auth key")
		return fault.NewUnauthorized("access token not provided")
	}

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.User.ID)
	if err != nil {
		return fault.NewBadRequest("failed to retrieve active session")
	}
	if sessRecord == nil {
		return fault.NewNotFound("active session not found")
	}

	sess := session.NewFromModel(*sessRecord)
	sess.Deactivate()

	err = s.sessionRepo.Update(ctx, sess.Model())
	if err != nil {
		return fault.NewBadRequest("failed to deactivate session")
	}

	return nil
}

func (s service) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewTokenResponse, error) {
	claims, err := token.Verify(s.RefreshKey, refreshToken)
	if err != nil {
		return nil, fault.NewUnauthorized("invalid refresh token")
	}

	sessionRecord, err := s.sessionRepo.GetByJTI(ctx, claims.ID)
	if err != nil || sessionRecord == nil || !sessionRecord.Active {
		return nil, fault.NewBadRequest("invalid or inactive session")
	}

	if sessionRecord.ExpiresAt.Before(time.Now()) {
		return nil, fault.NewUnauthorized("session expired")
	}

	user := claims.User

	accessToken, _, err := token.Generate(s.AccessKey, user, AccessTokenDuration)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to generate access token")
	}

	return &dto.RenewTokenResponse{
		AccessToken: accessToken,
	}, nil
}
