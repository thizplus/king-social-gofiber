package dto

import (
	"gofiber-social/domain/models"

	"github.com/google/uuid"
)

// UserToUserResponse converts User model to UserResponse DTO (old format for backward compatibility)
func UserToUserResponseOld(user *models.User) *UserResponseAdmin {
	if user == nil {
		return nil
	}
	return &UserResponseAdmin{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Website:   user.Website,
		Role:      user.Role,
		IsActive:  user.IsActive,
		IsVerified: user.IsVerified,
		IsPrivate:  user.IsPrivate,
		Stats: UserStats{
			Topics:    0,
			Videos:    0,
			Followers: user.FollowerCount,
			Following: user.FollowingCount,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserToUserResponse converts User model to UserResponse DTO
func UserToUserResponse(user *models.User, topicCount, videoCount int, isFollowing, isFollowedBy bool) *UserResponse {
	if user == nil {
		return nil
	}
	displayName := user.FirstName
	if user.LastName != "" {
		displayName = user.FirstName + " " + user.LastName
	}
	if displayName == "" {
		displayName = user.Username
	}

	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: displayName,
		Avatar:      user.Avatar,
		Bio:         user.Bio,
		Website:     user.Website,
		IsVerified:  user.IsVerified,
		IsPrivate:   user.IsPrivate,
		Stats: UserStats{
			Topics:    topicCount,
			Videos:    videoCount,
			Followers: user.FollowerCount,
			Following: user.FollowingCount,
		},
		IsFollowing:  isFollowing,
		IsFollowedBy: isFollowedBy,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

// UserToUserResponseAdmin converts User model to UserResponseAdmin DTO for admin views
func UserToUserResponseAdmin(user *models.User, topicCount, videoCount int) *UserResponseAdmin {
	if user == nil {
		return nil
	}
	return &UserResponseAdmin{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Website:   user.Website,
		Role:      user.Role,
		IsActive:  user.IsActive,
		IsVerified: user.IsVerified,
		IsPrivate:  user.IsPrivate,
		Stats: UserStats{
			Topics:    topicCount,
			Videos:    videoCount,
			Followers: user.FollowerCount,
			Following: user.FollowingCount,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func CreateUserRequestToUser(req *CreateUserRequest) *models.User {
	return &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}

func UpdateUserRequestToUser(req *UpdateUserRequest) *models.User {
	return &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Avatar:    req.Avatar,
	}
}

func TaskToTaskResponse(task *models.Task, user *models.User) *TaskResponse {
	if task == nil {
		return nil
	}
	taskResp := &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if user != nil {
		taskResp.User = *UserToUserResponseOld(user)
	}
	return taskResp
}

func CreateTaskRequestToTask(req *CreateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func UpdateTaskRequestToTask(req *UpdateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func JobToJobResponse(job *models.Job) *JobResponse {
	if job == nil {
		return nil
	}
	return &JobResponse{
		ID:        job.ID,
		Name:      job.Name,
		CronExpr:  job.CronExpr,
		Payload:   job.Payload,
		Status:    job.Status,
		LastRun:   job.LastRun,
		NextRun:   job.NextRun,
		IsActive:  job.IsActive,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

func CreateJobRequestToJob(req *CreateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
	}
}

func UpdateJobRequestToJob(req *UpdateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
		IsActive: req.IsActive,
	}
}

func FileToFileResponse(file *models.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		URL:       file.URL,
		CDNPath:   file.CDNPath,
		UserID:    file.UserID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

func ForumToForumResponse(forum *models.Forum) *ForumResponse {
	if forum == nil {
		return nil
	}
	return &ForumResponse{
		ID:          forum.ID,
		Name:        forum.Name,
		Slug:        forum.Slug,
		Description: forum.Description,
		Icon:        forum.Icon,
		Order:       forum.Order,
		IsActive:    forum.IsActive,
		TopicCount:  forum.TopicCount,
		CreatedAt:   forum.CreatedAt,
		UpdatedAt:   forum.UpdatedAt,
	}
}

func TopicToTopicResponse(topic *models.Topic) *TopicResponse {
	if topic == nil {
		return nil
	}
	resp := &TopicResponse{
		ID:         topic.ID,
		ForumID:    topic.ForumID,
		UserID:     topic.UserID,
		Title:      topic.Title,
		Content:    topic.Content,
		Thumbnail:  topic.Thumbnail,
		ViewCount:  topic.ViewCount,
		ReplyCount: topic.ReplyCount,
		IsPinned:   topic.IsPinned,
		IsLocked:   topic.IsLocked,
		CreatedAt:  topic.CreatedAt,
		UpdatedAt:  topic.UpdatedAt,
	}

	// Include Forum if loaded
	if topic.Forum.ID != uuid.Nil {
		resp.Forum = ForumToForumResponse(&topic.Forum)
	}

	// Include User if loaded
	if topic.User.ID != uuid.Nil {
		resp.User = UserToUserResponseOld(&topic.User)
	}

	// Include Tags if loaded
	if len(topic.Tags) > 0 {
		resp.Tags = make([]TagResponse, len(topic.Tags))
		for i, tag := range topic.Tags {
			resp.Tags[i] = *TagToTagResponse(&tag)
		}
	}

	return resp
}

func ReplyToReplyResponse(reply *models.Reply, includeNested bool) *ReplyResponse {
	if reply == nil {
		return nil
	}

	resp := &ReplyResponse{
		ID:        reply.ID,
		TopicID:   reply.TopicID,
		UserID:    reply.UserID,
		ParentID:  reply.ParentID,
		Content:   reply.Content,
		CreatedAt: reply.CreatedAt,
		UpdatedAt: reply.UpdatedAt,
	}

	// Include User if loaded
	if reply.User.ID != uuid.Nil {
		resp.User = UserToUserResponseOld(&reply.User)
	}

	// Include nested replies if requested
	if includeNested && len(reply.Replies) > 0 {
		resp.Replies = make([]ReplyResponse, len(reply.Replies))
		for i, nested := range reply.Replies {
			resp.Replies[i] = *ReplyToReplyResponse(&nested, true)
		}
	}

	return resp
}

func TagToTagResponse(tag *models.Tag) *TagResponse {
	if tag == nil {
		return nil
	}
	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		Color:       tag.Color,
		IsActive:    tag.IsActive,
		UsageCount:  tag.UsageCount,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}
}
