package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aridevk/dark-kitchen/apps/api-gateway/internal/routes"
	"github.com/aridevk/dark-kitchen/packages/go/common/config"
	"github.com/aridevk/dark-kitchen/packages/go/common/health"
	"github.com/aridevk/dark-kitchen/packages/go/common/logger"
	"github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/packages/go/common/response"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logging(log))

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: true,
	}))

	router.GET("/health", health.Handler(cfg.AppName))

	router.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{
			"message": "api-gateway is running",
		})
	})

	if err := routes.RegisterRoutes(router, cfg.AuthServiceURL, cfg.CatalogServiceURL); err != nil {
		log.Fatal("failed to register gateway routes", zap.Error(err))
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	log.Info("starting api-gateway",
		zap.String("addr", addr),
		zap.String("env", cfg.AppEnv),
		zap.String("auth_service_url", cfg.AuthServiceURL),
		zap.String("catalog_service_url", cfg.CatalogServiceURL),
	)

	if err := router.Run(addr); err != nil {
		log.Fatal("failed to start api-gateway", zap.Error(err))
	}
}
