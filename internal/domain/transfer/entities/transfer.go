package entities

import (
	"time"

	"github.com/google/uuid"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusCancelled TransferStatus = "cancelled"
)

type Transfer struct {
	ID          uuid.UUID      `json:"id"`
	TenantID    uuid.UUID      `json:"tenant_id"`
	ProductID   uuid.UUID      `json:"product_id"`
	FromStoreID uuid.UUID      `json:"from_store_id"`
	ToStoreID   uuid.UUID      `json:"to_store_id"`
	Quantity    int            `json:"quantity"`
	Status      TransferStatus `json:"status"`
	Notes       *string        `json:"notes,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (t *Transfer) IsPending() bool {
	return t.Status == TransferStatusPending
}

func (t *Transfer) IsCompleted() bool {
	return t.Status == TransferStatusCompleted
}

func (t *Transfer) IsCancelled() bool {
	return t.Status == TransferStatusCancelled
}

func (t *Transfer) CanUpdate() bool {
	return t.IsPending()
}

func (t *Transfer) CanDelete() bool {
	return t.IsPending()
}
