package dto

import (
	"time"

	"gofiber-social/domain/models"

	"github.com/google/uuid"
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
	UserGrowth      []DataPoint  `json:"userGrowth"`      // Last 30 days
	ContentActivity []DataPoint  `json:"contentActivity"` // Last 30 days
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
	ID             uuid.UUID  `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	FullName       string     `json:"fullName"`
	Role           string     `json:"role"`
	Avatar         string     `json:"avatar,omitempty"`
	IsActive       bool       `json:"isActive"`
	FollowerCount  int        `json:"followerCount"`
	FollowingCount int        `json:"followingCount"`
	TopicCount     int        `json:"topicCount"`
	VideoCount     int        `json:"videoCount"`
	CreatedAt      time.Time  `json:"createdAt"`
	LastLoginAt    *time.Time `json:"lastLoginAt,omitempty"`
	SuspendedUntil *time.Time `json:"suspendedUntil,omitempty"`
	SuspendReason  string     `json:"suspendReason,omitempty"`
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
	ResourceID  uuid.UUID         `json:"resourceId" validate:"required"`
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
