package category

import (
	"motico-api/config"
	"motico-api/internal/domain/category"
)

type Handler struct {
	service *category.Service
	config  *config.Config
}

func NewHandler(service *category.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}
