package entities

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU         *string   `json:"sku,omitempty" validate:"omitempty,max=100"`
	Price       *float64  `json:"price,omitempty" validate:"omitempty,gte=0"`
}

type UpdateProductRequest struct {
	StoreID     uuid.UUID `json:"store_id" validate:"required"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Name        string    `json:"name" validate:"required,max=255"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU         *string   `json:"sku,omitempty" validate:"omitempty,max=100"`
	Price       *float64  `json:"price,omitempty" validate:"omitempty,gte=0"`
}

type PartialUpdateProductRequest struct {
	StoreID     *uuid.UUID `json:"store_id,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=1000"`
	SKU         *string    `json:"sku,omitempty" validate:"omitempty,max=100"`
	Price       *float64   `json:"price,omitempty" validate:"omitempty,gte=0"`
}

type StockInfo struct {
	Quantity          int `json:"quantity"`
	ReservedQuantity  int `json:"reserved_quantity"`
	AvailableQuantity int `json:"available_quantity"`
}

type ProductResponse struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	StoreID     uuid.UUID  `json:"store_id"`
	CategoryID  uuid.UUID  `json:"category_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	SKU         *string    `json:"sku,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Stock       *StockInfo `json:"stock,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ListProductsResponse struct {
	Data       []ProductResponse `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
