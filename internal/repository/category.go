package repository

import (
	"context"
	"motico-api/internal/domain/category"
	"motico-api/internal/domain/category/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) category.Repository {
	return &categoryRepository{pool: pool}
}

func (r *categoryRepository) Create(ctx context.Context, category *entities.Category) error {
	query := `
		INSERT INTO categories (id, tenant_id, name, description, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, category.TenantID, category.Name, category.Description).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return entities.ErrCategoryNameExists
		}
		return err
	}

	return nil
}

func (r *categoryRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Category, error) {
	query := `
		SELECT id, tenant_id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1 AND tenant_id = $2
	`

	var category entities.Category
	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&category.ID,
		&category.TenantID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Category, error) {
	query := `
		SELECT id, tenant_id, name, description, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(
			&category.ID,
			&category.TenantID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}

	return categories, rows.Err()
}

func (r *categoryRepository) Update(ctx context.Context, category *entities.Category) error {
	query := `
		UPDATE categories
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND tenant_id = $4
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query, category.Name, category.Description, category.ID, category.TenantID).Scan(&category.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrCategoryNotFound
		}
		if isUniqueViolation(err) {
			return entities.ErrCategoryNameExists
		}
		return err
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1 AND tenant_id = $2`

	result, err := r.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entities.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) ExistsByName(ctx context.Context, tenantID uuid.UUID, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE tenant_id = $1 AND name = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) HasProducts(ctx context.Context, tenantID, categoryID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE tenant_id = $1 AND category_id = $2)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, tenantID, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
