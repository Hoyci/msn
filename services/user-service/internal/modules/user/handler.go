package user

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
	userService Service
}

func NewHandler(userService Service) *handler {
	Once.Do(func() {
		instance = &handler{
			userService: userService,
		}
	})
	return instance
}

func (h handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/users", func(r chi.Router) {
		// Public
		r.Post("/register", h.handleRegister)
	})
}

func (h handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateUser
	err := httputils.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	res, err := h.userService.CreateUser(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
