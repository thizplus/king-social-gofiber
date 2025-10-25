package repositories

import (
	"context"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type NotificationRepository interface {
	// Create
	Create(ctx context.Context, notification *models.Notification) error
	CreateBatch(ctx context.Context, notifications []models.Notification) error

	// Read
	FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) ([]models.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// Update
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkMultipleAsRead(ctx context.Context, ids []uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
