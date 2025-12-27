package entities

import (
	"time"

	"github.com/google/uuid"
)

type CreateStoreRequest struct {
	Name    string  `json:"name" validate:"required,max=255"`
	Address *string `json:"address,omitempty"`
}

type UpdateStoreRequest struct {
	Name    string  `json:"name" validate:"required,max=255"`
	Address *string `json:"address,omitempty"`
}

type PartialUpdateStoreRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Address *string `json:"address,omitempty"`
}

type StoreResponse struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Address   *string   `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListStoresResponse struct {
	Data       []StoreResponse `json:"data"`
	Pagination PaginationInfo  `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
