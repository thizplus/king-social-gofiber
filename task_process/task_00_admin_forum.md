# Task 00: Admin Forum Management (จัดการกระดาน)

## 📋 ภาพรวม
Admin สร้างและจัดการกระดาน (Forum/Board) สำหรับให้ User สร้างกระทู้

## 🎯 ความสำคัญ
⭐⭐⭐ **สำคัญมาก - ต้องทำก่อนทุก Task!**
เพราะ User ต้องเลือกกระดานก่อนสร้างกระทู้

## ⏱️ ระยะเวลา
**1 วัน**

---

## 📦 สิ่งที่ต้องสร้าง

### 1. Model (domain/models/)

#### 1.1 สร้าง `forum.go`
```go
package models

import (
	"time"
	"github.com/google/uuid"
)

type Forum struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `gorm:"type:varchar(100);not null"`
	Slug        string    `gorm:"type:varchar(100);uniqueIndex;not null"` // URL-friendly
	Description string    `gorm:"type:text;not null"`
	Icon        string    `gorm:"type:varchar(500)"` // URL ของไอคอน
	Order       int       `gorm:"default:0"`         // ลำดับการแสดงผล
	IsActive    bool      `gorm:"default:true"`
	TopicCount  int       `gorm:"default:0"`         // จำนวนกระทู้
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null"` // Admin ID
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Relations
	Admin  User    `gorm:"foreignKey:CreatedBy"`
	Topics []Topic `gorm:"foreignKey:ForumID"` // จะสร้างใน Task 01
}

func (Forum) TableName() string {
	return "forums"
}
```

---

### 2. DTOs (domain/dto/)

#### 2.1 สร้าง `forum.go`
```go
package dto

import (
	"time"
	"github.com/google/uuid"
)

// Request DTOs
type CreateForumRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"required,min=3,max=100,lowercase,alphanum"`
	Description string `json:"description" validate:"required,min=10,max=500"`
	Icon        string `json:"icon" validate:"omitempty,url"`
	Order       int    `json:"order" validate:"min=0"`
}

type UpdateForumRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,min=10,max=500"`
	Icon        string `json:"icon" validate:"omitempty,url"`
	Order       int    `json:"order" validate:"omitempty,min=0"`
	IsActive    *bool  `json:"isActive"` // pointer เพื่อรองรับ true/false
}

type ReorderForumsRequest struct {
	ForumOrders []ForumOrder `json:"forumOrders" validate:"required,min=1"`
}

type ForumOrder struct {
	ID    uuid.UUID `json:"id" validate:"required,uuid"`
	Order int       `json:"order" validate:"min=0"`
}

// Response DTOs
type ForumResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Order       int       `json:"order"`
	IsActive    bool      `json:"isActive"`
	TopicCount  int       `json:"topicCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ForumListResponse struct {
	Forums []ForumResponse `json:"forums"`
	Meta   PaginationMeta  `json:"meta"`
}
```

#### 2.2 อัปเดต `mappers.go`
```go
// เพิ่ม function ใน mappers.go
func ForumToForumResponse(forum *models.Forum) *ForumResponse {
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
```

---

### 3. Repository Interface (domain/repositories/)

#### 3.1 สร้าง `forum_repository.go`
```go
package repositories

import (
	"context"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ForumRepository interface {
	Create(ctx context.Context, forum *models.Forum) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Forum, error)
	GetBySlug(ctx context.Context, slug string) (*models.Forum, error)
	GetAll(ctx context.Context, includeInactive bool) ([]*models.Forum, error)
	GetActive(ctx context.Context) ([]*models.Forum, error)
	Update(ctx context.Context, id uuid.UUID, forum *models.Forum) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateOrder(ctx context.Context, id uuid.UUID, order int) error
	IncrementTopicCount(ctx context.Context, id uuid.UUID) error
	DecrementTopicCount(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
}
```

---

### 4. Repository Implementation (infrastructure/postgres/)

