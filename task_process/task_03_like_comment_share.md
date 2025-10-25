# Task 03: Like, Comment & Share System

## 📋 ภาพรวม
ระบบการโต้ตอบ (Interaction) สำหรับทั้งกระทู้ (Topic) และวิดีโอ (Video) รวมถึง Like, Comment, และ Share

## ในระบบมี bunny storage อยู่แล้ว
- infrastructure\storage\bunny_storage.go
- infrastructure\postgres\file_repository_impl.go
- interfaces\api\handlers\file_handler.go
- interfaces\api\routes\file_routes.go

## 🎯 ความสำคัญ
⭐⭐⭐ **หลัก - Social Interaction Features**

## ⏱️ ระยะเวลา
**2-3 วัน**

## 📦 Dependencies
- ✅ Task 01: Topic & Reply System
- ✅ Task 02: Video Upload System

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (domain/models/)

#### 1.1 สร้าง `like.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	TopicID   *uuid.UUID `gorm:"type:uuid;index"` // Nullable - for topic likes
	VideoID   *uuid.UUID `gorm:"type:uuid;index"` // Nullable - for video likes
	CreatedAt time.Time

	// Relations
	User  User   `gorm:"foreignKey:UserID"`
	Topic *Topic `gorm:"foreignKey:TopicID"`
	Video *Video `gorm:"foreignKey:VideoID"`
}

func (Like) TableName() string {
	return "likes"
}

// Ensure only one of TopicID or VideoID is set
func (l *Like) Validate() error {
	if l.TopicID == nil && l.VideoID == nil {
		return errors.New("either TopicID or VideoID must be set")
	}
	if l.TopicID != nil && l.VideoID != nil {
		return errors.New("cannot set both TopicID and VideoID")
	}
	return nil
}
```

#### 1.2 สร้าง `comment.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	VideoID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"` // For nested comments
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Soft delete

	// Relations
	User     User       `gorm:"foreignKey:UserID"`
	Video    Video      `gorm:"foreignKey:VideoID"`
	Parent   *Comment   `gorm:"foreignKey:ParentID"`
	Replies  []Comment  `gorm:"foreignKey:ParentID"`
}

func (Comment) TableName() string {
	return "comments"
}
```

#### 1.3 สร้าง `share.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Share struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	VideoID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Platform  string    `gorm:"type:varchar(50)"` // "facebook", "twitter", "line", "copy_link"
	CreatedAt time.Time

	// Relations
	User  User  `gorm:"foreignKey:UserID"`
	Video Video `gorm:"foreignKey:VideoID"`
}

func (Share) TableName() string {
	return "shares"
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `like.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type LikeTopicRequest struct {
	TopicID uuid.UUID `json:"topicId" validate:"required,uuid"`
}

type LikeVideoRequest struct {
	VideoID uuid.UUID `json:"videoId" validate:"required,uuid"`
}

// Response DTOs
type LikeResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	TopicID   *uuid.UUID `json:"topicId,omitempty"`
	VideoID   *uuid.UUID `json:"videoId,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type LikeStatusResponse struct {
	IsLiked   bool  `json:"isLiked"`
	LikeCount int64 `json:"likeCount"`
}
```

#### 2.2 สร้าง `comment.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type CreateCommentRequest struct {
	VideoID  uuid.UUID  `json:"videoId" validate:"required,uuid"`
	Content  string     `json:"content" validate:"required,min=1,max=1000"`
	ParentID *uuid.UUID `json:"parentId" validate:"omitempty,uuid"` // For nested comments
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// Response DTOs
type CommentResponse struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"userId"`
	User      UserSummary      `json:"user"`
	VideoID   uuid.UUID        `json:"videoId"`
	ParentID  *uuid.UUID       `json:"parentId,omitempty"`
	Content   string           `json:"content"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	Replies   []CommentResponse `json:"replies,omitempty"` // Nested replies
}

type CommentListResponse struct {
	Comments   []CommentResponse `json:"comments"`
	TotalCount int64             `json:"totalCount"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"totalPages"`
}

// UserSummary for comment responses (if not already exists)
type UserSummary struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar,omitempty"`
}
```

#### 2.3 สร้าง `share.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type ShareVideoRequest struct {
	VideoID  uuid.UUID `json:"videoId" validate:"required,uuid"`
	Platform string    `json:"platform" validate:"required,oneof=facebook twitter line copy_link"`
}

