package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupTopicRoutes(api fiber.Router, h *handlers.Handlers) {
	// Public routes
	topics := api.Group("/topics")
	topics.Get("/", h.TopicHandler.GetTopics)
	topics.Get("/search", h.TopicHandler.SearchTopics)
	topics.Get("/:id", h.TopicHandler.GetTopic)
	topics.Get("/:id/replies", h.ReplyHandler.GetReplies)

	// Protected routes
	topicsProtected := api.Group("/topics")
	topicsProtected.Use(middleware.Protected())
	topicsProtected.Post("/", h.TopicHandler.CreateTopic)
	topicsProtected.Put("/:id", h.TopicHandler.UpdateTopic)
	topicsProtected.Delete("/:id", h.TopicHandler.DeleteTopic)
	topicsProtected.Post("/:id/replies", h.ReplyHandler.CreateReply)

	// Forum topics
	api.Get("/forums/:id/topics", h.TopicHandler.GetTopicsByForum)
	api.Get("/forums/slug/:slug/topics", h.TopicHandler.GetTopicsByForumSlug)

	// Topics by tags
	api.Get("/topics/tags/:tag", h.TopicHandler.GetTopicsByTag)
	api.Get("/topics/tags", h.TopicHandler.GetTopicsByTags) // Multiple tags via query params

	// Admin routes
	adminTopics := api.Group("/admin/topics")
	adminTopics.Use(middleware.Protected())
	adminTopics.Use(middleware.AdminOnly())

	adminTopics.Put("/:id/pin", h.TopicHandler.PinTopic)
	adminTopics.Put("/:id/unpin", h.TopicHandler.UnpinTopic)
	adminTopics.Put("/:id/lock", h.TopicHandler.LockTopic)
	adminTopics.Put("/:id/unlock", h.TopicHandler.UnlockTopic)
	adminTopics.Delete("/:id", h.TopicHandler.DeleteTopicByAdmin)
}
