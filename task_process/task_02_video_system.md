# Task 02: Video Upload & Management System

## 📋 ภาพรวม
ระบบอัปโหลดและจัดการวิดีโอแบบ TikTok-style ที่ผู้ใช้สามารถอัปโหลดวิดีโอ, ดู Feed วิดีโอ, และจัดการวิดีโอของตัวเอง

**Phase:** 2 (Video System)
**ระยะเวลา:** 3-4 วัน
**Dependencies:** Task 00 (ต้องมีระบบ User)

---

## 🎯 Features

### User Features:
- ✅ อัปโหลดวิดีโอ (รองรับ multipart upload)
- ✅ สร้าง Thumbnail อัตโนมัติ
- ✅ ดูข้อมูล Video Metadata (duration, resolution)
- ✅ ดู Video Feed/List (แบบ infinite scroll)
- ✅ ดูวิดีโอของผู้ใช้คนอื่น
- ✅ แก้ไขข้อมูลวิดีโอ (เฉพาะเจ้าของ)
- ✅ ลบวิดีโอ (เฉพาะเจ้าของ)
- ✅ นับจำนวน View อัตโนมัติ

### Admin Features:
- ✅ ดูวิดีโอทั้งหมด
- ✅ ลบวิดีโอที่ไม่เหมาะสม
- ✅ ซ่อน/แสดงวิดีโอ

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (`domain/models/video.go`)

```go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Title        string         `gorm:"type:varchar(200);not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	VideoURL     string         `gorm:"type:varchar(500);not null" json:"videoUrl"`
	ThumbnailURL string         `gorm:"type:varchar(500)" json:"thumbnailUrl"`
	Duration     int            `gorm:"type:int;default:0" json:"duration"` // Duration in seconds
	Width        int            `gorm:"type:int" json:"width"`
	Height       int            `gorm:"type:int" json:"height"`
	FileSize     int64          `gorm:"type:bigint" json:"fileSize"` // File size in bytes
	ViewCount    int            `gorm:"type:int;default:0" json:"viewCount"`
	LikeCount    int            `gorm:"type:int;default:0" json:"likeCount"`
	CommentCount int            `gorm:"type:int;default:0" json:"commentCount"`
	IsActive     bool           `gorm:"type:boolean;default:true" json:"isActive"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

// BeforeCreate hook
func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
```

**Migration:**
```sql
-- Add to migrations
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    video_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    duration INT DEFAULT 0,
    width INT,
    height INT,
    file_size BIGINT,
    view_count INT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_video_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_created_at ON videos(created_at DESC);
CREATE INDEX idx_videos_deleted_at ON videos(deleted_at);
CREATE INDEX idx_videos_is_active ON videos(is_active);
```

---

### 2. DTOs (`domain/dto/video.go`)

```go
package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// ============= Request DTOs =============

type UploadVideoRequest struct {
	Title       string                `form:"title" validate:"required,min=3,max=200"`
	Description string                `form:"description" validate:"omitempty,max=1000"`
	Video       *multipart.FileHeader `form:"video" validate:"required"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail" validate:"omitempty"`
}

type UpdateVideoRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=200"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	IsActive    *bool  `json:"isActive"`
}

type VideoQueryParams struct {
	Page     int       `query:"page" validate:"omitempty,min=1"`
	Limit    int       `query:"limit" validate:"omitempty,min=1,max=100"`
	UserID   uuid.UUID `query:"userId" validate:"omitempty,uuid"`
	SortBy   string    `query:"sortBy" validate:"omitempty,oneof=newest oldest popular"`
	IsActive *bool     `query:"isActive"`
}

// ============= Response DTOs =============

type VideoResponse struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"userId"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	VideoURL     string       `json:"videoUrl"`
	ThumbnailURL string       `json:"thumbnailUrl"`
	Duration     int          `json:"duration"`
	Width        int          `json:"width"`
	Height       int          `json:"height"`
	FileSize     int64        `json:"fileSize"`
	ViewCount    int          `json:"viewCount"`
	LikeCount    int          `json:"likeCount"`
	CommentCount int          `json:"commentCount"`
	IsActive     bool         `json:"isActive"`
	CreatedAt    time.Time    `json:"createdAt"`
	User         *UserSummary `json:"user,omitempty"`
}

type VideoListResponse struct {
	Videos     []VideoResponse `json:"videos"`
	TotalCount int64           `json:"totalCount"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"totalPages"`
}

