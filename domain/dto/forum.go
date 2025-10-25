package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateForumRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"required,min=3,max=100,lowercase,alphanum"`
	Description string `json:"description" validate:"required,min=10,max=500"`
	Icon        string `json:"icon" validate:"omitempty,url"`
	Order       int    `json:"order" validate:"min=0"`
}

type UpdateForumRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,min=10,max=500"`
	Icon        string `json:"icon" validate:"omitempty,url"`
	Order       int    `json:"order" validate:"omitempty,min=0"`
	IsActive    *bool  `json:"isActive"` // pointer เพื่อรองรับ true/false
}

type ReorderForumsRequest struct {
	ForumOrders []ForumOrder `json:"forumOrders" validate:"required,min=1,dive"`
}

type ForumOrder struct {
	ID    string `json:"id" validate:"required,uuid4"`
	Order int    `json:"order" validate:"min=0"`
}

// Response DTOs
type ForumResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Order       int       `json:"order"`
	IsActive    bool      `json:"isActive"`
	TopicCount  int       `json:"topicCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ForumListResponse struct {
	Forums []ForumResponse `json:"forums"`
	Meta   PaginationMeta  `json:"meta"`
}
