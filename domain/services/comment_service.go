package services

import (
	"context"
	"gofiber-social/domain/dto"

	"github.com/google/uuid"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetCommentsByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) (*dto.CommentListResponse, error)
	UpdateComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error
	DeleteCommentByAdmin(ctx context.Context, commentID uuid.UUID) error
}
