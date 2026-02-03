package service

import (
	"context"

	"github.com/Monikanto/go-rest-backend/internal/model"
	"github.com/Monikanto/go-rest-backend/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description string, price float64, stock int) (*model.Product, error) {
	product := &model.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
	if err := s.repo.CreateProduct(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *ProductService) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, name, description string, price float64, stock int) (*model.Product, error) {
	product := &model.Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}
	if err := s.repo.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.DeleteProduct(ctx, id)
}
