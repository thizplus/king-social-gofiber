package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRoutes(api fiber.Router, h *handlers.Handlers) {
	// Protected routes - all notification routes require authentication
	notifications := api.Group("/notifications", middleware.Protected())

	notifications.Get("/", h.NotificationHandler.GetNotifications)                  // GET /api/v1/notifications
	notifications.Get("/unread/count", h.NotificationHandler.GetUnreadCount)        // GET /api/v1/notifications/unread/count
	notifications.Put("/:id/read", h.NotificationHandler.MarkAsRead)                // PUT /api/v1/notifications/:id/read
	notifications.Put("/read", h.NotificationHandler.MarkMultipleAsRead)            // PUT /api/v1/notifications/read
	notifications.Put("/read-all", h.NotificationHandler.MarkAllAsRead)             // PUT /api/v1/notifications/read-all
	notifications.Delete("/:id", h.NotificationHandler.DeleteNotification)          // DELETE /api/v1/notifications/:id
}
