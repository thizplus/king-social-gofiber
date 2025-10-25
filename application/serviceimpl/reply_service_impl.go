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

type ReplyServiceImpl struct {
	replyRepo           repositories.ReplyRepository
	topicRepo           repositories.TopicRepository
	notificationService services.NotificationService
}

func NewReplyService(
	replyRepo repositories.ReplyRepository,
	topicRepo repositories.TopicRepository,
	notificationService services.NotificationService,
) services.ReplyService {
	return &ReplyServiceImpl{
		replyRepo:           replyRepo,
		topicRepo:           topicRepo,
		notificationService: notificationService,
	}
}

func (s *ReplyServiceImpl) CreateReply(ctx context.Context, topicID, userID uuid.UUID, req *dto.CreateReplyRequest) (*models.Reply, error) {
	// ตรวจสอบว่า topic มีอยู่
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, errors.New("topic not found")
	}

	// ตรวจสอบว่า topic ถูก lock หรือไม่
	if topic.IsLocked {
		return nil, errors.New("topic is locked, cannot reply")
	}

	// ตรวจสอบ parent reply (ถ้ามี)
	var parentID *uuid.UUID
	if req.ParentID != nil {
		parsed, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent ID format")
		}
		parent, err := s.replyRepo.GetByID(ctx, parsed)
		if err != nil {
			return nil, errors.New("parent reply not found")
		}
		if parent.TopicID != topicID {
			return nil, errors.New("parent reply does not belong to this topic")
		}
		parentID = &parsed
	}

	reply := &models.Reply{
		ID:        uuid.New(),
		TopicID:   topicID,
		UserID:    userID,
		ParentID:  parentID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.replyRepo.Create(ctx, reply); err != nil {
		return nil, err
	}

	// Create notification for topic reply
	go func() {
		_ = s.notificationService.CreateTopicReplyNotification(context.Background(), topicID, userID)
	}()

	// เพิ่ม reply count
	s.topicRepo.IncrementReplyCount(ctx, topicID)

	return reply, nil
}

func (s *ReplyServiceImpl) GetReplies(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*dto.ReplyResponse, int64, error) {
	replies, err := s.replyRepo.GetByTopicID(ctx, topicID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.replyRepo.Count(ctx, topicID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*dto.ReplyResponse, len(replies))
	for i, reply := range replies {
		responses[i] = dto.ReplyToReplyResponse(reply, true)
	}

	return responses, total, nil
}

func (s *ReplyServiceImpl) UpdateReply(ctx context.Context, replyID, userID uuid.UUID, req *dto.UpdateReplyRequest) (*models.Reply, error) {
	reply, err := s.replyRepo.GetByID(ctx, replyID)
	if err != nil {
		return nil, errors.New("reply not found")
	}

	if reply.UserID != userID {
		return nil, errors.New("unauthorized to update this reply")
	}

	reply.Content = req.Content
	reply.UpdatedAt = time.Now()

	if err := s.replyRepo.Update(ctx, replyID, reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *ReplyServiceImpl) DeleteReply(ctx context.Context, replyID, userID uuid.UUID) error {
	reply, err := s.replyRepo.GetByID(ctx, replyID)
	if err != nil {
		return errors.New("reply not found")
	}

	if reply.UserID != userID {
		return errors.New("unauthorized to delete this reply")
	}

	// ลด reply count
	s.topicRepo.DecrementReplyCount(ctx, reply.TopicID)

	return s.replyRepo.Delete(ctx, replyID)
}

func (s *ReplyServiceImpl) DeleteReplyByAdmin(ctx context.Context, replyID uuid.UUID) error {
	reply, err := s.replyRepo.GetByID(ctx, replyID)
	if err != nil {
		return errors.New("reply not found")
	}

	s.topicRepo.DecrementReplyCount(ctx, reply.TopicID)
	return s.replyRepo.Delete(ctx, replyID)
}
