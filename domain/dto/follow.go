package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type FollowUserRequest struct {
	FollowingID uuid.UUID `json:"followingId" validate:"required,uuid"`
}

// Response DTOs
type FollowResponse struct {
	ID          uuid.UUID `json:"id,omitempty"`
	FollowerID  uuid.UUID `json:"followerId"`
	FollowingID uuid.UUID `json:"followingId"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	Message     string    `json:"message"`
}

type FollowStatusResponse struct {
	IsFollowing bool `json:"isFollowing"`
}

type FollowerResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Avatar         string    `json:"avatar,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	IsFollowing    bool      `json:"isFollowing"` // ว่าเราติดตามคนนี้หรือไม่
	FollowedAt     time.Time `json:"followedAt"`
}

type FollowingResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	FullName       string    `json:"fullName"`
	Avatar         string    `json:"avatar,omitempty"`
	Bio            string    `json:"bio,omitempty"`
	FollowerCount  int       `json:"followerCount"`
	FollowingCount int       `json:"followingCount"`
	IsFollowing    bool      `json:"isFollowing"` // ว่าเราติดตามคนนี้หรือไม่
	FollowedAt     time.Time `json:"followedAt"`
}

type FollowListResponse struct {
	Users      []FollowerResponse `json:"users"`
	TotalCount int64              `json:"totalCount"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

type FollowingListResponse struct {
	Users      []FollowingResponse `json:"users"`
	TotalCount int64               `json:"totalCount"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"totalPages"`
}

type UserStatsResponse struct {
	UserID         uuid.UUID `json:"userId"`
	FollowerCount  int64     `json:"followerCount"`
	FollowingCount int64     `json:"followingCount"`
}
