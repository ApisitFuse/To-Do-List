import { createRouter, createWebHistory } from 'vue-router'
import TodoApp from '../components/TodoApp.vue'
import TodoAppTwo from '../components/TodoAppTwo.vue'
import DragDrop from '../components/DragDrop.vue'

const routes = [
  { path: '/', name: 'TodoApp', component: TodoApp },
  { path: '/todo-two', name: 'TodoAppTwo', component: TodoAppTwo },
  { path: '/drag-drop', name: 'DragDrop', component: DragDrop }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
