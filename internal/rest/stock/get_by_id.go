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
