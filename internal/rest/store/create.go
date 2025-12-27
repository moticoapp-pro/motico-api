package store

import (
	"encoding/json"
	"motico-api/internal/domain/store"
	"motico-api/internal/domain/store/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/store/entities"
	"net/http"

	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var req restentities.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	createReq := store.CreateRequest{
		TenantID: tenantID,
		Name:     req.Name,
		Address:  req.Address,
	}

	store, err := h.service.Create(r.Context(), createReq)
	if err != nil {
		if err == entities.ErrStoreNameExists {
			response.Error(w, http.StatusConflict, "store name already exists", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create store", nil)
		return
	}

	response.JSON(w, http.StatusCreated, restentities.StoreResponse{
		ID:        store.ID,
		TenantID:  store.TenantID,
		Name:      store.Name,
		Address:   store.Address,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	})
}
