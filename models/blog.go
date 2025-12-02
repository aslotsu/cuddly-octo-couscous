package models

import (
	"time"
)

type Blog struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"` // JSONB content as string
	Author    string    `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BlogImage struct {
	ID       int       `json:"id"`
	BlogID   int       `json:"blog_id"`
	ImageKey string    `json:"image_key"` // S3 object key
	ImageURL string    `json:"image_url"` // Public URL
	AltText  string    `json:"alt_text,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ApiKey struct {
	ID      int       `json:"id"`
	KeyHash string    `json:"key_hash"`
	Name    string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBlogRequest struct {
	Title   string                 `json:"title" binding:"required"`
	Content map[string]any         `json:"content" binding:"required"`
	Author  string                 `json:"author,omitempty"`
}

type UpdateBlogRequest struct {
	Title   string                 `json:"title,omitempty"`
	Content map[string]any         `json:"content,omitempty"`
	Author  string                 `json:"author,omitempty"`
}