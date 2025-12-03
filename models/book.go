package models

import "time"

type Book struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle,omitempty"`
	Author          string    `json:"author"`
	ISBN            string    `json:"isbn,omitempty"`
	Description     string    `json:"description"`
	Publisher       string    `json:"publisher,omitempty"`
	PublicationDate time.Time `json:"publication_date,omitempty"`
	Pages           int       `json:"pages"`
	Language        string    `json:"language"`
	Category        string    `json:"category"` // spiritual, devotional, biblical, etc
	Price           float64   `json:"price"`
	SalePrice       *float64  `json:"sale_price,omitempty"`
	StockQuantity   int       `json:"stock_quantity"`
	Status          string    `json:"status"` // available, out_of_stock, pre_order, discontinued
	CoverImage      string    `json:"cover_image,omitempty"`
	GalleryImages   string    `json:"gallery_images,omitempty"` // JSONB
	PreviewURL      string    `json:"preview_url,omitempty"`
	PurchaseLinks   string    `json:"purchase_links,omitempty"` // JSONB - Amazon, Kindle, etc
	Tags            string    `json:"tags,omitempty"`           // JSONB
	IsFeatured      bool      `json:"is_featured"`
	IsPublished     bool      `json:"is_published"`
	TotalSales      int       `json:"total_sales"`
	AverageRating   float64   `json:"average_rating"`
	ReviewCount     int       `json:"review_count"`
	CreatedBy       string    `json:"created_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateBookRequest struct {
	Title           string    `json:"title" binding:"required"`
	Subtitle        string    `json:"subtitle,omitempty"`
	Author          string    `json:"author" binding:"required"`
	ISBN            string    `json:"isbn,omitempty"`
	Description     string    `json:"description" binding:"required"`
	Publisher       string    `json:"publisher,omitempty"`
	PublicationDate time.Time `json:"publication_date,omitempty"`
	Pages           int       `json:"pages"`
	Language        string    `json:"language"`
	Category        string    `json:"category" binding:"required"`
	Price           float64   `json:"price" binding:"required"`
	SalePrice       *float64  `json:"sale_price,omitempty"`
	StockQuantity   int       `json:"stock_quantity"`
	Status          string    `json:"status"`
	CoverImage      string    `json:"cover_image,omitempty"`
	GalleryImages   any       `json:"gallery_images,omitempty"`
	PreviewURL      string    `json:"preview_url,omitempty"`
	PurchaseLinks   any       `json:"purchase_links,omitempty"`
	Tags            any       `json:"tags,omitempty"`
	IsFeatured      bool      `json:"is_featured"`
	IsPublished     bool      `json:"is_published"`
	CreatedBy       string    `json:"created_by,omitempty"`
}

type UpdateBookRequest struct {
	Title           string     `json:"title,omitempty"`
	Subtitle        string     `json:"subtitle,omitempty"`
	Author          string     `json:"author,omitempty"`
	ISBN            string     `json:"isbn,omitempty"`
	Description     string     `json:"description,omitempty"`
	Publisher       string     `json:"publisher,omitempty"`
	PublicationDate *time.Time `json:"publication_date,omitempty"`
	Pages           *int       `json:"pages,omitempty"`
	Language        string     `json:"language,omitempty"`
	Category        string     `json:"category,omitempty"`
	Price           *float64   `json:"price,omitempty"`
	SalePrice       *float64   `json:"sale_price,omitempty"`
	StockQuantity   *int       `json:"stock_quantity,omitempty"`
	Status          string     `json:"status,omitempty"`
	CoverImage      string     `json:"cover_image,omitempty"`
	GalleryImages   any        `json:"gallery_images,omitempty"`
	PreviewURL      string     `json:"preview_url,omitempty"`
	PurchaseLinks   any        `json:"purchase_links,omitempty"`
	Tags            any        `json:"tags,omitempty"`
	IsFeatured      *bool      `json:"is_featured,omitempty"`
	IsPublished     *bool      `json:"is_published,omitempty"`
	TotalSales      *int       `json:"total_sales,omitempty"`
	AverageRating   *float64   `json:"average_rating,omitempty"`
	ReviewCount     *int       `json:"review_count,omitempty"`
}
