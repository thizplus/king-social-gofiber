package repositories

import (
	"context"

	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type ForumRepository interface {
	Create(ctx context.Context, forum *models.Forum) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Forum, error)
	GetBySlug(ctx context.Context, slug string) (*models.Forum, error)
	GetAll(ctx context.Context, includeInactive bool) ([]*models.Forum, error)
	GetActive(ctx context.Context) ([]*models.Forum, error)
	Update(ctx context.Context, id uuid.UUID, forum *models.Forum) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateOrder(ctx context.Context, id uuid.UUID, order int) error
	IncrementTopicCount(ctx context.Context, id uuid.UUID) error
	DecrementTopicCount(ctx context.Context, id uuid.UUID) error
	SyncTopicCount(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}
