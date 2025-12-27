package product

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/product/entities"
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

type CreateRequest struct {
	TenantID    uuid.UUID
	StoreID     uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Description *string
	SKU         *string
	Price       *float64
}

type UpdateRequest struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	StoreID     *uuid.UUID
	CategoryID  *uuid.UUID
	Name        *string
	Description *string
	SKU         *string
	Price       *float64
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*entities.Product, error) {
	if err := s.validateName(req.Name); err != nil {
		return nil, err
	}

	if req.SKU != nil && *req.SKU != "" {
		exists, err := s.repo.ExistsBySKU(ctx, req.TenantID, req.StoreID, *req.SKU)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, entities.ErrProductSKUExists
		}
	}

	product := &entities.Product{
		TenantID:    req.TenantID,
		StoreID:     req.StoreID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Product, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, storeID, categoryID *uuid.UUID, limit, offset int) ([]*entities.Product, error) {
	if limit <= 0 {
		limit = s.config.Pagination.DefaultLimit
	}
	if limit > s.config.Pagination.MaxLimit {
		limit = s.config.Pagination.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, tenantID, storeID, categoryID, limit, offset)
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Product, error) {
	product, err := s.repo.GetByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if err := s.validateName(*req.Name); err != nil {
			return nil, err
		}
		product.Name = *req.Name
	}

	if req.Description != nil {
		product.Description = req.Description
	}

	if req.SKU != nil {
		if *req.SKU != "" && (*req.SKU != getStringValue(product.SKU)) {
			exists, err := s.repo.ExistsBySKU(ctx, req.TenantID, product.StoreID, *req.SKU)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entities.ErrProductSKUExists
			}
		}
		product.SKU = req.SKU
	}

	if req.Price != nil {
		product.Price = req.Price
	}

	if req.StoreID != nil {
		product.StoreID = *req.StoreID
	}

	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	hasStock, err := s.repo.HasStock(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasStock {
		return entities.ErrProductHasStock
	}

	hasTransfers, err := s.repo.HasTransfers(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasTransfers {
		return entities.ErrProductHasTransfers
	}

	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) validateName(name string) error {
	if name == "" {
		return entities.ErrInvalidProductName
	}
	if len(name) > s.config.Validation.MaxNameLength {
		return entities.ErrInvalidProductName
	}
	return nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
