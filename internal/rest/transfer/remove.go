package transfer

import (
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	"motico-api/pkg/context"
	"motico-api/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Remove
// @Summary      Delete transfer
// @Description  Delete a transfer by ID (only pending transfers can be deleted)
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Transfer ID"
// @Success      204          "No Content"
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Transfer not found"
// @Security     BearerAuth
// @Router       /transfers/{id} [delete]
func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.logger.Warn("Invalid tenant ID in remove transfer request", logger.String("tenant_id", tenantIDStr))
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid transfer ID in remove request", logger.String("transfer_id", idStr), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusBadRequest, "invalid transfer ID", nil)
		return
	}

	h.logger.Info("Deleting transfer",
		logger.String("transfer_id", id.String()),
		logger.String("tenant_id", tenantID.String()))

	err = h.service.Delete(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			h.logger.Warn("Transfer not found for deletion", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusNotFound, "transfer not found", nil)
			return
		}
		if err == entities.ErrTransferNotPending {
			h.logger.Warn("Transfer is not pending, cannot delete", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusBadRequest, "transfer is not in pending status", nil)
			return
		}
		h.logger.Error("Failed to delete transfer", logger.Error(err), logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusInternalServerError, "failed to delete transfer", nil)
		return
	}

	h.logger.Info("Transfer deleted successfully",
		logger.String("transfer_id", id.String()),
		logger.String("tenant_id", tenantID.String()))

	w.WriteHeader(http.StatusNoContent)
}
