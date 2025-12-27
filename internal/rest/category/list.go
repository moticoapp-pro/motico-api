package category

import (
	restentities "motico-api/internal/rest/category/entities"
	"motico-api/internal/rest/response"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = h.config.Pagination.DefaultLimit
	}
	offset := (page - 1) * limit

	categories, err := h.service.List(r.Context(), tenantID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list categories", nil)
		return
	}

	responses := make([]restentities.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = restentities.CategoryResponse{
			ID:          cat.ID,
			TenantID:    cat.TenantID,
			Name:        cat.Name,
			Description: cat.Description,
			CreatedAt:   cat.CreatedAt,
			UpdatedAt:   cat.UpdatedAt,
		}
	}

	total := len(categories)
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	response.JSON(w, http.StatusOK, restentities.ListCategoriesResponse{
		Data: responses,
		Pagination: restentities.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
