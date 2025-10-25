# Tasks Updated - Social Platform (Webboard + Video)

## 📋 Overview
ระบบ Social Platform ที่รวม **Webboard (Pantip-style)** และ **Video Sharing (TikTok-style)**

---

## ✅ ที่มีอยู่แล้ว
- User Registration & Login
- JWT Authentication
- User Profile Management
- Role-based Access (user, admin)

---

## 🎯 Tasks ที่ต้องทำ (เรียงลำดับ)

### **Task 00: Admin - Forum Management** ⭐ เริ่มที่นี่
**ระยะเวลา:** 1 วัน

#### Models:
1. **Forum** - กระดาน/หมวดหมู่
   - ID, Name, Slug, Description, Icon, Order, IsActive
   - TopicCount, CreatedBy (AdminID), CreatedAt, UpdatedAt

#### Features (Admin Only):
- ✅ สร้างกระดาน (Forum)
- ✅ แก้ไขกระดาน
- ✅ ลบกระดาน
- ✅ เรียงลำดับกระดาน
- ✅ เปิด/ปิดกระดาน

#### Routes:
```
# Admin Only
POST   /api/v1/admin/forums              - สร้างกระดาน
GET    /api/v1/admin/forums              - ดูกระดานทั้งหมด (รวมที่ปิด)
PUT    /api/v1/admin/forums/:id          - แก้ไขกระดาน
DELETE /api/v1/admin/forums/:id          - ลบกระดาน
PUT    /api/v1/admin/forums/reorder      - เรียงลำดับ

# Public (User)
GET    /api/v1/forums                    - ดูกระดานที่เปิดใช้งาน
GET    /api/v1/forums/:id                - ดูกระดานและกระทู้ในนั้น
```

#### DTOs:
```go
type CreateForumRequest struct {
    Name        string `json:"name" validate:"required,min=3,max=100"`
    Slug        string `json:"slug" validate:"required,min=3,max=100,lowercase"`
    Description string `json:"description" validate:"required,min=10,max=500"`
    Icon        string `json:"icon" validate:"omitempty,url"`
    Order       int    `json:"order" validate:"min=0"`
}

type UpdateForumRequest struct {
    Name        string `json:"name" validate:"omitempty,min=3,max=100"`
    Description string `json:"description" validate:"omitempty,min=10,max=500"`
    Icon        string `json:"icon" validate:"omitempty,url"`
    Order       int    `json:"order" validate:"omitempty,min=0"`
    IsActive    *bool  `json:"isActive"`
}

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
}
```

---

### **Task 01: Topic & Reply System**
**ระยะเวลา:** 3-4 วัน

#### Models:
1. **Topic** - กระทู้
   - ID, UserID, ForumID, Title, Content, ViewCount, ReplyCount
   - IsPinned, IsLocked, CreatedAt, UpdatedAt
2. **Reply** - ตอบกระทู้
   - ID, TopicID, UserID, Content, ParentID (nested reply)
   - CreatedAt, UpdatedAt

#### Features:
- ✅ สร้างกระทู้
- ✅ ตอบกระทู้ (รองรับ nested replies)
- ✅ ดูกระทู้ (พร้อม view count)
- ✅ แก้ไข/ลบกระทู้ (เฉพาะเจ้าของ)
- ✅ กระทู้แนะนำ/ปักหมุด (Admin)

#### Routes:
```
# Topics
POST   /api/v1/topics                - สร้างกระทู้
GET    /api/v1/topics                - ดูกระทู้ทั้งหมด
GET    /api/v1/topics/:id            - ดูกระทู้
PUT    /api/v1/topics/:id            - แก้ไขกระทู้
DELETE /api/v1/topics/:id            - ลบกระทู้

# Replies
POST   /api/v1/topics/:id/replies    - ตอบกระทู้
GET    /api/v1/topics/:id/replies    - ดูคำตอบ
PUT    /api/v1/replies/:id           - แก้ไขคำตอบ
DELETE /api/v1/replies/:id           - ลบคำตอบ

# Forum Topics
GET    /api/v1/forums/:id/topics     - ดูกระทู้ในกระดาน
```

---

### **Task 02: Video Upload & Management System**
**ระยะเวลา:** 3-4 วัน

#### Models:
1. **Video** - วิดีโอ
   - ID, UserID, Title, Description, VideoURL, ThumbnailURL
   - Duration, ViewCount, LikeCount, CommentCount
   - IsActive, CreatedAt

#### Features:
- ✅ อัปโหลดวิดีโอ (รองรับ multipart upload)
- ✅ สร้าง Thumbnail อัตโนมัติ
- ✅ Video Metadata (duration, resolution)
- ✅ Video Feed/List
- ✅ แก้ไข/ลบวิดีโอ (เฉพาะเจ้าของ)

#### Storage:
- Local Storage (สำหรับ development)
- Bunny CDN / AWS S3 (สำหรับ production)

#### Routes:
```
POST   /api/v1/videos               - อัปโหลดวิดีโอ
GET    /api/v1/videos               - ดู Feed วิดีโอ
GET    /api/v1/videos/:id           - ดูวิดีโอ
PUT    /api/v1/videos/:id           - แก้ไขวิดีโอ
DELETE /api/v1/videos/:id           - ลบวิดีโอ
GET    /api/v1/videos/user/:userId  - ดูวิดีโอของ User
```

---

### **Task 03: Like, Comment & Share System**
**ระยะเวลา:** 2-3 วัน

#### Models:
1. **Like** - ไลค์
   - ID, UserID, TopicID (nullable), VideoID (nullable)
2. **Comment** - คอมเมนต์
   - ID, UserID, VideoID, Content, ParentID
