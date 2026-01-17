package transfer

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/stock"
	stockentities "motico-api/internal/domain/stock/entities"
	storedomain "motico-api/internal/domain/store"
	"motico-api/internal/domain/transfer/entities"
	"motico-api/pkg/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo         Repository
	stockService *stock.Service
	storeRepo    storedomain.Repository
	config       *config.Config
	logger       logger.Logger
}

func NewService(repo Repository, stockService *stock.Service, storeRepo storedomain.Repository, cfg *config.Config, log logger.Logger) *Service {
	return &Service{
		repo:         repo,
		stockService: stockService,
		storeRepo:    storeRepo,
		config:       cfg,
		logger:       log,
	}
}

type CreateRequest struct {
	TenantID    uuid.UUID
	ProductID   uuid.UUID
	FromStoreID uuid.UUID
	ToStoreID   uuid.UUID
	Quantity    int
	Notes       *string
}

type UpdateRequest struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ProductID   *uuid.UUID
	FromStoreID *uuid.UUID
	ToStoreID   *uuid.UUID
	Quantity    *int
	Notes       *string
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*entities.Transfer, error) {
	if err := s.validateCreateRequest(req); err != nil {
		s.logger.Warn("Create transfer request validation failed",
			logger.Error(err),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("product_id", req.ProductID.String()))
		return nil, err
	}

	if err := s.validateStoresBelongToTenant(ctx, req.TenantID, req.FromStoreID, req.ToStoreID); err != nil {
		s.logger.Warn("Stores validation failed",
			logger.Error(err),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("from_store_id", req.FromStoreID.String()),
			logger.String("to_store_id", req.ToStoreID.String()))
		return nil, err
	}

	stock, err := s.stockService.GetByProductIDAndStoreID(ctx, req.TenantID, req.ProductID, req.FromStoreID)
	if err != nil {
		if err == stockentities.ErrStockNotFound {
			s.logger.Warn("Stock not found for product in store",
				logger.String("tenant_id", req.TenantID.String()),
				logger.String("product_id", req.ProductID.String()),
				logger.String("store_id", req.FromStoreID.String()))
			return nil, entities.ErrInsufficientStock
		}
		s.logger.Error("Failed to get stock for product in store",
			logger.Error(err),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("product_id", req.ProductID.String()),
			logger.String("store_id", req.FromStoreID.String()))
		return nil, err
	}

	available := stock.AvailableQuantity()
	if available < req.Quantity {
		s.logger.Warn("Insufficient stock for transfer",
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("product_id", req.ProductID.String()),
			logger.String("store_id", req.FromStoreID.String()),
			logger.Int("available", available),
			logger.Int("requested", req.Quantity))
		return nil, entities.ErrInsufficientStock
	}

	transfer := &entities.Transfer{
		TenantID:    req.TenantID,
		ProductID:   req.ProductID,
		FromStoreID: req.FromStoreID,
		ToStoreID:   req.ToStoreID,
		Quantity:    req.Quantity,
		Status:      entities.TransferStatusPending,
		Notes:       req.Notes,
	}

	if err := s.repo.Create(ctx, transfer); err != nil {
		s.logger.Error("Failed to create transfer in repository",
			logger.Error(err),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("product_id", req.ProductID.String()))
		return nil, err
	}

	if err := s.stockService.Reserve(ctx, req.TenantID, req.ProductID, req.FromStoreID, req.Quantity); err != nil {
		s.logger.Error("Failed to reserve stock for transfer",
			logger.Error(err),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("product_id", req.ProductID.String()),
			logger.String("store_id", req.FromStoreID.String()),
			logger.String("transfer_id", transfer.ID.String()),
			logger.Int("quantity", req.Quantity))
		return nil, err
	}

	s.logger.Info("Transfer created successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", req.TenantID.String()),
		logger.String("product_id", req.ProductID.String()),
		logger.String("from_store_id", req.FromStoreID.String()),
		logger.String("to_store_id", req.ToStoreID.String()),
		logger.Int("quantity", req.Quantity),
		logger.String("status", string(transfer.Status)))

	return transfer, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			s.logger.Warn("Transfer not found",
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		} else {
			s.logger.Error("Failed to get transfer by ID",
				logger.Error(err),
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		}
		return nil, err
	}

	return transfer, nil
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, status *entities.TransferStatus, storeID *uuid.UUID, limit, offset int) ([]*entities.Transfer, error) {
	if limit <= 0 {
		limit = s.config.Pagination.DefaultLimit
	}
	if limit > s.config.Pagination.MaxLimit {
		limit = s.config.Pagination.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	transfers, err := s.repo.List(ctx, tenantID, status, storeID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to list transfers",
			logger.Error(err),
			logger.String("tenant_id", tenantID.String()))
		return nil, err
	}

	return transfers, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, req.TenantID, req.ID)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			s.logger.Warn("Transfer not found for update",
				logger.String("transfer_id", req.ID.String()),
				logger.String("tenant_id", req.TenantID.String()))
		} else {
			s.logger.Error("Failed to get transfer for update",
				logger.Error(err),
				logger.String("transfer_id", req.ID.String()),
				logger.String("tenant_id", req.TenantID.String()))
		}
		return nil, err
	}

	if !transfer.CanUpdate() {
		s.logger.Warn("Transfer cannot be updated, not in pending status",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("current_status", string(transfer.Status)))
		return nil, entities.ErrTransferNotPending
	}

	if req.ProductID != nil {
		transfer.ProductID = *req.ProductID
	}

	if req.FromStoreID != nil {
		transfer.FromStoreID = *req.FromStoreID
	}

	if req.ToStoreID != nil {
		transfer.ToStoreID = *req.ToStoreID
	}

	if req.FromStoreID != nil || req.ToStoreID != nil {
		fromStoreID := transfer.FromStoreID
		toStoreID := transfer.ToStoreID
		if req.FromStoreID != nil {
			fromStoreID = *req.FromStoreID
		}
		if req.ToStoreID != nil {
			toStoreID = *req.ToStoreID
		}
		if err := s.validateStoresBelongToTenant(ctx, req.TenantID, fromStoreID, toStoreID); err != nil {
			s.logger.Warn("Stores validation failed during update",
				logger.Error(err),
				logger.String("transfer_id", transfer.ID.String()),
				logger.String("tenant_id", req.TenantID.String()),
				logger.String("from_store_id", fromStoreID.String()),
				logger.String("to_store_id", toStoreID.String()))
			return nil, err
		}
	}

	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			s.logger.Warn("Invalid quantity in transfer update",
				logger.String("transfer_id", transfer.ID.String()),
				logger.String("tenant_id", req.TenantID.String()),
				logger.Int("quantity", *req.Quantity))
			return nil, entities.ErrInvalidQuantity
		}
		transfer.Quantity = *req.Quantity
	}

	if req.Notes != nil {
		transfer.Notes = req.Notes
	}

	if transfer.FromStoreID == transfer.ToStoreID {
		s.logger.Warn("Invalid transfer stores, from and to stores are the same",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", req.TenantID.String()),
			logger.String("store_id", transfer.FromStoreID.String()))
		return nil, entities.ErrInvalidTransferStores
	}

	if err := s.repo.Update(ctx, transfer); err != nil {
		s.logger.Error("Failed to update transfer in repository",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", req.TenantID.String()))
		return nil, err
	}

	s.logger.Info("Transfer updated successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", req.TenantID.String()),
		logger.String("status", string(transfer.Status)),
		logger.Int("quantity", transfer.Quantity))

	return transfer, nil
}

