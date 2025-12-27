package entities

import "errors"

var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryNameExists  = errors.New("category name already exists for this tenant")
	ErrCategoryHasProducts = errors.New("category has associated products and cannot be deleted")
	ErrInvalidCategoryName = errors.New("category name is invalid")
)
