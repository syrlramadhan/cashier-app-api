package dto

type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value"`
}

type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StoreSettings struct {
	StoreName     string `json:"store_name"`
	StoreAddress  string `json:"store_address"`
	StorePhone    string `json:"store_phone"`
	TaxRate       string `json:"tax_rate"`
	Currency      string `json:"currency"`
	ReceiptFooter string `json:"receipt_footer"`
}
