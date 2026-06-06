package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	commonResponse "github.com/aridevk/dark-kitchen/packages/go/common/response"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/dto"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/services"
)

type CatalogHandler struct {
	catalogService *services.CatalogService
}

func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

func (h *CatalogHandler) GetCategories(c *gin.Context) {
	categories, err := h.catalogService.GetCategories()
	if err != nil {
		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, categories)
}

func (h *CatalogHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	category, err := h.catalogService.CreateCategory(req)
	if err != nil {
		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.Created(c, category)
}

func (h *CatalogHandler) GetProducts(c *gin.Context) {
	onlyAvailable := c.DefaultQuery("available", "true") == "true"

	products, err := h.catalogService.GetProducts(onlyAvailable)
	if err != nil {
		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, products)
}

func (h *CatalogHandler) GetProductByID(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		return
	}

	product, err := h.catalogService.GetProductByID(id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			commonResponse.Error(c, 404, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, product)
}

func (h *CatalogHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	product, err := h.catalogService.CreateProduct(req)
	if err != nil {
		if errors.Is(err, services.ErrCategoryNotFound) {
			commonResponse.Error(c, 404, "CATEGORY_NOT_FOUND", "Category not found")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.Created(c, product)
}

func (h *CatalogHandler) UpdateProduct(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		return
	}

	var req dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	product, err := h.catalogService.UpdateProduct(id, req)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			commonResponse.Error(c, 404, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}

		if errors.Is(err, services.ErrCategoryNotFound) {
			commonResponse.Error(c, 404, "CATEGORY_NOT_FOUND", "Category not found")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, product)
}

func (h *CatalogHandler) DeleteProduct(c *gin.Context) {
	id, ok := getIDParam(c)
	if !ok {
		return
	}

	err := h.catalogService.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			commonResponse.Error(c, 404, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, gin.H{
		"message": "Product deleted successfully",
	})
}

func getIDParam(c *gin.Context) (uint, bool) {
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		commonResponse.BadRequest(c, "INVALID_ID", "Invalid id")
		return 0, false
	}

	return uint(id), true
}
