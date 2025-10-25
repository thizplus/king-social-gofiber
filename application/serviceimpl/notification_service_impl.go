package serviceimpl

import (
	"context"
	"errors"
	"fmt"
	"math"

	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"gofiber-social/infrastructure/websocket"

	"github.com/google/uuid"
)

type notificationServiceImpl struct {
	notificationRepo repositories.NotificationRepository
	userRepo         repositories.UserRepository
	topicRepo        repositories.TopicRepository
	videoRepo        repositories.VideoRepository
	commentRepo      repositories.CommentRepository
	replyRepo        repositories.ReplyRepository
}

func NewNotificationService(
	notificationRepo repositories.NotificationRepository,
	userRepo repositories.UserRepository,
	topicRepo repositories.TopicRepository,
	videoRepo repositories.VideoRepository,
	commentRepo repositories.CommentRepository,
	replyRepo repositories.ReplyRepository,
) services.NotificationService {
	return &notificationServiceImpl{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		topicRepo:        topicRepo,
		videoRepo:        videoRepo,
		commentRepo:      commentRepo,
		replyRepo:        replyRepo,
	}
}

// Helper function to broadcast notification via WebSocket
func (s *notificationServiceImpl) broadcastNotification(notification *models.Notification) {
	// Send real-time notification via WebSocket
	go func() {
		websocket.Manager.BroadcastToUser(notification.UserID, "notification", map[string]interface{}{
			"id":        notification.ID,
			"type":      notification.Type,
			"message":   notification.Message,
			"actorId":   notification.ActorID,
			"isRead":    notification.IsRead,
			"createdAt": notification.CreatedAt,
		})
	}()
}

