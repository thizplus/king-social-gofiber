package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	TopicID   *uuid.UUID `gorm:"type:uuid;index" json:"topicId,omitempty"`   // Nullable - for topic likes
	VideoID   *uuid.UUID `gorm:"type:uuid;index" json:"videoId,omitempty"`   // Nullable - for video likes
	ReplyID   *uuid.UUID `gorm:"type:uuid;index" json:"replyId,omitempty"`   // Nullable - for reply likes
	CommentID *uuid.UUID `gorm:"type:uuid;index" json:"commentId,omitempty"` // Nullable - for comment likes
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`

	// Relations
	User    *User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Topic   *Topic   `gorm:"foreignKey:TopicID;constraint:OnDelete:CASCADE" json:"topic,omitempty"`
	Video   *Video   `gorm:"foreignKey:VideoID;constraint:OnDelete:CASCADE" json:"video,omitempty"`
	Reply   *Reply   `gorm:"foreignKey:ReplyID;constraint:OnDelete:CASCADE" json:"reply,omitempty"`
	Comment *Comment `gorm:"foreignKey:CommentID;constraint:OnDelete:CASCADE" json:"comment,omitempty"`
}

func (Like) TableName() string {
	return "likes"
}

// Validate ensures exactly one of TopicID, VideoID, ReplyID, or CommentID is set
func (l *Like) Validate() error {
	count := 0
	if l.TopicID != nil {
		count++
	}
	if l.VideoID != nil {
		count++
	}
	if l.ReplyID != nil {
		count++
	}
	if l.CommentID != nil {
		count++
	}

	if count == 0 {
		return errors.New("one of TopicID, VideoID, ReplyID, or CommentID must be set")
	}
	if count > 1 {
		return errors.New("exactly one of TopicID, VideoID, ReplyID, or CommentID must be set")
	}
	return nil
}
