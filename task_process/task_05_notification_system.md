# Task 05: Notification System

## 📋 ภาพรวม
ระบบแจ้งเตือนสำหรับกิจกรรมต่างๆ เช่น ไลค์, คอมเมนต์, การติดตาม, ตอบกระทู้

## 🎯 ความสำคัญ
⭐⭐ **User Engagement - เพิ่ม Engagement ของผู้ใช้**

## ⏱️ ระยะเวลา
**2 วัน**

## 📦 Dependencies
- ✅ Task 01: Topic & Reply System
- ✅ Task 02: Video System
- ✅ Task 03: Like, Comment & Share
- ✅ Task 04: Follow System

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (domain/models/)

#### 1.1 สร้าง `notification.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeTopicReply   NotificationType = "topic_reply"    // มีคนตอบกระทู้
	NotificationTypeTopicLike    NotificationType = "topic_like"     // มีคนไลค์กระทู้
	NotificationTypeVideoLike    NotificationType = "video_like"     // มีคนไลค์วิดีโอ
	NotificationTypeVideoComment NotificationType = "video_comment"  // มีคนคอมเมนต์วิดีโอ
	NotificationTypeCommentReply NotificationType = "comment_reply"  // มีคนตอบกลับคอมเมนต์
	NotificationTypeNewFollower  NotificationType = "new_follower"   // มีคนติดตาม
)

type Notification struct {
	ID         uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null;index"` // ผู้รับการแจ้งเตือน
	ActorID    uuid.UUID        `gorm:"type:uuid;not null"`       // ผู้กระทำ (คนที่ไลค์, คอมเมนต์, ติดตาม)
	Type       NotificationType `gorm:"type:varchar(50);not null;index"`
	ResourceID *uuid.UUID       `gorm:"type:uuid"`                // ID ของ resource (topic_id, video_id, comment_id)
	Message    string           `gorm:"type:text"`
	IsRead     bool             `gorm:"type:boolean;default:false;index"`
	CreatedAt  time.Time

	// Relations
	User  User `gorm:"foreignKey:UserID"`
	Actor User `gorm:"foreignKey:ActorID"`
}

func (Notification) TableName() string {
	return "notifications"
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `notification.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

// Request DTOs
type MarkAsReadRequest struct {
	NotificationIDs []uuid.UUID `json:"notificationIds" validate:"required,min=1"`
}

type NotificationQueryParams struct {
	Type   string `query:"type" validate:"omitempty,oneof=topic_reply topic_like video_like video_comment comment_reply new_follower"`
	IsRead *bool  `query:"isRead"`
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

// Response DTOs
type ActorSummary struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar,omitempty"`
}

type NotificationResponse struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"userId"`
	Actor      ActorSummary           `json:"actor"`
	Type       models.NotificationType `json:"type"`
	ResourceID *uuid.UUID             `json:"resourceId,omitempty"`
	Message    string                 `json:"message"`
	IsRead     bool                   `json:"isRead"`
	CreatedAt  time.Time              `json:"createdAt"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	TotalCount    int64                  `json:"totalCount"`
	UnreadCount   int64                  `json:"unreadCount"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"totalPages"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type MarkAsReadResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"` // จำนวนที่ mark as read
}
```

---

### 3. Repository Interface (domain/repositories/)

#### 3.1 สร้าง `notification_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
)

type NotificationRepository interface {
	// Create
	Create(ctx context.Context, notification *models.Notification) error
	CreateBatch(ctx context.Context, notifications []models.Notification) error

	// Read
	FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) ([]models.Notification, int64, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// Update
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	MarkMultipleAsRead(ctx context.Context, ids []uuid.UUID) error
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
```

---

### 4. Repository Implementation (infrastructure/postgres/)

#### 4.1 สร้าง `notification_repository_impl.go`
```go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type notificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repositories.NotificationRepository {
	return &notificationRepositoryImpl{db: db}
}

func (r *notificationRepositoryImpl) Create(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *notificationRepositoryImpl) CreateBatch(ctx context.Context, notifications []models.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(notifications, 100).Error
}

func (r *notificationRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.WithContext(ctx).
		Preload("Actor").
		First(&notification, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Preload("Actor")

	// Filter by type
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// Filter by read status
	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifications).Error

	return notifications, totalCount, err
}

func (r *notificationRepositoryImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepositoryImpl) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) MarkMultipleAsRead(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id IN ?", ids).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Notification{}, id).Error
}

