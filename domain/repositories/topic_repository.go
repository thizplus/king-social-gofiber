package repositories

import (
	"context"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type TopicRepository interface {
	Create(ctx context.Context, topic *models.Topic) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Topic, error)
	GetByForumID(ctx context.Context, forumID uuid.UUID, offset, limit int) ([]*models.Topic, error)
	GetByTag(ctx context.Context, tag string, offset, limit int) ([]*models.Topic, error)
	GetByTags(ctx context.Context, tags []string, offset, limit int) ([]*models.Topic, error)
	List(ctx context.Context, offset, limit int) ([]*models.Topic, error)
	Update(ctx context.Context, id uuid.UUID, topic *models.Topic) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	IncrementReplyCount(ctx context.Context, id uuid.UUID) error
	DecrementReplyCount(ctx context.Context, id uuid.UUID) error
	Pin(ctx context.Context, id uuid.UUID) error
	Unpin(ctx context.Context, id uuid.UUID) error
	Lock(ctx context.Context, id uuid.UUID) error
	Unlock(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	CountByForumID(ctx context.Context, forumID uuid.UUID) (int64, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	CountByTag(ctx context.Context, tag string) (int64, error)
	CountByTags(ctx context.Context, tags []string) (int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Topic, int64, error)
	UpdateLikeCount(ctx context.Context, topicID uuid.UUID, count int) error
	GetTotalCount(ctx context.Context) (int64, error)
	AssociateTags(ctx context.Context, topicID uuid.UUID, tags []*models.Tag) error
	RemoveAllTags(ctx context.Context, topicID uuid.UUID) error
}
