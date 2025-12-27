package rest

import (
	"context"
	authdomain "motico-api/internal/domain/auth"
	"motico-api/internal/rest/response"
	ctxpkg "motico-api/pkg/context"
	"net/http"
	"strings"
)

func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get("X-Tenant-ID")
		if tenantID == "" {
			response.Error(w, http.StatusBadRequest, "X-Tenant-ID header is required", nil)
			return
		}

		ctx := context.WithValue(r.Context(), ctxpkg.TenantIDKey, tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(authService *authdomain.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "authorization header is required", nil)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(w, http.StatusUnauthorized, "invalid authorization header format", nil)
				return
			}

			tokenString := parts[1]
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), ctxpkg.UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar logging con el logger configurado
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response.Error(w, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