func (s *Service) Complete(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			s.logger.Warn("Transfer not found for completion",
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		} else {
			s.logger.Error("Failed to get transfer for completion",
				logger.Error(err),
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		}
		return nil, err
	}

	if transfer.IsCompleted() {
		s.logger.Warn("Transfer already completed",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, entities.ErrTransferAlreadyCompleted
	}

	if transfer.IsCancelled() {
		s.logger.Warn("Transfer already cancelled, cannot complete",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, entities.ErrTransferAlreadyCancelled
	}

	transfer.Status = entities.TransferStatusCompleted
	if err := s.repo.Update(ctx, transfer); err != nil {
		s.logger.Error("Failed to update transfer status to completed",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, err
	}

	if err := s.stockService.Release(ctx, tenantID, transfer.ProductID, transfer.FromStoreID, transfer.Quantity); err != nil {
		s.logger.Error("Failed to release stock for completed transfer",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()),
			logger.String("product_id", transfer.ProductID.String()),
			logger.String("from_store_id", transfer.FromStoreID.String()),
			logger.Int("quantity", transfer.Quantity))
		return nil, err
	}

	fromAdjustReq := stock.AdjustRequest{
		TenantID:  tenantID,
		ProductID: transfer.ProductID,
		StoreID:   transfer.FromStoreID,
		Amount:    -transfer.Quantity,
	}
	if _, err := s.stockService.Adjust(ctx, fromAdjustReq); err != nil {
		s.logger.Error("Failed to reduce stock in source store",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()),
			logger.String("product_id", transfer.ProductID.String()),
			logger.String("from_store_id", transfer.FromStoreID.String()),
			logger.Int("quantity", transfer.Quantity))
		return nil, err
	}

	toAdjustReq := stock.AdjustRequest{
		TenantID:  tenantID,
		ProductID: transfer.ProductID,
		StoreID:   transfer.ToStoreID,
		Amount:    transfer.Quantity,
	}
	if _, err := s.stockService.Adjust(ctx, toAdjustReq); err != nil {
		s.logger.Error("Failed to adjust stock in destination store",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()),
			logger.String("product_id", transfer.ProductID.String()),
			logger.String("to_store_id", transfer.ToStoreID.String()),
			logger.Int("quantity", transfer.Quantity))
		return nil, err
	}

	s.logger.Info("Transfer completed successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", transfer.ProductID.String()),
		logger.String("from_store_id", transfer.FromStoreID.String()),
		logger.String("to_store_id", transfer.ToStoreID.String()),
		logger.Int("quantity", transfer.Quantity))

	return transfer, nil
}

