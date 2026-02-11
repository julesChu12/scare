<template>
  <div class="dashboard">
    <!-- 欢迎信息 -->
    <div class="welcome-section">
      <h2>欢迎回来，{{ userName }}</h2>
      <p class="date">{{ currentDate }}</p>
    </div>

    <!-- 今日概览 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon new-requests">
              <el-icon :size="32"><Document /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ todayStats.new_requests }}</div>
              <div class="stat-label">今日新增需求</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon completed-tasks">
              <el-icon :size="32"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ todayStats.completed_tasks }}</div>
              <div class="stat-label">今日完成任务</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon pending-tasks">
              <el-icon :size="32"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ taskStats.dispatched }}</div>
              <div class="stat-label">待认领任务</div>
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
              <div class="stat-value">{{ taskStats.claimed }}</div>
              <div class="stat-label">进行中任务</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 我的任务统计（仅工作人员显示） -->
    <el-row v-if="myTaskStats" :gutter="20" class="my-stats-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>我的任务</span>
              <el-button type="primary" link @click="goToMyTasks">查看全部</el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <el-statistic title="待处理" :value="myTaskStats.claimed">
                <template #suffix>
                  <span class="stat-unit">个</span>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic title="已完成" :value="myTaskStats.completed">
                <template #suffix>
                  <span class="stat-unit">个</span>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic title="总计" :value="myTaskStats.total">
                <template #suffix>
                  <span class="stat-unit">个</span>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 任务和需求统计 -->
    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>任务状态分布</span>
            </div>
          </template>
          <div class="chart-container">
            <div class="status-list">
              <div class="status-item">
                <span class="status-dot dispatched"></span>
                <span class="status-name">待认领</span>
                <span class="status-value">{{ taskStats.dispatched }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot claimed"></span>
                <span class="status-name">进行中</span>
                <span class="status-value">{{ taskStats.claimed }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot completed"></span>
                <span class="status-name">已完成</span>
                <span class="status-value">{{ taskStats.completed }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot cancelled"></span>
                <span class="status-name">已取消</span>
                <span class="status-value">{{ taskStats.cancelled }}</span>
              </div>
            </div>
            <div class="total-info">
              <span class="total-label">总计</span>
              <span class="total-value">{{ taskStats.total }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>需求状态分布</span>
            </div>
          </template>
          <div class="chart-container">
            <div class="status-list">
              <div class="status-item">
                <span class="status-dot pending"></span>
                <span class="status-name">待处理</span>
                <span class="status-value">{{ requestStats.pending }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot dispatched"></span>
                <span class="status-name">已派发</span>
                <span class="status-value">{{ requestStats.dispatched }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot processing"></span>
                <span class="status-name">处理中</span>
                <span class="status-value">{{ requestStats.processing }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot completed"></span>
                <span class="status-name">已完成</span>
                <span class="status-value">{{ requestStats.completed }}</span>
              </div>
              <div class="status-item">
                <span class="status-dot cancelled"></span>
                <span class="status-name">已取消</span>
                <span class="status-value">{{ requestStats.cancelled }}</span>
              </div>
            </div>
            <div class="total-info">
              <span class="total-label">总计</span>
              <span class="total-value">{{ requestStats.total }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-row :gutter="20" class="quick-actions-row">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <el-button type="primary" :icon="List" @click="goToTaskPool">
              任务池
            </el-button>
            <el-button type="success" :icon="User" @click="goToMyTasks">
              我的任务
            </el-button>
            <el-button type="info" :icon="OfficeBuilding" @click="goToStations">
              站点管理
            </el-button>
            <el-button type="warning" :icon="UserFilled" @click="goToUsers">
              用户管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Document, CircleCheck, Clock, Loading, List, User, OfficeBuilding, UserFilled } from '@element-plus/icons-vue'
import { statisticsApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import type { TaskStatsData, RequestStatsData, TodayStatsData, MyTaskStatsData } from '@/types/api'
import dayjs from 'dayjs'

const router = useRouter()
const authStore = useAuthStore()

// 用户名
const userName = computed(() => authStore.user?.name || '用户')

// 当前日期
const currentDate = computed(() => dayjs().format('YYYY年MM月DD日 dddd'))

// 统计数据
const taskStats = ref<TaskStatsData>({
  total: 0,
  dispatched: 0,
  claimed: 0,
  completed: 0,
  cancelled: 0,
})

const requestStats = ref<RequestStatsData>({
  total: 0,
  pending: 0,
  dispatched: 0,
  processing: 0,
  completed: 0,
  cancelled: 0,
})

const todayStats = ref<TodayStatsData>({
  new_requests: 0,
  completed_tasks: 0,
  new_users: 0,
  avg_response_time: 0,
})

const myTaskStats = ref<MyTaskStatsData | null>(null)

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
  } catch (error) {
    console.error('Failed to load dashboard stats:', error)
  }
}

// 快捷导航
function goToTaskPool() {
  router.push('/services/tasks/pool')
}

function goToMyTasks() {
  router.push('/services/tasks/my')
}

function goToStations() {
  router.push('/stations/list')
}

function goToUsers() {
  router.push('/system/users')
}

onMounted(() => {
  loadDashboardStats()
})
</script>

<style scoped lang="scss">
.dashboard {
  padding: 20px;

  .welcome-section {
    margin-bottom: 24px;

    h2 {
      margin: 0 0 8px 0;
      font-size: 24px;
      font-weight: 500;
      color: #303133;
    }

    .date {
      margin: 0;
      font-size: 14px;
      color: #909399;
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

      &.new-requests {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.completed-tasks {
        background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
      }

      &.pending-tasks {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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

  .my-stats-row {
    margin-bottom: 20px;
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

    .status-list {
      .status-item {
        display: flex;
        align-items: center;
        padding: 10px 0;
        border-bottom: 1px solid #f0f0f0;

        &:last-child {
          border-bottom: none;
        }

        .status-dot {
          width: 10px;
          height: 10px;
          border-radius: 50%;
          margin-right: 12px;

          &.pending {
            background-color: #e6a23c;
          }

          &.dispatched {
            background-color: #409eff;
          }

          &.claimed,
          &.processing {
            background-color: #67c23a;
          }

          &.completed {
            background-color: #909399;
          }

          &.cancelled {
            background-color: #f56c6c;
          }
        }

        .status-name {
          flex: 1;
          font-size: 14px;
          color: #606266;
        }

        .status-value {
          font-size: 16px;
          font-weight: 500;
          color: #303133;
        }
      }
    }

    .total-info {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: 16px;
      padding-top: 16px;
      border-top: 2px solid #f0f0f0;

      .total-label {
        font-size: 14px;
        color: #909399;
      }

      .total-value {
        font-size: 24px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .quick-actions-row {
    .quick-actions {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
    }
  }

  .stat-unit {
    font-size: 14px;
    color: #909399;
    margin-left: 4px;
  }
}
</style>
