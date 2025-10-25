package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type TopicService interface {
	// User Actions
	CreateTopic(ctx context.Context, userID uuid.UUID, req *dto.CreateTopicRequest) (*models.Topic, error)
	GetTopic(ctx context.Context, topicID uuid.UUID) (*dto.TopicDetailResponse, error)
	GetTopics(ctx context.Context, offset, limit int) ([]*dto.TopicResponse, int64, error)
	GetTopicsByForum(ctx context.Context, forumID uuid.UUID, offset, limit int) ([]*dto.TopicResponse, int64, error)
	GetTopicsByForumSlug(ctx context.Context, slug string, offset, limit int) ([]*dto.TopicResponse, int64, error)
	GetTopicsByTag(ctx context.Context, tag string, offset, limit int) ([]*dto.TopicResponse, int64, error)
	GetTopicsByTags(ctx context.Context, tags []string, offset, limit int) ([]*dto.TopicResponse, int64, error)
	UpdateTopic(ctx context.Context, topicID, userID uuid.UUID, req *dto.UpdateTopicRequest) (*models.Topic, error)
	DeleteTopic(ctx context.Context, topicID, userID uuid.UUID) error
	SearchTopics(ctx context.Context, query string, offset, limit int) ([]*dto.TopicResponse, int64, error)

	// Admin Actions
	PinTopic(ctx context.Context, topicID uuid.UUID) error
	UnpinTopic(ctx context.Context, topicID uuid.UUID) error
	LockTopic(ctx context.Context, topicID uuid.UUID) error
	UnlockTopic(ctx context.Context, topicID uuid.UUID) error
	DeleteTopicByAdmin(ctx context.Context, topicID uuid.UUID) error
}
