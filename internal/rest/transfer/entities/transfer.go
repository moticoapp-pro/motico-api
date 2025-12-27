package entities

import (
	"time"

	"github.com/google/uuid"
)

type CreateTransferRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	FromStoreID uuid.UUID `json:"from_store_id" validate:"required"`
	ToStoreID   uuid.UUID `json:"to_store_id" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required,gt=0"`
	Notes       *string   `json:"notes,omitempty"`
}

type UpdateTransferRequest struct {
	ProductID   uuid.UUID `json:"product_id" validate:"required"`
	FromStoreID uuid.UUID `json:"from_store_id" validate:"required"`
	ToStoreID   uuid.UUID `json:"to_store_id" validate:"required"`
	Quantity    int       `json:"quantity" validate:"required,gt=0"`
	Notes       *string   `json:"notes,omitempty"`
}

type PartialUpdateTransferRequest struct {
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	FromStoreID *uuid.UUID `json:"from_store_id,omitempty"`
	ToStoreID   *uuid.UUID `json:"to_store_id,omitempty"`
	Quantity    *int       `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Notes       *string    `json:"notes,omitempty"`
}

type TransferResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	ProductID   uuid.UUID `json:"product_id"`
	FromStoreID uuid.UUID `json:"from_store_id"`
	ToStoreID   uuid.UUID `json:"to_store_id"`
	Quantity    int       `json:"quantity"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListTransfersResponse struct {
	Data       []TransferResponse `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
