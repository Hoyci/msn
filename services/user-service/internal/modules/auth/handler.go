package auth

import (
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/httputils"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	instance *handler
	Once     sync.Once
)

type handler struct {
	authService Service
}

func NewHandler(authService Service) *handler {
	Once.Do(func() {
		instance = &handler{
			authService: authService,
		}
	})

	return instance
}

func (h handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Public
		r.Post("/login", h.handleLogin)
	})
}

func (h handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.LoginRequest
	err := httputils.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.authService.Login(ctx, body.Email, body.Password)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
