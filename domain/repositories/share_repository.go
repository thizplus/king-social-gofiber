package repositories

import (
	"context"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type ShareRepository interface {
	// Basic operations
	Create(ctx context.Context, share *models.Share) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Share, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Query methods
	FindByVideoID(ctx context.Context, videoID uuid.UUID, offset, limit int) ([]*models.Share, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Share, int64, error)
	CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error)
	CountByPlatform(ctx context.Context, videoID uuid.UUID, platform string) (int64, error)

	// User-specific checks
	HasUserShared(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, platform string) (bool, error)
}
