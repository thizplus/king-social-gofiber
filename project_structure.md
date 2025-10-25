# โครงสร้างระบบ Social Platform - User & Admin

## 📋 ภาพรวมระบบ
ระบบนี้แบ่งเป็น 2 ส่วนหลัก: **User System** และ **Admin System** ใช้งานในเว็บเดียวกัน แยก routes ตาม role-based authentication

---

## 🎭 Role Types
```go
const (
    RoleUser  = "user"
    RoleAdmin = "admin"
)
```

---

## 👤 USER SYSTEM - สิทธิ์และฟีเจอร์

### 1. Authentication & Profile
- ✅ สมัครสมาชิก (Register)
- ✅ เข้าสู่ระบบ (Login)
- ✅ ออกจากระบบ (Logout)
- ✅ รีเซ็ตรหัสผ่าน (Forgot Password)
- ✅ ดูโปรไฟล์ตัวเอง (View Own Profile)
- ✅ แก้ไขโปรไฟล์ (Update Profile)
- ✅ เปลี่ยนรหัสผ่าน (Change Password)
- ✅ อัปโหลดรูปโปรไฟล์ (Upload Avatar)

### 2. Social Features
- ✅ ดูโปรไฟล์ผู้ใช้คนอื่น (View Other Users)
- ✅ ติดตามผู้ใช้ (Follow Users)
- ✅ เลิกติดตามผู้ใช้ (Unfollow Users)
- ✅ ดูรายชื่อผู้ติดตาม (View Followers)
- ✅ ดูรายชื่อที่กำลังติดตาม (View Following)
- ✅ ค้นหาผู้ใช้ (Search Users)

### 3. Post Management
- ✅ สร้างโพสต์ (Create Post)
- ✅ แก้ไขโพสต์ตัวเอง (Update Own Post)
- ✅ ลบโพสต์ตัวเอง (Delete Own Post)
- ✅ ดูโพสต์ทั้งหมด (View All Posts / Feed)
- ✅ ดูโพสต์ตัวเอง (View Own Posts)
- ✅ อัปโหลดรูปภาพ/วิดีโอในโพสต์ (Upload Media)

### 4. Interaction
- ✅ กดไลค์โพสต์ (Like Post)
- ✅ ยกเลิกไลค์ (Unlike Post)
- ✅ คอมเมนต์โพสต์ (Comment on Post)
- ✅ แก้ไขคอมเมนต์ตัวเอง (Edit Own Comment)
- ✅ ลบคอมเมนต์ตัวเอง (Delete Own Comment)
- ✅ ตอบกลับคอมเมนต์ (Reply to Comment)
- ✅ แชร์โพสต์ (Share/Repost)

### 5. Messaging
- ✅ ส่งข้อความถึงผู้ใช้อื่น (Send Message)
- ✅ อ่านข้อความ (Read Messages)
- ✅ ดูประวัติแชท (View Chat History)
- ✅ ลบข้อความตัวเอง (Delete Own Messages)
- ✅ ส่งไฟล์/รูปภาพในแชท (Send Attachments)

### 6. Notifications
- ✅ ดูการแจ้งเตือน (View Notifications)
- ✅ ทำเครื่องหมายว่าอ่านแล้ว (Mark as Read)
- ✅ ตั้งค่าการแจ้งเตือน (Notification Settings)

### 7. Content Discovery
- ✅ ดู Trending Posts (View Trending)
- ✅ ค้นหาโพสต์ (Search Posts)
- ✅ ค้นหาแฮชแท็ก (Search Hashtags)
- ✅ ดูโพสต์ตามแฮชแท็ก (View Posts by Hashtag)

### 8. Privacy & Settings
- ✅ ตั้งค่าความเป็นส่วนตัว (Privacy Settings)
- ✅ บล็อกผู้ใช้ (Block Users)
- ✅ ยกเลิกบล็อก (Unblock Users)
- ✅ รายงานโพสต์ (Report Post)
- ✅ รายงานผู้ใช้ (Report User)

---

