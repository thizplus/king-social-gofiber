# Task 01: Topic & Reply System (กระทู้และตอบกระทู้)

## 📋 ภาพรวม
User สร้างกระทู้ในกระดานและตอบกระทู้ (รองรับ nested replies)

## ในระบบมี bunny storage อยู่แล้ว 
- infrastructure\storage\bunny_storage.go
- infrastructure\postgres\file_repository_impl.go
- interfaces\api\handlers\file_handler.go
- interfaces\api\routes\file_routes.go

## 🎯 ความสำคัญ
⭐⭐⭐ **หลัก - ระบบ Webboard หลัก**

## ⏱️ ระยะเวลา
**3-4 วัน**

## 📦 Dependencies
- ✅ Task 00: Admin Forum Management (ต้องมีกระดานก่อน)

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (domain/models/)

#### 1.1 สร้าง `topic.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Topic struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ForumID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Title      string    `gorm:"type:varchar(200);not null"`
	Content    string    `gorm:"type:text;not null"`
	Thumbnail  string    `gorm:"type:varchar(500)"` // Optional thumbnail image URL
	ViewCount  int       `gorm:"default:0"`
	ReplyCount int       `gorm:"default:0"`
	IsPinned   bool      `gorm:"default:false"`
	IsLocked   bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"` // Soft delete

	// Relations
	Forum   Forum   `gorm:"foreignKey:ForumID"`
	User    User    `gorm:"foreignKey:UserID"`
	Replies []Reply `gorm:"foreignKey:TopicID"`
}

func (Topic) TableName() string {
	return "topics"
}
```

#### 1.2 สร้าง `reply.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Reply struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TopicID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"` // สำหรับ nested reply
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Soft delete

	// Relations
	Topic   Topic   `gorm:"foreignKey:TopicID"`
	User    User    `gorm:"foreignKey:UserID"`
	Parent  *Reply  `gorm:"foreignKey:ParentID"`
	Replies []Reply `gorm:"foreignKey:ParentID"` // Nested replies
}

func (Reply) TableName() string {
	return "replies"
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `topic.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type CreateTopicRequest struct {
	ForumID   uuid.UUID `json:"forumId" validate:"required,uuid"`
	Title     string    `json:"title" validate:"required,min=5,max=200"`
	Content   string    `json:"content" validate:"required,min=10,max=10000"`
	Thumbnail string    `json:"thumbnail" validate:"omitempty,url"` // Optional thumbnail image URL
}

type UpdateTopicRequest struct {
	Title     string `json:"title" validate:"omitempty,min=5,max=200"`
	Content   string `json:"content" validate:"omitempty,min=10,max=10000"`
	Thumbnail string `json:"thumbnail" validate:"omitempty,url"` // Optional thumbnail image URL
}

// Response DTOs
type TopicResponse struct {
	ID         uuid.UUID    `json:"id"`
	ForumID    uuid.UUID    `json:"forumId"`
	Forum      ForumResponse `json:"forum"`
	UserID     uuid.UUID    `json:"userId"`
	User       UserResponse `json:"user"`
	Title      string       `json:"title"`
	Content    string       `json:"content"`
	Thumbnail  string       `json:"thumbnail,omitempty"` // Optional thumbnail image URL
	ViewCount  int          `json:"viewCount"`
	ReplyCount int          `json:"replyCount"`
	IsPinned   bool         `json:"isPinned"`
	IsLocked   bool         `json:"isLocked"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}

type TopicListResponse struct {
	Topics []TopicResponse `json:"topics"`
	Meta   PaginationMeta  `json:"meta"`
}

type TopicDetailResponse struct {
	Topic   TopicResponse   `json:"topic"`
	Replies []ReplyResponse `json:"replies"`
}
```

#### 2.2 สร้าง `reply.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type CreateReplyRequest struct {
	Content  string     `json:"content" validate:"required,min=1,max=5000"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"` // สำหรับ nested reply
}

type UpdateReplyRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

