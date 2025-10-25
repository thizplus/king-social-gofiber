package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Slug        string    `gorm:"type:varchar(50);not null;uniqueIndex"`
	Description string    `gorm:"type:varchar(255)"`
	Color       string    `gorm:"type:varchar(7);default:'#3B82F6'"` // Hex color code
	IsActive    bool      `gorm:"default:true"`
	UsageCount  int       `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"` // Soft delete

	// Relations
	Topics []Topic `gorm:"many2many:topic_tags;"`
}

func (Tag) TableName() string {
	return "tags"
}