package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateReplyRequest struct {
	Content  string  `json:"content" validate:"required,min=1,max=5000"`
	ParentID *string `json:"parentId" validate:"omitempty,uuid4"` // สำหรับ nested reply
}

type UpdateReplyRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

// Response DTOs
type ReplyResponse struct {
	ID        uuid.UUID       `json:"id"`
	TopicID   uuid.UUID       `json:"topicId"`
	UserID    uuid.UUID       `json:"userId"`
	User      *UserResponseAdmin   `json:"user,omitempty"`
	ParentID  *uuid.UUID      `json:"parentId,omitempty"`
	Content   string          `json:"content"`
	Replies   []ReplyResponse `json:"replies,omitempty"` // Nested replies
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type ReplyListResponse struct {
	Replies []ReplyResponse `json:"replies"`
	Meta    PaginationMeta  `json:"meta"`
}
