package rest

import (
	authdomain "motico-api/internal/domain/auth"
	authhandler "motico-api/internal/rest/auth"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(authService *authdomain.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(LoggerMiddleware)
	router.Use(RecoveryMiddleware)
	router.Use(chimiddleware.AllowContentType("application/json"))

	authHandler := authhandler.NewHandler(authService)

	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(TenantMiddleware)
			r.Use(AuthMiddleware(authService))

			// Routes will be added here in future phases
		})
	})

	return router
}

