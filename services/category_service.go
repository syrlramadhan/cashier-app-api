package services

import (
	"errors"

	"github.com/syrlramadhan/cashier-app/dto"
	"github.com/syrlramadhan/cashier-app/models"
	"github.com/syrlramadhan/cashier-app/repositories"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetAllCategories() ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []dto.CategoryResponse
	for _, category := range categories {
		response = append(response, dto.CategoryResponse{
			ID:   category.ID,
			Name: category.Name,
		})
	}

	return response, nil
}

func (s *CategoryService) GetCategoryByID(id uint) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// Check if category name already exists
	existingCategory, _ := s.categoryRepo.FindByName(req.Name)
	if existingCategory != nil {
		return nil, errors.New("category name already exists")
	}

	category := &models.Category{
		Name: req.Name,
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, errors.New("failed to create category")
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (s *CategoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Check if new name already exists (excluding current category)
	if req.Name != category.Name {
		existingCategory, _ := s.categoryRepo.FindByName(req.Name)
		if existingCategory != nil {
			return nil, errors.New("category name already exists")
		}
	}

	category.Name = req.Name
	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, errors.New("failed to update category")
	}

	response := &dto.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}

	return response, nil
}

func (s *CategoryService) DeleteCategory(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.Delete(id)
}
