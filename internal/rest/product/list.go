package product

import (
	restentities "motico-api/internal/rest/product/entities"
	"motico-api/internal/rest/response"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// List
// @Summary      List products
// @Description  Get paginated list of products for the tenant
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        page         query     int     false "Page number" default(1)
// @Param        limit        query     int     false "Items per page" default(20)
// @Param        store_id     query     string  false "Filter by store ID"
// @Param        category_id  query     string  false "Filter by category ID"
// @Success      200          {object}  restentities.ListProductsResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid tenant ID"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Security     BearerAuth
// @Router       /products [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var storeID, categoryID *uuid.UUID
	if storeIDStr := r.URL.Query().Get("store_id"); storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err == nil {
			storeID = &id
		}
	}
	if categoryIDStr := r.URL.Query().Get("category_id"); categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &id
		}
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

	products, err := h.service.List(r.Context(), tenantID, storeID, categoryID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list products", nil)
		return
	}

	responses := make([]restentities.ProductResponse, len(products))
	for i, p := range products {
		stockInfo, _ := h.stockService.GetByProductID(r.Context(), tenantID, p.ID)
		response := restentities.ProductResponse{
			ID:          p.ID,
			TenantID:    p.TenantID,
			StoreID:     p.StoreID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			Price:       p.Price,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
		if stockInfo != nil {
			response.Stock = &restentities.StockInfo{
				Quantity:          stockInfo.Quantity,
				ReservedQuantity:  stockInfo.ReservedQuantity,
				AvailableQuantity: stockInfo.AvailableQuantity(),
			}
		}
		responses[i] = response
	}

	total := len(products)
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	response.JSON(w, http.StatusOK, restentities.ListProductsResponse{
		Data: responses,
		Pagination: restentities.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
