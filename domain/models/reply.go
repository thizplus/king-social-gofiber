package models

import (
	"time"

	"github.com/google/uuid"
)

type Reply struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	TopicID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"` // สำหรับ nested reply
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"` // Soft delete

	// Relations
	Topic   Topic   `gorm:"foreignKey:TopicID"`
	User    User    `gorm:"foreignKey:UserID"`
	Parent  *Reply  `gorm:"foreignKey:ParentID"`
	Replies []Reply `gorm:"foreignKey:ParentID"` // Nested replies
}

func (Reply) TableName() string {
	return "replies"
}
