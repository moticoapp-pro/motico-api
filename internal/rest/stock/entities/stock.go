package entities

import (
	"time"

	"github.com/google/uuid"
)

type UpdateStockRequest struct {
	StoreID  uuid.UUID `json:"store_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,gte=0"`
}

type AdjustStockRequest struct {
	StoreID uuid.UUID `json:"store_id" validate:"required"`
	Amount  int       `json:"amount" validate:"required"`
}

type StockResponse struct {
	ID                uuid.UUID    `json:"id"`
	TenantID          uuid.UUID    `json:"tenant_id"`
	ProductID         uuid.UUID    `json:"product_id"`
	StoreID           uuid.UUID    `json:"store_id"`
	Quantity          int          `json:"quantity"`
	ReservedQuantity  int          `json:"reserved_quantity"`
	AvailableQuantity int          `json:"available_quantity"`
	UpdatedAt         time.Time    `json:"updated_at"`
	Product           *ProductInfo `json:"product,omitempty"`
}

type StockTotalResponse struct {
	TenantID          uuid.UUID      `json:"tenant_id"`
	ProductID         uuid.UUID      `json:"product_id"`
	TotalQuantity     int            `json:"total_quantity"`
	TotalReserved     int            `json:"total_reserved"`
	AvailableQuantity int            `json:"available_quantity"`
	Stores            []StockByStore `json:"stores"`
	Product           *ProductInfo   `json:"product,omitempty"`
}

type StockByStore struct {
	StoreID           uuid.UUID `json:"store_id"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
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
