package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Monikanto/go-rest-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse Db config:", err)
	}

	// connection pool tuning(production basics)
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	// ping test
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatal("Database ping failed", err)
	}

	log.Println("Connected to PostgresSQL")
}