package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/handlers"
)

func RegisterRoutes(router *gin.Engine, catalogHandler *handlers.CatalogHandler) {
	catalog := router.Group("/catalog")

	catalog.GET("/categories", catalogHandler.GetCategories)
	catalog.POST("/categories", catalogHandler.CreateCategory)

	catalog.GET("/products", catalogHandler.GetProducts)
	catalog.GET("/products/:id", catalogHandler.GetProductByID)
	catalog.POST("/products", catalogHandler.CreateProduct)
	catalog.PATCH("/products/:id", catalogHandler.UpdateProduct)
	catalog.DELETE("/products/:id", catalogHandler.DeleteProduct)
}
