package models

import (
	"time"

	"github.com/google/uuid"
)

type Forum struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null"` // URL-friendly
	Description string    `gorm:"type:text;not null"`
	Icon        string    `gorm:"type:varchar(500)"` // URL ของไอคอน
	Order       int       `gorm:"default:0"`         // ลำดับการแสดงผล
	IsActive    bool      `gorm:"default:true"`
	TopicCount  int       `gorm:"default:0"`         // จำนวนกระทู้
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null"` // Admin ID
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	Admin User `gorm:"foreignKey:CreatedBy"`
}

func (Forum) TableName() string {
	return "forums"
}
