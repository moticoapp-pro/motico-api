package repository

import (
	"context"
	"motico-api/internal/domain/store"
	"motico-api/internal/domain/store/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type storeRepository struct {
	pool *pgxpool.Pool
}

func NewStoreRepository(pool *pgxpool.Pool) store.Repository {
	return &storeRepository{pool: pool}
}

func (r *storeRepository) Create(ctx context.Context, store *entities.Store) error {
	query := `
		INSERT INTO stores (id, tenant_id, name, address, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, store.TenantID, store.Name, store.Address).Scan(
		&store.ID,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return entities.ErrStoreNameExists
		}
		return err
	}

	return nil
}

func (r *storeRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Store, error) {
	query := `
		SELECT id, tenant_id, name, address, created_at, updated_at
		FROM stores
		WHERE id = $1 AND tenant_id = $2
	`

	var store entities.Store
	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&store.ID,
		&store.TenantID,
		&store.Name,
		&store.Address,
		&store.CreatedAt,
		&store.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrStoreNotFound
		}
		return nil, err
	}

	return &store, nil
}

func (r *storeRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Store, error) {
	query := `
		SELECT id, tenant_id, name, address, created_at, updated_at
		FROM stores
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []*entities.Store
	for rows.Next() {
		var store entities.Store
		if err := rows.Scan(
			&store.ID,
			&store.TenantID,
			&store.Name,
			&store.Address,
			&store.CreatedAt,
			&store.UpdatedAt,
		); err != nil {
			return nil, err
		}
		stores = append(stores, &store)
	}

	return stores, rows.Err()
}

func (r *storeRepository) Update(ctx context.Context, store *entities.Store) error {
	query := `
		UPDATE stores
		SET name = $1, address = $2, updated_at = NOW()
		WHERE id = $3 AND tenant_id = $4
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query, store.Name, store.Address, store.ID, store.TenantID).Scan(&store.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrStoreNotFound
		}
		if isUniqueViolation(err) {
			return entities.ErrStoreNameExists
		}
		return err
	}

	return nil
}

func (r *storeRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	query := `DELETE FROM stores WHERE id = $1 AND tenant_id = $2`

	result, err := r.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entities.ErrStoreNotFound
	}

	return nil
}

func (r *storeRepository) ExistsByName(ctx context.Context, tenantID uuid.UUID, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stores WHERE tenant_id = $1 AND name = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *storeRepository) HasProducts(ctx context.Context, tenantID, storeID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND store_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, storeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
