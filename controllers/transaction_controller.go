package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/services"
)

type TransactionController struct {
	transactionService *services.TransactionService
}

func NewTransactionController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

// GetAllTransactions godoc
// @Summary Get all transactions
// @Description Get list of all transactions with optional date filters
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param payment_method query string false "Filter by payment method"
// @Success 200 {object} dto.APIResponse{data=[]dto.TransactionResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /transactions [get]
func (c *TransactionController) GetAllTransactions(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	paymentMethod := ctx.Query("payment_method")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &t
		}
	}

	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			// Set to end of day
			t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			endDate = &t
		}
	}

	transactions, err := c.transactionService.GetAllTransactions(startDate, endDate, paymentMethod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get transactions",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    transactions,
	})
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Get transaction details by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} dto.APIResponse{data=dto.TransactionResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /transactions/{id} [get]
func (c *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid transaction ID",
			Error:   err.Error(),
		})
		return
	}

	transaction, err := c.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Transaction not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}

// GetTransactionByCode godoc
// @Summary Get transaction by code
// @Description Get transaction details by transaction code
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code path string true "Transaction Code"
// @Success 200 {object} dto.APIResponse{data=dto.TransactionResponse}
// @Failure 404 {object} dto.APIResponse
// @Router /transactions/code/{code} [get]
func (c *TransactionController) GetTransactionByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	transaction, err := c.transactionService.GetTransactionByCode(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Transaction not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transaction retrieved successfully",
		Data:    transaction,
	})
}

// CreateTransaction godoc
// @Summary Create new transaction
// @Description Create a new transaction (checkout)
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Create transaction request"
// @Success 201 {object} dto.APIResponse{data=dto.TransactionResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /transactions [post]
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	req.UserID = userID.(uint)

	transaction, err := c.transactionService.CreateTransaction(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to create transaction",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Transaction created successfully",
		Data:    transaction,
	})
}

// GetTodayTransactions godoc
// @Summary Get today's transactions
// @Description Get all transactions for today
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.TransactionResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /transactions/today [get]
func (c *TransactionController) GetTodayTransactions(ctx *gin.Context) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	transactions, err := c.transactionService.GetAllTransactions(&startOfDay, &endOfDay, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get transactions",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Today's transactions retrieved successfully",
		Data:    transactions,
	})
}

// GetTransactionsByUser godoc
// @Summary Get transactions by user
// @Description Get all transactions created by a specific user
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} dto.APIResponse{data=[]dto.TransactionResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /transactions/user/{user_id} [get]
func (c *TransactionController) GetTransactionsByUser(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("user_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
		return
	}

	transactions, err := c.transactionService.GetTransactionsByUser(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get transactions",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transactions retrieved successfully",
		Data:    transactions,
	})
}

// CancelTransaction godoc
// @Summary Cancel transaction
// @Description Cancel/void a transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /transactions/{id}/cancel [post]
func (c *TransactionController) CancelTransaction(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid transaction ID",
			Error:   err.Error(),
		})
		return
	}

	err = c.transactionService.CancelTransaction(uint(id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to cancel transaction",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transaction cancelled successfully",
	})
}
