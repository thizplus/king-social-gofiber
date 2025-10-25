package services

import (
	"context"
	"github.com/google/uuid"
	"gofiber-social/domain/dto"
)

type AdminService interface {
	// Dashboard
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
	GetDashboardCharts(ctx context.Context) (*dto.DashboardChartsResponse, error)

	// User Management
	GetUsers(ctx context.Context, params *dto.AdminUserListRequest) (*dto.AdminUserListResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.AdminUserResponse, error)
	SuspendUser(ctx context.Context, adminID, userID uuid.UUID, req *dto.SuspendUserRequest) error
	ActivateUser(ctx context.Context, adminID, userID uuid.UUID) error
	DeleteUser(ctx context.Context, adminID, userID uuid.UUID) error
	UpdateUserRole(ctx context.Context, adminID, userID uuid.UUID, req *dto.UpdateUserRoleRequest) error

	// Report Management
	GetReports(ctx context.Context, params *dto.ReportQueryParams) (*dto.ReportListResponse, error)
	GetReportByID(ctx context.Context, id uuid.UUID) (*dto.ReportResponse, error)
	ReviewReport(ctx context.Context, adminID, reportID uuid.UUID, req *dto.ReviewReportRequest) error

	// Activity Logs
	GetActivityLogs(ctx context.Context, page, limit int) (*dto.ActivityLogListResponse, error)
}

type ReportService interface {
	CreateReport(ctx context.Context, userID uuid.UUID, req *dto.CreateReportRequest) error
}
