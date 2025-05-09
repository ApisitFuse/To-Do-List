<template>
    <div class="flex flex-col items-center m-10">
        <h2 class="text-2xl font-semibold mb-4">Draggable List (Database Integrated)</h2>
        <div v-if="loading" class="text-lg">Loading items...</div>
        <div v-if="error" class="text-red-500 text-lg">Error loading items: {{ error }}</div>

        <draggable v-model="items" @change="handleDragEnd">
            <div class="list-group-item flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300 hover:bg-base-300"
                :class="{ 'opacity-50 bg-base-200': item.completed }" v-for="item in items" :key="item.ID">
                <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="item.completed"
                    @change="toggleTodo(item)" />

                <span class="flex-grow" :class="{ 'line-through': item.completed }">{{ item.title }}</span>

                <button @click="deleteTodo(item.ID)" class="btn btn-ghost btn-sm btn-square">❌</button>
            </div>
        </draggable>

        <router-link to="/" class="btn btn-primary mt-6">
            Back to To-Do List
        </router-link>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { VueDraggableNext as draggable } from 'vue-draggable-next'
import axios from 'axios'

const items = ref([])
const loading = ref(true)
const error = ref(null)


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