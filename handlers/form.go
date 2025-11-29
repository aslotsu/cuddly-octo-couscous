package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FormHandler struct {
	db *pgxpool.Pool
}

func NewFormHandler(db *pgxpool.Pool) *FormHandler {
	return &FormHandler{db: db}
}

// GetAllForms retrieves all forms
func (h *FormHandler) GetAllForms(c *gin.Context) {
	rows, err := h.db.Query(context.Background(), "SELECT id, title, data, created_at, updated_at FROM forms")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forms"})
		return
	}
	defer rows.Close()

	var forms []models.Form
	for rows.Next() {
		var form models.Form
		if err := rows.Scan(&form.ID, &form.Title, &form.Data, &form.CreatedAt, &form.UpdatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan form"})
			return
		}
		forms = append(forms, form)
	}

	if forms == nil {
		forms = []models.Form{}
	}

	c.JSON(http.StatusOK, forms)
}

// GetFormByID retrieves a single form by ID
func (h *FormHandler) GetFormByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var form models.Form
	err = h.db.QueryRow(
		context.Background(),
		"SELECT id, title, data, created_at, updated_at FROM forms WHERE id = $1",
		id,
	).Scan(&form.ID, &form.Title, &form.Data, &form.CreatedAt, &form.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	c.JSON(http.StatusOK, form)
}

// CreateForm creates a new form
func (h *FormHandler) CreateForm(c *gin.Context) {
	var req models.CreateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert data to JSON string
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
		return
	}

	var id int
	err = h.db.QueryRow(
		context.Background(),
		"INSERT INTO forms (title, data) VALUES ($1, $2) RETURNING id",
		req.Title,
		string(dataBytes),
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create form"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateForm updates an existing form
func (h *FormHandler) UpdateForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	var req models.UpdateFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if form exists
	var exists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM forms WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	// Build dynamic update query
	var dataStr string
	if req.Data != nil {
		dataBytes, err := json.Marshal(req.Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data format"})
			return
		}
		dataStr = string(dataBytes)
	}

	query := "UPDATE forms SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}

	if req.Title != "" {
		query += ", title = $" + strconv.Itoa(len(args)+1)
		args = append(args, req.Title)
	}
	if dataStr != "" {
		query += ", data = $" + strconv.Itoa(len(args)+1)
		args = append(args, dataStr)
	}

	query += " WHERE id = $" + strconv.Itoa(len(args)+1)
	args = append(args, id)

	_, err = h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form updated successfully"})
}

// DeleteForm deletes a form
func (h *FormHandler) DeleteForm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form ID"})
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM forms WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form deleted successfully"})
}