#### 4.1 สร้าง `forum_repository_impl.go`
```go
package postgres

import (
	"context"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ForumRepositoryImpl struct {
	db *gorm.DB
}

func NewForumRepository(db *gorm.DB) repositories.ForumRepository {
	return &ForumRepositoryImpl{db: db}
}

func (r *ForumRepositoryImpl) Create(ctx context.Context, forum *models.Forum) error {
	return r.db.WithContext(ctx).Create(forum).Error
}

func (r *ForumRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Forum, error) {
	var forum models.Forum
	err := r.db.WithContext(ctx).
		Preload("Admin").
		Where("id = ?", id).
		First(&forum).Error
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (r *ForumRepositoryImpl) GetBySlug(ctx context.Context, slug string) (*models.Forum, error) {
	var forum models.Forum
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&forum).Error
	if err != nil {
		return nil, err
	}
	return &forum, nil
}

func (r *ForumRepositoryImpl) GetAll(ctx context.Context, includeInactive bool) ([]*models.Forum, error) {
	var forums []*models.Forum
	query := r.db.WithContext(ctx)

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	err := query.Order("`order` ASC, created_at ASC").Find(&forums).Error
	return forums, err
}

func (r *ForumRepositoryImpl) GetActive(ctx context.Context) ([]*models.Forum, error) {
	return r.GetAll(ctx, false)
}

func (r *ForumRepositoryImpl) Update(ctx context.Context, id uuid.UUID, forum *models.Forum) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Updates(forum).Error
}

func (r *ForumRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.Forum{}, "id = ?", id).Error
}

func (r *ForumRepositoryImpl) UpdateOrder(ctx context.Context, id uuid.UUID, order int) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		Update("order", order).Error
}

func (r *ForumRepositoryImpl) IncrementTopicCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		UpdateColumn("topic_count", gorm.Expr("topic_count + ?", 1)).Error
}

func (r *ForumRepositoryImpl) DecrementTopicCount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Where("id = ?", id).
		Where("topic_count > ?", 0).
		UpdateColumn("topic_count", gorm.Expr("topic_count - ?", 1)).Error
}

func (r *ForumRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Forum{}).
		Count(&count).Error
	return count, err
}
```

---

### 5. Service Interface (domain/services/)

#### 5.1 สร้าง `forum_service.go`
```go
package services

import (
	"context"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"github.com/google/uuid"
)

type ForumService interface {
	// Admin Actions
	CreateForum(ctx context.Context, adminID uuid.UUID, req *dto.CreateForumRequest) (*models.Forum, error)
	UpdateForum(ctx context.Context, forumID uuid.UUID, req *dto.UpdateForumRequest) (*models.Forum, error)
	DeleteForum(ctx context.Context, forumID uuid.UUID) error
	ReorderForums(ctx context.Context, req *dto.ReorderForumsRequest) error
	GetAllForums(ctx context.Context, includeInactive bool) ([]*dto.ForumResponse, error)

	// Public Actions
	GetActiveForums(ctx context.Context) ([]*dto.ForumResponse, error)
	GetForumByID(ctx context.Context, forumID uuid.UUID) (*dto.ForumResponse, error)
	GetForumBySlug(ctx context.Context, slug string) (*dto.ForumResponse, error)
}
```

---

### 6. Service Implementation (application/serviceimpl/)

