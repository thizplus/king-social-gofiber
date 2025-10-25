package postgres

import (
	"context"
	"errors"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repositories.CommentRepository {
	return &commentRepositoryImpl{db: db}
}

// Create creates a new comment
func (r *commentRepositoryImpl) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// GetByID retrieves a comment by ID
func (r *commentRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Parent").
		Preload("Parent.User").
		Where("id = ?", id).
		First(&comment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}

	return &comment, nil
}

// Update updates a comment
func (r *commentRepositoryImpl) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// Delete soft deletes a comment
func (r *commentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Comment{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("comment not found")
	}

	return nil
}

// FindByVideoID retrieves all top-level comments for a video (excluding replies)
func (r *commentRepositoryImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var totalCount int64

	// Build query for top-level comments only (parent_id IS NULL)
	query := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("video_id = ? AND parent_id IS NULL", videoID)

	// Count total top-level comments
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get comments with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			// Preload first level of replies with their users
			return db.Preload("User").Order("created_at ASC").Limit(3)
		}).
		Where("video_id = ? AND parent_id IS NULL", videoID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, totalCount, err
}

// FindReplies retrieves all replies for a parent comment (nested comments)
func (r *commentRepositoryImpl) FindReplies(ctx context.Context, parentID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error) {
	var replies []*models.Comment
	var totalCount int64

	// Build query for replies
	query := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("parent_id = ?", parentID)

	// Count total replies
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get replies with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Parent").
		Preload("Parent.User").
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&replies).Error

	return replies, totalCount, err
}

// CountByVideoID returns the total number of comments for a video (including replies)
func (r *commentRepositoryImpl) CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	return count, err
}

// CountReplies returns the number of replies for a parent comment
func (r *commentRepositoryImpl) CountReplies(ctx context.Context, parentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("parent_id = ?", parentID).
		Count(&count).Error

	return count, err
}

// FindByUserID retrieves all comments by a user
func (r *commentRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var totalCount int64

	// Build query
	query := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("user_id = ?", userID)

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get comments with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Video.User").
		Preload("Parent").
		Preload("Parent.User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, totalCount, err
}

// Admin methods

func (r *commentRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Count(&count).Error
	return count, err
}
