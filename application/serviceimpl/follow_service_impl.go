package serviceimpl

import (
	"context"
	"errors"
	"math"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"

	"github.com/google/uuid"
)

type followServiceImpl struct {
	followRepo          repositories.FollowRepository
	userRepo            repositories.UserRepository
	notificationService services.NotificationService
}

func NewFollowService(
	followRepo repositories.FollowRepository,
	userRepo repositories.UserRepository,
	notificationService services.NotificationService,
) services.FollowService {
	return &followServiceImpl{
		followRepo:          followRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

func (s *followServiceImpl) FollowUser(ctx context.Context, followerID, followingID uuid.UUID) (*dto.FollowResponse, error) {
	// Verify both users exist
	_, err := s.userRepo.FindByID(ctx, followerID)
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

	// Create notification for new follower
	go func() {
		_ = s.notificationService.CreateNewFollowerNotification(context.Background(), followingID, followerID)
	}()

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
