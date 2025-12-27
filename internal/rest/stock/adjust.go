package stock

import (
	"encoding/json"
	"motico-api/internal/domain/stock"
	"motico-api/internal/domain/stock/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/stock/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

// Adjust
// @Summary      Adjust stock
// @Description  Add or subtract quantity from stock (amount can be positive or negative)
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                    true  "Tenant ID"
// @Param        id           path      string                    true  "Product ID"
// @Param        request      body      restentities.AdjustStockRequest  true  "Adjustment data"
// @Success      200          {object}  restentities.StockResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      409          {object}  map[string]interface{}  "Insufficient stock"
// @Security     BearerAuth
// @Router       /products/{id}/stock [patch]
func (h *Handler) Adjust(w http.ResponseWriter, r *http.Request) {
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

	var req restentities.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	adjustReq := stock.AdjustRequest{
		TenantID:  tenantID,
		ProductID: productID,
		Amount:    req.Amount,
	}

	stock, err := h.service.Adjust(r.Context(), adjustReq)
	if err != nil {
		if err == entities.ErrInvalidQuantity {
			response.Error(w, http.StatusBadRequest, "invalid quantity", nil)
			return
		}
		if err == entities.ErrInsufficientStock {
			response.Error(w, http.StatusConflict, "insufficient stock available", nil)
			return
		}
		if err == entities.ErrInvalidReservedAmount {
			response.Error(w, http.StatusBadRequest, "reserved quantity cannot exceed total quantity", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to adjust stock", nil)
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
