# Task 04: Follow System

## 📋 ภาพรวม
ระบบติดตาม (Follow/Unfollow) ผู้ใช้งาน เพื่อดู Content จากคนที่เราติดตาม

## 🎯 ความสำคัญ
⭐⭐ **Social Feature - การเชื่อมต่อระหว่างผู้ใช้**

## ⏱️ ระยะเวลา
**2 วัน**

## 📦 Dependencies
- ✅ User System (มีอยู่แล้ว)
- ✅ Task 01: Topic & Reply System (optional - สำหรับ feed)
- ✅ Task 02: Video System (optional - สำหรับ feed)

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (domain/models/)

#### 1.1 สร้าง `follow.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index"` // คนที่กดติดตาม
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index"` // คนที่ถูกติดตาม
	CreatedAt   time.Time

	// Relations
	Follower  User `gorm:"foreignKey:FollowerID"`
	Following User `gorm:"foreignKey:FollowingID"`
}

func (Follow) TableName() string {
	return "follows"
}
```

#### 1.2 Update `user.go` model
```go
// Add to User model
type User struct {
	// ... existing fields
	FollowerCount  int `gorm:"default:0" json:"followerCount"`  // จำนวนคนที่ติดตาม
	FollowingCount int `gorm:"default:0" json:"followingCount"` // จำนวนคนที่เราติดตาม
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `follow.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type FollowUserRequest struct {
	FollowingID uuid.UUID `json:"followingId" validate:"required,uuid"`
}

// Response DTOs
type FollowResponse struct {
	ID          uuid.UUID   `json:"id"`
	FollowerID  uuid.UUID   `json:"followerId"`
	FollowingID uuid.UUID   `json:"followingId"`
	CreatedAt   time.Time   `json:"createdAt"`
	Message     string      `json:"message"`
}

type FollowStatusResponse struct {
	IsFollowing bool `json:"isFollowing"`
}

type FollowerResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Avatar         string    `json:"avatar,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	IsFollowing    bool      `json:"isFollowing"` // ว่าเราติดตามคนนี้หรือไม่
	FollowedAt     time.Time `json:"followedAt"`
}

type FollowingResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Avatar         string    `json:"avatar,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	IsFollowing    bool      `json:"isFollowing"` // ว่าเราติดตามคนนี้หรือไม่
	FollowedAt     time.Time `json:"followedAt"`
}

type FollowListResponse struct {
	Users      []FollowerResponse `json:"users"`
	TotalCount int64              `json:"totalCount"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type FollowingListResponse struct {
	Users      []FollowingResponse `json:"users"`
	TotalCount int64               `json:"totalCount"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"totalPages"`
}

type UserStatsResponse struct {
	UserID         uuid.UUID `json:"userId"`
	FollowerCount  int64     `json:"followerCount"`
	FollowingCount int64     `json:"followingCount"`
}
```

---

### 3. Repository Interface (domain/repositories/)

#### 3.1 สร้าง `follow_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

type FollowRepository interface {
	// Follow/Unfollow
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error

	// Check status
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)

	// Get lists
	GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error)

	// Get counts
	GetFollowerCount(ctx context.Context, userID uuid.UUID) (int64, error)
	GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error)

	// Find
	FindByID(ctx context.Context, id uuid.UUID) (*models.Follow, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
```

---

### 4. Repository Implementation (infrastructure/postgres/)

#### 4.1 สร้าง `follow_repository_impl.go`
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

type followRepositoryImpl struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) repositories.FollowRepository {
	return &followRepositoryImpl{db: db}
}

func (r *followRepositoryImpl) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	// ป้องกันการติดตามตัวเอง
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// ตรวจสอบว่าติดตามอยู่แล้วหรือไม่
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("already following")
	}

	// สร้าง follow record
	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *followRepositoryImpl) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not following")
	}
	return nil
}

func (r *followRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *followRepositoryImpl) GetFollowers(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var totalCount int64

	// Get followers (คนที่ติดตาม userID)
	query := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Preload("Follower") // Load follower user data

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
		Find(&follows).Error

	return follows, totalCount, err
}

func (r *followRepositoryImpl) GetFollowing(ctx context.Context, userID uuid.UUID, page, limit int) ([]models.Follow, int64, error) {
	var follows []models.Follow
	var totalCount int64

	// Get following (คนที่ userID ติดตาม)
	query := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Preload("Following") // Load following user data

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
		Find(&follows).Error

	return follows, totalCount, err
}

func (r *followRepositoryImpl) GetFollowerCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *followRepositoryImpl) GetFollowingCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *followRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Follow, error) {
	var follow models.Follow
	err := r.db.WithContext(ctx).
		Preload("Follower").
		Preload("Following").
		First(&follow, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("follow not found")
		}
		return nil, err
	}
	return &follow, nil
}

func (r *followRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Follow{}, id).Error
}
```

---

### 5. Service Interface (domain/services/)

#### 5.1 สร้าง `follow_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
)

