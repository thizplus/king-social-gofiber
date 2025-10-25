package repositories

import (
	"context"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type CommentRepository interface {
	// Basic CRUD
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query methods
	FindByVideoID(ctx context.Context, videoID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error)
	FindReplies(ctx context.Context, parentID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error)
	CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error)
	CountReplies(ctx context.Context, parentID uuid.UUID) (int64, error)

	// User-specific queries
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error)

	// Admin methods
	GetTotalCount(ctx context.Context) (int64, error)
}
