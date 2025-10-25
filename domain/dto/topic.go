package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type CreateTopicRequest struct {
	ForumID   string   `json:"forumId" validate:"required,uuid4"`
	Title     string   `json:"title" validate:"required,min=5,max=200"`
	Content   string   `json:"content" validate:"required,min=10,max=10000"`
	Thumbnail string   `json:"thumbnail" validate:"omitempty,url"` // Optional thumbnail image URL
	TagIDs    []string `json:"tagIds" validate:"omitempty,dive,uuid4"` // Array of tag UUIDs
}

type UpdateTopicRequest struct {
	Title     string   `json:"title" validate:"omitempty,min=5,max=200"`
	Content   string   `json:"content" validate:"omitempty,min=10,max=10000"`
	Thumbnail string   `json:"thumbnail" validate:"omitempty,url"` // Optional thumbnail image URL
	TagIDs    []string `json:"tagIds" validate:"omitempty,dive,uuid4"` // Array of tag UUIDs
}

// Response DTOs
type TopicResponse struct {
	ID         uuid.UUID      `json:"id"`
	ForumID    uuid.UUID      `json:"forumId"`
	Forum      *ForumResponse `json:"forum,omitempty"`
	UserID     uuid.UUID      `json:"userId"`
	User       *UserResponseAdmin  `json:"user,omitempty"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	Thumbnail  string         `json:"thumbnail,omitempty"` // Optional thumbnail image URL
	ViewCount  int            `json:"viewCount"`
	ReplyCount int            `json:"replyCount"`
	IsPinned   bool           `json:"isPinned"`
	IsLocked   bool           `json:"isLocked"`
	Tags       []TagResponse  `json:"tags,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
}

type TopicListResponse struct {
	Topics []TopicResponse `json:"topics"`
	Meta   PaginationMeta  `json:"meta"`
}

type TopicDetailResponse struct {
	Topic   TopicResponse   `json:"topic"`
	Replies []ReplyResponse `json:"replies"`
}
