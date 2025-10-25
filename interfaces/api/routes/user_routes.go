package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, h *handlers.Handlers) {
	// Protected user routes
	users := api.Group("/users", middleware.Protected())
	users.Get("/profile", h.UserHandler.GetProfile)
	users.Put("/profile", h.UserHandler.UpdateProfile)
	users.Delete("/profile", h.UserHandler.DeleteUser)
	users.Get("/", middleware.AdminOnly(), h.UserHandler.ListUsers)
}
