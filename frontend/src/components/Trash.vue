<template>
    <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">รายการงานที่ถูกลบ</h1>

        <div v-if="loading" class="text-center">กำลังโหลด...</div>
        <div v-if="error" class="text-center text-error">{{ error }}</div>

        <div v-if="!loading && trashedItems.length === 0 && !error" class="text-center text-gray-500">
            ไม่มีงานในถังขยะ
        </div>

        <div v-for="item in trashedItems" :key="item.ID"
            class="flex items-center p-3 border border-base-300 rounded-box mb-2 hover:bg-base-200">
            <span class="flex-grow" :class="{ 'line-through': item.completed }">{{ item.title }}</span>
            <button @click="restoreTodo(item.ID)" class="btn btn-sm btn-success mr-2">
                กู้คืน
            </button>
            <!-- ตัวเลือก: สามารถเพิ่มปุ่มลบถาวรได้ที่นี่ หากต้องการ -->
            
            <button @click="permanentlyDeleteTodo(item.ID)" class="btn btn-sm btn-error">
                ลบถาวร
            </button> 
           
        </div>
        <div class="mt-6 text-center">
            <router-link to="/todo-two" class="btn btn-primary">กลับ</router-link>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const trashedItems = ref([]);
const loading = ref(true);
const error = ref(null);

const API_URL = 'http://localhost:8080/api/todos'; // ใช้ API base URL เดียวกัน

const loadTrashedItems = async () => {
    loading.value = true;
    error.value = null;
    try {
        // Endpoint นี้จำเป็นต้องสร้างขึ้นในฝั่ง backend
        // ควรจะคืนค่ารายการที่ถูกตั้งสถานะเป็น 'deleted' หรือ 'trashed'
        const response = await axios.get(`${API_URL}/trashed`); // Endpoint ใหม่
        if (Array.isArray(response.data)) {
            trashedItems.value = response.data;
        } else {
            trashedItems.value = [];
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
        // Endpoint นี้จำเป็นต้องสร้างขึ้นในฝั่ง backend
        // ควรจะตั้งสถานะรายการเป็น 'not deleted' หรือ 'active'
        await axios.put(`${API_URL}/${id}/restore`); // Endpoint ใหม่
        await loadTrashedItems(); // รีเฟรชรายการที่ถูกลบ
    } catch (err) {
        console.error('ไม่สามารถกู้คืนงาน:', err);
        alert('ไม่สามารถกู้คืนงานได้ กรุณาลองใหม่อีกครั้ง');
    }
};

const permanentlyDeleteTodo = async (id) => {
    // แสดง popup ยืนยันก่อนทำการลบถาวร
    if (!confirm('คุณแน่ใจหรือไม่ว่าต้องการลบงานนี้อย่างถาวร? การกระทำนี้ไม่สามารถย้อนกลับได้')) {
        return;
    }
    try {
        // Endpoint นี้จำเป็นต้องสร้างขึ้นในฝั่ง backend
        // และควรทำการลบข้อมูลออกจากฐานข้อมูลอย่างถาวร (hard delete)
        await axios.delete(`${API_URL}/${id}/permanent`); // Endpoint ใหม่สำหรับการลบถาวร
        await loadTrashedItems(); // รีเฟรชรายการที่ถูกลบ
    } catch (err) {
        console.error('ไม่สามารถลบงานอย่างถาวร:', err);
        alert('ไม่สามารถลบงานอย่างถาวรได้ กรุณาลองใหม่อีกครั้ง');
    }
};

onMounted(loadTrashedItems);
</script>

<style scoped>

</style>
