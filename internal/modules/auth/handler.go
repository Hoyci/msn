package auth

import (
	"crypto/rsa"
	"fmt"
	"msn/internal/infra/http/middleware"
	"msn/internal/infra/logging"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/httputils"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	authHandlerInstance *handler
	Once                sync.Once
)

type handler struct {
	authService AuthService
	accessKey   *rsa.PrivateKey
}

func NewHandler(authService AuthService, accessKey *rsa.PrivateKey) *handler {
	Once.Do(func() {
		authHandlerInstance = &handler{
			authService: authService,
			accessKey:   accessKey,
		}
	})

	return authHandlerInstance
}

func (h handler) RegisterRoutes(r *chi.Mux) {
	m := middleware.NewWithAuth(h.accessKey)
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Private
		r.With(m.WithAuth).Patch("/logout", h.handleLogout)

		// Public
		r.Post("/login", h.handleLogin)
		r.Post("/refresh", h.handleRenewToken)
	})
}

func (h handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)

	var body dto.LoginRequest
	err := httputils.ReadRequestBody(w, r, &body)
	if err != nil {
		logger.ErrorContext(ctx, "parse_body_failed", "error", err)
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.authService.Login(ctx, body.Email, body.Password)
	if err != nil {
		logger.ErrorContext(ctx, "login_failed", "error", err)
		fault.NewHTTPError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   int(RefreshTokenDuration.Seconds()),
	})

	logger.InfoContext(ctx, "login_success")
	httputils.WriteJSON(w, http.StatusOK, res)
}

func (h handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := h.authService.Logout(ctx)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	httputils.WriteSuccess(w, http.StatusOK)
}

func (h handler) handleRenewToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println("r.Cookie", r.Cookies())

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		fault.NewHTTPError(w, fault.NewUnauthorized("refresh token not found"))
		return
	}

	res, err := h.authService.RenewAccessToken(ctx, cookie.Value)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
