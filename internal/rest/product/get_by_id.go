package product

import (
	"motico-api/internal/domain/product/entities"
	restentities "motico-api/internal/rest/product/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// GetByID
// @Summary      Get product by ID
// @Description  Get a specific product by its ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Product ID"
// @Success      200         {object}  restentities.ProductResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      404         {object}  map[string]interface{}  "Product not found"
// @Security     BearerAuth
// @Router       /products/{id} [get]
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
		response.Error(w, http.StatusBadRequest, "invalid product ID", nil)
		return
	}

	product, err := h.service.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrProductNotFound {
			response.Error(w, http.StatusNotFound, "product not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get product", nil)
		return
	}

	stockInfo, _ := h.stockService.GetByProductID(r.Context(), tenantID, id)
	productResponse := restentities.ProductResponse{
		ID:          product.ID,
		TenantID:    product.TenantID,
		StoreID:     product.StoreID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
	if stockInfo != nil {
		productResponse.Stock = &restentities.StockInfo{
			Quantity:          stockInfo.Quantity,
			ReservedQuantity:  stockInfo.ReservedQuantity,
			AvailableQuantity: stockInfo.AvailableQuantity(),
		}
	}

	response.JSON(w, http.StatusOK, productResponse)
}
