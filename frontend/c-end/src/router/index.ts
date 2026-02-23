import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useTokenStore } from '@/store/tokenStore'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { requiresAuth: false, title: '首页' }
  },
  {
    path: '/quick',
    name: 'QuickStart',
    component: () => import('@/views/QuickStart.vue'),
    meta: { requiresAuth: false, title: '申请服务' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false, title: '登录' }
  },
  {
    path: '/services',
    name: 'ServiceList',
    component: () => import('@/views/ServiceListView.vue'),
    meta: { requiresAuth: false, title: '全部服务' }
  },
  {
    path: '/requests',
    name: 'RequestList',
    component: () => import('@/views/RequestList.vue'),
    meta: { requiresAuth: true, title: '我的服务' }
  },
  {
    path: '/requests/:id',
    name: 'RequestDetail',
    component: () => import('@/views/RequestDetail.vue'),
    meta: { requiresAuth: true, title: '服务详情' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { requiresAuth: true, title: '我的资料' }
  },
  {
    path: '/profile/basic',
    name: 'ProfileBasic',
    component: () => import('@/views/profile/BasicInfo.vue'),
    meta: { requiresAuth: true, title: '基本信息' }
  },
  {
    path: '/profile/contact',
    name: 'ProfileContact',
    component: () => import('@/views/profile/ContactInfo.vue'),
    meta: { requiresAuth: true, title: '联系信息' }
  },
  {
    path: '/profile/address',
    name: 'ProfileAddress',
    component: () => import('@/views/profile/AddressInfo.vue'),
    meta: { requiresAuth: true, title: '服务地址' }
  },
  {
    path: '/profile/health',
    name: 'ProfileHealth',
    component: () => import('@/views/profile/HealthInfo.vue'),
    meta: { requiresAuth: true, title: '健康档案' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { requiresAuth: true, title: '设置' }
  },
  {
    path: '/news',
    name: 'NewsList',
    component: () => import('@/views/NewsList.vue'),
    meta: { requiresAuth: false, title: '站点动态' }
  },
  {
    path: '/news/:id',
    name: 'NewsDetail',
    component: () => import('@/views/NewsDetail.vue'),
    meta: { requiresAuth: false, title: '详情' }
  },
  {
    path: '/notifications',
    name: 'NotificationList',
    component: () => import('@/views/NotificationList.vue'),
    meta: { requiresAuth: true, title: '消息通知' }
  },
  {
    path: '/mine',
    name: 'Mine',
    component: () => import('@/views/Mine.vue'),
    meta: { requiresAuth: false, title: '我的' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：检查登录状态
router.beforeEach((to, _from, next) => {
  const tokenStore = useTokenStore()

  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - sCare` : 'sCare 社区养老服务'

  // 检查是否需要登录
  if (to.meta.requiresAuth && !tokenStore.isLoggedIn) {
    // 保存原始目标路径
    next({
      name: 'Login',
      query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
})

export default router