type UploadVideoResponse struct {
	Video   *VideoResponse `json:"video"`
	Message string         `json:"message"`
}

// UserSummary for video responses
type UserSummary struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"fullName"`
	Avatar   string    `json:"avatar"`
}
```

---

### 3. Repository Interface (`domain/repositories/video_repository.go`)

```go
package repositories

import (
	"context"

	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
)

type VideoRepository interface {
	// Basic CRUD
	Create(ctx context.Context, video *models.Video) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error)
	Update(ctx context.Context, video *models.Video) error
	Delete(ctx context.Context, id uuid.UUID) error

	// List & Query
	FindAll(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) ([]models.Video, int64, error)

	// View Count
	IncrementViewCount(ctx context.Context, id uuid.UUID) error

	// Admin
	FindAllIncludingInactive(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error)
	SetActive(ctx context.Context, id uuid.UUID, isActive bool) error
}
```

---

### 4. Repository Implementation (`infrastructure/postgres/video_repository_impl.go`)

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

type videoRepositoryImpl struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) repositories.VideoRepository {
	return &videoRepositoryImpl{db: db}
}

func (r *videoRepositoryImpl) Create(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Create(video).Error
}

func (r *videoRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND is_active = ?", id, true).
		First(&video).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("video not found")
		}
		return nil, err
	}
	return &video, nil
}

func (r *videoRepositoryImpl) Update(ctx context.Context, video *models.Video) error {
	return r.db.WithContext(ctx).Save(video).Error
}

func (r *videoRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Video{}, id).Error
}

func (r *videoRepositoryImpl) FindAll(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).
		Preload("User").
		Where("is_active = ?", true)

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	orderBy := "created_at DESC" // default: newest
	if params.SortBy == "oldest" {
		orderBy = "created_at ASC"
	} else if params.SortBy == "popular" {
		orderBy = "view_count DESC"
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

	err := query.Order(orderBy).Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).
		Preload("User").
		Where("user_id = ? AND is_active = ?", userID, true)

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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *videoRepositoryImpl) FindAllIncludingInactive(ctx context.Context, params *dto.VideoQueryParams) ([]models.Video, int64, error) {
	var videos []models.Video
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Video{}).Preload("User")

	// Filter by IsActive if provided
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
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

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&videos).Error
	return videos, totalCount, err
}

func (r *videoRepositoryImpl) SetActive(ctx context.Context, id uuid.UUID, isActive bool) error {
	return r.db.WithContext(ctx).
		Model(&models.Video{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}
```

---

### 5. Service Interface (`domain/services/video_service.go`)

```go
package services

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"yourproject/domain/dto"
)

type VideoService interface {
	// User operations
	UploadVideo(ctx context.Context, userID uuid.UUID, req *dto.UploadVideoRequest) (*dto.VideoResponse, error)
	GetVideoByID(ctx context.Context, id uuid.UUID) (*dto.VideoResponse, error)
	GetVideos(ctx context.Context, params *dto.VideoQueryParams) (*dto.VideoListResponse, error)
	GetUserVideos(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) (*dto.VideoListResponse, error)
	UpdateVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, req *dto.UpdateVideoRequest) (*dto.VideoResponse, error)
	DeleteVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) error

	// Admin operations
	GetAllVideos(ctx context.Context, params *dto.VideoQueryParams) (*dto.VideoListResponse, error)
	HideVideo(ctx context.Context, videoID uuid.UUID) error
	ShowVideo(ctx context.Context, videoID uuid.UUID) error
	DeleteVideoByAdmin(ctx context.Context, videoID uuid.UUID) error
}
```

---

### 6. Service Implementation (`application/serviceimpl/video_service_impl.go`)

```go
package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type videoServiceImpl struct {
	videoRepo repositories.VideoRepository
	userRepo  repositories.UserRepository
	storage   utils.StorageService
}

func NewVideoService(
	videoRepo repositories.VideoRepository,
	userRepo repositories.UserRepository,
	storage utils.StorageService,
) services.VideoService {
	return &videoServiceImpl{
		videoRepo: videoRepo,
		userRepo:  userRepo,
		storage:   storage,
	}
}

