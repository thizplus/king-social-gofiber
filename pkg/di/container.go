package di

import (
	"context"
	"gofiber-social/application/serviceimpl"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"gofiber-social/infrastructure/postgres"
	"gofiber-social/infrastructure/redis"
	"gofiber-social/infrastructure/storage"
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/pkg/config"
	"gofiber-social/pkg/scheduler"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB             *gorm.DB
	RedisClient    *redis.RedisClient
	BunnyStorage   storage.BunnyStorage
	EventScheduler scheduler.EventScheduler

	// Repositories
	UserRepository         repositories.UserRepository
	TaskRepository         repositories.TaskRepository
	FileRepository         repositories.FileRepository
	JobRepository          repositories.JobRepository
	ForumRepository        repositories.ForumRepository
	TopicRepository        repositories.TopicRepository
	ReplyRepository        repositories.ReplyRepository
	TagRepository          repositories.TagRepository
	VideoRepository        repositories.VideoRepository
	LikeRepository         repositories.LikeRepository
	CommentRepository      repositories.CommentRepository
	ShareRepository        repositories.ShareRepository
	FollowRepository       repositories.FollowRepository
	NotificationRepository repositories.NotificationRepository
	ReportRepository       repositories.ReportRepository
	ActivityLogRepository  repositories.ActivityLogRepository

	// Services
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

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initInfrastructure(); err != nil {
		return err
	}

	if err := c.initRepositories(); err != nil {
		return err
	}

	if err := c.initServices(); err != nil {
		return err
	}

	if err := c.initScheduler(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.Config = cfg
	log.Println("✓ Configuration loaded")
	return nil
}

func (c *Container) initInfrastructure() error {
	// Initialize Database
	dbConfig := postgres.DatabaseConfig{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db
	log.Println("✓ Database connected")

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		return err
	}
	log.Println("✓ Database migrated")

	// Initialize Redis
	redisConfig := redis.RedisConfig{
		Host:     c.Config.Redis.Host,
		Port:     c.Config.Redis.Port,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	}
	c.RedisClient = redis.NewRedisClient(redisConfig)

	// Test Redis connection
	if err := c.RedisClient.Ping(context.Background()); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("✓ Redis connected")
	}

	// Initialize Bunny Storage
	bunnyConfig := storage.BunnyConfig{
		StorageZone: c.Config.Bunny.StorageZone,
		AccessKey:   c.Config.Bunny.AccessKey,
		BaseURL:     c.Config.Bunny.BaseURL,
		CDNUrl:      c.Config.Bunny.CDNUrl,
	}
	c.BunnyStorage = storage.NewBunnyStorage(bunnyConfig)
	log.Println("✓ Bunny Storage initialized")

	return nil
}

func (c *Container) initRepositories() error {
	c.UserRepository = postgres.NewUserRepository(c.DB)
	c.TaskRepository = postgres.NewTaskRepository(c.DB)
	c.FileRepository = postgres.NewFileRepository(c.DB)
	c.JobRepository = postgres.NewJobRepository(c.DB)
	c.ForumRepository = postgres.NewForumRepository(c.DB)
	c.TopicRepository = postgres.NewTopicRepository(c.DB)
	c.ReplyRepository = postgres.NewReplyRepository(c.DB)
	c.TagRepository = postgres.NewTagRepository(c.DB)
	c.VideoRepository = postgres.NewVideoRepository(c.DB)
	c.LikeRepository = postgres.NewLikeRepository(c.DB)
	c.CommentRepository = postgres.NewCommentRepository(c.DB)
	c.ShareRepository = postgres.NewShareRepository(c.DB)
	c.FollowRepository = postgres.NewFollowRepository(c.DB)
	c.NotificationRepository = postgres.NewNotificationRepository(c.DB)
	c.ReportRepository = postgres.NewReportRepository(c.DB)
	c.ActivityLogRepository = postgres.NewActivityLogRepository(c.DB)
	log.Println("✓ Repositories initialized")
	return nil
}