// Create notifications
func (s *notificationServiceImpl) CreateTopicReplyNotification(ctx context.Context, topicID, replyUserID uuid.UUID) error {
	// Get topic
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// Don't notify if replying to own topic
	if topic.UserID == replyUserID {
		return nil
	}

	// Get replier user
	replier, err := s.userRepo.FindByID(ctx, replyUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     topic.UserID,
		ActorID:    replyUserID,
		Type:       models.NotificationTypeTopicReply,
		ResourceID: &topicID,
		Message:    fmt.Sprintf("%s ตอบกระทู้ของคุณ: %s", replier.Username, topic.Title),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateTopicLikeNotification(ctx context.Context, topicID, likerUserID uuid.UUID) error {
	// Get topic
	topic, err := s.topicRepo.GetByID(ctx, topicID)
	if err != nil {
		return err
	}

	// Don't notify if liking own topic
	if topic.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     topic.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeTopicLike,
		ResourceID: &topicID,
		Message:    fmt.Sprintf("%s ถูกใจกระทู้ของคุณ: %s", liker.Username, topic.Title),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateVideoLikeNotification(ctx context.Context, videoID, likerUserID uuid.UUID) error {
	// Get video
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Don't notify if liking own video
	if video.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     video.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeVideoLike,
		ResourceID: &videoID,
		Message:    fmt.Sprintf("%s ถูกใจวิดีโอของคุณ: %s", liker.Username, video.Title),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateVideoCommentNotification(ctx context.Context, videoID, commenterUserID uuid.UUID) error {
	// Get video
	video, err := s.videoRepo.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	// Don't notify if commenting on own video
	if video.UserID == commenterUserID {
		return nil
	}

	// Get commenter user
	commenter, err := s.userRepo.FindByID(ctx, commenterUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     video.UserID,
		ActorID:    commenterUserID,
		Type:       models.NotificationTypeVideoComment,
		ResourceID: &videoID,
		Message:    fmt.Sprintf("%s แสดงความคิดเห็นในวิดีโอของคุณ: %s", commenter.Username, video.Title),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateCommentReplyNotification(ctx context.Context, commentID, replierUserID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Don't notify if replying to own comment
	if comment.UserID == replierUserID {
		return nil
	}

	// Get replier user
	replier, err := s.userRepo.FindByID(ctx, replierUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     comment.UserID,
		ActorID:    replierUserID,
		Type:       models.NotificationTypeCommentReply,
		ResourceID: &commentID,
		Message:    fmt.Sprintf("%s ตอบกลับความคิดเห็นของคุณ", replier.Username),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateReplyLikeNotification(ctx context.Context, replyID, likerUserID uuid.UUID) error {
	// Get reply
	reply, err := s.replyRepo.GetByID(ctx, replyID)
	if err != nil {
		return err
	}

	// Don't notify if liking own reply
	if reply.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     reply.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeReplyLike,
		ResourceID: &replyID,
		Message:    fmt.Sprintf("%s ถูกใจการตอบกลับของคุณ", liker.Username),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateCommentLikeNotification(ctx context.Context, commentID, likerUserID uuid.UUID) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	// Don't notify if liking own comment
	if comment.UserID == likerUserID {
		return nil
	}

	// Get liker user
	liker, err := s.userRepo.FindByID(ctx, likerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:     comment.UserID,
		ActorID:    likerUserID,
		Type:       models.NotificationTypeCommentLike,
		ResourceID: &commentID,
		Message:    fmt.Sprintf("%s ถูกใจความคิดเห็นของคุณ", liker.Username),
		IsRead:     false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

func (s *notificationServiceImpl) CreateNewFollowerNotification(ctx context.Context, followedUserID, followerUserID uuid.UUID) error {
	// Get follower user
	follower, err := s.userRepo.FindByID(ctx, followerUserID)
	if err != nil {
		return err
	}

	notification := &models.Notification{
		UserID:  followedUserID,
		ActorID: followerUserID,
		Type:    models.NotificationTypeNewFollower,
		Message: fmt.Sprintf("%s เริ่มติดตามคุณ", follower.Username),
		IsRead:  false,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	s.broadcastNotification(notification)
	return nil
}

// Read notifications
func (s *notificationServiceImpl) GetNotifications(ctx context.Context, userID uuid.UUID, params *dto.NotificationQueryParams) (*dto.NotificationListResponse, error) {
	notifications, totalCount, err := s.notificationRepo.FindByUserID(ctx, userID, params)
	if err != nil {
		return nil, err
	}

	// Get unread count
	unreadCount, _ := s.notificationRepo.GetUnreadCount(ctx, userID)

	// Convert to response
	notificationResponses := make([]dto.NotificationResponse, len(notifications))
	for i, notif := range notifications {
		notificationResponses[i] = dto.NotificationResponse{
			ID:     notif.ID,
			UserID: notif.UserID,
			Actor: dto.ActorSummary{
				ID:       notif.Actor.ID,
				Username: notif.Actor.Username,
				FullName: notif.Actor.FullName,
				Avatar:   notif.Actor.Avatar,
			},
			Type:       notif.Type,
			ResourceID: notif.ResourceID,
			Message:    notif.Message,
			IsRead:     notif.IsRead,
			CreatedAt:  notif.CreatedAt,
		}
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &dto.NotificationListResponse{
		Notifications: notificationResponses,
		TotalCount:    totalCount,
		UnreadCount:   unreadCount,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}, nil
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID uuid.UUID) (*dto.UnreadCountResponse, error) {
	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UnreadCountResponse{Count: count}, nil
}

// Update notifications
func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized to mark this notification")
	}

	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *notificationServiceImpl) MarkMultipleAsRead(ctx context.Context, userID uuid.UUID, notificationIDs []uuid.UUID) (*dto.MarkAsReadResponse, error) {
	if len(notificationIDs) == 0 {
		return &dto.MarkAsReadResponse{
			Message: "No notifications to mark as read",
			Count:   0,
		}, nil
	}

	// TODO: Verify all notifications belong to user (optional, depends on security requirements)

	if err := s.notificationRepo.MarkMultipleAsRead(ctx, notificationIDs); err != nil {
		return nil, err
	}

	return &dto.MarkAsReadResponse{
		Message: "Notifications marked as read",
		Count:   len(notificationIDs),
	}, nil
}

func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID uuid.UUID) (*dto.MarkAsReadResponse, error) {
	// Get unread count before marking
	count, _ := s.notificationRepo.GetUnreadCount(ctx, userID)

	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return nil, err
	}

	return &dto.MarkAsReadResponse{
		Message: "All notifications marked as read",
		Count:   int(count),
	}, nil
}

// Delete notifications
func (s *notificationServiceImpl) DeleteNotification(ctx context.Context, userID, notificationID uuid.UUID) error {
	// Verify notification belongs to user
	notification, err := s.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	if notification.UserID != userID {
		return errors.New("unauthorized to delete this notification")
	}

	return s.notificationRepo.Delete(ctx, notificationID)
}
