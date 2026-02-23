<template>
  <div class="notification-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>消息通知</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <!-- 骨架屏 -->
      <div v-if="loading && list.length === 0" class="skeleton-list">
        <div v-for="i in 5" :key="i" class="skeleton-item">
          <div class="skeleton-dot"></div>
          <div class="skeleton-body">
            <div class="skeleton-title"></div>
            <div class="skeleton-desc"></div>
            <div class="skeleton-date"></div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && list.length === 0" class="empty-state">
        <el-icon class="empty-icon"><Bell /></el-icon>
        <p class="empty-text">暂无通知</p>
      </div>

      <!-- 通知列表 -->
      <div v-else class="notification-list">
        <div
          v-for="item in list"
          :key="item.id"
          class="notification-item"
          :class="{ unread: !item.is_read }"
          @click="handleRead(item)"
        >
          <div class="dot" v-if="!item.is_read"></div>
          <div class="notification-body">
            <h4 class="notification-title">{{ item.title }}</h4>
            <p class="notification-desc">{{ item.content }}</p>
            <p class="notification-time">{{ formatTime(item.created_at) }}</p>
          </div>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="list.length > 0" class="load-more">
        <span v-if="loading" class="loading-text">
          <el-icon class="is-loading"><Loading /></el-icon>
          加载中...
        </span>
        <span v-else-if="noMore" class="no-more-text">没有更多了</span>
        <span v-else class="load-more-text" @click="loadMore">加载更多</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Bell, Loading } from '@element-plus/icons-vue'
import { notificationAPI, type Notification } from '@/api'

const router = useRouter()

const list = ref<Notification[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const noMore = computed(() => list.value.length >= total.value && total.value > 0)

const goBack = () => {
  router.back()
}

const fetchList = async (isLoadMore = false) => {
  if (loading.value) return
  loading.value = true

  try {
    const result = await notificationAPI.getList({
      page: page.value,
      page_size: pageSize
    })

    if (isLoadMore) {
      list.value = [...list.value, ...result.items]
    } else {
      list.value = result.items
    }
    total.value = result.total
  } catch (error) {
    console.error('获取通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (loading.value || noMore.value) return
  page.value++
  fetchList(true)
}

const handleRead = async (item: Notification) => {
  if (item.is_read) return
  try {
    await notificationAPI.markRead(item.id)
    item.is_read = true
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

const formatTime = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}天前`
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.notification-container {
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
  padding: 16px;
}

/* 骨架屏 */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
}

.skeleton-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  flex-shrink: 0;
  margin-top: 6px;
}

.skeleton-body {
  flex: 1;
}

.skeleton-title {
  width: 60%;
  height: 18px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 8px;
}

.skeleton-desc {
  width: 90%;
  height: 14px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 8px;
}

.skeleton-date {
  width: 30%;
  height: 12px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  background: white;
  border-radius: 12px;
}

.empty-icon {
  font-size: 64px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-text {
  font-size: var(--font-size-base, 16px);
  color: #909399;
}

/* 通知列表 */
.notification-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: transform 0.2s;
}

.notification-item:active {
  transform: scale(0.98);
}

.notification-item.unread {
  background: #f0f7ff;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #409EFF;
  flex-shrink: 0;
  margin-top: 6px;
}

.notification-body {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-size: var(--font-size-base, 16px);
  color: #303133;
  font-weight: 500;
  margin-bottom: 6px;
  line-height: 1.4;
}

.notification-item.unread .notification-title {
  font-weight: 600;
}

.notification-desc {
  font-size: var(--font-size-sm, 14px);
  color: #606266;
  margin-bottom: 6px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notification-time {
  font-size: 12px;
  color: #c0c4cc;
}

/* 加载更多 */
.load-more {
  text-align: center;
  padding: 16px;
  font-size: var(--font-size-sm, 14px);
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
</style>
