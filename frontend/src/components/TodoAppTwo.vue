<template>
    <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">To-Do List</h1>

        <input type="text" v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task..."
            class="input input-bordered w-full mb-4" />

        <draggable v-model="items" @change="handleDragEnd">
            <div class="list-group-item flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300 hover:bg-base-300"
                :class="{ 'opacity-50 bg-base-200': item.completed }" v-for="item in items" :key="item.ID">
                <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="item.completed"
                    @change="toggleTodo(item)" />

                <span class="flex-grow" :class="{ 'line-through': item.completed }">{{ item.title }}</span>

                <button @click="deleteTodo(item.ID)" class="btn btn-ghost btn-sm btn-square" aria-label="Delete task">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor" class="w-5 h-5">
                        <path stroke-linecap="round" stroke-linejoin="round"
                            d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12.56 0c.34-.059.678-.113 1.017-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                    </svg>
                </button>
            </div>
        </draggable>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { VueDraggableNext as draggable } from 'vue-draggable-next'
import axios from 'axios'

const items = ref([])
const loading = ref(true)
const error = ref(null)
const newTodo = ref('')


const API_URL = 'http://localhost:8080/api/todos'


const loadItems = async () => {
    loading.value = true
    error.value = null
    try {
        const response = await axios.get(`${API_URL}/`)
        // items.value = response.data
        // ตรวจสอบว่า response.data เป็น array และไม่ null ก่อนทำการ sort
        if (Array.isArray(response.data)) {
            // เรียงลำดับ items ตาม display_order จากน้อยไปมาก
            items.value = response.data.sort((a, b) => a.displayOrder - b.displayOrder);
            console.log('Loaded items:', items.value)
        } else {
            items.value = [];
        }
    } catch (err) {
        console.error('Failed to load items:', err)
        error.value = err.message || 'Could not fetch items from the server.'
    } finally {
        loading.value = false
    }
}

const handleDragEnd = async (event) => {
    console.log('Drag ended, new order:', items.value)

    if (event.moved) {
        const movedItemId = event.moved.element.ID; // สมมติว่า item มี property ID
        const newIndex = event.moved.newIndex;
        const oldIndex = event.moved.oldIndex;

        console.log('movedItemId:', movedItemId)
        console.log('New index:', newIndex)
        console.log('Old index:', oldIndex)

        try {
            await axios.put(`${API_URL}/order`, {
                itemId: movedItemId,
                newIndex,
                oldIndex
            });

            console.log('Order updated successfully on the server.')

        } catch (err) {
            console.error('Failed to update order:', err)
            alert('Failed to save the new order. Please try again.')
            await loadItems()
        }
    }
}

const addTodo = async () => {
    if (!newTodo.value.trim()) return
    // 1. หาค่า displayOrder ที่มากที่สุดจาก items ที่มีอยู่
    let maxDisplayOrder = 0;
    if (items.value && items.value.length > 0) {
        // ใช้ Math.max ร่วมกับ spread operator และ map เพื่อหาค่าสูงสุด
        // หรือจะ loop หากต้องการความชัดเจนมากขึ้น
        maxDisplayOrder = Math.max(...items.value.map(item => item.displayOrder || 0));
    }

    // 2. คำนวณ displayOrder ใหม่
    const newDisplayOrder = maxDisplayOrder + 1;

    try {
        await axios.post(`${API_URL}/`, {
            title: newTodo.value,
            completed: false,
            displayOrder: newDisplayOrder
        });
        newTodo.value = '';
        await loadItems();
    } catch (err) {
        console.error('Failed to add todo:', err);
        alert('Failed to add the new task. Please try again.');
    }
}

const toggleTodo = async (todo) => {
    await axios.put(`${API_URL}/${todo.ID}`, todo)

}

const deleteTodo = async (id) => {
    await axios.delete(`${API_URL}/${id}`)
    await loadItems()
}

onMounted(loadItems)
</script>

<style scoped>
.list-group-item:hover {
    cursor: pointer;
}

.sortable-chosen {
    background-color: #57c6fd !important;
    border-color: #99ddff !important;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    cursor: grabbing !important;

}

.sortable-chosen .flex-grow {
    color: #1f2937 !important;
    opacity: 1 !important;
}
</style>