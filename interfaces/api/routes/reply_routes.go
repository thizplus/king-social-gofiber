package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupReplyRoutes(api fiber.Router, h *handlers.Handlers) {
	replies := api.Group("/replies")
	replies.Use(middleware.Protected())

	replies.Put("/:id", h.ReplyHandler.UpdateReply)
	replies.Delete("/:id", h.ReplyHandler.DeleteReply)

	// Admin routes
	adminReplies := api.Group("/admin/replies")
	adminReplies.Use(middleware.Protected())
	adminReplies.Use(middleware.AdminOnly())

	adminReplies.Delete("/:id", h.ReplyHandler.DeleteReplyByAdmin)
}
