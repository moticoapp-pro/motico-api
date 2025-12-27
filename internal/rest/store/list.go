package store

import (
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/store/entities"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// List
// @Summary      List stores
// @Description  Get paginated list of stores for the tenant
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        page         query     int     false "Page number" default(1)
// @Param        limit        query     int     false "Items per page" default(20)
// @Success      200          {object}  restentities.ListStoresResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid tenant ID"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Security     BearerAuth
// @Router       /stores [get]
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

	stores, err := h.service.List(r.Context(), tenantID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list stores", nil)
		return
	}

	responses := make([]restentities.StoreResponse, len(stores))
	for i, s := range stores {
		responses[i] = restentities.StoreResponse{
			ID:        s.ID,
			TenantID:  s.TenantID,
			Name:      s.Name,
			Address:   s.Address,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}

	total := len(stores)
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	response.JSON(w, http.StatusOK, restentities.ListStoresResponse{
		Data: responses,
		Pagination: restentities.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
