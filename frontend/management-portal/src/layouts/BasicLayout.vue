<template>
  <el-container class="layout-container">
    <!-- 左侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="layout-aside">
      <!-- Logo -->
      <div class="logo" @click="router.push('/dashboard')" style="cursor: pointer;">
        <span v-if="!isCollapse">{{ appTitle }}</span>
        <span v-else>SC</span>
      </div>

      <!-- 侧边栏菜单 -->
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :unique-opened="true"
        router
      >
        <template v-for="menu in menuRoutes" :key="menu.id">
          <!-- 有子菜单 -->
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
            <template #title>
              <el-icon><component :is="menu.icon || 'Grid'" /></el-icon>
              <span>{{ menu.name }}</span>
            </template>
            <el-menu-item
              v-for="child in menu.children"
              :key="child.id"
              :index="child.path"
            >
              <el-icon v-if="child.icon"><component :is="child.icon" /></el-icon>
              <template #title>{{ child.name }}</template>
            </el-menu-item>
          </el-sub-menu>
          <!-- 无子菜单 -->
          <el-menu-item v-else :index="menu.path">
            <el-icon><component :is="menu.icon || 'Grid'" /></el-icon>
            <template #title>{{ menu.name }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>

    <!-- 右侧主体 -->
    <el-container>
      <!-- 顶部导航栏 -->
      <el-header class="layout-header">
        <!-- 折叠按钮 -->
        <div class="header-left">
          <el-icon class="collapse-icon" @click="toggleCollapse">
            <Expand v-if="isCollapse" />
            <Fold v-else />
          </el-icon>
        </div>

        <!-- 用户信息 -->
        <div class="header-right">
          <div class="header-action-item" @click="toggleNotificationDrawer">
            <el-badge :is-dot="unreadCount > 0" class="notification-badge">
              <el-icon class="notification-icon"><Bell /></el-icon>
            </el-badge>
          </div>
          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32">
                <el-icon><User /></el-icon>
              </el-avatar>
              <span class="user-name">{{ authStore.user?.name }}</span>
              <el-icon class="arrow-down"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>
                  <div class="user-detail">
                    <div>角色：{{ roleText }}</div>
                    <div>手机：{{ authStore.user?.phone }}</div>
                  </div>
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main class="layout-main">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>

    <!-- 通知抽屉 -->
    <el-drawer
      v-model="notificationDrawerVisible"
      title="通知中心"
      direction="rtl"
      size="350px"
    >
      <div class="notification-list" v-loading="drawerLoading">
        <div v-if="notificationList.length === 0" class="empty-text">
          暂无通知
        </div>
        <div
          v-for="item in notificationList"
          :key="item.id"
          class="notification-item"
          :class="{ unread: !item.is_read }"
          @click="handleViewNotification(item)"
        >
          <div class="notification-header">
            <span class="notification-title">{{ item.title }}</span>
            <span v-if="!item.is_read" class="unread-dot"></span>
          </div>
          <div class="notification-content">{{ item.content }}</div>
          <div class="notification-footer">
            <span class="notification-time">{{ formatDateTime(item.created_at) }}</span>
            <el-button
              v-if="!item.is_read"
              link
              type="primary"
              size="small"
              @click.stop="handleMarkRead(item)"
            >
              已读
            </el-button>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="drawer-footer">
          <el-button link type="primary" @click="viewAllNotifications">查看全部通知</el-button>
        </div>
      </template>
    </el-drawer>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Expand, Fold, User, ArrowDown, SwitchButton, Bell } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'
import { menuApi as menuAPI, notificationApi } from '@/api'
import type { Menu, Notification } from '@/types/api'
import dayjs from 'dayjs'

// 应用标题
const appTitle = import.meta.env.VITE_APP_TITLE || 'sCare'

// 路由
const router = useRouter()
const route = useRoute()

// 状态管理
const authStore = useAuthStore()

// 侧边栏折叠状态
const isCollapse = ref(false)

// 用户菜单数据
const userMenus = ref<Menu[]>([])
const menusLoading = ref(false)
const unreadCount = ref(0)

// 通知抽屉相关
const notificationDrawerVisible = ref(false)
const notificationList = ref<Notification[]>([])
const drawerLoading = ref(false)

// 当前激活的菜单
const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta.activeMenu) {
    return meta.activeMenu as string
  }
  return path
})

// 菜单路由列表（从API获取的动态菜单）
const menuRoutes = computed(() => {
  return userMenus.value.filter((menu) => !menu.hidden)
})

// 角色文本
const roleText = computed(() => {
  const roleMap: Record<string, string> = {
    staff: '工作人员',
    station_manager: '站点管理员',
    admin: '系统管理员',
  }
  const roles = authStore.user?.roles || []
  return roles.map((role) => roleMap[role] || role).join(', ') || '未知'
})

/**
 * 切换侧边栏折叠状态
 */
function toggleCollapse() {
  isCollapse.value = !isCollapse.value
}

/**
 * 获取用户菜单
 */
async function fetchUserMenus() {
  menusLoading.value = true
  try {
    const response = await menuAPI.getUserMenus()
    userMenus.value = response.data || []
  } catch (error) {
    console.error('获取用户菜单失败:', error)
    ElMessage.error('获取菜单失败')
  } finally {
    menusLoading.value = false
  }
}

/**
 * 检查未读通知
 */
