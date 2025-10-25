package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoHandler struct {
	videoService services.VideoService
}

func NewVideoHandler(videoService services.VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

// CreateVideo handles video creation (after file upload)
// POST /api/v1/videos
func (h *VideoHandler) CreateVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UploadVideoRequest
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

	video, err := h.videoService.CreateVideo(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video creation failed", err)
	}

	return utils.SuccessResponse(c, "Video created successfully", video)
}

// GetVideos handles listing videos
// GET /api/v1/videos
func (h *VideoHandler) GetVideos(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	sortBy := c.Query("sortBy", "newest")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid page parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid limit parameter")
	}

	params := &dto.VideoQueryParams{
		Page:   page,
		Limit:  limit,
		SortBy: sortBy,
	}

	videos, err := h.videoService.GetVideos(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve videos", err)
	}

	return utils.SuccessResponse(c, "Videos retrieved successfully", videos)
}

// GetVideoByID handles getting a single video
// GET /api/v1/videos/:id
func (h *VideoHandler) GetVideoByID(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	video, err := h.videoService.GetVideoByID(c.Context(), videoID)
	if err != nil {
		return utils.NotFoundResponse(c, "Video not found")
	}

	return utils.SuccessResponse(c, "Video retrieved successfully", video)
}

// GetUserVideos handles getting videos by user
// GET /api/v1/videos/user/:userId
func (h *VideoHandler) GetUserVideos(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
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

	params := &dto.VideoQueryParams{
		Page:  page,
		Limit: limit,
	}

	videos, err := h.videoService.GetUserVideos(c.Context(), userID, params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user videos", err)
	}

	return utils.SuccessResponse(c, "User videos retrieved successfully", videos)
}

// UpdateVideo handles updating video info
// PUT /api/v1/videos/:id
func (h *VideoHandler) UpdateVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	var req dto.UpdateVideoRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	video, err := h.videoService.UpdateVideo(c.Context(), user.ID, videoID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video update failed", err)
	}

	return utils.SuccessResponse(c, "Video updated successfully", video)
}

// DeleteVideo handles deleting a video
// DELETE /api/v1/videos/:id
func (h *VideoHandler) DeleteVideo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	if err := h.videoService.DeleteVideo(c.Context(), user.ID, videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video deletion failed", err)
	}

	return utils.SuccessResponse(c, "Video deleted successfully", nil)
}

// Admin handlers

// GetAllVideos handles getting all videos (including inactive)
// GET /api/v1/admin/videos
func (h *VideoHandler) GetAllVideos(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "20")
	isActiveStr := c.Query("isActive")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid page parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid limit parameter")
	}

	params := &dto.VideoQueryParams{
		Page:  page,
		Limit: limit,
	}

	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		params.IsActive = &isActive
	}

	videos, err := h.videoService.GetAllVideos(c.Context(), params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve videos", err)
	}

	return utils.SuccessResponse(c, "All videos retrieved successfully", videos)
}

// HideVideo handles hiding a video
// PUT /api/v1/admin/videos/:id/hide
func (h *VideoHandler) HideVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	if err := h.videoService.HideVideo(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to hide video", err)
	}

	return utils.SuccessResponse(c, "Video hidden successfully", nil)
}

// ShowVideo handles showing a video
// PUT /api/v1/admin/videos/:id/show
func (h *VideoHandler) ShowVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	if err := h.videoService.ShowVideo(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to show video", err)
	}

	return utils.SuccessResponse(c, "Video shown successfully", nil)
}

// DeleteVideoByAdmin handles deleting a video by admin
// DELETE /api/v1/admin/videos/:id
func (h *VideoHandler) DeleteVideoByAdmin(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid video ID")
	}

	if err := h.videoService.DeleteVideoByAdmin(c.Context(), videoID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video deletion failed", err)
	}

	return utils.SuccessResponse(c, "Video deleted successfully", nil)
}
