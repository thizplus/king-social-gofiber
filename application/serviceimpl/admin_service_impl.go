package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	"github.com/google/uuid"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
)

type adminServiceImpl struct {
	userRepo        repositories.UserRepository
	topicRepo       repositories.TopicRepository
	videoRepo       repositories.VideoRepository
	commentRepo     repositories.CommentRepository
	reportRepo      repositories.ReportRepository
	activityLogRepo repositories.ActivityLogRepository
	forumRepo       repositories.ForumRepository
}

func NewAdminService(
	userRepo repositories.UserRepository,
	topicRepo repositories.TopicRepository,
	videoRepo repositories.VideoRepository,
	commentRepo repositories.CommentRepository,
	reportRepo repositories.ReportRepository,
	activityLogRepo repositories.ActivityLogRepository,
	forumRepo repositories.ForumRepository,
) services.AdminService {
	return &adminServiceImpl{
		userRepo:        userRepo,
		topicRepo:       topicRepo,
		videoRepo:       videoRepo,
		commentRepo:     commentRepo,
		reportRepo:      reportRepo,
		activityLogRepo: activityLogRepo,
		forumRepo:       forumRepo,
	}
}

// Dashboard
func (s *adminServiceImpl) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	totalUsers, _ := s.userRepo.GetTotalCount(ctx)
	totalTopics, _ := s.topicRepo.GetTotalCount(ctx)
	totalVideos, _ := s.videoRepo.GetTotalCount(ctx)
	totalComments, _ := s.commentRepo.GetTotalCount(ctx)
	totalReports, _ := s.reportRepo.GetTotalCount(ctx)
	newUsersToday, _ := s.userRepo.GetNewUsersToday(ctx)
	activeUsers, _ := s.userRepo.GetActiveUsersCount(ctx)

	return &dto.DashboardStatsResponse{
		TotalUsers:    totalUsers,
		TotalTopics:   totalTopics,
		TotalVideos:   totalVideos,
		TotalComments: totalComments,
		TotalReports:  totalReports,
		NewUsersToday: newUsersToday,
		ActiveUsers:   activeUsers,
	}, nil
}

func (s *adminServiceImpl) GetDashboardCharts(ctx context.Context) (*dto.DashboardChartsResponse, error) {
	// TODO: Implement charts data
	// - User growth (last 30 days)
	// - Content activity (last 30 days)
	// - Top forums by activity

	return &dto.DashboardChartsResponse{
		UserGrowth:      []dto.DataPoint{},
		ContentActivity: []dto.DataPoint{},
		TopForums:       []dto.ForumStats{},
	}, nil
}

