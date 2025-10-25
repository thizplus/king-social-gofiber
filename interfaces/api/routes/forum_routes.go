package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupForumRoutes(api fiber.Router, h *handlers.Handlers) {
	// Admin Routes
	adminForums := api.Group("/admin/forums")
	adminForums.Use(middleware.Protected())
	adminForums.Use(middleware.AdminOnly())

	adminForums.Post("/", h.ForumHandler.CreateForum)
	adminForums.Get("/", h.ForumHandler.GetAllForums)
	adminForums.Put("/reorder", h.ForumHandler.ReorderForums)
	adminForums.Put("/:id", h.ForumHandler.UpdateForum)
	adminForums.Delete("/:id", h.ForumHandler.DeleteForum)
	adminForums.Post("/:id/sync-topic-count", h.ForumHandler.SyncTopicCount)
	adminForums.Post("/sync-all-topic-counts", h.ForumHandler.SyncAllTopicCounts)

	// Public Routes
	forums := api.Group("/forums")

	forums.Get("/", h.ForumHandler.GetActiveForums)
	forums.Get("/:id", h.ForumHandler.GetForumByID)
	forums.Get("/slug/:slug", h.ForumHandler.GetForumBySlug)
}