#### 6.1 สร้าง `forum_service_impl.go`
```go
package serviceimpl

import (
	"context"
	"errors"
	"time"
	"gofiber-social/domain/dto"
	"gofiber-social/domain/models"
	"gofiber-social/domain/repositories"
	"gofiber-social/domain/services"
	"github.com/google/uuid"
)

type ForumServiceImpl struct {
	forumRepo repositories.ForumRepository
}

func NewForumService(forumRepo repositories.ForumRepository) services.ForumService {
	return &ForumServiceImpl{
		forumRepo: forumRepo,
	}
}

func (s *ForumServiceImpl) CreateForum(ctx context.Context, adminID uuid.UUID, req *dto.CreateForumRequest) (*models.Forum, error) {
	// ตรวจสอบว่า slug ซ้ำหรือไม่
	existing, _ := s.forumRepo.GetBySlug(ctx, req.Slug)
	if existing != nil {
		return nil, errors.New("forum slug already exists")
	}

	forum := &models.Forum{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		Order:       req.Order,
		IsActive:    true,
		TopicCount:  0,
		CreatedBy:   adminID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.forumRepo.Create(ctx, forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (s *ForumServiceImpl) UpdateForum(ctx context.Context, forumID uuid.UUID, req *dto.UpdateForumRequest) (*models.Forum, error) {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	// Update only provided fields
	if req.Name != "" {
		forum.Name = req.Name
	}
	if req.Description != "" {
		forum.Description = req.Description
	}
	if req.Icon != "" {
		forum.Icon = req.Icon
	}
	if req.Order >= 0 {
		forum.Order = req.Order
	}
	if req.IsActive != nil {
		forum.IsActive = *req.IsActive
	}

	forum.UpdatedAt = time.Now()

	if err := s.forumRepo.Update(ctx, forumID, forum); err != nil {
		return nil, err
	}

	return forum, nil
}

func (s *ForumServiceImpl) DeleteForum(ctx context.Context, forumID uuid.UUID) error {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return errors.New("forum not found")
	}

	// ตรวจสอบว่ามีกระทู้อยู่หรือไม่
	if forum.TopicCount > 0 {
		return errors.New("cannot delete forum with existing topics")
	}

	return s.forumRepo.Delete(ctx, forumID)
}

func (s *ForumServiceImpl) ReorderForums(ctx context.Context, req *dto.ReorderForumsRequest) error {
	for _, item := range req.ForumOrders {
		if err := s.forumRepo.UpdateOrder(ctx, item.ID, item.Order); err != nil {
			return err
		}
	}
	return nil
}

func (s *ForumServiceImpl) GetAllForums(ctx context.Context, includeInactive bool) ([]*dto.ForumResponse, error) {
	forums, err := s.forumRepo.GetAll(ctx, includeInactive)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ForumResponse, len(forums))
	for i, forum := range forums {
		responses[i] = dto.ForumToForumResponse(forum)
	}

	return responses, nil
}

func (s *ForumServiceImpl) GetActiveForums(ctx context.Context) ([]*dto.ForumResponse, error) {
	return s.GetAllForums(ctx, false)
}

func (s *ForumServiceImpl) GetForumByID(ctx context.Context, forumID uuid.UUID) (*dto.ForumResponse, error) {
	forum, err := s.forumRepo.GetByID(ctx, forumID)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	return dto.ForumToForumResponse(forum), nil
}

func (s *ForumServiceImpl) GetForumBySlug(ctx context.Context, slug string) (*dto.ForumResponse, error) {
	forum, err := s.forumRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("forum not found")
	}

	return dto.ForumToForumResponse(forum), nil
}
```

---

### 7. Handler (interfaces/api/handlers/)

#### 7.1 สร้าง `forum_handler.go`
```go
package handlers

import (
	"gofiber-social/domain/dto"
	"gofiber-social/domain/services"
	"gofiber-social/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ForumHandler struct {
	forumService services.ForumService
}

func NewForumHandler(forumService services.ForumService) *ForumHandler {
	return &ForumHandler{forumService: forumService}
}

// Admin Handlers
func (h *ForumHandler) CreateForum(c *fiber.Ctx) error {
	admin, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateForumRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	forum, err := h.forumService.CreateForum(c.Context(), admin.ID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create forum", err)
	}

	return utils.SuccessResponse(c, "Forum created successfully", dto.ForumToForumResponse(forum))
}

func (h *ForumHandler) UpdateForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	var req dto.UpdateForumRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	forum, err := h.forumService.UpdateForum(c.Context(), forumID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update forum", err)
	}

	return utils.SuccessResponse(c, "Forum updated successfully", dto.ForumToForumResponse(forum))
}

func (h *ForumHandler) DeleteForum(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	if err := h.forumService.DeleteForum(c.Context(), forumID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete forum", err)
	}

	return utils.SuccessResponse(c, "Forum deleted successfully", nil)
}

func (h *ForumHandler) GetAllForums(c *fiber.Ctx) error {
	includeInactive := c.Query("includeInactive", "false") == "true"

	forums, err := h.forumService.GetAllForums(c.Context(), includeInactive)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get forums", err)
	}

	return utils.SuccessResponse(c, "Forums retrieved successfully", fiber.Map{
		"forums": forums,
		"total":  len(forums),
	})
}

func (h *ForumHandler) ReorderForums(c *fiber.Ctx) error {
	var req dto.ReorderForumsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := h.forumService.ReorderForums(c.Context(), &req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to reorder forums", err)
	}

	return utils.SuccessResponse(c, "Forums reordered successfully", nil)
}

// Public Handlers
func (h *ForumHandler) GetActiveForums(c *fiber.Ctx) error {
	forums, err := h.forumService.GetActiveForums(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get forums", err)
	}

	return utils.SuccessResponse(c, "Forums retrieved successfully", fiber.Map{
		"forums": forums,
		"total":  len(forums),
	})
}

func (h *ForumHandler) GetForumByID(c *fiber.Ctx) error {
	forumID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid forum ID")
	}

	forum, err := h.forumService.GetForumByID(c.Context(), forumID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Forum not found", err)
	}

	return utils.SuccessResponse(c, "Forum retrieved successfully", forum)
}

func (h *ForumHandler) GetForumBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	forum, err := h.forumService.GetForumBySlug(c.Context(), slug)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Forum not found", err)
	}

	return utils.SuccessResponse(c, "Forum retrieved successfully", forum)
}
```

