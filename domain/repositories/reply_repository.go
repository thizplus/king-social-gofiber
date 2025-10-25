package repositories

import (
	"context"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ReplyRepository interface {
	Create(ctx context.Context, reply *models.Reply) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Reply, error)
	GetByTopicID(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*models.Reply, error)
	GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*models.Reply, error)
	Update(ctx context.Context, id uuid.UUID, reply *models.Reply) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, topicID uuid.UUID) (int64, error)
}
