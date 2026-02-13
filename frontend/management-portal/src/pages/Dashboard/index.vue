<template>
  <div class="dashboard">
    <!-- 渐变欢迎横幅 -->
    <div class="welcome-banner">
      <div class="banner-left">
        <h2 class="greeting">{{ greeting }}，{{ userName }}</h2>
        <p class="date">{{ currentDate }}</p>
      </div>
      <div class="banner-right">
        <div class="banner-stat">
          <span class="banner-stat-value">{{ todayStats.new_requests }}</span>
          <span class="banner-stat-label">今日新增</span>
        </div>
        <div class="banner-divider"></div>
        <div class="banner-stat">
          <span class="banner-stat-value">{{ todayStats.completed_tasks }}</span>
          <span class="banner-stat-label">今日完成</span>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-content">
            <div class="stat-icon new-requests">
              <el-icon :size="28"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ todayStats.new_requests }}</div>
              <div class="stat-label">今日新增需求</div>
            </div>
          </div>
          <div class="stat-trend new-requests"></div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-content">
            <div class="stat-icon completed-tasks">
              <el-icon :size="28"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ todayStats.completed_tasks }}</div>
              <div class="stat-label">今日完成任务</div>
            </div>
          </div>
          <div class="stat-trend completed-tasks"></div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pending-tasks">
              <el-icon :size="28"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ taskStats.dispatched }}</div>
              <div class="stat-label">待认领任务</div>
            </div>
          </div>
          <div class="stat-trend pending-tasks"></div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-content">
            <div class="stat-icon in-progress">
              <el-icon :size="28"><Loading /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ taskStats.claimed }}</div>
              <div class="stat-label">进行中任务</div>
            </div>
          </div>
          <div class="stat-trend in-progress"></div>
        </div>
      </el-col>
    </el-row>

    <!-- 主体区域 -->
    <el-row :gutter="20" class="main-content">
      <!-- 左侧：待办任务列表 -->
      <el-col :xs="24" :lg="16">
        <div class="section-card task-list-card">
          <div class="section-header">
            <div class="header-title">
              <el-icon class="header-icon"><List /></el-icon>
              <span>最新待处理任务</span>
            </div>
            <el-button type="primary" link @click="goToTaskPool">查看全部</el-button>
          </div>
          <div class="task-card-list" v-if="recentTasks.length > 0">
            <div
              v-for="task in recentTasks"
              :key="task.id"
              class="task-item"
              @click="handleViewTask(task)"
            >
              <div class="task-item-left">
                <span class="service-tag" :class="task.request?.service_type">{{ getServiceTypeText(task.request?.service_type) }}</span>
                <span class="request-no">{{ task.request?.request_no || '-' }}</span>
              </div>
              <div class="task-item-right">
                <span class="time-text">{{ formatTimeAgo(task.created_at) }}</span>
                <el-icon class="arrow-icon"><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无待处理任务" :image-size="120" />
        </div>
      </el-col>

      <!-- 右侧：个人数据 -->
      <el-col :xs="24" :lg="8">
        <div class="section-card personal-card">
          <div class="section-header">
            <span>我的今日表现</span>
          </div>
          <div class="chart-wrapper">
            <v-chart class="chart" :option="pieOption" autoresize />
          </div>
          <div class="completion-rate" v-if="myTaskStats.total > 0">
            完成率 <strong>{{ completionRate }}%</strong>
          </div>
          <div class="personal-stats">
            <div class="p-stat-item">
              <div class="p-value claimed">{{ myTaskStats.claimed }}</div>
              <div class="p-label">进行中</div>
            </div>
            <div class="p-divider"></div>
            <div class="p-stat-item">
              <div class="p-value success">{{ myTaskStats.completed }}</div>
              <div class="p-label">已完成</div>
            </div>
            <div class="p-divider"></div>
            <div class="p-stat-item">
              <div class="p-value primary">{{ myTaskStats.total }}</div>
              <div class="p-label">累计参与</div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Document, CircleCheck, Clock, Loading, List, ArrowRight } from '@element-plus/icons-vue'
