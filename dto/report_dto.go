package dto

type DashboardResponse struct {
	TodayRevenue      float64 `json:"today_revenue"`
	TodayTransactions int     `json:"today_transactions"`
	TotalProducts     int     `json:"total_products"`
	LowStockCount     int     `json:"low_stock_count"`
}

type DailyRevenueResponse struct {
	Date             string  `json:"date"`
	Revenue          float64 `json:"revenue"`
	TransactionCount int     `json:"transaction_count"`
}

type PaymentDistributionResponse struct {
	PaymentMethod string  `json:"payment_method"`
	Count         int     `json:"count"`
	Percentage    float64 `json:"percentage"`
}

type TopProductResponse struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TotalSold    int     `json:"total_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

type TopProductData struct {
	ProductID     uint
	ProductName   string
	TotalQuantity int
	TotalRevenue  float64
}

type DashboardReport struct {
	TotalRevenue       float64               `json:"total_revenue"`
	TotalTransactions  int64                 `json:"total_transactions"`
	AverageTransaction float64               `json:"average_transaction"`
	TotalTax           float64               `json:"total_tax"`
	PaymentMethodStats []PaymentMethodStat   `json:"payment_method_stats"`
	DailyRevenue       []DailyRevenueStat    `json:"daily_revenue"`
	TopSellingProducts []TopProductStat      `json:"top_selling_products"`
	LowStockProducts   []LowStockProductStat `json:"low_stock_products"`
}

type PaymentMethodStat struct {
	Method      string  `json:"method"`
	Count       int64   `json:"count"`
	TotalAmount float64 `json:"total_amount"`
	Percentage  float64 `json:"percentage"`
}

type DailyRevenueStat struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Count   int64   `json:"count"`
}

type TopProductStat struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int64   `json:"quantity"`
	Revenue     float64 `json:"revenue"`
}

type LowStockProductStat struct {
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
}

type ReportFilter struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Period    string `form:"period"` // daily, weekly, monthly
}
