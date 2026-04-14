<template>
  <div class="register-container">
    <div class="register-card">
      <div class="header">
        <el-button class="back-btn" circle @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <h1>注册账号</h1>
        <div class="placeholder"></div>
      </div>

      <el-form :model="form" :rules="rules" ref="formRef" label-position="top" class="register-form">
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            maxlength="11"
            size="large"
          >
            <template #prefix>
              <el-icon><Phone /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="验证码" prop="code">
          <div class="code-input-group">
            <el-input
              v-model="form.code"
              placeholder="请输入验证码"
              maxlength="6"
              size="large"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
            <el-button :disabled="countdown > 0" @click="sendCode" size="large">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请设置密码（至少6位）"
            size="large"
            show-password
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            size="large"
            show-password
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入您的姓名"
            size="large"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            @click="handleRegister"
            :loading="loading"
            size="large"
            style="width: 100%; margin-top: 8px;"
          >
            注册
          </el-button>
        </el-form-item>
      </el-form>

      <div class="tips">
        <p>已有账号？<router-link to="/login">立即登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Phone, Message, Lock, User } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { authAPI } from '@/api'
import { useTokenStore, useUserStore } from '@/store'

const router = useRouter()
const tokenStore = useTokenStore()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const countdown = ref(0)

const form = reactive({
  phone: '',
  code: '',
  password: '',
  confirmPassword: '',
  name: ''
})

// 表单验证规则
const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码为6位数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ]
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 发送验证码
const sendCode = async () => {
  if (!form.phone) {
    ElMessage.warning('请先输入手机号')
    return
  }

  if (!/^1[3-9]\d{9}$/.test(form.phone)) {
    ElMessage.warning('手机号格式不正确')
    return
  }

  try {
    await authAPI.sendCode({ phone: form.phone })
    ElMessage.success('验证码已发送')

    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    console.error('发送验证码失败:', error)
  }
}

// 注册
const handleRegister = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const result = await authAPI.register({
        phone: form.phone,
        code: form.code,
        password: form.password,
        name: form.name
      })

      // 保存 Token 和用户信息
      tokenStore.setToken(result.token)
      tokenStore.setRefreshToken(result.refresh_token)
      userStore.setUser(result.user)
      userStore.setProfile(result.profile)

      ElMessage.success('注册成功！')
      router.replace('/home')
    } catch (error: any) {
      if (error?.response?.data?.msg) {
        ElMessage.error(error.response.data.msg)
      } else {
        ElMessage.error('注册失败')
      }
      console.error('注册失败:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 20px;
}

.register-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-top: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  font-size: var(--font-size-title, 20px);
  font-weight: bold;
  margin: 0;
  color: var(--text-color-primary, #303133);
}

.back-btn {
  background: rgba(102, 126, 234, 0.1);
  border: none;
  color: #667eea;
}

.placeholder {
  width: 32px;
}

.register-form :deep(.el-form-item__label) {
  font-size: var(--font-size-sm, 14px);
  color: var(--text-color-primary, #303133);
  font-weight: 500;
  padding-bottom: 4px !important;
}

.register-form :deep(.el-input__inner) {
  font-size: var(--font-size-base, 16px);
}

.code-input-group {
  display: flex;
  gap: 10px;
  width: 100%;
}

.code-input-group .el-input {
  flex: 1;
}

.tips {
  text-align: center;
  margin-top: 20px;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
}

.tips a {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}
</style>
