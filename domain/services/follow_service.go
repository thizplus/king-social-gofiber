package services

import (
	"context"
	"gofiber-social/domain/dto"

	"github.com/google/uuid"
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
