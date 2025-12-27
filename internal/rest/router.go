package rest

import (
	authdomain "motico-api/internal/domain/auth"
	authhandler "motico-api/internal/rest/auth"
	"motico-api/internal/rest/category"
	"motico-api/internal/rest/product"
	"motico-api/internal/rest/stock"
	"motico-api/internal/rest/store"
	"motico-api/internal/rest/transfer"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

type RouterDependencies struct {
	AuthService     *authdomain.Service
	CategoryHandler *category.Handler
	StoreHandler    *store.Handler
	ProductHandler  *product.Handler
	StockHandler    *stock.Handler
	TransferHandler *transfer.Handler
}

func NewRouter(deps RouterDependencies) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(LoggerMiddleware)
	router.Use(RecoveryMiddleware)
	router.Use(chimiddleware.AllowContentType("application/json"))

	authHandler := authhandler.NewHandler(deps.AuthService)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(TenantMiddleware)
			r.Use(AuthMiddleware(deps.AuthService))

			r.Route("/categories", func(r chi.Router) {
				r.Get("/", deps.CategoryHandler.List)
				r.Get("/{id}", deps.CategoryHandler.GetByID)
				r.Post("/", deps.CategoryHandler.Create)
				r.Put("/{id}", deps.CategoryHandler.Update)
				r.Patch("/{id}", deps.CategoryHandler.ParcialUpdate)
				r.Delete("/{id}", deps.CategoryHandler.Remove)
			})

			r.Route("/stores", func(r chi.Router) {
				r.Get("/", deps.StoreHandler.List)
				r.Get("/{id}", deps.StoreHandler.GetByID)
				r.Post("/", deps.StoreHandler.Create)
				r.Put("/{id}", deps.StoreHandler.Update)
				r.Patch("/{id}", deps.StoreHandler.ParcialUpdate)
				r.Delete("/{id}", deps.StoreHandler.Remove)
			})

			r.Route("/products", func(r chi.Router) {
				r.Get("/", deps.ProductHandler.List)
				r.Get("/{id}", deps.ProductHandler.GetByID)
				r.Post("/", deps.ProductHandler.Create)
				r.Put("/{id}", deps.ProductHandler.Update)
				r.Patch("/{id}", deps.ProductHandler.ParcialUpdate)
				r.Delete("/{id}", deps.ProductHandler.Remove)

				r.Route("/{id}/stock", func(r chi.Router) {
					r.Get("/", deps.StockHandler.GetByID)
					r.Put("/", deps.StockHandler.Update)
					r.Patch("/", deps.StockHandler.Adjust)
				})
			})

			r.Route("/transfers", func(r chi.Router) {
				r.Get("/", deps.TransferHandler.List)
				r.Get("/{id}", deps.TransferHandler.GetByID)
				r.Post("/", deps.TransferHandler.Create)
				r.Put("/{id}", deps.TransferHandler.Update)
				r.Patch("/{id}/complete", deps.TransferHandler.Complete)
				r.Patch("/{id}/cancel", deps.TransferHandler.Cancel)
				r.Delete("/{id}", deps.TransferHandler.Remove)
			})
		})
	})

	return router
}
