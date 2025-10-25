package serviceimpl

import (
	"context"
	"errors"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"

	"github.com/google/uuid"
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
	_, err := s.videoRepo.FindByID(ctx, req.VideoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	// Create share record
	share := &models.Share{
		UserID:   userID,
		VideoID:  req.VideoID,
		Platform: req.Platform,
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, err
	}

	return &dto.ShareResponse{
		ID:        share.ID,
		UserID:    share.UserID,
		VideoID:   share.VideoID,
		Platform:  share.Platform,
		CreatedAt: share.CreatedAt,
		Message:   "Video shared successfully",
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
