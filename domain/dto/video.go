package dto

import (
	"gofiber-social/domain/models"
	"time"

	"github.com/google/uuid"
)

// ============= Request DTOs =============

type UploadVideoRequest struct {
	Title        string    `json:"title" validate:"required,min=3,max=200"`
	Description  string    `json:"description" validate:"omitempty,max=1000"`
	VideoFileID  uuid.UUID `json:"videoFileId" validate:"required,uuid"`
	ThumbnailID  uuid.UUID `json:"thumbnailId" validate:"omitempty,uuid"`
	Duration     int       `json:"duration" validate:"omitempty,min=0"`
	Width        int       `json:"width" validate:"omitempty,min=0"`
	Height       int       `json:"height" validate:"omitempty,min=0"`
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
	Videos []VideoResponse `json:"videos"`
	Meta   PaginationMeta  `json:"meta"`
}

// UserSummary for video responses
type UserSummary struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Avatar    string    `json:"avatar"`
}

// ============= Converters =============

func VideoToVideoResponse(video *models.Video) *VideoResponse {
	resp := &VideoResponse{
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
		resp.User = &UserSummary{
			ID:        video.User.ID,
			Username:  video.User.Username,
			FirstName: video.User.FirstName,
			LastName:  video.User.LastName,
			Avatar:    video.User.Avatar,
		}
	}

	return resp
}
