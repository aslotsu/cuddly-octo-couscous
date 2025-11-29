package models

import (
	"time"
)

type Form struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFormRequest struct {
	Title string         `json:"title" binding:"required"`
	Data  map[string]any `json:"data" binding:"required"`
}

type UpdateFormRequest struct {
	Title string         `json:"title"`
	Data  map[string]any `json:"data"`
}
