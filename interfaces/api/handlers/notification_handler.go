package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

// GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var params dto.NotificationQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid query parameters")
	}

	if err := utils.ValidateStruct(&params); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	notifications, err := h.notificationService.GetNotifications(c.Context(), user.ID, &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get notifications", err)
	}

	return utils.SuccessResponse(c, "Notifications retrieved successfully", notifications)
}

// GET /api/v1/notifications/unread/count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	count, err := h.notificationService.GetUnreadCount(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get unread count", err)
	}

	return utils.SuccessResponse(c, "Unread count retrieved", count)
}

// PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid notification ID")
	}

	if err := h.notificationService.MarkAsRead(c.Context(), user.ID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark as read", err)
	}

	return utils.SuccessResponse(c, "Notification marked as read", nil)
}

// PUT /api/v1/notifications/read
func (h *NotificationHandler) MarkMultipleAsRead(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.MarkAsReadRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	result, err := h.notificationService.MarkMultipleAsRead(c.Context(), user.ID, req.NotificationIDs)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark as read", err)
	}

	return utils.SuccessResponse(c, "Notifications marked as read", result)
}

// PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	result, err := h.notificationService.MarkAllAsRead(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to mark all as read", err)
	}

	return utils.SuccessResponse(c, "All notifications marked as read", result)
}

// DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	notificationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid notification ID")
	}

	if err := h.notificationService.DeleteNotification(c.Context(), user.ID, notificationID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete notification", err)
	}

	return utils.SuccessResponse(c, "Notification deleted successfully", nil)
}
