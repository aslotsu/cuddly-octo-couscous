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

	createBlogsTableSQL := `
		CREATE TABLE IF NOT EXISTS blogs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content JSONB NOT NULL,
			author VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	createBlogImagesTableSQL := `
		CREATE TABLE IF NOT EXISTS blog_images (
			id SERIAL PRIMARY KEY,
			blog_id INTEGER REFERENCES blogs(id) ON DELETE CASCADE,
			image_key VARCHAR(500) NOT NULL,
			image_url VARCHAR(1000) NOT NULL,
			alt_text VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	createApiKeysTableSQL := `
		CREATE TABLE IF NOT EXISTS api_keys (
			id SERIAL PRIMARY KEY,
			key_hash VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := pool.Exec(context.Background(), createFormTableSQL)
	if err != nil {
		log.Fatalf("Failed to create forms table: %v", err)
	}

	_, err = pool.Exec(context.Background(), createBlogsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create blogs table: %v", err)
	}

	_, err = pool.Exec(context.Background(), createBlogImagesTableSQL)
	if err != nil {
		log.Fatalf("Failed to create blog_images table: %v", err)
	}

	_, err = pool.Exec(context.Background(), createApiKeysTableSQL)
	if err != nil {
		log.Fatalf("Failed to create api_keys table: %v", err)
	}
}
