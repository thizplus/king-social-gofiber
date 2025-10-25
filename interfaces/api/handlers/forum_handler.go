package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ForumHandler struct {
	forumService services.ForumService
}

func NewForumHandler(forumService services.ForumService) *ForumHandler {
	return &ForumHandler{forumService: forumService}
}

// Admin Handlers
func (h *ForumHandler) CreateForum(c *fiber.Ctx) error {
	admin, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateForumRequest
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

	forum, err := h.forumService.CreateForum(c.Context(), admin.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create forum", err)
	}

	return utils.SuccessResponse(c, "Forum created successfully", dto.ForumToForumResponse(forum))
}

func (h *ForumHandler) UpdateForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	var req dto.UpdateForumRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	forum, err := h.forumService.UpdateForum(c.Context(), forumID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update forum", err)
	}

	return utils.SuccessResponse(c, "Forum updated successfully", dto.ForumToForumResponse(forum))
}

func (h *ForumHandler) DeleteForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	if err := h.forumService.DeleteForum(c.Context(), forumID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete forum", err)
	}

	return utils.SuccessResponse(c, "Forum deleted successfully", nil)
}

func (h *ForumHandler) GetAllForums(c *fiber.Ctx) error {
	includeInactive := c.Query("includeInactive", "false") == "true"

	forums, err := h.forumService.GetAllForums(c.Context(), includeInactive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get forums", err)
	}

	return utils.SuccessResponse(c, "Forums retrieved successfully", fiber.Map{
		"forums": forums,
		"total":  len(forums),
	})
}

func (h *ForumHandler) ReorderForums(c *fiber.Ctx) error {
	var req dto.ReorderForumsRequest
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

	if err := h.forumService.ReorderForums(c.Context(), &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to reorder forums", err)
	}

	return utils.SuccessResponse(c, "Forums reordered successfully", nil)
}

// Public Handlers
func (h *ForumHandler) GetActiveForums(c *fiber.Ctx) error {
	forums, err := h.forumService.GetActiveForums(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get forums", err)
	}

	return utils.SuccessResponse(c, "Forums retrieved successfully", fiber.Map{
		"forums": forums,
		"total":  len(forums),
	})
}

func (h *ForumHandler) GetForumByID(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	forum, err := h.forumService.GetForumByID(c.Context(), forumID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Forum not found", err)
	}

	return utils.SuccessResponse(c, "Forum retrieved successfully", forum)
}

func (h *ForumHandler) GetForumBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	forum, err := h.forumService.GetForumBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Forum not found", err)
	}

	return utils.SuccessResponse(c, "Forum retrieved successfully", forum)
}

// Admin only - Sync topic count
func (h *ForumHandler) SyncTopicCount(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	if err := h.forumService.SyncTopicCount(c.Context(), forumID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to sync topic count", err)
	}

	return utils.SuccessResponse(c, "Topic count synchronized successfully", nil)
}

// Admin only - Sync all forums topic counts
func (h *ForumHandler) SyncAllTopicCounts(c *fiber.Ctx) error {
	if err := h.forumService.SyncAllTopicCounts(c.Context()); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to sync all topic counts", err)
	}

	return utils.SuccessResponse(c, "All topic counts synchronized successfully", nil)
}
