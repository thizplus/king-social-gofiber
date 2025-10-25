package handlers

import (
	"strconv"

	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FollowHandler struct {
	followService services.FollowService
}

func NewFollowHandler(followService services.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// FollowUser handles following a user
// POST /api/v1/users/:userId/follow
func (h *FollowHandler) FollowUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	result, err := h.followService.FollowUser(c.Context(), user.ID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to follow user", err)
	}

	return utils.SuccessResponse(c, "User followed successfully", result)
}

// UnfollowUser handles unfollowing a user
// DELETE /api/v1/users/:userId/follow
func (h *FollowHandler) UnfollowUser(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	if err := h.followService.UnfollowUser(c.Context(), user.ID, followingID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unfollow user", err)
	}

	return utils.SuccessResponse(c, "User unfollowed successfully", nil)
}

// GetFollowStatus handles getting follow status
// GET /api/v1/users/:userId/follow/status
func (h *FollowHandler) GetFollowStatus(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	followingID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	status, err := h.followService.GetFollowStatus(c.Context(), user.ID, followingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get follow status", err)
	}

	return utils.SuccessResponse(c, "Follow status retrieved successfully", status)
}

// GetFollowers handles getting a user's followers
// GET /api/v1/users/:userId/followers
func (h *FollowHandler) GetFollowers(c *fiber.Ctx) error {
	targetUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	// Get current user ID (optional - for checking if current user follows each follower)
	currentUserID := uuid.Nil
	user, err := utils.GetUserFromContext(c)
	if err == nil {
		currentUserID = user.ID
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	followers, err := h.followService.GetFollowers(c.Context(), currentUserID, targetUserID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get followers", err)
	}

	return utils.SuccessResponse(c, "Followers retrieved successfully", followers)
}

// GetFollowing handles getting users that a user is following
// GET /api/v1/users/:userId/following
func (h *FollowHandler) GetFollowing(c *fiber.Ctx) error {
	targetUserID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	// Get current user ID (optional - for checking if current user follows each user)
	currentUserID := uuid.Nil
	user, err := utils.GetUserFromContext(c)
	if err == nil {
		currentUserID = user.ID
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	following, err := h.followService.GetFollowing(c.Context(), currentUserID, targetUserID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get following", err)
	}

	return utils.SuccessResponse(c, "Following retrieved successfully", following)
}

// GetUserStats handles getting a user's follower and following stats
// GET /api/v1/users/:userId/stats
func (h *FollowHandler) GetUserStats(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	stats, err := h.followService.GetUserStats(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user stats", err)
	}

	return utils.SuccessResponse(c, "User stats retrieved successfully", stats)
}
