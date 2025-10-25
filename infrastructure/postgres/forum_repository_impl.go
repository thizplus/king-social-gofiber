package postgres

import (
	"context"

	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumRepositoryImpl struct {
	db *gorm.DB
}

func NewForumRepository(db *gorm.DB) repositories.ForumRepository {
	return &ForumRepositoryImpl{db: db}
}

func (r *ForumRepositoryImpl) Create(ctx context.Context, forum *models.Forum) error {
	return r.db.WithContext(ctx).Create(forum).Error
}

func (r *ForumRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Forum, error) {
	var forum models.Forum
	err := r.db.WithContext(ctx).
		Preload("Admin").
		Where("id = ?", id).
		First(&forum).Error
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (r *ForumRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.Forum, error) {
	var forum models.Forum
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&forum).Error
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (r *ForumRepositoryImpl) GetAll(ctx context.Context, includeInactive bool) ([]*models.Forum, error) {
	var forums []*models.Forum
	query := r.db.WithContext(ctx)

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("\"order\" ASC, created_at ASC").Find(&forums).Error
	return forums, err
}

func (r *ForumRepositoryImpl) GetActive(ctx context.Context) ([]*models.Forum, error) {
	return r.GetAll(ctx, false)
}

func (r *ForumRepositoryImpl) Update(ctx context.Context, id uuid.UUID, forum *models.Forum) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(forum).Error
}

func (r *ForumRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.Forum{}, "id = ?", id).Error
}

func (r *ForumRepositoryImpl) UpdateOrder(ctx context.Context, id uuid.UUID, order int) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		Update("order", order).Error
}

func (r *ForumRepositoryImpl) IncrementTopicCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		UpdateColumn("topic_count", gorm.Expr("topic_count + ?", 1)).Error
}

func (r *ForumRepositoryImpl) DecrementTopicCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		Where("topic_count > ?", 0).
		UpdateColumn("topic_count", gorm.Expr("topic_count - ?", 1)).Error
}

func (r *ForumRepositoryImpl) SyncTopicCount(ctx context.Context, id uuid.UUID) error {
	// Count actual non-deleted topics for this forum
	var actualCount int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("forum_id = ?", id).
		Where("deleted_at IS NULL").
		Count(&actualCount).Error

	if err != nil {
		return err
	}

	// Update forum's topic_count to match actual count
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		UpdateColumn("topic_count", actualCount).Error
}

func (r *ForumRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Count(&count).Error
	return count, err
}