## 🛡️ ADMIN SYSTEM - สิทธิ์และฟีเจอร์

### 1. Dashboard & Analytics
- ✅ ดู Dashboard สรุปภาพรวม (View Dashboard)
- ✅ สถิติผู้ใช้ทั้งหมด (Total Users Stats)
- ✅ สถิติโพสต์รายวัน/เดือน (Posts Statistics)
- ✅ กราฟการเติบโตของผู้ใช้ (User Growth Chart)
- ✅ กราฟการใช้งานระบบ (System Usage Analytics)
- ✅ รายงานการใช้งาน Storage (Storage Usage)

### 2. User Management
- ✅ ดูรายชื่อผู้ใช้ทั้งหมด (View All Users)
- ✅ ค้นหาผู้ใช้ (Search Users)
- ✅ ดูรายละเอียดผู้ใช้ (View User Details)
- ✅ แก้ไขข้อมูลผู้ใช้ (Edit User Info)
- ✅ ปิดการใช้งานบัญชี (Suspend/Deactivate Account)
- ✅ เปิดการใช้งานบัญชี (Activate Account)
- ✅ ลบบัญชีผู้ใช้ (Delete User Account)
- ✅ เปลี่ยน Role ผู้ใช้ (Change User Role)
- ✅ ดูประวัติการกระทำของผู้ใช้ (View User Activity Log)

### 3. Content Moderation
- ✅ ดูโพสต์ทั้งหมด (View All Posts)
- ✅ ค้นหาโพสต์ (Search Posts)
- ✅ ดูรายงานโพสต์ (View Reported Posts)
- ✅ อนุมัติ/ปฏิเสธโพสต์ (Approve/Reject Posts)
- ✅ ลบโพสต์ (Delete Posts)
- ✅ แก้ไขโพสต์ (Edit Posts)
- ✅ ซ่อนโพสต์ (Hide Posts)
- ✅ ปักหมุดโพสต์ (Pin Posts)

### 4. Comment Management
- ✅ ดูคอมเมนต์ทั้งหมด (View All Comments)
- ✅ ดูรายงานคอมเมนต์ (View Reported Comments)
- ✅ ลบคอมเมนต์ (Delete Comments)
- ✅ แก้ไขคอมเมนต์ (Edit Comments)
- ✅ ซ่อนคอมเมนต์ (Hide Comments)

### 5. Report Management
- ✅ ดูรายงานทั้งหมด (View All Reports)
- ✅ ดูรายงานโพสต์ (View Post Reports)
- ✅ ดูรายงานผู้ใช้ (View User Reports)
- ✅ ดูรายงานคอมเมนต์ (View Comment Reports)
- ✅ จัดการรายงาน (Process Reports)
- ✅ อนุมัติ/ปฏิเสธรายงาน (Approve/Reject Reports)
- ✅ ปิดรายงาน (Close Reports)

### 6. System Settings
- ✅ ตั้งค่าระบบทั่วไป (General Settings)
- ✅ จัดการ API Keys (Manage API Keys)
- ✅ ตั้งค่า Email Templates (Email Template Settings)
- ✅ ตั้งค่าการแจ้งเตือน (Notification Settings)
- ✅ จัดการ Storage (Storage Management)
- ✅ ตั้งค่าความปลอดภัย (Security Settings)

### 7. Admin Management
- ✅ ดูรายชื่อ Admin (View All Admins)
- ✅ เพิ่ม Admin ใหม่ (Add New Admin)
- ✅ แก้ไขสิทธิ์ Admin (Edit Admin Permissions)
- ✅ ลบ Admin (Remove Admin)
- ✅ ดู Activity Log ของ Admin (View Admin Activity Log)

### 8. Logs & Monitoring
- ✅ ดู System Logs (View System Logs)
- ✅ ดู Error Logs (View Error Logs)
- ✅ ดู Activity Logs (View Activity Logs)
- ✅ ดู Login History (View Login History)
- ✅ ติดตามการใช้งานแบบ Real-time (Real-time Monitoring)

---

## 🗂️ โครงสร้าง Routes แนะนำ

