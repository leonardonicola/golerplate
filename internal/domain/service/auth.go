package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/leonardonicola/golerplate/internal/domain/entity"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/pkg/constants"
)

type Claims struct {
	UserID string `json:"id"`
	Type   string `json:"type"`
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
	accessToken, err := s.accessToken(user)
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

func (s *authService) accessToken(user *entity.User) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		Type:   "access",
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
		Type:   "refresh",
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
			return nil, errors.New(constants.ErrMsgInvalidToken)
		}
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, errors.New(constants.ErrMsgInvalidToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New(constants.ErrMsgInvalidToken)
	}

	// Extract user information
	userID, ok := claims["id"].(string)

	if !ok {
		return nil, errors.New(constants.ErrMsgInvalidToken)
	}

	if claims["type"] != "refresh" {
		return nil, errors.New(constants.ErrMsgInvalidToken)
	}

	user := &entity.User{
		ID: userID,
	}

	// Generate access token
	accessToken, err := s.accessToken(user)
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