type FollowService interface {
	// Follow/Unfollow
	FollowUser(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowResponse, error)
	UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error

	// Check status
	GetFollowStatus(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowStatusResponse, error)

	// Get lists
	GetFollowers(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) (*dto.FollowListResponse, error)
	GetFollowing(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) (*dto.FollowingListResponse, error)

	// Get stats
	GetUserStats(ctx context.Context, userID uuid.UUID) (*dto.UserStatsResponse, error)
}
```

---

### 6. Service Implementation (application/serviceimpl/)

#### 6.1 สร้าง `follow_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"math"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
)

type followServiceImpl struct {
	followRepo repositories.FollowRepository
	userRepo   repositories.UserRepository
}

func NewFollowService(
	followRepo repositories.FollowRepository,
	userRepo repositories.UserRepository,
) services.FollowService {
	return &followServiceImpl{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

func (s *followServiceImpl) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowResponse, error) {
	// Verify both users exist
	follower, err := s.userRepo.FindByID(ctx, followerID)
	if err != nil {
		return nil, errors.New("follower user not found")
	}

	following, err := s.userRepo.FindByID(ctx, followingID)
	if err != nil {
		return nil, errors.New("following user not found")
	}

	// Follow
	if err := s.followRepo.Follow(ctx, followerID, followingID); err != nil {
		return nil, err
	}

	// Update counts asynchronously
	go func() {
		// Update follower's following count
		followingCount, _ := s.followRepo.GetFollowingCount(context.Background(), followerID)
		_ = s.userRepo.UpdateFollowingCount(context.Background(), followerID, int(followingCount))

		// Update following's follower count
		followerCount, _ := s.followRepo.GetFollowerCount(context.Background(), followingID)
		_ = s.userRepo.UpdateFollowerCount(context.Background(), followingID, int(followerCount))
	}()

	return &dto.FollowResponse{
		FollowerID:  followerID,
		FollowingID: followingID,
		Message:     "Successfully followed " + following.Username,
	}, nil
}

func (s *followServiceImpl) UnfollowUser(ctx context.Context, followerID, followingID uuid.UUID) error {
	if err := s.followRepo.Unfollow(ctx, followerID, followingID); err != nil {
		return err
	}

	// Update counts asynchronously
	go func() {
		// Update follower's following count
		followingCount, _ := s.followRepo.GetFollowingCount(context.Background(), followerID)
		_ = s.userRepo.UpdateFollowingCount(context.Background(), followerID, int(followingCount))

		// Update following's follower count
		followerCount, _ := s.followRepo.GetFollowerCount(context.Background(), followingID)
		_ = s.userRepo.UpdateFollowerCount(context.Background(), followingID, int(followerCount))
	}()

	return nil
}

func (s *followServiceImpl) GetFollowStatus(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowStatusResponse, error) {
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return nil, err
	}

	return &dto.FollowStatusResponse{
		IsFollowing: isFollowing,
	}, nil
}

