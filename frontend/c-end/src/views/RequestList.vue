<template>
  <div class="request-list-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>我的服务</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <!-- 筛选器 -->
      <div class="filters">
        <el-select v-model="statusFilter" placeholder="全部状态" @change="loadRequests">
          <el-option label="全部" value="" />
          <el-option label="待受理" value="pending" />
          <el-option label="已派发" value="dispatched" />
          <el-option label="已认领" value="claimed" />
          <el-option label="处理中" value="processing" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
      </div>

      <!-- 服务请求列表 -->
      <!-- 骨架屏 -->
      <div v-if="loading && requests.length === 0" class="skeleton-list">
        <div v-for="i in 3" :key="i" class="skeleton-card">
          <div class="skeleton-header">
            <div class="skeleton-number"></div>
            <div class="skeleton-tag"></div>
          </div>
          <div class="skeleton-body">
            <div class="skeleton-row"></div>
            <div class="skeleton-row short"></div>
            <div class="skeleton-row"></div>
          </div>
        </div>
      </div>

      <div v-else-if="requests.length === 0" class="empty">
        <el-icon class="empty-icon"><Document /></el-icon>
        <p class="empty-text">暂无服务记录</p>
        <p class="empty-hint">您还没有申请过服务</p>
        <el-button type="primary" round @click="goToQuickStart">发起服务请求</el-button>
      </div>

      <div v-else class="request-list">
        <el-card v-for="request in requests" :key="request.id" class="request-card" @click="viewDetail(request.id)">
          <div class="request-header">
            <span class="request-number">{{ request.request_no }}</span>
            <el-tag :type="getStatusType(request.status)">{{ getStatusText(request.status) }}</el-tag>
          </div>
          <div class="request-body">
            <div class="info-row">
              <span class="label">服务类型：</span>
              <span>{{ getServiceTypeText(request.service_type) }}</span>
            </div>
            <div class="info-row" v-if="request.description">
              <span class="label">服务描述：</span>
              <span class="description">{{ request.description }}</span>
            </div>
            <div class="info-row">
              <span class="label">申请时间：</span>
              <span>{{ formatTime(request.created_at) }}</span>
            </div>
          </div>
          <!-- 评价提示 -->
          <div v-if="request.status === 'completed' && !request.rating" class="rate-hint">
            <el-icon><Star /></el-icon>
            <span>点击评价服务</span>
          </div>
        </el-card>
      </div>
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
      <div class="nav-item" @click="goToProfile">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Star, HomeFilled, List, User, Document } from '@element-plus/icons-vue'
import { requestsAPI, type ServiceRequest } from '@/api'
import { getServiceTypeName } from '@/config/serviceTypes'

const router = useRouter()

const loading = ref(false)
const requests = ref<ServiceRequest[]>([])
const statusFilter = ref('')

const goBack = () => {
  router.back()
}

const goToHome = () => {
  router.push('/home')
}

const goToProfile = () => {
  router.push('/profile')
}

const goToQuickStart = () => {
  router.push('/quick')
}

const viewDetail = (id: number) => {
  router.push(`/requests/${id}`)
}

const loadRequests = async () => {
  loading.value = true
  try {
    const result = await requestsAPI.getMyRequests({
      status: statusFilter.value || undefined
    })
    requests.value = result.list || []
  } catch (error) {
    console.error('加载服务请求失败:', error)
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    pending: 'warning',
    dispatched: 'info',
    claimed: 'primary',
    processing: 'primary',
    completed: 'success',
    cancelled: 'info'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待受理',
    dispatched: '已派发',
    claimed: '已认领',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return textMap[status] || status
}

const getServiceTypeText = (type: string) => {
  return getServiceTypeName(type)
}

const formatTime = (time: string) => {
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(() => {
  loadRequests()
})
</script>

<style scoped>
.request-list-container {
  min-height: 100vh;
  background: var(--bg-color, #f5f5f5);
  padding-bottom: 80px;
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

.filters {
  margin-bottom: 20px;
}

.loading,
.empty {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-color-secondary, #909399);
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon {
  font-size: 64px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #909399;
  margin-bottom: 8px;
}

.empty-hint {
  font-size: 14px;
  color: #c0c4cc;
  margin-bottom: 24px;
}

/* 骨架屏 */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.skeleton-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.skeleton-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.skeleton-number {
  width: 140px;
  height: 18px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

.skeleton-tag {
  width: 60px;
  height: 22px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

.skeleton-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.skeleton-row {
  width: 100%;
  height: 16px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

.skeleton-row.short {
  width: 60%;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.request-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.request-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.request-card:hover {
  transform: translateY(-2px);
}

.request-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.request-number {
  font-weight: bold;
  color: var(--color-primary, #409EFF);
  font-size: var(--font-size-base, 16px);
}

.request-body .info-row {
  margin-bottom: 8px;
  font-size: var(--font-size-base, 16px);
}

.info-row .label {
  color: var(--text-color-secondary, #909399);
  margin-right: 8px;
}

.info-row .description {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.rate-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  color: var(--color-warning, #E6A23C);
  font-size: var(--font-size-sm, 14px);
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
  padding: 12px 0;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-xs, 12px);
}

.nav-item .el-icon {
  font-size: 24px;
}

.nav-item.active {
  color: var(--color-primary, #409EFF);
}
</style>
