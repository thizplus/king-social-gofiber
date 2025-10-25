package routes

import (
	"gofiber-social/interfaces/api/middleware"
	websocketHandler "gofiber-social/interfaces/api/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupWebSocketRoutes(app *fiber.App) {
	wsHandler := websocketHandler.NewWebSocketHandler()

	// WebSocket with optional authentication
	app.Use("/ws", middleware.Optional(), wsHandler.WebSocketUpgrade)
	app.Get("/ws", websocket.New(wsHandler.HandleWebSocket))
}
