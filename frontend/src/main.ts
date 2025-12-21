import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

// 预加载 composables 确保事件监听在应用启动时就注册
import './composables/useMediaList'
import './composables/useSelectedDir'
import './composables/useSimilarImages'

const app = createApp(App)

app.use(router)
app.mount('#app')
