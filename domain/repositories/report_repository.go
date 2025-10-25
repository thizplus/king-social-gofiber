package repositories

import (
	"context"
	"github.com/google/uuid"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
)

type ReportRepository interface {
	Create(ctx context.Context, report *models.Report) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Report, error)
	FindAll(ctx context.Context, params *dto.ReportQueryParams) ([]models.Report, int64, error)
	Update(ctx context.Context, report *models.Report) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCount(ctx context.Context) (int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
}
