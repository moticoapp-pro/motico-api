package repository

import (
	"context"
	"fmt"
	"motico-api/internal/domain/transfer"
	"motico-api/internal/domain/transfer/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transferRepository struct {
	pool *pgxpool.Pool
}

func NewTransferRepository(pool *pgxpool.Pool) transfer.Repository {
	return &transferRepository{pool: pool}
}

func (r *transferRepository) Create(ctx context.Context, transfer *entities.Transfer) error {
	query := `
		INSERT INTO transfers (id, tenant_id, product_id, from_store_id, to_store_id, quantity, status, notes, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		transfer.TenantID,
		transfer.ProductID,
		transfer.FromStoreID,
		transfer.ToStoreID,
		transfer.Quantity,
		transfer.Status,
		transfer.Notes,
	).Scan(
		&transfer.ID,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *transferRepository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	query := `
		SELECT id, tenant_id, product_id, from_store_id, to_store_id, quantity, status, notes, created_at, updated_at
		FROM transfers
		WHERE id = $1 AND tenant_id = $2
	`

	var transfer entities.Transfer
	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&transfer.ID,
		&transfer.TenantID,
		&transfer.ProductID,
		&transfer.FromStoreID,
		&transfer.ToStoreID,
		&transfer.Quantity,
		&transfer.Status,
		&transfer.Notes,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrTransferNotFound
		}
		return nil, err
	}

	return &transfer, nil
}

func (r *transferRepository) List(ctx context.Context, tenantID uuid.UUID, status *entities.TransferStatus, storeID *uuid.UUID, limit, offset int) ([]*entities.Transfer, error) {
	query := `
		SELECT id, tenant_id, product_id, from_store_id, to_store_id, quantity, status, notes, created_at, updated_at
		FROM transfers
		WHERE tenant_id = $1
	`
	args := []interface{}{tenantID}
	argPos := 2

	if status != nil {
		query += ` AND status = $` + fmt.Sprintf("%d", argPos)
		args = append(args, *status)
		argPos++
	}

	if storeID != nil {
		query += ` AND (from_store_id = $` + fmt.Sprintf("%d", argPos) + ` OR to_store_id = $` + fmt.Sprintf("%d", argPos) + `)`
		args = append(args, *storeID)
		argPos++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argPos) + ` OFFSET $` + fmt.Sprintf("%d", argPos+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*entities.Transfer
	for rows.Next() {
		var transfer entities.Transfer
		if err := rows.Scan(
			&transfer.ID,
			&transfer.TenantID,
			&transfer.ProductID,
			&transfer.FromStoreID,
			&transfer.ToStoreID,
			&transfer.Quantity,
			&transfer.Status,
			&transfer.Notes,
			&transfer.CreatedAt,
			&transfer.UpdatedAt,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, &transfer)
	}

	return transfers, rows.Err()
}

func (r *transferRepository) Update(ctx context.Context, transfer *entities.Transfer) error {
	query := `
		UPDATE transfers
		SET product_id = $1, from_store_id = $2, to_store_id = $3, quantity = $4, status = $5, notes = $6, updated_at = NOW()
		WHERE id = $7 AND tenant_id = $8
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		transfer.ProductID,
		transfer.FromStoreID,
		transfer.ToStoreID,
		transfer.Quantity,
		transfer.Status,
		transfer.Notes,
		transfer.ID,
		transfer.TenantID,
	).Scan(&transfer.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return entities.ErrTransferNotFound
		}
		return err
	}

	return nil
}

func (r *transferRepository) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	query := `DELETE FROM transfers WHERE id = $1 AND tenant_id = $2`

	result, err := r.pool.Exec(ctx, query, id, tenantID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return entities.ErrTransferNotFound
	}

	return nil
}
