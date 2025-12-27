package store

import (
	"encoding/json"
	"motico-api/internal/domain/store"
	"motico-api/internal/domain/store/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/store/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

func (h *Handler) ParcialUpdate(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid store ID", nil)
		return
	}

	var req restentities.PartialUpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	updateReq := store.UpdateRequest{
		ID:       id,
		TenantID: tenantID,
		Name:     req.Name,
		Address:  req.Address,
	}

	store, err := h.service.Update(r.Context(), updateReq)
	if err != nil {
		if err == entities.ErrStoreNotFound {
			response.Error(w, http.StatusNotFound, "store not found", nil)
			return
		}
		if err == entities.ErrStoreNameExists {
			response.Error(w, http.StatusConflict, "store name already exists", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update store", nil)
		return
	}

	response.JSON(w, http.StatusOK, restentities.StoreResponse{
		ID:        store.ID,
		TenantID:  store.TenantID,
		Name:      store.Name,
		Address:   store.Address,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	})
}
