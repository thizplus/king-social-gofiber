package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TagHandler struct {
	tagService services.TagService
}

func NewTagHandler(tagService services.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var req dto.CreateTagRequest
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

	tag, err := h.tagService.CreateTag(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create tag", err)
	}

	return utils.SuccessResponse(c, "Tag created successfully", dto.TagToTagResponse(tag))
}

func (h *TagHandler) GetTags(c *fiber.Ctx) error {
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	activeOnly := c.Query("activeOnly", "false") == "true"

	tags, total, err := h.tagService.GetTags(c.Context(), offset, limit, activeOnly)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get tags", err)
	}

	return utils.SuccessResponse(c, "Tags retrieved successfully", fiber.Map{
		"tags": tags,
		"meta": dto.NewPaginationMeta(int64(total), offset, limit),
	})
}

func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tag ID")
	}

	tag, err := h.tagService.GetTag(c.Context(), tagID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Tag not found", err)
	}

	return utils.SuccessResponse(c, "Tag retrieved successfully", dto.TagToTagResponse(tag))
}

func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tag ID")
	}

	var req dto.UpdateTagRequest
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

	tag, err := h.tagService.UpdateTag(c.Context(), tagID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update tag", err)
	}

	return utils.SuccessResponse(c, "Tag updated successfully", dto.TagToTagResponse(tag))
}

func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	tagID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tag ID")
	}

	if err := h.tagService.DeleteTag(c.Context(), tagID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete tag", err)
	}

	return utils.SuccessResponse(c, "Tag deleted successfully", nil)
}

func (h *TagHandler) SearchTags(c *fiber.Ctx) error {
	query := c.Query("q", "")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	tags, total, err := h.tagService.SearchTags(c.Context(), query, offset, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search tags", err)
	}

	return utils.SuccessResponse(c, "Tags found successfully", fiber.Map{
		"tags": tags,
		"meta": dto.NewPaginationMeta(int64(total), offset, limit),
	})
}