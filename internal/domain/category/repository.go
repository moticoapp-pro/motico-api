package category

import (
	"context"
	"motico-api/internal/domain/category/entities"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, category *entities.Category) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Category, error)
	List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	ExistsByName(ctx context.Context, tenantID uuid.UUID, name string) (bool, error)
	HasProducts(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error)
}
