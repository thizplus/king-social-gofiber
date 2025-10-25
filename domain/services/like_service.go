package services

import (
	"context"
	"gofiber-social/domain/dto"

	"github.com/google/uuid"
)

type LikeService interface {
	// Topic Likes
	LikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error)
	UnlikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error)
	GetTopicLikeStatus(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error)

	// Video Likes
	LikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error)
	UnlikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error)
	GetVideoLikeStatus(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error)

	// Reply Likes
	LikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error)
	UnlikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error)
	GetReplyLikeStatus(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error)

	// Comment Likes
	LikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error)
	UnlikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error)
	GetCommentLikeStatus(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error)
}
