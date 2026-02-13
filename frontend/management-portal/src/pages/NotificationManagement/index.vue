<template>
  <div class="notification-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>通知中心</h3>
            <p class="subtitle">查看系统公告、任务提醒及各类预警信息</p>
          </div>
          <div class="header-actions">
            <el-button :icon="Check" @click="handleMarkAllRead" :disabled="!hasUnread">
              全部标记已读
            </el-button>
            <el-button :icon="Refresh" :loading="loading" @click="loadNotifications">
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="filter-bar">
        <el-radio-group v-model="filterType" @change="handleFilter">
          <el-radio-button value="">全部</el-radio-button>
          <el-radio-button value="system">系统</el-radio-button>
          <el-radio-button value="task">任务</el-radio-button>
          <el-radio-button value="alert">告警</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 通知列表 -->
      <el-table
        v-loading="loading"
        :data="notificationList"
        stripe
        style="width: 100%"
        :row-class-name="tableRowClassName"
      >
        <!-- ID -->
        <el-table-column prop="id" label="ID" width="80" />

        <!-- 类型图标 -->
        <el-table-column width="50" align="center">
          <template #default="{ row }">
            <el-icon :class="['type-icon', row.type]">
              <component :is="getTypeIcon(row.type)" />
            </el-icon>
          </template>
        </el-table-column>

        <!-- 标题 -->
        <el-table-column prop="title" label="标题" min-width="200">
          <template #default="{ row }">
            <span :class="{ 'unread-title': !row.is_read }">{{ row.title }}</span>
            <el-tag v-if="!row.is_read" size="small" type="danger" effect="dark" style="margin-left: 8px">NEW</el-tag>
          </template>
        </el-table-column>

        <!-- 内容摘要 -->
        <el-table-column label="内容" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="content-preview">{{ row.content }}</span>
          </template>
        </el-table-column>

        <!-- 时间 -->
        <el-table-column label="接收时间" width="180">
          <template #default="{ row }">
            <span class="time-text">{{ formatDateTime(row.created_at) }}</span>
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="150" fixed="right" align="center">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              link
              @click="handleViewDetail(row)"
            >
              详情
            </el-button>
            <el-button
              v-if="!row.is_read"
              type="success"
              size="small"
              link
              @click="handleMarkRead(row)"
            >
              标记已读
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
      width="600px"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Check, Bell, Message, Warning, List } from '@element-plus/icons-vue'
import { notificationApi } from '@/api'
import type { Notification } from '@/types/api'
import dayjs from 'dayjs'

// 加载状态
const loading = ref(false)

// 筛选
const filterType = ref('')

// 通知列表
const notificationList = ref<Notification[]>([])

// 计算是否有未读
const hasUnread = computed(() => notificationList.value.some(n => !n.is_read))

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

/**
 * 获取图标
 */
function getTypeIcon(type: string) {
  const map: Record<string, any> = {
    system: Bell,
    task: List,
    alert: Warning,
    message: Message
  }
  return map[type] || Bell
}

/**
 * 筛选
 */
function handleFilter() {
  pagination.page = 1
  loadNotifications()
}

/**
 * 全部标记已读
 */
async function handleMarkAllRead() {
  try {
    await ElMessageBox.confirm('确定将所有通知标记为已读吗？', '确认操作', {
      type: 'info'
    })
    
    // 由于后端可能没有批量接口，我们对未读列表进行并行操作
    const unreads = notificationList.value.filter(n => !n.is_read)
    await Promise.all(unreads.map(n => notificationApi.markAsRead(n.id)))
    
    ElMessage.success('已全部标记为已读')
    loadNotifications()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to mark all as read:', error)
    }
  }
}

/**
 * 表格行样式
 */
function tableRowClassName({ row }: { row: Notification }) {
  if (!row.is_read) {
    return 'unread-row'
  }
  return ''
}

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
      type: filterType.value || undefined
    } as any)
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

    .header-left {
      h3 {
        margin: 0 0 4px 0;
        font-size: 18px;
        font-weight: 500;
      }
      .subtitle {
        margin: 0;
        font-size: 13px;
        color: #909399;
      }
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .filter-bar {
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f0f0f0;
  }

  .type-icon {
    font-size: 18px;
    &.system { color: #409eff; }
    &.task { color: #67c23a; }
    &.alert { color: #f56c6c; }
    &.message { color: #e6a23c; }
  }

  .unread-title {
    font-weight: 600;
    color: #303133;
  }

  .content-preview {
    color: #606266;
    font-size: 13px;
  }

  .time-text {
    font-size: 12px;
    color: #909399;
  }

  :deep(.unread-row) {
    background-color: #fdf6ec !important;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
