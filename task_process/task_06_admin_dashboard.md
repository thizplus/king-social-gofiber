# Task 06: Admin Dashboard & Management

## 📋 ภาพรวม
ระบบจัดการสำหรับ Admin ครอบคลุม Dashboard, User Management, Content Moderation, และ Statistics

## 🎯 ความสำคัญ
⭐⭐⭐ **สำคัญมาก - ระบบจัดการและควบคุมทั้งหมด**

## ⏱️ ระยะเวลา
**3-4 วัน**

## 📦 Dependencies
- ✅ All previous tasks (Task 00-05)

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Models (domain/models/)

#### 1.1 สร้าง `report.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type ReportType string
type ReportStatus string

const (
	ReportTypeTopic   ReportType = "topic"
	ReportTypeReply   ReportType = "reply"
	ReportTypeVideo   ReportType = "video"
	ReportTypeComment ReportType = "comment"
	ReportTypeUser    ReportType = "user"
)

const (
	ReportStatusPending   ReportStatus = "pending"
	ReportStatusReviewing ReportStatus = "reviewing"
	ReportStatusResolved  ReportStatus = "resolved"
	ReportStatusRejected  ReportStatus = "rejected"
)

type Report struct {
	ID          uuid.UUID    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ReporterID  uuid.UUID    `gorm:"type:uuid;not null;index"` // ผู้รายงาน
	Type        ReportType   `gorm:"type:varchar(50);not null;index"`
	ResourceID  uuid.UUID    `gorm:"type:uuid;not null;index"` // ID ของสิ่งที่ถูกรายงาน
	Reason      string       `gorm:"type:varchar(100);not null"` // spam, inappropriate, harassment, etc.
	Description string       `gorm:"type:text"`
	Status      ReportStatus `gorm:"type:varchar(50);default:'pending';index"`
	ReviewedBy  *uuid.UUID   `gorm:"type:uuid;index"` // Admin ที่ review
	ReviewNote  string       `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	Reporter User  `gorm:"foreignKey:ReporterID"`
	Reviewer *User `gorm:"foreignKey:ReviewedBy"`
}

func (Report) TableName() string {
	return "reports"
}
```

#### 1.2 สร้าง `activity_log.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type ActivityLog struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AdminID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Action      string    `gorm:"type:varchar(100);not null"` // delete_topic, ban_user, etc.
	ResourceType string   `gorm:"type:varchar(50);not null"`  // topic, user, video, etc.
	ResourceID  uuid.UUID `gorm:"type:uuid;not null"`
	Description string    `gorm:"type:text"`
	IPAddress   string    `gorm:"type:varchar(45)"`
	UserAgent   string    `gorm:"type:text"`
	CreatedAt   time.Time

	// Relations
	Admin User `gorm:"foreignKey:AdminID"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `admin.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

// Dashboard DTOs
type DashboardStatsResponse struct {
	TotalUsers    int64 `json:"totalUsers"`
	TotalTopics   int64 `json:"totalTopics"`
	TotalVideos   int64 `json:"totalVideos"`
	TotalComments int64 `json:"totalComments"`
	TotalReports  int64 `json:"totalReports"`
	NewUsersToday int64 `json:"newUsersToday"`
	ActiveUsers   int64 `json:"activeUsers"` // Users active in last 24h
}

type DashboardChartsResponse struct {
	UserGrowth      []DataPoint `json:"userGrowth"`      // Last 30 days
	ContentActivity []DataPoint `json:"contentActivity"` // Last 30 days
	TopForums       []ForumStats `json:"topForums"`
}

type DataPoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

type ForumStats struct {
	ForumID    uuid.UUID `json:"forumId"`
	ForumName  string    `json:"forumName"`
	TopicCount int64     `json:"topicCount"`
	ReplyCount int64     `json:"replyCount"`
}

// User Management DTOs
type AdminUserListRequest struct {
	Search   string `query:"search"`
	Role     string `query:"role" validate:"omitempty,oneof=user admin"`
	IsActive *bool  `query:"isActive"`
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type AdminUserResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	FullName       string    `json:"fullName"`
	Role           string    `json:"role"`
	Avatar         string    `json:"avatar,omitempty"`
	IsActive       bool      `json:"isActive"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	TopicCount     int       `json:"topicCount"`
	VideoCount     int       `json:"videoCount"`
	CreatedAt      time.Time `json:"createdAt"`
	LastLoginAt    *time.Time `json:"lastLoginAt,omitempty"`
}

type AdminUserListResponse struct {
	Users      []AdminUserResponse `json:"users"`
	TotalCount int64               `json:"totalCount"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"totalPages"`
}

type SuspendUserRequest struct {
	Reason   string `json:"reason" validate:"required,min=10,max=500"`
	Duration int    `json:"duration" validate:"required,min=1"` // days
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}

// Report Management DTOs
type ReportQueryParams struct {
	Type   string `query:"type" validate:"omitempty,oneof=topic reply video comment user"`
	Status string `query:"status" validate:"omitempty,oneof=pending reviewing resolved rejected"`
	Page   int    `query:"page" validate:"omitempty,min=1"`
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

type CreateReportRequest struct {
	Type        models.ReportType `json:"type" validate:"required,oneof=topic reply video comment user"`
	ResourceID  uuid.UUID         `json:"resourceId" validate:"required,uuid"`
	Reason      string            `json:"reason" validate:"required,oneof=spam inappropriate harassment misinformation other"`
	Description string            `json:"description" validate:"required,min=10,max=1000"`
}

type ReviewReportRequest struct {
	Status     models.ReportStatus `json:"status" validate:"required,oneof=resolved rejected"`
	ReviewNote string              `json:"reviewNote" validate:"required,min=10,max=500"`
}

type ReportResponse struct {
	ID          uuid.UUID           `json:"id"`
	Reporter    UserSummary         `json:"reporter"`
	Type        models.ReportType   `json:"type"`
	ResourceID  uuid.UUID           `json:"resourceId"`
	Reason      string              `json:"reason"`
	Description string              `json:"description"`
	Status      models.ReportStatus `json:"status"`
	Reviewer    *UserSummary        `json:"reviewer,omitempty"`
	ReviewNote  string              `json:"reviewNote,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

type ReportListResponse struct {
	Reports    []ReportResponse `json:"reports"`
	TotalCount int64            `json:"totalCount"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"totalPages"`
}

// Activity Log DTOs
type ActivityLogResponse struct {
	ID           uuid.UUID   `json:"id"`
	Admin        UserSummary `json:"admin"`
	Action       string      `json:"action"`
	ResourceType string      `json:"resourceType"`
	ResourceID   uuid.UUID   `json:"resourceId"`
	Description  string      `json:"description"`
	IPAddress    string      `json:"ipAddress"`
	CreatedAt    time.Time   `json:"createdAt"`
}

type ActivityLogListResponse struct {
	Logs       []ActivityLogResponse `json:"logs"`
	TotalCount int64                 `json:"totalCount"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"totalPages"`
}
```

---

### 3. Repository Interfaces (domain/repositories/)

#### 3.1 สร้าง `report_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
)

type ReportRepository interface {
	Create(ctx context.Context, report *models.Report) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Report, error)
	FindAll(ctx context.Context, params *dto.ReportQueryParams) ([]models.Report, int64, error)
	Update(ctx context.Context, report *models.Report) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCount(ctx context.Context) (int64, error)
	GetPendingCount(ctx context.Context) (int64, error)
}
```

#### 3.2 สร้าง `activity_log_repository.go`
```go
package repositories

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/models"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
	FindByAdminID(ctx context.Context, adminID uuid.UUID, page, limit int) ([]models.ActivityLog, int64, error)
	FindAll(ctx context.Context, page, limit int) ([]models.ActivityLog, int64, error)
}
```

#### 3.3 เพิ่ม Admin methods ใน existing repositories

**user_repository.go:**
```go
// Add to UserRepository interface
FindAllWithStats(ctx context.Context, params *dto.AdminUserListRequest) ([]models.User, int64, error)
GetTotalCount(ctx context.Context) (int64, error)
GetNewUsersToday(ctx context.Context) (int64, error)
GetActiveUsersCount(ctx context.Context) (int64, error)
SuspendUser(ctx context.Context, userID uuid.UUID, reason string, until time.Time) error
UpdateRole(ctx context.Context, userID uuid.UUID, role string) error
```

**topic_repository.go:**
```go
// Add to TopicRepository interface
GetTotalCount(ctx context.Context) (int64, error)
```

**video_repository.go:**
```go
// Add to VideoRepository interface
GetTotalCount(ctx context.Context) (int64, error)
```

**comment_repository.go:**
```go
// Add to CommentRepository interface
GetTotalCount(ctx context.Context) (int64, error)
```

---

### 4. Repository Implementations (infrastructure/postgres/)

#### 4.1 สร้าง `report_repository_impl.go`
```go
package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repositories.ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.WithContext(ctx).
		Preload("Reporter").
		Preload("Reviewer").
		First(&report, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepositoryImpl) FindAll(ctx context.Context, params *dto.ReportQueryParams) ([]models.Report, int64, error) {
	var reports []models.Report
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Report{}).
		Preload("Reporter").
		Preload("Reviewer")

	// Filter by type
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// Filter by status
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reports).Error

	return reports, totalCount, err
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Report{}, id).Error
}

