package userHandler

import (
	"msn/internal/modules/user"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/httputils"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	instance *UserHandler
	Once     sync.Once
)

type UserHandler struct {
	userService user.UserService
}

func NewHandler(userService user.UserService) *UserHandler {
	Once.Do(func() {
		instance = &UserHandler{
			userService: userService,
		}
	})
	return instance
}

func (h UserHandler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/users", func(r chi.Router) {
		// Public
		r.Post("/register", h.handleRegister)
	})
}

func (h UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body dto.CreateUser
	err := httputils.ReadRequestBody(w, r, &body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	if body.SubcategoryID != nil && *body.SubcategoryID == "" {
		body.SubcategoryID = nil
	}

	res, err := h.userService.CreateUser(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
