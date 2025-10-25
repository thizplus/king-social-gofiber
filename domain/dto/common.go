package dto

import "github.com/google/uuid"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    PaginationMeta `json:"meta"`
	Error   string      `json:"error,omitempty"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	Offset      int   `json:"offset"`
	Limit       int   `json:"limit"`
	Page        int   `json:"page"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

// NewPaginationMeta creates pagination metadata with calculated fields
func NewPaginationMeta(total int64, offset, limit int) PaginationMeta {
	page := (offset / limit) + 1
	if limit == 0 {
		limit = 1 // Prevent division by zero
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return PaginationMeta{
		Total:       total,
		Offset:      offset,
		Limit:       limit,
		Page:        page,
		TotalPages:  totalPages,
		HasNext:     offset+limit < int(total),
		HasPrevious: offset > 0,
	}
}

type IDRequest struct {
	ID uuid.UUID `json:"id" validate:"required" param:"id"`
}

type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}