package routes

import (
	"gofiber-social/interfaces/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	// Setup health and root routes
	SetupHealthRoutes(app)

	// API version group
	api := app.Group("/api/v1")

	// Setup all route groups
	SetupAuthRoutes(api, h)
	SetupFollowRoutes(api, h) // Must be before SetupUserRoutes to avoid middleware conflict
	SetupUserRoutes(api, h)
	SetupTaskRoutes(api, h)
	SetupFileRoutes(api, h)
	SetupJobRoutes(api, h)
	SetupForumRoutes(api, h)
	SetupTopicRoutes(api, h)
	SetupReplyRoutes(api, h)
	SetupTagRoutes(api, h)
	SetupVideoRoutes(api, h)
	SetupLikeRoutes(api, h)
	SetupCommentRoutes(api, h)
	SetupShareRoutes(api, h)
	SetupNotificationRoutes(api, h)
	SetupAdminRoutes(api, h)
	SetupReportRoutes(api, h)

	// Setup WebSocket routes (needs app, not api group)
	SetupWebSocketRoutes(app)
}
