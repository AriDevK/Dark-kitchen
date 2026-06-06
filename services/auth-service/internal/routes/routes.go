package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/packages/go/common/health"
	commonMiddleware "github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, appName string, jwtSecret string) {
	auth := router.Group("/auth")

	auth.GET("/health", health.Handler(appName))
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout", authHandler.Logout)

	protected := auth.Group("")
	protected.Use(commonMiddleware.JWT(jwtSecret))

	protected.GET("/me", authHandler.Me)

}