---

### 8. Routes (interfaces/api/routes/)

#### 8.1 สร้าง `forum_routes.go`
```go
package routes

import (
	"gofiber-social/interfaces/api/handlers"
	"gofiber-social/interfaces/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupForumRoutes(api fiber.Router, h *handlers.Handlers) {
	// Admin Routes
	adminForums := api.Group("/admin/forums")
	adminForums.Use(middleware.Protected())
	adminForums.Use(middleware.AdminOnly())

	adminForums.Post("/", h.ForumHandler.CreateForum)
	adminForums.Get("/", h.ForumHandler.GetAllForums)
	adminForums.Put("/:id", h.ForumHandler.UpdateForum)
	adminForums.Delete("/:id", h.ForumHandler.DeleteForum)
	adminForums.Put("/reorder", h.ForumHandler.ReorderForums)

	// Public Routes
	forums := api.Group("/forums")

	forums.Get("/", h.ForumHandler.GetActiveForums)
	forums.Get("/:id", h.ForumHandler.GetForumByID)
	forums.Get("/slug/:slug", h.ForumHandler.GetForumBySlug)
}
```

---

### 9. อัปเดต Container (pkg/di/container.go)

```go
// เพิ่มใน Container struct
type Container struct {
	// ... existing fields ...
	ForumRepository repositories.ForumRepository
	ForumService    services.ForumService
}

// เพิ่มใน initRepositories()
func (c *Container) initRepositories() error {
	// ... existing code ...
	c.ForumRepository = postgres.NewForumRepository(c.DB)
	return nil
}

// เพิ่มใน initServices()
func (c *Container) initServices() error {
	// ... existing code ...
	c.ForumService = serviceimpl.NewForumService(c.ForumRepository)
	return nil
}

// อัปเดต GetHandlerServices()
func (c *Container) GetHandlerServices() *handlers.Services {
	return &handlers.Services{
		// ... existing services ...
		ForumService: c.ForumService,
	}
}
```

---

### 10. อัปเดต Handlers (interfaces/api/handlers/handlers.go)

```go
type Services struct {
	// ... existing services ...
	ForumService services.ForumService
}

type Handlers struct {
	// ... existing handlers ...
	ForumHandler *ForumHandler
}

func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		// ... existing handlers ...
		ForumHandler: NewForumHandler(services.ForumService),
	}
}
```

---

### 11. อัปเดต Routes (interfaces/api/routes/routes.go)

```go
func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	SetupHealthRoutes(app)

	api := app.Group("/api/v1")

	SetupAuthRoutes(api, h)
	SetupUserRoutes(api, h)
	SetupForumRoutes(api, h) // เพิ่มบรรทัดนี้
	// ... existing routes ...

	SetupWebSocketRoutes(app)
}
```

---

