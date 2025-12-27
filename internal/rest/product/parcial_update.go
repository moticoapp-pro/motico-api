package product

import (
	"encoding/json"
	"motico-api/internal/domain/product"
	"motico-api/internal/domain/product/entities"
	restentities "motico-api/internal/rest/product/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

// ParcialUpdate
// @Summary      Partially update product
// @Description  Update specific fields of a product (partial update)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                            true  "Tenant ID"
// @Param        id           path      string                            true  "Product ID"
// @Param        request      body      restentities.PartialUpdateProductRequest  true  "Product data"
// @Success      200          {object}  restentities.ProductResponse
// @Failure      400          {object}  map[string]interface{}  "Invalid request"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Failure      404          {object}  map[string]interface{}  "Product not found"
// @Failure      409          {object}  map[string]interface{}  "Product SKU already exists"
// @Security     BearerAuth
// @Router       /products/{id} [patch]
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
		response.Error(w, http.StatusBadRequest, "invalid product ID", nil)
		return
	}

	var req restentities.PartialUpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	updateReq := product.UpdateRequest{
		ID:          id,
		TenantID:    tenantID,
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
	}

	product, err := h.service.Update(r.Context(), updateReq)
	if err != nil {
		if err == entities.ErrProductNotFound {
			response.Error(w, http.StatusNotFound, "product not found", nil)
			return
		}
		if err == entities.ErrProductSKUExists {
			response.Error(w, http.StatusConflict, "product SKU already exists for this store", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update product", nil)
		return
	}

	response.JSON(w, http.StatusOK, restentities.ProductResponse{
		ID:          product.ID,
		TenantID:    product.TenantID,
		StoreID:     product.StoreID,
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	})
}
