package services

import (
	"errors"

	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/models"
	"github.com/syrlramadhan/cashier-app/repositories"
)

type ProductService struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) GetAllProducts(categoryID *uint, search string) ([]dto.ProductResponse, error) {
	var products []models.Product
	var err error

	if categoryID != nil {
		products, err = s.productRepo.FindByCategoryID(*categoryID)
	} else {
		products, err = s.productRepo.FindAll()
	}

	if err != nil {
		return nil, err
	}

	var response []dto.ProductResponse
	for _, product := range products {
		// Apply search filter
		if search != "" {
			// Simple search - can be improved with database-level search
			continue
		}

		response = append(response, dto.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryID: product.CategoryID,
			Image:      product.Image,
		})
	}

	return response, nil
}

func (s *ProductService) GetProductByID(id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	response := &dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	}

	return response, nil
}

func (s *ProductService) GetProductsByCategory(categoryID uint) ([]dto.ProductResponse, error) {
	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	products, err := s.productRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	var response []dto.ProductResponse
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CategoryID: product.CategoryID,
			Image:      product.Image,
		})
	}

	return response, nil
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	_, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product := &models.Product{
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		Image:      req.Image,
	}

	err = s.productRepo.Create(product)
	if err != nil {
		return nil, errors.New("failed to create product")
	}

	response := &dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	}

	return response, nil
}

func (s *ProductService) UpdateProduct(id uint, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	_, err = s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Stock = req.Stock
	product.CategoryID = req.CategoryID
	product.Image = req.Image

	err = s.productRepo.Update(product)
	if err != nil {
		return nil, errors.New("failed to update product")
	}

	response := &dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	}

	return response, nil
}

func (s *ProductService) UpdateStock(id uint, quantity int) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	return s.productRepo.UpdateStock(id, newStock)
}

func (s *ProductService) DeleteProduct(id uint) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(id)
}