async function checkUnread() {
  try {
    const res = await notificationApi.getNotifications({ page: 1, page_size: 20 })
    if (res.data && res.data.items) {
      notificationList.value = res.data.items
      unreadCount.value = res.data.items.filter(n => !n.is_read).length
    }
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  }
}

/**
 * 切换通知抽屉
 */
function toggleNotificationDrawer() {
  notificationDrawerVisible.value = !notificationDrawerVisible.value
  if (notificationDrawerVisible.value) {
    fetchNotifications()
  }
}

/**
 * 加载通知列表（抽屉用）
 */
async function fetchNotifications() {
  drawerLoading.value = true
  try {
    await checkUnread()
  } finally {
    drawerLoading.value = false
  }
}

/**
 * 标记通知为已读
 */
async function handleMarkRead(notification: Notification) {
  try {
    await notificationApi.markAsRead(notification.id)
    // 更新本地状态
    const item = notificationList.value.find(n => n.id === notification.id)
    if (item) {
      item.is_read = true
    }
    unreadCount.value = notificationList.value.filter(n => !n.is_read).length
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

/**
 * 查看通知详情
 */
function handleViewNotification(notification: Notification) {
  handleMarkRead(notification)
  // 可以在这里添加跳转逻辑，如果需要的话
  // 目前已经在列表中展开显示了内容，或者跳转到完整列表
  // notificationDrawerVisible.value = false
  // router.push('/content/notifications')
}

/**
 * 查看全部通知
 */
function viewAllNotifications() {
  notificationDrawerVisible.value = false
  router.push('/notifications')
}

/**
 * 格式化时间
 */
function formatDateTime(time: string) {
  return dayjs(time).format('MM-DD HH:mm')
}

/**
 * 处理下拉菜单命令
 */
async function handleCommand(command: string) {
  if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })

      authStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    } catch {
      // 用户取消
    }
  }
}

// 初始化时获取用户菜单
onMounted(() => {
  fetchUserMenus()
  checkUnread()
})
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;
  overflow: hidden;

  .layout-aside {
    background-color: #001529;
    transition: width 0.3s;
    box-shadow: 2px 0 6px rgba(0, 21, 41, 0.35);

    .logo {
      height: 60px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      font-size: 18px;
      font-weight: bold;
      background-color: #002140;
      white-space: nowrap;
      overflow: hidden;
    }

    :deep(.el-menu) {
      border-right: none;
      background-color: transparent;

      .el-menu-item {
        color: rgba(255, 255, 255, 0.85);

        &:hover {
          color: #fff;
          background-color: rgba(255, 255, 255, 0.12);
        }

        &.is-active {
          color: #fff;
          background-color: #1890ff;
        }
      }

      .el-sub-menu {
        .el-sub-menu__title {
          color: rgba(255, 255, 255, 0.85);

          &:hover {
            color: #fff;
            background-color: rgba(255, 255, 255, 0.12);
          }
        }

        .el-menu-item {
          color: rgba(255, 255, 255, 0.85);

          &:hover {
            color: #fff;
            background-color: rgba(255, 255, 255, 0.12);
          }

          &.is-active {
            color: #fff;
            background-color: #1890ff;
          }
        }
      }
    }
  }

  .layout-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: #fff;
    border-bottom: 1px solid #f0f0f0;
    padding: 0 20px;

    .header-left {
      .collapse-icon {
        font-size: 20px;
        cursor: pointer;
        transition: color 0.3s;

        &:hover {
          color: #1890ff;
        }
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 12px;

      .header-action-item {
        display: flex;
        align-items: center;
        height: 100%;
        padding: 0 8px;
        cursor: pointer;
        transition: background-color 0.3s;

        &:hover {
          background-color: #f5f5f5;
        }

        .notification-icon {
          font-size: 20px;
          color: #606266;
        }
        
        :deep(.el-badge__content.is-fixed) {
          top: 0;
          right: 0;
          transform: translateY(-50%) translateX(50%) scale(0.8);
        }
      }

      .user-info {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
        padding: 5px 10px;
        border-radius: 4px;
        transition: background-color 0.3s;

        &:hover {
          background-color: #f5f5f5;
        }

        .user-name {
          font-size: 14px;
          color: #303133;
        }

        .arrow-down {
          font-size: 12px;
          color: #909399;
        }
      }
    }
  }

  .layout-main {
    background-color: #f0f2f5;
    padding: 20px;
    overflow-y: auto;
  }
}

.user-detail {
  font-size: 13px;
  color: #606266;
  line-height: 1.8;
  padding: 5px 0;
}

// 过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.notification-list {
  padding: 0 4px;
  
  .empty-text {
    text-align: center;
    color: #909399;
    padding: 20px 0;
  }

  .notification-item {
    padding: 12px;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
    transition: background-color 0.2s;
    border-radius: 4px;

    &:hover {
      background-color: #f5f7fa;
    }

    &.unread {
      background-color: #fdf6ec;
      
      .notification-title {
        font-weight: 600;
        color: #303133;
      }
    }

    .notification-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 6px;

      .notification-title {
        font-size: 14px;
        color: #606266;
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        margin-right: 8px;
      }

      .unread-dot {
        width: 6px;
        height: 6px;
        background-color: #f56c6c;
        border-radius: 50%;
      }
    }

    .notification-content {
      font-size: 13px;
      color: #909399;
      line-height: 1.5;
      margin-bottom: 8px;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    .notification-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .notification-time {
        font-size: 12px;
        color: #c0c4cc;
      }
    }
  }
}

.drawer-footer {
  text-align: center;
  padding-top: 10px;
}
</style>
