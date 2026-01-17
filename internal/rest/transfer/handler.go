package transfer

import (
	"motico-api/config"
	"motico-api/internal/domain/transfer"
	"motico-api/pkg/logger"
)

type Handler struct {
	service *transfer.Service
	config  *config.Config
	logger  logger.Logger
}

func NewHandler(service *transfer.Service, cfg *config.Config, log logger.Logger) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
		logger:  log,
	}
}
