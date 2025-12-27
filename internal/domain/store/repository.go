package store

import (
	"context"
	"motico-api/internal/domain/store/entities"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, store *entities.Store) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Store, error)
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Store, error)
	Update(ctx context.Context, store *entities.Store) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ExistsByName(ctx context.Context, tenantID uuid.UUID, name string) (bool, error)
	HasProducts(ctx context.Context, tenantID, storeID uuid.UUID) (bool, error)
}