func (r *notificationRepositoryImpl) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Notification{}).Error
}
```

---

### 5. Service Interface (domain/services/)

#### 5.1 สร้าง `notification_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
)

type NotificationService interface {
	// Create notifications
	CreateTopicReplyNotification(ctx context.Context, topicID, replyUserID uuid.UUID) error
	CreateTopicLikeNotification(ctx context.Context, topicID, likerUserID uuid.UUID) error
	CreateVideoLikeNotification(ctx context.Context, videoID, likerUserID uuid.UUID) error
	CreateVideoCommentNotification(ctx context.Context, videoID, commenterUserID uuid.UUID) error
	CreateCommentReplyNotification(ctx context.Context, commentID, replierUserID uuid.UUID) error
	CreateNewFollowerNotification(ctx context.Context, followedUserID, followerUserID uuid.UUID) error

	// Read notifications
	GetNotifications(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) (*dto.NotificationListResponse, error)
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error)

	// Update notifications
	MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error
	MarkMultipleAsRead(ctx context.Context, userID uuid.UUID, notificationIDs []uuid.UUID) (*dto.MarkAsReadResponse, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) (*dto.MarkAsReadResponse, error)

	// Delete notifications
	DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) error
}
```

---

### 6. Service Implementation (application/serviceimpl/)

#### 6.1 สร้าง `notification_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
)

type notificationServiceImpl struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
	topicRepo        repositories.TopicRepository
	videoRepo        repositories.VideoRepository
	commentRepo      repositories.CommentRepository
}

func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	topicRepo repositories.TopicRepository,
	videoRepo repositories.VideoRepository,
	commentRepo repositories.CommentRepository,
) services.NotificationService {
	return &notificationServiceImpl{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		topicRepo:        topicRepo,
		videoRepo:        videoRepo,
		commentRepo:      commentRepo,
	}
}

