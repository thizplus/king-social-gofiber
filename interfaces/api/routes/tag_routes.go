package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTagRoutes(api fiber.Router, h *handlers.Handlers) {
	// Public routes
	tags := api.Group("/tags")
	tags.Get("/", h.TagHandler.GetTags)
	tags.Get("/search", h.TagHandler.SearchTags)
	tags.Get("/:id", h.TagHandler.GetTag)

	// Admin routes
	adminTags := api.Group("/admin/tags")
	adminTags.Use(middleware.Protected())
	adminTags.Use(middleware.AdminOnly())

	adminTags.Post("/", h.TagHandler.CreateTag)
	adminTags.Put("/:id", h.TagHandler.UpdateTag)
	adminTags.Delete("/:id", h.TagHandler.DeleteTag)
}