### User Routes
```
/api/v1/user/
├── auth/
│   ├── POST   /register
│   ├── POST   /login
│   ├── POST   /logout
│   ├── POST   /forgot-password
│   └── POST   /reset-password
│
├── profile/
│   ├── GET    /me
│   ├── PUT    /me
│   ├── PUT    /me/password
│   ├── POST   /me/avatar
│   └── GET    /:userId
│
├── social/
│   ├── POST   /follow/:userId
│   ├── DELETE /unfollow/:userId
│   ├── GET    /followers
│   ├── GET    /following
│   └── GET    /search
│
├── posts/
│   ├── GET    /
│   ├── GET    /feed
│   ├── POST   /
│   ├── GET    /:postId
│   ├── PUT    /:postId
│   ├── DELETE /:postId
│   ├── POST   /:postId/like
│   ├── DELETE /:postId/unlike
│   └── GET    /trending
│
├── comments/
│   ├── POST   /posts/:postId/comments
│   ├── GET    /posts/:postId/comments
│   ├── PUT    /comments/:commentId
│   ├── DELETE /comments/:commentId
│   └── POST   /comments/:commentId/reply
│
├── messages/
│   ├── GET    /
│   ├── GET    /conversations
│   ├── POST   /send
│   ├── GET    /conversation/:userId
│   └── DELETE /:messageId
│
├── notifications/
│   ├── GET    /
│   ├── PUT    /:id/read
│   └── PUT    /mark-all-read
│
└── settings/
    ├── PUT    /privacy
    ├── POST   /block/:userId
    ├── DELETE /unblock/:userId
    ├── POST   /report/post/:postId
    └── POST   /report/user/:userId
```

### Admin Routes
```
/api/v1/admin/
├── auth/
│   ├── POST   /login
│   └── POST   /logout
│
├── dashboard/
│   ├── GET    /stats
│   ├── GET    /analytics
│   ├── GET    /user-growth
│   └── GET    /system-usage
│
├── users/
│   ├── GET    /
│   ├── GET    /search
│   ├── GET    /:userId
│   ├── PUT    /:userId
│   ├── DELETE /:userId
│   ├── PUT    /:userId/suspend
│   ├── PUT    /:userId/activate
│   ├── PUT    /:userId/role
│   └── GET    /:userId/activity-log
│
├── posts/
│   ├── GET    /
│   ├── GET    /search
│   ├── GET    /:postId
│   ├── DELETE /:postId
│   ├── PUT    /:postId
│   ├── PUT    /:postId/hide
│   ├── PUT    /:postId/pin
│   └── GET    /reported
│
├── comments/
│   ├── GET    /
│   ├── GET    /reported
│   ├── DELETE /:commentId
│   ├── PUT    /:commentId
│   └── PUT    /:commentId/hide
│
├── reports/
│   ├── GET    /
│   ├── GET    /posts
│   ├── GET    /users
│   ├── GET    /comments
│   ├── PUT    /:reportId/process
│   ├── PUT    /:reportId/approve
│   ├── PUT    /:reportId/reject
│   └── PUT    /:reportId/close
│
├── admins/
│   ├── GET    /
│   ├── POST   /
│   ├── PUT    /:adminId
│   ├── DELETE /:adminId
│   └── GET    /activity-log
│
├── settings/
│   ├── GET    /
│   ├── PUT    /general
│   ├── PUT    /email-templates
│   ├── PUT    /notifications
│   ├── PUT    /security
│   └── GET    /storage
│
└── logs/
    ├── GET    /system
    ├── GET    /errors
    ├── GET    /activity
    ├── GET    /login-history
    └── GET    /realtime
```

---

## 🗄️ Database Models ที่ต้องเพิ่ม

