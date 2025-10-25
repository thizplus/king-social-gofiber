package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// POST /api/v1/reports
func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized", err)
	}

	var req dto.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error(), err)
	}

	if err := h.reportService.CreateReport(c.Context(), userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error(), err)
	}

	return utils.SuccessResponse(c, "Report submitted successfully", nil)
}
