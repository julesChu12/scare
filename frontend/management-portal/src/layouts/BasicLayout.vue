<template>
  <el-container class="layout-container">
    <!-- 左侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="layout-aside">
      <!-- Logo -->
      <div class="logo">
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
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Expand, Fold, User, ArrowDown, SwitchButton } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'
import { menuApi as menuAPI } from '@/api'
import type { Menu } from '@/types/api'

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
</style>