func (r *reportRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Report{}).Count(&count).Error
	return count, err
}

func (r *reportRepositoryImpl) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Report{}).
		Where("status = ?", models.ReportStatusPending).
		Count(&count).Error
	return count, err
}
```

#### 4.2 สร้าง `activity_log_repository_impl.go`
```go
package postgres

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
)

type activityLogRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) repositories.ActivityLogRepository {
	return &activityLogRepositoryImpl{db: db}
}

func (r *activityLogRepositoryImpl) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *activityLogRepositoryImpl) FindByAdminID(ctx context.Context, adminID uuid.UUID, page, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{}).
		Where("admin_id = ?", adminID).
		Preload("Admin")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, totalCount, err
}

func (r *activityLogRepositoryImpl) FindAll(ctx context.Context, page, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{}).
		Preload("Admin")

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error

	return logs, totalCount, err
}
```

---

### 5. Service Interface (domain/services/)

#### 5.1 สร้าง `admin_service.go`
```go
package services

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
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
```

---

### 6. Service Implementation (application/serviceimpl/)

#### 6.1 สร้าง `admin_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
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
			Reporter:    dto.UserSummary{ID: report.Reporter.ID, Username: report.Reporter.Username, FullName: report.Reporter.FullName},
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
				ID:       report.Reviewer.ID,
				Username: report.Reviewer.Username,
				FullName: report.Reviewer.FullName,
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
		Reporter:    dto.UserSummary{ID: report.Reporter.ID, Username: report.Reporter.Username, FullName: report.Reporter.FullName},
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
			ID:       report.Reviewer.ID,
			Username: report.Reviewer.Username,
			FullName: report.Reviewer.FullName,
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
			Admin:        dto.UserSummary{ID: log.Admin.ID, Username: log.Admin.Username, FullName: log.Admin.FullName},
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
```

#### 6.2 สร้าง `report_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/models"
	"yourproject/domain/repositories"
	"yourproject/domain/services"
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
```

---

### 7. Handlers (interfaces/api/handlers/)

#### 7.1 สร้าง `admin_handler.go`
```go
package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Dashboard stats retrieved", stats)
}

