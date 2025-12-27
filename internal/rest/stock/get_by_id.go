package stock

import (
	"motico-api/internal/domain/stock/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/stock/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// GetByID
// @Summary      Get stock by product ID
// @Description  Get stock information for a specific product
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Product ID"
// @Success      200         {object}  restentities.StockResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      404         {object}  map[string]interface{}  "Stock not found"
// @Security     BearerAuth
// @Router       /products/{id}/stock [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	productIDStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid product ID", nil)
		return
	}

	stock, err := h.service.GetByProductID(r.Context(), tenantID, productID)
	if err != nil {
		if err == entities.ErrStockNotFound {
			response.Error(w, http.StatusNotFound, "stock not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get stock", nil)
		return
	}

	response.JSON(w, http.StatusOK, restentities.StockResponse{
		ID:                stock.ID,
		TenantID:          stock.TenantID,
		ProductID:         stock.ProductID,
		Quantity:          stock.Quantity,
		ReservedQuantity:  stock.ReservedQuantity,
		AvailableQuantity: stock.AvailableQuantity(),
		UpdatedAt:         stock.UpdatedAt,
	})
}
