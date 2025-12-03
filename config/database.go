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

	createEventsTableSQL := `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			event_type VARCHAR(100) NOT NULL,
			status VARCHAR(50) DEFAULT 'draft',
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			venue_name VARCHAR(255),
			venue_address TEXT,
			is_virtual BOOLEAN DEFAULT false,
			virtual_link VARCHAR(500),
			timezone VARCHAR(100),
			capacity INTEGER DEFAULT 0,
			expected_guests INTEGER DEFAULT 0,
			registered_count INTEGER DEFAULT 0,
			actual_guests INTEGER,
			waitlist_enabled BOOLEAN DEFAULT false,
			allow_walkins BOOLEAN DEFAULT true,
			ticket_price DECIMAL(10, 2) DEFAULT 0.00,
			early_bird_price DECIMAL(10, 2),
			organization_budget DECIMAL(10, 2) DEFAULT 0.00,
			expenses DECIMAL(10, 2) DEFAULT 0.00,
			revenue DECIMAL(10, 2) DEFAULT 0.00,
			registration_open_date TIMESTAMP,
			registration_close_date TIMESTAMP,
			registration_form_url VARCHAR(500),
			requires_approval BOOLEAN DEFAULT false,
			featured_image VARCHAR(500),
			gallery_images JSONB DEFAULT '[]',
			video_url VARCHAR(500),
			livestream_url VARCHAR(500),
			organizer_name VARCHAR(255) NOT NULL,
			organizer_email VARCHAR(255) NOT NULL,
			organizer_phone VARCHAR(50),
			speakers JSONB DEFAULT '[]',
			sponsors JSONB DEFAULT '[]',
			tags JSONB DEFAULT '[]',
			is_featured BOOLEAN DEFAULT false,
			is_public BOOLEAN DEFAULT true,
			created_by VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	createBooksTableSQL := `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			subtitle VARCHAR(255),
			author VARCHAR(255) NOT NULL,
			isbn VARCHAR(50),
			description TEXT NOT NULL,
			publisher VARCHAR(255),
			publication_date TIMESTAMP,
			pages INTEGER DEFAULT 0,
			language VARCHAR(50) DEFAULT 'English',
			category VARCHAR(100) NOT NULL,
			price DECIMAL(10, 2) NOT NULL,
			sale_price DECIMAL(10, 2),
			stock_quantity INTEGER DEFAULT 0,
			status VARCHAR(50) DEFAULT 'available',
			cover_image VARCHAR(500),
			gallery_images JSONB DEFAULT '[]',
			preview_url VARCHAR(500),
			purchase_links JSONB DEFAULT '{}',
			tags JSONB DEFAULT '[]',
			is_featured BOOLEAN DEFAULT false,
			is_published BOOLEAN DEFAULT false,
			total_sales INTEGER DEFAULT 0,
			average_rating DECIMAL(3, 2) DEFAULT 0.00,
			review_count INTEGER DEFAULT 0,
			created_by VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

	_, err = pool.Exec(context.Background(), createEventsTableSQL)
	if err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}

	_, err = pool.Exec(context.Background(), createBooksTableSQL)
	if err != nil {
		log.Fatalf("Failed to create books table: %v", err)
	}
}
