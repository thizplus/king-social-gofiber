const axios = require('axios');

// Configuration
const API_BASE_URL = 'http://localhost:3000/api/v1';

// Admin credentials for login
const ADMIN_CREDENTIALS = {
    email: 'admin@test.com',
    password: 'admin123'
};

// Default forums data
const DEFAULT_FORUMS = [
    { name: 'กล้อง', slug: 'camera', description: 'พูดคุยเกี่ยวกับกล้องถ่ายรูปและเทคนิคการถ่ายภาพ', icon: 'Camera', order: 1 },
    { name: 'วีดีโอ', slug: 'video', description: 'แชร์และพูดคุยเกี่ยวกับวีดีโอต่างๆ', icon: 'Code', order: 2 },
    { name: 'บทวิเคราะห์', slug: 'analysis', description: 'บทวิเคราะห์และความคิดเห็นต่างๆ', icon: 'BarChart2', order: 3 },
    { name: 'บ้านมือสอง', slug: 'second-hand-house', description: 'ซื้อขายบ้านมือสองและที่อยู่อาศัย', icon: 'Home', order: 4 },
    { name: 'การุศล', slug: 'charity', description: 'กิจกรรมการุศลและการช่วยเหลือสังคม', icon: 'Heart', order: 5 },
    { name: 'การศึกษา', slug: 'education', description: 'เรื่องราวเกี่ยวกับการศึกษาและการเรียนรู้', icon: 'Lightbulb', order: 6 },
    { name: 'ไทยรัฐ', slug: 'thairath', description: 'ข่าวสารและอัปเดตจากไทยรัฐ', icon: 'Newspaper', order: 7 },
    { name: 'เพลงแร็พ', slug: 'rap-music', description: 'พูดคุยเกี่ยวกับเพลงแร็พและฮิปฮอป', icon: 'Music', order: 8 },
    { name: 'ชาวบ้าน', slug: 'community', description: 'เรื่องราวของชาวบ้านและชุมชน', icon: 'Users', order: 9 },
    { name: 'ผู้ใช้ออนไลน์', slug: 'online-users', description: 'พูดคุยระหว่างผู้ใช้ออนไลน์', icon: 'UserRound', order: 10 },
    { name: 'ถนนคนเดิน', slug: 'walking-street', description: 'เรื่องราวจากถนนคนเดินและตลาดนัด', icon: 'Footprints', order: 11 },
    { name: 'ภาพถ่ายงาม', slug: 'beautiful-photos', description: 'แชร์ภาพถ่ายสวยๆ และเทคนิคการถ่ายภาพ', icon: 'ImageIcon', order: 12 },
    { name: 'มือถือสมาร์ท', slug: 'smartphones', description: 'รีวิวและพูดคุยเกี่ยวกับมือถือสมาร์ทโฟน', icon: 'Smartphone', order: 13 },
    { name: 'หัวเราะ', slug: 'funny', description: 'เรื่องตลกและความบันเทิง', icon: 'Sparkles', order: 14 },
    { name: 'กินเที่ยว', slug: 'food-travel', description: 'รีวิวอาหารและสถานที่ท่องเที่ยว', icon: 'UtensilsCrossed', order: 15 },
    { name: 'กีฬาไทย', slug: 'thai-sports', description: 'ข่าวและพูดคุยเกี่ยวกับกีฬาไทย', icon: 'Target', order: 16 },
    { name: 'เกมส์เด็ด', slug: 'games', description: 'รีวิวเกมและพูดคุยเกี่ยวกับเกมส์', icon: 'Gamepad2', order: 17 },
    { name: 'จตุจักร', slug: 'chatuchak', description: 'เรื่องราวจากตลาดจตุจักรและของดี', icon: 'CircleDollarSign', order: 18 },
    { name: 'เทปไทย', slug: 'thai-media', description: 'สื่อไทยและความบันเทิงท้องถิ่น', icon: 'Film', order: 19 },
    { name: 'ตลาดบ้าน', slug: 'local-market', description: 'ตลาดท้องถิ่นและของใช้ในบ้าน', icon: 'Building', order: 20 },
    { name: 'โต๊ะเครื่องดื่ม', slug: 'drinks', description: 'พูดคุยเกี่ยวกับเครื่องดื่มและคาเฟ่', icon: 'GlassWater', order: 21 },
    { name: 'มนุษย์เงินเดือน', slug: 'salary-man', description: 'เรื่องราวของคนทำงานและการเงิน', icon: 'Luggage', order: 22 },
    { name: 'ความรัก', slug: 'love', description: 'เรื่องราวความรักและความสัมพันธ์', icon: 'HeartHandshake', order: 23 }
];

// Login function
async function login() {
    try {
        console.log('🔐 Logging in...');
        const response = await axios.post(`${API_BASE_URL}/auth/login`, ADMIN_CREDENTIALS);

        if (response.data.success) {
            console.log('✅ Login successful');
            return response.data.data.token;
        } else {
            throw new Error('Login failed: ' + response.data.message);
        }
    } catch (error) {
        console.error('❌ Login error:', error.response?.data?.message || error.message);
        throw error;
    }
}

// Create forum function
async function createForum(token, forumData) {
    try {
        const response = await axios.post(
            `${API_BASE_URL}/admin/forums`,
            forumData,
            {
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            }
        );

        if (response.data.success) {
            console.log(`✅ Created forum: ${forumData.name}`);
            return response.data.data;
        } else {
            throw new Error('Create forum failed: ' + response.data.message);
        }
    } catch (error) {
        console.error(`❌ Error creating forum "${forumData.name}":`, error.response?.data?.message || error.message);
        throw error;
    }
}

// Main function
async function createDefaultForums() {
    try {
        console.log('🚀 Starting forum creation process...');

        // Step 1: Login
        const token = await login();

        // Step 2: Create forums
        console.log(`📝 Creating ${DEFAULT_FORUMS.length} forums...`);

        for (const forumData of DEFAULT_FORUMS) {
            try {
                await createForum(token, forumData);
                // Add small delay to avoid overwhelming the server
                await new Promise(resolve => setTimeout(resolve, 100));
            } catch (error) {
                console.error(`⚠️  Skipping forum "${forumData.name}" due to error`);
            }
        }

        console.log('🎉 Forum creation process completed!');

    } catch (error) {
        console.error('💥 Fatal error:', error.message);
        process.exit(1);
    }
}

// Run the script
if (require.main === module) {
    createDefaultForums();
}

module.exports = { createDefaultForums, login, createForum };