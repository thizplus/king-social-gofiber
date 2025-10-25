package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupShareRoutes(api fiber.Router, h *handlers.Handlers) {
	// Video share routes
	videos := api.Group("/videos")
	videos.Post("/:id/share", middleware.Protected(), h.ShareHandler.ShareVideo)  // POST /api/v1/videos/:id/share
	videos.Get("/:id/share/count", h.ShareHandler.GetShareCount)                  // GET /api/v1/videos/:id/share/count (public)
}
