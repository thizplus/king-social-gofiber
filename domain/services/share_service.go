package services

import (
	"context"
	"gofiber-social/domain/dto"

	"github.com/google/uuid"
)

type ShareService interface {
	ShareVideo(ctx context.Context, userID uuid.UUID, req *dto.ShareVideoRequest) (*dto.ShareResponse, error)
	GetShareCount(ctx context.Context, videoID uuid.UUID) (*dto.ShareCountResponse, error)
}
