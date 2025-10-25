package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(api fiber.Router, h *handlers.Handlers) {
	files := api.Group("/files")
	files.Use(middleware.Protected())
	files.Post("/upload", h.FileHandler.UploadFile)
	files.Get("/", middleware.AdminOnly(), h.FileHandler.ListFiles)
	files.Get("/my", h.FileHandler.GetUserFiles)
	files.Get("/:id", h.FileHandler.GetFile)
	files.Delete("/:id", middleware.OwnerOnly(), h.FileHandler.DeleteFile)
}
