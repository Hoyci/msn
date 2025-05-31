package userHandler

import (
	"fmt"
	"msn/internal/infra/storage"
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
	userService   user.UserService
	storageClient *storage.StorageClient
}

func NewHandler(userService user.UserService, storageClient *storage.StorageClient) *UserHandler {
	Once.Do(
		func() {
			instance = &UserHandler{
				userService:   userService,
				storageClient: storageClient,
			}
		},
	)
	return instance
}

func (h UserHandler) RegisterRoutes(r *chi.Mux) {
	r.Route(
		"/api/v1/users", func(r chi.Router) {
			// Public
			r.Post("/register", h.handleRegister)
			r.Get("/professionals", h.handleGetProfessionals)
		},
	)
}

func (h UserHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 6<<20)

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		fault.NewHTTPError(w, fmt.Errorf("invalid form data: %w", err))
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")
	userRole := r.FormValue("role")
	subcat := r.FormValue("subcategory_id")
	var subcategoryID *string
	if subcat != "" {
		subcategoryID = &subcat
	}

	file, header, err := r.FormFile("picture")
	if err != nil {
		fault.New(
			"picture is required",
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithTag(fault.BAD_REQUEST),
		)
		return
	}
	defer file.Close()

	if header.Size > 5<<20 {
		fault.New(
			"picture must be at most 5mb",
			fault.WithHTTPCode(http.StatusBadRequest),
			fault.WithTag(fault.BAD_REQUEST),
		)
		return
	}

	body := dto.CreateUser{
		Name:            name,
		Email:           email,
		Password:        password,
		ConfirmPassword: confirm,
		FileHeader:      header,
		Role:            userRole,
		SubcategoryID:   subcategoryID,
	}

	res, err := h.userService.CreateUser(ctx, body)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (h UserHandler) handleGetProfessionals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res, err := h.userService.GetProfessionalUsers(ctx)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, map[string][]*dto.ProfessionalUserResponse{"professionals": res})
}
