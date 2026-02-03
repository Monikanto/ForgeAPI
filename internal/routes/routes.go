package routes

import (
	"github.com/Monikanto/go-rest-backend/internal/config"
	"github.com/Monikanto/go-rest-backend/internal/handler"
	"github.com/Monikanto/go-rest-backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User    *handler.UserHandler
	Product *handler.ProductHandler
	Order   *handler.OrderHandler
}

func SetupRoutes(r *gin.Engine, cfg *config.Config, h *Handlers) {
	api := r.Group("/api")

	// Auth Routes
	api.POST("/register", h.User.Register)
	api.POST("/login", h.User.Login)

	// Public Product Routes
	api.GET("/products", h.Product.GetAllProducts)
	api.GET("/products/:id", h.Product.GetProductByID)

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Products (Admin/Protected)
		protected.POST("/products", h.Product.CreateProduct)
		protected.PUT("/products/:id", h.Product.UpdateProduct)
		protected.DELETE("/products/:id", h.Product.DeleteProduct)

		// Orders
		protected.POST("/orders", h.Order.CreateOrder)
		protected.GET("/orders", h.Order.GetOrders)
	}
}
