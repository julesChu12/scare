<template>
  <div class="request-detail-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>服务详情</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <div v-if="loading" class="loading">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载中...</p>
      </div>

      <div v-else-if="request" class="detail">
        <!-- 状态卡片 -->
        <div class="status-card" :class="'status-' + request.status">
          <div class="status-icon">{{ getServiceIcon(request.service_type) }}</div>
          <div class="status-info">
            <div class="status-type">{{ getServiceTypeText(request.service_type) }}</div>
            <div class="status-tag">{{ getStatusText(request.status) }}</div>
          </div>
        </div>

        <!-- 服务单号 -->
        <div class="order-no">
          <span class="label">服务单号</span>
          <span class="value">{{ request.request_no }}</span>
        </div>

        <!-- 服务信息 -->
        <div class="info-section">
          <div class="section-title">服务信息</div>
          <div class="info-list">
            <div class="info-row">
              <span class="info-label">联系人</span>
              <span class="info-value">{{ request.contact_name }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">联系电话</span>
              <span class="info-value">{{ request.contact_phone }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">服务地址</span>
              <span class="info-value">{{ request.address }}</span>
            </div>
            <div class="info-row" v-if="request.description">
              <span class="info-label">服务描述</span>
              <span class="info-value">{{ request.description }}</span>
            </div>
            <div class="info-row">
              <span class="info-label">申请时间</span>
              <span class="info-value">{{ formatTime(request.created_at) }}</span>
            </div>
            <div class="info-row" v-if="request.appointment_time">
              <span class="info-label">预约时间</span>
              <span class="info-value highlight">{{ formatTime(request.appointment_time) }}</span>
            </div>
          </div>
        </div>

        <!-- 服务站点 -->
        <div class="info-section" v-if="request.station_name">
          <div class="section-title">服务站点</div>
          <div class="info-list">
            <div class="info-row">
              <span class="info-label">站点名称</span>
              <span class="info-value">{{ request.station_name }}</span>
            </div>
          </div>
        </div>

        <!-- 服务照片 -->
        <div class="info-section" v-if="parsedImages.length">
          <div class="section-title">服务记录</div>
          <div class="photo-gallery">
            <el-image
              v-for="(photo, index) in parsedImages"
              :key="index"
              :src="photo"
              :preview-src-list="parsedImages"
              :initial-index="index"
              fit="cover"
              class="photo-item"
            >
              <template #error>
                <div class="photo-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </div>
        </div>

        <!-- 取消原因 -->
        <div class="info-section" v-if="request.status === 'cancelled' && request.reject_reason">
          <div class="section-title">取消原因</div>
          <div class="cancel-reason">{{ request.reject_reason }}</div>
        </div>

        <!-- 服务评价（已评价） -->
        <div class="info-section" v-if="request.rating">
          <div class="section-title">服务评价</div>
          <div class="rating-display">
            <el-rate v-model="request.rating" disabled show-score />
            <div class="feedback" v-if="request.feedback">{{ request.feedback }}</div>
          </div>
        </div>

        <!-- 底部操作区 -->
        <div class="action-area">
          <!-- 评价按钮（已完成但未评价） -->
          <el-button
            v-if="canRate"
            type="primary"
            size="large"
            class="action-btn"
            @click="showRatingDialog = true"
          >
            <el-icon><Star /></el-icon>
            评价服务
          </el-button>

          <!-- 取消按钮（待处理状态） -->
          <el-button
            v-if="canCancel"
            type="danger"
            plain
            size="large"
            class="action-btn"
            @click="handleCancel"
          >
            取消服务
          </el-button>

          <!-- 再次预约 -->
          <el-button
            v-if="request.status === 'completed' || request.status === 'cancelled'"
            type="primary"
            plain
            size="large"
            class="action-btn"
            @click="reorder"
          >
            再次预约
          </el-button>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="empty">
        <p>未找到服务记录</p>
        <el-button type="primary" @click="goBack">返回</el-button>
      </div>
    </div>

    <!-- 评价弹窗 -->
    <RatingDialog
      v-if="request"
      v-model="showRatingDialog"
      :request-id="request.id"
      @success="handleRatingSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Loading, Star, Picture } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { requestsAPI, type ServiceRequest } from '@/api'
import { getServiceTypeName, getServiceTypeIcon } from '@/config/serviceTypes'
import RatingDialog from '@/components/RatingDialog.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const request = ref<ServiceRequest | null>(null)
const showRatingDialog = ref(false)

// 是否可以评价（已完成且未评价）
const canRate = computed(() => {
  return request.value?.status === 'completed' && !request.value?.rating
})

// 是否可以取消（待处理状态）
const canCancel = computed(() => {
  return request.value?.status === 'pending'
})

const goBack = () => {
  router.back()
}

const getServiceIcon = (type: string) => {
  return getServiceTypeIcon(type)
}

// 解析图片字符串为数组
const parsedImages = computed(() => {
  if (!request.value?.images) return []
  try {
    const parsed = JSON.parse(request.value.images)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return request.value.images ? [request.value.images] : []
  }
})

const handleCancel = async () => {
  try {
    await ElMessageBox.confirm('确定要取消此服务吗？', '取消服务', {
      confirmButtonText: '确定取消',
      cancelButtonText: '再想想',
      type: 'warning'
    })

    await requestsAPI.cancelRequest(request.value!.id)
    ElMessage.success('服务已取消')
    loadDetail()
  } catch (e) {
    // 用户点击取消
  }
}

const reorder = () => {
  if (request.value) {
    router.push(`/quick?type=${request.value.service_type}`)
  }
}

const loadDetail = async () => {
  const id = parseInt(route.params.id as string, 10)
  if (!id) return

  loading.value = true
  try {
    request.value = await requestsAPI.getRequestDetail(id)
  } catch (error) {
    console.error('加载服务详情失败:', error)
  } finally {
    loading.value = false
  }
}

const handleRatingSuccess = (rating: number, feedback: string) => {
  if (request.value) {
    request.value.rating = rating
    request.value.feedback = feedback
  }
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    pending: '待处理',
    dispatched: '已派单',
    accepted: '已接单',
    in_progress: '服务中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝'
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
  loadDetail()
})
</script>

<style scoped>
.request-detail-container {
  min-height: 100vh;
  background: #f5f7fa;
}

.header {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header h1 {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.content {
  padding: 16px;
  padding-bottom: 100px;
}

.loading, .empty {
  text-align: center;
  padding: 60px 20px;
  color: #909399;
}

.detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 状态卡片 */
.status-card {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  color: white;
}

.status-card.status-completed {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
}

.status-card.status-pending {
  background: linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%);
}

.status-card.status-cancelled,
.status-card.status-rejected {
  background: linear-gradient(135deg, #909399 0%, #b4b4b4 100%);
}

.status-card.status-in_progress {
  background: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
}

.status-icon {
  font-size: 48px;
}

.status-info {
  flex: 1;
}

.status-type {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 4px;
}

.status-tag {
  font-size: 14px;
  opacity: 0.9;
}

/* 服务单号 */
.order-no {
  background: white;
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-no .label {
  color: #909399;
  font-size: 14px;
}

.order-no .value {
  color: #303133;
  font-size: 14px;
  font-family: monospace;
}

/* 信息区块 */
.info-section {
  background: white;
  border-radius: 12px;
  padding: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-row {
  display: flex;
  font-size: 14px;
}

.info-label {
  color: #909399;
  width: 80px;
  flex-shrink: 0;
}

.info-value {
  color: #303133;
  flex: 1;
  word-break: break-all;
}

.info-value.highlight {
  color: #409EFF;
  font-weight: 500;
}

/* 取消原因 */
.cancel-reason {
  color: #f56c6c;
  font-size: 14px;
  padding: 12px;
  background: #fef0f0;
  border-radius: 8px;
}

.photo-gallery {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.photo-item {
  width: 100%;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
}

.photo-error {
  width: 100%;
  height: 100%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 24px;
}

/* 评价展示 */
.rating-display {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rating-display .feedback {
  color: #606266;
  font-size: 14px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

/* 底部操作区 */
.action-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
  padding: 0;
}

.action-btn {
  width: 100%;
  font-size: 16px;
  height: 48px;
  border-radius: 24px;
  margin: 0 !important;
}
</style>
