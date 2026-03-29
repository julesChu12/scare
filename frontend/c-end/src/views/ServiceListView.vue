<template>
  <div class="service-list-container">
    <!-- 顶部分类筛选 - 图标网格 -->
    <div class="category-filter">
      <div class="category-grid">
        <div
          v-for="category in categories"
          :key="category.key"
          class="category-item"
          :class="{ active: selectedCategory === category.key }"
          @click="selectCategory(category.key)"
        >
          <span class="category-icon">{{ category.icon || '📋' }}</span>
          <span class="category-name">{{ category.name }}</span>
        </div>
      </div>
    </div>

    <!-- 服务记录列表 -->
    <div class="content" ref="contentRef">
      <!-- 未登录提示 -->
      <div v-if="!isLoggedIn" class="login-prompt">
        <el-icon class="prompt-icon"><User /></el-icon>
        <p class="prompt-text">登录后查看您的服务记录</p>
        <el-button type="primary" round @click="goToLogin">立即登录</el-button>
      </div>

      <!-- 已登录：服务记录列表 -->
      <template v-else>
        <!-- 下拉刷新提示 -->
        <div v-if="refreshing" class="refresh-tip">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>刷新中...</span>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && requests.length === 0" class="empty-state">
          <el-icon class="empty-icon"><Document /></el-icon>
          <p class="empty-text">暂无服务记录</p>
          <p class="empty-hint">您还没有申请过服务</p>
          <el-button type="primary" round @click="goToQuickStart">发起服务请求</el-button>
        </div>

        <!-- 骨架屏 -->
        <div v-else-if="loading && requests.length === 0" class="skeleton-list">
          <div v-for="i in 4" :key="i" class="skeleton-item">
            <div class="skeleton-icon"></div>
            <div class="skeleton-content">
              <div class="skeleton-header">
                <div class="skeleton-name"></div>
                <div class="skeleton-status"></div>
              </div>
              <div class="skeleton-meta"></div>
            </div>
          </div>
        </div>

        <!-- 服务记录列表 -->
        <div v-else class="request-list">
          <div
            v-for="request in requests"
            :key="request.id"
            class="request-item"
            @click="goToDetail(request.id)"
          >
            <div class="request-icon">{{ getServiceTypeIcon(request.service_type) }}</div>
            <div class="request-info">
              <div class="request-header">
                <span class="request-name">{{ getServiceTypeName(request.service_type) }}</span>
                <span class="request-status" :class="getStatusClass(request.status)">
                  {{ getStatusText(request.status) }}
                </span>
              </div>
              <div class="request-meta">
                <span class="request-no">{{ request.request_no }}</span>
                <span class="request-time">{{ formatTime(request.created_at) }}</span>
              </div>
            </div>
            <el-icon class="request-arrow"><ArrowRight /></el-icon>
          </div>
        </div>

        <!-- 加载更多 -->
        <div v-if="requests.length > 0" class="load-more">
          <span v-if="loading" class="loading-text">
            <el-icon class="is-loading"><Loading /></el-icon>
            加载中...
          </span>
          <span v-else-if="noMore" class="no-more-text">没有更多了</span>
          <span v-else class="load-more-text" @click="loadMore">加载更多</span>
        </div>
      </template>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <div class="nav-item" @click="goToHome">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </div>
      <div class="nav-item active">
        <el-icon><List /></el-icon>
        <span>服务</span>
      </div>
      <div class="nav-item" @click="goToMine">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowRight,
  HomeFilled,
  List,
  User,
  Loading,
  Document
} from '@element-plus/icons-vue'
import { useUserStore } from '@/store/userStore'
import { useTokenStore } from '@/store/tokenStore'
import { requestsAPI, type ServiceRequest } from '@/api/requests'
import { getServiceTypeName, getServiceTypeIcon } from '@/config/serviceTypes'

const router = useRouter()
const userStore = useUserStore()
const tokenStore = useTokenStore()

