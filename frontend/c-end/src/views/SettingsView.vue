<template>
  <div class="settings-container">
    <!-- Header -->
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>设置</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <!-- 字体大小设置 -->
      <el-card class="settings-card">
        <template #header>
          <div class="card-header">
            <el-icon><FontSize /></el-icon>
            <span>字体大小</span>
          </div>
        </template>
        <div class="font-size-options">
          <div
            v-for="(config, key) in fontSizeConfig"
            :key="key"
            class="font-option"
            :class="{ active: fontSize === key }"
            @click="setFontSize(key as FontSizeType)"
          >
            <span class="font-label" :style="{ fontSize: config.value }">{{ config.label }}</span>
            <el-icon v-if="fontSize === key" class="check-icon"><Check /></el-icon>
          </div>
        </div>
        <div class="preview-text">
          <p>预览效果：这是一段示例文字</p>
        </div>
      </el-card>

      <!-- 通知设置 -->
      <el-card class="settings-card">
        <template #header>
          <div class="card-header">
            <el-icon><Bell /></el-icon>
            <span>通知设置</span>
          </div>
        </template>
        <div class="setting-item">
          <span>接收服务通知</span>
          <el-switch v-model="notificationEnabled" @change="saveNotificationSetting" />
        </div>
      </el-card>

      <!-- 缓存管理 -->
      <el-card class="settings-card">
        <template #header>
          <div class="card-header">
            <el-icon><Delete /></el-icon>
            <span>缓存管理</span>
          </div>
        </template>
        <div class="setting-item clickable" @click="clearCache">
          <span>清除缓存</span>
          <el-icon><ArrowRight /></el-icon>
        </div>
      </el-card>

      <!-- 退出登录 -->
      <el-button
        type="danger"
        size="large"
        class="logout-btn"
        @click="handleLogout"
      >
        退出登录
      </el-button>

      <!-- 版本信息 -->
      <div class="version-info">
        <p>sCare 社区养老服务</p>
        <p>版本 1.0.0</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, ArrowRight, Check, Bell, Delete } from '@element-plus/icons-vue'
import { useTokenStore, useUserStore } from '@/store'
import { authAPI } from '@/api'
import { useFontSize, type FontSize as FontSizeType, fontSizeConfig } from '@/composables/useFontSize'

// 自定义 FontSize 图标组件
const FontSize = {
  template: `<svg viewBox="0 0 24 24" fill="currentColor"><path d="M9 4v3h5v12h3V7h5V4H9zm-6 8h3v7h3v-7h3V9H3v3z"/></svg>`
}

const router = useRouter()
const tokenStore = useTokenStore()
const userStore = useUserStore()

// 字体大小
const { fontSize, setFontSize } = useFontSize()

// 通知设置
const NOTIFICATION_KEY = 'c_notification_enabled'
const notificationEnabled = ref(localStorage.getItem(NOTIFICATION_KEY) !== 'false')

const saveNotificationSetting = (value: boolean) => {
  localStorage.setItem(NOTIFICATION_KEY, String(value))
  ElMessage.success(value ? '已开启通知' : '已关闭通知')
}

// 返回
const goBack = () => {
  router.back()
}

// 清除缓存
const clearCache = async () => {
  try {
    await ElMessageBox.confirm(
      '清除缓存将删除本地存储的临时数据，不会影响您的账号信息。确定要清除吗？',
      '清除缓存',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 保留需要保留的数据
    const token = localStorage.getItem('c_token')
    const refreshToken = localStorage.getItem('c_refresh_token')
    const fontSizeSetting = localStorage.getItem('c_font_size')
    const notificationSetting = localStorage.getItem('c_notification_enabled')

    // 清除所有 localStorage
    localStorage.clear()

    // 恢复保留的数据
    if (token) localStorage.setItem('c_token', token)
    if (refreshToken) localStorage.setItem('c_refresh_token', refreshToken)
    if (fontSizeSetting) localStorage.setItem('c_font_size', fontSizeSetting)
    if (notificationSetting) localStorage.setItem('c_notification_enabled', notificationSetting)

    ElMessage.success('缓存已清除')
  } catch {
    // 用户取消
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要退出登录吗？',
      '退出登录',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 调用后端使 Token 进入黑名单
    try {
      await authAPI.logout()
    } catch {
      // 即使后端调用失败也继续清除本地状态
    }
    tokenStore.clearToken()
    userStore.clearUser()
    router.push('/login')
    ElMessage.success('已退出登录')
  } catch {
    // 用户取消
  }
}
</script>

<style scoped>
.settings-container {
  min-height: 100vh;
  background: var(--bg-color, #f5f5f5);
}

.header {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header h1 {
  font-size: var(--font-size-subtitle, 18px);
  font-weight: bold;
}

.content {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.settings-card {
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: var(--font-size-base, 16px);
  font-weight: 500;
}

.font-size-options {
  display: flex;
  gap: 12px;
}

.font-option {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 16px;
  border: 2px solid var(--border-color, #DCDFE6);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.font-option:hover {
  border-color: var(--color-primary, #409EFF);
}

.font-option.active {
  border-color: var(--color-primary, #409EFF);
  background: rgba(64, 158, 255, 0.1);
}

.font-label {
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.check-icon {
  position: absolute;
  top: 8px;
  right: 8px;
  color: var(--color-primary, #409EFF);
}

.preview-text {
  margin-top: 16px;
  padding: 12px;
  background: var(--bg-color, #f5f5f5);
  border-radius: 8px;
  font-size: var(--font-size-base, 16px);
  color: var(--text-color-regular, #606266);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  font-size: var(--font-size-base, 16px);
}

.setting-item.clickable {
  cursor: pointer;
}

.setting-item.clickable:hover {
  color: var(--color-primary, #409EFF);
}

.logout-btn {
  width: 100%;
  margin-top: 24px;
  font-size: var(--font-size-base, 16px);
}

.version-info {
  text-align: center;
  margin-top: 40px;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
}

.version-info p {
  margin: 4px 0;
}
</style>
