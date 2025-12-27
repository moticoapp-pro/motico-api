package product

import (
	"motico-api/internal/domain/product/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

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
		response.Error(w, http.StatusBadRequest, "invalid product ID", nil)
		return
	}

	err = h.service.Delete(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrProductNotFound {
			response.Error(w, http.StatusNotFound, "product not found", nil)
			return
		}
		if err == entities.ErrProductHasStock {
			response.Error(w, http.StatusConflict, "product has stock and cannot be deleted", nil)
			return
		}
		if err == entities.ErrProductHasTransfers {
			response.Error(w, http.StatusConflict, "product has transfers and cannot be deleted", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete product", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
