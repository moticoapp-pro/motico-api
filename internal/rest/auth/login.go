package auth

import (
	"encoding/json"
	"net/http"
	authdomain "motico-api/internal/domain/auth"
	"motico-api/internal/rest/response"
)

type Handler struct {
	authService *authdomain.Service
}

func NewHandler(authService *authdomain.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	authReq := authdomain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	loginResp, err := h.authService.Login(authReq)
	if err != nil {
		if err == authdomain.ErrInvalidCredentials {
			response.Error(w, http.StatusUnauthorized, "invalid credentials", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.JSON(w, http.StatusOK, loginResp)
}

