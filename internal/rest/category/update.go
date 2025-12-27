package category

import (
	"encoding/json"
	"motico-api/internal/domain/category"
	"motico-api/internal/domain/category/entities"
	restentities "motico-api/internal/rest/category/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

// Update
// @Summary      Update category
// @Description  Update an existing category (full update)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                        true  "Tenant ID"
// @Param        id           path      string                        true  "Category ID"
// @Param        request      body      restentities.UpdateCategoryRequest  true  "Category data"
// @Success      200          {object}  restentities.CategoryResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Category not found"
// @Failure      409          {object}  map[string]interface{}  "Category name already exists"
// @Security     BearerAuth
// @Router       /categories/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req restentities.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	updateReq := category.UpdateRequest{
		ID:          id,
		TenantID:    tenantID,
		Name:        &req.Name,
		Description: req.Description,
	}

	category, err := h.service.Update(r.Context(), updateReq)
	if err != nil {
		if err == entities.ErrCategoryNotFound {
			response.Error(w, http.StatusNotFound, "category not found", nil)
			return
		}
		if err == entities.ErrCategoryNameExists {
			response.Error(w, http.StatusConflict, "category name already exists", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update category", nil)
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
