package auth

import (
	"errors"
	"motico-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service struct {
	config *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		config: cfg,
	}
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	TokenType string `json:"token_type"`
}

func (s *Service) GenerateToken(userID, email string) (string, int, error) {
	expirationTime, err := time.ParseDuration(s.config.JWT.ExpirationTime)
	if err != nil {
		return "", 0, err
	}

	expiresAt := time.Now().Add(expirationTime)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.SecretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, int(expirationTime.Seconds()), nil
}

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	// TODO: Implementar validación real contra base de datos
	// Por ahora, validación básica para scaffolding
	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidCredentials
	}

	// Placeholder: En producción, validar contra base de datos
	// userID := "user-id-from-db"
	// token, expiresIn, err := s.GenerateToken(userID, req.Email)
	// if err != nil {
	//     return nil, err
	// }

	// Por ahora, retornamos error indicando que necesita implementación
	return nil, errors.New("login implementation pending: database user validation required")
}

