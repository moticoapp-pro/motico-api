package transfer

import (
	"motico-api/internal/domain/transfer/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/transfer/entities"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// List
// @Summary      List transfers
// @Description  Get paginated list of transfers for the tenant
// @Tags         transfers
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        page         query     int     false "Page number" default(1)
// @Param        limit        query     int     false "Items per page" default(20)
// @Param        status       query     string  false "Filter by status (pending, completed, cancelled)"
// @Param        store_id     query     string  false "Filter by store ID"
// @Success      200          {object}  restentities.ListTransfersResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid tenant ID"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Security     BearerAuth
// @Router       /transfers [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var status *entities.TransferStatus
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		s := entities.TransferStatus(statusStr)
		if s == entities.TransferStatusPending || s == entities.TransferStatusCompleted || s == entities.TransferStatusCancelled {
			status = &s
		}
	}

	var storeID *uuid.UUID
	if storeIDStr := r.URL.Query().Get("store_id"); storeIDStr != "" {
		id, err := uuid.Parse(storeIDStr)
		if err == nil {
			storeID = &id
		}
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = h.config.Pagination.DefaultLimit
	}
	offset := (page - 1) * limit

	transfers, err := h.service.List(r.Context(), tenantID, status, storeID, limit, offset)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list transfers", nil)
		return
	}

	responses := make([]restentities.TransferResponse, len(transfers))
	for i, t := range transfers {
		responses[i] = restentities.TransferResponse{
			ID:          t.ID,
			TenantID:    t.TenantID,
			ProductID:   t.ProductID,
			FromStoreID: t.FromStoreID,
			ToStoreID:   t.ToStoreID,
			Quantity:    t.Quantity,
			Status:      string(t.Status),
			Notes:       t.Notes,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	total := len(transfers)
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	response.JSON(w, http.StatusOK, restentities.ListTransfersResponse{
		Data: responses,
		Pagination: restentities.PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}
