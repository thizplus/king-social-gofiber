package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReplyHandler struct {
	replyService services.ReplyService
}

func NewReplyHandler(replyService services.ReplyService) *ReplyHandler {
	return &ReplyHandler{replyService: replyService}
}

func (h *ReplyHandler) CreateReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	var req dto.CreateReplyRequest
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

	reply, err := h.replyService.CreateReply(c.Context(), topicID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create reply", err)
	}

	response := dto.ReplyToReplyResponse(reply, false)
	return utils.SuccessResponse(c, "Reply created successfully", response)
}

func (h *ReplyHandler) GetReplies(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	replies, total, err := h.replyService.GetReplies(c.Context(), topicID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get replies", err)
	}

	return utils.SuccessResponse(c, "Replies retrieved successfully", fiber.Map{
		"replies": replies,
		"meta":    dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *ReplyHandler) UpdateReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	var req dto.UpdateReplyRequest
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

	reply, err := h.replyService.UpdateReply(c.Context(), replyID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update reply", err)
	}

	response := dto.ReplyToReplyResponse(reply, false)
	return utils.SuccessResponse(c, "Reply updated successfully", response)
}

func (h *ReplyHandler) DeleteReply(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	if err := h.replyService.DeleteReply(c.Context(), replyID, user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete reply", err)
	}

	return utils.SuccessResponse(c, "Reply deleted successfully", nil)
}

func (h *ReplyHandler) DeleteReplyByAdmin(c *fiber.Ctx) error {
	replyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid reply ID")
	}

	if err := h.replyService.DeleteReplyByAdmin(c.Context(), replyID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete reply", err)
	}

	return utils.SuccessResponse(c, "Reply deleted successfully", nil)
}
