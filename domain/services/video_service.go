package services

import (
	"context"
	"gofiber-social/domain/dto"

	"github.com/google/uuid"
)

type VideoService interface {
	// User operations
	CreateVideo(ctx context.Context, userID uuid.UUID, req *dto.UploadVideoRequest) (*dto.VideoResponse, error)
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
