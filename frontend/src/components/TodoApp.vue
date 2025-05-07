<template>
    <!-- <div class="todo-app">
        <h1>To-Do List</h1>
        <input v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a task..." />
        <ul>
            <li v-for="todo in todos" :key="todo.ID" :class="{ 'completed': todo.completed }">
                
                <input type="checkbox" checked="checked" class="checkbox checkbox-success" v-model="todo.completed" @change="toggleTodo(todo)"/>
                <span :style="{ textDecoration: todo.completed ? 'line-through' : 'none' }">
                    {{ todo.title }}
                </span>
                <button @click="deleteTodo(todo.ID)">❌</button>
            </li>
        </ul>
    </div> -->

    <!-- ใช้ container ของ Tailwind และ DaisyUI card -->
    <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">To-Do List</h1>

        <input type="text" v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task..."
            class="input input-bordered w-full mb-4" />

        <ul class="space-y-2">

            <li v-for="todo in todos" :key="todo.ID"
                class="flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300"
                :class="{ 'opacity-50 bg-base-200': todo.completed }">

                <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="todo.completed"
                    @change="toggleTodo(todo)" />

                <span class="flex-grow" :class="{ 'line-through': todo.completed }">{{ todo.title }}</span>

                <button @click="deleteTodo(todo.ID)" class="btn btn-ghost btn-sm btn-square">❌</button>
            </li>
        </ul>
        <router-link to="/todo-two" class="btn btn-active btn-primary">
            <button type="button">Go to second page</button>
        </router-link>
        <router-link to="/drag-drop" class="btn btn-active btn-primary">
            <button type="button">Go to drag-drop page</button>
        </router-link>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'


const todos = ref([])
const newTodo = ref('')

const API_URL = 'http://localhost:8080/api/todos/'


const loadTodos = async () => {
    const res = await axios.get(API_URL)
    todos.value = res.data
}

const addTodo = async () => {
    if (!newTodo.value.trim()) return
    await axios.post(API_URL, {
        title: newTodo.value,
        completed: false
    })
    newTodo.value = ''
    await loadTodos()
}

const toggleTodo = async (todo) => {
    await axios.put(`${API_URL}${todo.ID}`, todo)

}

const deleteTodo = async (id) => {
    await axios.delete(`${API_URL}${id}`)
    await loadTodos()
}

onMounted(loadTodos)
</script>

<style scoped>
.todo-app {
    max-width: 500px;
    margin: auto;
    font-family: sans-serif;
}

input[type="text"] {
    width: 100%;
    padding: 8px;
}

ul {
    list-style: none;
    padding: 0;
}

li {
    display: flex;
    align-items: center;
    margin-top: 8px;
}

button {
    margin-left: auto;
    background: none;
    border: none;
    cursor: pointer;
}

/* --- สไตล์สำหรับกรอบรายการ Task --- */
.todo-item {
    border: 1px solid #ccc;
    /* สีขอบเริ่มต้น */
    padding: 10px 15px;
    /* ระยะห่างภายในกรอบ */
    border-radius: 5px;
    /* ทำให้ขอบมน */
    margin-bottom: 10px;
    /* ระยะห่างระหว่างรายการ */
    transition: border-color 0.3s ease, background-color 0.3s ease;
    /* เพิ่ม animation ตอนเปลี่ยนสี */
    background-color: white;
    /* สีพื้นหลังเริ่มต้น */
}

/* สไตล์เมื่อรายการถูกติ๊ก (completed) */
.todo-item.completed {
    border-color: #e0e0e0;
    /* สีขอบจางลง */
    background-color: #f9f9f9;
    /* สีพื้นหลังจางลงเล็กน้อย (ถ้าต้องการ) */
}
</style>