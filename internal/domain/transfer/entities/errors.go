package entities

import "errors"

var (
	ErrTransferNotFound         = errors.New("transfer not found")
	ErrTransferNotPending       = errors.New("transfer is not in pending status")
	ErrTransferAlreadyCompleted = errors.New("transfer is already completed")
	ErrTransferAlreadyCancelled = errors.New("transfer is already cancelled")
	ErrInvalidTransferStores    = errors.New("from_store and to_store must be different")
	ErrInvalidQuantity          = errors.New("quantity must be greater than zero")
	ErrInsufficientStock        = errors.New("insufficient stock available for transfer")
)