func (s *videoServiceImpl) UploadVideo(ctx context.Context, userID uuid.UUID, req *dto.UploadVideoRequest) (*dto.VideoResponse, error) {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Validate video file
	if req.Video == nil {
		return nil, errors.New("video file is required")
	}

	// Check file extension
	ext := filepath.Ext(req.Video.Filename)
	allowedExts := map[string]bool{".mp4": true, ".mov": true, ".avi": true, ".mkv": true}
	if !allowedExts[ext] {
		return nil, errors.New("invalid video format. Allowed: mp4, mov, avi, mkv")
	}

	// Check file size (max 500MB)
	maxSize := int64(500 * 1024 * 1024)
	if req.Video.Size > maxSize {
		return nil, errors.New("video file too large. Max 500MB")
	}

	// Upload video to storage
	videoURL, err := s.storage.UploadFile(ctx, req.Video, "videos")
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %w", err)
	}

	// Upload thumbnail if provided
	var thumbnailURL string
	if req.Thumbnail != nil {
		thumbnailURL, err = s.storage.UploadFile(ctx, req.Thumbnail, "thumbnails")
		if err != nil {
			// Log error but don't fail
			thumbnailURL = ""
		}
	}

	// TODO: Extract video metadata (duration, width, height) using ffmpeg
	// For now, set defaults
	duration := 0
	width := 0
	height := 0

	// Create video record
	video := &models.Video{
		UserID:       userID,
		Title:        req.Title,
		Description:  req.Description,
		VideoURL:     videoURL,
		ThumbnailURL: thumbnailURL,
		Duration:     duration,
		Width:        width,
		Height:       height,
		FileSize:     req.Video.Size,
		IsActive:     true,
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		// Cleanup uploaded files
		_ = s.storage.DeleteFile(ctx, videoURL)
		if thumbnailURL != "" {
			_ = s.storage.DeleteFile(ctx, thumbnailURL)
		}
		return nil, fmt.Errorf("failed to create video: %w", err)
	}

	// Load user for response
	video.User, _ = s.userRepo.FindByID(ctx, userID)

	return s.toVideoResponse(video), nil
}

func (s *videoServiceImpl) GetVideoByID(ctx context.Context, id uuid.UUID) (*dto.VideoResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously
	go func() {
		_ = s.videoRepo.IncrementViewCount(context.Background(), id)
	}()

	return s.toVideoResponse(video), nil
}

func (s *videoServiceImpl) GetVideos(ctx context.Context, params *dto.VideoQueryParams) (*dto.VideoListResponse, error) {
	videos, totalCount, err := s.videoRepo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toVideoListResponse(videos, totalCount, params), nil
}

func (s *videoServiceImpl) GetUserVideos(ctx context.Context, userID uuid.UUID, params *dto.VideoQueryParams) (*dto.VideoListResponse, error) {
	videos, totalCount, err := s.videoRepo.FindByUserID(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	return s.toVideoListResponse(videos, totalCount, params), nil
}

func (s *videoServiceImpl) UpdateVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID, req *dto.UpdateVideoRequest) (*dto.VideoResponse, error) {
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if video.UserID != userID {
		return nil, errors.New("you don't have permission to update this video")
	}

	// Update fields
	if req.Title != "" {
		video.Title = req.Title
	}
	if req.Description != "" {
		video.Description = req.Description
	}
	if req.IsActive != nil {
		video.IsActive = *req.IsActive
	}

	if err := s.videoRepo.Update(ctx, video); err != nil {
		return nil, err
	}

	return s.toVideoResponse(video), nil
}

func (s *videoServiceImpl) DeleteVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) error {
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Check ownership
	if video.UserID != userID {
		return errors.New("you don't have permission to delete this video")
	}

	// Delete from storage
	_ = s.storage.DeleteFile(ctx, video.VideoURL)
	if video.ThumbnailURL != "" {
		_ = s.storage.DeleteFile(ctx, video.ThumbnailURL)
	}

	return s.videoRepo.Delete(ctx, videoID)
}

