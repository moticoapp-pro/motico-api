package transfer

import (
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"motico-api/pkg/context"
	"motico-api/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Cancel
// @Summary      Cancel transfer
// @Description  Cancel a transfer and release reserved stock
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Transfer ID"
// @Success      200          {object}  restentities.TransferResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Transfer not found"
// @Security     BearerAuth
// @Router       /transfers/{id}/cancel [patch]
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.logger.Warn("Invalid tenant ID in cancel transfer request", logger.String("tenant_id", tenantIDStr))
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid transfer ID in cancel request", logger.String("transfer_id", idStr), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusBadRequest, "invalid transfer ID", nil)
		return
	}

	h.logger.Info("Cancelling transfer",
		logger.String("transfer_id", id.String()),
		logger.String("tenant_id", tenantID.String()))

	transfer, err := h.service.Cancel(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			h.logger.Warn("Transfer not found for cancellation", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusNotFound, "transfer not found", nil)
			return
		}
		if err == entities.ErrTransferAlreadyCompleted {
			h.logger.Warn("Transfer already completed, cannot cancel", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusBadRequest, "transfer is already completed", nil)
			return
		}
		if err == entities.ErrTransferAlreadyCancelled {
			h.logger.Warn("Transfer already cancelled", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusBadRequest, "transfer is already cancelled", nil)
			return
		}
		h.logger.Error("Failed to cancel transfer", logger.Error(err), logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusInternalServerError, "failed to cancel transfer", nil)
		return
	}

	h.logger.Info("Transfer cancelled successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", transfer.ProductID.String()),
		logger.Int("quantity", transfer.Quantity))

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
