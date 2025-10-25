package serviceimpl

import (
	"context"
	"errors"
	"time"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"

	"github.com/google/uuid"
)

type ForumServiceImpl struct {
	forumRepo repositories.ForumRepository
}

func NewForumService(forumRepo repositories.ForumRepository) services.ForumService {
	return &ForumServiceImpl{
		forumRepo: forumRepo,
	}
}

func (s *ForumServiceImpl) CreateForum(ctx context.Context, adminID uuid.UUID, req *dto.CreateForumRequest) (*models.Forum, error) {
	// ตรวจสอบว่า slug ซ้ำหรือไม่
	existing, _ := s.forumRepo.GetBySlug(ctx, req.Slug)
	if existing != nil {
		return nil, errors.New("forum slug already exists")
	}

	forum := &models.Forum{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		Order:       req.Order,
		IsActive:    true,
		TopicCount:  0,
		CreatedBy:   adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.forumRepo.Create(ctx, forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (s *ForumServiceImpl) UpdateForum(ctx context.Context, forumID uuid.UUID, req *dto.UpdateForumRequest) (*models.Forum, error) {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	// Update only provided fields
	if req.Name != "" {
		forum.Name = req.Name
	}
	if req.Description != "" {
		forum.Description = req.Description
	}
	if req.Icon != "" {
		forum.Icon = req.Icon
	}
	if req.Order >= 0 {
		forum.Order = req.Order
	}
	if req.IsActive != nil {
		forum.IsActive = *req.IsActive
	}

	forum.UpdatedAt = time.Now()

	if err := s.forumRepo.Update(ctx, forumID, forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (s *ForumServiceImpl) DeleteForum(ctx context.Context, forumID uuid.UUID) error {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return errors.New("forum not found")
	}

	// ตรวจสอบว่ามีกระทู้อยู่หรือไม่
	if forum.TopicCount > 0 {
		return errors.New("cannot delete forum with existing topics")
	}

	return s.forumRepo.Delete(ctx, forumID)
}

func (s *ForumServiceImpl) ReorderForums(ctx context.Context, req *dto.ReorderForumsRequest) error {
	// ตรวจสอบและ parse UUID, จากนั้นตรวจสอบว่าทุก forum ID มีอยู่จริง
	forumIDs := make([]uuid.UUID, 0, len(req.ForumOrders))

	for _, item := range req.ForumOrders {
		forumID, err := uuid.Parse(item.ID)
		if err != nil {
			return errors.New("invalid forum ID format")
		}

		// ตรวจสอบว่า forum มีอยู่จริง
		_, err = s.forumRepo.GetByID(ctx, forumID)
		if err != nil {
			return errors.New("forum not found: " + item.ID)
		}

		forumIDs = append(forumIDs, forumID)
	}

	// อัพเดท order
	for i, item := range req.ForumOrders {
		if err := s.forumRepo.UpdateOrder(ctx, forumIDs[i], item.Order); err != nil {
			return err
		}
	}

	return nil
}

func (s *ForumServiceImpl) GetAllForums(ctx context.Context, includeInactive bool) ([]*dto.ForumResponse, error) {
	forums, err := s.forumRepo.GetAll(ctx, includeInactive)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ForumResponse, len(forums))
	for i, forum := range forums {
		responses[i] = dto.ForumToForumResponse(forum)
	}

	return responses, nil
}

func (s *ForumServiceImpl) SyncTopicCount(ctx context.Context, forumID uuid.UUID) error {
	// Check if forum exists
	_, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return errors.New("forum not found")
	}

	// Sync the topic count
	return s.forumRepo.SyncTopicCount(ctx, forumID)
}

func (s *ForumServiceImpl) SyncAllTopicCounts(ctx context.Context) error {
	// Get all forums
	forums, err := s.forumRepo.GetAll(ctx, true) // Include inactive forums
	if err != nil {
		return err
	}

	// Sync topic count for each forum
	for _, forum := range forums {
		if err := s.forumRepo.SyncTopicCount(ctx, forum.ID); err != nil {
			// Log error but continue with other forums
			// In production, you might want proper logging here
			continue
		}
	}

	return nil
}

func (s *ForumServiceImpl) GetActiveForums(ctx context.Context) ([]*dto.ForumResponse, error) {
	// Sync all topic counts first to ensure accurate data
	if err := s.SyncAllTopicCounts(ctx); err != nil {
		// Log error but continue - don't fail the entire request
		// In production, you might want proper logging here
	}

	return s.GetAllForums(ctx, false)
}

func (s *ForumServiceImpl) GetForumByID(ctx context.Context, forumID uuid.UUID) (*dto.ForumResponse, error) {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	return dto.ForumToForumResponse(forum), nil
}

func (s *ForumServiceImpl) GetForumBySlug(ctx context.Context, slug string) (*dto.ForumResponse, error) {
	forum, err := s.forumRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	return dto.ForumToForumResponse(forum), nil
}
