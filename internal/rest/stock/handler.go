package stock

import (
	"motico-api/config"
	"motico-api/internal/domain/product"
	"motico-api/internal/domain/stock"
)

type Handler struct {
	service      *stock.Service
	productService *product.Service
	config       *config.Config
}

func NewHandler(service *stock.Service, productService *product.Service, cfg *config.Config) *Handler {
	return &Handler{
		service:        service,
		productService: productService,
		config:         cfg,
	}
}
