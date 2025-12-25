package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
