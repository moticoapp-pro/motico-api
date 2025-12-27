package category

import (
	"motico-api/internal/domain/category/entities"
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
		response.Error(w, http.StatusBadRequest, "invalid category ID", nil)
		return
	}

	err = h.service.Delete(r.Context(), tenantID, id)
	if err != nil {
		if err == entities.ErrCategoryNotFound {
			response.Error(w, http.StatusNotFound, "category not found", nil)
			return
		}
		if err == entities.ErrCategoryHasProducts {
			response.Error(w, http.StatusConflict, "category has associated products and cannot be deleted", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete category", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
