package store

import (
	"motico-api/internal/domain/store/entities"
	"motico-api/internal/rest/response"
	restentities "motico-api/internal/rest/store/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// GetByID
// @Summary      Get store by ID
// @Description  Get a specific store by its ID
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Store ID"
// @Success      200         {object}  restentities.StoreResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      404         {object}  map[string]interface{}  "Store not found"
// @Security     BearerAuth
// @Router       /stores/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	store, err := h.service.GetByID(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrStoreNotFound {
			response.Error(w, http.StatusNotFound, "store not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get store", nil)
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
