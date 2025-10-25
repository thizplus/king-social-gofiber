package serviceimpl

import (
	"context"
	"errors"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"math"

	"github.com/google/uuid"
)

type commentServiceImpl struct {
	commentRepo         repositories.CommentRepository
	videoRepo           repositories.VideoRepository
	userRepo            repositories.UserRepository
	notificationService services.NotificationService
}

func NewCommentService(
	commentRepo repositories.CommentRepository,
	videoRepo repositories.VideoRepository,
	userRepo repositories.UserRepository,
	notificationService services.NotificationService,
) services.CommentService {
	return &commentServiceImpl{
		commentRepo:         commentRepo,
		videoRepo:           videoRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

func (s *commentServiceImpl) CreateComment(ctx context.Context, userID uuid.UUID, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	// Verify video exists
	_, err := s.videoRepo.FindByID(ctx, req.VideoID)
	if err != nil {
		return nil, errors.New("video not found")
	}

	// If parent comment is provided, verify it exists and belongs to the same video
	if req.ParentID != nil {
		parentComment, err := s.commentRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		if parentComment.VideoID != req.VideoID {
			return nil, errors.New("parent comment does not belong to this video")
		}
	}

	// Create comment
	comment := &models.Comment{
		UserID:   userID,
		VideoID:  req.VideoID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}

	// Create notification for video comment or comment reply
	go func() {
		if req.ParentID != nil {
			// This is a reply to a comment
			_ = s.notificationService.CreateCommentReplyNotification(context.Background(), *req.ParentID, userID)
		} else {
			// This is a comment on a video
			_ = s.notificationService.CreateVideoCommentNotification(context.Background(), req.VideoID, userID)
		}
	}()

	// Update comment count in video table asynchronously
	go func() {
		count, err := s.commentRepo.CountByVideoID(context.Background(), req.VideoID)
		if err == nil {
			_ = s.videoRepo.UpdateCommentCount(context.Background(), req.VideoID, int(count))
		}
	}()

	// Load user for response
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment, user), nil
}

func (s *commentServiceImpl) GetCommentsByVideoID(ctx context.Context, videoID uuid.UUID, page, limit int) (*dto.CommentListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Get top-level comments (no parent)
	comments, totalCount, err := s.commentRepo.FindByVideoID(ctx, videoID, offset, limit)
	if err != nil {
		return nil, err
	}

	// Build response with replies
	commentResponses := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		// Load user for comment
		user, err := s.userRepo.GetByID(ctx, comment.UserID)
		if err != nil {
			continue
		}

		commentResp := s.toCommentResponse(comment, user)

		// Load replies if this is a top-level comment
		if comment.ParentID == nil {
			replies, _, err := s.commentRepo.FindReplies(ctx, comment.ID, 0, 10) // Load up to 10 replies
			if err == nil && len(replies) > 0 {
				commentResp.Replies = make([]dto.CommentResponse, 0, len(replies))
				for _, reply := range replies {
					replyUser, err := s.userRepo.GetByID(ctx, reply.UserID)
					if err != nil {
						continue
					}
					commentResp.Replies = append(commentResp.Replies, *s.toCommentResponse(reply, replyUser))
				}
			}
		}

		commentResponses = append(commentResponses, *commentResp)
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.CommentListResponse{
		Comments:   commentResponses,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *commentServiceImpl) UpdateComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, req *dto.UpdateCommentRequest) (*dto.CommentResponse, error) {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	// Check ownership
	if comment.UserID != userID {
		return nil, errors.New("you don't have permission to update this comment")
	}

	// Update content
	comment.Content = req.Content

	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, err
	}

	// Load user for response
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.toCommentResponse(comment, user), nil
}

func (s *commentServiceImpl) DeleteComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check ownership
	if comment.UserID != userID {
		return errors.New("you don't have permission to delete this comment")
	}

	// Delete comment
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	// Update comment count in video table asynchronously
	go func() {
		count, err := s.commentRepo.CountByVideoID(context.Background(), comment.VideoID)
		if err == nil {
			_ = s.videoRepo.UpdateCommentCount(context.Background(), comment.VideoID, int(count))
		}
	}()

	return nil
}

func (s *commentServiceImpl) DeleteCommentByAdmin(ctx context.Context, commentID uuid.UUID) error {
	// Get comment to get video ID for count update
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	// Delete comment
	if err := s.commentRepo.Delete(ctx, commentID); err != nil {
		return err
	}

	// Update comment count in video table asynchronously
	go func() {
		count, err := s.commentRepo.CountByVideoID(context.Background(), comment.VideoID)
		if err == nil {
			_ = s.videoRepo.UpdateCommentCount(context.Background(), comment.VideoID, int(count))
		}
	}()

	return nil
}

// Helper method to convert Comment to CommentResponse
func (s *commentServiceImpl) toCommentResponse(comment *models.Comment, user *models.User) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:       comment.ID,
		UserID:   comment.UserID,
		User: dto.UserSummaryComment{
			ID:        user.ID,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Avatar:    user.Avatar,
		},
		VideoID:   comment.VideoID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		Replies:   []dto.CommentResponse{}, // Initialize empty slice
	}
}
