package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/pkg/util"
)

type AuthHandler struct {
	userService  service.UserService
	tokenService service.AuthService
	log          *log.Logger
}

func NewAuthHandler(us service.UserService, ts service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService:  us,
		tokenService: ts,
		log:          log.Default(),
	}
}

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate a user and return access tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequestDTO		true	"Login credentials"
//	@Success		200		{object}	dto.TokenResponseDTO	"Successfully authenticated"
//	@Failure		400		{object}	dto.ErrorResponseDTO	"Bad request"
//	@Failure		401		{object}	dto.ErrorResponseDTO	"Unauthorized"
//	@Router			/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Print(err)
		c.JSON(http.StatusBadRequest, util.HandleValidationError(err))
		return
	}

	user, err := h.userService.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		h.log.Printf("USER SERVICE: %s", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	token, err := h.tokenService.GenerateToken(user)

	if err != nil {
		h.log.Printf("AUTH SERVICE: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Refresh Token godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh access token using refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshRequestDTO	true	"Refresh token request"
//	@Success		200		{object}	dto.TokenResponseDTO	"Successfully refreshed tokens"
//	@Failure		400		{object}	dto.ErrorResponseDTO	"Bad request"
//	@Failure		401		{object}	dto.ErrorResponseDTO	"Unauthorized"
//	@Router			/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := h.tokenService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}
