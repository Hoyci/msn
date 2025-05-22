package categoryRepository

import (
	"msn/internal/modules/category"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/httputils"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
)

var (
	categoryHandlerInstance *handler
	Once                    sync.Once
)

type handler struct {
	categoriesService category.Service
}

func NewHandler(categoriesService category.Service) *handler {
	Once.Do(func() {
		categoryHandlerInstance = &handler{
			categoriesService: categoriesService,
		}
	})
	return categoryHandlerInstance
}

func (h handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/categories", func(r chi.Router) {
		// Public
		r.Get("/", h.handleGetCategories)
	})
}

func (h handler) handleGetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	includeSubs := r.URL.Query().Get("include") == "subcategories"

	categories, err := h.categoriesService.GetCategories(ctx, includeSubs)
	if err != nil {
		fault.NewHTTPError(w, err)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, map[string][]*dto.Category{
		"categories": categories,
	})
}