// Response DTOs
type ShareResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	VideoID   uuid.UUID `json:"videoId"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
	Message   string    `json:"message"`
}

type ShareCountResponse struct {
	VideoID    uuid.UUID `json:"videoId"`
	ShareCount int64     `json:"shareCount"`
}
```

---

### 3. Repository Interfaces (domain/repositories/)

#### 3.1 สร้าง `like_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

type LikeRepository interface {
	// Topic Likes
	LikeTopic(ctx context.Context, userID, topicID uuid.UUID) error
	UnlikeTopic(ctx context.Context, userID, topicID uuid.UUID) error
	IsTopicLikedByUser(ctx context.Context, userID, topicID uuid.UUID) (bool, error)
	GetTopicLikeCount(ctx context.Context, topicID uuid.UUID) (int64, error)

	// Video Likes
	LikeVideo(ctx context.Context, userID, videoID uuid.UUID) error
	UnlikeVideo(ctx context.Context, userID, videoID uuid.UUID) error
	IsVideoLikedByUser(ctx context.Context, userID, videoID uuid.UUID) (bool, error)
	GetVideoLikeCount(ctx context.Context, videoID uuid.UUID) (int64, error)

	// General
	FindByID(ctx context.Context, id uuid.UUID) (*models.Like, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

#### 3.2 สร้าง `comment_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Comment, error)
	FindByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) ([]models.Comment, int64, error)
	FindReplies(ctx context.Context, parentID uuid.UUID) ([]models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error)
}
```

#### 3.3 สร้าง `share_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

type ShareRepository interface {
	Create(ctx context.Context, share *models.Share) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Share, error)
	CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error)
	FindByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) ([]models.Share, int64, error)
}
```

---

### 4. Repository Implementations (infrastructure/postgres/)

#### 4.1 สร้าง `like_repository_impl.go`
```go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type likeRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) repositories.LikeRepository {
	return &likeRepositoryImpl{db: db}
}

// Topic Likes
func (r *likeRepositoryImpl) LikeTopic(ctx context.Context, userID, topicID uuid.UUID) error {
	// Check if already liked
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already liked")
	}

	like := &models.Like{
		UserID:  userID,
		TopicID: &topicID,
	}

	return r.db.WithContext(ctx).Create(like).Error
}

func (r *likeRepositoryImpl) UnlikeTopic(ctx context.Context, userID, topicID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *likeRepositoryImpl) IsTopicLikedByUser(ctx context.Context, userID, topicID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepositoryImpl) GetTopicLikeCount(ctx context.Context, topicID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("topic_id = ?", topicID).
		Count(&count).Error
	return count, err
}

// Video Likes
func (r *likeRepositoryImpl) LikeVideo(ctx context.Context, userID, videoID uuid.UUID) error {
	// Check if already liked
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already liked")
	}

	like := &models.Like{
		UserID:  userID,
		VideoID: &videoID,
	}

	return r.db.WithContext(ctx).Create(like).Error
}

func (r *likeRepositoryImpl) UnlikeVideo(ctx context.Context, userID, videoID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}
	return nil
}

func (r *likeRepositoryImpl) IsVideoLikedByUser(ctx context.Context, userID, videoID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	return count > 0, err
}

func (r *likeRepositoryImpl) GetVideoLikeCount(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Like{}).
		Where("video_id = ?", videoID).
		Count(&count).Error
	return count, err
}

// General
func (r *likeRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Like, error) {
	var like models.Like
	err := r.db.WithContext(ctx).First(&like, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like not found")
		}
		return nil, err
	}
	return &like, nil
}

func (r *likeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Like{}, id).Error
}
```

#### 4.2 สร้าง `comment_repository_impl.go`
```go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repositories.CommentRepository {
	return &commentRepositoryImpl{db: db}
}

func (r *commentRepositoryImpl) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&comment, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comment not found")
		}
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepositoryImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Comment{}).
		Where("video_id = ? AND parent_id IS NULL", videoID). // Only top-level comments
		Preload("User")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error

	return comments, totalCount, err
}

func (r *commentRepositoryImpl) FindReplies(ctx context.Context, parentID uuid.UUID) ([]models.Comment, error) {
	var replies []models.Comment
	err := r.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Preload("User").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

func (r *commentRepositoryImpl) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}