3. **Share** - แชร์
   - ID, UserID, VideoID, Platform

#### Features:
- ✅ Like กระทู้/วิดีโอ
- ✅ Unlike
- ✅ Comment วิดีโอ (รองรับ nested)
- ✅ แก้ไข/ลบ Comment
- ✅ Share วิดีโอ (count)

#### Routes:
```
# Like
POST   /api/v1/forum/topics/:id/like      - ไลค์กระทู้
DELETE /api/v1/forum/topics/:id/unlike    - ยกเลิกไลค์กระทู้

POST   /api/v1/videos/:id/like            - ไลค์วิดีโอ
DELETE /api/v1/videos/:id/unlike          - ยกเลิกไลค์วิดีโอ

# Comment
POST   /api/v1/videos/:id/comments        - คอมเมนต์วิดีโอ
GET    /api/v1/videos/:id/comments        - ดูคอมเมนต์
PUT    /api/v1/videos/comments/:id        - แก้ไขคอมเมนต์
DELETE /api/v1/videos/comments/:id        - ลบคอมเมนต์

# Share
POST   /api/v1/videos/:id/share           - แชร์วิดีโอ
```

---

### **Task 04: Follow System**
**ระยะเวลา:** 2 วัน

#### Models:
1. **Follow** - ติดตาม
   - ID, FollowerID, FollowingID

#### Features:
- ✅ Follow/Unfollow User
- ✅ ดูรายชื่อ Followers
- ✅ ดูรายชื่อ Following
- ✅ แจ้งเตือนเมื่อมีคนติดตาม

#### Routes:
```
POST   /api/v1/users/:userId/follow    - ติดตาม
DELETE /api/v1/users/:userId/unfollow  - เลิกติดตาม
GET    /api/v1/users/:userId/followers - ดู Followers
GET    /api/v1/users/:userId/following - ดู Following
```

---

### **Task 05: Notification System**
**ระยะเวลา:** 2 วัน

#### Models:
1. **Notification** - การแจ้งเตือน
   - ID, UserID, Type, ActorID, ResourceID, IsRead

#### Types:
- `topic_reply` - มีคนตอบกระทู้
- `video_like` - มีคนไลค์วิดีโอ
- `video_comment` - มีคนคอมเมนต์วิดีโอ
- `new_follower` - มีคนติดตาม

#### Routes:
```
GET  /api/v1/notifications           - ดูการแจ้งเตือน
GET  /api/v1/notifications/unread    - ดูที่ยังไม่อ่าน
POST /api/v1/notifications/mark-read - ทำเครื่องหมายว่าอ่านแล้ว
```

---

### **Task 06: Admin Dashboard & Management**
**ระยะเวลา:** 3-4 วัน

#### Features:
- ✅ Dashboard (Stats, Charts)
- ✅ จัดการ Users (Suspend, Delete, Role)
- ✅ จัดการกระทู้ (Delete, Pin, Lock)
- ✅ จัดการวิดีโอ (Delete, Hide)
- ✅ Report & Moderation
- ✅ Activity Logs

#### Routes:
```
# Dashboard
GET /api/v1/admin/dashboard/stats

# User Management
GET    /api/v1/admin/users
PUT    /api/v1/admin/users/:id/suspend
DELETE /api/v1/admin/users/:id

# Forum Management
DELETE /api/v1/admin/forum/topics/:id
PUT    /api/v1/admin/forum/topics/:id/pin
PUT    /api/v1/admin/forum/topics/:id/lock

# Video Management
DELETE /api/v1/admin/videos/:id
PUT    /api/v1/admin/videos/:id/hide

# Reports
GET  /api/v1/admin/reports
POST /api/v1/admin/reports/:id/process
```

---

## 📊 Timeline Summary

| Week | Tasks | Status |
|------|-------|--------|
| 1 | Task 00: Admin Forum Management | 🏷️ |
| 1 | Task 01: Topic & Reply System | 📝 |
| 2 | Task 02: Video System | 🎬 |
| 3 | Task 03: Like/Comment/Share | ❤️ |
| 3 | Task 04: Follow System | 👥 |
| 4 | Task 05: Notification | 🔔 |
| 4 | Task 06: Admin Dashboard | 🛡️ |

**Total: ~4 สัปดาห์**

---

## 🚀 Next Steps

### ขั้นตอนที่ 1: Admin สร้างกระดานก่อน
1. ✅ ทำ Task 00: Admin Forum Management
2. ✅ Admin สร้างกระดาน (เช่น: ทั่วไป, เทคโนโลยี, บันเทิง)
3. ✅ ทดสอบ CRUD กระดาน

### ขั้นตอนที่ 2: User สร้างกระทู้
4. ✅ ทำ Task 01: Topic & Reply System
5. ✅ User เลือกกระดานและสร้างกระทู้
6. ✅ ทดสอบตั้งกระทู้และตอบกระทู้

### ขั้นตอนที่ 3: ต่อด้วย Video & Social Features
7. ✅ Task 02-06 ตามลำดับ

---

## 📝 Notes

- User เดียวกันใช้ได้ทั้งระบบกระทู้และวิดีโอ
- Like/Comment ทำงานกับทั้งกระทู้และวิดีโอ
- Admin จัดการได้ทั้ง 2 ระบบ
- Notification รวมทุกประเภท
- Follow System ใช้ร่วมกัน

---

**ค่อยๆ ทำทีละ Task นะครับ ไม่ต้องรีบ!** 🎯

แต่ละ Task จะมีไฟล์รายละเอียดแยกต่างหาก พร้อม code examples ครบถ้วน
