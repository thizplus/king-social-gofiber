package postgres

import (
	"context"
	"errors"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type shareRepositoryImpl struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) repositories.ShareRepository {
	return &shareRepositoryImpl{db: db}
}

// Create creates a new share record
func (r *shareRepositoryImpl) Create(ctx context.Context, share *models.Share) error {
	return r.db.WithContext(ctx).Create(share).Error
}

// GetByID retrieves a share by ID
func (r *shareRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Share, error) {
	var share models.Share
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Video.User").
		Where("id = ?", id).
		First(&share).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("share not found")
		}
		return nil, err
	}

	return &share, nil
}

// Delete removes a share by ID
func (r *shareRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Share{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("share not found")
	}

	return nil
}

// FindByVideoID retrieves all shares for a video
func (r *shareRepositoryImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID, offset, limit int) ([]*models.Share, int64, error) {
	var shares []*models.Share
	var totalCount int64

	// Build query
	query := r.db.WithContext(ctx).
		Model(&models.Share{}).
		Where("video_id = ?", videoID)

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get shares with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Video.User").
		Where("video_id = ?", videoID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&shares).Error

	return shares, totalCount, err
}

// FindByUserID retrieves all shares by a user
func (r *shareRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Share, int64, error) {
	var shares []*models.Share
	var totalCount int64

	// Build query
	query := r.db.WithContext(ctx).
		Model(&models.Share{}).
		Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get shares with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Video.User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&shares).Error

	return shares, totalCount, err
}

// CountByVideoID returns the total number of shares for a video
func (r *shareRepositoryImpl) CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Share{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	return count, err
}

// CountByPlatform returns the number of shares for a video on a specific platform
func (r *shareRepositoryImpl) CountByPlatform(ctx context.Context, videoID uuid.UUID, platform string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Share{}).
		Where("video_id = ? AND platform = ?", videoID, platform).
		Count(&count).Error

	return count, err
}

// HasUserShared checks if a user has already shared a video on a specific platform
func (r *shareRepositoryImpl) HasUserShared(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, platform string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Share{}).
		Where("user_id = ? AND video_id = ? AND platform = ?", userID, videoID, platform).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
