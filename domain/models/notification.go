package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeTopicReply   NotificationType = "topic_reply"    // มีคนตอบกระทู้
	NotificationTypeTopicLike    NotificationType = "topic_like"     // มีคนไลค์กระทู้
	NotificationTypeVideoLike    NotificationType = "video_like"     // มีคนไลค์วิดีโอ
	NotificationTypeVideoComment NotificationType = "video_comment"  // มีคนคอมเมนต์วิดีโอ
	NotificationTypeCommentReply NotificationType = "comment_reply"  // มีคนตอบกลับคอมเมนต์
	NotificationTypeReplyLike    NotificationType = "reply_like"     // มีคนไลค์การตอบกลับ
	NotificationTypeCommentLike  NotificationType = "comment_like"   // มีคนไลค์ความคิดเห็น
	NotificationTypeNewFollower  NotificationType = "new_follower"   // มีคนติดตาม
)

type Notification struct {
	ID         uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null;index"` // ผู้รับการแจ้งเตือน
	ActorID    uuid.UUID        `gorm:"type:uuid;not null"`       // ผู้กระทำ (คนที่ไลค์, คอมเมนต์, ติดตาม)
	Type       NotificationType `gorm:"type:varchar(50);not null;index"`
	ResourceID *uuid.UUID       `gorm:"type:uuid"` // ID ของ resource (topic_id, video_id, comment_id)
	Message    string           `gorm:"type:text"`
	IsRead     bool             `gorm:"type:boolean;default:false;index"`
	CreatedAt  time.Time

	// Relations
	User  User `gorm:"foreignKey:UserID"`
	Actor User `gorm:"foreignKey:ActorID"`
}

func (Notification) TableName() string {
	return "notifications"
}
