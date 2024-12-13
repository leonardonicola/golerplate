package config

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonardonicola/golerplate/docs"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/handler"
	"github.com/leonardonicola/golerplate/internal/infra/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(pool *pgxpool.Pool) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// User.
	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Auth.
	authService := service.NewAuthService(os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"), time.Hour, 2*time.Hour)
	authHandler := handler.NewAuthHandler(userService, authService)

	// jwtMiddleware := middleware.NewJWTAuthMiddleware(os.Getenv("JWT_SECRET"))

	docs.SwaggerInfo.BasePath = "/api"
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/docs/*any", func(c *gin.Context) {
			if c.Param("any") == "/" || c.Param("any") == "" {
				c.Redirect(http.StatusTemporaryRedirect, "/api/docs/index.html")
				return
			}

			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		})
		public.POST("/refresh", authHandler.Refresh)
	}

	// protected := r.Group("/api", jwtMiddleware.AuthRequired())

	return r

	// r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
