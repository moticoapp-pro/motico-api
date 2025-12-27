package repository

import (
	"context"
	"motico-api/internal/domain/tenant"
	"motico-api/internal/domain/tenant/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tenantRepository struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) tenant.Repository {
	return &tenantRepository{pool: pool}
}

func (r *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Tenant, error) {
	query := `SELECT id, name, created_at, updated_at FROM tenants WHERE id = $1`

	var tenant entities.Tenant
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrTenantNotFound
		}
		return nil, err
	}

	return &tenant, nil
}

func (r *tenantRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM tenants WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
