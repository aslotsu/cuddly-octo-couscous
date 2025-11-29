package main

import (
	"log"
	"os"

	"github.com/aslotsu/monkreflections-form-api/config"
	"github.com/aslotsu/monkreflections-form-api/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Connect to database
	db := config.ConnectDBFromURL(databaseURL)
	defer db.Close()

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(config.GetCORSConfig()))

	// Initialize handlers
	formHandler := handlers.NewFormHandler(db)

	// Register routes
	api := router.Group("/api")
	{
		forms := api.Group("/forms")
		{
			forms.GET("", formHandler.GetAllForms)
			forms.POST("", formHandler.CreateForm)
			forms.GET("/:id", formHandler.GetFormByID)
			forms.PUT("/:id", formHandler.UpdateForm)
			forms.DELETE("/:id", formHandler.DeleteForm)
		}
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