### 12. อัปเดต Migration (infrastructure/postgres/database.go)

```go
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Forum{}, // เพิ่มบรรทัดนี้
		// ... existing models ...
	)
}
```

---

## ✅ Checklist

- [ ] สร้าง `forum.go` model
- [ ] สร้าง `forum.go` DTO
- [ ] อัปเดต `mappers.go`
- [ ] สร้าง `forum_repository.go` interface
- [ ] สร้าง `forum_repository_impl.go`
- [ ] สร้าง `forum_service.go` interface
- [ ] สร้าง `forum_service_impl.go`
- [ ] สร้าง `forum_handler.go`
- [ ] สร้าง `forum_routes.go`
- [ ] อัปเดต `container.go`
- [ ] อัปเดต `handlers.go`
- [ ] อัปเดต `routes.go`
- [ ] อัปเดต `database.go` (migration)
- [ ] ทดสอบ Admin สร้างกระดาน
- [ ] ทดสอบ User ดูกระดาน

---

## 🧪 Testing Guide

### 1. Admin สร้างกระดาน
```bash
POST /api/v1/admin/forums
Authorization: Bearer {admin-token}
Content-Type: application/json

{
  "name": "เทคโนโลยี",
  "slug": "technology",
  "description": "พูดคุยเรื่องเทคโนโลยีและแก็ดเจ็ต",
  "icon": "🔧",
  "order": 1
}
```

### 2. Admin ดูกระดานทั้งหมด
```bash
GET /api/v1/admin/forums?includeInactive=true
Authorization: Bearer {admin-token}
```

### 3. Admin แก้ไขกระดาน
```bash
PUT /api/v1/admin/forums/{forumId}
Authorization: Bearer {admin-token}
Content-Type: application/json

{
  "name": "เทคโนโลยีและนวัตกรรม",
  "description": "พูดคุยเรื่องเทคโนโลยี แก็ดเจ็ต และนวัตกรรม"
}
```

### 4. Admin เปิด/ปิดกระดาน
```bash
PUT /api/v1/admin/forums/{forumId}
Authorization: Bearer {admin-token}
Content-Type: application/json

{
  "isActive": false
}
```

### 5. Admin เรียงลำดับกระดาน
```bash
PUT /api/v1/admin/forums/reorder
Authorization: Bearer {admin-token}
Content-Type: application/json

{
  "forumOrders": [
    { "id": "uuid-1", "order": 1 },
    { "id": "uuid-2", "order": 2 },
    { "id": "uuid-3", "order": 3 }
  ]
}
```

### 6. User ดูกระดานที่เปิดใช้งาน
```bash
GET /api/v1/forums
```

### 7. User ดูกระดานตาม ID
```bash
GET /api/v1/forums/{forumId}
```

### 8. User ดูกระดานตาม Slug
```bash
GET /api/v1/forums/slug/technology
```

---

## 📝 Notes

### ข้อควรระวัง:
1. **Slug ต้องไม่ซ้ำ** - ตรวจสอบก่อนสร้าง
2. **ไม่สามารถลบกระดานที่มีกระทู้** - ป้องกันการสูญหายของข้อมูล
3. **Order เริ่มจาก 0** - เรียงจากน้อยไปมาก
4. **IsActive = false** - กระดานจะไม่แสดงให้ User เห็น

### ตัวอย่างกระดาน:
```
Order 1: 🏠 ทั่วไป (general)
Order 2: 🔧 เทคโนโลยี (technology)
Order 3: 🎬 บันเทิง (entertainment)
Order 4: ⚽ กีฬา (sports)
Order 5: 🍔 อาหาร (food)
```

---

## 🎯 Expected Result

หลังทำ Task 00 เสร็จ:
- ✅ Admin สร้างกระดานได้
- ✅ Admin จัดการกระดานได้ (CRUD)
- ✅ Admin เรียงลำดับได้
- ✅ User ดูกระดานที่เปิดใช้งานได้
- ✅ พร้อมสำหรับ Task 01 (สร้างกระทู้)

---

**ระยะเวลาโดยประมาณ:** 1 วัน

**Next Task:** Task 01 - Topic & Reply System
