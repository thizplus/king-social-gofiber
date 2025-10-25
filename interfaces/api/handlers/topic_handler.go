package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TopicHandler struct {
	topicService services.TopicService
}

func NewTopicHandler(topicService services.TopicService) *TopicHandler {
	return &TopicHandler{topicService: topicService}
}

func (h *TopicHandler) CreateTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateTopicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	topic, err := h.topicService.CreateTopic(c.Context(), user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create topic", err)
	}

	response := dto.TopicToTopicResponse(topic)
	return utils.SuccessResponse(c, "Topic created successfully", response)
}

func (h *TopicHandler) GetTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	topic, err := h.topicService.GetTopic(c.Context(), topicID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Topic not found", err)
	}

	return utils.SuccessResponse(c, "Topic retrieved successfully", topic)
}

func (h *TopicHandler) GetTopics(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopics(c.Context(), offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) GetTopicsByForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopicsByForum(c.Context(), forumID, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) GetTopicsByForumSlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return utils.ValidationErrorResponse(c, "Forum slug is required")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopicsByForumSlug(c.Context(), slug, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) UpdateTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	var req dto.UpdateTopicRequest
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

	topic, err := h.topicService.UpdateTopic(c.Context(), topicID, user.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update topic", err)
	}

	return utils.SuccessResponse(c, "Topic updated successfully", topic)
}

func (h *TopicHandler) DeleteTopic(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.DeleteTopic(c.Context(), topicID, user.ID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete topic", err)
	}

	return utils.SuccessResponse(c, "Topic deleted successfully", nil)
}

func (h *TopicHandler) SearchTopics(c *fiber.Ctx) error {
	query := c.Query("q", "")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.SearchTopics(c.Context(), query, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search topics", err)
	}

	return utils.SuccessResponse(c, "Topics found successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) GetTopicsByTag(c *fiber.Ctx) error {
	tag := c.Params("tag")
	if tag == "" {
		return utils.ValidationErrorResponse(c, "Tag is required")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopicsByTag(c.Context(), tag, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics by tag", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

func (h *TopicHandler) GetTopicsByTags(c *fiber.Ctx) error {
	tagsQuery := c.Query("tags", "")
	if tagsQuery == "" {
		return utils.ValidationErrorResponse(c, "Tags parameter is required")
	}

	// Split tags by comma and trim spaces
	tags := make([]string, 0)
	for _, tag := range strings.Split(tagsQuery, ",") {
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			tags = append(tags, trimmedTag)
		}
	}

	if len(tags) == 0 {
		return utils.ValidationErrorResponse(c, "At least one valid tag is required")
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	topics, total, err := h.topicService.GetTopicsByTags(c.Context(), tags, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get topics by tags", err)
	}

	return utils.SuccessResponse(c, "Topics retrieved successfully", fiber.Map{
		"topics": topics,
		"meta":   dto.NewPaginationMeta(total, offset, limit),
	})
}

// Admin Handlers
func (h *TopicHandler) PinTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.PinTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to pin topic", err)
	}

	return utils.SuccessResponse(c, "Topic pinned successfully", nil)
}

func (h *TopicHandler) UnpinTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.UnpinTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unpin topic", err)
	}

	return utils.SuccessResponse(c, "Topic unpinned successfully", nil)
}

func (h *TopicHandler) LockTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.LockTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to lock topic", err)
	}

	return utils.SuccessResponse(c, "Topic locked successfully", nil)
}

func (h *TopicHandler) UnlockTopic(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.UnlockTopic(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to unlock topic", err)
	}

	return utils.SuccessResponse(c, "Topic unlocked successfully", nil)
}

func (h *TopicHandler) DeleteTopicByAdmin(c *fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid topic ID")
	}

	if err := h.topicService.DeleteTopicByAdmin(c.Context(), topicID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete topic", err)
	}

	return utils.SuccessResponse(c, "Topic deleted successfully", nil)
}
