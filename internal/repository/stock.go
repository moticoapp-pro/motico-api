package repository

import (
	"context"
	"motico-api/internal/domain/stock"
	"motico-api/internal/domain/stock/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stockRepository struct {
	pool *pgxpool.Pool
}

func NewStockRepository(pool *pgxpool.Pool) stock.Repository {
	return &stockRepository{pool: pool}
}

func (r *stockRepository) GetByProductIDAndStoreID(ctx context.Context, tenantID, productID, storeID uuid.UUID) (*entities.Stock, error) {
	query := `
		SELECT id, tenant_id, product_id, store_id, quantity, reserved_quantity, created_at, updated_at
		FROM stock
		WHERE tenant_id = $1 AND product_id = $2 AND store_id = $3
	`

	var stock entities.Stock
	err := r.pool.QueryRow(ctx, query, tenantID, productID, storeID).Scan(
		&stock.ID,
		&stock.TenantID,
		&stock.ProductID,
		&stock.StoreID,
		&stock.Quantity,
		&stock.ReservedQuantity,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrStockNotFound
		}
		return nil, err
	}

	return &stock, nil
}

func (r *stockRepository) ListByProductID(ctx context.Context, tenantID, productID uuid.UUID) ([]*entities.Stock, error) {
	query := `
		SELECT id, tenant_id, product_id, store_id, quantity, reserved_quantity, created_at, updated_at
		FROM stock
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY store_id
	`

	rows, err := r.pool.Query(ctx, query, tenantID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []*entities.Stock
	for rows.Next() {
		var stock entities.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.TenantID,
			&stock.ProductID,
			&stock.StoreID,
			&stock.Quantity,
			&stock.ReservedQuantity,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &stock)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *stockRepository) Create(ctx context.Context, stock *entities.Stock) error {
	query := `
		INSERT INTO stock (id, tenant_id, product_id, store_id, quantity, reserved_quantity, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		stock.TenantID,
		stock.ProductID,
		stock.StoreID,
		stock.Quantity,
		stock.ReservedQuantity,
	).Scan(
		&stock.ID,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *stockRepository) Update(ctx context.Context, stock *entities.Stock) error {
	query := `
		UPDATE stock
		SET quantity = $1, reserved_quantity = $2, updated_at = NOW()
		WHERE tenant_id = $3 AND product_id = $4 AND store_id = $5
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		stock.Quantity,
		stock.ReservedQuantity,
		stock.TenantID,
		stock.ProductID,
		stock.StoreID,
	).Scan(&stock.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrStockNotFound
		}
		return err
	}

	return nil
}

func (r *stockRepository) Reserve(ctx context.Context, tenantID, productID, storeID uuid.UUID, quantity int) error {
	query := `
		UPDATE stock
		SET reserved_quantity = reserved_quantity + $1, updated_at = NOW()
		WHERE tenant_id = $2 AND product_id = $3 AND store_id = $4
			AND (quantity - reserved_quantity) >= $1
		RETURNING id
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, quantity, tenantID, productID, storeID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrInsufficientStock
		}
		return err
	}

	return nil
}

func (r *stockRepository) Release(ctx context.Context, tenantID, productID, storeID uuid.UUID, quantity int) error {
	query := `
		UPDATE stock
		SET reserved_quantity = GREATEST(0, reserved_quantity - $1), updated_at = NOW()
		WHERE tenant_id = $2 AND product_id = $3 AND store_id = $4
		RETURNING id
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, quantity, tenantID, productID, storeID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrStockNotFound
		}
		return err
	}

	return nil
}
