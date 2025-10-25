package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ReplyService interface {
	CreateReply(ctx context.Context, topicID, userID uuid.UUID, req *dto.CreateReplyRequest) (*models.Reply, error)
	GetReplies(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*dto.ReplyResponse, int64, error)
	UpdateReply(ctx context.Context, replyID, userID uuid.UUID, req *dto.UpdateReplyRequest) (*models.Reply, error)
	DeleteReply(ctx context.Context, replyID, userID uuid.UUID) error
	DeleteReplyByAdmin(ctx context.Context, replyID uuid.UUID) error
}
