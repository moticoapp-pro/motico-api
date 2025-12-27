package category

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/category/entities"
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
	Name        string
	Description *string
}

type UpdateRequest struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        *string
	Description *string
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*entities.Category, error) {
	if err := s.validateName(req.Name); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByName(ctx, req.TenantID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrCategoryNameExists
	}

	category := &entities.Category{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Category, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Category, error) {
	if limit <= 0 {
		limit = s.config.Pagination.DefaultLimit
	}
	if limit > s.config.Pagination.MaxLimit {
		limit = s.config.Pagination.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	categories, err := s.repo.List(ctx, tenantID, limit, offset)
	if err != nil {
		s.logger.Error("Error listing categories", logger.Error(err), logger.String("tenant_id", tenantID.String()))
		return nil, err
	}

	return categories, nil
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Category, error) {
	category, err := s.repo.GetByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if err := s.validateName(*req.Name); err != nil {
			return nil, err
		}

		if *req.Name != category.Name {
			exists, err := s.repo.ExistsByName(ctx, req.TenantID, *req.Name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entities.ErrCategoryNameExists
			}
			category.Name = *req.Name
		}
	}

	if req.Description != nil {
		category.Description = req.Description
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	hasProducts, err := s.repo.HasProducts(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return entities.ErrCategoryHasProducts
	}

	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) validateName(name string) error {
	if name == "" {
		return entities.ErrInvalidCategoryName
	}
	if len(name) > s.config.Validation.MaxNameLength {
		return entities.ErrInvalidCategoryName
	}
	return nil
}
