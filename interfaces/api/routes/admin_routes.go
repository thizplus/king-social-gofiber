package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
)

func SetupAdminRoutes(api fiber.Router, h *handlers.Handlers) {
	// Admin routes
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))

	// Dashboard
	admin.Get("/dashboard/stats", h.AdminHandler.GetDashboardStats)     // GET /api/v1/admin/dashboard/stats
	admin.Get("/dashboard/charts", h.AdminHandler.GetDashboardCharts)   // GET /api/v1/admin/dashboard/charts

	// User Management
	users := admin.Group("/users")
	users.Get("/", h.AdminHandler.GetUsers)                  // GET /api/v1/admin/users
	users.Get("/:id", h.AdminHandler.GetUserByID)            // GET /api/v1/admin/users/:id
	users.Put("/:id/suspend", h.AdminHandler.SuspendUser)    // PUT /api/v1/admin/users/:id/suspend
	users.Put("/:id/activate", h.AdminHandler.ActivateUser)  // PUT /api/v1/admin/users/:id/activate
	users.Delete("/:id", h.AdminHandler.DeleteUser)          // DELETE /api/v1/admin/users/:id
	users.Put("/:id/role", h.AdminHandler.UpdateUserRole)    // PUT /api/v1/admin/users/:id/role

	// Report Management
	reports := admin.Group("/reports")
	reports.Get("/", h.AdminHandler.GetReports)              // GET /api/v1/admin/reports
	reports.Get("/:id", h.AdminHandler.GetReportByID)        // GET /api/v1/admin/reports/:id
	reports.Put("/:id/review", h.AdminHandler.ReviewReport)  // PUT /api/v1/admin/reports/:id/review

	// Activity Logs
	admin.Get("/activity-logs", h.AdminHandler.GetActivityLogs) // GET /api/v1/admin/activity-logs
}

func SetupReportRoutes(api fiber.Router, h *handlers.Handlers) {
	// User report routes (protected)
	reports := api.Group("/reports", middleware.Protected())
	reports.Post("/", h.ReportHandler.CreateReport) // POST /api/v1/reports
}
