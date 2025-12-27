package store

import (
	"motico-api/config"
	"motico-api/internal/domain/store"
)

type Handler struct {
	service *store.Service
	config  *config.Config
}

func NewHandler(service *store.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		config:  cfg,
	}
}
