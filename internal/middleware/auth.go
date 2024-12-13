package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/pkg/constants"
)

type JWTAuthMiddleware struct {
	secret string
}

var (
	ErrMissingHeader    = errors.New(constants.ErrMsgMissingHeader)
	ErrInvalidToken     = errors.New(constants.ErrMsgInvalidToken)
	ErrInvalidTokenType = errors.New(constants.ErrMsgInvalidTokenType)
)

func NewJWTAuthMiddleware(secret string) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		secret: secret,
	}
}

func (m *JWTAuthMiddleware) AuthRequired() gin.HandlerFunc {

	return func(c *gin.Context) {
		token, err := m.extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponseDTO{
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := m.validateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponseDTO{
				Message: constants.ErrMsgInvalidToken,
			})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)

		c.Next()
	}
}

func (m *JWTAuthMiddleware) extractToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")

	if header == "" {
		return "", ErrMissingHeader
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return parts[1], nil
}

func (m *JWTAuthMiddleware) validateToken(tokenString string) (*service.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, service.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*service.Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