// Response DTOs
type ReplyResponse struct {
	ID        uuid.UUID      `json:"id"`
	TopicID   uuid.UUID      `json:"topicId"`
	UserID    uuid.UUID      `json:"userId"`
	User      UserResponse   `json:"user"`
	ParentID  *uuid.UUID     `json:"parentId,omitempty"`
	Content   string         `json:"content"`
	Replies   []ReplyResponse `json:"replies,omitempty"` // Nested replies
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type ReplyListResponse struct {
	Replies []ReplyResponse `json:"replies"`
	Meta    PaginationMeta  `json:"meta"`
}
```

#### 2.3 อัปเดต `mappers.go`
```go
// เพิ่ม functions
func TopicToTopicResponse(topic *models.Topic) *TopicResponse {
	return &TopicResponse{
		ID:         topic.ID,
		ForumID:    topic.ForumID,
		Forum:      *ForumToForumResponse(&topic.Forum),
		UserID:     topic.UserID,
		User:       *UserToUserResponse(&topic.User, false, false),
		Title:      topic.Title,
		Content:    topic.Content,
		ViewCount:  topic.ViewCount,
		ReplyCount: topic.ReplyCount,
		IsPinned:   topic.IsPinned,
		IsLocked:   topic.IsLocked,
		CreatedAt:  topic.CreatedAt,
		UpdatedAt:  topic.UpdatedAt,
	}
}

func ReplyToReplyResponse(reply *models.Reply, includeNested bool) *ReplyResponse {
	resp := &ReplyResponse{
		ID:        reply.ID,
		TopicID:   reply.TopicID,
		UserID:    reply.UserID,
		User:      *UserToUserResponse(&reply.User, false, false),
		ParentID:  reply.ParentID,
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}

	// Include nested replies if requested
	if includeNested && len(reply.Replies) > 0 {
		resp.Replies = make([]ReplyResponse, len(reply.Replies))
		for i, nested := range reply.Replies {
			resp.Replies[i] = *ReplyToReplyResponse(&nested, true)
		}
	}

	return resp
}
```

---

### 3. Repository Interface (domain/repositories/)

#### 3.1 สร้าง `topic_repository.go`
```go
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
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Topic, int64, error)
}
```

#### 3.2 สร้าง `reply_repository.go`
```go
package repositories

import (
	"context"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ReplyRepository interface {
	Create(ctx context.Context, reply *models.Reply) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Reply, error)
	GetByTopicID(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*models.Reply, error)
	GetByParentID(ctx context.Context, parentID uuid.UUID) ([]*models.Reply, error)
	Update(ctx context.Context, id uuid.UUID, reply *models.Reply) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, topicID uuid.UUID) (int64, error)
}
```

---

### 4. Repository Implementation (infrastructure/postgres/)

#### 4.1 สร้าง `topic_repository_impl.go`
```go
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
		Where("id = ?", id).
		Delete(&models.Topic{}).Error
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
```

#### 4.2 สร้าง `reply_repository_impl.go`
```go
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
```

---

### 5. Service Interface (domain/services/)

#### 5.1 สร้าง `topic_service.go`
```go
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
```

#### 5.2 สร้าง `reply_service.go`
```go
package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ReplyService interface {
	CreateReply(ctx context.Context, topicID, userID uuid.UUID, req *dto.CreateReplyRequest) (*models.Reply, error)
	GetReplies(ctx context.Context, topicID uuid.UUID, offset, limit int) ([]*dto.ReplyResponse, int64, error)
	UpdateReply(ctx context.Context, replyID, userID uuid.UUID, req *dto.UpdateReplyRequest) (*models.Reply, error)
	DeleteReply(ctx context.Context, replyID, userID uuid.UUID) error
	DeleteReplyByAdmin(ctx context.Context, replyID uuid.UUID) error
}
```

---

### 6. Service Implementation (application/serviceimpl/)

#### 6.1 สร้าง `topic_service_impl.go`
```go
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
}

func NewTopicService(
	topicRepo repositories.TopicRepository,
	forumRepo repositories.ForumRepository,
	replyRepo repositories.ReplyRepository,
) services.TopicService {
	return &TopicServiceImpl{
		topicRepo: topicRepo,
		forumRepo: forumRepo,
		replyRepo: replyRepo,
	}
}

