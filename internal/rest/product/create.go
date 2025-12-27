package product

import (
	"encoding/json"
	"motico-api/internal/domain/product"
	"motico-api/internal/domain/product/entities"
	restentities "motico-api/internal/rest/product/entities"
	"motico-api/internal/rest/response"
	"net/http"

	"github.com/google/uuid"
	"motico-api/pkg/context"
	"motico-api/pkg/validator"
)

// Create
// @Summary      Create product
// @Description  Create a new product for the tenant
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        X-Tenant-ID  header    string                      true  "Tenant ID"
// @Param        request      body      restentities.CreateProductRequest  true  "Product data"
// @Success      201         {object}  restentities.ProductResponse
// @Failure      400         {object}  map[string]interface{}  "Invalid request"
// @Failure      401         {object}  map[string]interface{}  "Unauthorized"
// @Failure      409         {object}  map[string]interface{}  "Product SKU already exists"
// @Security     BearerAuth
// @Router       /products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantIDStr := context.GetTenantID(r.Context())
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid tenant ID", nil)
		return
	}

	var req restentities.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := validator.ValidateRequest(r, &req); err != nil {
		validator.HandleValidationError(w, err)
		return
	}

	createReq := product.CreateRequest{
		TenantID:    tenantID,
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
	}

	product, err := h.service.Create(r.Context(), createReq)
	if err != nil {
		if err == entities.ErrProductSKUExists {
			response.Error(w, http.StatusConflict, "product SKU already exists for this store", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create product", nil)
		return
	}

	response.JSON(w, http.StatusCreated, restentities.ProductResponse{
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
