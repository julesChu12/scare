<template>
  <div class="statistics-overview">
    <!-- 页面头部 -->
    <div class="page-header">
      <h2>统计概览</h2>
      <el-select v-model="timeRange" placeholder="选择时间范围" style="width: 150px">
        <el-option label="最近7天" value="7" />
        <el-option label="最近30天" value="30" />
        <el-option label="最近90天" value="90" />
      </el-select>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon total-requests">
              <el-icon :size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overviewStats.total_requests }}</div>
              <div class="stat-label">总需求数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon pending">
              <el-icon :size="32"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overviewStats.pending }}</div>
              <div class="stat-label">待处理</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon completed">
              <el-icon :size="32"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overviewStats.completed }}</div>
              <div class="stat-label">已完成</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon in-progress">
              <el-icon :size="32"><Loading /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overviewStats.in_progress }}</div>
              <div class="stat-label">进行中</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>服务类型分布</span>
            </div>
          </template>
          <div class="chart-container service-distribution">
            <div v-if="serviceTypeStats.length" class="service-type-list">
              <div v-for="item in serviceTypeStats" :key="item.type" class="service-type-item">
                <div class="type-info">
                  <el-tag :type="getServiceTypeTag(item.type)" size="small">
                    {{ getServiceTypeText(item.type) }}
                  </el-tag>
                  <span class="type-count">{{ item.count }} 次</span>
                </div>
                <el-progress
                  :percentage="item.percentage"
                  :stroke-width="8"
                  :show-text="false"
                  :color="getProgressColor(item.type)"
                />
              </div>
            </div>
            <el-empty v-else description="暂无服务类型数据" :image-size="90" />
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>需求趋势</span>
            </div>
          </template>
          <div class="chart-container trend-chart">
            <v-chart v-if="trendData.length" class="trend-chart-canvas" :option="trendChartOption" autoresize />
            <el-empty v-else description="暂无趋势数据" :image-size="90" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 处理效率 -->
    <el-card class="efficiency-card">
      <template #header>
        <div class="card-header">
          <span>处理效率</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :xs="24" :md="8">
          <div class="efficiency-item">
            <div class="efficiency-value">{{ efficiencyStats.avg_response_time }} 分钟</div>
            <div class="efficiency-label">平均响应时间</div>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="efficiency-item">
            <div class="efficiency-value">{{ efficiencyStats.avg_process_time }} 分钟</div>
            <div class="efficiency-label">平均处理时间</div>
          </div>
        </el-col>
        <el-col :xs="24" :md="8">
          <div class="efficiency-item">
            <div class="efficiency-value">{{ efficiencyStats.satisfaction_rate }}%</div>
            <div class="efficiency-label">满意度</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 服务人员排行 -->
    <el-card class="ranking-card">
      <template #header>
        <div class="card-header">
          <span>服务人员排行</span>
        </div>
      </template>
      <div v-if="staffRanking.length" class="ranking-list">
        <div
          v-for="(staff, index) in staffRanking"
          :key="staff.id"
          :class="['ranking-item', `rank-${Math.min(index + 1, 4)}`]"
        >
          <div class="ranking-badge">
            <span class="ranking-number">{{ index + 1 }}</span>
            <span class="ranking-label">{{ getRankLabel(index) }}</span>
          </div>

          <div class="ranking-main">
            <div class="ranking-topline">
              <div class="ranking-person">
                <span class="person-name">{{ staff.name }}</span>
                <el-tag
                  :type="staff.is_online ? 'success' : 'info'"
                  size="small"
                  effect="plain"
                >
                  {{ staff.is_online ? '在线' : '离线' }}
                </el-tag>
              </div>
              <div class="ranking-score">
                <span class="score-value">{{ staff.completed_count }}</span>
                <span class="score-unit">单</span>
              </div>
            </div>

            <div class="ranking-progress">
              <div
                class="ranking-progress-fill"
                :style="{ width: `${getRankProgress(staff.completed_count)}%` }"
              />
            </div>

            <div class="ranking-metrics">
              <span class="metric-pill">完成量 {{ staff.completed_count }}</span>
              <span class="metric-pill rating-pill">
                <el-icon><StarFilled /></el-icon>
                {{ formatRating(staff.avg_rating) }}
              </span>
            </div>
          </div>
        </div>
      </div>
      <el-empty v-else description="暂无服务人员排行数据" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Clock, CircleCheck, Loading, StarFilled } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { statisticsApi } from '@/api'
