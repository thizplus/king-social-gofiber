package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Title        string         `gorm:"type:varchar(200);not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	VideoURL     string         `gorm:"type:varchar(500);not null" json:"videoUrl"`
	ThumbnailURL string         `gorm:"type:varchar(500)" json:"thumbnailUrl"`
	Duration     int            `gorm:"type:int;default:0" json:"duration"` // Duration in seconds
	Width        int            `gorm:"type:int" json:"width"`
	Height       int            `gorm:"type:int" json:"height"`
	FileSize     int64          `gorm:"type:bigint" json:"fileSize"` // File size in bytes
	ViewCount    int            `gorm:"type:int;default:0" json:"viewCount"`
	LikeCount    int            `gorm:"type:int;default:0" json:"likeCount"`
	CommentCount int            `gorm:"type:int;default:0" json:"commentCount"`
	IsActive     bool           `gorm:"type:boolean;default:true" json:"isActive"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

// BeforeCreate hook
func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
