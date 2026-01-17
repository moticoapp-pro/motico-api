package transfer

import (
	"encoding/json"
	"motico-api/internal/domain/transfer"
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"motico-api/pkg/context"
	"motico-api/pkg/logger"
	"motico-api/pkg/validator"
	"net/http"

	"github.com/google/uuid"
)

// Create
// @Summary      Create transfer
// @Description  Create a new transfer between stores
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                      true  "Tenant ID"
// @Param        request      body      restentities.CreateTransferRequest  true  "Transfer data"
// @Success      201         {object}  restentities.TransferResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      409         {object}  map[string]interface{}  "Insufficient stock"
// @Security     BearerAuth
// @Router       /transfers [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		h.logger.Warn("Invalid tenant ID in create transfer request", logger.String("tenant_id", tenantIDStr))
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var req restentities.CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body in create transfer", logger.Error(err), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		h.logger.Warn("Validation error in create transfer request", logger.Error(err), logger.String("tenant_id", tenantID.String()))
		validator.HandleValidationError(w, err)
		return
	}

	h.logger.Info("Creating transfer request",
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", req.ProductID.String()),
		logger.String("from_store_id", req.FromStoreID.String()),
		logger.String("to_store_id", req.ToStoreID.String()),
		logger.Int("quantity", req.Quantity))

	createReq := transfer.CreateRequest{
		TenantID:    tenantID,
		ProductID:   req.ProductID,
		FromStoreID: req.FromStoreID,
		ToStoreID:   req.ToStoreID,
		Quantity:    req.Quantity,
		Notes:       req.Notes,
	}

	transfer, err := h.service.Create(r.Context(), createReq)
	if err != nil {
		if err == entities.ErrInvalidTransferStores {
			h.logger.Warn("Invalid transfer stores", logger.String("tenant_id", tenantID.String()), logger.String("from_store_id", req.FromStoreID.String()), logger.String("to_store_id", req.ToStoreID.String()))
			response.Error(w, http.StatusBadRequest, "from_store and to_store must be different", nil)
			return
		}
		if err == entities.ErrInvalidQuantity {
			h.logger.Warn("Invalid quantity in transfer request", logger.String("tenant_id", tenantID.String()), logger.Int("quantity", req.Quantity))
			response.Error(w, http.StatusBadRequest, "quantity must be greater than zero", nil)
			return
		}
		if err == entities.ErrInsufficientStock {
			h.logger.Warn("Insufficient stock for transfer", logger.String("tenant_id", tenantID.String()), logger.String("product_id", req.ProductID.String()), logger.Int("quantity", req.Quantity))
			response.Error(w, http.StatusConflict, "insufficient stock available for transfer", nil)
			return
		}
		h.logger.Error("Failed to create transfer", logger.Error(err), logger.String("tenant_id", tenantID.String()))
		response.Error(w, http.StatusInternalServerError, "failed to create transfer", nil)
		return
	}

	h.logger.Info("Transfer created successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", transfer.ProductID.String()),
		logger.String("status", string(transfer.Status)))

	response.JSON(w, http.StatusCreated, restentities.TransferResponse{
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
