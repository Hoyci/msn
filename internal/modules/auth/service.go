package auth

import (
	"context"
	"msn/internal/infra/http/middlewares"
	"msn/internal/infra/jwt"
	"msn/internal/infra/logging"
	"msn/internal/modules/session"
	"msn/internal/modules/user"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"net/http"
)

type ServiceConfig struct {
	UserRepo       user.UserRepository
	SessionService session.SessionService
	TokenProvider  jwt.JWTProvider
}

type service struct {
	userRepo       user.UserRepository
	sessionService session.SessionService
	tokenProvider  jwt.JWTProvider
}

func NewService(c ServiceConfig) AuthService {
	return &service{
		userRepo:       c.UserRepo,
		sessionService: c.SessionService,
		tokenProvider:  c.TokenProvider,
	}
}

func (s *service) Login(ctx context.Context, email string, password string) (*dto.LoginResponse, error) {
	logger := logging.FromContext(ctx)

	logger.DebugContext(ctx, "login_attempt", "email", email)

	err := ValidateCredentials(email, password)
	if err != nil {
		return nil, fault.New(
			"invalid user data",
			fault.WithTag(fault.BAD_REQUEST),
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithError(err),
		)
	}

	enrichedUser, err := s.userRepo.GetEnrichedByEmail(ctx, email)
	if err != nil {
		logger.ErrorContext(
			ctx, "db_error",
			"operation", "userRepo.GetByEmail",
			"error", err,
		)
		return nil, fault.New(
			"failed to get user",
			fault.WithTag(fault.DB_RESOURCE_NOT_FOUND),
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithError(err),
		)
	}

	if enrichedUser == nil {
		logger.DebugContext(ctx, "user_not_found", "email", email)
		return nil, fault.NewUnauthorized("invalid credentials")
	}

	err = ValidateUser(email, password, enrichedUser.HashedPassword, enrichedUser.DeletedAt)
	if err != nil {
		logger.DebugContext(ctx, "failed_validate_user", "email", email, "error", err)
		return nil, fault.New(
			"failed to validate user",
			fault.WithTag(fault.UNAUTHORIZED),
			fault.WithHTTPCode(http.StatusUnauthorized),
			fault.WithError(err),
		)
	}

	err = s.sessionService.DeactivateAllSessions(ctx, enrichedUser.ID)
	if err != nil {
		return nil, fault.New(
			"failed to deactivate user sessions",
			fault.WithTag(fault.INTERNAL_SERVER_ERROR),
			fault.WithHTTPCode(http.StatusInternalServerError),
			fault.WithError(err),
		)
	}

	accessToken, _, err := s.tokenProvider.GenerateAccessToken(enrichedUser)
	if err != nil {
		logger.ErrorContext(ctx, "access_token_generation_failed", "error", err)
		return nil, fault.NewInternalServerError("failed to login")
	}
	refreshToken, refreshTokenClaims, err := s.tokenProvider.GenerateRefreshToken(enrichedUser)
	if err != nil {
		logger.ErrorContext(ctx, "refresh_token_generation_failed", "error", err)
		return nil, fault.NewInternalServerError("failed to login")
	}

	session, err := s.sessionService.CreateSession(
		ctx, dto.CreateSession{UserID: enrichedUser.ID, JTI: refreshTokenClaims.ID},
	)
	if err != nil {
		logger.ErrorContext(ctx, "session_generation_failed", "error", err)
		return nil, fault.NewInternalServerError("failed to login")
	}

	logger.InfoContext(
		ctx, "login_successful",
		"user_id", enrichedUser.ID,
		"session_id", session.ID,
	)

	return &dto.LoginResponse{
		SessionID:    session.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context) error {
	logger := logging.FromContext(ctx)

	c, ok := ctx.Value(middlewares.AuthKey{}).(*jwt.Claims)
	if !ok {
		logger.ErrorContext(ctx, "missing_auth_context")
		return fault.NewUnauthorized("access token not provided")
	}

	logger.DebugContext(
		ctx, "logout_attempt",
		"user_id", c.User.ID,
		"jti", c.ID,
	)

	activeSession, err := s.sessionService.GetActiveSessionByUserID(ctx, c.User.ID)
	if err != nil {
		logger.ErrorContext(
			ctx, "sessionServiceError",
			"operation", "GetActiveSessionByUserID",
			"user_id", c.User.ID,
			"error", err.Error(),
		)
		return fault.NewBadRequest("failed to retrieve active session")
	}

	if activeSession == nil {
		logger.WarnContext(
			ctx, "no_active_session",
			"user_id", c.User.ID,
		)
		return fault.NewNotFound("active session not found")
	}

	activeSession.Deactivate()

	sess, err := s.sessionService.UpdateSession(ctx, activeSession)
	if err != nil {
		logger.ErrorContext(
			ctx, "sessionServiceError",
			"operation", "UpdateSession",
			"error", err.Error(),
		)
		return fault.NewInternalServerError("failed to login")
	}

	logger.InfoContext(
		ctx, "logout_success",
		"user_id", c.User.ID,
		"session_id", sess.ID,
	)
	return nil
}

// func (s *authService) RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewTokenResponse, error) {
// 	logger := logging.FromContext(ctx)
// 	logger.DebugContext(ctx, "token_renewal_attempt",
//
// 		"token_prefix", "refresh_"+refreshToken[:6])
//
// 	claims, err := token.Verify(s.RefreshKey, refreshToken)
// 	if err != nil {
// 		logger.ErrorContext(ctx, "invalid_refresh_token",
// 			"error", err,
// 		)
// 		return nil, fault.NewUnauthorized("invalid refresh token")
// 	}
//
// 	sessionRecord, err := s.sessionRepo.GetByJTI(ctx, claims.ID)
//
// 	if err != nil || sessionRecord == nil || !sessionRecord.Active {
// 		logger.WarnContext(ctx, "invalid_session",
// 			"jti", claims.ID,
// 			"active", sessionRecord != nil && sessionRecord.Active,
// 			"error", err,
// 		)
// 		return nil, fault.NewBadRequest("invalid or inactive session")
// 	}
//
// 	if sessionRecord.ExpiresAt.Before(time.Now()) {
// 		logger.WarnContext(ctx, "expired_session",
// 			"session_id", sessionRecord.ID,
// 			"expires_at", sessionRecord.ExpiresAt,
// 		)
// 		return nil, fault.NewUnauthorized("session expired")
// 	}
//
// 	accessToken, _, err := token.GenerateToken(s.AccessKey, claims.User, AccessTokenDuration)
// 	if err != nil {
// 		logger.ErrorContext(ctx, "access_token_generation_failed",
// 			"user_id", claims.User.ID,
// 			"error", err,
// 		)
// 		return nil, fault.NewInternalServerError("failed to generate access token")
// 	}
//
// 	logger.InfoContext(ctx, "token_renewed",
//
// 		"user_id", claims.User.ID,
// 		"session_id", sessionRecord.ID,
// 	)
//
// 	return &dto.RenewTokenResponse{
// 		AccessToken: accessToken,
// 	}, nil
// }