// Admin operations
func (s *videoServiceImpl) GetAllVideos(ctx context.Context, params *dto.VideoQueryParams) (*dto.VideoListResponse, error) {
	videos, totalCount, err := s.videoRepo.FindAllIncludingInactive(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.toVideoListResponse(videos, totalCount, params), nil
}

func (s *videoServiceImpl) HideVideo(ctx context.Context, videoID uuid.UUID) error {
	return s.videoRepo.SetActive(ctx, videoID, false)
}

func (s *videoServiceImpl) ShowVideo(ctx context.Context, videoID uuid.UUID) error {
	return s.videoRepo.SetActive(ctx, videoID, true)
}

func (s *videoServiceImpl) DeleteVideoByAdmin(ctx context.Context, videoID uuid.UUID) error {
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Delete from storage
	_ = s.storage.DeleteFile(ctx, video.VideoURL)
	if video.ThumbnailURL != "" {
		_ = s.storage.DeleteFile(ctx, video.ThumbnailURL)
	}

	return s.videoRepo.Delete(ctx, videoID)
}

// Helper methods
func (s *videoServiceImpl) toVideoResponse(video *models.Video) *dto.VideoResponse {
	resp := &dto.VideoResponse{
		ID:           video.ID,
		UserID:       video.UserID,
		Title:        video.Title,
		Description:  video.Description,
		VideoURL:     video.VideoURL,
		ThumbnailURL: video.ThumbnailURL,
		Duration:     video.Duration,
		Width:        video.Width,
		Height:       video.Height,
		FileSize:     video.FileSize,
		ViewCount:    video.ViewCount,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
		IsActive:     video.IsActive,
		CreatedAt:    video.CreatedAt,
	}

	if video.User != nil {
		resp.User = &dto.UserSummary{
			ID:       video.User.ID,
			Username: video.User.Username,
			FullName: video.User.FullName,
			Avatar:   video.User.Avatar,
		}
	}

	return resp
}

func (s *videoServiceImpl) toVideoListResponse(videos []models.Video, totalCount int64, params *dto.VideoQueryParams) *dto.VideoListResponse {
	videoResponses := make([]dto.VideoResponse, len(videos))
	for i, video := range videos {
		videoResponses[i] = *s.toVideoResponse(&video)
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

	return &dto.VideoListResponse{
		Videos:     videoResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}
```

---

### 7. Storage Service (`pkg/utils/storage.go`)

```go
package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type StorageService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}

type localStorageService struct {
	uploadDir string
	baseURL   string
}

func NewLocalStorageService(uploadDir, baseURL string) StorageService {
	return &localStorageService{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *localStorageService) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// Create folder if not exists
	folderPath := filepath.Join(s.uploadDir, folder)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", err
	}

	// Full file path
	filePath := filepath.Join(folderPath, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	// Return URL
	fileURL := fmt.Sprintf("%s/%s/%s", s.baseURL, folder, filename)
	return fileURL, nil
}

func (s *localStorageService) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract file path from URL
	// For production, implement proper URL parsing
	return nil
}
```

---

### 8. Handlers (`interfaces/api/handlers/video_handler.go`)

```go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
)

type VideoHandler struct {
	videoService services.VideoService
}

