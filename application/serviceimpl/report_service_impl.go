package serviceimpl

import (
	"context"
	"github.com/google/uuid"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
)

type reportServiceImpl struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) services.ReportService {
	return &reportServiceImpl{reportRepo: reportRepo}
}

func (s *reportServiceImpl) CreateReport(ctx context.Context, userID uuid.UUID, req *dto.CreateReportRequest) error {
	report := &models.Report{
		ReporterID:  userID,
		Type:        req.Type,
		ResourceID:  req.ResourceID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      models.ReportStatusPending,
	}

	return s.reportRepo.Create(ctx, report)
}
