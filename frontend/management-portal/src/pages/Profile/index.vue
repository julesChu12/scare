<template>
  <div class="profile">
    <!-- 基本信息 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
          <el-button type="primary" @click="showEditDialog = true">编辑</el-button>
        </div>
      </template>
      <div class="profile-content">
        <el-avatar :size="80" :src="userInfo.avatar">
          <el-icon :size="40"><User /></el-icon>
        </el-avatar>
        <el-descriptions :column="2" border class="info-descriptions">
          <el-descriptions-item label="姓名">{{ userInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ maskPhone(userInfo.phone) }}</el-descriptions-item>
          <el-descriptions-item label="角色">{{ getRoleName(userInfo.role) }}</el-descriptions-item>
          <el-descriptions-item label="所属站点">{{ userInfo.station_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="入职时间">{{ formatDate(userInfo.created_at) }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>

    <!-- 账号安全 -->
    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <span>账号安全</span>
        </div>
      </template>
      <div class="security-content">
        <div class="security-item">
          <div class="security-info">
            <span class="label">登录密码</span>
            <span class="value">••••••••</span>
          </div>
          <el-button type="primary" link @click="showPasswordDialog = true">修改密码</el-button>
        </div>
        <div class="security-item">
          <div class="security-info">
            <span class="label">最近登录</span>
            <span class="value">{{ formatDateTime(userInfo.last_login_at) }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 工作统计（仅 staff 角色显示） -->
    <el-card v-if="userInfo.role === 'staff'" class="stats-card">
      <template #header>
        <div class="card-header">
          <span>工作统计</span>
        </div>
      </template>
      <el-row :gutter="20">
        <el-col :span="8">
          <div class="stat-item">
            <div class="stat-value">{{ workStats.total_tasks }}</div>
            <div class="stat-label">累计任务</div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="stat-item">
            <div class="stat-value">{{ workStats.avg_rating || '-' }}</div>
            <div class="stat-label">平均评分</div>
          </div>
        </el-col>
        <el-col :span="8">
          <div class="stat-item">
            <div class="stat-value">{{ workStats.completion_rate }}%</div>
            <div class="stat-label">完成率</div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 编辑信息弹窗 -->
    <el-dialog v-model="showEditDialog" title="编辑个人信息" width="500px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="姓名">
          <el-input v-model="editForm.name" placeholder="请输入姓名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveProfile">保存</el-button>
      </template>
    </el-dialog>

    <!-- 修改密码弹窗 -->
    <el-dialog v-model="showPasswordDialog" title="修改密码" width="500px">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
        <el-form-item label="当前密码" prop="old_password">
          <el-input v-model="passwordForm.old_password" type="password" show-password placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password placeholder="请输入新密码" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input v-model="passwordForm.confirm_password" type="password" show-password placeholder="请再次输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" :loading="changingPassword" @click="handleChangePassword">确认修改</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { User } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/store/modules/auth'
import dayjs from 'dayjs'

const authStore = useAuthStore()

// 用户信息
const userInfo = ref({
  name: '',
  phone: '',
  role: '',
  station_name: '',
  avatar: '',
  created_at: '',
  last_login_at: '',
})

// 工作统计
const workStats = ref({
  total_tasks: 0,
  avg_rating: 0,
  completion_rate: 0,
})

// 弹窗控制
const showEditDialog = ref(false)
const showPasswordDialog = ref(false)
const saving = ref(false)
const changingPassword = ref(false)

// 编辑表单
const editForm = reactive({
  name: '',
})

// 密码表单
const passwordFormRef = ref<FormInstance>()
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
})

// 密码验证规则
const passwordRules: FormRules = {
  old_password: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.new_password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

/**
 * 加载用户信息
 */
async function loadUserInfo() {
  try {
    // 从 store 获取基本信息
    const user = authStore.user as any
    if (user) {
      userInfo.value = {
        name: user.name || '',
        phone: user.phone || '',
        role: user.roles?.[0] || '',
        station_name: user.station_id ? `站点ID: ${user.station_id}` : '-', // 暂时使用ID
        avatar: '',
        created_at: user.created_at || '',
        last_login_at: user.last_login_at || '',
      }
      editForm.name = user.name || ''
    }
  } catch (error) {
    console.error('Failed to load user info:', error)
  }
}

/**
 * 加载工作统计（仅 staff）
 */
async function loadWorkStats() {
  if (userInfo.value.role !== 'staff') return
  try {
    // 这里可以调用统计 API
    workStats.value = {
      total_tasks: 0,
      avg_rating: 0,
      completion_rate: 0,
    }
  } catch (error) {
    console.error('Failed to load work stats:', error)
  }
}

/**
 * 保存个人信息
 */
async function handleSaveProfile() {
  try {
    saving.value = true
    // await userApi.updateProfile({ name: editForm.name })
    console.warn('updateProfile API not implemented')
    ElMessage.success('保存成功 (模拟)')
    showEditDialog.value = false
    // 刷新用户信息
    await authStore.fetchUserPermissions()
    await loadUserInfo()
  } catch (error) {
    console.error('Failed to save profile:', error)
  } finally {
    saving.value = false
  }
}

/**
 * 修改密码
 */
async function handleChangePassword() {
  if (!passwordFormRef.value) return
  await passwordFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      changingPassword.value = true
      // await userApi.changePassword({
      //   old_password: passwordForm.old_password,
      //   new_password: passwordForm.new_password,
      // })
      console.warn('changePassword API not implemented')
      ElMessage.success('密码修改成功 (模拟)')
      showPasswordDialog.value = false
      // 重置表单
      passwordForm.old_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
    } catch (error) {
      console.error('Failed to change password:', error)
    } finally {
      changingPassword.value = false
    }
  })
}

/**
 * 手机号脱敏
 */
function maskPhone(phone: string): string {
  if (!phone || phone.length < 7) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 获取角色名称
 */
function getRoleName(role: string): string {
  const roleMap: Record<string, string> = {
    admin: '系统管理员',
    station_manager: '社区工作人员',
    staff: '服务人员',
  }
  return roleMap[role] || role
}

/**
 * 格式化日期
 */
function formatDate(date?: string): string {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  loadUserInfo()
  loadWorkStats()
})
</script>

<style scoped lang="scss">
.profile {
  padding: 20px;

  .el-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .info-card {
    .profile-content {
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

  .security-card {
    .security-content {
      .security-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px 0;
        border-bottom: 1px solid #f0f0f0;

        &:last-child {
          border-bottom: none;
        }

        .security-info {
          .label {
            color: #909399;
            margin-right: 16px;
          }

          .value {
            color: #303133;
          }
        }
      }
    }
  }

  .stats-card {
    .stat-item {
      text-align: center;
      padding: 20px 0;

      .stat-value {
        font-size: 32px;
        font-weight: 600;
        color: #303133;
        line-height: 1.2;
      }

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-top: 8px;
      }
    }
  }
}
</style>