func NewVideoHandler(videoService services.VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

// UploadVideo handles video upload
// POST /api/v1/videos
func (h *VideoHandler) UploadVideo(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid form data")
	}

	// Get title and description
	title := c.FormValue("title")
	description := c.FormValue("description")

	// Get video file
	videoFiles := form.File["video"]
	if len(videoFiles) == 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video file is required")
	}

	// Get thumbnail (optional)
	var thumbnail *multipart.FileHeader
	thumbnailFiles := form.File["thumbnail"]
	if len(thumbnailFiles) > 0 {
		thumbnail = thumbnailFiles[0]
	}

	req := &dto.UploadVideoRequest{
		Title:       title,
		Description: description,
		Video:       videoFiles[0],
		Thumbnail:   thumbnail,
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Upload video
	video, err := h.videoService.UploadVideo(c.Context(), userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Video uploaded successfully", video)
}

// GetVideos handles listing videos
// GET /api/v1/videos
func (h *VideoHandler) GetVideos(c *fiber.Ctx) error {
	var params dto.VideoQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	videos, err := h.videoService.GetVideos(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Videos retrieved successfully", videos)
}

// GetVideoByID handles getting a single video
// GET /api/v1/videos/:id
func (h *VideoHandler) GetVideoByID(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	video, err := h.videoService.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video retrieved successfully", video)
}

// GetUserVideos handles getting videos by user
// GET /api/v1/videos/user/:userId
func (h *VideoHandler) GetUserVideos(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var params dto.VideoQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	videos, err := h.videoService.GetUserVideos(c.Context(), userID, &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User videos retrieved successfully", videos)
}

// UpdateVideo handles updating video info
// PUT /api/v1/videos/:id
func (h *VideoHandler) UpdateVideo(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	var req dto.UpdateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	video, err := h.videoService.UpdateVideo(c.Context(), userID, videoID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video updated successfully", video)
}

// DeleteVideo handles deleting a video
// DELETE /api/v1/videos/:id
func (h *VideoHandler) DeleteVideo(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.videoService.DeleteVideo(c.Context(), userID, videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video deleted successfully", nil)
}

// Admin handlers
// GET /api/v1/admin/videos
func (h *VideoHandler) GetAllVideos(c *fiber.Ctx) error {
	var params dto.VideoQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	videos, err := h.videoService.GetAllVideos(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "All videos retrieved successfully", videos)
}

// PUT /api/v1/admin/videos/:id/hide
func (h *VideoHandler) HideVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.videoService.HideVideo(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video hidden successfully", nil)
}

// PUT /api/v1/admin/videos/:id/show
func (h *VideoHandler) ShowVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.videoService.ShowVideo(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video shown successfully", nil)
}

// DELETE /api/v1/admin/videos/:id
func (h *VideoHandler) DeleteVideoByAdmin(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.videoService.DeleteVideoByAdmin(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Video deleted successfully", nil)
}
```

---

### 9. Routes (`interfaces/api/routes/video_routes.go`)

```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupVideoRoutes(app *fiber.App, videoHandler *handlers.VideoHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// Public routes (with optional auth for view tracking)
	videos := api.Group("/videos")
	videos.Get("/", videoHandler.GetVideos)                  // GET /api/v1/videos
	videos.Get("/:id", videoHandler.GetVideoByID)            // GET /api/v1/videos/:id
	videos.Get("/user/:userId", videoHandler.GetUserVideos)  // GET /api/v1/videos/user/:userId

	// Protected user routes
	videosAuth := videos.Use(authMiddleware.Protected())
	videosAuth.Post("/", videoHandler.UploadVideo)           // POST /api/v1/videos
	videosAuth.Put("/:id", videoHandler.UpdateVideo)         // PUT /api/v1/videos/:id
	videosAuth.Delete("/:id", videoHandler.DeleteVideo)      // DELETE /api/v1/videos/:id

	// Admin routes
	admin := api.Group("/admin", authMiddleware.Protected(), authMiddleware.RequireRole("admin"))
	adminVideos := admin.Group("/videos")
	adminVideos.Get("/", videoHandler.GetAllVideos)              // GET /api/v1/admin/videos
	adminVideos.Put("/:id/hide", videoHandler.HideVideo)         // PUT /api/v1/admin/videos/:id/hide
	adminVideos.Put("/:id/show", videoHandler.ShowVideo)         // PUT /api/v1/admin/videos/:id/show
	adminVideos.Delete("/:id", videoHandler.DeleteVideoByAdmin)  // DELETE /api/v1/admin/videos/:id
}
```

---

### 10. Container Updates (`pkg/di/container.go`)

```go
// Add to container.go

func (c *Container) InitializeVideoComponents() {
	// Repository
	c.VideoRepository = postgres.NewVideoRepository(c.DB)

	// Service
	c.VideoService = serviceimpl.NewVideoService(
		c.VideoRepository,
		c.UserRepository,
		c.StorageService,
	)

	// Handler
	c.VideoHandler = handlers.NewVideoHandler(c.VideoService)
}

// Add to Container struct
type Container struct {
	// ... existing fields
	VideoRepository repositories.VideoRepository
	VideoService    services.VideoService
	VideoHandler    *handlers.VideoHandler
	StorageService  utils.StorageService
}
```

---

### 11. Main Updates (`cmd/api/main.go`)

```go
// Add to main.go

func main() {
	// ... existing code

	// Initialize storage service
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:3000/uploads"
	}
	container.StorageService = utils.NewLocalStorageService(uploadDir, baseURL)

	// Initialize video components
	container.InitializeVideoComponents()

	// Setup routes
	routes.SetupVideoRoutes(app, container.VideoHandler, container.AuthMiddleware)

	// Serve static files (for uploaded videos)
	app.Static("/uploads", uploadDir)

	// ... rest of the code
}
```

---

### 12. Environment Variables (`.env`)

```env
# Storage
UPLOAD_DIR=./uploads
BASE_URL=http://localhost:3000/uploads

# For production with Bunny CDN
# BUNNY_STORAGE_ZONE=your-storage-zone
# BUNNY_API_KEY=your-api-key
# BUNNY_CDN_URL=https://your-cdn.b-cdn.net
```

---

## ✅ Checklist

### Models & Database
- [ ] สร้าง `Video` model
- [ ] เพิ่ม migration สำหรับ `videos` table
- [ ] สร้าง indexes ที่จำเป็น
- [ ] ทดสอบ relationships (Video -> User)

### DTOs
- [ ] สร้าง `UploadVideoRequest`
- [ ] สร้าง `UpdateVideoRequest`
- [ ] สร้าง `VideoQueryParams`
- [ ] สร้าง `VideoResponse` และ `VideoListResponse`
- [ ] สร้าง `UserSummary` (ถ้ายังไม่มี)

### Repository
- [ ] สร้าง `VideoRepository` interface
- [ ] Implement `VideoRepositoryImpl`
- [ ] ทดสอบทุก methods

### Service
- [ ] สร้าง `VideoService` interface
- [ ] Implement `VideoServiceImpl`
- [ ] สร้าง `StorageService` (Local/Bunny CDN)
- [ ] ทดสอบ upload logic
- [ ] ทดสอบ view count increment

### Handlers
- [ ] สร้าง `VideoHandler`
- [ ] Implement upload endpoint
- [ ] Implement list/get endpoints
- [ ] Implement update/delete endpoints
- [ ] Implement admin endpoints

### Routes
- [ ] Setup video routes
- [ ] ใช้ auth middleware
- [ ] ใช้ role middleware สำหรับ admin
- [ ] Setup static file serving

### Integration
- [ ] Register components ใน Container
- [ ] Update main.go
- [ ] Setup environment variables
- [ ] ทดสอบ end-to-end

---

## 🧪 Testing Guide

### 1. Upload Video (User)
```bash
curl -X POST http://localhost:3000/api/v1/videos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "title=My First Video" \
  -F "description=This is a test video" \
  -F "video=@/path/to/video.mp4" \
  -F "thumbnail=@/path/to/thumb.jpg"
```

### 2. Get Video Feed
```bash
curl http://localhost:3000/api/v1/videos?page=1&limit=20&sortBy=newest
```

### 3. Get Video by ID
```bash
curl http://localhost:3000/api/v1/videos/{video-id}
```

### 4. Get User's Videos
```bash
curl http://localhost:3000/api/v1/videos/user/{user-id}?page=1&limit=10
```

### 5. Update Video
```bash
curl -X PUT http://localhost:3000/api/v1/videos/{video-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "description": "Updated description"
  }'
```

### 6. Delete Video
```bash
curl -X DELETE http://localhost:3000/api/v1/videos/{video-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 7. Admin: Get All Videos
```bash
curl http://localhost:3000/api/v1/admin/videos?isActive=false \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### 8. Admin: Hide Video
```bash
curl -X PUT http://localhost:3000/api/v1/admin/videos/{video-id}/hide \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

---

## 📝 Notes

### File Upload:
- รองรับ formats: `.mp4`, `.mov`, `.avi`, `.mkv`
- ขนาดไฟล์สูงสุด: 500MB
- Thumbnail เป็น optional

### Storage:
- **Development**: Local file system
- **Production**: แนะนำใช้ Bunny CDN หรือ AWS S3
- ต้องสร้าง folder `./uploads/videos` และ `./uploads/thumbnails`

## ในระบบมี bunny storage อยู่แล้ว 
- infrastructure\storage\bunny_storage.go
- infrastructure\postgres\file_repository_impl.go
- interfaces\api\handlers\file_handler.go
- interfaces\api\routes\file_routes.go

### Video Metadata:
- Duration, Width, Height ต้องใช้ `ffmpeg` ในการ extract
- สามารถเพิ่ม package `github.com/xfrr/goffmpeg` สำหรับการประมวลผล

### Performance:
- View count increment ทำแบบ asynchronous
- ใช้ pagination สำหรับ video feed
- สามารถเพิ่ม caching (Redis) สำหรับ popular videos

### Security:
- ตรวจสอบ file type และขนาดก่อน upload
- ใช้ virus scanning สำหรับ production
- จำกัด rate limit สำหรับ upload endpoint

### TODO (Advanced):
- [ ] Video transcoding (convert to multiple resolutions)
- [ ] Generate thumbnail automatically with ffmpeg
- [ ] Video streaming with HLS
- [ ] CDN integration (Bunny, Cloudflare)
- [ ] Video compression
- [ ] Watermark

---

**ระยะเวลาโดยประมาณ:** 3-4 วัน

**Dependencies:**
- User System (มีอยู่แล้ว)
- File upload library
- Optional: ffmpeg สำหรับ video processing

**ไฟล์ที่เกี่ยวข้อง:**
- `task_00_admin_forum.md` (สำหรับ reference)
- `task_01_topic_reply.md` (สำหรับ reference)
