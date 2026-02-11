import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { setupPermissionGuard } from './router/guards/permission.guard'
import { useAuthStore } from './store/modules/auth'
import { vPermission } from './directives/permission'

const app = createApp(App)

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 注册权限指令
app.directive('permission', vPermission)

// 状态管理
const pinia = createPinia()
app.use(pinia)

// 初始化认证状态（从 localStorage 恢复）
const authStore = useAuthStore()
authStore.init()

// 路由
app.use(router)

// 设置权限守卫
setupPermissionGuard(router)

// Element Plus
app.use(ElementPlus)

app.mount('#app')

