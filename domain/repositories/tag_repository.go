package repositories

import (
	"context"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type TagRepository interface {
	Create(ctx context.Context, tag *models.Tag) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*models.Tag, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Tag, error)
	GetAll(ctx context.Context, offset, limit int, activeOnly bool) ([]*models.Tag, int, error)
	Update(ctx context.Context, tag *models.Tag) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Tag, int, error)
	IncrementUsageCount(ctx context.Context, tagID uuid.UUID) error
	DecrementUsageCount(ctx context.Context, tagID uuid.UUID) error
}