package handlers

import (
	"gofiber-social/domain/services"
)

// Services contains all the services needed for handlers
type Services struct {
	UserService         services.UserService
	TaskService         services.TaskService
	FileService         services.FileService
	JobService          services.JobService
	ForumService        services.ForumService
	TopicService        services.TopicService
	ReplyService        services.ReplyService
	TagService          services.TagService
	VideoService        services.VideoService
	LikeService         services.LikeService
	CommentService      services.CommentService
	ShareService        services.ShareService
	FollowService       services.FollowService
	NotificationService services.NotificationService
	AdminService        services.AdminService
	ReportService       services.ReportService
}

// Handlers contains all HTTP handlers
type Handlers struct {
	UserHandler         *UserHandler
	TaskHandler         *TaskHandler
	FileHandler         *FileHandler
	JobHandler          *JobHandler
	ForumHandler        *ForumHandler
	TopicHandler        *TopicHandler
	ReplyHandler        *ReplyHandler
	TagHandler          *TagHandler
	VideoHandler        *VideoHandler
	LikeHandler         *LikeHandler
	CommentHandler      *CommentHandler
	ShareHandler        *ShareHandler
	FollowHandler       *FollowHandler
	NotificationHandler *NotificationHandler
	AdminHandler        *AdminHandler
	ReportHandler       *ReportHandler
}

// NewHandlers creates a new instance of Handlers with all dependencies
func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		UserHandler:         NewUserHandler(services.UserService),
		TaskHandler:         NewTaskHandler(services.TaskService),
		FileHandler:         NewFileHandler(services.FileService),
		JobHandler:          NewJobHandler(services.JobService),
		ForumHandler:        NewForumHandler(services.ForumService),
		TopicHandler:        NewTopicHandler(services.TopicService),
		ReplyHandler:        NewReplyHandler(services.ReplyService),
		TagHandler:          NewTagHandler(services.TagService),
		VideoHandler:        NewVideoHandler(services.VideoService),
		LikeHandler:         NewLikeHandler(services.LikeService),
		CommentHandler:      NewCommentHandler(services.CommentService),
		ShareHandler:        NewShareHandler(services.ShareService),
		FollowHandler:       NewFollowHandler(services.FollowService),
		NotificationHandler: NewNotificationHandler(services.NotificationService),
		AdminHandler:        NewAdminHandler(services.AdminService),
		ReportHandler:       NewReportHandler(services.ReportService),
	}
}
