package main

import (
	"log"
	"os"

	"github.com/aslotsu/monkreflections-form-api/config"
	"github.com/aslotsu/monkreflections-form-api/handlers"
	"github.com/aslotsu/monkreflections-form-api/middleware"
	"github.com/aslotsu/monkreflections-form-api/services"
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

	// Initialize S3 service (optional - for blog image uploads)
	var s3Service *services.S3Service
	s3Service, err = services.NewS3Service()
	if err != nil {
		log.Printf("Warning: S3 service not initialized (blog image uploads disabled): %v", err)
		s3Service = nil
	} else {
		log.Println("S3 service initialized successfully")
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(config.GetCORSConfig()))

	// Initialize handlers
	formHandler := handlers.NewFormHandler(db)
	blogHandler := handlers.NewBlogHandler(db, s3Service)
	authMiddleware := middleware.NewAuthMiddleware(db)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "monkreflections-form-api",
			"version": "1.0.0",
			"s3":      s3Service != nil,
		})
	})

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

		// Blog routes
		blogs := api.Group("/blogs")
		{
			blogs.GET("", blogHandler.GetAllBlogs)
			blogs.GET("/:id", blogHandler.GetBlogByID)

			// Protected blog routes (require API key)
			blogs.POST("", authMiddleware.RequireAPIKey(), blogHandler.CreateBlog)
			blogs.PUT("/:id", authMiddleware.RequireAPIKey(), blogHandler.UpdateBlog)
			blogs.DELETE("/:id", authMiddleware.RequireAPIKey(), blogHandler.DeleteBlog)
			blogs.POST("/:id/upload-image", authMiddleware.RequireAPIKey(), blogHandler.UploadBlogImage)
		}
	}

	// Start server (use Railway's PORT env var if available)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
