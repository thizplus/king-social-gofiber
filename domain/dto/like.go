package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type LikeTopicRequest struct {
	TopicID uuid.UUID `json:"topicId" validate:"required,uuid"`
}

type LikeVideoRequest struct {
	VideoID uuid.UUID `json:"videoId" validate:"required,uuid"`
}

// Response DTOs
type LikeResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	TopicID   *uuid.UUID `json:"topicId,omitempty"`
	VideoID   *uuid.UUID `json:"videoId,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type LikeStatusResponse struct {
	IsLiked   bool  `json:"isLiked"`
	LikeCount int64 `json:"likeCount"`
}