func (s *TopicServiceImpl) CreateTopic(ctx context.Context, userID uuid.UUID, req *dto.CreateTopicRequest) (*models.Topic, error) {
	// ตรวจสอบว่า forum มีอยู่และเปิดใช้งาน
	forum, err := s.forumRepo.GetByID(ctx, req.ForumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}
	if !forum.IsActive {
		return nil, errors.New("forum is not active")
	}

	topic := &models.Topic{
		ID:        uuid.New(),
		ForumID:   req.ForumID,
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

	// เพิ่ม topic count ใน forum
	s.forumRepo.IncrementTopicCount(ctx, req.ForumID)

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
```

#### 6.2 สร้าง `reply_service_impl.go`
```go
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
	replyRepo repositories.ReplyRepository
	topicRepo repositories.TopicRepository
}

func NewReplyService(
	replyRepo repositories.ReplyRepository,
	topicRepo repositories.TopicRepository,
) services.ReplyService {
	return &ReplyServiceImpl{
		replyRepo: replyRepo,
		topicRepo: topicRepo,
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
	if req.ParentID != nil {
		parent, err := s.replyRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("parent reply not found")
		}
		if parent.TopicID != topicID {
			return nil, errors.New("parent reply does not belong to this topic")
		}
	}

	reply := &models.Reply{
		ID:        uuid.New(),
		TopicID:   topicID,
		UserID:    userID,
		ParentID:  req.ParentID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.replyRepo.Create(ctx, reply); err != nil {
		return nil, err
	}

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
```

---

### 7. Handlers (interfaces/api/handlers/)

#### 7.1 สร้าง `topic_handler.go`
```go
package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TopicHandler struct {
	topicService services.TopicService
}

func NewTopicHandler(topicService services.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

func (h *TopicHandler) CreateTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	topic, err := h.topicService.CreateTopic(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create topic", err)
	}

	return utils.SuccessResponse(c, "Topic created successfully", topic)
}

func (h *TopicHandler) GetTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	topic, err := h.topicService.GetTopic(c.Context(), topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Topic not found", err)
	}

	return utils.SuccessResponse(c, "Topic retrieved successfully", topic)
}

func (h *TopicHandler) GetTopics(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopics(c.Context(), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) GetTopicsByForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopicsByForum(c.Context(), forumID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) UpdateTopic(c *fiber.Ctx) error {
	user, _ := utils.GetUserFromContext(c)
	topicID, _ := uuid.Parse(c.Params("id"))

	var req dto.UpdateTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	topic, err := h.topicService.UpdateTopic(c.Context(), topicID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update topic", err)
	}

	return utils.SuccessResponse(c, "Topic updated successfully", topic)
}

func (h *TopicHandler) DeleteTopic(c *fiber.Ctx) error {
	user, _ := utils.GetUserFromContext(c)
	topicID, _ := uuid.Parse(c.Params("id"))

	if err := h.topicService.DeleteTopic(c.Context(), topicID, user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete topic", err)
	}

	return utils.SuccessResponse(c, "Topic deleted successfully", nil)
}

func (h *TopicHandler) SearchTopics(c *fiber.Ctx) error {
	query := c.Query("q", "")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.SearchTopics(c.Context(), query, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search topics", err)
	}

	return utils.SuccessResponse(c, "Topics found successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

// Admin Handlers
func (h *TopicHandler) PinTopic(c *fiber.Ctx) error {
	topicID, _ := uuid.Parse(c.Params("id"))
	if err := h.topicService.PinTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to pin topic", err)
	}
	return utils.SuccessResponse(c, "Topic pinned successfully", nil)
}

func (h *TopicHandler) UnpinTopic(c *fiber.Ctx) error {
	topicID, _ := uuid.Parse(c.Params("id"))
	if err := h.topicService.UnpinTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unpin topic", err)
	}
	return utils.SuccessResponse(c, "Topic unpinned successfully", nil)
}

func (h *TopicHandler) LockTopic(c *fiber.Ctx) error {
	topicID, _ := uuid.Parse(c.Params("id"))
	if err := h.topicService.LockTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to lock topic", err)
	}
	return utils.SuccessResponse(c, "Topic locked successfully", nil)
}

func (h *TopicHandler) UnlockTopic(c *fiber.Ctx) error {
	topicID, _ := uuid.Parse(c.Params("id"))
	if err := h.topicService.UnlockTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlock topic", err)
	}
	return utils.SuccessResponse(c, "Topic unlocked successfully", nil)
}

func (h *TopicHandler) DeleteTopicByAdmin(c *fiber.Ctx) error {
	topicID, _ := uuid.Parse(c.Params("id"))
	if err := h.topicService.DeleteTopicByAdmin(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete topic", err)
	}
	return utils.SuccessResponse(c, "Topic deleted successfully", nil)
}
```

#### 7.2 สร้าง `reply_handler.go`
```go
package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReplyHandler struct {
	replyService services.ReplyService
}

func NewReplyHandler(replyService services.ReplyService) *ReplyHandler {
	return &ReplyHandler{replyService: replyService}
}

func (h *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	var req dto.CreateReplyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	reply, err := h.replyService.CreateReply(c.Context(), topicID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create reply", err)
	}

	return utils.SuccessResponse(c, "Reply created successfully", reply)
}

func (h *ReplyHandler) GetReplies(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	replies, total, err := h.replyService.GetReplies(c.Context(), topicID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get replies", err)
	}

	return utils.SuccessResponse(c, "Replies retrieved successfully", fiber.Map{
		"replies": replies,
		"meta":    dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *ReplyHandler) UpdateReply(c *fiber.Ctx) error {
	user, _ := utils.GetUserFromContext(c)
	replyID, _ := uuid.Parse(c.Params("id"))

	var req dto.UpdateReplyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	reply, err := h.replyService.UpdateReply(c.Context(), replyID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update reply", err)
	}

	return utils.SuccessResponse(c, "Reply updated successfully", reply)
}

func (h *ReplyHandler) DeleteReply(c *fiber.Ctx) error {
	user, _ := utils.GetUserFromContext(c)
	replyID, _ := uuid.Parse(c.Params("id"))

	if err := h.replyService.DeleteReply(c.Context(), replyID, user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete reply", err)
	}

	return utils.SuccessResponse(c, "Reply deleted successfully", nil)
}

func (h *ReplyHandler) DeleteReplyByAdmin(c *fiber.Ctx) error {
	replyID, _ := uuid.Parse(c.Params("id"))

	if err := h.replyService.DeleteReplyByAdmin(c.Context(), replyID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete reply", err)
	}

	return utils.SuccessResponse(c, "Reply deleted successfully", nil)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `topic_routes.go`
```go
package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupTopicRoutes(api fiber.Router, h *handlers.Handlers) {
	// Public routes
	topics := api.Group("/topics")
	topics.Get("/", h.TopicHandler.GetTopics)
	topics.Get("/search", h.TopicHandler.SearchTopics)
	topics.Get("/:id", h.TopicHandler.GetTopic)

	// Protected routes
	topics.Use(middleware.Protected())
	topics.Post("/", h.TopicHandler.CreateTopic)
	topics.Put("/:id", h.TopicHandler.UpdateTopic)
	topics.Delete("/:id", h.TopicHandler.DeleteTopic)

	// Reply routes
	topics.Post("/:id/replies", h.ReplyHandler.CreateReply)
	topics.Get("/:id/replies", h.ReplyHandler.GetReplies)

	// Forum topics
	api.Get("/forums/:id/topics", h.TopicHandler.GetTopicsByForum)

	// Admin routes
	adminTopics := api.Group("/admin/topics")
	adminTopics.Use(middleware.Protected())
	adminTopics.Use(middleware.AdminOnly())

	adminTopics.Put("/:id/pin", h.TopicHandler.PinTopic)
	adminTopics.Put("/:id/unpin", h.TopicHandler.UnpinTopic)
	adminTopics.Put("/:id/lock", h.TopicHandler.LockTopic)
	adminTopics.Put("/:id/unlock", h.TopicHandler.UnlockTopic)
	adminTopics.Delete("/:id", h.TopicHandler.DeleteTopicByAdmin)
}
```

#### 8.2 สร้าง `reply_routes.go`
```go
package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupReplyRoutes(api fiber.Router, h *handlers.Handlers) {
	replies := api.Group("/replies")
	replies.Use(middleware.Protected())

	replies.Put("/:id", h.ReplyHandler.UpdateReply)
	replies.Delete("/:id", h.ReplyHandler.DeleteReply)

	// Admin routes
	adminReplies := api.Group("/admin/replies")
	adminReplies.Use(middleware.Protected())
	adminReplies.Use(middleware.AdminOnly())

	adminReplies.Delete("/:id", h.ReplyHandler.DeleteReplyByAdmin)
}
```

---

### 9-12. อัปเดตไฟล์อื่นๆ

เหมือน Task 00 ต้องอัปเดต:
- ✅ `container.go` - เพิ่ม TopicRepository, ReplyRepository, TopicService, ReplyService
- ✅ `handlers.go` - เพิ่ม TopicHandler, ReplyHandler
- ✅ `routes.go` - เพิ่ม SetupTopicRoutes, SetupReplyRoutes
- ✅ `database.go` - เพิ่ม &models.Topic{}, &models.Reply{}

---

## ✅ Checklist

- [ ] สร้าง models: `topic.go`, `reply.go`
- [ ] สร้าง DTOs: `topic.go`, `reply.go`
- [ ] อัปเดต `mappers.go`
- [ ] สร้าง repositories (2 interfaces)
- [ ] สร้าง implementations (2 files)
- [ ] สร้าง services (2 interfaces)
- [ ] สร้าง service implementations (2 files)
- [ ] สร้าง handlers (2 files)
- [ ] สร้าง routes (2 files)
- [ ] อัปเดต container
- [ ] อัปเดต handlers struct
- [ ] อัปเดต routes registration
- [ ] อัปเดต migration
- [ ] ทดสอบสร้างกระทู้
- [ ] ทดสอบตอบกระทู้
- [ ] ทดสอบ nested reply

---

## 🧪 Testing Guide

### 1. User สร้างกระทู้
```bash
POST /api/v1/topics
Authorization: Bearer {token}

{
  "forumId": "uuid-ของกระดานเทคโนโลยี",
  "title": "แนะนำโทรศัพท์ราคา 10,000",
  "content": "อยากได้โทรศัพท์ดีๆ ราคาประมาณ 10,000 บาท มีอะไรแนะนำบ้างครับ",
  "thumbnail": "https://example.com/images/phone-thumbnail.jpg"  // Optional
}
```



### 2. User ตอบกระทู้
```bash
POST /api/v1/topics/{topicId}/replies
Authorization: Bearer {token}

{
  "content": "ผมแนะนำ iPhone SE ครับ ราคาพอดี"
}
```

### 3. User ตอบกลับ reply (nested)
```bash
POST /api/v1/topics/{topicId}/replies
Authorization: Bearer {token}

{
  "content": "iPhone SE รุ่นไหนครับ?",
  "parentId": "uuid-ของ-reply-ที่จะตอบกลับ"
}
```

### 4. ดูกระทู้พร้อมคำตอบ
```bash
GET /api/v1/topics/{topicId}
```

### 5. ค้นหากระทู้
```bash
GET /api/v1/topics/search?q=โทรศัพท์&offset=0&limit=20
```

### 6. Admin ปักหมุดกระทู้
```bash
PUT /api/v1/admin/topics/{topicId}/pin
Authorization: Bearer {admin-token}
```

### 7. Admin ล็อคกระทู้
```bash
PUT /api/v1/admin/topics/{topicId}/lock
Authorization: Bearer {admin-token}
```

---

## 📝 Notes

- **View Count** เพิ่มทุกครั้งที่เปิดดูกระทู้
- **Reply Count** เพิ่มเมื่อมีคนตอบ
- **Pinned Topics** แสดงก่อนเสมอ
- **Locked Topics** ไม่สามารถตอบได้
- **Nested Replies** รองรับได้หลายระดับ
- **Soft Delete** ลบแบบไม่ถาวร
- **Thumbnail** ภาพขนาดย่อเป็น optional สามารถใส่หรือไม่ใส่ก็ได้ (ใช้ URL จาก Bunny Storage)

---

**ระยะเวลาโดยประมาณ:** 3-4 วัน

**Next Task:** Task 02 - Video Upload & Management System
