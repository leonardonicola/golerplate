package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/dto"
	"github.com/leonardonicola/golerplate/pkg/util"
)

type UserHandler struct {
	userService service.UserService
	log         *log.Logger
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
		log:         log.Default(),
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user in the system
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterUserDTO		true	"User registration details"
//	@Success		201		{object}	dto.RegisterResponseDTO	"Successfully created user"
//	@Failure		422		{object}	dto.ErrorResponseDTO	"Validation error"
//	@Failure		500		{object}	dto.ErrorResponseDTO	"Internal server error"
//	@Router			/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterUserDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Print(err)
		c.JSON(http.StatusUnprocessableEntity, util.HandleValidationError(err))
		return
	}

	user, err := h.userService.Create(ctx, req)

	if err != nil {
		h.log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
