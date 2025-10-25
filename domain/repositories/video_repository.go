package repositories

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type VideoRepository interface {
	// Basic CRUD
	Create(ctx context.Context, video *models.Video) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	Update(ctx context.Context, video *models.Video) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List & Query
	FindAll(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) ([]models.Video, int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)

	// View Count
	IncrementViewCount(ctx context.Context, id uuid.UUID) error

	// Like & Comment Count
	UpdateLikeCount(ctx context.Context, id uuid.UUID, count int) error
	UpdateCommentCount(ctx context.Context, id uuid.UUID, count int) error

	// Admin
	FindAllIncludingInactive(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error)
	SetActive(ctx context.Context, id uuid.UUID, isActive bool) error
	GetTotalCount(ctx context.Context) (int64, error)
}
