package serviceimpl

import (
	"context"
	"errors"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo    repositories.UserRepository
	topicRepo   repositories.TopicRepository
	videoRepo   repositories.VideoRepository
	followRepo  repositories.FollowRepository
	jwtSecret   string
}

func NewUserService(userRepo repositories.UserRepository, topicRepo repositories.TopicRepository, videoRepo repositories.VideoRepository, followRepo repositories.FollowRepository, jwtSecret string) services.UserService {
	return &UserServiceImpl{
		userRepo:   userRepo,
		topicRepo:  topicRepo,
		videoRepo:  videoRepo,
		followRepo: followRepo,
		jwtSecret:  jwtSecret,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, req *dto.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	existingUser, _ = s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, req *dto.LoginRequest) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return "", nil, errors.New("account is disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := s.GenerateJWT(user)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *UserServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserServiceImpl) GetProfileWithStats(ctx context.Context, userID uuid.UUID, viewerID *uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Get topic and video counts
	topicCount, videoCount, err := s.GetUserStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check follow relationships
	var isFollowing, isFollowedBy bool
	if viewerID != nil && *viewerID != userID {
		isFollowing, _ = s.followRepo.IsFollowing(ctx, *viewerID, userID)
		isFollowedBy, _ = s.followRepo.IsFollowing(ctx, userID, *viewerID)
	}

	return dto.UserToUserResponse(user, topicCount, videoCount, isFollowing, isFollowedBy), nil
}

func (s *UserServiceImpl) GetUserStats(ctx context.Context, userID uuid.UUID) (topicCount, videoCount int, err error) {
	// Count topics
	topicCountInt64, err := s.topicRepo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	topicCount = int(topicCountInt64)

	// Count videos
	videoCountInt64, err := s.videoRepo.CountByUserID(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	videoCount = int(videoCountInt64)

	return topicCount, videoCount, nil
}

func (s *UserServiceImpl) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Website != "" {
		user.Website = req.Website
	}
	if req.IsPrivate != nil {
		user.IsPrivate = *req.IsPrivate
	}

	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, offset, limit int) ([]*models.User, int64, error) {
	users, err := s.userRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (s *UserServiceImpl) GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserServiceImpl) ValidateJWT(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user ID in token")
		}

		user, err := s.userRepo.GetByID(context.Background(), userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}