func (r *commentRepositoryImpl) CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Comment{}).
		Where("video_id = ?", videoID).
		Count(&count).Error
	return count, err
}
```

#### 4.3 สร้าง `share_repository_impl.go`
```go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type shareRepositoryImpl struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) repositories.ShareRepository {
	return &shareRepositoryImpl{db: db}
}

func (r *shareRepositoryImpl) Create(ctx context.Context, share *models.Share) error {
	return r.db.WithContext(ctx).Create(share).Error
}

func (r *shareRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Share, error) {
	var share models.Share
	err := r.db.WithContext(ctx).First(&share, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("share not found")
		}
		return nil, err
	}
	return &share, nil
}

func (r *shareRepositoryImpl) CountByVideoID(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Share{}).
		Where("video_id = ?", videoID).
		Count(&count).Error
	return count, err
}

func (r *shareRepositoryImpl) FindByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) ([]models.Share, int64, error) {
	var shares []models.Share
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Share{}).
		Where("video_id = ?", videoID).
		Preload("User")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&shares).Error

	return shares, totalCount, err
}
```

---

### 5. Service Interfaces (domain/services/)

#### 5.1 สร้าง `like_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
)

type LikeService interface {
	// Topic Likes
	LikeTopic(ctx context.Context, userID, topicID uuid.UUID) error
	UnlikeTopic(ctx context.Context, userID, topicID uuid.UUID) error
	GetTopicLikeStatus(ctx context.Context, userID, topicID uuid.UUID) (*dto.LikeStatusResponse, error)

	// Video Likes
	LikeVideo(ctx context.Context, userID, videoID uuid.UUID) error
	UnlikeVideo(ctx context.Context, userID, videoID uuid.UUID) error
	GetVideoLikeStatus(ctx context.Context, userID, videoID uuid.UUID) (*dto.LikeStatusResponse, error)
}
```

#### 5.2 สร้าง `comment_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error)
	GetCommentsByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) (*dto.CommentListResponse, error)
	UpdateComment(ctx context.Context, userID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error)
	DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error
	DeleteCommentByAdmin(ctx context.Context, commentID uuid.UUID) error
}
```

#### 5.3 สร้าง `share_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
)

type ShareService interface {
	ShareVideo(ctx context.Context, userID uuid.UUID, req *dto.ShareVideoRequest) (*dto.ShareResponse, error)
	GetShareCount(ctx context.Context, videoID uuid.UUID) (*dto.ShareCountResponse, error)
}
```

---

### 6. Service Implementations (application/serviceimpl/)

#### 6.1 สร้าง `like_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
)

type likeServiceImpl struct {
	likeRepo  repositories.LikeRepository
	topicRepo repositories.TopicRepository
	videoRepo repositories.VideoRepository
}

func NewLikeService(
	likeRepo repositories.LikeRepository,
	topicRepo repositories.TopicRepository,
	videoRepo repositories.VideoRepository,
) services.LikeService {
	return &likeServiceImpl{
		likeRepo:  likeRepo,
		topicRepo: topicRepo,
		videoRepo: videoRepo,
	}
}

// Topic Likes
func (s *likeServiceImpl) LikeTopic(ctx context.Context, userID, topicID uuid.UUID) error {
	// Verify topic exists
	_, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// Create like
	if err := s.likeRepo.LikeTopic(ctx, userID, topicID); err != nil {
		return err
	}

	// Update topic like count asynchronously
	go func() {
		count, _ := s.likeRepo.GetTopicLikeCount(context.Background(), topicID)
		_ = s.topicRepo.UpdateLikeCount(context.Background(), topicID, int(count))
	}()

	return nil
}

func (s *likeServiceImpl) UnlikeTopic(ctx context.Context, userID, topicID uuid.UUID) error {
	if err := s.likeRepo.UnlikeTopic(ctx, userID, topicID); err != nil {
		return err
	}

	// Update topic like count asynchronously
	go func() {
		count, _ := s.likeRepo.GetTopicLikeCount(context.Background(), topicID)
		_ = s.topicRepo.UpdateLikeCount(context.Background(), topicID, int(count))
	}()

	return nil
}

