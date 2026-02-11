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
        <el-card>
          <template #header>
            <div class="card-header">
              <span>服务类型分布</span>
            </div>
          </template>
          <div class="chart-container">
            <div class="service-type-list">
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
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>需求趋势</span>
            </div>
          </template>
          <div class="chart-container trend-chart">
            <div v-for="(item, index) in trendData" :key="index" class="trend-item">
              <span class="trend-date">{{ item.date }}</span>
              <el-progress
                :percentage="item.percentage"
                :stroke-width="12"
                :format="() => item.count"
              />
            </div>
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
        <el-col :span="8">
          <div class="efficiency-item">
            <div class="efficiency-value">{{ efficiencyStats.avg_response_time }} 分钟</div>
            <div class="efficiency-label">平均响应时间</div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="efficiency-item">
            <div class="efficiency-value">{{ efficiencyStats.avg_process_time }} 分钟</div>
            <div class="efficiency-label">平均处理时间</div>
          </div>
        </el-col>
        <el-col :span="8">
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
      <el-table :data="staffRanking" stripe style="width: 100%">
        <el-table-column label="排名" width="80">
          <template #default="{ $index }">
            <el-tag v-if="$index < 3" :type="getRankType($index)" effect="dark" round>
              {{ $index + 1 }}
            </el-tag>
            <span v-else>{{ $index + 1 }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="completed_count" label="完成数" width="100" />
        <el-table-column prop="avg_rating" label="平均分" width="100">
          <template #default="{ row }">
            <span class="rating">{{ row.avg_rating || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_online ? 'success' : 'info'" size="small">
              {{ row.is_online ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

// 时间范围
const timeRange = ref('7')

// 概览统计
const overviewStats = ref({
  total_requests: 0,
  pending: 0,
  completed: 0,
  in_progress: 0,
})

// 服务类型统计
const serviceTypeStats = ref<Array<{
  type: string
  count: number
  percentage: number
}>>([])

// 趋势数据
const trendData = ref<Array<{
  date: string
  count: number
  percentage: number
}>>([])

// 效率统计
const efficiencyStats = ref({
  avg_response_time: 0,
  avg_process_time: 0,
  satisfaction_rate: 0,
})

// 服务人员排行
const staffRanking = ref<Array<{
  name: string
  completed_count: number
  avg_rating: number
  is_online: boolean
}>>([])

/**
 * 加载统计数据
 */
async function loadStatistics() {
  try {
    // 这里调用统计 API
    // const response = await statisticsApi.getOverview({ days: timeRange.value })
    // 模拟数据
    overviewStats.value = {
      total_requests: 256,
      pending: 42,
      completed: 198,
      in_progress: 16,
    }

    serviceTypeStats.value = [
      { type: 'meal', count: 89, percentage: 35 },
      { type: 'cleaning', count: 64, percentage: 25 },
      { type: 'medical', count: 51, percentage: 20 },
      { type: 'accompany', count: 32, percentage: 12 },
      { type: 'shopping', count: 20, percentage: 8 },
    ]

    trendData.value = [
      { date: '02-04', count: 12, percentage: 60 },
      { date: '02-03', count: 18, percentage: 90 },
      { date: '02-02', count: 15, percentage: 75 },
      { date: '02-01', count: 20, percentage: 100 },
      { date: '01-31', count: 14, percentage: 70 },
      { date: '01-30', count: 16, percentage: 80 },
      { date: '01-29', count: 11, percentage: 55 },
    ]

    efficiencyStats.value = {
      avg_response_time: 15,
      avg_process_time: 45,
      satisfaction_rate: 96,
    }

    staffRanking.value = [
      { name: '李师傅', completed_count: 45, avg_rating: 4.8, is_online: true },
      { name: '王师傅', completed_count: 38, avg_rating: 4.7, is_online: true },
      { name: '张师傅', completed_count: 32, avg_rating: 4.9, is_online: false },
      { name: '赵师傅', completed_count: 28, avg_rating: 4.6, is_online: true },
      { name: '刘师傅', completed_count: 25, avg_rating: 4.5, is_online: false },
    ]
  } catch (error) {
    console.error('Failed to load statistics:', error)
  }
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
 * 获取进度条颜色
 */
function getProgressColor(type: string): string {
  const colorMap: Record<string, string> = {
    meal: '#67c23a',
    medical: '#f56c6c',
    cleaning: '#e6a23c',
    shopping: '#409eff',
    accompany: '#909399',
  }
  return colorMap[type] || '#409eff'
}

/**
 * 获取排名标签类型
 */
function getRankType(index: number): string {
  const types = ['warning', 'info', 'success']
  return types[index] || ''
}

// 监听时间范围变化
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
    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;
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

  .chart-container {
    padding: 12px 0;

    .service-type-list {
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

    &.trend-chart {
      .trend-item {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 12px;

        &:last-child {
          margin-bottom: 0;
        }

        .trend-date {
          width: 50px;
          font-size: 13px;
          color: #909399;
        }

        .el-progress {
          flex: 1;
        }
      }
    }
  }

  .efficiency-card {
    margin-bottom: 20px;

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

  .ranking-card {
    .rating {
      color: #e6a23c;
      font-weight: 500;
    }
  }
}
</style>
