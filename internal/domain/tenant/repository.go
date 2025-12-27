package tenant

import (
	"context"
	"motico-api/internal/domain/tenant/entities"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Tenant, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
