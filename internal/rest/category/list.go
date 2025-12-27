package category

import (
	"fmt"
	restentities "motico-api/internal/rest/category/entities"
	"motico-api/internal/rest/response"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// List
// @Summary      List categories
// @Description  Get paginated list of categories for the tenant
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        page         query     int     false "Page number" default(1)
// @Param        limit        query     int     false "Items per page" default(20)
// @Success      200          {object}  restentities.ListCategoriesResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid tenant ID"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      500          {object}  map[string]interface{}  "Internal server error"
// @Security     BearerAuth
// @Router       /categories [get]
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
		// Log del error para debugging (temporal)
		// En producción, usar el logger del handler
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to list categories: %v", err), nil)
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
