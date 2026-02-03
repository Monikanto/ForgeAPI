package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Monikanto/go-rest-backend/internal/model"
	"github.com/Monikanto/go-rest-backend/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []OrderItemRequest) (*model.Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	var orderItems []model.OrderItem
	var totalAmount float64

	for _, item := range items {
		// Fetch product to get price
		product, err := s.productRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		// Calculate item total
		price := product.Price
		itemTotal := price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     price,
		})
	}

	order := &model.Order{
		UserID:      userID,
		Status:      "pending",
		TotalAmount: totalAmount,
		Items:       orderItems,
	}

	if err := s.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string) ([]model.Order, error) {
	return s.orderRepo.GetOrdersByUserID(ctx, userID)
}
