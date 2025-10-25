package models

import (
	"time"

	"github.com/google/uuid"
)

type TopicTag struct {
	TopicID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	TagID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time

	// Relations
	Topic Topic `gorm:"foreignKey:TopicID"`
	Tag   Tag   `gorm:"foreignKey:TagID"`
}

func (TopicTag) TableName() string {
	return "topic_tags"
}