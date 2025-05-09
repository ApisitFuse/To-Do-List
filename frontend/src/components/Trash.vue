<template>
    <div class="max-w-xl mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">รายการงานที่ถูกลบ</h1>

        <div v-if="loading" class="text-center">กำลังโหลด...</div>
        <div v-if="error" class="text-center text-error">{{ error }}</div>

        <div v-if="!loading && trashedItems.length === 0 && !error" class="text-center text-gray-500">
            ไม่มีงานในถังขยะ
        </div>

        <div v-if="!loading && !error && trashedItemGroups.length === 0 && trashedItems.length > 0"
            class="text-center text-gray-500">
            ไม่มีงานที่จัดกลุ่มได้ (อาจจะไม่มีวันที่ถูกลบ)
        </div>

        <div v-for="group in trashedItemGroups" :key="group.key" class="mb-6">
            <h2 v-if="group.items.length > 0" class="text-xl font-semibold mb-3 text-primary">{{ group.label }}</h2>
            <div v-if="group.items.length > 0">
                <div v-for="item in group.items" :key="item.ID"
                    class="flex items-center p-3 border border-base-300 rounded-box mb-2 hover:bg-base-200 transition-colors duration-150">
                    <span class="flex-grow" :class="{ 'line-through text-gray-500': item.completed }">{{ item.title
                        }}</span>
                    <button @click="restoreTodo(item.ID)" class="btn btn-sm btn-success mr-2">
                        กู้คืน
                    </button>
                    <button @click="permanentlyDeleteTodo(item.ID)" class="btn btn-sm btn-error">
                        ลบถาวร
                    </button>
                </div>
            </div>

        </div>
        <div class="mt-6 text-center">
            <router-link to="/todo-two" class="btn btn-primary">กลับ</router-link>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';

const trashedItems = ref([]);
const loading = ref(true);
const error = ref(null);

const API_URL = 'http://localhost:8080/api/todos';

const trashedItemGroups = ref([]);

const loadTrashedItems = async () => {
    loading.value = true;
    error.value = null;
    trashedItems.value = [];
    trashedItemGroups.value = [];
    try {
        const response = await axios.get(`${API_URL}/trashed`); // Endpoint ใหม่
        if (Array.isArray(response.data)) {
            const rawItems = response.data;
            trashedItems.value = rawItems;

            const groups = {
                today: [],
                yesterday: [],
                lastWeek: [],
                older: []
            };

            const now = new Date();
            const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            const yesterdayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1);
            const sevenDaysAgoStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 7);

            rawItems.forEach(item => {
                if (!item.DeletedAt) {
                    console.warn(`Item ID ${item.ID} is missing deleted_at, placing in 'older'`);
                    groups.older.push(item); // Fallback for items without a deletion date
                    return;
                }
                const itemDeletedDate = new Date(item.DeletedAt);

                if (itemDeletedDate >= todayStart) {
                    groups.today.push(item);
                } else if (itemDeletedDate >= yesterdayStart) {
                    groups.yesterday.push(item);
                } else if (itemDeletedDate >= sevenDaysAgoStart) { // Items from 2 to 7 days ago
                    groups.lastWeek.push(item);
                } else { // Items older than 7 days
                    groups.older.push(item);
                }
            });

            // Sort items within each group by deletion date (newest first)
            for (const key in groups) {
                groups[key].sort((a, b) => new Date(b.deleted_at) - new Date(a.deleted_at));
            }

            const resultGroups = [];
            if (groups.today.length) resultGroups.push({ key: 'today', label: 'วันนี้', items: groups.today });
            if (groups.yesterday.length) resultGroups.push({ key: 'yesterday', label: 'เมื่อวาน', items: groups.yesterday });
            if (groups.lastWeek.length) resultGroups.push({ key: 'lastWeek', label: 'สัปดาห์ที่แล้ว', items: groups.lastWeek });
            if (groups.older.length) resultGroups.push({ key: 'older', label: 'หนึ่งเดือนขึ้นไป (เก่ากว่า 7 วัน)', items: groups.older });

            trashedItemGroups.value = resultGroups;
        } else {
            console.warn('ได้รับข้อมูลที่ไม่ใช่ array สำหรับรายการที่ถูกลบ:', response.data);
        }
    } catch (err) {
        console.error('ไม่สามารถโหลดรายการที่ถูกลบ:', err);
        error.value = err.message || 'ไม่สามารถดึงข้อมูลรายการที่ถูกลบได้';
    } finally {
        loading.value = false;
    }
};

const restoreTodo = async (id) => {
    try {
        await axios.put(`${API_URL}/${id}/restore`);
        await loadTrashedItems();
    } catch (err) {
        console.error('ไม่สามารถกู้คืนงาน:', err);
        alert('ไม่สามารถกู้คืนงานได้ กรุณาลองใหม่อีกครั้ง');
    }
};

const permanentlyDeleteTodo = async (id) => {
    if (!confirm('คุณแน่ใจหรือไม่ว่าต้องการลบงานนี้อย่างถาวร? การกระทำนี้ไม่สามารถย้อนกลับได้')) {
        return;
    }
    try {
        await axios.delete(`${API_URL}/${id}/permanent`);
        await loadTrashedItems();
    } catch (err) {
        console.error('ไม่สามารถลบงานอย่างถาวร:', err);
        alert('ไม่สามารถลบงานอย่างถาวรได้ กรุณาลองใหม่อีกครั้ง');
    }
};

onMounted(loadTrashedItems);
</script>

<style scoped></style>
