package auth

import (
	"context"
	"crypto/rsa"
	"msn/internal/infra/http/middleware"
	"msn/internal/infra/http/token"
	"msn/internal/infra/logging"
	"msn/internal/modules/session"
	"msn/internal/modules/user"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"net/http"
	"time"
)

const (
	AccessTokenDuration  = time.Minute * 2
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
	logger := logging.FromContext(ctx)

	logger.DebugContext(ctx, "login_attempt", "email", email)

	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "userRepo.GetByEmail",
			"error", err,
		)
		return nil, fault.NewBadRequest("failed to get user by email")
	}

	if userRecord == nil {
		logger.DebugContext(ctx, "user_not_found", "email", email)
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	if !crypto.PasswordMatches(password, userRecord.Password) {
		logger.DebugContext(ctx, "invalid_password",
			"user_id", userRecord.ID,
			"email", email,
		)
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	if userRecord.DeletedAt != nil {
		logger.WarnContext(ctx, "disabled_user_login_attempt",
			"user_id", userRecord.ID,
			"email", email,
		)

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
		logger.ErrorContext(ctx, "db_error",
			"operation", "sessionRepo.DeactivateAll",
			"error", err,
		)
		return nil, fault.NewBadRequest("failed to deactivate user sessions")
	}

	accessToken, _, err := token.Generate(s.AccessKey, user, AccessTokenDuration)
	if err != nil {
		logger.ErrorContext(ctx, "access_token_generation_failed", "error", err)
		return nil, fault.NewInternalServerError("failed to login")
	}

	refreshToken, refreshTokenClaims, err := token.Generate(s.RefreshKey, user, RefreshTokenDuration)
	if err != nil {
		logger.ErrorContext(ctx, "refresh_token_generation_failed", "error", err)
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
		logger.ErrorContext(ctx, "session_creation_failed", "error", err)
		return nil, err
	}

	logger.InfoContext(ctx, "login_successful",
		"user_id", userRecord.ID,
		"session_id", session.ID,
	)

	return &dto.LoginResponse{
		SessionID:    session.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s service) Logout(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	c, ok := ctx.Value(middleware.AuthKey{}).(*token.Claims)
	if !ok {
		logger.ErrorContext(ctx, "missing_auth_context")
		return fault.NewUnauthorized("access token not provided")
	}

	logger.DebugContext(ctx, "logout_attempt",
		"user_id", c.User.ID,
		"jti", c.ID,
	)

	sessRecord, err := s.sessionRepo.GetActiveByUserID(ctx, c.User.ID)
	if err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "sessionRepo.GetActiveByUserID",
			"user_id", c.User.ID,
			"error", err,
		)
		return fault.NewBadRequest("failed to retrieve active session")
	}

	if sessRecord == nil {
		logger.WarnContext(ctx, "no_active_session",
			"user_id", c.User.ID,
		)
		return fault.NewNotFound("active session not found")
	}

	sess := session.NewFromModel(*sessRecord)
	sess.Deactivate()

	err = s.sessionRepo.Update(ctx, sess.Model())
	if err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "sessionRepo.Update",
			"session_id", sess.ID,
			"error", err,
		)
		return fault.NewBadRequest("failed to deactivate session")
	}

	logger.InfoContext(ctx, "logout_success",
		"user_id", c.User.ID,
		"session_id", sess.ID,
	)

	return nil
}

func (s service) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewTokenResponse, error) {
	logger := logging.FromContext(ctx)

	logger.DebugContext(ctx, "token_renewal_attempt",
		"token_prefix", "refresh_"+refreshToken[:6])

	claims, err := token.Verify(s.RefreshKey, refreshToken)
	if err != nil {
		logger.ErrorContext(ctx, "invalid_refresh_token",
			"error", err,
		)
		return nil, fault.NewUnauthorized("invalid refresh token")
	}

	sessionRecord, err := s.sessionRepo.GetByJTI(ctx, claims.ID)
	if err != nil || sessionRecord == nil || !sessionRecord.Active {
		logger.WarnContext(ctx, "invalid_session",
			"jti", claims.ID,
			"active", sessionRecord != nil && sessionRecord.Active,
			"error", err,
		)
		return nil, fault.NewBadRequest("invalid or inactive session")
	}

	if sessionRecord.ExpiresAt.Before(time.Now()) {
		logger.WarnContext(ctx, "expired_session",
			"session_id", sessionRecord.ID,
			"expires_at", sessionRecord.ExpiresAt,
		)
		return nil, fault.NewUnauthorized("session expired")
	}

	accessToken, _, err := token.Generate(s.AccessKey, claims.User, AccessTokenDuration)
	if err != nil {
		logger.ErrorContext(ctx, "access_token_generation_failed",
			"user_id", claims.User.ID,
			"error", err,
		)
		return nil, fault.NewInternalServerError("failed to generate access token")
	}

	logger.InfoContext(ctx, "token_renewed",
		"user_id", claims.User.ID,
		"session_id", sessionRecord.ID,
	)

	return &dto.RenewTokenResponse{
		AccessToken: accessToken,
	}, nil
}
