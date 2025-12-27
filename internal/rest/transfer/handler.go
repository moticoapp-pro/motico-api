package transfer

import (
	"motico-api/config"
	"motico-api/internal/domain/transfer"
)

type Handler struct {
	service *transfer.Service
	config  *config.Config
}

func NewHandler(service *transfer.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}
