<template>
  <div class="elderly-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
      <h2>老年人档案详情</h2>
    </div>

    <div v-loading="loading">
      <!-- 基本信息 -->
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <div class="basic-info">
          <el-avatar :size="100">
            <el-icon :size="50"><User /></el-icon>
          </el-avatar>
          <el-descriptions :column="2" border class="info-descriptions">
            <el-descriptions-item label="姓名">{{ elderlyInfo.name }}</el-descriptions-item>
            <el-descriptions-item label="性别">{{ elderlyInfo.gender === 'male' ? '男' : elderlyInfo.gender === 'female' ? '女' : '-' }}</el-descriptions-item>
            <el-descriptions-item label="年龄">{{ elderlyInfo.birth_date ? calculateAge(elderlyInfo.birth_date) + '岁' : '-' }}</el-descriptions-item>
            <el-descriptions-item label="出生日期">{{ formatDate(elderlyInfo.birth_date) }}</el-descriptions-item>
            <el-descriptions-item label="身份证号">{{ maskIdCard(elderlyInfo.id_card) }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ elderlyInfo.phone }}</el-descriptions-item>
            <el-descriptions-item label="地址" :span="2">{{ elderlyInfo.address || '-' }}</el-descriptions-item>
            <el-descriptions-item label="所属站点">{{ elderlyInfo.station_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="建档时间">{{ formatDate(elderlyInfo.created_at) }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </el-card>

      <!-- 健康信息 -->
      <el-card class="health-card">
        <template #header>
          <div class="card-header">
            <span>健康信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="健康状况">
            <el-tag :type="getHealthStatusType(elderlyInfo.health_status)" size="small">
              {{ getHealthStatusText(elderlyInfo.health_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="失能等级">{{ elderlyInfo.disability_level || '未评估' }}</el-descriptions-item>
          <el-descriptions-item label="病史" :span="2">{{ elderlyInfo.medical_history || '无' }}</el-descriptions-item>
          <el-descriptions-item label="特殊需求" :span="2">{{ elderlyInfo.special_needs || '无' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 服务记录 -->
      <el-card class="records-card">
        <template #header>
          <div class="card-header">
            <span>服务记录</span>
          </div>
        </template>
        <el-table v-loading="recordsLoading" :data="serviceRecords" stripe style="width: 100%" empty-text="暂无服务记录">
          <el-table-column prop="request_no" label="工单号" width="160" />
          <el-table-column label="日期" width="120">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="服务类型" width="120">
            <template #default="{ row }">
              <el-tag :type="getServiceTypeTag(row.service_type)" size="small">
                {{ getServiceTypeText(row.service_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="服务人员" width="100">
            <template #default="{ row }">
              {{ row.staff_name || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="评分" width="80">
            <template #default="{ row }">
              <span v-if="row.rating" class="rating">
                <el-icon><Star /></el-icon>
                {{ row.rating }}
              </span>
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="150" show-overflow-tooltip />
        </el-table>

        <!-- 服务记录分页 -->
        <div v-if="recordsPagination.total > 0" class="pagination-container">
          <el-pagination
            v-model:current-page="recordsPagination.page"
            v-model:page-size="recordsPagination.pageSize"
            :page-sizes="[5, 10, 20]"
            :total="recordsPagination.total"
            layout="total, sizes, prev, pager, next"
            small
            @size-change="handleRecordsSizeChange"
            @current-change="handleRecordsPageChange"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, User, Star } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { elderlyApi } from '@/api'
import type { ElderlyProfile, ElderlyServiceRecord } from '@/types/api'

const route = useRoute()
const router = useRouter()

// 加载状态
const loading = ref(false)
const recordsLoading = ref(false)

// 老年人ID
const elderlyId = ref<number>(Number(route.params.id))

// 老年人信息
const elderlyInfo = ref<ElderlyProfile>({
  id: 0,
  name: '',
  phone: '',
  gender: '',
  birth_date: '',
  id_card: '',
  address: '',
  station_id: 0,
  station_name: '',
  health_status: '',
  disability_level: '',
  medical_history: '',
  special_needs: '',
  customer_type: '',
  created_at: '',
})

// 服务记录
const serviceRecords = ref<ElderlyServiceRecord[]>([])

// 服务记录分页
const recordsPagination = reactive({
  page: 1,
  pageSize: 5,
  total: 0,
})

/**
 * 加载老年人详情
 */
async function loadElderlyDetail() {
  try {
    loading.value = true
    const res = await elderlyApi.getDetail(elderlyId.value)
    if (res.data) {
      elderlyInfo.value = res.data
    }
  } catch (error) {
    console.error('加载老人详情失败:', error)
    ElMessage.error('加载老人详情失败')
  } finally {
    loading.value = false
  }
}

/**
 * 加载服务记录
 */
async function loadServiceRecords() {
  try {
    recordsLoading.value = true
    const res = await elderlyApi.getServiceRecords(elderlyId.value, {
      page: recordsPagination.page,
      page_size: recordsPagination.pageSize,
    })
    serviceRecords.value = res.data?.items || []
    recordsPagination.total = res.data?.total || 0
  } catch (error) {
    console.error('加载服务记录失败:', error)
  } finally {
    recordsLoading.value = false
  }
}

/**
 * 返回列表
 */
function goBack() {
  router.push('/residents/elderly')
}

/**
 * 服务记录分页大小变化
 */
function handleRecordsSizeChange(size: number) {
  recordsPagination.pageSize = size
  recordsPagination.page = 1
  loadServiceRecords()
}

/**
 * 服务记录页码变化
 */
function handleRecordsPageChange(page: number) {
  recordsPagination.page = page
  loadServiceRecords()
}

/**
 * 计算年龄
 */
function calculateAge(birthDate: string): number {
  if (!birthDate) return 0
  return dayjs().diff(dayjs(birthDate), 'year')
}

/**
 * 身份证号脱敏
 */
function maskIdCard(idCard: string): string {
  if (!idCard || idCard.length < 15) return idCard || '-'
  return idCard.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2')
}

/**
 * 格式化日期
 */
function formatDate(date: string): string {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

/**
 * 获取健康状况标签类型
 */
function getHealthStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    good: 'success',
    normal: 'warning',
    poor: 'danger',
  }
  return typeMap[status] || ''
}

/**
 * 获取健康状况文本
 */
function getHealthStatusText(status: string): string {
  const textMap: Record<string, string> = {
    good: '良好',
    normal: '一般',
    poor: '较差',
  }
  return textMap[status] || status || '-'
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
 * 获取服务类型文本
 */
function getServiceTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    meal: '助餐服务',
    medical: '医疗服务',
    cleaning: '清洁服务',
    shopping: '代购服务',
    accompany: '陪护服务',
  }
  return typeMap[type] || type
}

/**
 * 获取状态标签类型
 */
function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    dispatched: 'info',
    processing: 'primary',
    completed: 'success',
    cancelled: 'info',
  }
  return typeMap[status] || ''
}

/**
 * 获取状态文本
 */
function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    pending: '待处理',
    dispatched: '已派单',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return textMap[status] || status
}

onMounted(() => {
  loadElderlyDetail()
  loadServiceRecords()
})
</script>

<style scoped lang="scss">
.elderly-detail {
  padding: 20px;

  .page-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      flex: 1;
      margin: 0 20px;
      font-size: 20px;
      font-weight: 500;
      color: #303133;
    }
  }

  .el-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .info-card {
    .basic-info {
      display: flex;
      align-items: flex-start;
      gap: 24px;

      .el-avatar {
        flex-shrink: 0;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      .info-descriptions {
        flex: 1;
      }
    }
  }

  .records-card {
    .rating {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      color: #e6a23c;
      font-weight: 500;

      .el-icon {
        font-size: 14px;
      }
    }

    .pagination-container {
      margin-top: 16px;
      display: flex;
      justify-content: flex-end;
    }
  }
}
</style>
