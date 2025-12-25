package repositories

import (
	"time"

	"github.com/syrlramadhan/cashier-app/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindAll() ([]models.Transaction, error)
	FindAllWithDetails() ([]models.Transaction, error)
	FindByID(id uint) (*models.Transaction, error)
	FindByIDWithDetails(id uint) (*models.Transaction, error)
	FindByCode(code string) (*models.Transaction, error)
	FindByUserID(userID uint) ([]models.Transaction, error)
	FindByDateRange(startDate, endDate time.Time) ([]models.Transaction, error)
	FindByPaymentMethod(method string) ([]models.Transaction, error)
	FindWithFilters(startDate, endDate *time.Time, paymentMethod string, limit, offset int) ([]models.Transaction, int64, error)
	Create(transaction *models.Transaction) error
	Update(transaction *models.Transaction) error
	Delete(id uint) error
	Count() (int64, error)
	CountByDateRange(startDate, endDate time.Time) (int64, error)
	CountByPaymentMethod(method string) (int64, error)
	GetTotalRevenue(startDate, endDate time.Time) (float64, error)
	GetTotalRevenueByDateRange(startDate, endDate time.Time) (float64, error)
	GetRevenueByPaymentMethod() ([]map[string]interface{}, error)
	GetDailyRevenue(days int) ([]map[string]interface{}, error)
	GenerateTransactionCode() (string, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindAllWithDetails() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Preload("Items.Product").Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByIDWithDetails(id uint) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Items").Preload("Items.Product").First(&transaction, id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByCode(code string) (*models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("transaction_code = ?", code).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) FindByUserID(userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("created_at BETWEEN ? AND ?", startDate, endDate).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByPaymentMethod(method string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Items").Where("payment_method = ?", method).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindWithFilters(startDate, endDate *time.Time, paymentMethod string, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{})

	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("User").Preload("Items").Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}

func (r *transactionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Count(&count).Error
	return count, err
}

func (r *transactionRepository) CountByDateRange(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Count(&count).Error
	return count, err
}

func (r *transactionRepository) CountByPaymentMethod(method string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).Where("payment_method = ? AND status = ?", method, "completed").Count(&count).Error
	return count, err
}

func (r *transactionRepository) GetTotalRevenue(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Transaction{}).Where("created_at BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").Select("COALESCE(SUM(total), 0)").Scan(&total).Error
	return total, err
}

func (r *transactionRepository) GetTotalRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Transaction{}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Select("COALESCE(SUM(total), 0)").Scan(&total).Error
	return total, err
}

func (r *transactionRepository) GetRevenueByPaymentMethod() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.Transaction{}).
		Select("payment_method, COUNT(*) as count, COALESCE(SUM(total), 0) as total_amount").
		Group("payment_method").
		Find(&results).Error
	return results, err
}

func (r *transactionRepository) GetDailyRevenue(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.Transaction{}).
		Select("DATE(created_at) as date, COUNT(*) as count, COALESCE(SUM(total), 0) as revenue").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -days)).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&results).Error
	return results, err
}

func (r *transactionRepository) GenerateTransactionCode() (string, error) {
	var count int64
	today := time.Now().Format("20060102")
	r.db.Model(&models.Transaction{}).Where("DATE(created_at) = CURRENT_DATE").Count(&count)
	code := "TRX-" + today + "-" + padNumber(int(count+1), 4)
	return code, nil
}

func padNumber(num int, length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result = "0" + result
	}
	numStr := ""
	for num > 0 {
		numStr = string(rune('0'+num%10)) + numStr
		num /= 10
	}
	if len(numStr) >= length {
		return numStr
	}
	return result[:length-len(numStr)] + numStr
}
