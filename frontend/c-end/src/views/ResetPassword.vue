<template>
  <div class="reset-password-container">
    <el-card class="reset-password-card" shadow="hover">
      <div class="reset-password-header">
        <h2>sCare 社区养老服务</h2>
        <p class="subtitle">验证码重置密码</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" size="large" label-width="0">
        <el-form-item prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            maxlength="11"
            :prefix-icon="Phone"
            clearable
          />
        </el-form-item>

        <el-form-item prop="code">
          <div class="code-input-group">
            <el-input
              v-model="form.code"
              placeholder="请输入验证码"
              maxlength="6"
              :prefix-icon="Message"
              clearable
            />
            <el-button class="send-code-btn" :disabled="countdown > 0" plain @click="sendCode">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item prop="new_password">
          <el-input
            v-model="form.new_password"
            type="password"
            placeholder="请输入新密码，至少 6 位"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item prop="confirm_password">
          <el-input
            v-model="form.confirm_password"
            type="password"
            placeholder="请再次输入新密码"
            :prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            :loading="loading"
            style="width: 100%"
            @click="handleSubmit"
          >
            {{ loading ? '提交中...' : '重置密码' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="auth-footer">
        <router-link class="extra-link" :to="loginLink">返回登录</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { Lock, Message, Phone } from '@element-plus/icons-vue'
import { authAPI } from '@/api'

const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const countdown = ref(0)

const form = reactive({
  phone: '',
  code: '',
  new_password: '',
  confirm_password: ''
})

const validateConfirmPassword = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== form.new_password) {
    callback(new Error('两次输入的新密码不一致'))
    return
  }
  callback()
}

const rules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码为 6 位数字', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const loginLink = computed(() => {
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined
  return redirect ? { path: '/login', query: { redirect } } : { path: '/login' }
})

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

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await authAPI.resetPassword({
        phone: form.phone,
        code: form.code,
        new_password: form.new_password
      })

      ElMessage.success('密码已重置，请使用新密码登录')
      router.replace(loginLink.value)
    } catch (error) {
      console.error('重置密码失败:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.reset-password-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.reset-password-card {
  width: 100%;
  max-width: 460px;
}

.reset-password-card :deep(.el-card__body) {
  padding: 28px 22px 22px;
}

.reset-password-card :deep(.el-form-item) {
  margin-bottom: 20px;
}

.reset-password-card :deep(.el-input__wrapper) {
  border-radius: 10px;
}

.reset-password-card :deep(.el-button) {
  border-radius: 10px;
}

.reset-password-header {
  text-align: center;
  margin-bottom: 24px;
}

.reset-password-header h2 {
  margin: 0 0 8px;
  color: #303133;
  font-size: var(--font-size-title, 24px);
  font-weight: 600;
}

.subtitle {
  margin: 0;
  font-size: var(--font-size-sm, 14px);
  color: #909399;
}

.code-input-group {
  display: flex;
  gap: 8px;
  width: 100%;
}

.code-input-group .el-input {
  flex: 1;
}

.send-code-btn {
  min-width: 108px;
  padding: 0 10px;
}

.auth-footer {
  text-align: center;
  margin-top: 10px;
  font-size: var(--font-size-sm, 14px);
}

.auth-footer a {
  color: var(--color-primary, #409EFF);
  text-decoration: none;
  font-weight: 500;
}

.extra-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
</style>
