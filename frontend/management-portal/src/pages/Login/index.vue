<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <!-- 标题 -->
      <div class="login-header">
        <h2>{{ appTitle }}</h2>
        <p class="subtitle">工作人员登录</p>
      </div>

      <!-- 登录表单 -->
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        size="large"
        @keyup.enter="handleLogin"
      >
        <!-- 手机号 -->
        <el-form-item prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            :prefix-icon="Phone"
            clearable
          />
        </el-form-item>

        <!-- 密码 -->
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <!-- 登录按钮 -->
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 测试账号提示 -->
      <div class="test-accounts">
        <el-divider>测试账号</el-divider>
        <div class="account-list">
          <div class="account-item" @click="fillTestAccount('13800000001')">
            <span class="label">管理员：</span>
            <span class="value">13800000001 / Test@123</span>
          </div>
          <div class="account-item" @click="fillTestAccount('13800000002')">
            <span class="label">站点管理员：</span>
            <span class="value">13800000002 / Test@123</span>
          </div>
          <div class="account-item" @click="fillTestAccount('13800000004')">
            <span class="label">工作人员：</span>
            <span class="value">13800000004 / Test@123</span>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Phone, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'

// 应用标题
const appTitle = import.meta.env.VITE_APP_TITLE || 'sCare 管理后台'

// 路由
const router = useRouter()
const route = useRoute()

// 状态管理
const authStore = useAuthStore()

// 表单引用
const formRef = ref<FormInstance>()

// 表单数据
const form = reactive({
  phone: '',
  password: '',
})

// 加载状态
const loading = ref(false)

// 表单验证规则
const rules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少6位', trigger: 'blur' },
  ],
}

/**
 * 处理登录
 */
async function handleLogin() {
  if (!formRef.value) return

  try {
    // 表单验证
    await formRef.value.validate()

    loading.value = true

    // 调用登录接口
    await authStore.login({
      phone: form.phone,
      password: form.password,
    })

    ElMessage.success('登录成功')

    // 跳转到重定向页面或首页
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error: any) {
    console.error('Login failed:', error)
    // Axios 拦截器已经处理了错误提示
  } finally {
    loading.value = false
  }
}

/**
 * 填充测试账号
 */
function fillTestAccount(phone: string) {
  form.phone = phone
  form.password = 'Test@123'
}
</script>

<style scoped lang="scss">
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

  .login-card {
    width: 450px;
    padding: 20px;

    :deep(.el-card__body) {
      padding: 40px 30px;
    }
  }

  .login-header {
    text-align: center;
    margin-bottom: 40px;

    h2 {
      font-size: 28px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 10px 0;
    }

    .subtitle {
      font-size: 14px;
      color: #909399;
      margin: 0;
    }
  }

  .test-accounts {
    margin-top: 30px;

    .el-divider {
      margin: 20px 0;
    }

    .account-list {
      .account-item {
        padding: 10px;
        border-radius: 4px;
        cursor: pointer;
        transition: all 0.3s;
        font-size: 13px;
        color: #606266;

        &:hover {
          background-color: #f5f7fa;
        }

        .label {
          font-weight: 500;
          color: #303133;
        }

        .value {
          color: #409eff;
          margin-left: 5px;
        }
      }
    }
  }
}
</style>
