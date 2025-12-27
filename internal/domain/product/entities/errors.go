package entities

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrProductSKUExists    = errors.New("product SKU already exists for this store")
	ErrProductHasStock     = errors.New("product has stock and cannot be deleted")
	ErrProductHasTransfers = errors.New("product has transfers and cannot be deleted")
	ErrInvalidProductName  = errors.New("product name is invalid")
)
