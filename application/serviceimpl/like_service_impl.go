package serviceimpl

import (
	"context"
	"errors"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"

	"github.com/google/uuid"
)

type likeServiceImpl struct {
	likeRepo            repositories.LikeRepository
	topicRepo           repositories.TopicRepository
	videoRepo           repositories.VideoRepository
	replyRepo           repositories.ReplyRepository
	commentRepo         repositories.CommentRepository
	notificationService services.NotificationService
}

func NewLikeService(
	likeRepo repositories.LikeRepository,
	topicRepo repositories.TopicRepository,
	videoRepo repositories.VideoRepository,
	replyRepo repositories.ReplyRepository,
	commentRepo repositories.CommentRepository,
	notificationService services.NotificationService,
) services.LikeService {
	return &likeServiceImpl{
		likeRepo:            likeRepo,
		topicRepo:           topicRepo,
		videoRepo:           videoRepo,
		replyRepo:           replyRepo,
		commentRepo:         commentRepo,
		notificationService: notificationService,
	}
}

// Topic Likes
func (s *likeServiceImpl) LikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Verify topic exists
	_, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return nil, errors.New("topic not found")
	}

	// Check if already liked
	isLiked, err := s.likeRepo.IsTopicLikedByUser(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		return nil, errors.New("topic already liked")
	}

	// Create like
	_, err = s.likeRepo.LikeTopic(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	// Create notification for topic like
	go func() {
		_ = s.notificationService.CreateTopicLikeNotification(context.Background(), topicID, userID)
	}()

	// Update like count in topic table asynchronously
	go func() {
		count, err := s.likeRepo.CountTopicLikes(context.Background(), topicID)
		if err == nil {
			_ = s.topicRepo.UpdateLikeCount(context.Background(), topicID, int(count))
		}
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountTopicLikes(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   true,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) UnlikeTopic(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsTopicLikedByUser(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	if !isLiked {
		return nil, errors.New("topic not liked")
	}

	// Remove like
	err = s.likeRepo.UnlikeTopic(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	// Update like count in topic table asynchronously
	go func() {
		count, err := s.likeRepo.CountTopicLikes(context.Background(), topicID)
		if err == nil {
			_ = s.topicRepo.UpdateLikeCount(context.Background(), topicID, int(count))
		}
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountTopicLikes(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   false,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) GetTopicLikeStatus(ctx context.Context, userID uuid.UUID, topicID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsTopicLikedByUser(ctx, userID, topicID)
	if err != nil {
		return nil, err
	}

	// Get like count
	likeCount, err := s.likeRepo.CountTopicLikes(ctx, topicID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// Video Likes
func (s *likeServiceImpl) LikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Verify video exists
	_, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	// Check if already liked
	isLiked, err := s.likeRepo.IsVideoLikedByUser(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		return nil, errors.New("video already liked")
	}

	// Create like
	_, err = s.likeRepo.LikeVideo(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	// Create notification for video like
	go func() {
		_ = s.notificationService.CreateVideoLikeNotification(context.Background(), videoID, userID)
	}()

	// Update like count in video table asynchronously
	go func() {
		count, err := s.likeRepo.CountVideoLikes(context.Background(), videoID)
		if err == nil {
			_ = s.videoRepo.UpdateLikeCount(context.Background(), videoID, int(count))
		}
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountVideoLikes(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   true,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) UnlikeVideo(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsVideoLikedByUser(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	if !isLiked {
		return nil, errors.New("video not liked")
	}

	// Remove like
	err = s.likeRepo.UnlikeVideo(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	// Update like count in video table asynchronously
	go func() {
		count, err := s.likeRepo.CountVideoLikes(context.Background(), videoID)
		if err == nil {
			_ = s.videoRepo.UpdateLikeCount(context.Background(), videoID, int(count))
		}
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountVideoLikes(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   false,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) GetVideoLikeStatus(ctx context.Context, userID uuid.UUID, videoID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsVideoLikedByUser(ctx, userID, videoID)
	if err != nil {
		return nil, err
	}

	// Get like count
	likeCount, err := s.likeRepo.CountVideoLikes(ctx, videoID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// Reply Likes
func (s *likeServiceImpl) LikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Verify reply exists
	_, err := s.replyRepo.GetByID(ctx, replyID)
	if err != nil {
		return nil, errors.New("reply not found")
	}

	// Check if already liked
	isLiked, err := s.likeRepo.IsReplyLikedByUser(ctx, userID, replyID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		return nil, errors.New("reply already liked")
	}

	// Create like
	_, err = s.likeRepo.LikeReply(ctx, userID, replyID)
	if err != nil {
		return nil, err
	}

	// Create notification for reply like
	go func() {
		_ = s.notificationService.CreateReplyLikeNotification(context.Background(), replyID, userID)
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountReplyLikes(ctx, replyID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   true,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) UnlikeReply(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsReplyLikedByUser(ctx, userID, replyID)
	if err != nil {
		return nil, err
	}

	if !isLiked {
		return nil, errors.New("reply not liked")
	}

	// Remove like
	err = s.likeRepo.UnlikeReply(ctx, userID, replyID)
	if err != nil {
		return nil, err
	}

	// Get updated like count
	likeCount, err := s.likeRepo.CountReplyLikes(ctx, replyID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   false,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) GetReplyLikeStatus(ctx context.Context, userID uuid.UUID, replyID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsReplyLikedByUser(ctx, userID, replyID)
	if err != nil {
		return nil, err
	}

	// Get like count
	likeCount, err := s.likeRepo.CountReplyLikes(ctx, replyID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}

// Comment Likes
func (s *likeServiceImpl) LikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Verify comment exists
	_, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Check if already liked
	isLiked, err := s.likeRepo.IsCommentLikedByUser(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	if isLiked {
		return nil, errors.New("comment already liked")
	}

	// Create like
	_, err = s.likeRepo.LikeComment(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	// Create notification for comment like
	go func() {
		_ = s.notificationService.CreateCommentLikeNotification(context.Background(), commentID, userID)
	}()

	// Get updated like count
	likeCount, err := s.likeRepo.CountCommentLikes(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   true,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) UnlikeComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsCommentLikedByUser(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	if !isLiked {
		return nil, errors.New("comment not liked")
	}

	// Remove like
	err = s.likeRepo.UnlikeComment(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	// Get updated like count
	likeCount, err := s.likeRepo.CountCommentLikes(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   false,
		LikeCount: likeCount,
	}, nil
}

func (s *likeServiceImpl) GetCommentLikeStatus(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) (*dto.LikeStatusResponse, error) {
	// Check if liked
	isLiked, err := s.likeRepo.IsCommentLikedByUser(ctx, userID, commentID)
	if err != nil {
		return nil, err
	}

	// Get like count
	likeCount, err := s.likeRepo.CountCommentLikes(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return &dto.LikeStatusResponse{
		IsLiked:   isLiked,
		LikeCount: likeCount,
	}, nil
}
