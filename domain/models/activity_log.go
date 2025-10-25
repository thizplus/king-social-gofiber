package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityLog struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AdminID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Action       string    `gorm:"type:varchar(100);not null"` // delete_topic, ban_user, etc.
	ResourceType string    `gorm:"type:varchar(50);not null"`  // topic, user, video, etc.
	ResourceID   uuid.UUID `gorm:"type:uuid;not null"`
	Description  string    `gorm:"type:text"`
	IPAddress    string    `gorm:"type:varchar(45)"`
	UserAgent    string    `gorm:"type:text"`
	CreatedAt    time.Time

	// Relations
	Admin User `gorm:"foreignKey:AdminID"`
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}
