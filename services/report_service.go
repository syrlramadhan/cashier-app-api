package services

import (
	"time"

	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/repositories"
)

type ReportService struct {
	transactionRepo     repositories.TransactionRepository
	transactionItemRepo repositories.TransactionItemRepository
	productRepo         repositories.ProductRepository
	categoryRepo        repositories.CategoryRepository
}

func NewReportService(
	transactionRepo repositories.TransactionRepository,
	transactionItemRepo repositories.TransactionItemRepository,
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
) *ReportService {
	return &ReportService{
		transactionRepo:     transactionRepo,
		transactionItemRepo: transactionItemRepo,
		productRepo:         productRepo,
		categoryRepo:        categoryRepo,
	}
}

func (s *ReportService) GetDashboard() (*dto.DashboardResponse, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	// Get today's revenue
	todayRevenue, _ := s.transactionRepo.GetTotalRevenue(startOfDay, endOfDay)

	// Get today's transactions count
	transactions, _ := s.transactionRepo.FindByDateRange(startOfDay, endOfDay)
	todayTransactions := len(transactions)

	// Get total products
	products, _ := s.productRepo.FindAll()
	totalProducts := len(products)

	// Get low stock count (stock < 10)
	lowStockCount := 0
	for _, p := range products {
		if p.Stock < 10 {
			lowStockCount++
		}
	}

	return &dto.DashboardResponse{
		TodayRevenue:      todayRevenue,
		TodayTransactions: todayTransactions,
		TotalProducts:     totalProducts,
		LowStockCount:     lowStockCount,
	}, nil
}

func (s *ReportService) GetDailyRevenue(days int) ([]dto.DailyRevenueResponse, error) {
	var result []dto.DailyRevenueResponse

	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

		revenue, _ := s.transactionRepo.GetTotalRevenue(startOfDay, endOfDay)
		transactions, _ := s.transactionRepo.FindByDateRange(startOfDay, endOfDay)

		result = append(result, dto.DailyRevenueResponse{
			Date:             date.Format("2006-01-02"),
			Revenue:          revenue,
			TransactionCount: len(transactions),
		})
	}

	return result, nil
}

func (s *ReportService) GetPaymentDistribution() ([]dto.PaymentDistributionResponse, error) {
	paymentMethods := []string{"cash", "card", "qris"}
	var result []dto.PaymentDistributionResponse

	totalCount := 0
	counts := make(map[string]int)

	for _, method := range paymentMethods {
		count, _ := s.transactionRepo.CountByPaymentMethod(method)
		counts[method] = int(count)
		totalCount += int(count)
	}

	for _, method := range paymentMethods {
		percentage := 0.0
		if totalCount > 0 {
			percentage = float64(counts[method]) / float64(totalCount) * 100
		}
		result = append(result, dto.PaymentDistributionResponse{
			PaymentMethod: method,
			Count:         counts[method],
			Percentage:    percentage,
		})
	}

	return result, nil
}

func (s *ReportService) GetTopProducts(limit int) ([]dto.TopProductResponse, error) {
	topProducts, err := s.transactionItemRepo.GetTopProducts(limit)
	if err != nil {
		return nil, err
	}

	var result []dto.TopProductResponse
	for _, tp := range topProducts {
		result = append(result, dto.TopProductResponse{
			ProductID:    tp.ProductID,
			ProductName:  tp.ProductName,
			TotalSold:    tp.TotalQuantity,
			TotalRevenue: tp.TotalRevenue,
		})
	}

	return result, nil
}

func (s *ReportService) GetRevenueByDateRange(startDate, endDate time.Time) (float64, error) {
	return s.transactionRepo.GetTotalRevenue(startDate, endDate)
}

func (s *ReportService) ExportTransactions(startDate, endDate time.Time) ([]dto.TransactionResponse, error) {
	transactions, err := s.transactionRepo.FindByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var result []dto.TransactionResponse
	for _, t := range transactions {
		var items []dto.TransactionItemResponse
		for _, item := range t.Items {
			items = append(items, dto.TransactionItemResponse{
				ID:          item.ID,
				ProductID:   item.ProductID,
				ProductName: item.ProductName,
				Price:       item.Price,
				Quantity:    item.Quantity,
				Subtotal:    item.Subtotal,
			})
		}

		cashierName := ""
		if t.User.ID > 0 {
			cashierName = t.User.Name
		}

		result = append(result, dto.TransactionResponse{
			ID:              t.ID,
			TransactionCode: t.TransactionCode,
			CashierName:     cashierName,
			Subtotal:        t.Subtotal,
			Tax:             t.Tax,
			Total:           t.Total,
			PaymentMethod:   t.PaymentMethod,
			Status:          t.Status,
			Items:           items,
			CreatedAt:       t.CreatedAt,
		})
	}

	return result, nil
}
