import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import './styles/variables.css'
import { initFontSize } from './composables/useFontSize'
import { useUserStore } from './store/userStore'
import { useTokenStore } from './store/tokenStore'

initFontSize()

const app = createApp(App)

// Pinia状态管理
const pinia = createPinia()
app.use(pinia)

// 开发模式：Mock 用户数据（方便测试页面效果）
if (import.meta.env.DEV && import.meta.env.VITE_MOCK_USER === 'true') {
  const userStore = useUserStore()
  const tokenStore = useTokenStore()

  // 设置 mock token
  tokenStore.setToken('mock_token_for_dev')
  tokenStore.setRefreshToken('mock_refresh_token_for_dev')

  // 设置 mock 用户
  userStore.setUser({
    id: 1,
    phone: '13800138000',
    role: 'customer',
    has_password: true
  })

  // 设置 mock 个人资料
  userStore.setProfile({
    id: 1,
    name: '张三',
    id_number: '110101199001011234',
    address: '北京市朝阳区望京街道',
    user_type: 'elderly'
  })

  console.log('🔧 开发模式：已加载 Mock 用户数据')
}

// Vue Router
app.use(router)

// Element Plus
app.use(ElementPlus)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')
