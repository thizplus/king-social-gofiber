package postgres

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateFollowerCount(ctx context.Context, userID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("follower_count", count).Error
}

func (r *UserRepositoryImpl) UpdateFollowingCount(ctx context.Context, userID uuid.UUID, count int) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("following_count", count).Error
}

// Admin methods

func (r *UserRepositoryImpl) FindAllWithStats(ctx context.Context, params *dto.AdminUserListRequest) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Apply role filter
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	// Apply isActive filter
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get users with ordering
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, totalCount, err
}

func (r *UserRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) GetNewUsersToday(ctx context.Context) (int64, error) {
	var count int64
	startOfDay := time.Now().Truncate(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("created_at >= ?", startOfDay).
		Count(&count).Error

	return count, err
}

func (r *UserRepositoryImpl) GetActiveUsersCount(ctx context.Context) (int64, error) {
	var count int64
	last24Hours := time.Now().Add(-24 * time.Hour)

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("last_login_at >= ?", last24Hours).
		Count(&count).Error

	return count, err
}

func (r *UserRepositoryImpl) SuspendUser(ctx context.Context, userID uuid.UUID, reason string, until time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"suspended_until": until,
			"suspend_reason":  reason,
		}).Error
}

func (r *UserRepositoryImpl) UpdateRole(ctx context.Context, userID uuid.UUID, role string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("role", role).Error
}