// Create notifications
func (s *notificationServiceImpl) CreateTopicReplyNotification(ctx context.Context, topicID, replyUserID uuid.UUID) error {
	// Get topic
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// Don't notify if replying to own topic
	if topic.UserID == replyUserID {
		return nil
	}

	// Get replier user
	replier, err := s.userRepo.FindByID(ctx, replyUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     topic.UserID,
		ActorID:    replyUserID,
		Type:       models.NotificationTypeTopicReply,
		ResourceID: &topicID,
		Message:    fmt.Sprintf("%s ตอบกระทู้ของคุณ: %s", replier.Username, topic.Title),
		IsRead:     false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationServiceImpl) CreateTopicLikeNotification(ctx context.Context, topicID, likerUserID uuid.UUID) error {
	// Get topic
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// Don't notify if liking own topic
	if topic.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     topic.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeTopicLike,
		ResourceID: &topicID,
		Message:    fmt.Sprintf("%s ถูกใจกระทู้ของคุณ: %s", liker.Username, topic.Title),
		IsRead:     false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationServiceImpl) CreateVideoLikeNotification(ctx context.Context, videoID, likerUserID uuid.UUID) error {
	// Get video
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Don't notify if liking own video
	if video.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     video.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeVideoLike,
		ResourceID: &videoID,
		Message:    fmt.Sprintf("%s ถูกใจวิดีโอของคุณ: %s", liker.Username, video.Title),
		IsRead:     false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationServiceImpl) CreateVideoCommentNotification(ctx context.Context, videoID, commenterUserID uuid.UUID) error {
	// Get video
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Don't notify if commenting on own video
	if video.UserID == commenterUserID {
		return nil
	}

	// Get commenter user
	commenter, err := s.userRepo.FindByID(ctx, commenterUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     video.UserID,
		ActorID:    commenterUserID,
		Type:       models.NotificationTypeVideoComment,
		ResourceID: &videoID,
		Message:    fmt.Sprintf("%s แสดงความคิดเห็นในวิดีโอของคุณ: %s", commenter.Username, video.Title),
		IsRead:     false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationServiceImpl) CreateCommentReplyNotification(ctx context.Context, commentID, replierUserID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Don't notify if replying to own comment
	if comment.UserID == replierUserID {
		return nil
	}

	// Get replier user
	replier, err := s.userRepo.FindByID(ctx, replierUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     comment.UserID,
		ActorID:    replierUserID,
		Type:       models.NotificationTypeCommentReply,
		ResourceID: &commentID,
		Message:    fmt.Sprintf("%s ตอบกลับความคิดเห็นของคุณ", replier.Username),
		IsRead:     false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

func (s *notificationServiceImpl) CreateNewFollowerNotification(ctx context.Context, followedUserID, followerUserID uuid.UUID) error {
	// Get follower user
	follower, err := s.userRepo.FindByID(ctx, followerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:  followedUserID,
		ActorID: followerUserID,
		Type:    models.NotificationTypeNewFollower,
		Message: fmt.Sprintf("%s เริ่มติดตามคุณ", follower.Username),
		IsRead:  false,
	}

	return s.notificationRepo.Create(ctx, notification)
}

// Read notifications
func (s *notificationServiceImpl) GetNotifications(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) (*dto.NotificationListResponse, error) {
	notifications, totalCount, err := s.notificationRepo.FindByUserID(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	// Get unread count
	unreadCount, _ := s.notificationRepo.GetUnreadCount(ctx, userID)

	// Convert to response
	notificationResponses := make([]dto.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		notificationResponses[i] = dto.NotificationResponse{
			ID:         notif.ID,
			UserID:     notif.UserID,
			Actor: dto.ActorSummary{
				ID:       notif.Actor.ID,
				Username: notif.Actor.Username,
				FullName: notif.Actor.FullName,
				Avatar:   notif.Actor.Avatar,
			},
			Type:       notif.Type,
			ResourceID: notif.ResourceID,
			Message:    notif.Message,
			IsRead:     notif.IsRead,
			CreatedAt:  notif.CreatedAt,
		}
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.NotificationListResponse{
		Notifications: notificationResponses,
		TotalCount:    totalCount,
		UnreadCount:   unreadCount,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}, nil
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error) {
	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UnreadCountResponse{Count: count}, nil
}

// Update notifications
func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized to mark this notification")
	}

	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *notificationServiceImpl) MarkMultipleAsRead(ctx context.Context, userID uuid.UUID, notificationIDs []uuid.UUID) (*dto.MarkAsReadResponse, error) {
	if len(notificationIDs) == 0 {
		return &dto.MarkAsReadResponse{
			Message: "No notifications to mark as read",
			Count:   0,
		}, nil
	}

	// TODO: Verify all notifications belong to user (optional, depends on security requirements)

	if err := s.notificationRepo.MarkMultipleAsRead(ctx, notificationIDs); err != nil {
		return nil, err
	}

	return &dto.MarkAsReadResponse{
		Message: "Notifications marked as read",
		Count:   len(notificationIDs),
	}, nil
}

func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (*dto.MarkAsReadResponse, error) {
	// Get unread count before marking
	count, _ := s.notificationRepo.GetUnreadCount(ctx, userID)

	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return nil, err
	}

	return &dto.MarkAsReadResponse{
		Message: "All notifications marked as read",
		Count:   int(count),
	}, nil
}

// Delete notifications
func (s *notificationServiceImpl) DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized to delete this notification")
	}

	return s.notificationRepo.Delete(ctx, notificationID)
}
```

---

### 7. Handlers (interfaces/api/handlers/)

#### 7.1 สร้าง `notification_handler.go`
```go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var params dto.NotificationQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	notifications, err := h.notificationService.GetNotifications(c.Context(), userID, &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notifications retrieved successfully", notifications)
}

// GET /api/v1/notifications/unread/count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	count, err := h.notificationService.GetUnreadCount(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Unread count retrieved", count)
}

// PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.notificationService.MarkAsRead(c.Context(), userID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notification marked as read", nil)
}

// PUT /api/v1/notifications/read
func (h *NotificationHandler) MarkMultipleAsRead(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.notificationService.MarkMultipleAsRead(c.Context(), userID, req.NotificationIDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notifications marked as read", result)
}

// PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	result, err := h.notificationService.MarkAllAsRead(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "All notifications marked as read", result)
}

// DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid notification ID")
	}

	if err := h.notificationService.DeleteNotification(c.Context(), userID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Notification deleted successfully", nil)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `notification_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupNotificationRoutes(app *fiber.App, notificationHandler *handlers.NotificationHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Protected routes
	auth := api.Use(authMiddleware.Protected())

	notifications := auth.Group("/notifications")
	notifications.Get("/", notificationHandler.GetNotifications)                    // GET /api/v1/notifications
	notifications.Get("/unread/count", notificationHandler.GetUnreadCount)          // GET /api/v1/notifications/unread/count
	notifications.Put("/:id/read", notificationHandler.MarkAsRead)                  // PUT /api/v1/notifications/:id/read
	notifications.Put("/read", notificationHandler.MarkMultipleAsRead)              // PUT /api/v1/notifications/read
	notifications.Put("/read-all", notificationHandler.MarkAllAsRead)               // PUT /api/v1/notifications/read-all
	notifications.Delete("/:id", notificationHandler.DeleteNotification)            // DELETE /api/v1/notifications/:id
}
```

---

### 9. Integration with Other Services

ต้องเรียก notification service จาก services อื่นๆ:

#### 9.1 Update `reply_service_impl.go`
```go
// เพิ่ม notificationService ใน struct
type replyServiceImpl struct {
	// ... existing fields
	notificationService services.NotificationService
}

// เพิ่มใน CreateReply method
func (s *replyServiceImpl) CreateReply(ctx context.Context, userID uuid.UUID, req *dto.CreateReplyRequest) (*models.Reply, error) {
	// ... existing code

	// Create notification asynchronously
	go func() {
		_ = s.notificationService.CreateTopicReplyNotification(context.Background(), req.TopicID, userID)
	}()

	return reply, nil
}
```

#### 9.2 Update `like_service_impl.go`
```go
// เพิ่ม notificationService ใน struct
type likeServiceImpl struct {
	// ... existing fields
	notificationService services.NotificationService
}

// เพิ่มใน LikeTopic method
func (s *likeServiceImpl) LikeTopic(ctx context.Context, userID, topicID uuid.UUID) error {
	// ... existing code

	// Create notification asynchronously
	go func() {
		_ = s.notificationService.CreateTopicLikeNotification(context.Background(), topicID, userID)
	}()

	return nil
}

// เพิ่มใน LikeVideo method
func (s *likeServiceImpl) LikeVideo(ctx context.Context, userID, videoID uuid.UUID) error {
	// ... existing code

	// Create notification asynchronously
	go func() {
		_ = s.notificationService.CreateVideoLikeNotification(context.Background(), videoID, userID)
	}()

	return nil
}
```

#### 9.3 Update `comment_service_impl.go`
```go
// เพิ่ม notificationService ใน struct
type commentServiceImpl struct {
	// ... existing fields
	notificationService services.NotificationService
}