// User Management
func (s *adminServiceImpl) GetUsers(ctx context.Context, params *dto.AdminUserListRequest) (*dto.AdminUserListResponse, error) {
	users, totalCount, err := s.userRepo.FindAllWithStats(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert to response
	userResponses := make([]dto.AdminUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.AdminUserResponse{
			ID:             user.ID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			Role:           user.Role,
			Avatar:         user.Avatar,
			IsActive:       user.IsActive,
			FollowerCount:  user.FollowerCount,
			FollowingCount: user.FollowingCount,
			CreatedAt:      user.CreatedAt,
			LastLoginAt:    user.LastLoginAt,
		}
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.AdminUserListResponse{
		Users:      userResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *adminServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.AdminUserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.AdminUserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		FullName:       user.FullName,
		Role:           user.Role,
		Avatar:         user.Avatar,
		IsActive:       user.IsActive,
		FollowerCount:  user.FollowerCount,
		FollowingCount: user.FollowingCount,
		CreatedAt:      user.CreatedAt,
		LastLoginAt:    user.LastLoginAt,
	}, nil
}

func (s *adminServiceImpl) SuspendUser(ctx context.Context, adminID, userID uuid.UUID, req *dto.SuspendUserRequest) error {
	// Don't allow suspending admins
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Role == "admin" {
		return errors.New("cannot suspend admin users")
	}

	until := time.Now().Add(time.Duration(req.Duration) * 24 * time.Hour)
	if err := s.userRepo.SuspendUser(ctx, userID, req.Reason, until); err != nil {
		return err
	}

	// Log activity
	s.logActivity(ctx, adminID, "suspend_user", "user", userID, fmt.Sprintf("Suspended for %d days: %s", req.Duration, req.Reason))

	return nil
}

func (s *adminServiceImpl) ActivateUser(ctx context.Context, adminID, userID uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.IsActive = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Log activity
	s.logActivity(ctx, adminID, "activate_user", "user", userID, "User activated")

	return nil
}

func (s *adminServiceImpl) DeleteUser(ctx context.Context, adminID, userID uuid.UUID) error {
	// Don't allow deleting admins
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Role == "admin" {
		return errors.New("cannot delete admin users")
	}

	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// Log activity
	s.logActivity(ctx, adminID, "delete_user", "user", userID, fmt.Sprintf("Deleted user: %s", user.Username))

	return nil
}

func (s *adminServiceImpl) UpdateUserRole(ctx context.Context, adminID, userID uuid.UUID, req *dto.UpdateUserRoleRequest) error {
	if err := s.userRepo.UpdateRole(ctx, userID, req.Role); err != nil {
		return err
	}

	// Log activity
	s.logActivity(ctx, adminID, "update_user_role", "user", userID, fmt.Sprintf("Changed role to: %s", req.Role))

	return nil
}

// Report Management
func (s *adminServiceImpl) GetReports(ctx context.Context, params *dto.ReportQueryParams) (*dto.ReportListResponse, error) {
	reports, totalCount, err := s.reportRepo.FindAll(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert to response
	reportResponses := make([]dto.ReportResponse, len(reports))
	for i, report := range reports {
		resp := dto.ReportResponse{
			ID:          report.ID,
			Reporter:    dto.UserSummary{ID: report.Reporter.ID, Username: report.Reporter.Username, FirstName: report.Reporter.FirstName, LastName: report.Reporter.LastName, Avatar: report.Reporter.Avatar},
			Type:        report.Type,
			ResourceID:  report.ResourceID,
			Reason:      report.Reason,
			Description: report.Description,
			Status:      report.Status,
			ReviewNote:  report.ReviewNote,
			CreatedAt:   report.CreatedAt,
			UpdatedAt:   report.UpdatedAt,
		}

		if report.Reviewer != nil {
			resp.Reviewer = &dto.UserSummary{
				ID:        report.Reviewer.ID,
				Username:  report.Reviewer.Username,
				FirstName: report.Reviewer.FirstName,
				LastName:  report.Reviewer.LastName,
				Avatar:    report.Reviewer.Avatar,
			}
		}

		reportResponses[i] = resp
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.ReportListResponse{
		Reports:    reportResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *adminServiceImpl) GetReportByID(ctx context.Context, id uuid.UUID) (*dto.ReportResponse, error) {
	report, err := s.reportRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &dto.ReportResponse{
		ID:          report.ID,
		Reporter:    dto.UserSummary{ID: report.Reporter.ID, Username: report.Reporter.Username, FirstName: report.Reporter.FirstName, LastName: report.Reporter.LastName, Avatar: report.Reporter.Avatar},
		Type:        report.Type,
		ResourceID:  report.ResourceID,
		Reason:      report.Reason,
		Description: report.Description,
		Status:      report.Status,
		ReviewNote:  report.ReviewNote,
		CreatedAt:   report.CreatedAt,
		UpdatedAt:   report.UpdatedAt,
	}

	if report.Reviewer != nil {
		resp.Reviewer = &dto.UserSummary{
			ID:        report.Reviewer.ID,
			Username:  report.Reviewer.Username,
			FirstName: report.Reviewer.FirstName,
			LastName:  report.Reviewer.LastName,
			Avatar:    report.Reviewer.Avatar,
		}
	}

	return resp, nil
}

func (s *adminServiceImpl) ReviewReport(ctx context.Context, adminID, reportID uuid.UUID, req *dto.ReviewReportRequest) error {
	report, err := s.reportRepo.FindByID(ctx, reportID)
	if err != nil {
		return err
	}

	report.Status = req.Status
	report.ReviewedBy = &adminID
	report.ReviewNote = req.ReviewNote

	if err := s.reportRepo.Update(ctx, report); err != nil {
		return err
	}

	// Log activity
	s.logActivity(ctx, adminID, "review_report", "report", reportID, fmt.Sprintf("Status: %s - %s", req.Status, req.ReviewNote))

	return nil
}

// Activity Logs
func (s *adminServiceImpl) GetActivityLogs(ctx context.Context, page, limit int) (*dto.ActivityLogListResponse, error) {
	logs, totalCount, err := s.activityLogRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert to response
	logResponses := make([]dto.ActivityLogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = dto.ActivityLogResponse{
			ID:           log.ID,
			Admin:        dto.UserSummary{ID: log.Admin.ID, Username: log.Admin.Username, FirstName: log.Admin.FirstName, LastName: log.Admin.LastName, Avatar: log.Admin.Avatar},
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Description:  log.Description,
			IPAddress:    log.IPAddress,
			CreatedAt:    log.CreatedAt,
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.ActivityLogListResponse{
		Logs:       logResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Helper
func (s *adminServiceImpl) logActivity(ctx context.Context, adminID uuid.UUID, action, resourceType string, resourceID uuid.UUID, description string) {
	log := &models.ActivityLog{
		AdminID:      adminID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Description:  description,
	}
	_ = s.activityLogRepo.Create(context.Background(), log)
}
