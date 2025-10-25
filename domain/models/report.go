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
