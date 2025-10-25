package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupLikeRoutes(api fiber.Router, h *handlers.Handlers) {
	// Topic like routes (all protected)
	topics := api.Group("/topics")
	topics.Post("/:id/like", middleware.Protected(), h.LikeHandler.LikeTopic)          // POST /api/v1/topics/:id/like
	topics.Delete("/:id/like", middleware.Protected(), h.LikeHandler.UnlikeTopic)      // DELETE /api/v1/topics/:id/like
	topics.Get("/:id/like", middleware.Protected(), h.LikeHandler.GetTopicLikeStatus)  // GET /api/v1/topics/:id/like

	// Video like routes (all protected)
	videos := api.Group("/videos")
	videos.Post("/:id/like", middleware.Protected(), h.LikeHandler.LikeVideo)         // POST /api/v1/videos/:id/like
	videos.Delete("/:id/like", middleware.Protected(), h.LikeHandler.UnlikeVideo)     // DELETE /api/v1/videos/:id/like
	videos.Get("/:id/like", middleware.Protected(), h.LikeHandler.GetVideoLikeStatus) // GET /api/v1/videos/:id/like

	// Reply like routes (all protected)
	replies := api.Group("/replies")
	replies.Post("/:id/like", middleware.Protected(), h.LikeHandler.LikeReply)         // POST /api/v1/replies/:id/like
	replies.Delete("/:id/like", middleware.Protected(), h.LikeHandler.UnlikeReply)     // DELETE /api/v1/replies/:id/like
	replies.Get("/:id/like", middleware.Protected(), h.LikeHandler.GetReplyLikeStatus) // GET /api/v1/replies/:id/like

	// Comment like routes (all protected)
	comments := api.Group("/comments")
	comments.Post("/:id/like", middleware.Protected(), h.LikeHandler.LikeComment)         // POST /api/v1/comments/:id/like
	comments.Delete("/:id/like", middleware.Protected(), h.LikeHandler.UnlikeComment)     // DELETE /api/v1/comments/:id/like
	comments.Get("/:id/like", middleware.Protected(), h.LikeHandler.GetCommentLikeStatus) // GET /api/v1/comments/:id/like
}
