import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

/**
 * 路由配置
 * 注意：路由路径需要与后端菜单配置的 path 一致
 */
const routes: RouteRecordRaw[] = [
  // 登录页
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/Login/index.vue'),
    meta: {
      title: '登录',
      public: true, // 公开访问，无需认证
    },
  },
  // 主应用布局
  {
    path: '/',
    component: () => import('@/layouts/BasicLayout.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true, // 需要认证
    },
    children: [
      // 工作台
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/pages/Dashboard/index.vue'),
        meta: {
          title: '工作台',
          icon: 'Odometer',
          permission_code: 'dashboard',
        },
      },
      // ========== 服务管理 ==========
      // 需求管理
      {
        path: '/services/requests',
        name: 'RequestManagement',
        component: () => import('@/pages/RequestManagement/index.vue'),
        meta: {
          title: '需求管理',
          icon: 'Document',
          permission_code: 'service:request:list',
        },
      },
      // 任务管理
      {
        path: '/services/tasks',
        name: 'TaskManagement',
        component: () => import('@/pages/TaskManagement/index.vue'),
        meta: {
          title: '任务管理',
          icon: 'List',
          permission_code: 'service:task',
        },
      },
      // 任务详情（不在菜单显示）
      {
        path: '/services/tasks/:id(\\d+)',
        name: 'TaskDetail',
        component: () => import('@/pages/TaskDetail/index.vue'),
        meta: {
          title: '任务详情',
          hidden: true,
          // 详情页权限较为复杂（可能是任务池、列表或我的任务来源），暂时移除路由级权限控制，依靠接口权限
          // permission_code: 'service:task',
          activeMenu: '/services/tasks',
        },
      },
      // ========== 数据中心 ==========
      // 统计分析
      {
        path: '/data/statistics',
        name: 'StatisticsOverview',
        component: () => import('@/pages/StatisticsOverview/index.vue'),
        meta: {
          title: '统计分析',
          icon: 'DataLine',
          permission_code: 'data:statistics:view',
        },
      },
      // 报表管理
      {
        path: '/data/reports',
        name: 'StatisticsReports',
        component: () => import('@/pages/StatisticsReports/index.vue'),
        meta: {
          title: '报表管理',
          icon: 'Download',
          permission_code: 'data:report:list',
        },
      },
      // ========== 居民管理 ==========
      // 老年人档案
      {
        path: '/residents/elderly',
        name: 'ElderlyManagement',
        component: () => import('@/pages/ElderlyManagement/index.vue'),
        meta: {
          title: '老年人档案',
          icon: 'Avatar',
          permission_code: 'resident:elderly:list',
        },
      },
      // 老年人档案详情（不在菜单显示）
      {
        path: '/residents/elderly/:id',
        name: 'ElderlyDetail',
        component: () => import('@/pages/ElderlyDetail/index.vue'),
        meta: {
          title: '档案详情',
          hidden: true,
          permission_code: 'resident:elderly:detail',
          activeMenu: '/residents/elderly',
        },
      },
      // ========== 站点管理 ==========
      // 站点列表
      {
        path: '/stations/list',
        name: 'StationManagement',
        component: () => import('@/pages/StationManagement/index.vue'),
        meta: {
          title: '站点列表',
          icon: 'OfficeBuilding',
          permission_code: 'station:list:view',
        },
      },
      // 围栏管理
      {
        path: '/stations/zones',
        name: 'ZoneManagement',
        component: () => import('@/pages/ZoneManagement/index.vue'),
        meta: {
          title: '服务区域',
          icon: 'MapLocation',
          permission_code: 'station:zone:list',
        },
      },
      // ========== 系统管理 ==========
      // 用户管理
      {
        path: '/system/users',
        name: 'UserManagement',
        component: () => import('@/pages/UserManagement/index.vue'),
        meta: {
          title: '用户管理',
          icon: 'UserFilled',
          permission_code: 'system:user:list',
        },
      },
      // 角色权限管理
      {
        path: '/system/roles',
        name: 'RolePermission',
        component: () => import('@/pages/RolePermission/index.vue'),
        meta: {
          title: '角色管理',
          icon: 'Lock',
          permission_code: 'system:role:list',
        },
      },
      // 菜单管理
      {
        path: '/system/menus',
        name: 'MenuManagement',
        component: () => import('@/pages/MenuManagement/index.vue'),
        meta: {
          title: '菜单管理',
          icon: 'Menu',
          permission_code: 'system:menu:list',
        },
      },
      // ========== 内容管理 ==========
      // 轮播图管理
      {
        path: '/content/banners',
        name: 'BannerManagement',
        component: () => import('@/pages/BannerManagement/index.vue'),
        meta: {
          title: '轮播图管理',
          icon: 'Picture',
          permission_code: 'content:banner:list',
        },
      },
      // 新闻管理
      {
        path: '/content/news',
        name: 'NewsManagement',
        component: () => import('@/pages/NewsManagement/index.vue'),
        meta: {
          title: '新闻管理',
          icon: 'Document',
          permission_code: 'content:news:list',
        },
      },
      // 通知管理
      {
        path: '/content/notifications',
        name: 'NotificationManagement',
        component: () => import('@/pages/NotificationManagement/index.vue'),
        meta: {
          title: '通知管理',
          icon: 'Bell',
          permission_code: 'content:notification:list',
        },
      },
      // 通知中心（菜单入口，指向同一页面）
      {
        path: '/notifications',
        name: 'NotificationCenter',
        component: () => import('@/pages/NotificationManagement/index.vue'),
        meta: {
          title: '通知中心',
          icon: 'Bell',
          permission_code: 'public:notification',
        },
      },
      // ========== 个人中心 ==========
      // 个人信息
      {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/pages/Profile/index.vue'),
        meta: {
          title: '个人信息',
          icon: 'User',
          permission_code: 'profile',
        },
      },
    ],
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
