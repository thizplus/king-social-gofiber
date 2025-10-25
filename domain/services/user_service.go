package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, *models.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetProfileWithStats(ctx context.Context, userID uuid.UUID, viewerID *uuid.UUID) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, offset, limit int) ([]*models.User, int64, error)
	GetUserStats(ctx context.Context, userID uuid.UUID) (topicCount, videoCount int, err error)
	GenerateJWT(user *models.User) (string, error)
	ValidateJWT(token string) (*models.User, error)
}
