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

type TopicServiceImpl struct {
	topicRepo repositories.TopicRepository
	forumRepo repositories.ForumRepository
	replyRepo repositories.ReplyRepository
	tagService services.TagService
}

func NewTopicService(
	topicRepo repositories.TopicRepository,
	forumRepo repositories.ForumRepository,
	replyRepo repositories.ReplyRepository,
	tagService services.TagService,
) services.TopicService {
	return &TopicServiceImpl{
		topicRepo: topicRepo,
		forumRepo: forumRepo,
		replyRepo: replyRepo,
		tagService: tagService,
	}
}

func (s *TopicServiceImpl) CreateTopic(ctx context.Context, userID uuid.UUID, req *dto.CreateTopicRequest) (*models.Topic, error) {
	// Parse ForumID from string to UUID
	forumID, err := uuid.Parse(req.ForumID)
	if err != nil {
		return nil, errors.New("invalid forum ID format")
	}

	// ตรวจสอบว่า forum มีอยู่และเปิดใช้งาน
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}
	if !forum.IsActive {
		return nil, errors.New("forum is not active")
	}

	topic := &models.Topic{
		ID:        uuid.New(),
		ForumID:   forumID,
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		Thumbnail: req.Thumbnail,
		ViewCount:  0,
		ReplyCount: 0,
		IsPinned:   false,
		IsLocked:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.topicRepo.Create(ctx, topic); err != nil {
		return nil, err
	}

	// Handle tags if provided
	if len(req.TagIDs) > 0 {
		tagUUIDs := make([]uuid.UUID, len(req.TagIDs))
		for i, tagID := range req.TagIDs {
			parsedID, err := uuid.Parse(tagID)
			if err != nil {
				continue // Skip invalid tag IDs
			}
			tagUUIDs[i] = parsedID
		}

		// Get valid tags
		tags, err := s.tagService.GetTagsByIDs(ctx, tagUUIDs)
		if err == nil && len(tags) > 0 {
			// Associate tags with topic
			if err := s.topicRepo.AssociateTags(ctx, topic.ID, tags); err != nil {
				// Log error but don't fail the topic creation
				// TODO: Add proper logging
			}
		}
	}

	// เพิ่ม topic count ใน forum
	s.forumRepo.IncrementTopicCount(ctx, forumID)

	return topic, nil
}

func (s *TopicServiceImpl) GetTopic(ctx context.Context, topicID uuid.UUID) (*dto.TopicDetailResponse, error) {
	// Get topic
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, errors.New("topic not found")
	}

	// Increment view count
	s.topicRepo.IncrementViewCount(ctx, topicID)

	// Get replies
	replies, _ := s.replyRepo.GetByTopicID(ctx, topicID, 0, 100)

	// Convert to responses
	topicResp := dto.TopicToTopicResponse(topic)

	replyResps := make([]dto.ReplyResponse, len(replies))
	for i, reply := range replies {
		replyResps[i] = *dto.ReplyToReplyResponse(reply, true)
	}

	return &dto.TopicDetailResponse{
		Topic:   *topicResp,
		Replies: replyResps,
	}, nil
}

func (s *TopicServiceImpl) GetTopics(ctx context.Context, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	topics, err := s.topicRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.topicRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

func (s *TopicServiceImpl) GetTopicsByForum(ctx context.Context, forumID uuid.UUID, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	topics, err := s.topicRepo.GetByForumID(ctx, forumID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.topicRepo.CountByForumID(ctx, forumID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

func (s *TopicServiceImpl) GetTopicsByForumSlug(ctx context.Context, slug string, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	// Get forum by slug first
	forum, err := s.forumRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, 0, errors.New("forum not found")
	}

	// Check if forum is active
	if !forum.IsActive {
		return nil, 0, errors.New("forum is not active")
	}

	// Get topics by forum ID
	topics, err := s.topicRepo.GetByForumID(ctx, forum.ID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.topicRepo.CountByForumID(ctx, forum.ID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

func (s *TopicServiceImpl) GetTopicsByTag(ctx context.Context, tag string, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	topics, err := s.topicRepo.GetByTag(ctx, tag, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.topicRepo.CountByTag(ctx, tag)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

func (s *TopicServiceImpl) GetTopicsByTags(ctx context.Context, tags []string, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	if len(tags) == 0 {
		return nil, 0, errors.New("at least one tag is required")
	}

	topics, err := s.topicRepo.GetByTags(ctx, tags, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.topicRepo.CountByTags(ctx, tags)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

func (s *TopicServiceImpl) UpdateTopic(ctx context.Context, topicID, userID uuid.UUID, req *dto.UpdateTopicRequest) (*models.Topic, error) {
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, errors.New("topic not found")
	}

	// ตรวจสอบว่าเป็นเจ้าของ
	if topic.UserID != userID {
		return nil, errors.New("unauthorized to update this topic")
	}

	// ตรวจสอบว่า topic ถูก lock หรือไม่
	if topic.IsLocked {
		return nil, errors.New("topic is locked")
	}

	// Update fields
	if req.Title != "" {
		topic.Title = req.Title
	}
	if req.Content != "" {
		topic.Content = req.Content
	}
	if req.Thumbnail != "" {
		topic.Thumbnail = req.Thumbnail
	}
	topic.UpdatedAt = time.Now()

	if err := s.topicRepo.Update(ctx, topicID, topic); err != nil {
		return nil, err
	}

	return topic, nil
}

func (s *TopicServiceImpl) DeleteTopic(ctx context.Context, topicID, userID uuid.UUID) error {
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return errors.New("topic not found")
	}

	if topic.UserID != userID {
		return errors.New("unauthorized to delete this topic")
	}

	// ลด topic count ใน forum
	s.forumRepo.DecrementTopicCount(ctx, topic.ForumID)

	return s.topicRepo.Delete(ctx, topicID)
}

func (s *TopicServiceImpl) SearchTopics(ctx context.Context, query string, offset, limit int) ([]*dto.TopicResponse, int64, error) {
	topics, total, err := s.topicRepo.Search(ctx, query, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.TopicResponse, len(topics))
	for i, topic := range topics {
		responses[i] = dto.TopicToTopicResponse(topic)
	}

	return responses, total, nil
}

// Admin Actions
func (s *TopicServiceImpl) PinTopic(ctx context.Context, topicID uuid.UUID) error {
	return s.topicRepo.Pin(ctx, topicID)
}

func (s *TopicServiceImpl) UnpinTopic(ctx context.Context, topicID uuid.UUID) error {
	return s.topicRepo.Unpin(ctx, topicID)
}

func (s *TopicServiceImpl) LockTopic(ctx context.Context, topicID uuid.UUID) error {
	return s.topicRepo.Lock(ctx, topicID)
}

func (s *TopicServiceImpl) UnlockTopic(ctx context.Context, topicID uuid.UUID) error {
	return s.topicRepo.Unlock(ctx, topicID)
}

func (s *TopicServiceImpl) DeleteTopicByAdmin(ctx context.Context, topicID uuid.UUID) error {
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return errors.New("topic not found")
	}

	s.forumRepo.DecrementTopicCount(ctx, topic.ForumID)
	return s.topicRepo.Delete(ctx, topicID)
}
