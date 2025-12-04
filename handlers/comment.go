package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentHandler struct {
	db *pgxpool.Pool
}

func NewCommentHandler(db *pgxpool.Pool) *CommentHandler {
	return &CommentHandler{db: db}
}

// GetCommentsByBlogID retrieves all approved comments for a blog post
func (h *CommentHandler) GetCommentsByBlogID(c *gin.Context) {
	blogID, err := strconv.Atoi(c.Param("blog_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	query := `
		SELECT id, blog_id, blog_slug, author_name, author_email, content, status, parent_id, created_at, updated_at
		FROM comments
		WHERE blog_id = $1 AND status = 'approved'
		ORDER BY created_at ASC
	`

	rows, err := h.db.Query(context.Background(), query, blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID, &comment.BlogID, &comment.BlogSlug, &comment.AuthorName,
			&comment.AuthorEmail, &comment.Content, &comment.Status, &comment.ParentID,
			&comment.CreatedAt, &comment.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// GetCommentsByBlogSlug retrieves all approved comments for a blog post by slug
func (h *CommentHandler) GetCommentsByBlogSlug(c *gin.Context) {
	slug := c.Param("slug")

	query := `
		SELECT id, blog_id, blog_slug, author_name, author_email, content, status, parent_id, created_at, updated_at
		FROM comments
		WHERE blog_slug = $1 AND status = 'approved'
		ORDER BY created_at ASC
	`

	rows, err := h.db.Query(context.Background(), query, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID, &comment.BlogID, &comment.BlogSlug, &comment.AuthorName,
			&comment.AuthorEmail, &comment.Content, &comment.Status, &comment.ParentID,
			&comment.CreatedAt, &comment.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// GetAllComments retrieves all comments (for admin, includes pending/rejected)
func (h *CommentHandler) GetAllComments(c *gin.Context) {
	status := c.Query("status") // optional filter by status

	query := `
		SELECT id, blog_id, blog_slug, author_name, author_email, content, status, parent_id, created_at, updated_at
		FROM comments
	`

	var args []interface{}
	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(
			&comment.ID, &comment.BlogID, &comment.BlogSlug, &comment.AuthorName,
			&comment.AuthorEmail, &comment.Content, &comment.Status, &comment.ParentID,
			&comment.CreatedAt, &comment.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan comment"})
			return
		}
		comments = append(comments, comment)
	}

	if comments == nil {
		comments = []models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment creates a new comment (status: pending by default)
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO comments (blog_id, blog_slug, author_name, author_email, content, parent_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'pending')
		RETURNING id
	`

	var id int
	err := h.db.QueryRow(
		context.Background(),
		query,
		req.BlogID, req.BlogSlug, req.AuthorName, req.AuthorEmail, req.Content, req.ParentID,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Comment submitted for moderation"})
}

// UpdateComment updates comment status or content (admin only)
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if comment exists
	var exists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Build dynamic update query
	query := "UPDATE comments SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}
	argCount := 1

	if req.Status != "" {
		query += ", status = $" + strconv.Itoa(argCount)
		args = append(args, req.Status)
		argCount++
	}
	if req.Content != "" {
		query += ", content = $" + strconv.Itoa(argCount)
		args = append(args, req.Content)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, id)

	_, err = h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

// DeleteComment deletes a comment (admin only)
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
