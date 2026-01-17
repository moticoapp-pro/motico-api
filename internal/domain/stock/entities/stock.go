package entities

import (
	"time"

	"github.com/google/uuid"
)

type Stock struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	ProductID        uuid.UUID `json:"product_id"`
	StoreID          uuid.UUID `json:"store_id"`
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

type StockTotal struct {
	TenantID          uuid.UUID `json:"tenant_id"`
	ProductID         uuid.UUID `json:"product_id"`
	TotalQuantity     int       `json:"total_quantity"`
	TotalReserved     int       `json:"total_reserved"`
	AvailableQuantity int       `json:"available_quantity"`
	Stores            []Stock   `json:"stores"`
}

func (st *StockTotal) CalculateAvailable() int {
	available := st.TotalQuantity - st.TotalReserved
	if available < 0 {
		return 0
	}
	return available
}
