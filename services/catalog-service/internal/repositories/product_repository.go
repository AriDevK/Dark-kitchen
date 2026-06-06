package repositories

import (
	"gorm.io/gorm"

	"github.com/aridevk/dark-kitchen/services/catalog-service/internal/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(onlyAvailable bool) ([]models.Product, error) {
	var products []models.Product

	query := r.db.Preload("Category").Order("name ASC")

	if onlyAvailable {
		query = query.Where("is_available = ?", true)
	}

	err := query.Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product

	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}