func (s *Service) Cancel(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			s.logger.Warn("Transfer not found for cancellation",
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		} else {
			s.logger.Error("Failed to get transfer for cancellation",
				logger.Error(err),
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		}
		return nil, err
	}

	if transfer.IsCompleted() {
		s.logger.Warn("Transfer already completed, cannot cancel",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, entities.ErrTransferAlreadyCompleted
	}

	if transfer.IsCancelled() {
		s.logger.Warn("Transfer already cancelled",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, entities.ErrTransferAlreadyCancelled
	}

	transfer.Status = entities.TransferStatusCancelled
	if err := s.repo.Update(ctx, transfer); err != nil {
		s.logger.Error("Failed to update transfer status to cancelled",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()))
		return nil, err
	}

	if err := s.stockService.Release(ctx, tenantID, transfer.ProductID, transfer.FromStoreID, transfer.Quantity); err != nil {
		s.logger.Error("Failed to release stock for cancelled transfer",
			logger.Error(err),
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()),
			logger.String("product_id", transfer.ProductID.String()),
			logger.Int("quantity", transfer.Quantity))
		return nil, err
	}

	s.logger.Info("Transfer cancelled successfully",
		logger.String("transfer_id", transfer.ID.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", transfer.ProductID.String()),
		logger.String("from_store_id", transfer.FromStoreID.String()),
		logger.String("to_store_id", transfer.ToStoreID.String()),
		logger.Int("quantity", transfer.Quantity))

	return transfer, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		if err == entities.ErrTransferNotFound {
			s.logger.Warn("Transfer not found for deletion",
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		} else {
			s.logger.Error("Failed to get transfer for deletion",
				logger.Error(err),
				logger.String("transfer_id", id.String()),
				logger.String("tenant_id", tenantID.String()))
		}
		return err
	}

	if !transfer.CanDelete() {
		s.logger.Warn("Transfer cannot be deleted, not in pending status",
			logger.String("transfer_id", transfer.ID.String()),
			logger.String("tenant_id", tenantID.String()),
			logger.String("current_status", string(transfer.Status)))
		return entities.ErrTransferNotPending
	}

	err = s.repo.Delete(ctx, tenantID, id)
	if err != nil {
		s.logger.Error("Failed to delete transfer from repository",
			logger.Error(err),
			logger.String("transfer_id", id.String()),
			logger.String("tenant_id", tenantID.String()))
		return err
	}

	s.logger.Info("Transfer deleted successfully",
		logger.String("transfer_id", id.String()),
		logger.String("tenant_id", tenantID.String()),
		logger.String("product_id", transfer.ProductID.String()),
		logger.Int("quantity", transfer.Quantity))

	return nil
}

func (s *Service) validateCreateRequest(req CreateRequest) error {
	if req.FromStoreID == req.ToStoreID {
		return entities.ErrInvalidTransferStores
	}

	if req.Quantity <= 0 {
		return entities.ErrInvalidQuantity
	}

	return nil
}

func (s *Service) validateStoresBelongToTenant(ctx context.Context, tenantID, fromStoreID, toStoreID uuid.UUID) error {
	fromStore, err := s.storeRepo.GetByID(ctx, tenantID, fromStoreID)
	if err != nil {
		return entities.ErrInvalidTransferStores
	}
	if fromStore.TenantID != tenantID {
		return entities.ErrInvalidTransferStores
	}

	toStore, err := s.storeRepo.GetByID(ctx, tenantID, toStoreID)
	if err != nil {
		return entities.ErrInvalidTransferStores
	}
	if toStore.TenantID != tenantID {
		return entities.ErrInvalidTransferStores
	}

	return nil
}
