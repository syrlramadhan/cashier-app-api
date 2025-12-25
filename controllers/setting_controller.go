package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/services"
)

type SettingController struct {
	settingService *services.SettingService
}

func NewSettingController(settingService *services.SettingService) *SettingController {
	return &SettingController{settingService: settingService}
}

// GetAllSettings godoc
// @Summary Get all settings
// @Description Get list of all settings
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.SettingResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /settings [get]
func (c *SettingController) GetAllSettings(ctx *gin.Context) {
	settings, err := c.settingService.GetAllSettings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get settings",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Settings retrieved successfully",
		Data:    settings,
	})
}

// GetSettingByKey godoc
// @Summary Get setting by key
// @Description Get setting value by key
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "Setting key"
// @Success 200 {object} dto.APIResponse{data=dto.SettingResponse}
// @Failure 404 {object} dto.APIResponse
// @Router /settings/{key} [get]
func (c *SettingController) GetSettingByKey(ctx *gin.Context) {
	key := ctx.Param("key")

	setting, err := c.settingService.GetSettingByKey(key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Setting not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Setting retrieved successfully",
		Data:    setting,
	})
}

// UpdateSetting godoc
// @Summary Update setting
// @Description Update or create a setting
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateSettingRequest true "Update setting request"
// @Success 200 {object} dto.APIResponse{data=dto.SettingResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /settings [put]
func (c *SettingController) UpdateSetting(ctx *gin.Context) {
	var req dto.UpdateSettingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	setting, err := c.settingService.UpdateSetting(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to update setting",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Setting updated successfully",
		Data:    setting,
	})
}

// UpdateSettings godoc
// @Summary Update multiple settings
// @Description Update multiple settings at once
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body []dto.UpdateSettingRequest true "Update settings request"
// @Success 200 {object} dto.APIResponse{data=[]dto.SettingResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /settings/batch [put]
func (c *SettingController) UpdateSettings(ctx *gin.Context) {
	var requests []dto.UpdateSettingRequest
	if err := ctx.ShouldBindJSON(&requests); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	var results []dto.SettingResponse
	for _, req := range requests {
		setting, err := c.settingService.UpdateSetting(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, dto.APIResponse{
				Success: false,
				Message: "Failed to update setting: " + req.Key,
				Error:   err.Error(),
			})
			return
		}
		results = append(results, *setting)
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Settings updated successfully",
		Data:    results,
	})
}

// GetStoreSettings godoc
// @Summary Get store settings
// @Description Get all store-related settings
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=map[string]string}
// @Failure 500 {object} dto.APIResponse
// @Router /settings/store [get]
func (c *SettingController) GetStoreSettings(ctx *gin.Context) {
	storeKeys := []string{"store_name", "store_address", "store_phone", "store_email", "tax_rate", "currency"}

	settings := make(map[string]string)
	for _, key := range storeKeys {
		setting, err := c.settingService.GetSettingByKey(key)
		if err == nil {
			settings[key] = setting.Value
		}
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Store settings retrieved successfully",
		Data:    settings,
	})
}

// GetPaymentSettings godoc
// @Summary Get payment settings
// @Description Get all payment-related settings
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=map[string]string}
// @Failure 500 {object} dto.APIResponse
// @Router /settings/payment [get]
func (c *SettingController) GetPaymentSettings(ctx *gin.Context) {
	paymentKeys := []string{"payment_cash_enabled", "payment_card_enabled", "payment_qris_enabled"}

	settings := make(map[string]string)
	for _, key := range paymentKeys {
		setting, err := c.settingService.GetSettingByKey(key)
		if err == nil {
			settings[key] = setting.Value
		}
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Payment settings retrieved successfully",
		Data:    settings,
	})
}
