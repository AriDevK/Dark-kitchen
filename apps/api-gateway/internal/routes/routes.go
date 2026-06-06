package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/apps/api-gateway/internal/proxy"
)

func RegisterRoutes(router *gin.Engine, authServiceURL string, catalogServiceURL string) error {
	authProxy, err := proxy.NewReverseProxy(authServiceURL, "/api/")
	if err != nil {
		return err
	}

	catalogProxy, err := proxy.NewReverseProxy(catalogServiceURL, "/api/")
	if err != nil {
		return err
	}

	router.Any("/api/auth", authProxy)
	router.Any("/api/auth/*path", authProxy)
	router.Any("/api/catalog", catalogProxy)
	router.Any("/api/catalog/*path", catalogProxy)

	return nil
}
