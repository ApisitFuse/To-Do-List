<template>

    <!-- ใช้ container ของ Tailwind และ DaisyUI card -->
    <!-- <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">To-Do List</h1>

        <input type="text" v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task..."
            class="input input-bordered w-full mb-4" />

        <draggable v-model="todos" item-key="ID" tag="ul" class="space-y-2" ghost-class="opacity-50" animation="150">
            <template #item="{ element: todo }">

                <li class="flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300 cursor-move"
                    :class="{ 'opacity-50 bg-base-200': todo.completed }">

                    <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="todo.completed"
                        @change="toggleTodo(todo)" />

                    <span class="flex-grow" :class="{ 'line-through': todo.completed }">{{ todo.title }}</span>

                    <button @click="deleteTodo(todo.ID)" class="btn btn-ghost btn-sm btn-square">❌</button>
                </li>
            </template>
</draggable>
</div> -->

    <!-- <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">To-Do List</h1>

        <input type="text" v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task..."
            class="input input-bordered w-full mb-4" />

        <draggable v-model="todos" item-key="ID" tag="ul" class="space-y-2" ghost-class="opacity-50" animation="150">
            <template #item="{ element }">
                <li class="flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300 cursor-move"
                    :class="{ 'opacity-50 bg-base-200': element.completed }">
                    <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="element.completed"
                        @change="toggleTodo(element)" />
                    <span class="flex-grow" :class="{ 'line-through': element.completed }">
                        {{ element.title }}
                    </span>
                    <button @click="deleteTodo(element.ID)" class="btn btn-ghost btn-sm btn-square">
                        ❌
                    </button>
                </li>
            </template>
        </draggable>
    </div> -->

    <div class="max-w-lg mx-auto mt-10 p-6 bg-base-100 rounded-box shadow-xl">
        <h1 class="text-3xl font-bold text-center mb-6">To-Do List</h1>

        <input type="text" v-model="newTodo" @keyup.enter="addTodo" placeholder="Add a new task..."
            class="input input-bordered w-full mb-4" />

        <draggable v-model="todos" item-key="ID" tag="ul" class="space-y-2" ghost-class="opacity-50" animation="150">
            <template #item="{ element }">
                <li class="flex items-center p-3 border border-base-300 rounded-box transition-opacity duration-300 cursor-move"
                    :class="{ 'opacity-50 bg-base-200': element.completed }">

                    <input type="checkbox" class="checkbox checkbox-success mr-3" v-model="element.completed"
                        @change="toggleTodo(element)" />

                    <span class="flex-grow" :class="{ 'line-through': element.completed }">{{ element.title }}</span>

                    <button @click="deleteTodo(element.ID)" class="btn btn-ghost btn-sm btn-square">❌</button>
                </li>
            </template>
        </draggable>
    </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { VueDraggableNext } from 'vue-draggable-next'
// import draggable from 'vuedraggable'

const draggable = VueDraggableNext

const todos = ref([])
const newTodo = ref('')

const loadTodos = async () => {
    const res = await axios.get('http://localhost:8080/api/todos/')
    todos.value = res.data
}

const addTodo = async () => {
    if (!newTodo.value.trim()) return
    await axios.post('http://localhost:8080/api/todos/', {
        title: newTodo.value,
        completed: false
    })
    newTodo.value = ''
    await loadTodos()
}

const toggleTodo = async (todo) => {
    await axios.put(`http://localhost:8080/api/todos/${todo.ID}`, todo)

}

const deleteTodo = async (id) => {
    await axios.delete(`http://localhost:8080/api/todos/${id}`)
    await loadTodos()
}

onMounted(loadTodos)
</script>

<style scoped></style>