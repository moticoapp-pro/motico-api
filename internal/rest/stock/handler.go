package stock

import (
	"motico-api/config"
	"motico-api/internal/domain/stock"
)

type Handler struct {
	service *stock.Service
	config  *config.Config
}

func NewHandler(service *stock.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}
