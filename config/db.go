package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var DB *pgxpool.Pool

// InitDB initializes the database connection
func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://pmqc_user:pmqc_password@postgres:5432/pmqc_db?sslmode=disable"
	}

	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		panic(fmt.Sprintf("❌ Unable to connect to database: %v", err))
	}

	fmt.Println("✅ Connected to PostgreSQL")
}
