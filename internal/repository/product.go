package repository

import (
	"context"
	"fmt"
	"motico-api/internal/domain/product"
	"motico-api/internal/domain/product/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) product.Repository {
	return &productRepository{pool: pool}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	query := `
		INSERT INTO products (id, tenant_id, store_id, category_id, name, description, sku, price, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		product.TenantID,
		product.StoreID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
	).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return entities.ErrProductSKUExists
		}
		return err
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Product, error) {
	query := `
		SELECT id, tenant_id, store_id, category_id, name, description, sku, price, created_at, updated_at
		FROM products
		WHERE id = $1 AND tenant_id = $2
	`

	var product entities.Product
	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&product.ID,
		&product.TenantID,
		&product.StoreID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) List(ctx context.Context, tenantID uuid.UUID, storeID, categoryID *uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	query := `
		SELECT id, tenant_id, store_id, category_id, name, description, sku, price, created_at, updated_at
		FROM products
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argPos := 2

	if storeID != nil {
		query += ` AND store_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *storeID)
		argPos++
	}

	if categoryID != nil {
		query += ` AND category_id = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *categoryID)
		argPos++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(
			&product.ID,
			&product.TenantID,
			&product.StoreID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.SKU,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	query := `
		UPDATE products
		SET store_id = $1, category_id = $2, name = $3, description = $4, sku = $5, price = $6, updated_at = NOW()
		WHERE id = $7 AND tenant_id = $8
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		product.StoreID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.ID,
		product.TenantID,
	).Scan(&product.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrProductNotFound
		}
		if isUniqueViolation(err) {
			return entities.ErrProductSKUExists
		}
		return err
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1 AND tenant_id = $2`

	result, err := r.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entities.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) ExistsBySKU(ctx context.Context, tenantID, storeID uuid.UUID, sku string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND store_id = $2 AND sku = $3)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, storeID, sku).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *productRepository) HasStock(ctx context.Context, tenantID, productID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM stock WHERE tenant_id = $1 AND product_id = $2 AND quantity > 0)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *productRepository) HasTransfers(ctx context.Context, tenantID, productID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM transfers WHERE tenant_id = $1 AND product_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, productID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
