package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func generateRandomKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Generate API key
	apiKey, err := generateRandomKey()
	if err != nil {
		log.Fatalf("Failed to generate API key: %v", err)
	}

	// Hash the key
	keyHash := hashKey(apiKey)

	// Insert into database
	var id int
	err = pool.QueryRow(
		context.Background(),
		"INSERT INTO api_keys (key_hash, name) VALUES ($1, $2) RETURNING id",
		keyHash,
		"Admin Dashboard Key",
	).Scan(&id)

	if err != nil {
		log.Fatalf("Failed to insert API key: %v", err)
	}

	fmt.Println("=====================================")
	fmt.Println("API Key generated successfully!")
	fmt.Println("=====================================")
	fmt.Printf("Key ID: %d\n", id)
	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Println("=====================================")
	fmt.Println("IMPORTANT: Copy this key now!")
	fmt.Println("Add it to your .env file:")
	fmt.Printf("NEXT_PUBLIC_API_KEY=%s\n", apiKey)
	fmt.Println("=====================================")
}
