package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	FirstName string
	LastName  string
	FullName  string `gorm:"type:varchar(200)"`
	Avatar    string
	Bio       string `gorm:"type:text"`
	Website   string `gorm:"type:varchar(255)"`
	Role      string `gorm:"default:'user'"`
	IsActive  bool   `gorm:"default:true"`
	IsVerified bool  `gorm:"default:false"`
	IsPrivate  bool  `gorm:"default:false"`

	// Follow System (Task 04)
	FollowerCount  int `gorm:"default:0"`
	FollowingCount int `gorm:"default:0"`

	// Admin System (Task 06)
	SuspendedUntil *time.Time
	SuspendReason  string `gorm:"type:text"`
	LastLoginAt    *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (User) TableName() string {
	return "users"
}
