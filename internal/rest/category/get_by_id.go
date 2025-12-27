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

// GetByID
// @Summary      Get category by ID
// @Description  Get a specific category by its ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Category ID"
// @Success      200         {object}  restentities.CategoryResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      404         {object}  map[string]interface{}  "Category not found"
// @Security     BearerAuth
// @Router       /categories/{id} [get]
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
