package handlers

import (
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LikeHandler struct {
	likeService services.LikeService
}

func NewLikeHandler(likeService services.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// Topic Like Handlers

// LikeTopic handles liking a topic
// POST /api/v1/topics/:id/like
func (h *LikeHandler) LikeTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	status, err := h.likeService.LikeTopic(c.Context(), user.ID, topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to like topic", err)
	}

	return utils.SuccessResponse(c, "Topic liked successfully", status)
}

// UnlikeTopic handles unliking a topic
// DELETE /api/v1/topics/:id/like
func (h *LikeHandler) UnlikeTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	status, err := h.likeService.UnlikeTopic(c.Context(), user.ID, topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlike topic", err)
	}

	return utils.SuccessResponse(c, "Topic unliked successfully", status)
}

// GetTopicLikeStatus handles getting topic like status
// GET /api/v1/topics/:id/like
func (h *LikeHandler) GetTopicLikeStatus(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	status, err := h.likeService.GetTopicLikeStatus(c.Context(), user.ID, topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get like status", err)
	}

	return utils.SuccessResponse(c, "Topic like status retrieved successfully", status)
}

// Video Like Handlers

// LikeVideo handles liking a video
// POST /api/v1/videos/:id/like
func (h *LikeHandler) LikeVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	status, err := h.likeService.LikeVideo(c.Context(), user.ID, videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to like video", err)
	}

	return utils.SuccessResponse(c, "Video liked successfully", status)
}

// UnlikeVideo handles unliking a video
// DELETE /api/v1/videos/:id/like
func (h *LikeHandler) UnlikeVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	status, err := h.likeService.UnlikeVideo(c.Context(), user.ID, videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlike video", err)
	}

	return utils.SuccessResponse(c, "Video unliked successfully", status)
}

// GetVideoLikeStatus handles getting video like status
// GET /api/v1/videos/:id/like
func (h *LikeHandler) GetVideoLikeStatus(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	status, err := h.likeService.GetVideoLikeStatus(c.Context(), user.ID, videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get like status", err)
	}

	return utils.SuccessResponse(c, "Video like status retrieved successfully", status)
}

// Reply Like Handlers

// LikeReply handles liking a reply
// POST /api/v1/replies/:id/like
func (h *LikeHandler) LikeReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	status, err := h.likeService.LikeReply(c.Context(), user.ID, replyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to like reply", err)
	}

	return utils.SuccessResponse(c, "Reply liked successfully", status)
}

// UnlikeReply handles unliking a reply
// DELETE /api/v1/replies/:id/like
func (h *LikeHandler) UnlikeReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	status, err := h.likeService.UnlikeReply(c.Context(), user.ID, replyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlike reply", err)
	}

	return utils.SuccessResponse(c, "Reply unliked successfully", status)
}

// GetReplyLikeStatus handles getting reply like status
// GET /api/v1/replies/:id/like
func (h *LikeHandler) GetReplyLikeStatus(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	status, err := h.likeService.GetReplyLikeStatus(c.Context(), user.ID, replyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get like status", err)
	}

	return utils.SuccessResponse(c, "Reply like status retrieved successfully", status)
}

// Comment Like Handlers

// LikeComment handles liking a comment
// POST /api/v1/comments/:id/like
func (h *LikeHandler) LikeComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	status, err := h.likeService.LikeComment(c.Context(), user.ID, commentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to like comment", err)
	}

	return utils.SuccessResponse(c, "Comment liked successfully", status)
}

// UnlikeComment handles unliking a comment
// DELETE /api/v1/comments/:id/like
func (h *LikeHandler) UnlikeComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	status, err := h.likeService.UnlikeComment(c.Context(), user.ID, commentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlike comment", err)
	}

	return utils.SuccessResponse(c, "Comment unliked successfully", status)
}

// GetCommentLikeStatus handles getting comment like status
// GET /api/v1/comments/:id/like
func (h *LikeHandler) GetCommentLikeStatus(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	status, err := h.likeService.GetCommentLikeStatus(c.Context(), user.ID, commentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get like status", err)
	}

	return utils.SuccessResponse(c, "Comment like status retrieved successfully", status)
}
