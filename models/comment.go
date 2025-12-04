package models

import "time"

type Comment struct {
	ID          int       `json:"id"`
	BlogID      int       `json:"blog_id"`
	BlogSlug    string    `json:"blog_slug,omitempty"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Content     string    `json:"content"`
	Status      string    `json:"status"` // pending, approved, rejected, spam
	ParentID    *int      `json:"parent_id,omitempty"` // for nested replies
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	BlogID      int    `json:"blog_id" binding:"required"`
	BlogSlug    string `json:"blog_slug,omitempty"`
	AuthorName  string `json:"author_name" binding:"required"`
	AuthorEmail string `json:"author_email" binding:"required,email"`
	Content     string `json:"content" binding:"required"`
	ParentID    *int   `json:"parent_id,omitempty"`
}

type UpdateCommentRequest struct {
	Status  string `json:"status,omitempty"`
	Content string `json:"content,omitempty"`
}