func (s *likeServiceImpl) GetTopicLikeStatus(ctx context.Context, userID, topicID uuid.UUID) (*dto.LikeStatusResponse, error) {
	isLiked, err := s.likeRepo.IsTopicLikedByUser(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	count, err := s.likeRepo.GetTopicLikeCount(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: count,
	}, nil
}

// Video Likes
func (s *likeServiceImpl) LikeVideo(ctx context.Context, userID, videoID uuid.UUID) error {
	// Verify video exists
	_, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Create like
	if err := s.likeRepo.LikeVideo(ctx, userID, videoID); err != nil {
		return err
	}

	// Update video like count asynchronously
	go func() {
		count, _ := s.likeRepo.GetVideoLikeCount(context.Background(), videoID)
		_ = s.videoRepo.UpdateLikeCount(context.Background(), videoID, int(count))
	}()

	return nil
}

func (s *likeServiceImpl) UnlikeVideo(ctx context.Context, userID, videoID uuid.UUID) error {
	if err := s.likeRepo.UnlikeVideo(ctx, userID, videoID); err != nil {
		return err
	}

	// Update video like count asynchronously
	go func() {
		count, _ := s.likeRepo.GetVideoLikeCount(context.Background(), videoID)
		_ = s.videoRepo.UpdateLikeCount(context.Background(), videoID, int(count))
	}()

	return nil
}

func (s *likeServiceImpl) GetVideoLikeStatus(ctx context.Context, userID, videoID uuid.UUID) (*dto.LikeStatusResponse, error) {
	isLiked, err := s.likeRepo.IsVideoLikedByUser(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	count, err := s.likeRepo.GetVideoLikeCount(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: count,
	}, nil
}
```

#### 6.2 สร้าง `comment_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"math"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
)

type commentServiceImpl struct {
	commentRepo repositories.CommentRepository
	videoRepo   repositories.VideoRepository
	userRepo    repositories.UserRepository
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	videoRepo repositories.VideoRepository,
	userRepo repositories.UserRepository,
) services.CommentService {
	return &commentServiceImpl{
		commentRepo: commentRepo,
		videoRepo:   videoRepo,
		userRepo:    userRepo,
	}
}

func (s *commentServiceImpl) CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// Verify video exists
	_, err := s.videoRepo.FindByID(ctx, req.VideoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	// If parent comment exists, verify it
	if req.ParentID != nil {
		parent, err := s.commentRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		// Ensure parent belongs to same video
		if parent.VideoID != req.VideoID {
			return nil, errors.New("parent comment does not belong to this video")
		}
	}

	comment := &models.Comment{
		UserID:   userID,
		VideoID:  req.VideoID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Update video comment count asynchronously
	go func() {
		count, _ := s.commentRepo.CountByVideoID(context.Background(), req.VideoID)
		_ = s.videoRepo.UpdateCommentCount(context.Background(), req.VideoID, int(count))
	}()

	// Load user for response
	comment.User, _ = s.userRepo.FindByID(ctx, userID)

	return s.toCommentResponse(comment), nil
}

func (s *commentServiceImpl) GetCommentsByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) (*dto.CommentListResponse, error) {
	comments, totalCount, err := s.commentRepo.FindByVideoID(ctx, videoID, page, limit)
	if err != nil {
		return nil, err
	}

	// Load replies for each comment
	commentResponses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResp := s.toCommentResponse(&comment)

		// Load nested replies
		replies, _ := s.commentRepo.FindReplies(ctx, comment.ID)
		replyResps := make([]dto.CommentResponse, len(replies))
		for j, reply := range replies {
			replyResps[j] = *s.toCommentResponse(&reply)
		}
		commentResp.Replies = replyResps

		commentResponses[i] = *commentResp
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.CommentListResponse{
		Comments:   commentResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *commentServiceImpl) UpdateComment(ctx context.Context, userID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if comment.UserID != userID {
		return nil, errors.New("you don't have permission to update this comment")
	}

	comment.Content = req.Content
	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment), nil
}

func (s *commentServiceImpl) DeleteComment(ctx context.Context, userID, commentID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Check ownership
	if comment.UserID != userID {
		return errors.New("you don't have permission to delete this comment")
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	// Update video comment count asynchronously
	go func() {
		count, _ := s.commentRepo.CountByVideoID(context.Background(), comment.VideoID)
		_ = s.videoRepo.UpdateCommentCount(context.Background(), comment.VideoID, int(count))
	}()

	return nil
}

func (s *commentServiceImpl) DeleteCommentByAdmin(ctx context.Context, commentID uuid.UUID) error {
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	// Update video comment count asynchronously
	go func() {
		count, _ := s.commentRepo.CountByVideoID(context.Background(), comment.VideoID)
		_ = s.videoRepo.UpdateCommentCount(context.Background(), comment.VideoID, int(count))
	}()

	return nil
}

func (s *commentServiceImpl) toCommentResponse(comment *models.Comment) *dto.CommentResponse {
	resp := &dto.CommentResponse{
		ID:        comment.ID,
		UserID:    comment.UserID,
		VideoID:   comment.VideoID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	if comment.User != nil {
		resp.User = dto.UserSummary{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			FullName: comment.User.FullName,
			Avatar:   comment.User.Avatar,
		}
	}

	return resp
}
```

#### 6.3 สร้าง `share_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
)

type shareServiceImpl struct {
	shareRepo repositories.ShareRepository
	videoRepo repositories.VideoRepository
}

func NewShareService(
	shareRepo repositories.ShareRepository,
	videoRepo repositories.VideoRepository,
) services.ShareService {
	return &shareServiceImpl{
		shareRepo: shareRepo,
		videoRepo: videoRepo,
	}
}

func (s *shareServiceImpl) ShareVideo(ctx context.Context, userID uuid.UUID, req *dto.ShareVideoRequest) (*dto.ShareResponse, error) {
	// Verify video exists
	video, err := s.videoRepo.FindByID(ctx, req.VideoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	share := &models.Share{
		UserID:   userID,
		VideoID:  req.VideoID,
		Platform: req.Platform,
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, err
	}

	// Update video share count (if you have this field)
	// For now, just track in shares table

	return &dto.ShareResponse{
		ID:        share.ID,
		UserID:    share.UserID,
		VideoID:   share.VideoID,
		Platform:  share.Platform,
		CreatedAt: share.CreatedAt,
		Message:   "Video shared successfully to " + req.Platform,
	}, nil
}

func (s *shareServiceImpl) GetShareCount(ctx context.Context, videoID uuid.UUID) (*dto.ShareCountResponse, error) {
	count, err := s.shareRepo.CountByVideoID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return &dto.ShareCountResponse{
		VideoID:    videoID,
		ShareCount: count,
	}, nil
}
```

---

### 7. Handlers (interfaces/api/handlers/)

#### 7.1 สร้าง `like_handler.go`
```go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type LikeHandler struct {
	likeService services.LikeService
}

func NewLikeHandler(likeService services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// Topic Likes
// POST /api/v1/topics/:id/like
func (h *LikeHandler) LikeTopic(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid topic ID")
	}

	if err := h.likeService.LikeTopic(c.Context(), userID, topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Topic liked successfully", nil)
}

// DELETE /api/v1/topics/:id/like
func (h *LikeHandler) UnlikeTopic(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid topic ID")
	}

	if err := h.likeService.UnlikeTopic(c.Context(), userID, topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Topic unliked successfully", nil)
}

// GET /api/v1/topics/:id/like/status
func (h *LikeHandler) GetTopicLikeStatus(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid topic ID")
	}

	status, err := h.likeService.GetTopicLikeStatus(c.Context(), userID, topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Topic like status retrieved", status)
}

// Video Likes
// POST /api/v1/videos/:id/like
func (h *LikeHandler) LikeVideo(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.likeService.LikeVideo(c.Context(), userID, videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video liked successfully", nil)
}

// DELETE /api/v1/videos/:id/like
func (h *LikeHandler) UnlikeVideo(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.likeService.UnlikeVideo(c.Context(), userID, videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video unliked successfully", nil)
}

// GET /api/v1/videos/:id/like/status
func (h *LikeHandler) GetVideoLikeStatus(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	status, err := h.likeService.GetVideoLikeStatus(c.Context(), userID, videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video like status retrieved", status)
}
```

#### 7.2 สร้าง `comment_handler.go`
```go
package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// POST /api/v1/comments
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	comment, err := h.commentService.CreateComment(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Comment created successfully", comment)
}

// GET /api/v1/videos/:id/comments
func (h *CommentHandler) GetCommentsByVideoID(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	comments, err := h.commentService.GetCommentsByVideoID(c.Context(), videoID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Comments retrieved successfully", comments)
}

// PUT /api/v1/comments/:id
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid comment ID")
	}

	var req dto.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	comment, err := h.commentService.UpdateComment(c.Context(), userID, commentID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Comment updated successfully", comment)
}

// DELETE /api/v1/comments/:id
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid comment ID")
	}

	if err := h.commentService.DeleteComment(c.Context(), userID, commentID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Comment deleted successfully", nil)
}

// Admin: DELETE /api/v1/admin/comments/:id
func (h *CommentHandler) DeleteCommentByAdmin(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid comment ID")
	}

	if err := h.commentService.DeleteCommentByAdmin(c.Context(), commentID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Comment deleted successfully", nil)
}
```

#### 7.3 สร้าง `share_handler.go`
```go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type ShareHandler struct {
	shareService services.ShareService
}

func NewShareHandler(shareService services.ShareService) *ShareHandler {
	return &ShareHandler{shareService: shareService}
}

// POST /api/v1/videos/:id/share
func (h *ShareHandler) ShareVideo(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	var req dto.ShareVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}
	req.VideoID = videoID

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	share, err := h.shareService.ShareVideo(c.Context(), userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Video shared successfully", share)
}

// GET /api/v1/videos/:id/share/count
func (h *ShareHandler) GetShareCount(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	count, err := h.shareService.GetShareCount(c.Context(), videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Share count retrieved", count)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `like_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupLikeRoutes(app *fiber.App, likeHandler *handlers.LikeHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Protected routes
	auth := api.Use(authMiddleware.Protected())

	// Topic Likes
	auth.Post("/topics/:id/like", likeHandler.LikeTopic)                 // POST /api/v1/topics/:id/like
	auth.Delete("/topics/:id/like", likeHandler.UnlikeTopic)              // DELETE /api/v1/topics/:id/like
	auth.Get("/topics/:id/like/status", likeHandler.GetTopicLikeStatus)  // GET /api/v1/topics/:id/like/status

	// Video Likes
	auth.Post("/videos/:id/like", likeHandler.LikeVideo)                 // POST /api/v1/videos/:id/like
	auth.Delete("/videos/:id/like", likeHandler.UnlikeVideo)              // DELETE /api/v1/videos/:id/like
	auth.Get("/videos/:id/like/status", likeHandler.GetVideoLikeStatus)  // GET /api/v1/videos/:id/like/status
}
```

#### 8.2 สร้าง `comment_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupCommentRoutes(app *fiber.App, commentHandler *handlers.CommentHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/videos/:id/comments", commentHandler.GetCommentsByVideoID)  // GET /api/v1/videos/:id/comments

	// Protected routes
	auth := api.Use(authMiddleware.Protected())
	auth.Post("/comments", commentHandler.CreateComment)                  // POST /api/v1/comments
	auth.Put("/comments/:id", commentHandler.UpdateComment)               // PUT /api/v1/comments/:id
	auth.Delete("/comments/:id", commentHandler.DeleteComment)            // DELETE /api/v1/comments/:id

	// Admin routes
	admin := api.Group("/admin", authMiddleware.Protected(), authMiddleware.RequireRole("admin"))
	admin.Delete("/comments/:id", commentHandler.DeleteCommentByAdmin)    // DELETE /api/v1/admin/comments/:id
}
```

#### 8.3 สร้าง `share_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupShareRoutes(app *fiber.App, shareHandler *handlers.ShareHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/videos/:id/share/count", shareHandler.GetShareCount)  // GET /api/v1/videos/:id/share/count

	// Protected routes
	auth := api.Use(authMiddleware.Protected())
	auth.Post("/videos/:id/share", shareHandler.ShareVideo)         // POST /api/v1/videos/:id/share
}
```

---

### 9. Container Updates (`pkg/di/container.go`)

```go
// Add to container.go

func (c *Container) InitializeLikeComponents() {
	// Repository
	c.LikeRepository = postgres.NewLikeRepository(c.DB)

	// Service
	c.LikeService = serviceimpl.NewLikeService(
		c.LikeRepository,
		c.TopicRepository,
		c.VideoRepository,
	)

	// Handler
	c.LikeHandler = handlers.NewLikeHandler(c.LikeService)
}

func (c *Container) InitializeCommentComponents() {
	// Repository
	c.CommentRepository = postgres.NewCommentRepository(c.DB)

	// Service
	c.CommentService = serviceimpl.NewCommentService(
		c.CommentRepository,
		c.VideoRepository,
		c.UserRepository,
	)

	// Handler
	c.CommentHandler = handlers.NewCommentHandler(c.CommentService)
}

func (c *Container) InitializeShareComponents() {
	// Repository
	c.ShareRepository = postgres.NewShareRepository(c.DB)

	// Service
	c.ShareService = serviceimpl.NewShareService(
		c.ShareRepository,
		c.VideoRepository,
	)

	// Handler
	c.ShareHandler = handlers.NewShareHandler(c.ShareService)
}

// Add to Container struct
type Container struct {
	// ... existing fields

	// Like
	LikeRepository repositories.LikeRepository
	LikeService    services.LikeService
	LikeHandler    *handlers.LikeHandler

	// Comment
	CommentRepository repositories.CommentRepository
	CommentService    services.CommentService
	CommentHandler    *handlers.CommentHandler

	// Share
	ShareRepository repositories.ShareRepository
	ShareService    services.ShareService
	ShareHandler    *handlers.ShareHandler
}
```

---

### 10. Main Updates (`cmd/api/main.go`)

```go
// Add to main.go

func main() {
	// ... existing code

	// Initialize components
	container.InitializeLikeComponents()
	container.InitializeCommentComponents()
	container.InitializeShareComponents()

	// Setup routes
	routes.SetupLikeRoutes(app, container.LikeHandler, container.AuthMiddleware)
	routes.SetupCommentRoutes(app, container.CommentHandler, container.AuthMiddleware)
	routes.SetupShareRoutes(app, container.ShareHandler, container.AuthMiddleware)

	// ... rest of the code
}
```

---

### 11. Database Migrations

```sql
-- likes table
CREATE TABLE IF NOT EXISTS likes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE CASCADE,
    video_id UUID REFERENCES videos(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_like_target CHECK (
        (topic_id IS NOT NULL AND video_id IS NULL) OR
        (topic_id IS NULL AND video_id IS NOT NULL)
    ),
    UNIQUE(user_id, topic_id),
    UNIQUE(user_id, video_id)
);

CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_topic_id ON likes(topic_id);
CREATE INDEX idx_likes_video_id ON likes(video_id);

-- comments table
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_video_id ON comments(video_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_deleted_at ON comments(deleted_at);

-- shares table
CREATE TABLE IF NOT EXISTS shares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    platform VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_shares_user_id ON shares(user_id);
CREATE INDEX idx_shares_video_id ON shares(video_id);

-- Update topics table to add like_count
ALTER TABLE topics ADD COLUMN IF NOT EXISTS like_count INT DEFAULT 0;

-- Update videos table if not already added
-- (Already has like_count and comment_count from Task 02)
```

---

### 12. Additional Repository Methods

เพิ่ม methods ใน `TopicRepository` และ `VideoRepository`:

```go
// topic_repository.go interface
UpdateLikeCount(ctx context.Context, topicID uuid.UUID, count int) error

// topic_repository_impl.go
func (r *TopicRepositoryImpl) UpdateLikeCount(ctx context.Context, topicID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Topic{}).
		Where("id = ?", topicID).
		Update("like_count", count).Error
}

// video_repository.go interface
UpdateLikeCount(ctx context.Context, videoID uuid.UUID, count int) error
UpdateCommentCount(ctx context.Context, videoID uuid.UUID, count int) error

// video_repository_impl.go
func (r *videoRepositoryImpl) UpdateLikeCount(ctx context.Context, videoID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", videoID).
		Update("like_count", count).Error
}

func (r *videoRepositoryImpl) UpdateCommentCount(ctx context.Context, videoID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", videoID).
		Update("comment_count", count).Error
}
```

---

## ✅ Checklist

### Models & Database
- [ ] สร้าง `Like` model
- [ ] สร้าง `Comment` model
- [ ] สร้าง `Share` model
- [ ] เพิ่ม migrations สำหรับ tables
- [ ] เพิ่ม `like_count` ใน `topics` table
- [ ] สร้าง indexes ที่จำเป็น

### DTOs
- [ ] สร้าง Like DTOs (Request/Response)
- [ ] สร้าง Comment DTOs (Request/Response)
- [ ] สร้าง Share DTOs (Request/Response)

### Repositories
- [ ] สร้าง `LikeRepository` interface
- [ ] Implement `LikeRepositoryImpl`
- [ ] สร้าง `CommentRepository` interface
- [ ] Implement `CommentRepositoryImpl`
- [ ] สร้าง `ShareRepository` interface
- [ ] Implement `ShareRepositoryImpl`
- [ ] เพิ่ม methods ใน TopicRepository และ VideoRepository

### Services
- [ ] สร้าง `LikeService` interface
- [ ] Implement `LikeServiceImpl`
- [ ] สร้าง `CommentService` interface
- [ ] Implement `CommentServiceImpl`
- [ ] สร้าง `ShareService` interface
- [ ] Implement `ShareServiceImpl`

### Handlers
- [ ] สร้าง `LikeHandler`
- [ ] สร้าง `CommentHandler`
- [ ] สร้าง `ShareHandler`

### Routes
- [ ] Setup like routes
- [ ] Setup comment routes
- [ ] Setup share routes

### Integration
- [ ] Register components ใน Container
- [ ] Update main.go
- [ ] ทดสอบ end-to-end

---

## 🧪 Testing Guide

### 1. Like Topic
```bash
# Like
POST /api/v1/topics/{topic-id}/like
Authorization: Bearer {token}

# Unlike
DELETE /api/v1/topics/{topic-id}/like
Authorization: Bearer {token}

# Get status
GET /api/v1/topics/{topic-id}/like/status
Authorization: Bearer {token}
```

### 2. Like Video
```bash
# Like
POST /api/v1/videos/{video-id}/like
Authorization: Bearer {token}

# Unlike
DELETE /api/v1/videos/{video-id}/like
Authorization: Bearer {token}

# Get status
GET /api/v1/videos/{video-id}/like/status
Authorization: Bearer {token}
```

### 3. Comment on Video
```bash
# Create comment
POST /api/v1/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "videoId": "uuid",
  "content": "Great video!"
}

# Reply to comment (nested)
POST /api/v1/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "videoId": "uuid",
  "content": "I agree!",
  "parentId": "parent-comment-uuid"
}

# Get comments
GET /api/v1/videos/{video-id}/comments?page=1&limit=20

# Update comment
PUT /api/v1/comments/{comment-id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Updated content"
}

# Delete comment
DELETE /api/v1/comments/{comment-id}
Authorization: Bearer {token}
```

### 4. Share Video
```bash
# Share
POST /api/v1/videos/{video-id}/share
Authorization: Bearer {token}
Content-Type: application/json

{
  "platform": "facebook"
}

# Get share count
GET /api/v1/videos/{video-id}/share/count
```

---

## 📝 Notes

### Like System:
- User สามารถไลค์ได้ทั้งกระทู้ (Topic) และวิดีโอ (Video)
- ไม่สามารถไลค์ซ้ำได้ (Unique constraint)
- Like count อัพเดทแบบ asynchronous
- Unlike จะลด count ลง

### Comment System:
- Comment สำหรับวิดีโอเท่านั้น (ไม่ใช่กระทู้)
- รองรับ nested comments (reply to comment)
- เฉพาะเจ้าของเท่านั้นที่แก้ไข/ลบได้
- Admin สามารถลบได้
- Soft delete

### Share System:
- Share สำหรับวิดีโอเท่านั้น
- รองรับหลาย platforms: facebook, twitter, line, copy_link
- Track จำนวนการแชร์
- ไม่จำกัดจำนวนครั้งที่แชร์

### Performance:
- Like/Comment count updates เป็น asynchronous
- ใช้ database indexes สำหรับ performance
- Nested comments load แบบ recursive

### Security:
- ทุก endpoints ต้อง authenticate
- เฉพาะเจ้าของแก้ไข/ลบได้
- Admin มี permission พิเศษ

---

**ระยะเวลาโดยประมาณ:** 2-3 วัน

**Dependencies:**
- Task 01: Topic & Reply System
- Task 02: Video Upload System

**Next Task:** Task 04 - Follow System
