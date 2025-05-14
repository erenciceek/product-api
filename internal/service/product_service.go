package service

import (
	"context"
	"product-api/internal/dto"
	"product-api/internal/model"
	"product-api/internal/repository"
	"time"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error)
	GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error
	SearchProducts(ctx context.Context, req *dto.SearchProductRequest) ([]*dto.ProductResponse, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return mapProductToResponse(product), nil
}

func (s *productService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapProductToResponse(product), nil
}

func (s *productService) GetAllProducts(ctx context.Context) ([]*dto.ProductResponse, error) {
	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ProductResponse
	for _, product := range products {
		responses = append(responses, mapProductToResponse(product))
	}

	return responses, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	existingProduct, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existingProduct.Name = req.Name
	}
	if req.Description != "" {
		existingProduct.Description = req.Description
	}
	if req.Price > 0 {
		existingProduct.Price = req.Price
	}

	if err := s.repo.Update(ctx, id, existingProduct); err != nil {
		return nil, err
	}

	return mapProductToResponse(existingProduct), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) SearchProducts(ctx context.Context, req *dto.SearchProductRequest) ([]*dto.ProductResponse, error) {
	products, err := s.repo.Search(ctx, req.Name, req.ExactMatch, req.MinPrice, req.MaxPrice, req.SortByPrice)
	if err != nil {
		return nil, err
	}

	var responses []*dto.ProductResponse
	for _, product := range products {
		responses = append(responses, mapProductToResponse(product))
	}

	return responses, nil
}

func mapProductToResponse(product *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID.Hex(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   product.UpdatedAt.Format(time.RFC3339),
	}
}