import type {
  OverviewStatsData,
  ServiceTypeStatsData,
  TrendItemData,
  EfficiencyStatsData,
  StaffRankingItemData,
} from '@/types/api'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const timeRange = ref('7')

const overviewStats = ref<OverviewStatsData>({
  total_requests: 0,
  pending: 0,
  completed: 0,
  in_progress: 0,
})

const serviceTypeStats = ref<ServiceTypeStatsData[]>([])
const trendData = ref<TrendItemData[]>([])

const efficiencyStats = ref<EfficiencyStatsData>({
  avg_response_time: 0,
  avg_process_time: 0,
  satisfaction_rate: 0,
})

const staffRanking = ref<StaffRankingItemData[]>([])
const topCompletedCount = computed(() => {
  return Math.max(...staffRanking.value.map(item => item.completed_count), 0)
})

const trendChartOption = computed(() => {
  const axisData = trendData.value.map(item => formatTrendDate(item.date))
  const seriesData = trendData.value.map(item => item.count)

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      borderWidth: 0,
      textStyle: {
        color: '#f8fafc',
      },
      formatter: (params: Array<{ axisValue: string; data: number }>) => {
        const [{ axisValue, data }] = params
        return `${axisValue}<br/>新增需求：${data}`
      },
    },
    grid: {
      left: 12,
      right: 12,
      top: 28,
      bottom: 12,
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: axisData,
      axisTick: {
        show: false,
      },
      axisLine: {
        lineStyle: {
          color: '#dbe3ef',
        },
      },
      axisLabel: {
        color: '#6b7280',
        hideOverlap: true,
      },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: {
        lineStyle: {
          color: '#e8edf5',
          type: 'dashed',
        },
      },
      axisLabel: {
        color: '#94a3b8',
      },
    },
    series: [
      {
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: seriesData,
        lineStyle: {
          width: 3,
          color: '#409eff',
        },
        itemStyle: {
          color: '#409eff',
          borderColor: '#ffffff',
          borderWidth: 2,
        },
        areaStyle: {
          color: 'rgba(64, 158, 255, 0.16)',
        },
      },
    ],
  }
})

async function loadStatistics() {
  const days = parseInt(timeRange.value)
  try {
    const [overviewRes, serviceTypeRes, trendRes, efficiencyRes, rankingRes] = await Promise.all([
      statisticsApi.getOverviewStats({ days }),
      statisticsApi.getServiceTypeStats({ days }),
      statisticsApi.getRequestTrend({ days }),
      statisticsApi.getEfficiencyStats({ days }),
      statisticsApi.getStaffRanking({ days, limit: 10 }),
    ])

    if (overviewRes.msg === 'ok') {
      overviewStats.value = overviewRes.data
    }
    if (serviceTypeRes.msg === 'ok') {
      serviceTypeStats.value = serviceTypeRes.data || []
    }
    if (trendRes.msg === 'ok') {
      trendData.value = trendRes.data || []
    }
    if (efficiencyRes.msg === 'ok') {
      efficiencyStats.value = efficiencyRes.data
    }
    if (rankingRes.msg === 'ok') {
      staffRanking.value = rankingRes.data || []
    }
  } catch (error) {
    console.error('Failed to load statistics:', error)
    ElMessage.error('加载统计数据失败')
  }
}

