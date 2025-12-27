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
