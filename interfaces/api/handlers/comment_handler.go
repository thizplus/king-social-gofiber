package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentHandler struct {
	commentService services.CommentService
}

func NewCommentHandler(commentService services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment handles creating a comment
// POST /api/v1/comments
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	comment, err := h.commentService.CreateComment(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create comment", err)
	}

	return utils.SuccessResponse(c, "Comment created successfully", comment)
}

// GetCommentsByVideoID handles getting comments for a video
// GET /api/v1/videos/:id/comments
func (h *CommentHandler) GetCommentsByVideoID(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid page parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid limit parameter")
	}

	comments, err := h.commentService.GetCommentsByVideoID(c.Context(), videoID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve comments", err)
	}

	return utils.SuccessResponse(c, "Comments retrieved successfully", comments)
}

// UpdateComment handles updating a comment
// PUT /api/v1/comments/:id
func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	var req dto.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	comment, err := h.commentService.UpdateComment(c.Context(), user.ID, commentID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update comment", err)
	}

	return utils.SuccessResponse(c, "Comment updated successfully", comment)
}

// DeleteComment handles deleting a comment
// DELETE /api/v1/comments/:id
func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	if err := h.commentService.DeleteComment(c.Context(), user.ID, commentID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete comment", err)
	}

	return utils.SuccessResponse(c, "Comment deleted successfully", nil)
}

// DeleteCommentByAdmin handles deleting a comment by admin
// DELETE /api/v1/admin/comments/:id
func (h *CommentHandler) DeleteCommentByAdmin(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid comment ID")
	}

	if err := h.commentService.DeleteCommentByAdmin(c.Context(), commentID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete comment", err)
	}

	return utils.SuccessResponse(c, "Comment deleted successfully", nil)
}
