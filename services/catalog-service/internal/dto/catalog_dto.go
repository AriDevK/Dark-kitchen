package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=120"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=120"`
	Description string `json:"description"`
}

type CreateProductRequest struct {
	CategoryID  uint    `json:"categoryId" binding:"required"`
	Name        string  `json:"name" binding:"required,min=2,max=160"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	ImageURL    string  `json:"imageUrl"`
	IsAvailable *bool   `json:"isAvailable"`
}

type UpdateProductRequest struct {
	CategoryID  *uint    `json:"categoryId"`
	Name        string   `json:"name" binding:"omitempty,min=2,max=160"`
	Description string   `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	ImageURL    string   `json:"imageUrl"`
	IsAvailable *bool    `json:"isAvailable"`
}
