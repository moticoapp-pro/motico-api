package product

import (
	"context"
	"motico-api/internal/domain/product/entities"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Product, error)
	List(ctx context.Context, tenantID uuid.UUID, storeID, categoryID *uuid.UUID, limit, offset int) ([]*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ExistsBySKU(ctx context.Context, tenantID, storeID uuid.UUID, sku string) (bool, error)
	HasStock(ctx context.Context, tenantID, productID uuid.UUID) (bool, error)
	HasTransfers(ctx context.Context, tenantID, productID uuid.UUID) (bool, error)
}