func (s *followServiceImpl) GetFollowers(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) (*dto.FollowListResponse, error) {
	follows, totalCount, err := s.followRepo.GetFollowers(ctx, targetUserID, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	followers := make([]dto.FollowerResponse, len(follows))
	for i, follow := range follows {
		// Check if current user follows this follower
		isFollowing := false
		if currentUserID != uuid.Nil {
			isFollowing, _ = s.followRepo.IsFollowing(ctx, currentUserID, follow.Follower.ID)
		}

		followers[i] = dto.FollowerResponse{
			ID:             follow.Follower.ID,
			Username:       follow.Follower.Username,
			FullName:       follow.Follower.FullName,
			Avatar:         follow.Follower.Avatar,
			Bio:            follow.Follower.Bio,
			FollowerCount:  follow.Follower.FollowerCount,
			FollowingCount: follow.Follower.FollowingCount,
			IsFollowing:    isFollowing,
			FollowedAt:     follow.CreatedAt,
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.FollowListResponse{
		Users:      followers,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *followServiceImpl) GetFollowing(ctx context.Context, currentUserID, targetUserID uuid.UUID, page, limit int) (*dto.FollowingListResponse, error) {
	follows, totalCount, err := s.followRepo.GetFollowing(ctx, targetUserID, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	following := make([]dto.FollowingResponse, len(follows))
	for i, follow := range follows {
		// Check if current user follows this user
		isFollowing := false
		if currentUserID != uuid.Nil {
			isFollowing, _ = s.followRepo.IsFollowing(ctx, currentUserID, follow.Following.ID)
		}

		following[i] = dto.FollowingResponse{
			ID:             follow.Following.ID,
			Username:       follow.Following.Username,
			FullName:       follow.Following.FullName,
			Avatar:         follow.Following.Avatar,
			Bio:            follow.Following.Bio,
			FollowerCount:  follow.Following.FollowerCount,
			FollowingCount: follow.Following.FollowingCount,
			IsFollowing:    isFollowing,
			FollowedAt:     follow.CreatedAt,
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.FollowingListResponse{
		Users:      following,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *followServiceImpl) GetUserStats(ctx context.Context, userID uuid.UUID) (*dto.UserStatsResponse, error) {
	followerCount, err := s.followRepo.GetFollowerCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	followingCount, err := s.followRepo.GetFollowingCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserStatsResponse{
		UserID:         userID,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
	}, nil
}
```

---

### 7. Handlers (interfaces/api/handlers/)

#### 7.1 สร้าง `follow_handler.go`
```go
package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type FollowHandler struct {
	followService services.FollowService
}

func NewFollowHandler(followService services.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// POST /api/v1/users/:userId/follow
func (h *FollowHandler) FollowUser(c *fiber.Ctx) error {
	followerID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	result, err := h.followService.FollowUser(c.Context(), followerID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User followed successfully", result)
}

// DELETE /api/v1/users/:userId/follow
func (h *FollowHandler) UnfollowUser(c *fiber.Ctx) error {
	followerID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.followService.UnfollowUser(c.Context(), followerID, followingID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User unfollowed successfully", nil)
}

// GET /api/v1/users/:userId/follow/status
func (h *FollowHandler) GetFollowStatus(c *fiber.Ctx) error {
	followerID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	status, err := h.followService.GetFollowStatus(c.Context(), followerID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Follow status retrieved", status)
}

// GET /api/v1/users/:userId/followers
func (h *FollowHandler) GetFollowers(c *fiber.Ctx) error {
	targetUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get current user ID (optional - for checking if current user follows each follower)
	currentUserID, _ := utils.GetUserIDFromContext(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	followers, err := h.followService.GetFollowers(c.Context(), currentUserID, targetUserID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Followers retrieved successfully", followers)
}

// GET /api/v1/users/:userId/following
func (h *FollowHandler) GetFollowing(c *fiber.Ctx) error {
	targetUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	// Get current user ID (optional - for checking if current user follows each user)
	currentUserID, _ := utils.GetUserIDFromContext(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	following, err := h.followService.GetFollowing(c.Context(), currentUserID, targetUserID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Following retrieved successfully", following)
}

// GET /api/v1/users/:userId/stats
func (h *FollowHandler) GetUserStats(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	stats, err := h.followService.GetUserStats(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User stats retrieved", stats)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `follow_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupFollowRoutes(app *fiber.App, followHandler *handlers.FollowHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/users/:userId/followers", followHandler.GetFollowers)  // GET /api/v1/users/:userId/followers
	api.Get("/users/:userId/following", followHandler.GetFollowing)  // GET /api/v1/users/:userId/following
	api.Get("/users/:userId/stats", followHandler.GetUserStats)      // GET /api/v1/users/:userId/stats

	// Protected routes
	auth := api.Use(authMiddleware.Protected())
	auth.Post("/users/:userId/follow", followHandler.FollowUser)           // POST /api/v1/users/:userId/follow
	auth.Delete("/users/:userId/follow", followHandler.UnfollowUser)       // DELETE /api/v1/users/:userId/follow
	auth.Get("/users/:userId/follow/status", followHandler.GetFollowStatus) // GET /api/v1/users/:userId/follow/status
}
```

---

### 9. Container Updates (`pkg/di/container.go`)

```go
// Add to container.go

func (c *Container) InitializeFollowComponents() {
	// Repository
	c.FollowRepository = postgres.NewFollowRepository(c.DB)

	// Service
	c.FollowService = serviceimpl.NewFollowService(
		c.FollowRepository,
		c.UserRepository,
	)

	// Handler
	c.FollowHandler = handlers.NewFollowHandler(c.FollowService)
}

// Add to Container struct
type Container struct {
	// ... existing fields

	// Follow
	FollowRepository repositories.FollowRepository
	FollowService    services.FollowService
	FollowHandler    *handlers.FollowHandler
}
```

---

### 10. Main Updates (`cmd/api/main.go`)

```go
// Add to main.go

func main() {
	// ... existing code

	// Initialize follow components
	container.InitializeFollowComponents()

	// Setup routes
	routes.SetupFollowRoutes(app, container.FollowHandler, container.AuthMiddleware)

	// ... rest of the code
}
```

---

### 11. Database Migrations

```sql
-- follows table
CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
CREATE INDEX idx_follows_created_at ON follows(created_at);

-- Update users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS follower_count INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS following_count INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS bio TEXT;

CREATE INDEX idx_users_follower_count ON users(follower_count);
CREATE INDEX idx_users_following_count ON users(following_count);
```

---

### 12. Additional Repository Methods

เพิ่ม methods ใน `UserRepository`:

```go
// user_repository.go interface
UpdateFollowerCount(ctx context.Context, userID uuid.UUID, count int) error
UpdateFollowingCount(ctx context.Context, userID uuid.UUID, count int) error

// user_repository_impl.go
func (r *userRepositoryImpl) UpdateFollowerCount(ctx context.Context, userID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("follower_count", count).Error
}

func (r *userRepositoryImpl) UpdateFollowingCount(ctx context.Context, userID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("following_count", count).Error
}
```

---

## ✅ Checklist

### Models & Database
- [ ] สร้าง `Follow` model
- [ ] Update `User` model (เพิ่ม follower_count, following_count, bio)
- [ ] เพิ่ม migration สำหรับ `follows` table
- [ ] Update `users` table migration
- [ ] สร้าง indexes ที่จำเป็น

### DTOs
- [ ] สร้าง Follow DTOs (Request/Response)
- [ ] สร้าง FollowStatus DTO
- [ ] สร้าง FollowerList และ FollowingList DTOs
- [ ] สร้าง UserStats DTO

### Repository
- [ ] สร้าง `FollowRepository` interface
- [ ] Implement `FollowRepositoryImpl`
- [ ] เพิ่ม methods ใน UserRepository
- [ ] ทดสอบทุก methods

### Service
- [ ] สร้าง `FollowService` interface
- [ ] Implement `FollowServiceImpl`
- [ ] ทดสอบ follow/unfollow logic

### Handlers
- [ ] สร้าง `FollowHandler`
- [ ] Implement ทุก endpoints

### Routes
- [ ] Setup follow routes
- [ ] ใช้ auth middleware

### Integration
- [ ] Register components ใน Container
- [ ] Update main.go
- [ ] ทดสอบ end-to-end

---

## 🧪 Testing Guide

### 1. Follow User
```bash
POST /api/v1/users/{user-id}/follow
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "User followed successfully",
  "data": {
    "followerId": "uuid",
    "followingId": "uuid",
    "message": "Successfully followed john_doe"
  }
}
```

### 2. Unfollow User
```bash
DELETE /api/v1/users/{user-id}/follow
Authorization: Bearer {token}

Response:
{
  "success": true,
  "message": "User unfollowed successfully"
}
```

### 3. Check Follow Status
```bash
GET /api/v1/users/{user-id}/follow/status
Authorization: Bearer {token}

Response:
{
  "success": true,
  "data": {
    "isFollowing": true
  }
}
```

### 4. Get Followers
```bash
GET /api/v1/users/{user-id}/followers?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "john_doe",
        "fullName": "John Doe",
        "avatar": "url",
        "bio": "Hello world",
        "followerCount": 100,
        "followingCount": 50,
        "isFollowing": false,
        "followedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "totalCount": 100,
    "page": 1,
    "limit": 20,
    "totalPages": 5
  }
}
```

### 5. Get Following
```bash
GET /api/v1/users/{user-id}/following?page=1&limit=20

Response:
{
  "success": true,
  "data": {
    "users": [
      {
        "id": "uuid",
        "username": "jane_doe",
        "fullName": "Jane Doe",
        "avatar": "url",
        "bio": "Developer",
        "followerCount": 200,
        "followingCount": 100,
        "isFollowing": true,
        "followedAt": "2024-01-01T00:00:00Z"
      }
    ],
    "totalCount": 50,
    "page": 1,
    "limit": 20,
    "totalPages": 3
  }
}
```

### 6. Get User Stats
```bash
GET /api/v1/users/{user-id}/stats

Response:
{
  "success": true,
  "data": {
    "userId": "uuid",
    "followerCount": 100,
    "followingCount": 50
  }
}
```

---

## 📝 Notes

### Follow System:
- **Follower**: คนที่กดติดตาม (follower_id)
- **Following**: คนที่ถูกติดตาม (following_id)
- ไม่สามารถติดตามตัวเองได้
- ไม่สามารถติดตามซ้ำได้ (Unique constraint)
- Follower/Following counts อัพเดทแบบ asynchronous

### Database:
- ใช้ composite unique index (follower_id, following_id)
- Check constraint ป้องกันการติดตามตัวเอง
- Cascade delete เมื่อลบ user

### Features:
- ✅ Follow/Unfollow users
- ✅ Check follow status
- ✅ ดูรายชื่อ Followers
- ✅ ดูรายชื่อ Following
- ✅ ดู Stats (follower/following counts)
- ✅ แสดงสถานะ "isFollowing" ในแต่ละ user

### Performance:
- Count updates เป็น asynchronous
- ใช้ indexes สำหรับ queries
- Pagination สำหรับ lists
- Cache counts ใน user table

### Use Cases:
- **Social Feed**: แสดง content จากคนที่เราติดตาม
- **Recommendations**: แนะนำ users ที่น่าสนใจ
- **Notifications**: แจ้งเตือนเมื่อมีคนติดตาม
- **Profile**: แสดง follower/following counts

### Future Enhancements:
- [ ] Follow suggestions (คนที่คุณอาจรู้จัก)
- [ ] Mutual followers (เพื่อนร่วมกัน)
- [ ] Block/Mute users
- [ ] Private accounts (ต้องขออนุมัติก่อน follow)
- [ ] Activity feed from following

---

**ระยะเวลาโดยประมาณ:** 2 วัน

**Dependencies:**
- User System (มีอยู่แล้ว)

**Next Task:** Task 05 - Notification System
