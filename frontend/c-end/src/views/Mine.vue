<template>
  <div class="mine-container">
    <!-- 顶部用户信息卡片 -->
    <div class="user-card" @click="handleUserCardClick">
      <div class="user-avatar">
        <el-icon v-if="!userStore.profile?.name" :size="40"><UserFilled /></el-icon>
        <span v-else class="avatar-text">{{ avatarText }}</span>
      </div>
      <div class="user-info">
        <div class="user-name">{{ displayName }}</div>
        <div class="user-phone">{{ displayPhone }}</div>
      </div>
      <el-icon class="arrow-icon"><ArrowRight /></el-icon>
    </div>

    <!-- 功能入口列表 -->
    <div class="menu-section">
      <div class="menu-item" @click="goToNotifications">
        <div class="menu-left">
          <span class="menu-icon">🔔</span>
          <span class="menu-text">消息通知</span>
        </div>
        <el-icon class="menu-arrow"><ArrowRight /></el-icon>
      </div>

      <div class="menu-item" @click="goToAddress">
        <div class="menu-left">
          <span class="menu-icon">📍</span>
          <span class="menu-text">我的地址</span>
        </div>
        <el-icon class="menu-arrow"><ArrowRight /></el-icon>
      </div>

      <div class="menu-item" @click="goToSettings">
        <div class="menu-left">
          <span class="menu-icon">⚙️</span>
          <span class="menu-text">设置</span>
        </div>
        <el-icon class="menu-arrow"><ArrowRight /></el-icon>
      </div>

      <div class="menu-item" @click="showServiceDialog">
        <div class="menu-left">
          <span class="menu-icon">📞</span>
          <span class="menu-text">联系客服</span>
        </div>
        <el-icon class="menu-arrow"><ArrowRight /></el-icon>
      </div>

      <div class="menu-item" @click="showAboutDialog">
        <div class="menu-left">
          <span class="menu-icon">ℹ️</span>
          <span class="menu-text">关于我们</span>
        </div>
        <el-icon class="menu-arrow"><ArrowRight /></el-icon>
      </div>
    </div>

    <!-- 底部占位 -->
    <div class="bottom-placeholder"></div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <div class="nav-item" @click="goToHome">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </div>
      <div class="nav-item" @click="goToServices">
        <el-icon><List /></el-icon>
        <span>服务</span>
      </div>
      <div class="nav-item active">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>

    <!-- 客服电话弹窗 -->
    <el-dialog
      v-model="serviceDialogVisible"
      title="联系客服"
      width="80%"
      center
    >
      <div class="service-dialog-content">
        <p class="service-tip">如有问题，请拨打客服热线：</p>
        <p class="service-phone">400-888-9999</p>
        <p class="service-time">服务时间：周一至周日 8:00-20:00</p>
      </div>
      <template #footer>
        <el-button type="primary" @click="callService">
          <el-icon><Phone /></el-icon>
          立即拨打
        </el-button>
      </template>
    </el-dialog>

    <!-- 关于我们弹窗 -->
    <el-dialog
      v-model="aboutDialogVisible"
      title="关于我们"
      width="80%"
      center
    >
      <div class="about-dialog-content">
        <div class="app-logo">🏠</div>
        <h3 class="app-name">sCare 社区养老服务</h3>
        <p class="app-version">版本 1.0.0</p>
        <p class="app-desc">致力于为社区老人提供便捷、贴心的居家养老服务</p>
        <p class="app-copyright">© 2025 sCare. All rights reserved.</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowRight,
  HomeFilled,
  List,
  User,
  UserFilled,
  Phone
} from '@element-plus/icons-vue'
import { useUserStore, useTokenStore } from '@/store'
import { authAPI } from '@/api'

const router = useRouter()
const userStore = useUserStore()
const tokenStore = useTokenStore()

// 弹窗状态
const serviceDialogVisible = ref(false)
const aboutDialogVisible = ref(false)

// 计算属性
const isLoggedIn = computed(() => tokenStore.isLoggedIn)

const displayName = computed(() => {
  if (!isLoggedIn.value) {
    return '点击登录'
  }
  return userStore.profile?.name || '未设置昵称'
})

const displayPhone = computed(() => {
  if (!isLoggedIn.value) {
    return '登录后查看更多功能'
  }
  const phone = userStore.user?.phone
  if (phone && phone.length === 11) {
    return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
  }
  return phone || ''
})

