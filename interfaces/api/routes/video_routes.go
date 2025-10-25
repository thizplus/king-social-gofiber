package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupVideoRoutes(api fiber.Router, h *handlers.Handlers) {
	// Public routes
	videos := api.Group("/videos")
	videos.Get("/", h.VideoHandler.GetVideos)                   // GET /api/v1/videos
	videos.Get("/:id", h.VideoHandler.GetVideoByID)             // GET /api/v1/videos/:id
	videos.Get("/user/:userId", h.VideoHandler.GetUserVideos)   // GET /api/v1/videos/user/:userId

	// Protected user routes (requires authentication)
	videos.Post("/", middleware.Protected(), h.VideoHandler.CreateVideo)            // POST /api/v1/videos
	videos.Put("/:id", middleware.Protected(), h.VideoHandler.UpdateVideo)          // PUT /api/v1/videos/:id
	videos.Delete("/:id", middleware.Protected(), h.VideoHandler.DeleteVideo)       // DELETE /api/v1/videos/:id

	// Admin routes
	adminVideos := api.Group("/admin/videos")
	adminVideos.Use(middleware.Protected(), middleware.AdminOnly())
	adminVideos.Get("/", h.VideoHandler.GetAllVideos)              // GET /api/v1/admin/videos
	adminVideos.Put("/:id/hide", h.VideoHandler.HideVideo)         // PUT /api/v1/admin/videos/:id/hide
	adminVideos.Put("/:id/show", h.VideoHandler.ShowVideo)         // PUT /api/v1/admin/videos/:id/show
	adminVideos.Delete("/:id", h.VideoHandler.DeleteVideoByAdmin)  // DELETE /api/v1/admin/videos/:id
}
