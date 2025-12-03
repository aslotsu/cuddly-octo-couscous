package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookHandler struct {
	db *pgxpool.Pool
}

func NewBookHandler(db *pgxpool.Pool) *BookHandler {
	return &BookHandler{db: db}
}

// GetAllBooks retrieves all books
func (h *BookHandler) GetAllBooks(c *gin.Context) {
	query := `
		SELECT id, title, subtitle, author, isbn, description, publisher, publication_date,
		       pages, language, category, price, sale_price, stock_quantity, status,
		       cover_image, gallery_images, preview_url, purchase_links, tags,
		       is_featured, is_published, total_sales, average_rating, review_count,
		       created_by, created_at, updated_at
		FROM books
		ORDER BY created_at DESC
	`

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID, &book.Title, &book.Subtitle, &book.Author, &book.ISBN,
			&book.Description, &book.Publisher, &book.PublicationDate,
			&book.Pages, &book.Language, &book.Category, &book.Price, &book.SalePrice,
			&book.StockQuantity, &book.Status, &book.CoverImage, &book.GalleryImages,
			&book.PreviewURL, &book.PurchaseLinks, &book.Tags,
			&book.IsFeatured, &book.IsPublished, &book.TotalSales, &book.AverageRating,
			&book.ReviewCount, &book.CreatedBy, &book.CreatedAt, &book.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan book"})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []models.Book{}
	}

	c.JSON(http.StatusOK, books)
}

// GetBookByID retrieves a single book by ID
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	query := `
		SELECT id, title, subtitle, author, isbn, description, publisher, publication_date,
		       pages, language, category, price, sale_price, stock_quantity, status,
		       cover_image, gallery_images, preview_url, purchase_links, tags,
		       is_featured, is_published, total_sales, average_rating, review_count,
		       created_by, created_at, updated_at
		FROM books
		WHERE id = $1
	`

	var book models.Book
	err = h.db.QueryRow(context.Background(), query, id).Scan(
		&book.ID, &book.Title, &book.Subtitle, &book.Author, &book.ISBN,
		&book.Description, &book.Publisher, &book.PublicationDate,
		&book.Pages, &book.Language, &book.Category, &book.Price, &book.SalePrice,
		&book.StockQuantity, &book.Status, &book.CoverImage, &book.GalleryImages,
		&book.PreviewURL, &book.PurchaseLinks, &book.Tags,
		&book.IsFeatured, &book.IsPublished, &book.TotalSales, &book.AverageRating,
		&book.ReviewCount, &book.CreatedBy, &book.CreatedAt, &book.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}

// CreateBook creates a new book
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert JSONB fields to strings
	galleryImagesStr := marshalJSONB(req.GalleryImages)
	purchaseLinksStr := marshalJSONB(req.PurchaseLinks)
	tagsStr := marshalJSONB(req.Tags)

	query := `
		INSERT INTO books (
			title, subtitle, author, isbn, description, publisher, publication_date,
			pages, language, category, price, sale_price, stock_quantity, status,
			cover_image, gallery_images, preview_url, purchase_links, tags,
			is_featured, is_published, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22
		) RETURNING id
	`

	var id int
	err := h.db.QueryRow(
		context.Background(),
		query,
		req.Title, req.Subtitle, req.Author, req.ISBN, req.Description,
		req.Publisher, req.PublicationDate, req.Pages, req.Language, req.Category,
		req.Price, req.SalePrice, req.StockQuantity, req.Status,
		req.CoverImage, galleryImagesStr, req.PreviewURL, purchaseLinksStr, tagsStr,
		req.IsFeatured, req.IsPublished, req.CreatedBy,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if book exists
	var exists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Build dynamic update query
	query := "UPDATE books SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}
	argCount := 1

	if req.Title != "" {
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, req.Title)
		argCount++
	}
	if req.Subtitle != "" {
		query += ", subtitle = $" + strconv.Itoa(argCount)
		args = append(args, req.Subtitle)
		argCount++
	}
	if req.Author != "" {
		query += ", author = $" + strconv.Itoa(argCount)
		args = append(args, req.Author)
		argCount++
	}
	if req.ISBN != "" {
		query += ", isbn = $" + strconv.Itoa(argCount)
		args = append(args, req.ISBN)
		argCount++
	}
	if req.Description != "" {
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, req.Description)
		argCount++
	}
	if req.Publisher != "" {
		query += ", publisher = $" + strconv.Itoa(argCount)
		args = append(args, req.Publisher)
		argCount++
	}
	if req.PublicationDate != nil {
		query += ", publication_date = $" + strconv.Itoa(argCount)
		args = append(args, req.PublicationDate)
		argCount++
	}
	if req.Pages != nil {
		query += ", pages = $" + strconv.Itoa(argCount)
		args = append(args, req.Pages)
		argCount++
	}
	if req.Language != "" {
		query += ", language = $" + strconv.Itoa(argCount)
		args = append(args, req.Language)
		argCount++
	}
	if req.Category != "" {
		query += ", category = $" + strconv.Itoa(argCount)
		args = append(args, req.Category)
		argCount++
	}
	if req.Price != nil {
		query += ", price = $" + strconv.Itoa(argCount)
		args = append(args, req.Price)
		argCount++
	}
	if req.SalePrice != nil {
		query += ", sale_price = $" + strconv.Itoa(argCount)
		args = append(args, req.SalePrice)
		argCount++
	}
	if req.StockQuantity != nil {
		query += ", stock_quantity = $" + strconv.Itoa(argCount)
		args = append(args, req.StockQuantity)
		argCount++
	}
	if req.Status != "" {
		query += ", status = $" + strconv.Itoa(argCount)
		args = append(args, req.Status)
		argCount++
	}
	if req.CoverImage != "" {
		query += ", cover_image = $" + strconv.Itoa(argCount)
		args = append(args, req.CoverImage)
		argCount++
	}
	if req.GalleryImages != nil {
		query += ", gallery_images = $" + strconv.Itoa(argCount)
		args = append(args, marshalJSONB(req.GalleryImages))
		argCount++
	}
	if req.PreviewURL != "" {
		query += ", preview_url = $" + strconv.Itoa(argCount)
		args = append(args, req.PreviewURL)
		argCount++
	}
	if req.PurchaseLinks != nil {
		query += ", purchase_links = $" + strconv.Itoa(argCount)
		args = append(args, marshalJSONB(req.PurchaseLinks))
		argCount++
	}
	if req.Tags != nil {
		query += ", tags = $" + strconv.Itoa(argCount)
		args = append(args, marshalJSONB(req.Tags))
		argCount++
	}
	if req.IsFeatured != nil {
		query += ", is_featured = $" + strconv.Itoa(argCount)
		args = append(args, req.IsFeatured)
		argCount++
	}
	if req.IsPublished != nil {
		query += ", is_published = $" + strconv.Itoa(argCount)
		args = append(args, req.IsPublished)
		argCount++
	}
	if req.TotalSales != nil {
		query += ", total_sales = $" + strconv.Itoa(argCount)
		args = append(args, req.TotalSales)
		argCount++
	}
	if req.AverageRating != nil {
		query += ", average_rating = $" + strconv.Itoa(argCount)
		args = append(args, req.AverageRating)
		argCount++
	}
	if req.ReviewCount != nil {
		query += ", review_count = $" + strconv.Itoa(argCount)
		args = append(args, req.ReviewCount)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, id)

	_, err = h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
