package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupCommentRoutes(api fiber.Router, h *handlers.Handlers) {
	// Comment routes
	comments := api.Group("/comments")
	comments.Post("/", middleware.Protected(), h.CommentHandler.CreateComment)       // POST /api/v1/comments
	comments.Put("/:id", middleware.Protected(), h.CommentHandler.UpdateComment)     // PUT /api/v1/comments/:id
	comments.Delete("/:id", middleware.Protected(), h.CommentHandler.DeleteComment)  // DELETE /api/v1/comments/:id

	// Video comments route (public)
	videos := api.Group("/videos")
	videos.Get("/:id/comments", h.CommentHandler.GetCommentsByVideoID)  // GET /api/v1/videos/:id/comments

	// Admin routes
	adminComments := api.Group("/admin/comments")
	adminComments.Use(middleware.Protected(), middleware.AdminOnly())
	adminComments.Delete("/:id", h.CommentHandler.DeleteCommentByAdmin)  // DELETE /api/v1/admin/comments/:id
}