// เพิ่มใน CreateComment method
func (s *commentServiceImpl) CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// ... existing code

	// Create notification asynchronously
	if req.ParentID != nil {
		// Reply to comment
		go func() {
			_ = s.notificationService.CreateCommentReplyNotification(context.Background(), *req.ParentID, userID)
		}()
	} else {
		// Comment on video
		go func() {
			_ = s.notificationService.CreateVideoCommentNotification(context.Background(), req.VideoID, userID)
		}()
	}

	return s.toCommentResponse(comment), nil
}
```

#### 9.4 Update `follow_service_impl.go`
```go
// เพิ่ม notificationService ใน struct
type followServiceImpl struct {
	// ... existing fields
	notificationService services.NotificationService
}

// เพิ่มใน FollowUser method
func (s *followServiceImpl) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowResponse, error) {
	// ... existing code

	// Create notification asynchronously
	go func() {
		_ = s.notificationService.CreateNewFollowerNotification(context.Background(), followingID, followerID)
	}()

	return &dto.FollowResponse{...}, nil
}
```

---

### 10. Container Updates (`pkg/di/container.go`)

```go
// Add to container.go

func (c *Container) InitializeNotificationComponents() {
	// Repository
	c.NotificationRepository = postgres.NewNotificationRepository(c.DB)

	// Service
	c.NotificationService = serviceimpl.NewNotificationService(
		c.NotificationRepository,
		c.UserRepository,
		c.TopicRepository,
		c.VideoRepository,
		c.CommentRepository,
	)

	// Handler
	c.NotificationHandler = handlers.NewNotificationHandler(c.NotificationService)
}

// Update existing service initializations to include NotificationService
func (c *Container) InitializeLikeComponents() {
	// ...
	c.LikeService = serviceimpl.NewLikeService(
		c.LikeRepository,
		c.TopicRepository,
		c.VideoRepository,
		c.NotificationService, // Add this
	)
	// ...
}

// Similar updates for Comment, Reply, and Follow services

// Add to Container struct
type Container struct {
	// ... existing fields

	// Notification
	NotificationRepository repositories.NotificationRepository
	NotificationService    services.NotificationService
	NotificationHandler    *handlers.NotificationHandler
}
```

---

### 11. Main Updates (`cmd/api/main.go`)

```go
// Add to main.go

