package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aridevk/dark-kitchen/packages/go/common/config"
	"github.com/aridevk/dark-kitchen/packages/go/common/database"
	"github.com/aridevk/dark-kitchen/packages/go/common/health"
	"github.com/aridevk/dark-kitchen/packages/go/common/logger"
	"github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/packages/go/common/response"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/handlers"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/repositories"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/routes"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/services"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := database.Connect(cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logging(log))

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTAccessTokenTTLMinutes)
	authHandler := handlers.NewAuthHandler(authService)
	routes.RegisterRoutes(router, authHandler)

	router.GET("/health", health.Handler(cfg.AppName))

	router.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{
			"message": "auth-service is running",
		})
	})

	router.GET("/users", func(c *gin.Context) {
		users, err := userRepo.GetUsers()
		if err != nil {
			response.Error(
				c,
				http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				"Failed to fetch users",
			)
			return
		}

		response.OK(c, gin.H{
			"users": users,
			"count": len(users),
		})
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	log.Info("starting auth-service",
		zap.String("addr", addr),
		zap.String("env", cfg.AppEnv),
	)

	if err := router.Run(addr); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}
