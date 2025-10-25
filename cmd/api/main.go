package main

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"gofiber-social/interfaces/api/routes"
	"gofiber-social/pkg/di"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize DI container
	container := di.NewContainer()

	// Initialize all dependencies
	if err := container.Initialize(); err != nil {
		log.Fatal("Failed to initialize container:", err)
	}

	// Setup graceful shutdown
	setupGracefulShutdown(container)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(),
		AppName:      container.GetConfig().App.Name,
	})

	// Setup middleware
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

	// Create handlers from services
	services := container.GetHandlerServices()
	h := handlers.NewHandlers(services)

	// Setup routes
	routes.SetupRoutes(app, h)

	// Start server
	port := container.GetConfig().App.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌍 Environment: %s", container.GetConfig().App.Env)
	log.Printf("📚 Health check: http://localhost:%s/health", port)
	log.Printf("📖 API docs: http://localhost:%s/api/v1", port)
	log.Printf("🔌 WebSocket: ws://localhost:%s/ws", port)

	log.Fatal(app.Listen(":" + port))
}

func setupGracefulShutdown(container *di.Container) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\n🛑 Gracefully shutting down...")

		if err := container.Cleanup(); err != nil {
			log.Printf("❌ Error during cleanup: %v", err)
		}

		log.Println("👋 Shutdown complete")
		os.Exit(0)
	}()
}
