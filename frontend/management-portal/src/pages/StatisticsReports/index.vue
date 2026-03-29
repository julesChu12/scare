<template>
  <div class="statistics-reports">
    <el-card class="generate-card">
      <template #header>
        <div class="card-header">
          <div>
            <div class="card-title">生成报表</div>
            <div class="card-subtitle">按时间范围导出统计数据，生成后会自动下载并写入历史记录。</div>
          </div>
        </div>
      </template>

      <div v-if="!isAdmin" class="scope-banner">
        当前账号仅能生成所属站点报表：{{ currentScopeText }}
      </div>

      <el-form :model="reportForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :xs="24" :md="8">
            <el-form-item label="报表类型">
              <el-select v-model="reportForm.type" placeholder="请选择报表类型" style="width: 100%">
                <el-option v-for="option in reportTypeOptions" :key="option.value" :label="option.label" :value="option.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col v-if="isAdmin" :xs="24" :md="8">
            <el-form-item label="站点范围">
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
          <el-col :xs="24" :md="8">
            <el-form-item label="导出格式">
              <el-radio-group v-model="reportForm.format">
                <el-radio value="xlsx">Excel (.xlsx)</el-radio>
                <el-radio value="csv">CSV</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :xs="24" :md="16">
            <el-form-item label="时间范围">
              <el-date-picker
                v-model="reportForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                style="width: 100%"
                :shortcuts="dateShortcuts"
                :disabled-date="disableFutureDate"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="8">
            <el-form-item label="快捷范围">
              <div class="quick-range-actions">
                <el-button size="small" @click="setRecentRange(7)">近7天</el-button>
                <el-button size="small" @click="setRecentRange(30)">近30天</el-button>
                <el-button size="small" @click="setCurrentMonthRange">本月</el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <div class="generate-footer">
          <div class="scope-summary">
            <span>导出范围</span>
            <strong>{{ currentScopeText }}</strong>
            <span>{{ formatDateRange(reportForm.dateRange) }}</span>
          </div>
          <div class="actions">
            <el-button v-if="canGenerateReport" :loading="previewing" @click="handlePreview">预览</el-button>
            <el-button v-if="canGenerateReport" type="primary" :loading="generating" @click="handleGenerate">
              生成并下载
            </el-button>
          </div>
        </div>
      </el-form>
    </el-card>

    <div class="history-overview">
      <div v-for="item in historyOverview" :key="item.label" class="overview-card">
        <div class="overview-label">{{ item.label }}</div>
        <div class="overview-value">{{ item.value }}</div>
        <div class="overview-hint">{{ item.hint }}</div>
      </div>
    </div>

    <el-card class="history-card">
      <template #header>
        <div class="card-header">
          <div>
            <div class="card-title">历史报表</div>
            <div class="card-subtitle">只展示当前账号可见的报表记录。</div>
          </div>
          <el-button :icon="RefreshRight" @click="loadHistoryReports">刷新</el-button>
        </div>
      </template>

      <div class="history-toolbar">
        <el-select v-model="historyFilter.type" placeholder="全部类型" clearable style="width: 220px" @change="handleHistoryFilterChange">
          <el-option label="全部类型" value="" />
          <el-option v-for="option in reportTypeOptions" :key="option.value" :label="option.label" :value="option.value" />
        </el-select>
        <div class="history-toolbar__meta">
          当前筛选：{{ historyFilter.type ? getReportTypeText(historyFilter.type) : '全部类型' }}
        </div>
      </div>

      <el-table v-loading="loading" :data="historyReports" stripe style="width: 100%" empty-text="暂无报表记录">
        <el-table-column prop="name" label="报表名称" min-width="240">
          <template #default="{ row }">
            <div class="report-name">
              <el-icon class="file-icon"><Document /></el-icon>
              <div class="report-name__content">
                <span>{{ row.name }}</span>
                <small>{{ getScopeText(row.station_id) }}</small>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="130">
          <template #default="{ row }">
            <el-tag size="small" :type="getReportTagType(row.type)">{{ getReportTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="统计周期" min-width="220">
          <template #default="{ row }">
            {{ formatReportPeriod(row.start_date, row.end_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="format" label="格式" width="100">
          <template #default="{ row }">
            <span class="format-badge">{{ row.format.toUpperCase() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="生成时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="file_size" label="文件大小" width="110">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button v-if="canDownloadReport" type="primary" link @click="handleDownload(row)">
                <el-icon><Download /></el-icon>
                下载
              </el-button>
              <el-button
                v-if="canDeleteReport && canDeleteRow(row)"
                type="danger"
                link
                @click="handleDelete(row)"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

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

    <el-dialog v-model="showPreview" title="报表预览" width="820px">
      <div class="preview-content">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="报表类型">{{ getReportTypeText(reportForm.type) }}</el-descriptions-item>
          <el-descriptions-item label="站点范围">{{ currentScopeText }}</el-descriptions-item>
          <el-descriptions-item label="时间范围">{{ formatDateRange(reportForm.dateRange) }}</el-descriptions-item>
          <el-descriptions-item label="导出格式">{{ reportForm.format.toUpperCase() }}</el-descriptions-item>
        </el-descriptions>

        <div class="preview-summary">
          <div v-for="item in previewMetrics" :key="item.label" class="preview-card">
            <div class="preview-card__label">{{ item.label }}</div>
            <div class="preview-card__value">{{ item.value }}</div>
            <div class="preview-card__hint">{{ item.hint }}</div>
          </div>
        </div>

        <el-alert
          type="info"
          :closable="false"
          title="预览基于当前时间范围和站点范围统计，确认后会立即生成文件并下载到本地。"
        />
      </div>
      <template #footer>
        <el-button @click="showPreview = false">取消</el-button>
        <el-button type="primary" :loading="generating" @click="handleGenerate">确认生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Document, Download, RefreshRight } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { reportApi, stationApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import type { ReportData, ReportPreviewData } from '@/types/api'

type ReportType = 'service' | 'performance' | 'request' | 'station'
type ReportFormat = 'xlsx' | 'csv'

const authStore = useAuthStore()

const loading = ref(false)
const generating = ref(false)
const previewing = ref(false)
const showPreview = ref(false)

const reportTypeOptions: Array<{ label: string; value: ReportType }> = [
  { label: '服务统计报表', value: 'service' },
  { label: '人员绩效报表', value: 'performance' },
  { label: '需求分析报表', value: 'request' },
  { label: '站点运营报表', value: 'station' },
]

const reportForm = reactive({
  type: 'service' as ReportType,
  station_id: null as number | null,
  dateRange: [] as Date[],
  format: 'xlsx' as ReportFormat,
})

const historyFilter = reactive({
  type: '' as '' | ReportType,
})

const stationList = ref<Array<{ id: number; name: string }>>([])
const historyReports = ref<ReportData[]>([])
const previewData = ref<ReportPreviewData | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

const isAdmin = computed(() => authStore.hasRole('admin'))
const canGenerateReport = computed(() => authStore.hasPermission('data:report:generate'))
const canDownloadReport = computed(() => authStore.hasPermission('data:report:download'))
const canDeleteReport = computed(() => authStore.hasPermission('data:report:delete'))
const currentUserId = computed(() => authStore.user?.id ?? 0)
const currentUserStationId = computed(() => authStore.user?.station_id ?? null)

const currentScopeText = computed(() => {
  if (isAdmin.value) {
    return getScopeText(reportForm.station_id)
  }
  return getScopeText(currentUserStationId.value)
})

const historyOverview = computed(() => {
  const currentPageSize = historyReports.value.reduce((total, item) => total + (item.file_size || 0), 0)
  const latestReport = historyReports.value[0]

  return [
    {
      label: '可见报表数',
      value: `${pagination.total}`,
      hint: historyFilter.type ? getReportTypeText(historyFilter.type) : '全部类型',
    },
    {
      label: '当前页体积',
      value: formatFileSize(currentPageSize),
      hint: `${historyReports.value.length} 个文件`,
    },
    {
      label: '最近生成',
      value: latestReport ? formatDateTime(latestReport.created_at) : '-',
      hint: latestReport ? latestReport.name : '暂无历史记录',
    },
  ]
})

const previewMetrics = computed(() => {
  const preview = previewData.value
  if (!preview) return []

  switch (reportForm.type) {
    case 'service':
      return [
        {
          label: '需求总数',
          value: formatInteger(preview.request_count),
          hint: '纳入服务统计主表的需求记录',
        },
        {
          label: '完成率',
          value: formatPercent(preview.completion_rate),
          hint: `${formatInteger(preview.completed_request_count)} 条已完成`,
        },
        {
          label: '服务类型数',
          value: formatInteger(preview.service_type_count),
          hint: '服务类型分布工作表会按此维度展开',
        },
      ]
    case 'performance':
      return [
        {
          label: '上榜人员数',
          value: formatInteger(preview.ranked_staff_count),
          hint: '按完成任务数排序展示',
        },
        {
          label: '完成任务数',
          value: formatInteger(preview.completed_task_count),
          hint: '统计周期内已完成任务',
        },
        {
          label: '平均评分',
          value: formatRating(preview.avg_rating),
          hint: '按已完成任务加权计算',
        },
      ]
    case 'request':
      return [
        {
          label: '新增需求数',
          value: formatInteger(preview.request_count),
          hint: '需求趋势主表会覆盖该时间段',
        },
        {
          label: '趋势天数',
          value: formatInteger(preview.trend_days),
          hint: '按天输出新增需求趋势',
        },
        {
          label: '已完成需求',
          value: formatInteger(preview.completed_request_count),
          hint: '状态分布会同步展示',
        },
      ]
    case 'station':
      return [
        {
          label: '站点数量',
          value: formatInteger(preview.station_count),
          hint: '站点运营表中的统计范围',
        },
        {
          label: '需求总数',
          value: formatInteger(preview.request_count),
          hint: '站点维度汇总需求量',
        },
        {
          label: '整体完成率',
          value: formatPercent(preview.completion_rate),
          hint: `${formatInteger(preview.completed_request_count)} 条已完成`,
        },
      ]
    default:
      return []
  }
})

const dateShortcuts = [
  {
    text: '近7天',
    value: () => buildRecentRange(7),
  },
  {
    text: '近30天',
    value: () => buildRecentRange(30),
  },
  {
    text: '本月',
    value: () => [dayjs().startOf('month').toDate(), dayjs().endOf('day').toDate()],
  },
]

function buildRecentRange(days: number): [Date, Date] {
  const end = dayjs().endOf('day')
  const start = end.subtract(days - 1, 'day').startOf('day')
  return [start.toDate(), end.toDate()]
}

function setRecentRange(days: number) {
  reportForm.dateRange = buildRecentRange(days)
}

function setCurrentMonthRange() {
  reportForm.dateRange = [dayjs().startOf('month').toDate(), dayjs().endOf('day').toDate()]
}

function disableFutureDate(date: Date) {
  return dayjs(date).isAfter(dayjs().endOf('day'))
}

function getValidatedRange(): [Date, Date] | null {
  if (!reportForm.dateRange || reportForm.dateRange.length !== 2) {
    ElMessage.warning('请选择时间范围')
    return null
  }
  return [reportForm.dateRange[0], reportForm.dateRange[1]]
}

function buildReportPayload() {
  const range = getValidatedRange()
  if (!range) return null

  return {
    type: reportForm.type,
    format: reportForm.format,
    station_id: isAdmin.value ? reportForm.station_id : currentUserStationId.value,
    start_date: dayjs(range[0]).format('YYYY-MM-DD'),
    end_date: dayjs(range[1]).format('YYYY-MM-DD'),
  }
}

async function loadStations() {
  if (!authStore.hasPermission('station:list:view') && !isAdmin.value) {
    return
  }

  try {
    const res = await stationApi.getStations({ page: 1, page_size: 100 })
    if (res.msg === 'ok') {
      stationList.value = res.data.items.map((station) => ({ id: station.id, name: station.name }))
    }
  } catch (error) {
    console.error('Failed to load stations:', error)
  }
}

async function loadHistoryReports() {
  try {
    loading.value = true
    const res = await reportApi.getReports({
      page: pagination.page,
      page_size: pagination.pageSize,
      type: historyFilter.type || undefined,
    })
    if (res.msg === 'ok') {
      historyReports.value = res.data.items
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error('Failed to load history reports:', error)
    ElMessage.error('历史报表加载失败')
  } finally {
    loading.value = false
  }
}

async function handlePreview() {
  const payload = buildReportPayload()
  if (!payload) return

  try {
    previewing.value = true
    const res = await reportApi.previewReport(payload)
    if (res.msg === 'ok') {
      previewData.value = res.data
      showPreview.value = true
    }
  } catch (error) {
    console.error('Failed to preview report:', error)
    ElMessage.error('报表预览失败')
  } finally {
    previewing.value = false
  }
}

async function handleGenerate() {
  const payload = buildReportPayload()
  if (!payload) return

  try {
    generating.value = true
    const blob = await reportApi.generateReport(payload)

    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${getReportTypeText(reportForm.type)}_${payload.start_date.replace(/-/g, '')}_${payload.end_date.replace(/-/g, '')}.${reportForm.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('报表生成成功')
    showPreview.value = false
    await loadHistoryReports()
  } catch (error) {
    console.error('Failed to generate report:', error)
    ElMessage.error('报表生成失败')
  } finally {
    generating.value = false
  }
}

async function handleDownload(report: ReportData) {
  try {
    const blob = await reportApi.downloadReport(report.id)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${report.name}.${report.format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('下载成功')
  } catch (error) {
    console.error('Failed to download report:', error)
    ElMessage.error('下载失败')
  }
}

async function handleDelete(report: ReportData) {
  try {
    await ElMessageBox.confirm(`确定删除报表“${report.name}”吗？删除后无法恢复。`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })

    await reportApi.deleteReport(report.id)
    ElMessage.success('删除成功')

    if (historyReports.value.length === 1 && pagination.page > 1) {
      pagination.page -= 1
    }

    await loadHistoryReports()
  } catch (error) {
    if (error === 'cancel') return
    console.error('Failed to delete report:', error)
    ElMessage.error('删除失败')
  }
}

function canDeleteRow(report: ReportData) {
  return isAdmin.value || report.created_by === currentUserId.value
}

function handleHistoryFilterChange() {
  pagination.page = 1
  loadHistoryReports()
}

function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadHistoryReports()
}

function handlePageChange(page: number) {
  pagination.page = page
  loadHistoryReports()
}

function getReportTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    service: '服务统计',
    performance: '人员绩效',
    request: '需求分析',
    station: '站点运营',
  }
  return typeMap[type] || type
}

function getReportTagType(type: string): 'success' | 'warning' | 'info' | 'primary' {
  const tagMap: Record<string, 'success' | 'warning' | 'info' | 'primary'> = {
    service: 'primary',
    performance: 'success',
    request: 'warning',
    station: 'info',
  }
  return tagMap[type] || 'info'
}

function getStationName(stationId: number | null): string {
  if (!stationId) return '全部站点'
  const station = stationList.value.find((item) => item.id === stationId)
  if (station) return station.name
  if (stationId === currentUserStationId.value) return `当前站点 #${stationId}`
  return `站点 #${stationId}`
}

function getScopeText(stationId: number | null): string {
  return stationId ? getStationName(stationId) : '全部站点'
}

function formatDateTime(dateTime?: string | null): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

function formatDateRange(dateRange: Date[]): string {
  if (!dateRange || dateRange.length !== 2) return '-'
  return `${dayjs(dateRange[0]).format('YYYY-MM-DD')} 至 ${dayjs(dateRange[1]).format('YYYY-MM-DD')}`
}

function formatReportPeriod(startDate: string, endDate: string): string {
  return `${dayjs(startDate).format('YYYY-MM-DD')} 至 ${dayjs(endDate).format('YYYY-MM-DD')}`
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`
}

function formatInteger(value: number): string {
  return String(value || 0)
}

function formatRating(value: number): string {
  if (!value) return '暂无评分'
  return `${value.toFixed(1)} 分`
}

onMounted(() => {
  reportForm.dateRange = buildRecentRange(30)
  loadStations()
  loadHistoryReports()
})
</script>

<style scoped lang="scss">
.statistics-reports {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
  }

  .card-title {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }

  .card-subtitle {
    margin-top: 6px;
    font-size: 13px;
    color: #6b7280;
  }

  .scope-banner {
    margin-bottom: 18px;
    padding: 12px 14px;
    border-radius: 12px;
    background: linear-gradient(90deg, #eff6ff 0%, #f8fafc 100%);
    border: 1px solid #dbeafe;
    color: #1d4ed8;
    font-size: 13px;
  }

  .generate-card {
    .quick-range-actions {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .generate-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 16px;
      padding-top: 6px;
      border-top: 1px solid #f1f5f9;
    }

    .scope-summary {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 10px;
      color: #64748b;
      font-size: 13px;

      strong {
        color: #0f172a;
        font-weight: 600;
      }
    }

    .actions {
      display: flex;
      gap: 12px;
      flex-shrink: 0;
    }
  }

  .history-overview {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 16px;
  }

  .overview-card {
    padding: 18px 20px;
    border-radius: 16px;
    background: linear-gradient(145deg, #ffffff 0%, #f8fafc 100%);
    border: 1px solid #e5e7eb;
    box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
  }

  .overview-label {
    font-size: 13px;
    color: #64748b;
  }

  .overview-value {
    margin-top: 10px;
    font-size: 28px;
    font-weight: 700;
    color: #111827;
  }

  .overview-hint {
    margin-top: 8px;
    font-size: 12px;
    color: #94a3b8;
  }

  .history-card {
    .history-toolbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
    }

    .history-toolbar__meta {
      font-size: 13px;
      color: #64748b;
    }

    .report-name {
      display: flex;
      align-items: center;
      gap: 10px;

      .file-icon {
        color: #2563eb;
        font-size: 18px;
      }
    }

    .report-name__content {
      display: flex;
      flex-direction: column;
      gap: 4px;

      span {
        color: #111827;
      }

      small {
        color: #94a3b8;
      }
    }

    .format-badge {
      display: inline-block;
      min-width: 54px;
      text-align: center;
      padding: 4px 8px;
      background: #eef2ff;
      border-radius: 999px;
      font-size: 12px;
      color: #4338ca;
      font-weight: 600;
    }

    .table-actions {
      display: flex;
      gap: 12px;
      align-items: center;
    }

    .pagination-container {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .preview-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .preview-summary {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 14px;
  }

  .preview-card {
    padding: 18px;
    border-radius: 16px;
    background: linear-gradient(180deg, #f8fafc 0%, #ffffff 100%);
    border: 1px solid #e2e8f0;
  }

  .preview-card__label {
    font-size: 13px;
    color: #64748b;
  }

  .preview-card__value {
    margin-top: 10px;
    font-size: 26px;
    font-weight: 700;
    color: #0f172a;
  }

  .preview-card__hint {
    margin-top: 8px;
    font-size: 12px;
    line-height: 1.5;
    color: #94a3b8;
  }
}

@media (max-width: 992px) {
  .statistics-reports {
    .history-overview,
    .preview-summary {
      grid-template-columns: 1fr;
    }

    .generate-card .generate-footer,
    .history-card .history-toolbar,
    .card-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .history-card .pagination-container {
      justify-content: flex-start;
    }
  }
}
</style>