import { statisticsApi, taskApi, requestApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { PieChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'

use([CanvasRenderer, PieChart, TitleComponent, TooltipComponent, LegendComponent])
import type { TaskStatsData, RequestStatsData, TodayStatsData, MyTaskStatsData, Task, TaskListParams } from '@/types/api'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const authStore = useAuthStore()

// 用户名
const userName = computed(() => authStore.user?.name || '用户')

// 时段问候语
const greeting = computed(() => {
  const hour = dayjs().hour()
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

// 当前日期
const currentDate = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))

// 完成率
const completionRate = computed(() => {
  if (myTaskStats.value.total === 0) return 0
  return Math.round((myTaskStats.value.completed / myTaskStats.value.total) * 100)
})

// 统计数据
const taskStats = ref<TaskStatsData>({
  total: 0, dispatched: 0, claimed: 0, completed: 0, cancelled: 0,
})

const requestStats = ref<RequestStatsData>({
  total: 0, pending: 0, dispatched: 0, processing: 0, completed: 0, cancelled: 0,
})

const todayStats = ref<TodayStatsData>({
  new_requests: 0, completed_tasks: 0, new_users: 0, avg_response_time: 0,
})

const myTaskStats = ref<MyTaskStatsData>({
  claimed: 0, completed: 0, total: 0
})

// 图表配置
const pieOption = ref({
  tooltip: { trigger: 'item' },
  legend: { bottom: '0%', left: 'center' },
  series: [{
    name: '任务状态',
    type: 'pie',
    radius: ['45%', '70%'],
    avoidLabelOverlap: false,
    itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
    label: { show: false, position: 'center' },
    emphasis: { label: { show: true, fontSize: 18, fontWeight: 'bold' } },
    labelLine: { show: false },
    data: [
      { value: 0, name: '进行中', itemStyle: { color: '#409eff' } },
      { value: 0, name: '已完成', itemStyle: { color: '#67c23a' } },
    ],
  }],
})

// 最近任务
interface RecentTask extends Omit<Task, 'request'> {
  request?: { request_no: string; service_type: string }
}
const recentTasks = ref<RecentTask[]>([])

// 加载统计数据
async function loadDashboardStats() {
  try {
    const response = await statisticsApi.getDashboardStats()
    const data = response.data
    taskStats.value = data.task_stats
    requestStats.value = data.request_stats
    todayStats.value = data.today_stats
    if (data.my_task_stats) {
      myTaskStats.value = data.my_task_stats
    }
    pieOption.value.series[0].data = [
      { value: myTaskStats.value.claimed, name: '进行中', itemStyle: { color: '#409eff' } },
      { value: myTaskStats.value.completed, name: '已完成', itemStyle: { color: '#67c23a' } },
    ]
  } catch (error) {
    console.error('Failed to load dashboard stats:', error)
  }
}

// 加载最近任务
async function loadRecentTasks() {
  try {
    const response = await taskApi.getTaskPool({
      page: 1, page_size: 5, sort_by: 'created_at', sort_order: 'desc'
    } as TaskListParams & { sort_by: string; sort_order: string })
    const items = response.data.items || []
    const tasksWithRequests = await Promise.all(items.map(async (task: Task) => {
      if (task.request) return task
      if (task.request_id) {
        try {
          const reqRes = await requestApi.getRequest(task.request_id)
          return { ...task, request: reqRes.data }
        } catch { return task }
      }
      return task
    }))
    recentTasks.value = tasksWithRequests as RecentTask[]
  } catch (error) {
    console.error('Failed to load recent tasks:', error)
  }
}

// 快捷导航
function goToTaskPool() { router.push('/services/tasks') }
function handleViewTask(task: RecentTask) { router.push(`/services/tasks/${task.id}`) }

function getServiceTypeText(type?: string): string {
  const map: Record<string, string> = {
    meal: '送餐', medical: '医疗', cleaning: '清洁', shopping: '代购', accompany: '陪护',
  }
  return map[type || ''] || type || '服务'
}

function formatTimeAgo(time?: string): string {
  if (!time) return ''
  return dayjs(time).fromNow()
}

onMounted(() => {
  loadDashboardStats()
  loadRecentTasks()
})
</script>

<style scoped lang="scss">
.dashboard {
  padding: 20px;

  .welcome-banner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 28px 32px;
    margin-bottom: 24px;
    border-radius: 12px;
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
    color: #fff;

    .greeting {
      margin: 0 0 6px 0;
      font-size: 22px;
      font-weight: 600;
      color: #fff;
    }

    .date {
      margin: 0;
      font-size: 14px;
      color: rgba(255, 255, 255, 0.7);
    }

    .banner-right {
      display: flex;
      align-items: center;
      gap: 20px;
    }

    .banner-stat {
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    .banner-stat-value {
      font-size: 28px;
      font-weight: 700;
      line-height: 1.2;
    }

    .banner-stat-label {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.65);
      margin-top: 4px;
    }

    .banner-divider {
      width: 1px;
      height: 36px;
      background: rgba(255, 255, 255, 0.2);
    }
  }

  .stats-row {
    margin-bottom: 20px;

    .el-col {
      margin-bottom: 16px;
    }
  }

  .stat-card {
    background: #fff;
    border-radius: 12px;
    padding: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
    transition: all 0.3s ease;
    cursor: default;
    overflow: hidden;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    }

    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .stat-icon {
      width: 56px;
      height: 56px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      flex-shrink: 0;

      &.new-requests { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
      &.completed-tasks { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
      &.pending-tasks { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
      &.in-progress { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
    }

    .stat-info {
      flex: 1;

      .stat-value {
        font-size: 28px;
        font-weight: 700;
        color: #303133;
        line-height: 1.2;
        transition: all 0.6s ease;
      }

      .stat-label {
        font-size: 13px;
        color: #909399;
        margin-top: 4px;
      }
    }

    .stat-trend {
      height: 4px;
      border-radius: 2px;
      margin-top: 16px;
      opacity: 0.6;

      &.new-requests { background: linear-gradient(90deg, #667eea, #764ba2); }
      &.completed-tasks { background: linear-gradient(90deg, #11998e, #38ef7d); }
      &.pending-tasks { background: linear-gradient(90deg, #f093fb, #f5576c); }
      &.in-progress { background: linear-gradient(90deg, #4facfe, #00f2fe); }
    }
  }

  .main-content {
    .el-col { margin-bottom: 20px; }
  }

  .section-card {
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
    padding: 20px 24px;
    transition: box-shadow 0.3s ease;
    height: 100%;

    &:hover {
      box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
    }
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    font-size: 16px;
    font-weight: 500;
    color: #303133;

    .header-title {
      display: flex;
      align-items: center;
      gap: 8px;

      .header-icon { color: #409eff; }
    }
  }

  .task-list-card {
    min-height: 400px;

    .task-card-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }

    .task-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 14px 16px;
      border-radius: 8px;
      background: #f9fafb;
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover {
        background: #f0f5ff;
        transform: translateX(4px);
      }
    }

    .task-item-left {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .service-tag {
      display: inline-block;
      font-size: 12px;
      padding: 3px 10px;
      border-radius: 12px;
      background: #f0f2f5;
      color: #909399;
      font-weight: 500;

      &.meal { color: #67c23a; background: #f0f9eb; }
      &.medical { color: #f56c6c; background: #fef0f0; }
      &.cleaning { color: #e6a23c; background: #fdf6ec; }
      &.shopping { color: #409eff; background: #ecf5ff; }
      &.accompany { color: #9b59b6; background: #f5eef8; }
    }

    .request-no {
      font-family: monospace;
      font-size: 13px;
      color: #606266;
    }

    .task-item-right {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .time-text {
      color: #909399;
      font-size: 13px;
    }

    .arrow-icon {
      color: #c0c4cc;
      font-size: 14px;
    }
  }

  .personal-card {
    .chart-wrapper {
      height: 200px;
      width: 100%;
      margin-bottom: 8px;

      .chart { height: 100%; width: 100%; }
    }

    .completion-rate {
      text-align: center;
      font-size: 14px;
      color: #909399;
      margin-bottom: 12px;

      strong {
        color: #67c23a;
        font-size: 18px;
      }
    }

    .personal-stats {
      display: flex;
      justify-content: space-around;
      align-items: center;
      padding: 14px 0 4px;
      border-top: 1px solid #f0f0f0;

      .p-stat-item {
        text-align: center;

        .p-value {
          font-size: 22px;
          font-weight: 700;
          color: #303133;
          margin-bottom: 4px;

          &.claimed { color: #409eff; }
          &.success { color: #67c23a; }
          &.primary { color: #e6a23c; }
        }

        .p-label {
          font-size: 12px;
          color: #909399;
        }
      }

      .p-divider {
        width: 1px;
        height: 28px;
        background: #e4e7ed;
      }
    }
  }
}

// 响应式适配
@media (max-width: 768px) {
  .dashboard {
    .welcome-banner {
      flex-direction: column;
      align-items: flex-start;
      gap: 16px;
      padding: 20px 24px;

      .banner-right { align-self: flex-start; }
    }
  }
}
</style>

