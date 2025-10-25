package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateCommentRequest struct {
	VideoID  uuid.UUID  `json:"videoId" validate:"required,uuid"`
	Content  string     `json:"content" validate:"required,min=1,max=1000"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"` // For nested comments
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// Response DTOs
type CommentResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"userId"`
	User      UserSummaryComment `json:"user"`
	VideoID   uuid.UUID          `json:"videoId"`
	ParentID  *uuid.UUID         `json:"parentId,omitempty"`
	Content   string             `json:"content"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Replies   []CommentResponse  `json:"replies,omitempty"` // Nested replies
}

type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	TotalCount int64             `json:"totalCount"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

// UserSummaryComment for comment responses
type UserSummaryComment struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Avatar    string    `json:"avatar,omitempty"`
}
