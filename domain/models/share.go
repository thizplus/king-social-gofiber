package models

import (
	"time"

	"github.com/google/uuid"
)

type Share struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	VideoID   uuid.UUID `gorm:"type:uuid;not null;index" json:"videoId"`
	Platform  string    `gorm:"type:varchar(50)" json:"platform"` // "facebook", "twitter", "line", "copy_link"
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`

	// Relations
	User  *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Video *Video `gorm:"foreignKey:VideoID;constraint:OnDelete:CASCADE" json:"video,omitempty"`
}

func (Share) TableName() string {
	return "shares"
}
