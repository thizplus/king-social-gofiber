package routes

import (
	"gofiber-social/interfaces/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, h *handlers.Handlers) {
	auth := api.Group("/auth")
	auth.Post("/register", h.UserHandler.Register)
	auth.Post("/login", h.UserHandler.Login)
}
