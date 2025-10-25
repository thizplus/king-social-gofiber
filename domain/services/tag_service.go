package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type TagService interface {
	CreateTag(ctx context.Context, req *dto.CreateTagRequest) (*models.Tag, error)
	GetTags(ctx context.Context, offset, limit int, activeOnly bool) ([]*dto.TagResponse, int, error)
	GetTag(ctx context.Context, tagID uuid.UUID) (*models.Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, req *dto.UpdateTagRequest) (*models.Tag, error)
	DeleteTag(ctx context.Context, tagID uuid.UUID) error
	SearchTags(ctx context.Context, query string, offset, limit int) ([]*dto.TagResponse, int, error)
	GetTagsByIDs(ctx context.Context, tagIDs []uuid.UUID) ([]*models.Tag, error)
}