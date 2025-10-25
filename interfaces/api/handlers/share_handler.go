package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ShareHandler struct {
	shareService services.ShareService
}

func NewShareHandler(shareService services.ShareService) *ShareHandler {
	return &ShareHandler{shareService: shareService}
}

// ShareVideo handles sharing a video
// POST /api/v1/videos/:id/share
func (h *ShareHandler) ShareVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	var req dto.ShareVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set video ID from URL parameter
	req.VideoID = videoID

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	share, err := h.shareService.ShareVideo(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to share video", err)
	}

	return utils.SuccessResponse(c, "Video shared successfully", share)
}

// GetShareCount handles getting share count for a video
// GET /api/v1/videos/:id/share/count
func (h *ShareHandler) GetShareCount(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	count, err := h.shareService.GetShareCount(c.Context(), videoID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get share count", err)
	}

	return utils.SuccessResponse(c, "Share count retrieved successfully", count)
}
