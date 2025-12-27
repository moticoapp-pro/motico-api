package category

import (
	"encoding/json"
	"motico-api/internal/domain/category"
	"motico-api/internal/domain/category/entities"
	restentities "motico-api/internal/rest/category/entities"
	"motico-api/internal/rest/response"
	"motico-api/pkg/validator"
	"net/http"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var req restentities.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	createReq := category.CreateRequest{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	category, err := h.service.Create(r.Context(), createReq)
	if err != nil {
		if err == entities.ErrCategoryNameExists {
			response.Error(w, http.StatusConflict, "category name already exists", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create category", nil)
		return
	}

	response.JSON(w, http.StatusCreated, restentities.CategoryResponse{
		ID:          category.ID,
		TenantID:    category.TenantID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	})
}
