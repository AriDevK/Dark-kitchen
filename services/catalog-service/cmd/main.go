package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aridevk/dark-kitchen/packages/go/common/config"
	"github.com/aridevk/dark-kitchen/packages/go/common/database"
	"github.com/aridevk/dark-kitchen/packages/go/common/health"
	"github.com/aridevk/dark-kitchen/packages/go/common/logger"
	"github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/packages/go/common/response"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/handlers"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/repositories"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/routes"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/services"
)

func main() {
	cfg := config.Load()

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	db, err := database.Connect(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	catalogService := services.NewCatalogService(categoryRepo, productRepo)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logging(log))

	router.GET("/health", health.Handler(cfg.AppName))

	router.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{
			"message": "catalog-service is running",
		})
	})

	routes.RegisterRoutes(router, catalogHandler, cfg.JWTSecret)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	log.Info("starting catalog-service",
		zap.String("addr", addr),
		zap.String("env", cfg.AppEnv),
	)

	if err := router.Run(addr); err != nil {
		log.Fatal("failed to start catalog-service", zap.Error(err))
	}
}
