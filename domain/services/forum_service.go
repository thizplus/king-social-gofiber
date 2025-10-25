package services

import (
	"context"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type ForumService interface {
	// Admin Actions
	CreateForum(ctx context.Context, adminID uuid.UUID, req *dto.CreateForumRequest) (*models.Forum, error)
	UpdateForum(ctx context.Context, forumID uuid.UUID, req *dto.UpdateForumRequest) (*models.Forum, error)
	DeleteForum(ctx context.Context, forumID uuid.UUID) error
	ReorderForums(ctx context.Context, req *dto.ReorderForumsRequest) error
	GetAllForums(ctx context.Context, includeInactive bool) ([]*dto.ForumResponse, error)
	SyncTopicCount(ctx context.Context, forumID uuid.UUID) error
	SyncAllTopicCounts(ctx context.Context) error

	// Public Actions
	GetActiveForums(ctx context.Context) ([]*dto.ForumResponse, error)
	GetForumByID(ctx context.Context, forumID uuid.UUID) (*dto.ForumResponse, error)
	GetForumBySlug(ctx context.Context, slug string) (*dto.ForumResponse, error)
}
