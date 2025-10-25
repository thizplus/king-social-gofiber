package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// Dashboard
// GET /api/v1/admin/dashboard/stats
func (h *AdminHandler) GetDashboardStats(c *fiber.Ctx) error {
	stats, err := h.adminService.GetDashboardStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve dashboard stats", err)
	}

	return utils.SuccessResponse(c, "Dashboard stats retrieved", stats)
}

// GET /api/v1/admin/dashboard/charts
func (h *AdminHandler) GetDashboardCharts(c *fiber.Ctx) error {
	charts, err := h.adminService.GetDashboardCharts(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve dashboard charts", err)
	}

	return utils.SuccessResponse(c, "Dashboard charts retrieved", charts)
}

// User Management
// GET /api/v1/admin/users
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	var params dto.AdminUserListRequest
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", err)
	}

	users, err := h.adminService.GetUsers(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve users", err)
	}

	return utils.SuccessResponse(c, "Users retrieved successfully", users)
}

// GET /api/v1/admin/users/:id
func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	user, err := h.adminService.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve user", err)
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}

// PUT /api/v1/admin/users/:id/suspend
func (h *AdminHandler) SuspendUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	var req dto.SuspendUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	if err := h.adminService.SuspendUser(c.Context(), adminID, userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to suspend user", err)
	}

	return utils.SuccessResponse(c, "User suspended successfully", nil)
}

// PUT /api/v1/admin/users/:id/activate
func (h *AdminHandler) ActivateUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	if err := h.adminService.ActivateUser(c.Context(), adminID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to activate user", err)
	}

	return utils.SuccessResponse(c, "User activated successfully", nil)
}

// DELETE /api/v1/admin/users/:id
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	if err := h.adminService.DeleteUser(c.Context(), adminID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete user", err)
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

// PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	var req dto.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	if err := h.adminService.UpdateUserRole(c.Context(), adminID, userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update user role", err)
	}

	return utils.SuccessResponse(c, "User role updated successfully", nil)
}

// Report Management
// GET /api/v1/admin/reports
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	var params dto.ReportQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters", err)
	}

	reports, err := h.adminService.GetReports(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve reports", err)
	}

	return utils.SuccessResponse(c, "Reports retrieved successfully", reports)
}

// GET /api/v1/admin/reports/:id
func (h *AdminHandler) GetReportByID(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid report ID", err)
	}

	report, err := h.adminService.GetReportByID(c.Context(), reportID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve report", err)
	}

	return utils.SuccessResponse(c, "Report retrieved successfully", report)
}

// PUT /api/v1/admin/reports/:id/review
func (h *AdminHandler) ReviewReport(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid report ID", err)
	}

	var req dto.ReviewReportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	if err := h.adminService.ReviewReport(c.Context(), adminID, reportID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to review report", err)
	}

	return utils.SuccessResponse(c, "Report reviewed successfully", nil)
}

// Activity Logs
// GET /api/v1/admin/activity-logs
func (h *AdminHandler) GetActivityLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	logs, err := h.adminService.GetActivityLogs(c.Context(), page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve activity logs", err)
	}

	return utils.SuccessResponse(c, "Activity logs retrieved successfully", logs)
}