// GET /api/v1/admin/dashboard/charts
func (h *AdminHandler) GetDashboardCharts(c *fiber.Ctx) error {
	charts, err := h.adminService.GetDashboardCharts(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Dashboard charts retrieved", charts)
}

// User Management
// GET /api/v1/admin/users
func (h *AdminHandler) GetUsers(c *fiber.Ctx) error {
	var params dto.AdminUserListRequest
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	users, err := h.adminService.GetUsers(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Users retrieved successfully", users)
}

// GET /api/v1/admin/users/:id
func (h *AdminHandler) GetUserByID(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.adminService.GetUserByID(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User retrieved successfully", user)
}

// PUT /api/v1/admin/users/:id/suspend
func (h *AdminHandler) SuspendUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req dto.SuspendUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.SuspendUser(c.Context(), adminID, userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User suspended successfully", nil)
}

// PUT /api/v1/admin/users/:id/activate
func (h *AdminHandler) ActivateUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.adminService.ActivateUser(c.Context(), adminID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User activated successfully", nil)
}

// DELETE /api/v1/admin/users/:id
func (h *AdminHandler) DeleteUser(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	if err := h.adminService.DeleteUser(c.Context(), adminID, userID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User deleted successfully", nil)
}

// PUT /api/v1/admin/users/:id/role
func (h *AdminHandler) UpdateUserRole(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	var req dto.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.UpdateUserRole(c.Context(), adminID, userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User role updated successfully", nil)
}

// Report Management
// GET /api/v1/admin/reports
func (h *AdminHandler) GetReports(c *fiber.Ctx) error {
	var params dto.ReportQueryParams
	if err := c.QueryParser(&params); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid query parameters")
	}

	reports, err := h.adminService.GetReports(c.Context(), &params)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Reports retrieved successfully", reports)
}

// GET /api/v1/admin/reports/:id
func (h *AdminHandler) GetReportByID(c *fiber.Ctx) error {
	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid report ID")
	}

	report, err := h.adminService.GetReportByID(c.Context(), reportID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Report retrieved successfully", report)
}

// PUT /api/v1/admin/reports/:id/review
func (h *AdminHandler) ReviewReport(c *fiber.Ctx) error {
	adminID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	reportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid report ID")
	}

	var req dto.ReviewReportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.ReviewReport(c.Context(), adminID, reportID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Report reviewed successfully", nil)
}

// Activity Logs
// GET /api/v1/admin/activity-logs
func (h *AdminHandler) GetActivityLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	logs, err := h.adminService.GetActivityLogs(c.Context(), page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Activity logs retrieved successfully", logs)
}
```

#### 7.2 สร้าง `report_handler.go`
```go
package handlers

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/domain/dto"
	"yourproject/domain/services"
	"yourproject/pkg/utils"
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
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	var req dto.CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.reportService.CreateReport(c.Context(), userID, &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Report submitted successfully", nil)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `admin_routes.go`
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"yourproject/interfaces/api/handlers"
	"yourproject/interfaces/api/middleware"
)

func SetupAdminRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, reportHandler *handlers.ReportHandler, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/v1")

	// User report routes (protected)
	auth := api.Use(authMiddleware.Protected())
	auth.Post("/reports", reportHandler.CreateReport) // POST /api/v1/reports

	// Admin routes
	admin := api.Group("/admin", authMiddleware.Protected(), authMiddleware.RequireRole("admin"))

	// Dashboard
	admin.Get("/dashboard/stats", adminHandler.GetDashboardStats)     // GET /api/v1/admin/dashboard/stats
	admin.Get("/dashboard/charts", adminHandler.GetDashboardCharts)   // GET /api/v1/admin/dashboard/charts

	// User Management
	users := admin.Group("/users")
	users.Get("/", adminHandler.GetUsers)                  // GET /api/v1/admin/users
	users.Get("/:id", adminHandler.GetUserByID)            // GET /api/v1/admin/users/:id
	users.Put("/:id/suspend", adminHandler.SuspendUser)    // PUT /api/v1/admin/users/:id/suspend
	users.Put("/:id/activate", adminHandler.ActivateUser)  // PUT /api/v1/admin/users/:id/activate
	users.Delete("/:id", adminHandler.DeleteUser)          // DELETE /api/v1/admin/users/:id
	users.Put("/:id/role", adminHandler.UpdateUserRole)    // PUT /api/v1/admin/users/:id/role

	// Report Management
	reports := admin.Group("/reports")
	reports.Get("/", adminHandler.GetReports)              // GET /api/v1/admin/reports
	reports.Get("/:id", adminHandler.GetReportByID)        // GET /api/v1/admin/reports/:id
	reports.Put("/:id/review", adminHandler.ReviewReport)  // PUT /api/v1/admin/reports/:id/review

	// Activity Logs
	admin.Get("/activity-logs", adminHandler.GetActivityLogs) // GET /api/v1/admin/activity-logs
}
```

---

### 9. Container Updates (`pkg/di/container.go`)

```go
// Add to container.go

func (c *Container) InitializeAdminComponents() {
	// Repositories
	c.ReportRepository = postgres.NewReportRepository(c.DB)
	c.ActivityLogRepository = postgres.NewActivityLogRepository(c.DB)

	// Services
	c.AdminService = serviceimpl.NewAdminService(
		c.UserRepository,
		c.TopicRepository,
		c.VideoRepository,
		c.CommentRepository,
		c.ReportRepository,
		c.ActivityLogRepository,
		c.ForumRepository,
	)

	c.ReportService = serviceimpl.NewReportService(c.ReportRepository)

	// Handlers
	c.AdminHandler = handlers.NewAdminHandler(c.AdminService)
	c.ReportHandler = handlers.NewReportHandler(c.ReportService)
}

// Add to Container struct
type Container struct {
	// ... existing fields

	// Admin
	ReportRepository      repositories.ReportRepository
	ActivityLogRepository repositories.ActivityLogRepository
	AdminService          services.AdminService
	ReportService         services.ReportService
	AdminHandler          *handlers.AdminHandler
	ReportHandler         *handlers.ReportHandler
}
```

---

### 10. Main Updates (`cmd/api/main.go`)

```go
// Add to main.go

func main() {
	// ... existing code

	// Initialize admin components
	container.InitializeAdminComponents()

	// Setup routes
	routes.SetupAdminRoutes(app, container.AdminHandler, container.ReportHandler, container.AuthMiddleware)

	// ... rest of the code
}
```

---

### 11. Database Migrations

```sql
-- reports table
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    reason VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    review_note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_type ON reports(type);
CREATE INDEX idx_reports_resource_id ON reports(resource_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_reviewed_by ON reports(reviewed_by);

-- activity_logs table
CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    description TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activity_logs_admin_id ON activity_logs(admin_id);
CREATE INDEX idx_activity_logs_action ON activity_logs(action);
CREATE INDEX idx_activity_logs_resource_type ON activity_logs(resource_type);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at DESC);

-- Update users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS suspended_until TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS suspend_reason TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP;
```

---

## ✅ Checklist

### Models & Database
- [ ] สร้าง `Report` model
- [ ] สร้าง `ActivityLog` model
- [ ] เพิ่ม migrations
- [ ] Update User model (suspended_until, last_login_at)
- [ ] สร้าง indexes

### DTOs
- [ ] สร้าง Dashboard DTOs
- [ ] สร้าง User Management DTOs
- [ ] สร้าง Report DTOs
- [ ] สร้าง Activity Log DTOs

### Repositories
- [ ] สร้าง `ReportRepository`
- [ ] สร้าง `ActivityLogRepository`
- [ ] เพิ่ม Admin methods ใน existing repositories

### Services
- [ ] สร้าง `AdminService`
- [ ] สร้าง `ReportService`
- [ ] Implement ทุก methods

### Handlers
- [ ] สร้าง `AdminHandler`
- [ ] สร้าง `ReportHandler`

### Routes
- [ ] Setup admin routes
- [ ] Setup report routes

### Integration
- [ ] Register components ใน Container
- [ ] Update main.go
- [ ] ทดสอบ end-to-end

---

## 🧪 Testing Guide

### 1. Dashboard Stats
```bash
GET /api/v1/admin/dashboard/stats
Authorization: Bearer {admin-token}
```

### 2. Get Users
```bash
GET /api/v1/admin/users?page=1&limit=20&search=john&role=user
Authorization: Bearer {admin-token}
```

### 3. Suspend User
```bash
PUT /api/v1/admin/users/{user-id}/suspend
Authorization: Bearer {admin-token}

{
  "reason": "Violated community guidelines",
  "duration": 7
}
```

### 4. User Report Content
```bash
POST /api/v1/reports
Authorization: Bearer {token}

{
  "type": "video",
  "resourceId": "uuid",
  "reason": "spam",
  "description": "This video contains spam content"
}
```

### 5. Admin Review Report
```bash
PUT /api/v1/admin/reports/{report-id}/review
Authorization: Bearer {admin-token}

{
  "status": "resolved",
  "reviewNote": "Content removed and user warned"
}
```

### 6. Activity Logs
```bash
GET /api/v1/admin/activity-logs?page=1&limit=50
Authorization: Bearer {admin-token}
```

---

## 📝 Notes

### Admin Features:
- ✅ Dashboard with statistics
- ✅ User management (suspend, delete, role change)
- ✅ Report system & moderation
- ✅ Activity logging
- ✅ Content management

### Security:
- Admin-only endpoints (role middleware)
- Activity logging for audit trail
- Cannot suspend/delete other admins
- IP address and user agent tracking

### Report System:
- Users can report: topics, replies, videos, comments, users
- Report reasons: spam, inappropriate, harassment, misinformation, other
- Report workflow: pending → reviewing → resolved/rejected
- Admin can add review notes

### Activity Logs:
- Track all admin actions
- Store IP address and user agent
- Useful for audit and security

### Future Enhancements:
- [ ] Advanced analytics and charts
- [ ] Bulk operations
- [ ] Export data (CSV, PDF)
- [ ] Email notifications to users
- [ ] Automated moderation (AI)
- [ ] Custom admin roles with permissions
- [ ] Scheduled tasks (cleanup, reports)

---

**ระยะเวลาโดยประมาณ:** 3-4 วัน

**Dependencies:**
- All previous tasks (Task 00-05)

**🎉 นี่คือ Task สุดท้าย!** หลังจากทำ Task 06 เสร็จ ระบบจะครบถ้วนและพร้อมใช้งาน!
