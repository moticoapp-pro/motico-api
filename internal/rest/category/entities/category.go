package entities

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type PartialUpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListCategoriesResponse struct {
	Data       []CategoryResponse `json:"data"`
	Pagination PaginationInfo     `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
