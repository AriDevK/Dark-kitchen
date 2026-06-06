package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/apps/api-gateway/internal/proxy"
)

func RegisterRoutes(router *gin.Engine, authServiceURL string) error {
	authProxy, err := proxy.NewReverseProxy(authServiceURL, "/api/")
	if err != nil {
		return err
	}

	router.Any("/api/auth", authProxy)
	router.Any("/api/auth/*path", authProxy)

	return nil
}
