package transfer

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/stock"
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
		return nil, err
	}

	if err := s.validateStoresBelongToTenant(ctx, req.TenantID, req.FromStoreID, req.ToStoreID); err != nil {
		return nil, err
	}

	stock, err := s.stockService.GetByProductID(ctx, req.TenantID, req.ProductID)
	if err != nil {
		return nil, err
	}

	available := stock.AvailableQuantity()
	if available < req.Quantity {
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
		return nil, err
	}

	if err := s.stockService.Reserve(ctx, req.TenantID, req.ProductID, req.Quantity); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	return s.repo.GetByID(ctx, tenantID, id)
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

	return s.repo.List(ctx, tenantID, status, storeID, limit, offset)
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if !transfer.CanUpdate() {
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
			return nil, err
		}
	}

	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			return nil, entities.ErrInvalidQuantity
		}
		transfer.Quantity = *req.Quantity
	}

	if req.Notes != nil {
		transfer.Notes = req.Notes
	}

	if transfer.FromStoreID == transfer.ToStoreID {
		return nil, entities.ErrInvalidTransferStores
	}

	if err := s.repo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *Service) Complete(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if transfer.IsCompleted() {
		return nil, entities.ErrTransferAlreadyCompleted
	}

	if transfer.IsCancelled() {
		return nil, entities.ErrTransferAlreadyCancelled
	}

	transfer.Status = entities.TransferStatusCompleted
	if err := s.repo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	if err := s.stockService.Release(ctx, tenantID, transfer.ProductID, transfer.Quantity); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *Service) Cancel(ctx context.Context, tenantID, id uuid.UUID) (*entities.Transfer, error) {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if transfer.IsCompleted() {
		return nil, entities.ErrTransferAlreadyCompleted
	}

	if transfer.IsCancelled() {
		return nil, entities.ErrTransferAlreadyCancelled
	}

	transfer.Status = entities.TransferStatusCancelled
	if err := s.repo.Update(ctx, transfer); err != nil {
		return nil, err
	}

	if err := s.stockService.Release(ctx, tenantID, transfer.ProductID, transfer.Quantity); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	transfer, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}

	if !transfer.CanDelete() {
		return entities.ErrTransferNotPending
	}

	return s.repo.Delete(ctx, tenantID, id)
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
