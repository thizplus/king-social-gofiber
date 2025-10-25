package postgres

import (
	"context"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repositories.TagRepository {
	return &TagRepositoryImpl{db: db}
}

func (r *TagRepositoryImpl) Create(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.WithContext(ctx).First(&tag, "slug = ?", slug).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepositoryImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *TagRepositoryImpl) GetAll(ctx context.Context, offset, limit int, activeOnly bool) ([]*models.Tag, int, error) {
	var tags []*models.Tag
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Tag{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Order("usage_count DESC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&tags).Error

	return tags, int(total), err
}

func (r *TagRepositoryImpl) Update(ctx context.Context, tag *models.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *TagRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Tag{}, "id = ?", id).Error
}

func (r *TagRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Tag, int, error) {
	var tags []*models.Tag
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&models.Tag{}).
		Where("is_active = ?", true).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := dbQuery.Order("usage_count DESC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&tags).Error

	return tags, int(total), err
}

func (r *TagRepositoryImpl) IncrementUsageCount(ctx context.Context, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ?", tagID).
		Update("usage_count", gorm.Expr("usage_count + 1")).Error
}

func (r *TagRepositoryImpl) DecrementUsageCount(ctx context.Context, tagID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Tag{}).
		Where("id = ? AND usage_count > 0", tagID).
		Update("usage_count", gorm.Expr("usage_count - 1")).Error
}