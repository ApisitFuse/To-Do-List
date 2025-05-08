import { createRouter, createWebHistory } from 'vue-router'
import TodoApp from '../components/TodoApp.vue'
import TodoAppTwo from '../components/TodoAppTwo.vue'
import DragDrop from '../components/DragDrop.vue'
import Trash from '../components/Trash.vue'

const routes = [
  { path: '/', name: 'TodoApp', component: TodoApp },
  { path: '/todo-two', name: 'TodoAppTwo', component: TodoAppTwo },
  { path: '/drag-drop', name: 'DragDrop', component: DragDrop },
  { path: '/trash', name: 'Trash', component: Trash }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
