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

// GetByID
// @Summary      Get transfer by ID
// @Description  Get a specific transfer by its ID
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Transfer ID"
// @Success      200         {object}  restentities.TransferResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      404         {object}  map[string]interface{}  "Transfer not found"
// @Security     BearerAuth
// @Router       /transfers/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.logger.Warn("Invalid tenant ID in get transfer by ID request", logger.String("tenant_id", tenantIDStr))
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid transfer ID in get by ID request", logger.String("transfer_id", idStr), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusBadRequest, "invalid transfer ID", nil)
		return
	}

	h.logger.Info("Getting transfer by ID",
		logger.String("transfer_id", id.String()),
		logger.String("tenant_id", tenantID.String()))

	transfer, err := h.service.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			h.logger.Warn("Transfer not found", logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
			response.Error(w, http.StatusNotFound, "transfer not found", nil)
			return
		}
		h.logger.Error("Failed to get transfer by ID", logger.Error(err), logger.String("transfer_id", id.String()), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusInternalServerError, "failed to get transfer", nil)
		return
	}

	h.logger.Info("Transfer retrieved successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("status", string(transfer.Status)))

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
