import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'

// createApp(App).mount('#app')
const app = createApp(App)

app.use(router) // 2. Tell the app to use the router
app.mount('#app')
