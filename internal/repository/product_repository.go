package repository

import (
	"context"
	"fmt"

	"github.com/Monikanto/go-rest-backend/internal/db"
	"github.com/Monikanto/go-rest-backend/internal/model"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *model.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := db.DB.QueryRow(ctx, query, p.Name, p.Description, p.Price, p.Stock).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*model.Product, error) {
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	var p model.Product
	err := db.DB.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &p, nil
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products`
	rows, err := db.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, p *model.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := db.DB.QueryRow(ctx, query, p.Name, p.Description, p.Price, p.Stock, p.ID).Scan(&p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
