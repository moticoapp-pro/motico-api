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
		response.Error(w, http.StatusBadRequest, "invalid product ID", nil)
		return
	}

	var req restentities.UpdateProductRequest
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
		StoreID:     &req.StoreID,
		CategoryID:  &req.CategoryID,
		Name:        &req.Name,
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