### 1. Enhanced User Model
```go
type User struct {
    ID              uuid.UUID
    Email           string
    Username        string
    Password        string
    FirstName       string
    LastName        string
    Avatar          string
    Bio             string       // เพิ่ม
    Role            string       // user, admin
    IsActive        bool
    IsSuspended     bool         // เพิ่ม
    SuspendedAt     *time.Time   // เพิ่ม
    SuspendReason   string       // เพิ่ม
    FollowersCount  int          // เพิ่ม
    FollowingCount  int          // เพิ่ม
    PostsCount      int          // เพิ่ม
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 2. Post Model
```go
type Post struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    Content         string
    MediaURLs       []string     // JSON array
    Type            string       // text, image, video
    LikesCount      int
    CommentsCount   int
    SharesCount     int
    IsHidden        bool
    IsPinned        bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 3. Comment Model
```go
type Comment struct {
    ID              uuid.UUID
    PostID          uuid.UUID
    UserID          uuid.UUID
    ParentID        *uuid.UUID   // for replies
    Content         string
    LikesCount      int
    IsHidden        bool
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 4. Follow Model
```go
type Follow struct {
    ID              uuid.UUID
    FollowerID      uuid.UUID
    FollowingID     uuid.UUID
    CreatedAt       time.Time
}
```

### 5. Like Model
```go
type Like struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    PostID          *uuid.UUID   // nullable
    CommentID       *uuid.UUID   // nullable
    CreatedAt       time.Time
}
```

### 6. Message Model
```go
type Message struct {
    ID              uuid.UUID
    SenderID        uuid.UUID
    ReceiverID      uuid.UUID
    Content         string
    AttachmentURL   string
    IsRead          bool
    CreatedAt       time.Time
}
```

### 7. Notification Model
```go
type Notification struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    Type            string       // like, comment, follow, message
    ActorID         uuid.UUID
    PostID          *uuid.UUID
    CommentID       *uuid.UUID
    Content         string
    IsRead          bool
    CreatedAt       time.Time
}
```

### 8. Report Model
```go
type Report struct {
    ID              uuid.UUID
    ReporterID      uuid.UUID
    ReportedUserID  *uuid.UUID
    ReportedPostID  *uuid.UUID
    ReportedCommentID *uuid.UUID
    Reason          string
    Description     string
    Status          string       // pending, approved, rejected, closed
    ProcessedBy     *uuid.UUID   // admin ID
    ProcessedAt     *time.Time
    CreatedAt       time.Time
}
```

### 9. BlockedUser Model
```go
type BlockedUser struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    BlockedUserID   uuid.UUID
    CreatedAt       time.Time
}
```

### 10. ActivityLog Model
```go
type ActivityLog struct {
    ID              uuid.UUID
    UserID          uuid.UUID
    Action          string
    Resource        string
    ResourceID      uuid.UUID
    IPAddress       string
    UserAgent       string
    CreatedAt       time.Time
}
```

### 11. AdminSettings Model
```go
type AdminSettings struct {
    ID              uuid.UUID
    Key             string
    Value           string
    Type            string       // string, int, bool, json
    UpdatedBy       uuid.UUID
    UpdatedAt       time.Time
}
```

---

## 🔐 Middleware Structure

### 1. Auth Middleware
```go
// interfaces/api/middleware/auth_middleware.go
func RequireAuth() fiber.Handler
func RequireRole(roles ...string) fiber.Handler
func RequireAdmin() fiber.Handler
```

### 2. Permission Middleware
```go
// interfaces/api/middleware/permission_middleware.go
func CanEditPost() fiber.Handler
func CanDeleteComment() fiber.Handler
func IsPostOwner() fiber.Handler
```

### 3. Rate Limiting
```go
// interfaces/api/middleware/rate_limit_middleware.go
func UserRateLimit() fiber.Handler
func AdminRateLimit() fiber.Handler
```

---

## 📁 File Structure แนะนำ

```
gofiber/
├── domain/
│   ├── models/
│   │   ├── user.go              ✅ มีแล้ว
│   │   ├── post.go              ❌ ต้องสร้างใหม่
│   │   ├── comment.go           ❌ ต้องสร้างใหม่
│   │   ├── follow.go            ❌ ต้องสร้างใหม่
│   │   ├── like.go              ❌ ต้องสร้างใหม่
│   │   ├── message.go           ❌ ต้องสร้างใหม่
│   │   ├── notification.go      ❌ ต้องสร้างใหม่
│   │   ├── report.go            ❌ ต้องสร้างใหม่
│   │   ├── blocked_user.go      ❌ ต้องสร้างใหม่
│   │   ├── activity_log.go      ❌ ต้องสร้างใหม่
│   │   └── admin_settings.go    ❌ ต้องสร้างใหม่
│   │
│   ├── dto/
│   │   ├── user.go              ✅ มีแล้ว
│   │   ├── post.go              ❌ ต้องสร้างใหม่
│   │   ├── comment.go           ❌ ต้องสร้างใหม่
│   │   ├── message.go           ❌ ต้องสร้างใหม่
│   │   ├── notification.go      ❌ ต้องสร้างใหม่
│   │   ├── report.go            ❌ ต้องสร้างใหม่
│   │   └── admin.go             ❌ ต้องสร้างใหม่
│   │
│   ├── repositories/
│   │   ├── user_repository.go   ✅ มีแล้ว
│   │   ├── post_repository.go   ❌ ต้องสร้างใหม่
│   │   ├── comment_repository.go ❌ ต้องสร้างใหม่
│   │   ├── follow_repository.go ❌ ต้องสร้างใหม่
│   │   ├── like_repository.go   ❌ ต้องสร้างใหม่
│   │   ├── message_repository.go ❌ ต้องสร้างใหม่
│   │   ├── notification_repository.go ❌ ต้องสร้างใหม่
│   │   ├── report_repository.go ❌ ต้องสร้างใหม่
│   │   └── activity_log_repository.go ❌ ต้องสร้างใหม่
│   │
│   └── services/
│       ├── user_service.go      ✅ มีแล้ว
│       ├── post_service.go      ❌ ต้องสร้างใหม่
│       ├── comment_service.go   ❌ ต้องสร้างใหม่
│       ├── follow_service.go    ❌ ต้องสร้างใหม่
│       ├── message_service.go   ❌ ต้องสร้างใหม่
│       ├── notification_service.go ❌ ต้องสร้างใหม่
│       ├── admin_service.go     ❌ ต้องสร้างใหม่
│       └── report_service.go    ❌ ต้องสร้างใหม่
│
├── application/
│   └── serviceimpl/
│       ├── user_service_impl.go ✅ มีแล้ว
│       ├── post_service_impl.go ❌ ต้องสร้างใหม่
│       ├── comment_service_impl.go ❌ ต้องสร้างใหม่
│       ├── follow_service_impl.go ❌ ต้องสร้างใหม่
│       ├── message_service_impl.go ❌ ต้องสร้างใหม่
│       ├── notification_service_impl.go ❌ ต้องสร้างใหม่
│       ├── admin_service_impl.go ❌ ต้องสร้างใหม่
│       └── report_service_impl.go ❌ ต้องสร้างใหม่
│
├── infrastructure/
│   └── postgres/
│       ├── user_repository_impl.go ✅ มีแล้ว
│       ├── post_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── comment_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── follow_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── like_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── message_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── notification_repository_impl.go ❌ ต้องสร้างใหม่
│       ├── report_repository_impl.go ❌ ต้องสร้างใหม่
│       └── activity_log_repository_impl.go ❌ ต้องสร้างใหม่
│
└── interfaces/
    └── api/
        ├── handlers/
        │   ├── user_handler.go   ✅ มีแล้ว (ต้องแก้ไข)
        │   ├── post_handler.go   ❌ ต้องสร้างใหม่
        │   ├── comment_handler.go ❌ ต้องสร้างใหม่
        │   ├── follow_handler.go ❌ ต้องสร้างใหม่
        │   ├── message_handler.go ❌ ต้องสร้างใหม่
        │   ├── notification_handler.go ❌ ต้องสร้างใหม่
        │   ├── admin_handler.go  ❌ ต้องสร้างใหม่
        │   └── report_handler.go ❌ ต้องสร้างใหม่
        │
        ├── routes/
        │   ├── user_routes.go    ✅ มีแล้ว (ต้องแก้ไข)
        │   ├── post_routes.go    ❌ ต้องสร้างใหม่
        │   ├── comment_routes.go ❌ ต้องสร้างใหม่
        │   ├── follow_routes.go  ❌ ต้องสร้างใหม่
        │   ├── message_routes.go ❌ ต้องสร้างใหม่
        │   ├── notification_routes.go ❌ ต้องสร้างใหม่
        │   ├── admin_routes.go   ❌ ต้องสร้างใหม่
        │   └── report_routes.go  ❌ ต้องสร้างใหม่
        │
        └── middleware/
            ├── auth_middleware.go ✅ มีแล้ว (ต้องแก้ไข)
            ├── role_middleware.go ❌ ต้องสร้างใหม่
            ├── permission_middleware.go ❌ ต้องสร้างใหม่
            └── rate_limit_middleware.go ❌ ต้องสร้างใหม่
