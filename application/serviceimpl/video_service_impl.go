package serviceimpl

import (
	"context"
	"errors"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"math"

	"github.com/google/uuid"
)

type videoServiceImpl struct {
	videoRepo repositories.VideoRepository
	fileRepo  repositories.FileRepository
	userRepo  repositories.UserRepository
}

func NewVideoService(
	videoRepo repositories.VideoRepository,
	fileRepo repositories.FileRepository,
	userRepo repositories.UserRepository,
) services.VideoService {
	return &videoServiceImpl{
		videoRepo: videoRepo,
		fileRepo:  fileRepo,
		userRepo:  userRepo,
	}
}

func (s *videoServiceImpl) CreateVideo(ctx context.Context, userID uuid.UUID, req *dto.UploadVideoRequest) (*dto.VideoResponse, error) {
	// Verify user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify video file exists
	videoFile, err := s.fileRepo.GetByID(ctx, req.VideoFileID)
	if err != nil {
		return nil, errors.New("video file not found")
	}

	// Verify file belongs to user
	if videoFile.UserID != userID {
		return nil, errors.New("you don't have permission to use this file")
	}

	// Get thumbnail URL if provided
	var thumbnailURL string
	if req.ThumbnailID != uuid.Nil {
		thumbnailFile, err := s.fileRepo.GetByID(ctx, req.ThumbnailID)
		if err == nil && thumbnailFile.UserID == userID {
			thumbnailURL = thumbnailFile.URL
		}
	}

	// Create video record
	video := &models.Video{
		UserID:       userID,
		Title:        req.Title,
		Description:  req.Description,
		VideoURL:     videoFile.URL,
		ThumbnailURL: thumbnailURL,
		Duration:     req.Duration,
		Width:        req.Width,
		Height:       req.Height,
		FileSize:     videoFile.FileSize,
		IsActive:     true,
	}

	if err := s.videoRepo.Create(ctx, video); err != nil {
		return nil, err
	}

	// Load user for response
	video.User = user

	return dto.VideoToVideoResponse(video), nil
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

	return dto.VideoToVideoResponse(video), nil
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

	return dto.VideoToVideoResponse(video), nil
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
	return s.videoRepo.Delete(ctx, videoID)
}

// Helper methods
func (s *videoServiceImpl) toVideoListResponse(videos []models.Video, totalCount int64, params *dto.VideoQueryParams) *dto.VideoListResponse {
	videoResponses := make([]dto.VideoResponse, len(videos))
	for i, video := range videos {
		videoResponses[i] = *dto.VideoToVideoResponse(&video)
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
	hasNext := page < totalPages
	hasPrevious := page > 1

	return &dto.VideoListResponse{
		Videos: videoResponses,
		Meta: dto.PaginationMeta{
			Total:       totalCount,
			Offset:      (page - 1) * limit,
			Limit:       limit,
			Page:        page,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
	}
}
