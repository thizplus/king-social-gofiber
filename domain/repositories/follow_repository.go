package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-social/domain/models"
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