func main() {
	// ... existing code

	// Initialize notification components BEFORE other services
	container.InitializeNotificationComponents()

	// Then initialize other services (they depend on notification service)
	container.InitializeLikeComponents()
	container.InitializeCommentComponents()
	container.InitializeReplyComponents()
	container.InitializeFollowComponents()

	// Setup routes
	routes.SetupNotificationRoutes(app, container.NotificationHandler, container.AuthMiddleware)

	// ... rest of the code
}
```

---

### 12. Database Migrations

```sql
-- notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    resource_id UUID,
    message TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_actor_id ON notifications(actor_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;
```

---

## ✅ Checklist

### Models & Database
- [ ] สร้าง `Notification` model
- [ ] กำหนด `NotificationType` constants
- [ ] เพิ่ม migration สำหรับ `notifications` table
- [ ] สร้าง indexes ที่จำเป็น

### DTOs
- [ ] สร้าง Notification DTOs (Request/Response)
- [ ] สร้าง Query params DTO
- [ ] สร้าง Count และ Status DTOs

### Repository
- [ ] สร้าง `NotificationRepository` interface
- [ ] Implement `NotificationRepositoryImpl`
- [ ] ทดสอบทุก methods

### Service
- [ ] สร้าง `NotificationService` interface
- [ ] Implement `NotificationServiceImpl`
- [ ] สร้าง methods สำหรับแต่ละประเภท notification

### Handlers
- [ ] สร้าง `NotificationHandler`
- [ ] Implement ทุก endpoints

### Routes
- [ ] Setup notification routes

### Integration
- [ ] Update Like service
- [ ] Update Comment service
- [ ] Update Reply service
- [ ] Update Follow service
- [ ] Register components ใน Container
- [ ] Update main.go

### Testing
- [ ] ทดสอบ notification creation
- [ ] ทดสอบ read/unread functionality
- [ ] ทดสอบ filters
- [ ] ทดสอบ end-to-end

---

## 🧪 Testing Guide

### 1. Get Notifications
```bash
GET /api/v1/notifications?page=1&limit=20
Authorization: Bearer {token}

# Filter by type
GET /api/v1/notifications?type=video_like

# Filter by read status
GET /api/v1/notifications?isRead=false

Response:
{
  "success": true,
  "data": {
    "notifications": [
      {
        "id": "uuid",
        "userId": "uuid",
        "actor": {
          "id": "uuid",
          "username": "john_doe",
          "fullName": "John Doe",
          "avatar": "url"
        },
        "type": "video_like",
        "resourceId": "video-uuid",
        "message": "john_doe ถูกใจวิดีโอของคุณ: My Video Title",
        "isRead": false,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "totalCount": 50,
    "unreadCount": 10,
    "page": 1,
    "limit": 20,
    "totalPages": 3
  }
}
```

### 2. Get Unread Count
```bash
GET /api/v1/notifications/unread/count
Authorization: Bearer {token}

Response:
{
  "success": true,
  "data": {
    "count": 10
  }
}
```

### 3. Mark Single Notification as Read
```bash
PUT /api/v1/notifications/{notification-id}/read
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "Notification marked as read"
}
```

### 4. Mark Multiple as Read
```bash
PUT /api/v1/notifications/read
Authorization: Bearer {token}
Content-Type: application/json

{
  "notificationIds": ["uuid1", "uuid2", "uuid3"]
}

Response:
{
  "success": true,
  "data": {
    "message": "Notifications marked as read",
    "count": 3
  }
}
```

### 5. Mark All as Read
```bash
PUT /api/v1/notifications/read-all
Authorization: Bearer {token}

Response:
{
  "success": true,
  "data": {
    "message": "All notifications marked as read",
    "count": 10
  }
}
```

### 6. Delete Notification
```bash
DELETE /api/v1/notifications/{notification-id}
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "Notification deleted successfully"
}
```

---

## 📝 Notes

### Notification Types:
1. **topic_reply** - มีคนตอบกระทู้ของคุณ
2. **topic_like** - มีคนไลค์กระทู้ของคุณ
3. **video_like** - มีคนไลค์วิดีโอของคุณ
4. **video_comment** - มีคนคอมเมนต์วิดีโอของคุณ
5. **comment_reply** - มีคนตอบกลับคอมเมนต์ของคุณ
6. **new_follower** - มีคนติดตามคุณ

### Business Logic:
- ไม่แจ้งเตือนเมื่อกระทำกับ content ของตัวเอง
- สร้าง notification แบบ asynchronous (ไม่ block main flow)
- ResourceID เก็บ ID ของ resource ที่เกี่ยวข้อง (topic, video, comment)
- Message เป็น human-readable string

### Performance:
- ใช้ indexes สำหรับ query ที่บ่อย
- Composite index สำหรับ (user_id, is_read)
- Pagination สำหรับ notification list
- Asynchronous notification creation

### Features:
- ✅ Filter by type
- ✅ Filter by read status
- ✅ Mark single/multiple/all as read
- ✅ Delete notifications
- ✅ Unread count badge
- ✅ Pagination

### Future Enhancements:
- [ ] Real-time notifications (WebSocket/SSE)
- [ ] Push notifications (FCM/APNS)
- [ ] Email notifications
- [ ] Notification preferences (เลือกว่าจะรับแจ้งเตือนอะไรบ้าง)
- [ ] Notification grouping (รวม notifications ประเภทเดียวกัน)
- [ ] Notification batching (ส่งสรุปแทนการส่งทีละรายการ)

---

**ระยะเวลาโดยประมาณ:** 2 วัน

**Dependencies:**
- Task 01: Topic & Reply System
- Task 02: Video System
- Task 03: Like, Comment & Share
- Task 04: Follow System

**Next Task:** Task 06 - Admin Dashboard & Management
