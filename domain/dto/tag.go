package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateTagRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Slug        string `json:"slug" validate:"required,min=1,max=50,lowercase,alphanum"`
	Description string `json:"description" validate:"omitempty,max=255"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
}

type UpdateTagRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=50"`
	Description string `json:"description" validate:"omitempty,max=255"`
	Color       string `json:"color" validate:"omitempty,hexcolor"`
	IsActive    *bool  `json:"isActive"`
}

// Response DTOs
type TagResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	IsActive    bool      `json:"isActive"`
	UsageCount  int       `json:"usageCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TagListResponse struct {
	Tags []TagResponse  `json:"tags"`
	Meta PaginationMeta `json:"meta"`
}