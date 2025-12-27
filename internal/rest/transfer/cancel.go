package transfer

import (
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
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
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid transfer ID", nil)
		return
	}

	transfer, err := h.service.Cancel(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			response.Error(w, http.StatusNotFound, "transfer not found", nil)
			return
		}
		if err == entities.ErrTransferAlreadyCompleted {
			response.Error(w, http.StatusBadRequest, "transfer is already completed", nil)
			return
		}
		if err == entities.ErrTransferAlreadyCancelled {
			response.Error(w, http.StatusBadRequest, "transfer is already cancelled", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to cancel transfer", nil)
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
