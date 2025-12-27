package transfer

import (
	"encoding/json"
	"motico-api/internal/domain/transfer"
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"net/http"

	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
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
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var req restentities.CreateTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

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
			response.Error(w, http.StatusBadRequest, "from_store and to_store must be different", nil)
			return
		}
		if err == entities.ErrInvalidQuantity {
			response.Error(w, http.StatusBadRequest, "quantity must be greater than zero", nil)
			return
		}
		if err == entities.ErrInsufficientStock {
			response.Error(w, http.StatusConflict, "insufficient stock available for transfer", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create transfer", nil)
		return
	}

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
