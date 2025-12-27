package entities

import "errors"

var (
	ErrStockNotFound         = errors.New("stock not found")
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrInvalidQuantity       = errors.New("quantity must be greater than or equal to zero")
	ErrInvalidReservedAmount = errors.New("reserved quantity cannot exceed total quantity")
)
