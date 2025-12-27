package transfer

import (
	"encoding/json"
	"motico-api/internal/domain/transfer"
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

// Update
// @Summary      Update transfer
// @Description  Update an existing transfer (only pending transfers can be updated)
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                      true  "Tenant ID"
// @Param        id           path      string                      true  "Transfer ID"
// @Param        request      body      restentities.UpdateTransferRequest  true  "Transfer data"
// @Success      200          {object}  restentities.TransferResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Transfer not found"
// @Security     BearerAuth
// @Router       /transfers/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid transfer ID", nil)
		return
	}

	var req restentities.UpdateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	updateReq := transfer.UpdateRequest{
		ID:          id,
		TenantID:    tenantID,
		ProductID:   &req.ProductID,
		FromStoreID: &req.FromStoreID,
		ToStoreID:   &req.ToStoreID,
		Quantity:    &req.Quantity,
		Notes:       req.Notes,
	}

	transfer, err := h.service.Update(r.Context(), updateReq)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			response.Error(w, http.StatusNotFound, "transfer not found", nil)
			return
		}
		if err == entities.ErrTransferNotPending {
			response.Error(w, http.StatusBadRequest, "transfer is not in pending status", nil)
			return
		}
		if err == entities.ErrInvalidTransferStores {
			response.Error(w, http.StatusBadRequest, "from_store and to_store must be different", nil)
			return
		}
		if err == entities.ErrInvalidQuantity {
			response.Error(w, http.StatusBadRequest, "quantity must be greater than zero", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update transfer", nil)
		return
	}

	response.JSON(w, http.StatusOK, restentities.TransferResponse{
		ID:          transfer.ID,
		TenantID:    transfer.TenantID,
		ProductID:   transfer.ProductID,
		FromStoreID: transfer.FromStoreID,
		ToStoreID:   transfer.ToStoreID,
		Quantity:    transfer.Quantity,
		Status:      string(transfer.Status),
		Notes:       transfer.Notes,
		CreatedAt:   transfer.CreatedAt,
		UpdatedAt:   transfer.UpdatedAt,
	})
}
