package store

import (
	"motico-api/internal/domain/store/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// Remove
// @Summary      Delete store
// @Description  Delete a store by ID
// @Tags         stores
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Store ID"
// @Success      204          "No Content"
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Store not found"
// @Failure      409          {object}  map[string]interface{}  "Store has associated products"
// @Security     BearerAuth
// @Router       /stores/{id} [delete]
func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.Delete(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrStoreNotFound {
			response.Error(w, http.StatusNotFound, "store not found", nil)
			return
		}
		if err == entities.ErrStoreHasProducts {
			response.Error(w, http.StatusConflict, "store has associated products and cannot be deleted", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete store", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
