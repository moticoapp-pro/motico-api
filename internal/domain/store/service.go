package store

import (
	"context"
	"motico-api/config"
	"motico-api/internal/domain/store/entities"
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
	TenantID uuid.UUID
	Name     string
	Address  *string
}

type UpdateRequest struct {
	ID       uuid.UUID
	TenantID uuid.UUID
	Name     *string
	Address  *string
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*entities.Store, error) {
	if err := s.validateName(req.Name); err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByName(ctx, req.TenantID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, entities.ErrStoreNameExists
	}

	store := &entities.Store{
		TenantID: req.TenantID,
		Name:     req.Name,
		Address:  req.Address,
	}

	if err := s.repo.Create(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*entities.Store, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

func (s *Service) List(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*entities.Store, error) {
	if limit <= 0 {
		limit = s.config.Pagination.DefaultLimit
	}
	if limit > s.config.Pagination.MaxLimit {
		limit = s.config.Pagination.MaxLimit
	}
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, tenantID, limit, offset)
}

func (s *Service) Update(ctx context.Context, req UpdateRequest) (*entities.Store, error) {
	store, err := s.repo.GetByID(ctx, req.TenantID, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if err := s.validateName(*req.Name); err != nil {
			return nil, err
		}

		if *req.Name != store.Name {
			exists, err := s.repo.ExistsByName(ctx, req.TenantID, *req.Name)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, entities.ErrStoreNameExists
			}
			store.Name = *req.Name
		}
	}

	if req.Address != nil {
		store.Address = req.Address
	}

	if err := s.repo.Update(ctx, store); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Service) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	hasProducts, err := s.repo.HasProducts(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return entities.ErrStoreHasProducts
	}

	return s.repo.Delete(ctx, tenantID, id)
}

func (s *Service) validateName(name string) error {
	if name == "" {
		return entities.ErrInvalidStoreName
	}
	if len(name) > s.config.Validation.MaxNameLength {
		return entities.ErrInvalidStoreName
	}
	return nil
}
