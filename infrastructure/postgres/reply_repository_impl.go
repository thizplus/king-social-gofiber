package postgres

import (
	"context"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReplyRepositoryImpl struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) repositories.ReplyRepository {
	return &ReplyRepositoryImpl{db: db}
}

func (r *ReplyRepositoryImpl) Create(ctx context.Context, reply *models.Reply) error {
	return r.db.WithContext(ctx).Create(reply).Error
}

func (r *ReplyRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Reply, error) {
	var reply models.Reply
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Topic").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&reply).Error
	return &reply, err
}

func (r *ReplyRepositoryImpl) GetByTopicID(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*models.Reply, error) {
	var replies []*models.Reply
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Replies.User"). // Nested replies
		Preload("Replies.Replies.User"). // Level 2 nested
		Where("topic_id = ?", topicID).
		Where("parent_id IS NULL"). // Only top-level replies
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&replies).Error
	return replies, err
}

func (r *ReplyRepositoryImpl) GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*models.Reply, error) {
	var replies []*models.Reply
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("parent_id = ?", parentID).
		Where("deleted_at IS NULL").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

func (r *ReplyRepositoryImpl) Update(ctx context.Context, id uuid.UUID, reply *models.Reply) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Updates(reply).Error
}

func (r *ReplyRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&models.Reply{}).Error
}

func (r *ReplyRepositoryImpl) Count(ctx context.Context, topicID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Reply{}).
		Where("topic_id = ?", topicID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}
