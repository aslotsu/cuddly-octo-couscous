package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/aslotsu/monkreflections-form-api/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogHandler struct {
	db       *pgxpool.Pool
	s3Service *services.S3Service
}

func NewBlogHandler(db *pgxpool.Pool, s3Service *services.S3Service) *BlogHandler {
	return &BlogHandler{
		db:       db,
		s3Service: s3Service,
	}
}

// GetAllBlogs retrieves all blogs
func (bh *BlogHandler) GetAllBlogs(c *gin.Context) {
	rows, err := bh.db.Query(context.Background(), "SELECT id, title, content, author, created_at, updated_at FROM blogs ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs"})
		return
	}
	defer rows.Close()

	var blogs []models.Blog
	for rows.Next() {
		var blog models.Blog
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan blog"})
			return
		}
		blogs = append(blogs, blog)
	}

	if blogs == nil {
		blogs = []models.Blog{}
	}

	c.JSON(http.StatusOK, blogs)
}

// GetBlogByID retrieves a single blog by ID
func (bh *BlogHandler) GetBlogByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	err = bh.db.QueryRow(
		context.Background(),
		"SELECT id, title, content, author, created_at, updated_at FROM blogs WHERE id = $1",
		id,
	).Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Author, &blog.CreatedAt, &blog.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// CreateBlog creates a new blog
func (bh *BlogHandler) CreateBlog(c *gin.Context) {
	var req models.CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert content to JSON string
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
		return
	}

	var id int
	err = bh.db.QueryRow(
		context.Background(),
		"INSERT INTO blogs (title, content, author) VALUES ($1, $2, $3) RETURNING id",
		req.Title,
		string(contentBytes),
		req.Author,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateBlog updates an existing blog
func (bh *BlogHandler) UpdateBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var req models.UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if blog exists
	var exists bool
	err = bh.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM blogs WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	// Build dynamic update query
	var contentStr string
	if req.Content != nil {
		contentBytes, err := json.Marshal(req.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content format"})
			return
		}
		contentStr = string(contentBytes)
	}

	query := "UPDATE blogs SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}

	if req.Title != "" {
		query += ", title = $" + strconv.Itoa(len(args)+1)
		args = append(args, req.Title)
	}
	if contentStr != "" {
		query += ", content = $" + strconv.Itoa(len(args)+1)
		args = append(args, contentStr)
	}
	if req.Author != "" {
		query += ", author = $" + strconv.Itoa(len(args)+1)
		args = append(args, req.Author)
	}

	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	_, err = bh.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog updated successfully"})
}

// DeleteBlog deletes a blog
func (bh *BlogHandler) DeleteBlog(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	result, err := bh.db.Exec(context.Background(), "DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

// UploadBlogImage handles image uploads for a blog
func (bh *BlogHandler) UploadBlogImage(c *gin.Context) {
	// Check if S3 service is available
	if bh.s3Service == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Image upload service is not available. S3 is not configured.",
		})
		return
	}

	blogID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	// Check if blog exists
	var exists bool
	err = bh.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM blogs WHERE id = $1)", blogID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	defer file.Close()

	// Upload to S3
	imageKey, imageURL, err := bh.s3Service.UploadImage(file, header.Size, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image: " + err.Error()})
		return
	}

	// Store image reference in database
	var imageID int
	err = bh.db.QueryRow(
		context.Background(),
		"INSERT INTO blog_images (blog_id, image_key, image_url, alt_text) VALUES ($1, $2, $3, $4) RETURNING id",
		blogID, imageKey, imageURL, c.PostForm("alt_text"),
	).Scan(&imageID)

	if err != nil {
		// If database insertion fails, try to delete the uploaded image from S3
		deleteErr := bh.s3Service.DeleteImage(imageKey)
		if deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":          "Failed to store image reference in database",
				"cleanup_error":  "Also failed to cleanup S3 image",
				"cleanup_detail": deleteErr.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":         "Failed to store image reference in database",
				"cleanup_msg":   "Successfully cleaned up S3 image",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        imageID,
		"image_url": imageURL,
		"image_key": imageKey,
	})
}