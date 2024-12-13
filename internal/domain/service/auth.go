package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/internal/dto"
)

type Claims struct {
	UserID string `json:"id"`
	Email  string `json:"email"`
	// Embedding
	jwt.RegisteredClaims
}

type AuthService interface {
	GenerateToken(user *entity.User) (*dto.TokenResponseDTO, error)
	RefreshToken(refreshToken string) (*dto.TokenResponseDTO, error)
}

type authService struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewAuthService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *authService {
	return &authService{
		accessSecret:  accessSecret,
		accessTTL:     accessTTL,
		refreshSecret: refreshSecret,
		refreshTTL:    refreshTTL,
	}
}

func (s *authService) GenerateToken(user *entity.User) (*dto.TokenResponseDTO, error) {
	// Generate access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.refreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &dto.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) generateToken(user *entity.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})

	return token.SignedString([]byte(s.accessSecret))
}

func (s *authService) refreshToken(user *entity.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	})

	return token.SignedString([]byte(s.refreshSecret))
}

func (s *authService) RefreshToken(refreshToken string) (*dto.TokenResponseDTO, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Extract user information
	userID, ok := claims["id"].(string)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	user := &entity.User{
		ID: userID,
	}

	// Generate access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err = s.refreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Generate new token pair
	return &dto.TokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
