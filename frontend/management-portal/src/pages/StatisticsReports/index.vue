<template>
  <div class="statistics-reports">
    <!-- 生成报表 -->
    <el-card class="generate-card">
      <template #header>
        <div class="card-header">
          <span>生成报表</span>
        </div>
      </template>
      <el-form :model="reportForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="报表类型">
              <el-select v-model="reportForm.type" placeholder="请选择报表类型" style="width: 100%">
                <el-option label="服务统计报表" value="service" />
                <el-option label="人员绩效报表" value="performance" />
                <el-option label="需求分析报表" value="request" />
                <el-option label="站点运营报表" value="station" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="站点筛选">
              <el-select v-model="reportForm.station_id" placeholder="全部站点" clearable style="width: 100%">
                <el-option
                  v-for="station in stationList"
                  :key="station.id"
                  :label="station.name"
                  :value="station.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="reportForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="导出格式">
              <el-radio-group v-model="reportForm.format">
                <el-radio value="xlsx">Excel (.xlsx)</el-radio>
                <el-radio value="pdf">PDF</el-radio>
                <el-radio value="csv">CSV</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item>
          <el-button @click="handlePreview">预览</el-button>
          <el-button type="primary" :loading="generating" @click="handleGenerate">
            生成并下载
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 历史报表 -->
    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <span>历史报表</span>
        </div>
      </template>
      <el-table v-loading="loading" :data="historyReports" stripe style="width: 100%">
        <el-table-column prop="name" label="报表名称" min-width="250">
          <template #default="{ row }">
            <div class="report-name">
              <el-icon class="file-icon"><Document /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="报表类型" width="150">
          <template #default="{ row }">
            <el-tag size="small">{{ getReportTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="format" label="格式" width="100">
          <template #default="{ row }">
            <span class="format-badge">{{ row.format.toUpperCase() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="生成时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="size" label="文件大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">
              <el-icon><Download /></el-icon>
              下载
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

    <!-- 预览弹窗 -->
    <el-dialog v-model="showPreview" title="报表预览" width="800px">
      <div class="preview-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="报表类型">{{ getReportTypeText(reportForm.type) }}</el-descriptions-item>
          <el-descriptions-item label="站点">{{ getStationName(reportForm.station_id) }}</el-descriptions-item>
          <el-descriptions-item label="时间范围">{{ formatDateRange(reportForm.dateRange) }}</el-descriptions-item>
          <el-descriptions-item label="导出格式">{{ reportForm.format.toUpperCase() }}</el-descriptions-item>
        </el-descriptions>
        <div class="preview-summary">
          <h4>预计包含数据</h4>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-value">{{ previewData.request_count }}</div>
                <div class="summary-label">需求记录</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-value">{{ previewData.task_count }}</div>
                <div class="summary-label">任务记录</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="summary-item">
                <div class="summary-value">{{ previewData.staff_count }}</div>
                <div class="summary-label">服务人员</div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
      <template #footer>
        <el-button @click="showPreview = false">取消</el-button>
        <el-button type="primary" :loading="generating" @click="handleGenerate">
          确认生成
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Download } from '@element-plus/icons-vue'
import dayjs from 'dayjs'

// 加载状态
const loading = ref(false)
const generating = ref(false)
const showPreview = ref(false)

// 报表表单
const reportForm = reactive({
  type: 'service',
  station_id: null as number | null,
  dateRange: [] as Date[],
  format: 'xlsx',
})

// 站点列表
const stationList = ref<Array<{ id: number; name: string }>>([])

// 历史报表
const historyReports = ref<Array<{
  id: number
  name: string
  type: string
  format: string
  created_at: string
  size: number
  url: string
}>>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 预览数据
const previewData = ref({
  request_count: 0,
  task_count: 0,
  staff_count: 0,
})

/**
 * 加载站点列表
 */
async function loadStations() {
  try {
    // 模拟数据
    stationList.value = [
      { id: 1, name: '霍营街道第一服务站' },
      { id: 2, name: '霍营街道第二服务站' },
      { id: 3, name: '回龙观服务站' },
    ]
  } catch (error) {
    console.error('Failed to load stations:', error)
  }
}

/**
 * 加载历史报表
 */
async function loadHistoryReports() {
  try {
    loading.value = true
    // 模拟数据
    historyReports.value = [
      {
        id: 1,
        name: '2026年1月服务统计报表.xlsx',
        type: 'service',
        format: 'xlsx',
        created_at: '2026-02-01 10:30:00',
        size: 125440,
        url: '#',
      },
      {
        id: 2,
        name: '2025年12月服务统计报表.xlsx',
        type: 'service',
        format: 'xlsx',
        created_at: '2026-01-01 09:15:00',
        size: 118784,
        url: '#',
      },
      {
        id: 3,
        name: 'Q4季度汇总报表.pdf',
        type: 'request',
        format: 'pdf',
        created_at: '2026-01-05 14:20:00',
        size: 256000,
        url: '#',
      },
      {
        id: 4,
        name: '2025年度人员绩效报表.xlsx',
        type: 'performance',
        format: 'xlsx',
        created_at: '2026-01-02 11:00:00',
        size: 89600,
        url: '#',
      },
    ]
    pagination.total = 4
  } catch (error) {
    console.error('Failed to load history reports:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 预览报表
 */
function handlePreview() {
  if (!reportForm.dateRange || reportForm.dateRange.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  // 模拟预览数据
  previewData.value = {
    request_count: 156,
    task_count: 142,
    staff_count: 12,
  }
  showPreview.value = true
}

/**
 * 生成报表
 */
async function handleGenerate() {
  if (!reportForm.dateRange || reportForm.dateRange.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return
  }
  try {
    generating.value = true
    // 模拟生成
    await new Promise(resolve => setTimeout(resolve, 1500))
    ElMessage.success('报表生成成功，开始下载')
    showPreview.value = false
    // 刷新历史列表
    await loadHistoryReports()
  } catch (error) {
    console.error('Failed to generate report:', error)
  } finally {
    generating.value = false
  }
}

/**
 * 下载报表
 */
function handleDownload(report: { url: string; name: string }) {
  ElMessage.success(`开始下载: ${report.name}`)
  // 实际下载逻辑
  // window.open(report.url)
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadHistoryReports()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadHistoryReports()
}

/**
 * 获取报表类型文本
 */
function getReportTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    service: '服务统计',
    performance: '人员绩效',
    request: '需求分析',
    station: '站点运营',
  }
  return typeMap[type] || type
}

/**
 * 获取站点名称
 */
function getStationName(stationId: number | null): string {
  if (!stationId) return '全部站点'
  const station = stationList.value.find(s => s.id === stationId)
  return station?.name || '-'
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime: string): string {
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

/**
 * 格式化日期范围
 */
function formatDateRange(dateRange: Date[]): string {
  if (!dateRange || dateRange.length !== 2) return '-'
  return `${dayjs(dateRange[0]).format('YYYY-MM-DD')} 至 ${dayjs(dateRange[1]).format('YYYY-MM-DD')}`
}

/**
 * 格式化文件大小
 */
function formatFileSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
}

onMounted(() => {
  loadStations()
  loadHistoryReports()
})
</script>

<style scoped lang="scss">
.statistics-reports {
  padding: 20px;

  .el-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .generate-card {
    .el-form-item {
      margin-bottom: 18px;
    }
  }

  .history-card {
    .report-name {
      display: flex;
      align-items: center;
      gap: 8px;

      .file-icon {
        color: #409eff;
        font-size: 18px;
      }
    }

    .format-badge {
      display: inline-block;
      padding: 2px 8px;
      background: #f0f2f5;
      border-radius: 4px;
      font-size: 12px;
      color: #606266;
    }

    .pagination-container {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .preview-content {
    .preview-summary {
      margin-top: 24px;

      h4 {
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 500;
        color: #303133;
      }

      .summary-item {
        text-align: center;
        padding: 16px;
        background: #f5f7fa;
        border-radius: 8px;

        .summary-value {
          font-size: 28px;
          font-weight: 600;
          color: #303133;
        }

        .summary-label {
          font-size: 14px;
          color: #909399;
          margin-top: 4px;
        }
      }
    }
  }
}
</style>
