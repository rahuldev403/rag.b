package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahuldev403/forgeflow/internal/config"
)

var DB *pgxpool.Pool

func Connect() {
	dbUrl := config.GetEnv("DATABASE_URL", "")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL environment variable is missing")
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database didn't respond to ping: %v", err)
	}
	DB = pool
	log.Println("Successfully connected to PostgreSQL via pgxpool")

}

func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
