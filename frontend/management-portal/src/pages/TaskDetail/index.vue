<template>
  <div class="task-detail">
    <el-card v-loading="loading">
      <!-- 页面头部 -->
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
            <div class="title-group">
              <h3>任务详情</h3>
              <el-tag v-if="task" :type="getStatusType(task.status)" size="large">
                {{ getStatusText(task.status) }}
              </el-tag>
            </div>
          </div>
          <div class="header-right">
            <!-- 认领按钮 -->
            <el-button
              v-if="task?.status === 'dispatched'"
              v-permission="'service:task:claim'"
              type="primary"
              size="large"
              :loading="claiming"
              @click="handleClaimTask"
            >
              认领任务
            </el-button>
            <!-- 完成按钮 -->
            <el-button
              v-if="task?.status === 'claimed'"
              v-permission="'service:task:complete'"
              type="success"
              size="large"
              @click="handleCompleteTask"
            >
              完成任务
            </el-button>
          </div>
        </div>
      </template>

      <!-- 任务信息 -->
      <div v-if="task && serviceRequest" class="detail-content">
        <el-row :gutter="20">
          <!-- 左侧：任务信息 -->
          <el-col :span="12">
            <div class="info-section">
              <h4 class="section-title">任务信息</h4>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="任务ID">{{ task.id }}</el-descriptions-item>
                <el-descriptions-item label="需求编号">
                  <el-tag type="info">{{ serviceRequest.request_no }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="服务类型">
                  <el-tag :type="getServiceTypeTag(serviceRequest.service_type)">
                    {{ getServiceTypeText(serviceRequest.service_type) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="优先级">
                  <el-tag :type="getPriorityTag(serviceRequest.priority)">
                    {{ getPriorityText(serviceRequest.priority) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="任务状态">
                  <el-tag :type="getStatusType(task.status)">
                    {{ getStatusText(task.status) }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="创建时间">
                  {{ formatDateTime(task.created_at) }}
                </el-descriptions-item>
                <el-descriptions-item v-if="task.claimed_at" label="认领时间">
                  {{ formatDateTime(task.claimed_at) }}
                </el-descriptions-item>
                <el-descriptions-item v-if="task.completed_at" label="完成时间">
                  {{ formatDateTime(task.completed_at) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>

            <!-- 工作人员备注 -->
            <div v-if="task.staff_notes" class="info-section">
              <h4 class="section-title">工作人员备注</h4>
              <el-alert type="info" :closable="false">
                {{ task.staff_notes }}
              </el-alert>
            </div>
          </el-col>

          <!-- 右侧：服务需求信息 -->
          <el-col :span="12">
            <div class="info-section">
              <h4 class="section-title">服务需求信息</h4>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="联系人">
                  {{ serviceRequest.contact_name }}
                </el-descriptions-item>
                <el-descriptions-item label="联系电话">
                  <a :href="`tel:${serviceRequest.contact_phone}`" class="phone-link">
                    <el-icon><Phone /></el-icon>
                    {{ serviceRequest.contact_phone }}
                  </a>
                </el-descriptions-item>
                <el-descriptions-item label="服务地址">
                  {{ serviceRequest.address }}
                </el-descriptions-item>
                <el-descriptions-item label="需求描述">
                  {{ serviceRequest.description || '无' }}
                </el-descriptions-item>
                <el-descriptions-item v-if="serviceRequest.scheduled_at" label="预约时间">
                  {{ formatDateTime(serviceRequest.scheduled_at) }}
                </el-descriptions-item>
                <el-descriptions-item label="提交时间">
                  {{ formatDateTime(serviceRequest.created_at) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>

            <!-- 地图 -->
            <div class="info-section">
              <h4 class="section-title">服务位置</h4>
              <div class="map-wrapper">
                <MapViewer
                  :longitude="serviceRequest.submit_location_lng"
                  :latitude="serviceRequest.submit_location_lat"
                />
              </div>
            </div>
          </el-col>
        </el-row>

        <!-- 服务照片 -->
        <div v-if="task.images && task.images.length > 0" class="info-section">
          <h4 class="section-title">服务照片</h4>
          <div class="image-list">
            <el-image
              v-for="(img, index) in task.images"
              :key="index"
              :src="getImageUrl(img)"
              :preview-src-list="task.images.map(i => getImageUrl(i))"
              :initial-index="index"
              preview-teleported
              fit="cover"
              class="service-image"
            />
          </div>
        </div>
      </div>

      <!-- 加载中或无数据 -->
      <el-empty v-else-if="!loading" description="任务不存在或已被删除" />
    </el-card>

    <!-- 完成任务对话框 -->
    <el-dialog
      v-model="showCompleteDialog"
      title="完成任务"
      width="600px"
    >
      <el-form :model="completeForm" label-width="100px">
        <!-- 服务照片上传 -->
        <el-form-item label="服务照片">
          <ImageUpload v-model="completeForm.images" :limit="9" />
          <div class="form-tip">请上传服务过程照片（最多9张，每张不超过5MB）</div>
        </el-form-item>

        <!-- 工作备注 -->
        <el-form-item label="工作备注">
          <el-input
            v-model="completeForm.staff_notes"
            type="textarea"
            :rows="4"
            placeholder="请输入服务过程中的备注信息（选填）"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCompleteDialog = false">取消</el-button>
        <el-button
          type="primary"
          :loading="completing"
          @click="submitCompleteTask"
        >
          确认完成
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Phone } from '@element-plus/icons-vue'
import { taskApi, requestApi } from '@/api'
import type { Task, ServiceRequest } from '@/types/api'
import dayjs from 'dayjs'
import ImageUpload from '@/components/ImageUpload.vue'
import MapViewer from '@/components/MapViewer.vue'

const router = useRouter()
const route = useRoute()

// 加载状态
const loading = ref(false)
const claiming = ref(false)

// 任务数据
const task = ref<Task | null>(null)
const serviceRequest = ref<ServiceRequest | null>(null)

// 完成任务对话框
const showCompleteDialog = ref(false)
const completing = ref(false)
const completeForm = ref({
  images: [] as string[],
  staff_notes: '',
})

/**
 * 加载任务详情
 */
async function loadTaskDetail() {
  try {
    loading.value = true

    const taskId = Number(route.params.id)
    if (!taskId) {
      throw new Error('Invalid task ID')
    }

    // 直接调用详情接口
    const res = await taskApi.getTask(taskId)
    task.value = res.data

    if (task.value) {
       // 加载服务需求详情
       if (task.value.request) {
         serviceRequest.value = task.value.request
       } else {
         const requestResponse = await requestApi.getRequest(task.value.request_id)
         serviceRequest.value = requestResponse.data
       }
    }

  } catch (error: any) {
    console.error('Failed to load task detail:', error)
    ElMessage.error('加载任务详情失败')
  } finally {
    loading.value = false
  }
}

/**
 * 返回上一页
 */
function goBack() {
  router.back()
}

/**
 * 认领任务
 */
async function handleClaimTask() {
  if (!task.value) return

  try {
    await ElMessageBox.confirm('确定要认领此任务吗？', '认领确认', {
      confirmButtonText: '确定认领',
      cancelButtonText: '取消',
      type: 'info',
    })

    claiming.value = true
    await taskApi.claimTask(task.value.id)

    ElMessage.success('认领成功！')

    // 重新加载任务详情
    await loadTaskDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to claim task:', error)
    }
  } finally {
    claiming.value = false
  }
}

/**
 * 完成任务
 */
async function handleCompleteTask() {
  if (!task.value) return

  // 显示完成任务对话框
  showCompleteDialog.value = true
}

/**
 * 提交完成任务
 */
async function submitCompleteTask() {
  if (!task.value) return

  try {
    completing.value = true

    await taskApi.completeTask(task.value.id, {
      images: completeForm.value.images,
      staff_notes: completeForm.value.staff_notes,
    })

    ElMessage.success('任务已完成！')

    // 关闭对话框
    showCompleteDialog.value = false
    await navigateAfterComplete()
  } catch (error: any) {
    console.error('Failed to complete task:', error)
  } finally {
    completing.value = false
  }
}

async function navigateAfterComplete() {
  const from = typeof route.query.from === 'string' ? route.query.from : ''
  const tab = route.query.tab === 'pool' || route.query.tab === 'my' ? route.query.tab : 'my'

  if (from === 'task-management') {
    await router.push({
      path: '/services/tasks',
      query: { tab },
    })
    return
  }

  await router.push({
    path: '/services/tasks',
    query: { tab: 'my' },
  })
}

/**
 * 获取任务状态文本
 */
function getStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    dispatched: '待认领',
    claimed: '已认领',
    processing: '服务中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return statusMap[status] || status
}

/**
 * 获取任务状态标签类型
 */
function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    dispatched: 'warning',
    claimed: 'primary',
    processing: 'warning',
    completed: 'success',
    cancelled: 'info',
  }
  return typeMap[status] || ''
}

/**
 * 获取服务类型文本
 */
function getServiceTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    meal: '送餐服务',
    medical: '医疗服务',
    cleaning: '清洁服务',
    shopping: '代购服务',
    accompany: '陪护服务',
  }
  return typeMap[type] || type
}

/**
 * 获取服务类型标签颜色
 */
function getServiceTypeTag(type: string): string {
  const tagMap: Record<string, string> = {
    meal: 'success',
    medical: 'danger',
    cleaning: 'warning',
    shopping: 'primary',
    accompany: 'info',
  }
  return tagMap[type] || ''
}

/**
 * 获取优先级文本
 */
function getPriorityText(priority?: string): string {
  const priorityMap: Record<string, string> = {
    low: '低',
    normal: '普通',
    high: '高',
    urgent: '紧急',
  }
  return priorityMap[priority || ''] || priority || '普通'
}

/**
 * 获取优先级标签颜色
 */
function getPriorityTag(priority?: string): string {
  const tagMap: Record<string, string> = {
    low: 'info',
    normal: 'primary',
    high: 'warning',
    urgent: 'danger',
  }
  return tagMap[priority || ''] || ''
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string | null): string {
  if (!dateTime) return '-'

  const parsed = dayjs(dateTime)
  if (!parsed.isValid() || parsed.year() <= 1) return '-'

  return parsed.format('YYYY-MM-DD HH:mm:ss')
}

/**
 * 获取图片完整URL
 */
function getImageUrl(path: string): string {
  if (path.startsWith('http')) return path
  // 开发环境通过 proxy 代理 /static，生产环境直接访问
  return path
}

// 组件挂载时加载数据
onMounted(() => {
  loadTaskDetail()
})
</script>

<style scoped lang="scss">
.task-detail {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .header-left {
      display: flex;
      align-items: center;
      gap: 16px;

      .title-group {
        display: flex;
        align-items: center;
        gap: 12px;

        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 500;
        }
      }
    }
  }

  .detail-content {
    .info-section {
      margin-bottom: 24px;

      .section-title {
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 500;
        color: #303133;
        border-left: 3px solid #409eff;
        padding-left: 12px;
      }
    }

    .phone-link {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      color: #409eff;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }

    .map-wrapper {
      height: 300px;
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      overflow: hidden;
    }

    .image-list {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;

      .service-image {
        width: 150px;
        height: 150px;
        border-radius: 4px;
        cursor: pointer;
      }
    }
  }

  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 8px;
  }
}
</style>
