package categories

import (
	"fmt"
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
	categoriesService Service
}

func NewHandler(categoriesService Service) *handler {
	Once.Do(func() {
		instance = &handler{
			categoriesService: categoriesService,
		}
	})
	return instance
}

func (h handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/categories", func(r chi.Router) {
		// Public
		r.Get("/", h.handleGetCategories)
	})
}

func (h handler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	include := r.URL.Query().Get("include")
	fmt.Println(include)

	switch include {
	case "subcategories":
		res, err := h.categoriesService.GetCategoriesWithSubcategories(ctx)
		if err != nil {
			fault.NewHTTPError(w, err)
			return
		}
		httputils.WriteJSON(w, http.StatusOK, res)
		return

	default:
		res, err := h.categoriesService.GetCategoriesWithUserCount(ctx)
		if err != nil {
			fault.NewHTTPError(w, err)
			return
		}
		httputils.WriteJSON(w, http.StatusOK, res)
		return
	}
}
