package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Username  string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	FirstName string `json:"firstName" validate:"required,min=1,max=50"`
	LastName  string `json:"lastName" validate:"required,min=1,max=50"`
}

type UpdateUserRequest struct {
	FirstName  string `json:"firstName" validate:"omitempty,min=1,max=50"`
	LastName   string `json:"lastName" validate:"omitempty,min=1,max=50"`
	Avatar     string `json:"avatar" validate:"omitempty,url,max=500"`
	Bio        string `json:"bio" validate:"omitempty,max=500"`
	Website    string `json:"website" validate:"omitempty,url,max=255"`
	IsPrivate  *bool  `json:"isPrivate" validate:"omitempty"`
}

type UserStats struct {
	Topics    int `json:"topics"`
	Videos    int `json:"videos"`
	Followers int `json:"followers"`
	Following int `json:"following"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	Avatar      string    `json:"avatar"`
	Bio         string    `json:"bio"`
	Website     string    `json:"website"`
	IsVerified  bool      `json:"isVerified"`
	IsPrivate   bool      `json:"isPrivate"`
	Stats       UserStats `json:"stats"`
	IsFollowing bool      `json:"isFollowing"`
	IsFollowedBy bool     `json:"isFollowedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserResponseAdmin struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	Website   string    `json:"website"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	IsVerified bool     `json:"isVerified"`
	IsPrivate  bool     `json:"isPrivate"`
	Stats      UserStats `json:"stats"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UserListResponse struct {
	Users []UserResponseAdmin `json:"users"`
	Meta  PaginationMeta      `json:"meta"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
}

