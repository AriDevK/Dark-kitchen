package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/services/auth-service/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, authHandler *handlers.AuthHandler) {
	auth := router.Group("/auth")

	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
