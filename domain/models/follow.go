package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index" json:"followerId"`  // คนที่กดติดตาม
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index" json:"followingId"` // คนที่ถูกติดตาม
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`

	// Relations
	Follower  User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID;constraint:OnDelete:CASCADE" json:"following,omitempty"`
}

func (Follow) TableName() string {
	return "follows"
}
