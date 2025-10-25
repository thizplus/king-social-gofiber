package postgres

import (
	"context"
	"errors"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type notificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repositories.NotificationRepository {
	return &notificationRepositoryImpl{db: db}
}

func (r *notificationRepositoryImpl) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepositoryImpl) CreateBatch(ctx context.Context, notifications []models.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(notifications, 100).Error
}

func (r *notificationRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).
		Preload("Actor").
		First(&notification, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Preload("Actor")

	// Filter by type
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// Filter by read status
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	return notifications, totalCount, err
}

func (r *notificationRepositoryImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepositoryImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) MarkMultipleAsRead(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id IN ?", ids).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, id).Error
}

func (r *notificationRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Notification{}).Error
}
