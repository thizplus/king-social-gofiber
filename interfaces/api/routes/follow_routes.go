package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupFollowRoutes(api fiber.Router, h *handlers.Handlers) {
	// Public routes (anyone can view, optional auth for isFollowing status)
	api.Get("/users/:userId/followers", middleware.Optional(), h.FollowHandler.GetFollowers)   // GET /api/v1/users/:userId/followers
	api.Get("/users/:userId/following", middleware.Optional(), h.FollowHandler.GetFollowing)   // GET /api/v1/users/:userId/following
	api.Get("/users/:userId/stats", h.FollowHandler.GetUserStats)                              // GET /api/v1/users/:userId/stats

	// Protected routes (requires authentication)
	api.Post("/users/:userId/follow", middleware.Protected(), h.FollowHandler.FollowUser)                 // POST /api/v1/users/:userId/follow
	api.Delete("/users/:userId/follow", middleware.Protected(), h.FollowHandler.UnfollowUser)             // DELETE /api/v1/users/:userId/follow
	api.Get("/users/:userId/follow/status", middleware.Protected(), h.FollowHandler.GetFollowStatus)      // GET /api/v1/users/:userId/follow/status
}
