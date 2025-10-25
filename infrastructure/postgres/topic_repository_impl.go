package postgres

import (
	"context"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicRepositoryImpl struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) repositories.TopicRepository {
	return &TopicRepositoryImpl{db: db}
}

func (r *TopicRepositoryImpl) Create(ctx context.Context, topic *models.Topic) error {
	return r.db.WithContext(ctx).Create(topic).Error
}

func (r *TopicRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error) {
	var topic models.Topic
	err := r.db.WithContext(ctx).
		Preload("Forum").
		Preload("User").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&topic).Error
	return &topic, err
}

func (r *TopicRepositoryImpl) GetByForumID(ctx context.Context, forumID uuid.UUID, offset, limit int) ([]*models.Topic, error) {
	var topics []*models.Topic
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Forum").
		Where("forum_id = ?", forumID).
		Where("deleted_at IS NULL").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *TopicRepositoryImpl) GetByTag(ctx context.Context, tag string, offset, limit int) ([]*models.Topic, error) {
	var topics []*models.Topic
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Forum").
		Preload("Tags").
		Joins("JOIN topic_tags ON topics.id = topic_tags.topic_id").
		Joins("JOIN tags ON topic_tags.tag_id = tags.id").
		Where("tags.name = ?", tag).
		Where("topics.deleted_at IS NULL").
		Order("topics.is_pinned DESC, topics.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *TopicRepositoryImpl) GetByTags(ctx context.Context, tags []string, offset, limit int) ([]*models.Topic, error) {
	var topics []*models.Topic

	// For multiple tags, we want topics that have ALL the specified tags (AND logic)
	// We'll use subquery to find topics that have all the specified tags
	subQuery := r.db.Table("topic_tags").
		Select("topic_id").
		Joins("JOIN tags ON topic_tags.tag_id = tags.id").
		Where("tags.name IN ?", tags).
		Group("topic_id").
		Having("COUNT(DISTINCT tags.name) = ?", len(tags))

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Forum").
		Preload("Tags").
		Where("topics.id IN (?)", subQuery).
		Where("topics.deleted_at IS NULL").
		Order("topics.is_pinned DESC, topics.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *TopicRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Topic, error) {
	var topics []*models.Topic
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Forum").
		Where("deleted_at IS NULL").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error
	return topics, err
}

func (r *TopicRepositoryImpl) Update(ctx context.Context, id uuid.UUID, topic *models.Topic) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Updates(topic).Error
}

func (r *TopicRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.Topic{}, id).Error
}

func (r *TopicRepositoryImpl) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *TopicRepositoryImpl) IncrementReplyCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1)).Error
}

func (r *TopicRepositoryImpl) DecrementReplyCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		Where("reply_count > ?", 0).
		UpdateColumn("reply_count", gorm.Expr("reply_count - ?", 1)).Error
}

func (r *TopicRepositoryImpl) Pin(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		Update("is_pinned", true).Error
}

func (r *TopicRepositoryImpl) Unpin(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		Update("is_pinned", false).Error
}

func (r *TopicRepositoryImpl) Lock(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		Update("is_locked", true).Error
}

func (r *TopicRepositoryImpl) Unlock(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", id).
		Update("is_locked", false).Error
}

func (r *TopicRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *TopicRepositoryImpl) CountByForumID(ctx context.Context, forumID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("forum_id = ?", forumID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *TopicRepositoryImpl) CountByTag(ctx context.Context, tag string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Joins("JOIN topic_tags ON topics.id = topic_tags.topic_id").
		Joins("JOIN tags ON topic_tags.tag_id = tags.id").
		Where("tags.name = ?", tag).
		Where("topics.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *TopicRepositoryImpl) CountByTags(ctx context.Context, tags []string) (int64, error) {
	var count int64

	// Count topics that have ALL the specified tags
	subQuery := r.db.Table("topic_tags").
		Select("topic_id").
		Joins("JOIN tags ON topic_tags.tag_id = tags.id").
		Where("tags.name IN ?", tags).
		Group("topic_id").
		Having("COUNT(DISTINCT tags.name) = ?", len(tags))

	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("topics.id IN (?)", subQuery).
		Where("topics.deleted_at IS NULL").
		Count(&count).Error
	return count, err
}

func (r *TopicRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Topic, int64, error) {
	var topics []*models.Topic
	var count int64

	searchQuery := "%" + query + "%"

	dbQuery := r.db.WithContext(ctx).
		Preload("User").
		Preload("Forum").
		Where("deleted_at IS NULL").
		Where("title LIKE ? OR content LIKE ?", searchQuery, searchQuery)

	// Count
	if err := dbQuery.Model(&models.Topic{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get topics
	err := dbQuery.
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&topics).Error

	return topics, count, err
}

func (r *TopicRepositoryImpl) UpdateLikeCount(ctx context.Context, topicID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", topicID).
		Update("like_count", count).Error
}

func (r *TopicRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Count(&count).Error
	return count, err
}

func (r *TopicRepositoryImpl) AssociateTags(ctx context.Context, topicID uuid.UUID, tags []*models.Tag) error {
	if len(tags) == 0 {
		return nil
	}

	// First get the topic to ensure it exists
	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, "id = ?", topicID).Error; err != nil {
		return err
	}

	// Associate tags using GORM's Association
	return r.db.WithContext(ctx).Model(&topic).Association("Tags").Append(tags)
}

func (r *TopicRepositoryImpl) RemoveAllTags(ctx context.Context, topicID uuid.UUID) error {
	// First get the topic to ensure it exists
	var topic models.Topic
	if err := r.db.WithContext(ctx).First(&topic, "id = ?", topicID).Error; err != nil {
		return err
	}

	// Clear all tag associations
	return r.db.WithContext(ctx).Model(&topic).Association("Tags").Clear()
}

func (r *TopicRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(&count).Error
	return count, err
}
