package product

import (
	"motico-api/config"
	"motico-api/internal/domain/product"
	"motico-api/internal/domain/stock"
)

type Handler struct {
	service      *product.Service
	stockService *stock.Service
	config       *config.Config
}

func NewHandler(service *product.Service, stockService *stock.Service, cfg *config.Config) *Handler {
	return &Handler{
		service:      service,
		stockService: stockService,
		config:       cfg,
	}
}
