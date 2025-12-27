package entities

import "errors"

var (
	ErrStoreNotFound    = errors.New("store not found")
	ErrStoreNameExists  = errors.New("store name already exists for this tenant")
	ErrStoreHasProducts = errors.New("store has associated products and cannot be deleted")
	ErrInvalidStoreName = errors.New("store name is invalid")
)
