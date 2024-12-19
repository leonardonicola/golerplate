package config

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonardonicola/golerplate/docs"
	"github.com/leonardonicola/golerplate/internal/domain/service"
	"github.com/leonardonicola/golerplate/internal/handler"
	"github.com/leonardonicola/golerplate/internal/infra/repository"
	"github.com/leonardonicola/golerplate/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(pool *pgxpool.Pool) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.TracingMiddleware())

	// User.
	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Auth.
	accessSecret, accessExists := os.LookupEnv("ACCESS_SECRET")
	refreshSecret, refreshExists := os.LookupEnv("REFRESH_SECRET")
	if !accessExists || !refreshExists {
		log.Panic("SECRETS variables are not defined correctly.")
	}

	authService := service.NewAuthService(accessSecret, refreshSecret, time.Hour, 2*time.Hour)
	authHandler := handler.NewAuthHandler(userService, authService)

	jwtMiddleware := middleware.NewJWTAuthMiddleware(accessSecret)

	docs.SwaggerInfo.BasePath = "/api"
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.Refresh)
	}

	protected := r.Group("/api", jwtMiddleware.AuthRequired())
	{
		protected.GET("/docs/*any", func(c *gin.Context) {
			if c.Param("any") == "/" || c.Param("any") == "" {
				c.Redirect(http.StatusTemporaryRedirect, "/api/docs/index.html")
				return
			}

			ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
		})
	}

	return r
}
