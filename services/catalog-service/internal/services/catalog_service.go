package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/dto"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/models"
	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/repositories"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
)

type CatalogService struct {
	categoryRepo *repositories.CategoryRepository
	productRepo  *repositories.ProductRepository
}

func NewCatalogService(
	categoryRepo *repositories.CategoryRepository,
	productRepo *repositories.ProductRepository,
) *CatalogService {
	return &CatalogService{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

func (s *CatalogService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *CatalogService) CreateCategory(req dto.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CatalogService) GetProducts(onlyAvailable bool) ([]models.Product, error) {
	return s.productRepo.FindAll(onlyAvailable)
}

func (s *CatalogService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (s *CatalogService) CreateProduct(req dto.CreateProductRequest) (*models.Product, error) {
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	isAvailable := true
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	product := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsAvailable: isAvailable,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return s.productRepo.FindByID(product.ID)
}

func (s *CatalogService) UpdateProduct(id uint, req dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if req.CategoryID != nil {
		if _, err := s.categoryRepo.FindByID(*req.CategoryID); err != nil {
			return nil, ErrCategoryNotFound
		}

		product.CategoryID = *req.CategoryID
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	product.Description = req.Description

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	if req.IsAvailable != nil {
		product.IsAvailable = *req.IsAvailable
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return s.productRepo.FindByID(product.ID)
}

func (s *CatalogService) DeleteProduct(id uint) error {
	if _, err := s.productRepo.FindByID(id); err != nil {
		return ErrProductNotFound
	}

	return s.productRepo.Delete(id)
}
