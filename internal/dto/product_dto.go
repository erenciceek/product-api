package dto

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"`
}

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type SearchProductRequest struct {
	Name        string  `query:"name"`
	ExactMatch  bool    `query:"exact_match"`
	MinPrice    float64 `query:"min_price"`
	MaxPrice    float64 `query:"max_price"`
	SortByPrice string  `query:"sort_by_price" validate:"omitempty,oneof=asc desc"`
}
