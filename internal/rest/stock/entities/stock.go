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
	ID                uuid.UUID `json:"id"`
	TenantID          uuid.UUID `json:"tenant_id"`
	ProductID         uuid.UUID `json:"product_id"`
	Quantity          int       `json:"quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	UpdatedAt         time.Time `json:"updated_at"`
}
