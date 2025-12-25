package repositories

import (
	"github.com/syrlramadhan/cashier-app/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindAllWithCategory() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	FindByIDWithCategory(id uint) (*models.Product, error)
	FindByCategoryID(categoryID uint) ([]models.Product, error)
	FindLowStock(threshold int) ([]models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	UpdateStock(id uint, stock int) error
	Delete(id uint) error
	Count() (int64, error)
	Search(keyword string) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) FindAllWithCategory() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByIDWithCategory(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *productRepository) FindLowStock(threshold int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("stock <= ?", threshold).Order("stock ASC").Find(&products).Error
	return products, err
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateStock(id uint, stock int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", stock).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *productRepository) Search(keyword string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Where("name ILIKE ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}
