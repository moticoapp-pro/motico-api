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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req restentities.UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	updateReq := stock.UpdateRequest{
		TenantID:  tenantID,
		ProductID: productID,
		Quantity:  req.Quantity,
	}

	stock, err := h.service.Update(r.Context(), updateReq)
	if err != nil {
		if err == entities.ErrInvalidQuantity {
			response.Error(w, http.StatusBadRequest, "invalid quantity", nil)
			return
		}
		if err == entities.ErrInvalidReservedAmount {
			response.Error(w, http.StatusBadRequest, "reserved quantity cannot exceed total quantity", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update stock", nil)
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
