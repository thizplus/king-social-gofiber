package postgres

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
)

type activityLogRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) repositories.ActivityLogRepository {
	return &activityLogRepositoryImpl{db: db}
}

func (r *activityLogRepositoryImpl) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *activityLogRepositoryImpl) FindByAdminID(ctx context.Context, adminID uuid.UUID, page, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{}).
		Where("admin_id = ?", adminID).
		Preload("Admin")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, totalCount, err
}

func (r *activityLogRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{}).
		Preload("Admin")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, totalCount, err
}
