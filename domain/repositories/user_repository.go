package repositories

import (
	"context"
	"time"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateFollowerCount(ctx context.Context, userID uuid.UUID, count int) error
	UpdateFollowingCount(ctx context.Context, userID uuid.UUID, count int) error

	// Admin methods
	FindAllWithStats(ctx context.Context, params *dto.AdminUserListRequest) ([]models.User, int64, error)
	GetTotalCount(ctx context.Context) (int64, error)
	GetNewUsersToday(ctx context.Context) (int64, error)
	GetActiveUsersCount(ctx context.Context) (int64, error)
	SuspendUser(ctx context.Context, userID uuid.UUID, reason string, until time.Time) error
	UpdateRole(ctx context.Context, userID uuid.UUID, role string) error
}
