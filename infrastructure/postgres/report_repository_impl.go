package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
)

type reportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) repositories.ReportRepository {
	return &reportRepositoryImpl{db: db}
}

func (r *reportRepositoryImpl) Create(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *reportRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*models.Report, error) {
	var report models.Report
	err := r.db.WithContext(ctx).
		Preload("Reporter").
		Preload("Reviewer").
		First(&report, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, err
	}
	return &report, nil
}

func (r *reportRepositoryImpl) FindAll(ctx context.Context, params *dto.ReportQueryParams) ([]models.Report, int64, error) {
	var reports []models.Report
	var totalCount int64

	query := r.db.WithContext(ctx).Model(&models.Report{}).
		Preload("Reporter").
		Preload("Reviewer")

	// Filter by type
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	// Filter by status
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
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

	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reports).Error

	return reports, totalCount, err
}

func (r *reportRepositoryImpl) Update(ctx context.Context, report *models.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r *reportRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Report{}, id).Error
}

func (r *reportRepositoryImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Report{}).Count(&count).Error
	return count, err
}

func (r *reportRepositoryImpl) GetPendingCount(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Report{}).
		Where("status = ?", models.ReportStatusPending).
		Count(&count).Error
	return count, err
}
