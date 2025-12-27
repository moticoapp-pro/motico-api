package auth

import (
	"encoding/json"
	authdomain "motico-api/internal/domain/auth"
	"motico-api/internal/rest/response"
	"net/http"
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
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest   true  "Login credentials"
// @Success      200      {object}  LoginResponse  "Successfully authenticated"
// @Failure      400      {object}  map[string]interface{}  "Invalid request body"
// @Failure      401      {object}  map[string]interface{}  "Invalid credentials"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/login [post]
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
