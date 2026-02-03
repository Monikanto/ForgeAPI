package repository

import (
	"context"
	"fmt"

	"github.com/Monikanto/go-rest-backend/internal/db"
	"github.com/Monikanto/go-rest-backend/internal/model"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create Order
	queryOrder := `
		INSERT INTO orders (user_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, queryOrder, order.UserID, order.Status, order.TotalAmount).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	// Create Order Items
	queryItem := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	for i, item := range order.Items {
		item.OrderID = order.ID // Link item to order
		err = tx.QueryRow(ctx, queryItem, item.OrderID, item.ProductID, item.Quantity, item.Price).
			Scan(&item.ID)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
		order.Items[i].ID = item.ID // Update ID in struct
		order.Items[i].OrderID = item.OrderID
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]model.Order, error) {
	// Simple fetch, generally ideally we join items but for list view maybe just order info
	// For "Order flow", getting details usually includes items.
	// I'll just fetch orders for now. User can query detail separately if needed, or I can JOIN.
	// Let's just fetch orders first.
	query := `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE user_id = $1`
	rows, err := db.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
