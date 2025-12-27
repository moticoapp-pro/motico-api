package transfer

import (
	"context"
	"motico-api/internal/domain/transfer/entities"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, transfer *entities.Transfer) error
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error)
	List(ctx context.Context, tenantID uuid.UUID, status *entities.TransferStatus, storeID *uuid.UUID, limit, offset int) ([]*entities.Transfer, error)
	Update(ctx context.Context, transfer *entities.Transfer) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
}
