package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	ID         uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ForumID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title      string     `gorm:"type:varchar(200);not null"`
	Content    string     `gorm:"type:text;not null"`
	Thumbnail  string     `gorm:"type:varchar(500)"` // Optional thumbnail image URL
	ViewCount  int        `gorm:"default:0"`
	ReplyCount int        `gorm:"default:0"`
	LikeCount  int        `gorm:"default:0"`
	IsPinned   bool       `gorm:"default:false"`
	IsLocked   bool       `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"` // Soft delete

	// Relations
	Forum   Forum   `gorm:"foreignKey:ForumID"`
	User    User    `gorm:"foreignKey:UserID"`
	Replies []Reply `gorm:"foreignKey:TopicID"`
	Tags    []Tag   `gorm:"many2many:topic_tags;"`
}

func (Topic) TableName() string {
	return "topics"
}
