package config

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (config *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)
}

func ConnectDB(config *DBConfig) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), config.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create tables
	CreateTables(pool)

	return pool
}

func ConnectDBFromURL(dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create tables
	CreateTables(pool)

	return pool
}

func CreateTables(pool *pgxpool.Pool) {
	createFormTableSQL := `
		CREATE TABLE IF NOT EXISTS forms (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := pool.Exec(context.Background(), createFormTableSQL)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}
