<template>
  <div class="request-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>需求管理</h3>
            <p class="subtitle">服务需求列表</p>
          </div>
          <div class="header-actions">
            <el-select
              v-model="filterStationId"
              placeholder="站点筛选"
              clearable
              filterable
              style="width: 160px; margin-right: 12px"
              @change="handleFilterChange"
            >
              <el-option
                v-for="station in stationList"
                :key="station.id"
                :label="station.name"
                :value="station.id"
              />
            </el-select>
            <el-select
              v-model="filterStatus"
              placeholder="状态筛选"
              clearable
              style="width: 140px; margin-right: 12px"
              @change="handleFilterChange"
            >
              <el-option label="全部" value="" />
              <el-option label="待处理" value="pending" />
              <el-option label="已派发" value="dispatched" />
              <el-option label="已认领" value="claimed" />
              <el-option label="处理中" value="processing" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
              <el-option label="已拒绝" value="rejected" />
            </el-select>
            <el-button
              type="primary"
              :icon="Plus"
              @click="handleCreate"
            >
              新建需求
            </el-button>
            <el-button
              :icon="Refresh"
              :loading="loading"
              @click="loadRequests"
            >
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 需求列表 -->
      <el-table
        v-loading="loading"
        :data="requestList"
        stripe
        style="width: 100%"
        empty-text="暂无服务需求"
      >
        <!-- 需求编号 -->
        <el-table-column label="需求编号" width="180">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.request_no }}</el-tag>
          </template>
        </el-table-column>

        <!-- 服务类型 -->
        <el-table-column label="服务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getServiceTypeTag(row.service_type)">
              {{ getServiceTypeText(row.service_type) }}
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
              <div>{{ row.contact_name || '-' }}</div>
              <div class="phone">{{ row.contact_phone || '-' }}</div>
            </div>
          </template>
        </el-table-column>

        <!-- 所属站点 -->
        <el-table-column label="所属站点" width="180">
          <template #default="{ row }">
            {{ row.station_name || '未分配' }}
          </template>
        </el-table-column>

        <!-- 服务地址 -->
        <el-table-column label="服务地址" min-width="200">
          <template #default="{ row }">
            {{ row.address || '-' }}
          </template>
        </el-table-column>

        <!-- 优先级 -->
        <el-table-column label="优先级" width="90">
          <template #default="{ row }">
            <el-tag :type="getPriorityTag(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 提交时间 -->
        <el-table-column label="提交时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="200" fixed="right">
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
              type="primary"
              size="small"
              link
              :icon="Edit"
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              link
              :icon="Delete"
              @click="handleDelete(row)"
            >
              删除
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
      title="需求详情"
      width="600px"
    >
      <el-descriptions v-if="currentRequest" :column="2" border>
        <el-descriptions-item label="需求编号" :span="2">
          {{ currentRequest.request_no }}
        </el-descriptions-item>
        <el-descriptions-item label="服务类型">
          <el-tag :type="getServiceTypeTag(currentRequest.service_type)">
            {{ getServiceTypeText(currentRequest.service_type) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTag(currentRequest.status)">
            {{ getStatusText(currentRequest.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="优先级">
          <el-tag :type="getPriorityTag(currentRequest.priority)" size="small">
            {{ getPriorityText(currentRequest.priority) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="站点ID">
          {{ currentRequest.station_id || '未分配' }}
        </el-descriptions-item>
        <el-descriptions-item label="联系人" :span="2">
          {{ currentRequest.contact_name }} - {{ currentRequest.contact_phone }}
        </el-descriptions-item>
        <el-descriptions-item label="服务地址" :span="2">
          {{ currentRequest.address }}
        </el-descriptions-item>
        <el-descriptions-item label="需求描述" :span="2">
          {{ currentRequest.description || '无' }}
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">
          {{ formatDateTime(currentRequest.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ formatDateTime(currentRequest.updated_at) }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 编辑/创建对话框 -->
    <el-dialog
      v-model="formVisible"
      :title="formTitle"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="服务类型" prop="service_type">
          <el-select v-model="formData.service_type" placeholder="请选择服务类型" style="width: 100%">
            <el-option label="送餐服务" value="meal" />
            <el-option label="就医陪护" value="medical" />
            <el-option label="日常照护" value="care" />
            <el-option label="居家维修" value="repair" />
            <el-option label="家政保洁" value="cleaning" />
            <el-option label="陪伴聊天" value="company" />
            <el-option label="紧急救助" value="emergency" />
            <el-option label="代买代购" value="shopping" />
            <el-option label="出行接送" value="transport" />
            <el-option label="康复理疗" value="rehab" />
            <el-option label="心理慰藉" value="psychology" />
            <el-option label="法律援助" value="legal_aid" />
            <el-option label="其他服务" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="formData.priority">
            <el-radio label="low">低</el-radio>
            <el-radio label="normal">普通</el-radio>
            <el-radio label="high">高</el-radio>
            <el-radio label="urgent">紧急</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="联系人" prop="contact_name">
              <el-input v-model="formData.contact_name" placeholder="姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="contact_phone">
              <el-input v-model="formData.contact_phone" placeholder="电话号码" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="服务地址" prop="address">
          <el-input v-model="formData.address" placeholder="详细地址" />
        </el-form-item>

        <el-form-item label="预约时间" prop="scheduled_at">
          <el-date-picker
            v-model="formData.scheduled_at"
            type="datetime"
            placeholder="选择预约时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="指派站点" prop="station_id">
          <el-select
            v-model="formData.station_id"
            placeholder="选择站点（可选）"
            clearable
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="station in stationList"
              :key="station.id"
              :label="station.name"
              :value="station.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="需求描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入详细需求描述"
          />
        </el-form-item>

        <el-form-item v-if="isEdit" label="状态" prop="status">
          <el-select v-model="formData.status" placeholder="选择状态" style="width: 100%">
            <el-option label="待处理" value="pending" />
            <el-option label="已派发" value="dispatched" />
            <el-option label="已认领" value="claimed" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" :loading="formLoading" @click="submitForm">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, nextTick } from 'vue'
import { Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { requestApi, stationApi } from '@/api'
import type { ServiceRequest, Station } from '@/types/api'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

import dayjs from 'dayjs'

// 加载状态
const loading = ref(false)

// 需求列表
const requestList = ref<ServiceRequest[]>([])

// 筛选状态
const filterStatus = ref('')
const filterStationId = ref<number | undefined>(undefined)

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 详情对话框
const detailDialogVisible = ref(false)
const currentRequest = ref<ServiceRequest | null>(null)

/**
 * 加载需求列表
 */
async function loadRequests() {
  try {
    loading.value = true

    const response = await requestApi.getRequests({
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filterStatus.value || undefined,
      station_id: filterStationId.value || undefined,
    } as any)

    const { items, total } = response.data
    pagination.total = total
    requestList.value = items
  } catch (error) {
    console.error('Failed to load requests:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 查看详情
 */
function handleViewDetail(request: ServiceRequest) {
  currentRequest.value = request
  detailDialogVisible.value = true
}

/**
 * 筛选变化
 */
function handleFilterChange() {
  pagination.page = 1
  loadRequests()
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadRequests()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadRequests()
}

/**
 * 获取服务类型文本
 */
function getServiceTypeText(type?: string): string {
  const typeMap: Record<string, string> = {
    meal: '送餐服务',
    medical: '就医陪护',
    care: '日常照护',
    repair: '居家维修',
    cleaning: '家政保洁',
    company: '陪伴聊天',
    emergency: '紧急救助',
    shopping: '代买代购',
    transport: '出行接送',
    rehab: '康复理疗',
    psychology: '心理慰藉',
    legal_aid: '法律援助',
    other: '其他服务',
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
    care: 'warning',
    repair: 'danger',
    cleaning: 'success',
    company: 'warning',
    emergency: 'danger',
    shopping: 'success',
    transport: 'primary',
    rehab: 'primary',
    psychology: 'warning',
    legal_aid: 'primary',
    other: 'info',
  }
  return tagMap[type || ''] || ''
}

/**
 * 获取状态文本
 */
function getStatusText(status?: string): string {
  const statusMap: Record<string, string> = {
    pending: '待处理',
    dispatched: '已派发',
    claimed: '已认领',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消',
    rejected: '已拒绝',
  }
  return statusMap[status || ''] || status || '未知'
}

/**
 * 获取状态标签颜色
 */
function getStatusTag(status?: string): string {
  const tagMap: Record<string, string> = {
    pending: 'warning',
    dispatched: 'primary',
    claimed: 'primary',
    processing: 'success',
    completed: 'success',
    cancelled: 'info',
    rejected: 'danger',
  }
  return tagMap[status || ''] || ''
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
    normal: '',
    high: 'warning',
    urgent: 'danger',
  }
  return tagMap[priority || ''] || ''
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}


// 站点列表
const stationList = ref<Station[]>([])

// 表单相关
const formVisible = ref(false)
const formLoading = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()
const formTitle = ref('')

const formData = reactive({
  id: undefined as number | undefined,
  service_type: '',
  contact_name: '',
  contact_phone: '',
  address: '',
  priority: 'normal',
  description: '',
  scheduled_at: '',
  station_id: undefined as number | undefined,
  status: 'pending'
})

const rules = reactive<FormRules>({
  service_type: [{ required: true, message: '请选择服务类型', trigger: 'change' }],
  contact_name: [{ required: true, message: '请输入联系人姓名', trigger: 'blur' }],
  contact_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  address: [{ required: true, message: '请输入服务地址', trigger: 'blur' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }]
})

/**
 * 加载站点列表
 */
async function loadStations() {
  try {
    const res = await stationApi.getStations({ page: 1, page_size: 100 })
    stationList.value = res.data.items
  } catch (error) {
    console.error('Failed to load stations:', error)
  }
}

/**
 * 新建需求
 */
function handleCreate() {
  isEdit.value = false
  formTitle.value = '新建服务需求'
  resetForm()
  formVisible.value = true
}

/**
 * 编辑需求
 */
function handleEdit(row: ServiceRequest) {
  isEdit.value = true
  formTitle.value = '编辑服务需求'
  
  formData.id = row.id
  formData.service_type = row.service_type
  formData.contact_name = row.contact_name
  formData.contact_phone = row.contact_phone
  formData.address = row.address
  formData.priority = row.priority
  formData.description = row.description
  formData.scheduled_at = row.scheduled_at || ''
  formData.station_id = row.station_id || undefined
  formData.status = row.status

  formVisible.value = true
}

/**
 * 重置表单
 */
function resetForm() {
  formData.id = undefined
  formData.service_type = ''
  formData.contact_name = ''
  formData.contact_phone = ''
  formData.address = ''
  formData.priority = 'normal'
  formData.description = ''
  formData.scheduled_at = ''
  formData.station_id = undefined
  formData.status = 'pending'
  
  nextTick(() => {
    formRef.value?.clearValidate()
  })
}

/**
 * 提交表单
 */
async function submitForm() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        formLoading.value = true
        const data = {
          service_type: formData.service_type,
          contact_name: formData.contact_name,
          contact_phone: formData.contact_phone,
          address: formData.address,
          priority: formData.priority,
          description: formData.description,
          scheduled_at: formData.scheduled_at || undefined,
          station_id: formData.station_id || undefined
        }

        if (isEdit.value && formData.id) {
           await requestApi.updateRequest(formData.id, {
             ...data,
             status: formData.status as any
           })
           ElMessage.success('更新成功')
        } else {
          await requestApi.createRequest(data)
          ElMessage.success('创建成功')
        }
        
        formVisible.value = false
        loadRequests()
      } catch (error) {
        console.error('Submit failed:', error)
        ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
      } finally {
        formLoading.value = false
      }
    }
  })
}

/**
 * 删除需求
 */
function handleDelete(row: ServiceRequest) {
  ElMessageBox.confirm(
    '确认删除该服务需求吗？此操作不可恢复。',
    '警告',
    {
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await requestApi.deleteRequest(row.id)
      ElMessage.success('删除成功')
      // 如果当前页只有一条数据且不是第一页，则跳转到上一页
      if (requestList.value.length === 1 && pagination.page > 1) {
        pagination.page--
      }
      loadRequests()
    } catch (error) {
      console.error('Delete failed:', error)
      ElMessage.error('删除失败')
    }
  })
}

// 组件挂载时加载数据
onMounted(() => {
  loadRequests()
  loadStations()
})
</script>

<style scoped lang="scss">
.request-management {
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
      align-items: center;
    }
  }

  .contact-info {
    .phone {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
    }
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
