package product

import (
	"motico-api/internal/domain/product/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
)

// Remove
// @Summary      Delete product
// @Description  Delete a product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string  true  "Tenant ID"
// @Param        id           path      string  true  "Product ID"
// @Success      204          "No Content"
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Product not found"
// @Failure      409          {object}  map[string]interface{}  "Product has stock or transfers"
// @Security     BearerAuth
// @Router       /products/{id} [delete]
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
