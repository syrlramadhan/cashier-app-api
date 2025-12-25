package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/models"
	"github.com/syrlramadhan/cashier-app/repositories"
)

type TransactionService struct {
	transactionRepo     repositories.TransactionRepository
	transactionItemRepo repositories.TransactionItemRepository
	productRepo         repositories.ProductRepository
}

func NewTransactionService(
	transactionRepo repositories.TransactionRepository,
	transactionItemRepo repositories.TransactionItemRepository,
	productRepo repositories.ProductRepository,
) *TransactionService {
	return &TransactionService{
		transactionRepo:     transactionRepo,
		transactionItemRepo: transactionItemRepo,
		productRepo:         productRepo,
	}
}

func (s *TransactionService) GetAllTransactions(startDate, endDate *time.Time, paymentMethod string) ([]dto.TransactionResponse, error) {
	transactions, err := s.transactionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.TransactionResponse
	for _, transaction := range transactions {
		// Apply date filter
		if startDate != nil && transaction.CreatedAt.Before(*startDate) {
			continue
		}
		if endDate != nil && transaction.CreatedAt.After(*endDate) {
			continue
		}
		// Apply payment method filter
		if paymentMethod != "" && transaction.PaymentMethod != paymentMethod {
			continue
		}

		response = append(response, *s.mapTransactionToResponse(&transaction))
	}

	return response, nil
}

func (s *TransactionService) GetTransactionByID(id uint) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return s.mapTransactionToResponse(transaction), nil
}

func (s *TransactionService) GetTransactionByCode(code string) (*dto.TransactionResponse, error) {
	transaction, err := s.transactionRepo.FindByCode(code)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	return s.mapTransactionToResponse(transaction), nil
}

func (s *TransactionService) GetTransactionsByUser(userID uint) ([]dto.TransactionResponse, error) {
	transactions, err := s.transactionRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.TransactionResponse
	for _, transaction := range transactions {
		response = append(response, *s.mapTransactionToResponse(&transaction))
	}

	return response, nil
}

func (s *TransactionService) CreateTransaction(req *dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("transaction must have at least one item")
	}

	taxRate := 0.11 // 11% tax rate

	// Validate products and calculate totals
	var subtotal float64
	var items []models.TransactionItem

	for _, itemReq := range req.Items {
		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %d", itemReq.ProductID)
		}

		if product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s", product.Name)
		}

		itemSubtotal := product.Price * float64(itemReq.Quantity)
		subtotal += itemSubtotal

		items = append(items, models.TransactionItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Price:       product.Price,
			Quantity:    itemReq.Quantity,
			Subtotal:    itemSubtotal,
		})
	}

	tax := subtotal * taxRate
	total := subtotal + tax

	// Generate transaction code
	transactionCode := fmt.Sprintf("TRX%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)

	// Create transaction
	transaction := &models.Transaction{
		TransactionCode: transactionCode,
		UserID:          req.UserID,
		Subtotal:        subtotal,
		Tax:             tax,
		Total:           total,
		PaymentMethod:   req.PaymentMethod,
		Status:          "completed",
		Items:           items,
	}

	err := s.transactionRepo.Create(transaction)
	if err != nil {
		return nil, errors.New("failed to create transaction")
	}

	// Update product stock
	for _, itemReq := range req.Items {
		product, _ := s.productRepo.FindByID(itemReq.ProductID)
		newStock := product.Stock - itemReq.Quantity
		s.productRepo.UpdateStock(itemReq.ProductID, newStock)
	}

	return s.mapTransactionToResponse(transaction), nil
}

func (s *TransactionService) CancelTransaction(id uint) error {
	transaction, err := s.transactionRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	if transaction.Status == "cancelled" {
		return errors.New("transaction is already cancelled")
	}

	// Restore stock
	for _, item := range transaction.Items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err == nil {
			newStock := product.Stock + item.Quantity
			s.productRepo.UpdateStock(item.ProductID, newStock)
		}
	}

	// Update transaction status
	transaction.Status = "cancelled"
	return s.transactionRepo.Update(transaction)
}

func (s *TransactionService) mapTransactionToResponse(transaction *models.Transaction) *dto.TransactionResponse {
	var itemResponses []dto.TransactionItemResponse
	for _, item := range transaction.Items {
		itemResponses = append(itemResponses, dto.TransactionItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		})
	}

	cashierName := ""
	if transaction.User.ID > 0 {
		cashierName = transaction.User.Name
	}

	return &dto.TransactionResponse{
		ID:              transaction.ID,
		TransactionCode: transaction.TransactionCode,
		CashierName:     cashierName,
		Subtotal:        transaction.Subtotal,
		Tax:             transaction.Tax,
		Total:           transaction.Total,
		PaymentMethod:   transaction.PaymentMethod,
		Status:          transaction.Status,
		Items:           itemResponses,
		CreatedAt:       transaction.CreatedAt,
	}
}
