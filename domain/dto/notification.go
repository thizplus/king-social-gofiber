package dto

import (
	"time"

	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

// Request DTOs
type MarkAsReadRequest struct {
	NotificationIDs []uuid.UUID `json:"notificationIds" validate:"required,min=1"`
}

type NotificationQueryParams struct {
	Type   string `query:"type" validate:"omitempty,oneof=topic_reply topic_like video_like video_comment comment_reply new_follower"`
	IsRead *bool  `query:"isRead"`
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

// Response DTOs
type ActorSummary struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar,omitempty"`
}

type NotificationResponse struct {
	ID         uuid.UUID               `json:"id"`
	UserID     uuid.UUID               `json:"userId"`
	Actor      ActorSummary            `json:"actor"`
	Type       models.NotificationType `json:"type"`
	ResourceID *uuid.UUID              `json:"resourceId,omitempty"`
	Message    string                  `json:"message"`
	IsRead     bool                    `json:"isRead"`
	CreatedAt  time.Time               `json:"createdAt"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	TotalCount    int64                  `json:"totalCount"`
	UnreadCount   int64                  `json:"unreadCount"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"totalPages"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type MarkAsReadResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // จำนวนที่ mark as read
}
