package services

import (
	"context"

	"gofiber-social/domain/dto"

	"github.com/google/uuid"
)

type NotificationService interface {
	// Create notifications
	CreateTopicReplyNotification(ctx context.Context, topicID, replyUserID uuid.UUID) error
	CreateTopicLikeNotification(ctx context.Context, topicID, likerUserID uuid.UUID) error
	CreateVideoLikeNotification(ctx context.Context, videoID, likerUserID uuid.UUID) error
	CreateVideoCommentNotification(ctx context.Context, videoID, commenterUserID uuid.UUID) error
	CreateCommentReplyNotification(ctx context.Context, commentID, replierUserID uuid.UUID) error
	CreateReplyLikeNotification(ctx context.Context, replyID, likerUserID uuid.UUID) error
	CreateCommentLikeNotification(ctx context.Context, commentID, likerUserID uuid.UUID) error
	CreateNewFollowerNotification(ctx context.Context, followedUserID, followerUserID uuid.UUID) error

	// Read notifications
	GetNotifications(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) (*dto.NotificationListResponse, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error)

	// Update notifications
	MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	MarkMultipleAsRead(ctx context.Context, userID uuid.UUID, notificationIDs []uuid.UUID) (*dto.MarkAsReadResponse, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) (*dto.MarkAsReadResponse, error)

	// Delete notifications
	DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) error
}
