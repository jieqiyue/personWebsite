import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'
import { initTheme } from './utils/theme'

// 初始化主题
initTheme()

const app = createApp(App)
app.use(router)
app.mount('#app')
