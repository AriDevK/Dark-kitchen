package routes

import (
	"github.com/gin-gonic/gin"

	commonMiddleware "github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
	auth := router.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := auth.Group("")
	protected.Use(commonMiddleware.JWT(jwtSecret))

	protected.GET("/me", authHandler.Me)
}
