package repositories

import (
	"context"
	"github.com/google/uuid"
	"gofiber-social/domain/models"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
	FindByAdminID(ctx context.Context, adminID uuid.UUID, page, limit int) ([]models.ActivityLog, int64, error)
	FindAll(ctx context.Context, page, limit int) ([]models.ActivityLog, int64, error)
}