```

---

## 🚀 Implementation Priority

### Phase 1: Foundation (สัปดาห์ที่ 1-2)
1. ✅ User Model Enhancement
2. ✅ Post, Comment, Like Models
3. ✅ Role-based Middleware
4. ✅ Basic User Routes
5. ✅ Basic Post Routes

### Phase 2: Social Features (สัปดาห์ที่ 3-4)
6. ✅ Follow System
7. ✅ Comment System
8. ✅ Like System
9. ✅ Feed Algorithm
10. ✅ Search Functionality

### Phase 3: Communication (สัปดาห์ที่ 5-6)
11. ✅ Message System
12. ✅ Notification System
13. ✅ Real-time Updates (WebSocket)

### Phase 4: Admin System (สัปดาห์ที่ 7-8)
14. ✅ Admin Dashboard
15. ✅ User Management
16. ✅ Content Moderation
17. ✅ Report System
18. ✅ Activity Logs

### Phase 5: Enhancement (สัปดาห์ที่ 9-10)
19. ✅ Analytics & Statistics
20. ✅ Privacy Settings
21. ✅ Block System
22. ✅ Advanced Search
23. ✅ System Monitoring

---

## 🔒 Security Considerations

1. **Authentication**
   - JWT Token with refresh token
   - Password hashing (bcrypt)
   - Rate limiting on login

2. **Authorization**
   - Role-based access control (RBAC)
   - Resource ownership validation
   - Admin permission levels

3. **Data Protection**
   - Input validation
   - SQL injection prevention (use parameterized queries)
   - XSS protection
   - CSRF protection

4. **Privacy**
   - User data encryption
   - Secure file upload
   - Private message encryption

---

## 📊 Monitoring & Logging

1. **Activity Logs**
   - User actions
   - Admin actions
   - System events

2. **Error Logs**
   - Application errors
   - Database errors
   - API errors

3. **Analytics**
   - User engagement
   - System performance
   - Resource usage

---

## 🎯 Next Steps

1. **ทบทวนโครงสร้างนี้** - ดูว่าตรงตามความต้องการหรือไม่
2. **เริ่มจาก Phase 1** - สร้าง Models และ Middleware พื้นฐาน
3. **ทำทีละ Feature** - อย่าพยายามทำทั้งหมดพร้อมกัน
4. **Test แต่ละส่วน** - ให้แน่ใจว่าทำงานถูกต้อง
5. **Deploy Gradually** - Deploy ทีละส่วน

---

## 📝 Notes

- ใช้ **Clean Architecture** ที่มีอยู่แล้ว
- แยก routes ตาม role ชัดเจน (`/api/v1/user/` vs `/api/v1/admin/`)
- ใช้ middleware ตรวจสอบ role ทุก route
- Admin มีสิทธิ์เกือบทุกอย่างที่ user ทำได้ + สิทธิ์พิเศษ
- เก็บ Activity Log ทุก action ของ admin
- ใช้ Soft Delete สำหรับข้อมูลสำคัญ

---

**Created:** 2025-10-01
**Version:** 1.0
**Status:** Planning Phase