// 分类数据 - 显示常用的9种 + 全部 = 10个（2行5列）
const filterCategories = [
  { key: '', name: '全部', icon: '📋' },
  { key: 'meal', name: '送餐', icon: '🍱' },
  { key: 'medical', name: '陪护', icon: '🏥' },
  { key: 'care', name: '照护', icon: '👴' },
  { key: 'cleaning', name: '保洁', icon: '🧹' },
  { key: 'company', name: '陪聊', icon: '💬' },
  { key: 'repair', name: '维修', icon: '🔧' },
  { key: 'shopping', name: '代购', icon: '🛒' },
  { key: 'transport', name: '接送', icon: '🚗' },
  { key: 'other', name: '其他', icon: '📦' }
]

const categories = computed(() => filterCategories)

// 状态
const selectedCategory = ref('')
const requests = ref<ServiceRequest[]>([])
const loading = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const noMore = computed(() => requests.value.length >= total.value && total.value > 0)
const contentRef = ref<HTMLElement | null>(null)

// 是否已登录
const isLoggedIn = computed(() => {
  return tokenStore.isLoggedIn || !!userStore.user
})

// 选择分类
const selectCategory = (key: string) => {
  if (selectedCategory.value === key) return
  selectedCategory.value = key
  page.value = 1
  requests.value = []
  fetchRequests()
}

// 获取服务请求列表
const fetchRequests = async (isRefresh = false) => {
  if (!isLoggedIn.value) return

  if (isRefresh) {
    refreshing.value = true
    page.value = 1
  } else {
    loading.value = true
  }

  try {
    let items: ServiceRequest[] = []

    const params: { page: number; page_size: number; status?: string } = {
      page: page.value,
      page_size: pageSize
    }
    const result = await requestsAPI.getMyRequests(params)
    items = result.items || []
    total.value = result.total

    // 前端按服务类型筛选
    if (selectedCategory.value) {
      items = items.filter(item => item.service_type === selectedCategory.value)
    }

    if (isRefresh || page.value === 1) {
      requests.value = items
    } else {
      requests.value = [...requests.value, ...items]
    }
  } catch (error) {
    console.error('获取服务记录失败:', error)
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 加载更多
const loadMore = () => {
  if (loading.value || noMore.value) return
  page.value++
  fetchRequests()
}

// 状态文本映射（与后端 consts/status.go 保持一致）
const statusTextMap: Record<string, string> = {
  pending: '待受理',
  dispatched: '已派发',
  claimed: '已认领',
  processing: '处理中',
  completed: '已完成',
  cancelled: '已取消',
  rejected: '已拒绝'
}

const getStatusText = (status: string) => {
  return statusTextMap[status] || status
}

// 状态样式映射（与后端枚举保持一致）
const getStatusClass = (status: string) => {
  const classMap: Record<string, string> = {
    pending: 'status-pending',
    dispatched: 'status-dispatched',
    claimed: 'status-claimed',
    processing: 'status-processing',
    completed: 'status-completed',
    cancelled: 'status-cancelled',
    rejected: 'status-rejected'
  }
  return classMap[status] || ''
}

// 格式化时间
const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    return `今天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  } else if (days === 1) {
    return `昨天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return `${date.getMonth() + 1}月${date.getDate()}日`
  }
}

// 导航
const goToHome = () => router.push({ name: 'Home' })
const goToMine = () => router.push({ name: 'Mine' })
const goToLogin = () => router.push({ name: 'Login' })
const goToQuickStart = () => router.push({ name: 'QuickStart' })
const goToDetail = (id: number) => router.push({ name: 'RequestDetail', params: { id } })

// 监听滚动实现上拉加载
const handleScroll = () => {
  if (!contentRef.value || loading.value || noMore.value || !isLoggedIn.value) return

  const { scrollTop, scrollHeight, clientHeight } = document.documentElement
  if (scrollTop + clientHeight >= scrollHeight - 100) {
    loadMore()
  }
}

