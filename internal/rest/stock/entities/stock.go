package entities

import (
	"time"

	"github.com/google/uuid"
)

type UpdateStockRequest struct {
	Quantity int `json:"quantity" validate:"required,gte=0"`
}

type AdjustStockRequest struct {
	Amount int `json:"amount" validate:"required"`
}

type StockResponse struct {
	ID                uuid.UUID    `json:"id"`
	TenantID          uuid.UUID    `json:"tenant_id"`
	ProductID         uuid.UUID    `json:"product_id"`
	Quantity          int          `json:"quantity"`
	ReservedQuantity  int          `json:"reserved_quantity"`
	AvailableQuantity int          `json:"available_quantity"`
	UpdatedAt         time.Time    `json:"updated_at"`
	Product           *ProductInfo `json:"product,omitempty"`
}

type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	StoreID     uuid.UUID `json:"store_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	SKU         *string   `json:"sku,omitempty"`
	Price       *float64  `json:"price,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
