package stock

import (
	"context"
	"motico-api/internal/domain/stock/entities"

	"github.com/google/uuid"
)

type Repository interface {
	GetByProductID(ctx context.Context, tenantID, productID uuid.UUID) (*entities.Stock, error)
	Create(ctx context.Context, stock *entities.Stock) error
	Update(ctx context.Context, stock *entities.Stock) error
	Reserve(ctx context.Context, tenantID, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, tenantID, productID uuid.UUID, quantity int) error
}