// 初始化
onMounted(() => {
  if (isLoggedIn.value) {
    fetchRequests()
  }
  window.addEventListener('scroll', handleScroll)
})

// 监听登录状态变化
watch(isLoggedIn, (newVal) => {
  if (newVal) {
    fetchRequests()
  } else {
    requests.value = []
  }
})
</script>

<style scoped>
.service-list-container {
  min-height: 100vh;
  max-height: 100vh;
  background: #f5f7fa;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 分类筛选 - 图标网格 */
.category-filter {
  background: white;
  padding: 16px;
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.category-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 12px;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 10px 4px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:active {
  transform: scale(0.95);
}

.category-item.active {
  background: #ecf5ff;
}

.category-item.active .category-icon {
  transform: scale(1.1);
}

.category-item.active .category-name {
  color: #409EFF;
  font-weight: 500;
}

.category-icon {
  font-size: calc(28px * var(--font-scale));
  transition: transform 0.2s;
}

.category-name {
  font-size: calc(12px * var(--font-scale));
  color: #606266;
  text-align: center;
  line-height: 1.2;
}

/* 内容区域 */
.content {
  padding: 16px;
  padding-bottom: 80px; /* 为底部导航留出空间 */
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

/* 未登录提示 */
.login-prompt {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: white;
  border-radius: 12px;
  flex: 1;
}

.prompt-icon {
  font-size: calc(64px * var(--font-scale));
  color: #c0c4cc;
  margin-bottom: 16px;
}

.prompt-text {
  font-size: calc(16px * var(--font-scale));
  color: #909399;
  margin-bottom: 24px;
}

/* 刷新提示 */
.refresh-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  color: #909399;
  font-size: calc(14px * var(--font-scale));
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: white;
  border-radius: 12px;
  flex: 1;
}

.empty-icon {
  font-size: calc(64px * var(--font-scale));
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-text {
  font-size: calc(16px * var(--font-scale));
  color: #909399;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: calc(14px * var(--font-scale));
  color: #c0c4cc;
  margin-bottom: 24px;
}

/* 骨架屏 */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
}

.skeleton-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 12px;
  flex-shrink: 0;
}

.skeleton-content {
  flex: 1;
}

.skeleton-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.skeleton-name {
  width: 80px;
  height: 18px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

.skeleton-status {
  width: 50px;
  height: 18px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 10px;
}

.skeleton-meta {
  width: 60%;
  height: 14px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 服务记录列表 */
.request-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.request-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.request-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.request-item:active {
  transform: scale(0.98);
}

.request-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  border-radius: 12px;
  font-size: calc(24px * var(--font-scale));
  flex-shrink: 0;
}

.request-info {
  flex: 1;
  min-width: 0;
}

.request-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.request-name {
  font-size: calc(16px * var(--font-scale));
  font-weight: 500;
  color: #303133;
}

.request-status {
  font-size: calc(12px * var(--font-scale));
  padding: 2px 8px;
  border-radius: 10px;
  flex-shrink: 0;
}

.status-pending {
  background: #fef0f0;
  color: #f56c6c;
}

.status-dispatched,
.status-claimed {
  background: #fdf6ec;
  color: #e6a23c;
}

.status-processing {
  background: #ecf5ff;
  color: #409eff;
}

.status-completed {
  background: #f0f9eb;
  color: #67c23a;
}

.status-cancelled,
.status-rejected {
  background: #f4f4f5;
  color: #909399;
}

.request-meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: calc(13px * var(--font-scale));
  color: #909399;
}

.request-no {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.request-time {
  flex-shrink: 0;
}

.request-arrow {
  color: #c0c4cc;
  flex-shrink: 0;
}

/* 加载更多 */
.load-more {
  text-align: center;
  padding: 16px;
  font-size: calc(14px * var(--font-scale));
  color: #909399;
}

.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.load-more-text {
  color: #409EFF;
  cursor: pointer;
}

.no-more-text {
  color: #c0c4cc;
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
  padding: 10px 0 20px;
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
</style>
