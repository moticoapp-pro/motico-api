package stock

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/stock/entities"
	"motico-api/pkg/logger"

	"github.com/google/uuid"
)

type Service struct {
	repo   Repository
	config *config.Config
	logger logger.Logger
}

func NewService(repo Repository, cfg *config.Config, log logger.Logger) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
		logger: log,
	}
}

type UpdateRequest struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Quantity  int
}

type AdjustRequest struct {
	TenantID  uuid.UUID
	ProductID uuid.UUID
	Amount    int
}

func (s *Service) GetByProductID(ctx context.Context, tenantID, productID uuid.UUID) (*entities.Stock, error) {
	stock, err := s.repo.GetByProductID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Stock, error) {
	if req.Quantity < 0 {
		return nil, entities.ErrInvalidQuantity
	}

	stock, err := s.repo.GetByProductID(ctx, req.TenantID, req.ProductID)
	if err != nil {
		if err == entities.ErrStockNotFound {
			stock = &entities.Stock{
				TenantID:         req.TenantID,
				ProductID:        req.ProductID,
				Quantity:         req.Quantity,
				ReservedQuantity: 0,
			}
			if err := s.repo.Create(ctx, stock); err != nil {
				return nil, err
			}
			return stock, nil
		}
		return nil, err
	}

	if req.Quantity < stock.ReservedQuantity {
		return nil, entities.ErrInvalidReservedAmount
	}

	stock.Quantity = req.Quantity
	if err := s.repo.Update(ctx, stock); err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *Service) Adjust(ctx context.Context, req AdjustRequest) (*entities.Stock, error) {
	stock, err := s.repo.GetByProductID(ctx, req.TenantID, req.ProductID)
	if err != nil {
		if err == entities.ErrStockNotFound {
			if req.Amount < 0 {
				return nil, entities.ErrInsufficientStock
			}
			stock = &entities.Stock{
				TenantID:         req.TenantID,
				ProductID:        req.ProductID,
				Quantity:         req.Amount,
				ReservedQuantity: 0,
			}
			if err := s.repo.Create(ctx, stock); err != nil {
				return nil, err
			}
			return stock, nil
		}
		return nil, err
	}

	newQuantity := stock.Quantity + req.Amount
	if newQuantity < 0 {
		return nil, entities.ErrInsufficientStock
	}

	if newQuantity < stock.ReservedQuantity {
		return nil, entities.ErrInvalidReservedAmount
	}

	stock.Quantity = newQuantity
	if err := s.repo.Update(ctx, stock); err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *Service) Reserve(ctx context.Context, tenantID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return entities.ErrInvalidQuantity
	}

	return s.repo.Reserve(ctx, tenantID, productID, quantity)
}

func (s *Service) Release(ctx context.Context, tenantID, productID uuid.UUID, quantity int) error {
	if quantity <= 0 {
		return entities.ErrInvalidQuantity
	}

	return s.repo.Release(ctx, tenantID, productID, quantity)
}
