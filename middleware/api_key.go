package middleware

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthMiddleware struct {
	db *pgxpool.Pool
}

func NewAuthMiddleware(db *pgxpool.Pool) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

func (am *AuthMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Expect format: "Bearer YOUR_API_KEY" or just "YOUR_API_KEY"
		var apiKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			apiKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			apiKey = authHeader
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Hash the API key for comparison
		hash := sha256.Sum256([]byte(apiKey))
		hashString := hex.EncodeToString(hash[:])

		// Check if the hashed API key exists in the database
		var apiKEyRecord models.ApiKey
		err := am.db.QueryRow(
			context.Background(),
			"SELECT id, key_hash, name, created_at FROM api_keys WHERE key_hash = $1",
			hashString,
		).Scan(&apiKEyRecord.ID, &apiKEyRecord.KeyHash, &apiKEyRecord.Name, &apiKEyRecord.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error validating API key"})
			}
			c.Abort()
			return
		}

		// API key is valid, continue with the request
		c.Next()
	}
}