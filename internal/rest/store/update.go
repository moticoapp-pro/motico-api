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

// Update
// @Summary      Update store
// @Description  Update an existing store (full update)
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                    true  "Tenant ID"
// @Param        id           path      string                    true  "Store ID"
// @Param        request      body      restentities.UpdateStoreRequest  true  "Store data"
// @Success      200          {object}  restentities.StoreResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Store not found"
// @Failure      409          {object}  map[string]interface{}  "Store name already exists"
// @Security     BearerAuth
// @Router       /stores/{id} [put]
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
		response.Error(w, http.StatusBadRequest, "invalid store ID", nil)
		return
	}

	var req restentities.UpdateStoreRequest
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
		Name:     &req.Name,
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
