package serviceimpl

import (
	"context"
	"fmt"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagServiceImpl struct {
	tagRepo repositories.TagRepository
	db      *gorm.DB
}

func NewTagService(tagRepo repositories.TagRepository, db *gorm.DB) services.TagService {
	return &TagServiceImpl{
		tagRepo: tagRepo,
		db:      db,
	}
}

func (s *TagServiceImpl) CreateTag(ctx context.Context, req *dto.CreateTagRequest) (*models.Tag, error) {
	// Check if tag with same name or slug already exists
	existing, err := s.tagRepo.GetBySlug(ctx, req.Slug)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing tag: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("tag with slug '%s' already exists", req.Slug)
	}

	tag := &models.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
		IsActive:    true,
		UsageCount:  0,
	}

	if tag.Color == "" {
		tag.Color = "#3B82F6" // Default blue color
	}

	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (s *TagServiceImpl) GetTags(ctx context.Context, offset, limit int, activeOnly bool) ([]*dto.TagResponse, int, error) {
	tags, total, err := s.tagRepo.GetAll(ctx, offset, limit, activeOnly)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tags: %w", err)
	}

	responses := make([]*dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = dto.TagToTagResponse(tag)
	}

	return responses, total, nil
}

func (s *TagServiceImpl) GetTag(ctx context.Context, tagID uuid.UUID) (*models.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}

func (s *TagServiceImpl) UpdateTag(ctx context.Context, tagID uuid.UUID, req *dto.UpdateTagRequest) (*models.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Description != "" {
		tag.Description = req.Description
	}
	if req.Color != "" {
		tag.Color = req.Color
	}
	if req.IsActive != nil {
		tag.IsActive = *req.IsActive
	}

	if err := s.tagRepo.Update(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return tag, nil
}

func (s *TagServiceImpl) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	if err := s.tagRepo.Delete(ctx, tagID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (s *TagServiceImpl) SearchTags(ctx context.Context, query string, offset, limit int) ([]*dto.TagResponse, int, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return s.GetTags(ctx, offset, limit, true)
	}

	tags, total, err := s.tagRepo.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search tags: %w", err)
	}

	responses := make([]*dto.TagResponse, len(tags))
	for i, tag := range tags {
		responses[i] = dto.TagToTagResponse(tag)
	}

	return responses, total, nil
}

func (s *TagServiceImpl) GetTagsByIDs(ctx context.Context, tagIDs []uuid.UUID) ([]*models.Tag, error) {
	if len(tagIDs) == 0 {
		return []*models.Tag{}, nil
	}

	tags, err := s.tagRepo.GetByIDs(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by IDs: %w", err)
	}

	return tags, nil
}