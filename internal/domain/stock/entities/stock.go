package entities

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	ProductID        uuid.UUID `json:"product_id"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (s *Stock) AvailableQuantity() int {
	available := s.Quantity - s.ReservedQuantity
	if available < 0 {
		return 0
	}
	return available
}
