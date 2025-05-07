<template>
    <div class="flex flex-col items-center m-10">
        <h2 class="text-2xl font-semibold mb-4">Draggable List (Database Integrated)</h2>
        <div v-if="loading" class="text-lg">Loading items...</div>
        <div v-if="error" class="text-red-500 text-lg">Error loading items: {{ error }}</div>

        <!-- <draggable class="dragArea list-group w-full" :list="items" @change="handleDragEnd"> -->
        <draggable class="dragArea list-group w-full" v-model="items" @change="handleDragEnd">
            <div class="list-group-item bg-[#c15925] m-1 p-3 rounded-md text-center" v-for="item in items"
                :key="item.ID">
                {{ item.title }}
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
            items.value = []; // หรือจัดการ error ตามความเหมาะสม
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
.list-group {
    min-height: 100px;
}

.list-group-item:hover {
    background-color: #d2855f;
}
</style>






<!-- <template>
    <div class="flex m-10">
        <draggable v-model="myList">
            <div class="list-group-item bg-amber-700 m-1 p-3 rounded-md text-center" v-for="element in list"
                :key="element.name">
                {{ element.name }}
            </div>
        </draggable>
    </div>
</template>

<script>
import { defineComponent } from 'vue'
import { VueDraggableNext } from 'vue-draggable-next'
export default defineComponent({
    components: {
        draggable: VueDraggableNext,
    },
    data() {
        return {
            enabled: true,
            list: [
                { name: 'John', id: 1 },
                { name: 'Joao', id: 2 },
                { name: 'Jean', id: 3 },
                { name: 'Gerard', id: 4 },
            ],
            dragging: false,
        }
    },
    methods: {
        log(event) {
            console.log(event)
        },
    },
})
</script> -->