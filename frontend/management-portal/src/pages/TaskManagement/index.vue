<template>
  <div class="task-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>任务管理</h3>
            <p class="subtitle">所有任务列表及管理</p>
          </div>
          <el-button
            type="primary"
            :icon="Refresh"
            :loading="loading"
            @click="loadTasks"
          >
            刷新
          </el-button>
        </div>
      </template>

      <!-- 筛选栏 -->
      <div class="filter-container">
        <el-form :inline="true" :model="filters" class="demo-form-inline">
          <el-form-item label="需求编号">
            <el-input v-model="filters.request_no" placeholder="输入需求编号" clearable />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 150px">
              <el-option label="已分发" value="dispatched" />
              <el-option label="已认领" value="claimed" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
          </el-form-item>
          <el-form-item label="服务类型">
             <el-select v-model="filters.service_type" placeholder="全部类型" clearable style="width: 150px">
              <el-option label="送餐服务" value="meal" />
              <el-option label="医疗服务" value="medical" />
              <el-option label="清洁服务" value="cleaning" />
              <el-option label="代购服务" value="shopping" />
              <el-option label="陪护服务" value="accompany" />
            </el-select>
          </el-form-item>
          <el-form-item label="所属站点" v-if="isSystemAdmin">
             <el-select v-model="filters.station_id" placeholder="全部站点" clearable filterable style="width: 200px">
              <el-option
                v-for="station in stationList"
                :key="station.id"
                :label="station.name"
                :value="station.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleFilter">查询</el-button>
            <el-button @click="resetFilter">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 任务视图切换 -->
      <el-tabs v-model="activeTab" class="task-tabs">
        <el-tab-pane label="任务池" name="pool" v-if="canViewTaskPool" />
        <el-tab-pane label="我的任务" name="my" v-if="canViewMyTasks" />
      </el-tabs>

      <!-- 任务列表 -->
      <el-table
        v-loading="loading"
        :data="taskList"
        stripe
        style="width: 100%"
        empty-text="暂无任务"
      >
        <!-- 任务ID -->
        <el-table-column prop="id" label="任务ID" width="80" />

        <!-- 服务需求编号 -->
        <el-table-column label="需求编号" width="180">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.request?.request_no || '-' }}</el-tag>
          </template>
        </el-table-column>

        <!-- 服务类型 -->
        <el-table-column label="服务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getServiceTypeTag(row.request?.service_type)">
              {{ getServiceTypeText(row.request?.service_type) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 状态 -->
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 联系人 -->
        <el-table-column label="联系人" width="150">
          <template #default="{ row }">
            <div class="contact-info">
              <div>{{ row.request?.contact_name || '-' }}</div>
              <div class="phone">{{ row.request?.contact_phone || '-' }}</div>
            </div>
          </template>
        </el-table-column>

        <!-- 服务地址 -->
        <el-table-column label="服务地址" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.request?.address || '-' }}
          </template>
        </el-table-column>
        
        <!-- 服务人员（仅任务池需要显示） -->
         <el-table-column v-if="activeTab !== 'my'" label="服务人员" width="120">
           <template #default="{ row }">
              <span v-if="row.staff_id">{{ row.staff_name || '已指派' }}</span>
              <span v-else class="text-gray">-</span>
           </template>
         </el-table-column>

        <!-- 创建/更新时间 -->
        <el-table-column label="更新时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.updated_at) }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              link
              @click="handleViewDetail(row)"
            >
              查看
            </el-button>
            <el-button
              v-if="canManageTask(row)"
              v-permission="'service:task:assign'"
              type="primary"
              size="small"
              link
              @click="handleAssign(row)"
            >
              {{ row.status === 'dispatched' ? '指派' : '转派' }}
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

    <!-- 指派/转派弹窗 -->
    <el-dialog
      v-model="assignDialogVisible"
      :title="currentTask?.status === 'dispatched' ? '任务指派' : '任务转派'"
      width="500px"
    >
      <el-form label-width="100px">
        <el-form-item label="选择服务人员">
          <el-select
            v-model="selectedStaffId"
            placeholder="请选择服务人员"
            filterable
            style="width: 100%"
            :loading="staffLoading"
          >
            <el-option
              v-for="staff in staffList"
              :key="staff.id"
              :label="`${staff.name} (${staff.phone})`"
              :value="staff.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="assignLoading"
          :disabled="!selectedStaffId"
          @click="submitAssign"
        >
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { taskApi, stationApi, userApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import type { Task, ServiceRequest, Station, User } from '@/types/api'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 判断是否为系统管理员
const isSystemAdmin = computed(() => {
  return authStore.user?.roles?.includes('admin') || false
})

// 判断是否为站点管理员
const isStationManager = computed(() => {
  return authStore.user?.roles?.includes('station_manager') || false
})

// 判断是否为普通员工
const isRegularStaff = computed(() => {
  return authStore.user?.roles?.includes('staff') && !isSystemAdmin.value && !isStationManager.value
})

// 加载状态
const loading = ref(false)

// 任务列表（扩展类型包含 request 信息）
interface TaskWithRequest extends Task {
  request?: ServiceRequest
  staff_name?: string // 暂无后端返回，预留
}
const taskList = ref<TaskWithRequest[]>([])
const stationList = ref<Station[]>([])

// 筛选参数
const filters = reactive({
  request_no: '',
  status: '',
  service_type: '',
  station_id: undefined as number | undefined,
})

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 权限检查辅助函数
const hasPermission = (permission: string) => authStore.hasPermission(permission)
const canViewTaskPool = computed(() => hasPermission('service:task:pool'))
const canViewMyTasks = computed(() => hasPermission('service:task:my'))

function isAllowedTab(tab?: string | null): tab is 'pool' | 'my' {
  if (tab === 'pool') return canViewTaskPool.value
  if (tab === 'my') return canViewMyTasks.value
  return false
}

function resolveDefaultTab(): 'pool' | 'my' {
  const queryTab = typeof route.query.tab === 'string' ? route.query.tab : undefined
  if (isAllowedTab(queryTab)) {
    return queryTab
  }

  if (isRegularStaff.value && canViewMyTasks.value) {
    return 'my'
  }

  if (canViewTaskPool.value) {
    return 'pool'
  }

  if (canViewMyTasks.value) {
    return 'my'
  }

  return 'pool'
}

// 标签页状态
const activeTab = ref<'pool' | 'my'>(resolveDefaultTab())

// 监听标签页变化，重新加载数据
watch(activeTab, () => {
  pagination.page = 1
  loadTasks()
})


// 指派相关状态
const assignDialogVisible = ref(false)
const assignLoading = ref(false)
const staffLoading = ref(false)
const currentTask = ref<TaskWithRequest | null>(null)
const selectedStaffId = ref<number | undefined>(undefined)
const staffList = ref<User[]>([])

/**
 * 判断是否有权管理任务
 */
function canManageTask(task: Task): boolean {
  // 只有管理员和站点管理员可以指派
  if (!isSystemAdmin.value && !isStationManager.value) return false
  // 只有分发和已认领状态可以操作
  return ['dispatched', 'claimed'].includes(task.status)
}

/**
 * 打开指派弹窗
 */
async function handleAssign(task: TaskWithRequest) {
  currentTask.value = task
  selectedStaffId.value = task.staff_id || undefined
  assignDialogVisible.value = true
  
  // 加载该站点的服务人员
  staffLoading.value = true
  try {
    const res = await userApi.getUsers({
      page: 1,
      page_size: 100, // 获取足够多的人员
      role: 'staff',
      station_id: task.station_id
    } as any)
    staffList.value = res.data.items
  } catch (error) {
    console.error('Failed to load staff:', error)
    ElMessage.error('加载服务人员列表失败')
  } finally {
    staffLoading.value = false
  }
}

/**
 * 提交指派
 */
async function submitAssign() {
  if (!currentTask.value || !selectedStaffId.value) return
  
  assignLoading.value = true
  try {
    await taskApi.transferTask(currentTask.value.id, selectedStaffId.value)
    ElMessage.success('操作成功')
    assignDialogVisible.value = false
    loadTasks() // 刷新列表
  } catch (error) {
    console.error('Failed to assign task:', error)
    ElMessage.error('操作失败')
  } finally {
    assignLoading.value = false
  }
}

/**
 * 加载站点列表（仅管理员）
 */
async function loadStations() {
  if (!isSystemAdmin.value) return
  try {
    const res = await stationApi.getStations({ page: 1, page_size: 100 })
    stationList.value = res.data.items
  } catch (error) {
    console.error('Failed to load stations:', error)
  }
}

/**
 * 加载任务列表
 */
async function loadTasks() {
  try {
    loading.value = true

    let response
    const commonParams = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filters.status || undefined,
      service_type: filters.service_type || undefined,
      station_id: filters.station_id,
      request_no: filters.request_no || undefined,
    }

    if (activeTab.value === 'my') {
      response = await taskApi.getMyTasks(commonParams)
    } else {
      // 任务池 (涵盖所有可查看的任务，含指派给他人或已完成的)
      response = await taskApi.getTasks(commonParams)
    }

    const { items, total } = response.data
    pagination.total = total

    // 后端已返回 TaskWithRequest，包含 request 数据，直接使用
    taskList.value = items
  } catch (error) {
    console.error('Failed to load tasks:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 处理查询
 */
function handleFilter() {
  pagination.page = 1
  loadTasks()
}

/**
 * 重置筛选
 */
function resetFilter() {
  filters.request_no = ''
  filters.status = ''
  filters.service_type = ''
  filters.station_id = undefined
  handleFilter()
}

/**
 * 查看任务详情
 */
function handleViewDetail(task: TaskWithRequest) {
  router.push({
    path: `/services/tasks/${task.id}`,
    query: {
      from: 'task-management',
      tab: activeTab.value,
    },
  })
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadTasks()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadTasks()
}

/**
 * 获取服务类型文本
 */
function getServiceTypeText(type?: string): string {
  const typeMap: Record<string, string> = {
    meal: '送餐服务',
    medical: '医疗服务',
    cleaning: '清洁服务',
    shopping: '代购服务',
    accompany: '陪护服务',
  }
  return typeMap[type || ''] || type || '未知'
}

/**
 * 获取服务类型标签颜色
 */
function getServiceTypeTag(type?: string): string {
  const tagMap: Record<string, string> = {
    meal: 'success',
    medical: 'danger',
    cleaning: 'warning',
    shopping: 'primary',
    accompany: 'info',
  }
  return tagMap[type || ''] || ''
}

/**
 * 获取状态文本
 */
function getStatusText(status?: string): string {
  const map: Record<string, string> = {
    dispatched: '已分发',
    claimed: '已认领',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status || ''] || status || '未知'
}

/**
 * 获取状态标签
 */
function getStatusTag(status?: string): string {
  const map: Record<string, string> = {
    dispatched: 'info',
    claimed: 'primary',
    completed: 'success',
    cancelled: 'danger'
  }
  return map[status || ''] || ''
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

// 组件挂载时加载数据
onMounted(() => {
  loadTasks()
  loadStations()
})
</script>

<style scoped lang="scss">
.task-management {
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
  }

  .filter-container {
    margin-bottom: 20px;
    padding: 16px;
    background-color: #f5f7fa;
    border-radius: 4px;
  }

  .contact-info {
    .phone {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
    }
  }

  .text-gray {
    color: #909399;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
