package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/services"
)

type ReportController struct {
	reportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

// GetDashboard godoc
// @Summary Get dashboard data
// @Description Get dashboard statistics and summary
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.DashboardResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /reports/dashboard [get]
func (c *ReportController) GetDashboard(ctx *gin.Context) {
	dashboard, err := c.reportService.GetDashboard()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get dashboard data",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Dashboard data retrieved successfully",
		Data:    dashboard,
	})
}

// GetDailyRevenue godoc
// @Summary Get daily revenue
// @Description Get revenue for last N days
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days (default 7)"
// @Success 200 {object} dto.APIResponse{data=[]dto.DailyRevenueResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /reports/revenue/daily [get]
func (c *ReportController) GetDailyRevenue(ctx *gin.Context) {
	daysStr := ctx.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		days = 7
	}

	revenue, err := c.reportService.GetDailyRevenue(days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get daily revenue",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Daily revenue retrieved successfully",
		Data:    revenue,
	})
}

// GetPaymentDistribution godoc
// @Summary Get payment method distribution
// @Description Get transaction count by payment method
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.PaymentDistributionResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /reports/payment-distribution [get]
func (c *ReportController) GetPaymentDistribution(ctx *gin.Context) {
	distribution, err := c.reportService.GetPaymentDistribution()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get payment distribution",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Payment distribution retrieved successfully",
		Data:    distribution,
	})
}

// GetTopProducts godoc
// @Summary Get top selling products
// @Description Get top N best selling products
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of products (default 10)"
// @Success 200 {object} dto.APIResponse{data=[]dto.TopProductResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /reports/products/top [get]
func (c *ReportController) GetTopProducts(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	products, err := c.reportService.GetTopProducts(limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get top products",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Top products retrieved successfully",
		Data:    products,
	})
}

// GetRevenueByDateRange godoc
// @Summary Get revenue by date range
// @Description Get total revenue for a specific date range
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /reports/revenue/range [get]
func (c *ReportController) GetRevenueByDateRange(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
			Error:   err.Error(),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
			Error:   err.Error(),
		})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	revenue, err := c.reportService.GetRevenueByDateRange(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get revenue",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Revenue retrieved successfully",
		Data: map[string]interface{}{
			"start_date":    startDateStr,
			"end_date":      endDateStr,
			"total_revenue": revenue,
		},
	})
}

// GetMonthlySummary godoc
// @Summary Get monthly summary
// @Description Get monthly summary for current month
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=map[string]interface{}}
// @Failure 500 {object} dto.APIResponse
// @Router /reports/summary/monthly [get]
func (c *ReportController) GetMonthlySummary(ctx *gin.Context) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	revenue, err := c.reportService.GetRevenueByDateRange(startOfMonth, endOfMonth)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get monthly summary",
			Error:   err.Error(),
		})
		return
	}

	// Get transaction count for the month
	dashboard, _ := c.reportService.GetDashboard()

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Monthly summary retrieved successfully",
		Data: map[string]interface{}{
			"month":              now.Month().String(),
			"year":               now.Year(),
			"total_revenue":      revenue,
			"total_transactions": dashboard.TodayTransactions, // This would need adjustment for monthly
		},
	})
}

// ExportTransactions godoc
// @Summary Export transactions
// @Description Export transactions data for a date range
// @Tags reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} dto.APIResponse{data=[]dto.TransactionResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /reports/export/transactions [get]
func (c *ReportController) ExportTransactions(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "start_date and end_date are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid start_date format. Use YYYY-MM-DD",
			Error:   err.Error(),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid end_date format. Use YYYY-MM-DD",
			Error:   err.Error(),
		})
		return
	}

	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	transactions, err := c.reportService.ExportTransactions(startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to export transactions",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Transactions exported successfully",
		Data:    transactions,
	})
}
