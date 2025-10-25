package postgres

import (
	"context"
	"errors"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type likeRepositoryImpl struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) repositories.LikeRepository {
	return &likeRepositoryImpl{db: db}
}

// LikeTopic creates a like for a topic
func (r *likeRepositoryImpl) LikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*models.Like, error) {
	// Check if already liked
	var existingLike models.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		First(&existingLike).Error

	if err == nil {
		// Already liked
		return &existingLike, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new like
	like := &models.Like{
		UserID:  userID,
		TopicID: &topicID,
	}

	if err := like.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		return nil, err
	}

	return like, nil
}

// UnlikeTopic removes a like from a topic
func (r *likeRepositoryImpl) UnlikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

// IsTopicLikedByUser checks if a user has liked a topic
func (r *likeRepositoryImpl) IsTopicLikedByUser(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountTopicLikes returns the number of likes for a topic
func (r *likeRepositoryImpl) CountTopicLikes(ctx context.Context, topicID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("topic_id = ?", topicID).
		Count(&count).Error

	return count, err
}

// GetTopicLikesByUserID retrieves all topic likes by a user
func (r *likeRepositoryImpl) GetTopicLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var totalCount int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND topic_id IS NOT NULL", userID).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get likes with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Topic").
		Preload("Topic.User").
		Preload("Topic.Forum").
		Where("user_id = ? AND topic_id IS NOT NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error

	return likes, totalCount, err
}

// LikeVideo creates a like for a video
func (r *likeRepositoryImpl) LikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*models.Like, error) {
	// Check if already liked
	var existingLike models.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		First(&existingLike).Error

	if err == nil {
		// Already liked
		return &existingLike, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new like
	like := &models.Like{
		UserID:  userID,
		VideoID: &videoID,
	}

	if err := like.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		return nil, err
	}

	return like, nil
}

// UnlikeVideo removes a like from a video
func (r *likeRepositoryImpl) UnlikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

// IsVideoLikedByUser checks if a user has liked a video
func (r *likeRepositoryImpl) IsVideoLikedByUser(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CountVideoLikes returns the number of likes for a video
func (r *likeRepositoryImpl) CountVideoLikes(ctx context.Context, videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("video_id = ?", videoID).
		Count(&count).Error

	return count, err
}

// GetVideoLikesByUserID retrieves all video likes by a user
func (r *likeRepositoryImpl) GetVideoLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var totalCount int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND video_id IS NOT NULL", userID).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get likes with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Video").
		Preload("Video.User").
		Where("user_id = ? AND video_id IS NOT NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error

	return likes, totalCount, err
}

// GetByID retrieves a like by ID
func (r *likeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Like, error) {
	var like models.Like
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Topic").
		Preload("Video").
		Where("id = ?", id).
		First(&like).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("like not found")
		}
		return nil, err
	}

	return &like, nil
}

// Delete removes a like by ID
func (r *likeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Like{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

// Reply Likes
func (r *likeRepositoryImpl) LikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*models.Like, error) {
	// Check if already liked
	var existingLike models.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND reply_id = ?", userID, replyID).
		First(&existingLike).Error

	if err == nil {
		// Already liked
		return &existingLike, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new like
	like := &models.Like{
		UserID:  userID,
		ReplyID: &replyID,
	}

	if err := like.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		return nil, err
	}

	return like, nil
}

func (r *likeRepositoryImpl) UnlikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND reply_id = ?", userID, replyID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *likeRepositoryImpl) IsReplyLikedByUser(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND reply_id = ?", userID, replyID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *likeRepositoryImpl) CountReplyLikes(ctx context.Context, replyID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("reply_id = ?", replyID).
		Count(&count).Error

	return count, err
}

func (r *likeRepositoryImpl) GetReplyLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var totalCount int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND reply_id IS NOT NULL", userID).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get likes with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Reply").
		Preload("Reply.User").
		Where("user_id = ? AND reply_id IS NOT NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error

	return likes, totalCount, err
}

// Comment Likes
func (r *likeRepositoryImpl) LikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*models.Like, error) {
	// Check if already liked
	var existingLike models.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		First(&existingLike).Error

	if err == nil {
		// Already liked
		return &existingLike, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new like
	like := &models.Like{
		UserID:    userID,
		CommentID: &commentID,
	}

	if err := like.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		return nil, err
	}

	return like, nil
}

func (r *likeRepositoryImpl) UnlikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Delete(&models.Like{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("like not found")
	}

	return nil
}

func (r *likeRepositoryImpl) IsCommentLikedByUser(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND comment_id = ?", userID, commentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *likeRepositoryImpl) CountCommentLikes(ctx context.Context, commentID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("comment_id = ?", commentID).
		Count(&count).Error

	return count, err
}

func (r *likeRepositoryImpl) GetCommentLikesByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Like, int64, error) {
	var likes []*models.Like
	var totalCount int64

	// Count total
	if err := r.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND comment_id IS NOT NULL", userID).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get likes with preloaded relationships
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Comment").
		Preload("Comment.User").
		Where("user_id = ? AND comment_id IS NOT NULL", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error

	return likes, totalCount, err
}
