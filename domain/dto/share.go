package dto

import (
	"time"

	"github.com/google/uuid"
)

// Request DTOs
type ShareVideoRequest struct {
	VideoID  uuid.UUID `json:"videoId" validate:"required,uuid"`
	Platform string    `json:"platform" validate:"required,oneof=facebook twitter line copy_link"`
}

// Response DTOs
type ShareResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	VideoID   uuid.UUID `json:"videoId"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
	Message   string    `json:"message"`
}

type ShareCountResponse struct {
	VideoID    uuid.UUID `json:"videoId"`
	ShareCount int64     `json:"shareCount"`
}
