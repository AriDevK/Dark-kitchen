package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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

	router.GET("/health", health.Handler(cfg.AppName))

	router.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{
			"message": "auth-service is running",
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
