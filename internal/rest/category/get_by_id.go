package category

import (
	"motico-api/internal/domain/category/entities"
	restentities "motico-api/internal/rest/category/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid category ID", nil)
		return
	}

	category, err := h.service.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrCategoryNotFound {
			response.Error(w, http.StatusNotFound, "category not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get category", nil)
		return
	}

	response.JSON(w, http.StatusOK, restentities.CategoryResponse{
		ID:          category.ID,
		TenantID:    category.TenantID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	})
}
