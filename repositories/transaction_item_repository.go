package repositories

import (
	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/models"
	"gorm.io/gorm"
)

type TransactionItemRepository interface {
	FindByTransactionID(transactionID uint) ([]models.TransactionItem, error)
	Create(item *models.TransactionItem) error
	CreateBatch(items []models.TransactionItem) error
	Delete(id uint) error
	DeleteByTransactionID(transactionID uint) error
	GetTopProducts(limit int) ([]dto.TopProductData, error)
}

type transactionItemRepository struct {
	db *gorm.DB
}

func NewTransactionItemRepository(db *gorm.DB) TransactionItemRepository {
	return &transactionItemRepository{db: db}
}

func (r *transactionItemRepository) FindByTransactionID(transactionID uint) ([]models.TransactionItem, error) {
	var items []models.TransactionItem
	err := r.db.Preload("Product").Where("transaction_id = ?", transactionID).Find(&items).Error
	return items, err
}

func (r *transactionItemRepository) Create(item *models.TransactionItem) error {
	return r.db.Create(item).Error
}

func (r *transactionItemRepository) CreateBatch(items []models.TransactionItem) error {
	return r.db.Create(&items).Error
}

func (r *transactionItemRepository) Delete(id uint) error {
	return r.db.Delete(&models.TransactionItem{}, id).Error
}

func (r *transactionItemRepository) DeleteByTransactionID(transactionID uint) error {
	return r.db.Where("transaction_id = ?", transactionID).Delete(&models.TransactionItem{}).Error
}

func (r *transactionItemRepository) GetTopProducts(limit int) ([]dto.TopProductData, error) {
	var results []dto.TopProductData
	err := r.db.Model(&models.TransactionItem{}).
		Select("product_id, product_name, SUM(quantity) as total_quantity, SUM(subtotal) as total_revenue").
		Group("product_id, product_name").
		Order("total_quantity DESC").
		Limit(limit).
		Scan(&results).Error
	return results, err
}
