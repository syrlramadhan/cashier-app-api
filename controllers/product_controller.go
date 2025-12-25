package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/services"
)

type ProductController struct {
	productService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Get list of all products with optional filters
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id query int false "Filter by category ID"
// @Param search query string false "Search by name or SKU"
// @Success 200 {object} dto.APIResponse{data=[]dto.ProductResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /products [get]
func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	categoryIDStr := ctx.Query("category_id")
	search := ctx.Query("search")

	var categoryID *uint
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}

	products, err := c.productService.GetAllProducts(categoryID, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get products",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} dto.APIResponse{data=dto.ProductResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /products/{id} [get]
func (c *ProductController) GetProductByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
			Error:   err.Error(),
		})
		return
	}

	product, err := c.productService.GetProductByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Product not found",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Product retrieved successfully",
		Data:    product,
	})
}

// CreateProduct godoc
// @Summary Create new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateProductRequest true "Create product request"
// @Success 201 {object} dto.APIResponse{data=dto.ProductResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /products [post]
func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var req dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	product, err := c.productService.CreateProduct(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to create product",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product details
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body dto.UpdateProductRequest true "Update product request"
// @Success 200 {object} dto.APIResponse{data=dto.ProductResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /products/{id} [put]
func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	product, err := c.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to update product",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /products/{id} [delete]
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
			Error:   err.Error(),
		})
		return
	}

	err = c.productService.DeleteProduct(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Failed to delete product",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}

// UpdateStock godoc
// @Summary Update product stock
// @Description Update stock quantity for a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param request body object{quantity int} true "Stock update request"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Router /products/{id}/stock [patch]
func (c *ProductController) UpdateStock(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid product ID",
			Error:   err.Error(),
		})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	err = c.productService.UpdateStock(uint(id), req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Failed to update stock",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Stock updated successfully",
	})
}

// GetProductsByCategory godoc
// @Summary Get products by category
// @Description Get all products in a specific category
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id path int true "Category ID"
// @Success 200 {object} dto.APIResponse{data=[]dto.ProductResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /products/category/{category_id} [get]
func (c *ProductController) GetProductsByCategory(ctx *gin.Context) {
	categoryID, err := strconv.ParseUint(ctx.Param("category_id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
		})
		return
	}

	products, err := c.productService.GetProductsByCategory(uint(categoryID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get products",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Products retrieved successfully",
		Data:    products,
	})
}