const serviceTypeNameMap: Record<string, string> = {
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

function getServiceTypeText(type: string): string {
  return serviceTypeNameMap[type] || type
}

const serviceTypeTagMap: Record<string, string> = {
  meal: 'success',
  medical: 'danger',
  cleaning: 'warning',
  shopping: 'primary',
  care: 'info',
}

function getServiceTypeTag(type: string): string {
  return serviceTypeTagMap[type] || ''
}

const serviceTypeColorMap: Record<string, string> = {
  meal: '#67c23a',
  medical: '#f56c6c',
  cleaning: '#e6a23c',
  shopping: '#409eff',
  care: '#909399',
}

function getProgressColor(type: string): string {
  return serviceTypeColorMap[type] || '#409eff'
}

function formatTrendDate(date: string): string {
  const parsedDate = dayjs(date)
  if (parsedDate.isValid()) {
    return parsedDate.format('MM-DD')
  }

  return date.length >= 10 ? date.slice(5, 10) : date
}

function getRankLabel(index: number): string {
  const labels = ['TOP 1', 'TOP 2', 'TOP 3']
  return labels[index] || `NO. ${index + 1}`
}

function getRankProgress(completedCount: number): number {
  if (!topCompletedCount.value) return 0
  return Math.max((completedCount / topCompletedCount.value) * 100, 8)
}

function formatRating(rating: number): string {
  if (!rating) return '暂无评分'
  return `${rating.toFixed(1)} 分`
}

watch(timeRange, () => {
  loadStatistics()
})

onMounted(() => {
  loadStatistics()
})
</script>

<style scoped lang="scss">
.statistics-overview {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    flex-wrap: wrap;
    margin-bottom: 24px;

    h2 {
      margin: 0;
      font-size: 24px;
      font-weight: 500;
      color: #303133;
    }
  }

  .stats-row {
    margin-bottom: 20px;

    .el-col {
      margin-bottom: 12px;
    }
  }

  .stat-card {
    height: 100%;

    :deep(.el-card__body) {
      height: 100%;
    }

    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;
      min-width: 0;
    }

    .stat-icon {
      width: 64px;
      height: 64px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;

      &.total-requests {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.pending {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }

      &.completed {
        background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
      }

      &.in-progress {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      }
    }

    .stat-info {
      flex: 1;
      min-width: 0;

      .stat-value {
        font-size: 28px;
        font-weight: 600;
        color: #303133;
        line-height: 1.2;
      }

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-top: 4px;
        line-height: 1.5;
        word-break: break-word;
      }
    }
  }

  .charts-row {
    margin-bottom: 20px;

    .el-col {
      margin-bottom: 12px;
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .chart-card {
    height: 100%;
    border: 1px solid #ebeef5;

    :deep(.el-card__body) {
      height: calc(100% - 57px);
    }
  }

  .chart-container {
    padding: 12px 0;
    min-height: 280px;

    .service-type-list {
      max-height: 256px;
      padding-right: 4px;
      overflow-y: auto;

      .service-type-item {
        margin-bottom: 16px;

        &:last-child {
          margin-bottom: 0;
        }

        .type-info {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;

          .type-count {
            font-size: 14px;
            color: #606266;
          }
        }
      }
    }

    &.service-distribution {
      display: flex;
      align-items: center;
    }

    &.trend-chart {
      padding: 0;

      .trend-chart-canvas {
        width: 100%;
        height: 280px;
      }
    }
  }

  .efficiency-card {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }

    .efficiency-item {
      text-align: center;
      padding: 20px 0;

      .efficiency-value {
        font-size: 32px;
        font-weight: 600;
        color: #303133;
        line-height: 1.2;
      }

      .efficiency-label {
        font-size: 14px;
        color: #909399;
        margin-top: 8px;
      }
    }
  }

  @media (max-width: 768px) {
    padding: 16px;

    .page-header {
      align-items: flex-start;
      margin-bottom: 20px;

      h2 {
        font-size: 22px;
      }

      :deep(.el-select) {
        width: 100% !important;
        max-width: 180px;
      }
    }

    .stat-card {
      .stat-content {
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
      }

      .stat-icon {
        width: 56px;
        height: 56px;
      }

      .stat-info {
        width: 100%;

        .stat-value {
          font-size: 24px;
        }

        .stat-label {
          margin-top: 6px;
        }
      }
    }

    .chart-container {
      min-height: 240px;

      &.trend-chart {
        .trend-chart-canvas {
          height: 240px;
        }
      }
    }

    .efficiency-card {
      .efficiency-item {
        padding: 10px 0;

        .efficiency-value {
          font-size: 28px;
        }
      }
    }

    .ranking-card {
      .ranking-item {
        flex-direction: column;
      }

      .ranking-badge {
        width: 100%;
        min-width: 0;
        flex-direction: row;
        justify-content: space-between;
      }
    }
  }

  .ranking-card {
    .ranking-list {
      display: flex;
      flex-direction: column;
      gap: 14px;
    }

    .ranking-item {
      display: flex;
      align-items: stretch;
      gap: 16px;
      padding: 16px 18px;
      border-radius: 16px;
      border: 1px solid #ebeef5;
      background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
      transition: transform 0.2s ease, box-shadow 0.2s ease;

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 10px 24px rgba(15, 23, 42, 0.08);
      }

      &.rank-1 {
        border-color: #f7d26a;
        background: linear-gradient(135deg, #fff8de 0%, #fffef8 100%);

        .ranking-badge {
          background: linear-gradient(180deg, #f6c54f 0%, #d9a31a 100%);
        }

        .ranking-progress-fill {
          background: linear-gradient(90deg, #f6c54f 0%, #e19a1f 100%);
        }
      }

      &.rank-2 {
        .ranking-badge {
          background: linear-gradient(180deg, #b9c3d0 0%, #8995a6 100%);
        }

        .ranking-progress-fill {
          background: linear-gradient(90deg, #9aa8ba 0%, #74839a 100%);
        }
      }

      &.rank-3 {
        .ranking-badge {
          background: linear-gradient(180deg, #d59b72 0%, #b97445 100%);
        }

        .ranking-progress-fill {
          background: linear-gradient(90deg, #d59b72 0%, #c57c4d 100%);
        }
      }

      &.rank-4 {
        .ranking-badge {
          background: linear-gradient(180deg, #7c94b5 0%, #5f7494 100%);
        }

        .ranking-progress-fill {
          background: linear-gradient(90deg, #4f87ff 0%, #57b3ff 100%);
        }
      }
    }

    .ranking-badge {
      width: 76px;
      min-width: 76px;
      border-radius: 14px;
      color: #fff;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 10px 8px;
      box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.25);
    }

    .ranking-number {
      font-size: 28px;
      font-weight: 700;
      line-height: 1;
    }

    .ranking-label {
      margin-top: 6px;
      font-size: 11px;
      letter-spacing: 0.08em;
      opacity: 0.92;
    }

    .ranking-main {
      flex: 1;
      min-width: 0;
    }

    .ranking-topline {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      margin-bottom: 12px;
    }

    .ranking-person {
      display: flex;
      align-items: center;
      gap: 10px;
      min-width: 0;
    }

    .person-name {
      font-size: 17px;
      font-weight: 600;
      color: #1f2937;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .ranking-score {
      display: flex;
      align-items: baseline;
      gap: 4px;
      color: #111827;
      flex-shrink: 0;
    }

    .score-value {
      font-size: 24px;
      font-weight: 700;
      line-height: 1;
    }

    .score-unit {
      font-size: 13px;
      color: #6b7280;
    }

    .ranking-progress {
      width: 100%;
      height: 10px;
      border-radius: 999px;
      background: #edf2f7;
      overflow: hidden;
    }

    .ranking-progress-fill {
      height: 100%;
      border-radius: inherit;
      transition: width 0.3s ease;
    }

    .ranking-metrics {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
      margin-top: 12px;
    }

    .metric-pill {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 6px 10px;
      border-radius: 999px;
      background: #f3f6fb;
      color: #475569;
      font-size: 13px;
      line-height: 1;
    }

    .rating-pill {
      color: #b7791f;
      background: #fff7e6;
    }

    @media (max-width: 768px) {
      .ranking-item {
        padding: 14px;
        gap: 12px;
      }

      .ranking-badge {
        width: 64px;
        min-width: 64px;
      }

      .ranking-number {
        font-size: 24px;
      }

      .person-name {
        font-size: 15px;
      }

      .score-value {
        font-size: 20px;
      }

      .ranking-topline {
        align-items: flex-start;
      }
    }
  }
}
</style>
