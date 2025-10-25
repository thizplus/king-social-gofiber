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

type videoRepositoryImpl struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) repositories.VideoRepository {
	return &videoRepositoryImpl{db: db}
}

func (r *videoRepositoryImpl) Create(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND is_active = ?", id, true).
		First(&video).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video not found")
		}
		return nil, err
	}
	return &video, nil
}

func (r *videoRepositoryImpl) Update(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Video{}, id).Error
}

func (r *videoRepositoryImpl) FindAll(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).
		Preload("User").
		Where("is_active = ?", true)

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderBy := "created_at DESC" // default: newest
	if params.SortBy == "oldest" {
		orderBy = "created_at ASC"
	} else if params.SortBy == "popular" {
		orderBy = "view_count DESC"
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

	err := query.Order(orderBy).Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).
		Preload("User").
		Where("user_id = ? AND is_active = ?", userID, true)

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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *videoRepositoryImpl) FindAllIncludingInactive(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).Preload("User")

	// Filter by IsActive if provided
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) SetActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}

func (r *videoRepositoryImpl) UpdateLikeCount(ctx context.Context, id uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		Update("like_count", count).Error
}

func (r *videoRepositoryImpl) UpdateCommentCount(ctx context.Context, id uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		Update("comment_count", count).Error
}

func (r *videoRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&count).Error
	return count, err
}

func (r *videoRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Video{}).
		Count(&count).Error
	return count, err
}
