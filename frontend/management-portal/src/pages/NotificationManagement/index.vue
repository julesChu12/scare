<template>
  <div class="notification-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>通知管理</h3>
            <p class="subtitle">系统通知列表</p>
          </div>
          <div class="header-actions">
            <el-button
              type="primary"
              :icon="Refresh"
              :loading="loading"
              @click="loadNotifications"
            >
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 通知列表 -->
      <el-table
        v-loading="loading"
        :data="notificationList"
        stripe
        style="width: 100%"
        empty-text="暂无通知"
      >
        <!-- ID -->
        <el-table-column prop="id" label="ID" width="80" />

        <!-- 标题 -->
        <el-table-column prop="title" label="标题" min-width="200" />

        <!-- 类型 -->
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 内容 -->
        <el-table-column label="内容" min-width="250">
          <template #default="{ row }">
            <el-tooltip :content="row.content" placement="top" :show-after="500">
              <span class="content-text">{{ row.content }}</span>
            </el-tooltip>
          </template>
        </el-table-column>

        <!-- 状态 -->
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_read ? 'info' : 'warning'" size="small">
              {{ row.is_read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 创建时间 -->
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <!-- 阅读时间 -->
        <el-table-column label="阅读时间" width="160">
          <template #default="{ row }">
            {{ row.read_at ? formatDateTime(row.read_at) : '-' }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_read"
              type="primary"
              size="small"
              link
              @click="handleMarkRead(row)"
            >
              标记已读
            </el-button>
            <el-button
              type="info"
              size="small"
              link
              @click="handleViewDetail(row)"
            >
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="通知详情"
      width="500px"
    >
      <el-descriptions v-if="currentNotification" :column="1" border>
        <el-descriptions-item label="标题">
          {{ currentNotification.title }}
        </el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="getTypeTag(currentNotification.type)" size="small">
            {{ getTypeText(currentNotification.type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="内容">
          {{ currentNotification.content }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="currentNotification.is_read ? 'info' : 'warning'" size="small">
            {{ currentNotification.is_read ? '已读' : '未读' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(currentNotification.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item v-if="currentNotification.read_at" label="阅读时间">
          {{ formatDateTime(currentNotification.read_at) }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button
          v-if="currentNotification && !currentNotification.is_read"
          type="primary"
          @click="handleMarkReadAndClose"
        >
          标记已读
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { notificationApi } from '@/api'
import type { Notification } from '@/types/api'
import dayjs from 'dayjs'

// 加载状态
const loading = ref(false)

// 通知列表
const notificationList = ref<Notification[]>([])

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 详情对话框
const detailDialogVisible = ref(false)
const currentNotification = ref<Notification | null>(null)

/**
 * 加载通知列表
 */
async function loadNotifications() {
  try {
    loading.value = true
    const response = await notificationApi.getNotifications({
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    const { items, total } = response.data
    pagination.total = total
    notificationList.value = items
  } catch (error) {
    console.error('Failed to load notifications:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 标记已读
 */
async function handleMarkRead(notification: Notification) {
  try {
    await notificationApi.markAsRead(notification.id)
    ElMessage.success('已标记为已读')
    loadNotifications()
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

/**
 * 查看详情
 */
function handleViewDetail(notification: Notification) {
  currentNotification.value = notification
  detailDialogVisible.value = true
}

/**
 * 标记已读并关闭
 */
async function handleMarkReadAndClose() {
  if (!currentNotification.value) return

  try {
    await notificationApi.markAsRead(currentNotification.value.id)
    ElMessage.success('已标记为已读')
    detailDialogVisible.value = false
    loadNotifications()
  } catch (error) {
    console.error('Failed to mark as read:', error)
  }
}

/**
 * 获取类型文本
 */
function getTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    system: '系统通知',
    task: '任务通知',
    request: '需求通知',
    alert: '告警通知',
  }
  return typeMap[type] || type
}

/**
 * 获取类型标签颜色
 */
function getTypeTag(type: string): string {
  const tagMap: Record<string, string> = {
    system: 'info',
    task: 'primary',
    request: 'success',
    alert: 'danger',
  }
  return tagMap[type] || ''
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadNotifications()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadNotifications()
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped lang="scss">
.notification-management {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    h3 {
      margin: 0 0 5px 0;
      font-size: 18px;
      font-weight: 500;
    }

    .subtitle {
      margin: 0;
      font-size: 13px;
      color: #909399;
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .content-text {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
