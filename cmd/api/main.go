package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Monikanto/go-rest-backend/internal/config"
	"github.com/Monikanto/go-rest-backend/internal/db"
	"github.com/Monikanto/go-rest-backend/internal/handler"
	"github.com/Monikanto/go-rest-backend/internal/repository"
	"github.com/Monikanto/go-rest-backend/internal/routes"
	"github.com/Monikanto/go-rest-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to Database
	db.Connect(cfg)
	defer db.DB.Close()

	// Initialize Repositories
	userRepo := repository.NewUserRepository()
	productRepo := repository.NewProductRepository()
	orderRepo := repository.NewOrderRepository()

	// Initialize Services
	userCmd := service.NewUserService(userRepo)
	productSvc := service.NewProductService(productRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo)

	// Initialize Handlers
	userHandler := handler.NewUserHandler(userCmd)
	productHandler := handler.NewProductHandler(productSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// Setup Router
	r := gin.Default()

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup Routes
	routes.SetupRoutes(r, cfg, &routes.Handlers{
		User:    userHandler,
		Product: productHandler,
		Order:   orderHandler,
	})

	srv := &http.Server{
		Addr: ":" + cfg.DBPort, // Using DBPort temporarily or hardcoded?
		// Actually config usually has ServerPort. Let's use 8000 default or env.
		// Re-checking config.go, it only has DB ports. I should probably default to 8000.
		Handler: r,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	srv.Addr = ":" + port

	// Graceful Shutdown
	go func() {
		log.Printf("Server running on %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
