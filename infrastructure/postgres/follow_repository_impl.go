package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gorm.io/gorm"
)

type followRepositoryImpl struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) repositories.FollowRepository {
	return &followRepositoryImpl{db: db}
}

// Follow creates a new follow relationship
func (r *followRepositoryImpl) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return errors.New("users cannot follow themselves")
	}

	// Check if already following
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("already following this user")
	}

	follow := &models.Follow{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	}

	return r.db.WithContext(ctx).Create(follow).Error
}

// Unfollow removes a follow relationship
func (r *followRepositoryImpl) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}

// IsFollowing checks if a user is following another user
func (r *followRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetFollowers retrieves the list of users following a specific user
func (r *followRepositoryImpl) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloaded follower user data
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Preload("Follower").
		Preload("Following").
		Where("following_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

// GetFollowing retrieves the list of users that a specific user is following
func (r *followRepositoryImpl) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloaded following user data
	offset := (page - 1) * limit
	if err := r.db.WithContext(ctx).
		Preload("Follower").
		Preload("Following").
		Where("follower_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&follows).Error; err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

// GetFollowerCount returns the number of followers for a user
func (r *followRepositoryImpl) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error

	return count, err
}

// GetFollowingCount returns the number of users a user is following
func (r *followRepositoryImpl) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error

	return count, err
}

// FindByID retrieves a follow relationship by its ID
func (r *followRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Follow, error) {
	var follow models.Follow
	err := r.db.WithContext(ctx).
		Preload("Follower").
		Preload("Following").
		Where("id = ?", id).
		First(&follow).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("follow relationship not found")
		}
		return nil, err
	}

	return &follow, nil
}

// Delete removes a follow relationship by its ID
func (r *followRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Follow{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("follow relationship not found")
	}

	return nil
}
