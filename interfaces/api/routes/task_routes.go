package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTaskRoutes(api fiber.Router, h *handlers.Handlers) {
	tasks := api.Group("/tasks")
	tasks.Use(middleware.Protected())
	tasks.Post("/", h.TaskHandler.CreateTask)
	tasks.Get("/", middleware.AdminOnly(), h.TaskHandler.ListTasks)
	tasks.Get("/my", h.TaskHandler.GetUserTasks)
	tasks.Get("/:id", h.TaskHandler.GetTask)
	tasks.Put("/:id", middleware.OwnerOnly(), h.TaskHandler.UpdateTask)
	tasks.Delete("/:id", middleware.OwnerOnly(), h.TaskHandler.DeleteTask)
}
