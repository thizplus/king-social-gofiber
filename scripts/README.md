# Forum Management Scripts

Scripts สำหรับจัดการข้อมูล forum

## การติดตั้ง

1. เข้าไปในโฟลเดอร์ scripts:
```bash
cd scripts
```

2. ติดตั้ง dependencies:
```bash
npm install
```

## การใช้งาน

### สร้าง Default Forums

1. ให้แน่ใจว่า server กำลัง run อยู่ที่ port 3000
2. ให้แน่ใจว่ามี admin user พร้อม credentials:
   - Email: `admin@test.com`
   - Password: `admin123`

3. รัน script:
```bash
npm run create-forums
```

หรือ:
```bash
node create_default_forums.js
```

### ปรับแต่ง

หากต้องการเปลี่ยน credentials หรือ API URL สามารถแก้ไขในไฟล์ `create_default_forums.js`:

```javascript
// Admin credentials for login
const ADMIN_CREDENTIALS = {
    email: 'your-admin@email.com',
    password: 'your-password'
};

// API Base URL
const API_BASE_URL = 'http://localhost:3000/api/v1';
```

## Forums ที่จะถูกสร้าง

Script จะสร้าง 23 forums ดังนี้:

1. กล้อง (camera)
2. วีดีโอ (video)
3. บทวิเคราะห์ (analysis)
4. บ้านมือสอง (second-hand-house)
5. การุศล (charity)
6. การศึกษา (education)
7. ไทยรัฐ (thairath)
8. เพลงแร็พ (rap-music)
9. ชาวบ้าน (community)
10. ผู้ใช้ออนไลน์ (online-users)
11. ถนนคนเดิน (walking-street)
12. ภาพถ่ายงาม (beautiful-photos)
13. มือถือสมาร์ท (smartphones)
14. หัวเราะ (funny)
15. กินเที่ยว (food-travel)
16. กีฬาไทย (thai-sports)
17. เกมส์เด็ด (games)
18. จตุจักร (chatuchak)
19. เทปไทย (thai-media)
20. ตลาดบ้าน (local-market)
21. โต๊ะเครื่องดื่ม (drinks)
22. มนุษย์เงินเดือน (salary-man)
23. ความรัก (love)

## หมายเหตุ

- Script จะ login เพื่อขอ JWT token ก่อนสร้าง forums
- หาก forum ใดสร้างไม่สำเร็จ script จะข้ามไปยัง forum ถัดไป
- มี delay 100ms ระหว่างการสร้างแต่ละ forum เพื่อไม่ให้ server โหลดหนัก