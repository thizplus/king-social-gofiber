package repositories

import (
	"context"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type LikeRepository interface {
	// Topic Likes
	LikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*models.Like, error)
	UnlikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) error
	IsTopicLikedByUser(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (bool, error)
	CountTopicLikes(ctx context.Context, topicID uuid.UUID) (int64, error)
	GetTopicLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error)

	// Video Likes
	LikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*models.Like, error)
	UnlikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) error
	IsVideoLikedByUser(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (bool, error)
	CountVideoLikes(ctx context.Context, videoID uuid.UUID) (int64, error)
	GetVideoLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error)

	// Reply Likes
	LikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*models.Like, error)
	UnlikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) error
	IsReplyLikedByUser(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (bool, error)
	CountReplyLikes(ctx context.Context, replyID uuid.UUID) (int64, error)
	GetReplyLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error)

	// Comment Likes
	LikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*models.Like, error)
	UnlikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error
	IsCommentLikedByUser(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (bool, error)
	CountCommentLikes(ctx context.Context, commentID uuid.UUID) (int64, error)
	GetCommentLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error)

	// General
	GetByID(ctx context.Context, id uuid.UUID) (*models.Like, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
