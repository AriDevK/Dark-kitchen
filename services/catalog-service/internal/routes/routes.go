package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, catalogHandler *handlers.CatalogHandler, jwtSecret string) {
	catalog := router.Group("/catalog")

	// Public routes
	catalog.GET("/categories", catalogHandler.GetCategories)
	catalog.GET("/products", catalogHandler.GetProducts)
	catalog.GET("/products/:id", catalogHandler.GetProductByID)
	// Admin routes
	admin := catalog.Group("")
	admin.Use(middleware.JWT(jwtSecret))
	admin.Use(middleware.RequireRole("admin"))
	admin.POST("/categories", catalogHandler.CreateCategory)
	admin.POST("/products", catalogHandler.CreateProduct)
	admin.PATCH("/products/:id", catalogHandler.UpdateProduct)
	admin.DELETE("/products/:id", catalogHandler.DeleteProduct)
}