func (c *Container) initServices() error {
	// Notification service - must be initialized first to be used by other services
	c.NotificationService = serviceimpl.NewNotificationService(
		c.NotificationRepository,
		c.UserRepository,
		c.TopicRepository,
		c.VideoRepository,
		c.CommentRepository,
		c.ReplyRepository,
	)

	// Initialize other services with notification service where needed
	c.UserService = serviceimpl.NewUserService(c.UserRepository, c.TopicRepository, c.VideoRepository, c.FollowRepository, c.Config.JWT.Secret)
	c.TaskService = serviceimpl.NewTaskService(c.TaskRepository, c.UserRepository)
	c.FileService = serviceimpl.NewFileService(c.FileRepository, c.UserRepository, c.BunnyStorage)
	c.ForumService = serviceimpl.NewForumService(c.ForumRepository)
	c.TagService = serviceimpl.NewTagService(c.TagRepository, c.DB)
	c.TopicService = serviceimpl.NewTopicService(c.TopicRepository, c.ForumRepository, c.ReplyRepository, c.TagService)
	c.ReplyService = serviceimpl.NewReplyService(c.ReplyRepository, c.TopicRepository, c.NotificationService)
	c.VideoService = serviceimpl.NewVideoService(c.VideoRepository, c.FileRepository, c.UserRepository)
	c.LikeService = serviceimpl.NewLikeService(c.LikeRepository, c.TopicRepository, c.VideoRepository, c.ReplyRepository, c.CommentRepository, c.NotificationService)
	c.CommentService = serviceimpl.NewCommentService(c.CommentRepository, c.VideoRepository, c.UserRepository, c.NotificationService)
	c.FollowService = serviceimpl.NewFollowService(c.FollowRepository, c.UserRepository, c.NotificationService)
	c.ShareService = serviceimpl.NewShareService(c.ShareRepository, c.VideoRepository)

	// Admin services
	c.ReportService = serviceimpl.NewReportService(c.ReportRepository)
	c.AdminService = serviceimpl.NewAdminService(
		c.UserRepository,
		c.TopicRepository,
		c.VideoRepository,
		c.CommentRepository,
		c.ReportRepository,
		c.ActivityLogRepository,
		c.ForumRepository,
	)

	log.Println("✓ Services initialized")
	return nil
}

func (c *Container) initScheduler() error {
	c.EventScheduler = scheduler.NewEventScheduler()
	c.JobService = serviceimpl.NewJobService(c.JobRepository, c.EventScheduler)

	// Start the scheduler
	c.EventScheduler.Start()
	log.Println("✓ Event scheduler started")

	// Load and schedule existing active jobs
	ctx := context.Background()
	jobs, _, err := c.JobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		log.Printf("Warning: Failed to load existing jobs: %v", err)
		return nil
	}

	activeJobCount := 0
	for _, job := range jobs {
		if job.IsActive {
			err := c.EventScheduler.AddJob(job.ID.String(), job.CronExpr, func() {
				c.JobService.ExecuteJob(ctx, job)
			})
			if err != nil {
				log.Printf("Warning: Failed to schedule job %s: %v", job.Name, err)
			} else {
				activeJobCount++
			}
		}
	}

	if activeJobCount > 0 {
		log.Printf("✓ Scheduled %d active jobs", activeJobCount)
	}

	return nil
}

func (c *Container) Cleanup() error {
	log.Println("Starting cleanup...")

	// Stop scheduler
	if c.EventScheduler != nil {
		if c.EventScheduler.IsRunning() {
			c.EventScheduler.Stop()
			log.Println("✓ Event scheduler stopped")
		} else {
			log.Println("✓ Event scheduler was already stopped")
		}
	}

	// Close Redis connection
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			log.Printf("Warning: Failed to close Redis connection: %v", err)
		} else {
			log.Println("✓ Redis connection closed")
		}
	}

	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	log.Println("✓ Cleanup completed")
	return nil
}

func (c *Container) GetServices() (services.UserService, services.TaskService, services.FileService, services.JobService) {
	return c.UserService, c.TaskService, c.FileService, c.JobService
}

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetHandlerServices() *handlers.Services {
	return &handlers.Services{
		UserService:         c.UserService,
		TaskService:         c.TaskService,
		FileService:         c.FileService,
		JobService:          c.JobService,
		ForumService:        c.ForumService,
		TopicService:        c.TopicService,
		ReplyService:        c.ReplyService,
		TagService:          c.TagService,
		VideoService:        c.VideoService,
		LikeService:         c.LikeService,
		CommentService:      c.CommentService,
		ShareService:        c.ShareService,
		FollowService:       c.FollowService,
		NotificationService: c.NotificationService,
		AdminService:        c.AdminService,
		ReportService:       c.ReportService,
	}
}