const avatarText = computed(() => {
  const name = userStore.profile?.name
  if (name) {
    return name.slice(0, 1)
  }
  return ''
})

// 用户卡片点击
const handleUserCardClick = () => {
  if (isLoggedIn.value) {
    router.push('/profile')
  } else {
    router.push('/login')
  }
}

// 导航方法
const goToHome = () => {
  router.push('/home')
}

const goToServices = () => {
  router.push('/services')
}

const goToAddress = () => {
  router.push('/profile#address')
}

const goToSettings = () => {
  router.push('/settings')
}

const goToNotifications = () => {
  router.push('/notifications')
}

// 弹窗方法
const showServiceDialog = () => {
  serviceDialogVisible.value = true
}

const showAboutDialog = () => {
  aboutDialogVisible.value = true
}

// 获取用户信息
onMounted(async () => {
  if (isLoggedIn.value && !userStore.profile) {
    try {
      const result = await authAPI.checkToken()
      userStore.setUser(result.user)
      if (result.profile) {
        userStore.setProfile(result.profile)
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
    }
  }
})

const callService = () => {
  window.location.href = 'tel:400-888-9999'
}
</script>

<style scoped>
.mine-container {
  min-height: 100vh;
  background: #f5f7fa;
}

/* 用户信息卡片 */
.user-card {
  display: flex;
  align-items: center;
  padding: 32px 20px;
  background: linear-gradient(135deg, #409EFF 0%, #53a8ff 50%, #66b1ff 100%);
  color: white;
  cursor: pointer;
  border-radius: 0 0 24px 24px;
}

.user-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  flex-shrink: 0;
}

.user-avatar .el-icon {
  color: white;
}

.avatar-text {
  font-size: calc(28px * var(--font-scale));
  font-weight: bold;
  color: white;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: calc(20px * var(--font-scale));
  font-weight: bold;
  margin-bottom: 4px;
}

.user-phone {
  font-size: calc(14px * var(--font-scale));
  opacity: 0.9;
}

.arrow-icon {
  font-size: calc(20px * var(--font-scale));
  opacity: 0.8;
}

/* 功能入口列表 */
.menu-section {
  background: white;
  margin-top: 12px;
}

.menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:active {
  background-color: #f5f7fa;
}

.menu-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-icon {
  font-size: calc(20px * var(--font-scale));
}

.menu-text {
  font-size: calc(16px * var(--font-scale));
  color: #303133;
}

.menu-arrow {
  color: #c0c4cc;
  font-size: calc(16px * var(--font-scale));
}

/* 底部占位 */
.bottom-placeholder {
  height: calc(60px + env(safe-area-inset-bottom, 20px));
}

/* 底部导航 */
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  justify-content: space-around;
  padding: 10px 0 calc(10px + env(safe-area-inset-bottom, 20px));
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: #909399;
  font-size: calc(12px * var(--font-scale));
}

.nav-item .el-icon {
  font-size: calc(24px * var(--font-scale));
}

.nav-item.active {
  color: #409EFF;
}

/* 客服弹窗 */
.service-dialog-content {
  text-align: center;
  padding: 16px 0;
}

.service-tip {
  font-size: calc(14px * var(--font-scale));
  color: #606266;
  margin-bottom: 12px;
}

.service-phone {
  font-size: calc(28px * var(--font-scale));
  font-weight: bold;
  color: #409EFF;
  margin-bottom: 8px;
}

.service-time {
  font-size: calc(12px * var(--font-scale));
  color: #909399;
}

/* 关于我们弹窗 */
.about-dialog-content {
  text-align: center;
  padding: 16px 0;
}

.app-logo {
  font-size: calc(64px * var(--font-scale));
  margin-bottom: 12px;
}

.app-name {
  font-size: calc(20px * var(--font-scale));
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.app-version {
  font-size: calc(14px * var(--font-scale));
  color: #909399;
  margin-bottom: 16px;
}

.app-desc {
  font-size: calc(14px * var(--font-scale));
  color: #606266;
  margin-bottom: 16px;
  line-height: 1.6;
}

.app-copyright {
  font-size: calc(12px * var(--font-scale));
  color: #c0c4cc;
}
</style>
