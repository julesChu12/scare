<template>
  <div class="login-container">
    <el-card class="login-card" shadow="hover">
      <div class="login-header">
        <!-- 返回按钮已隐藏 -->
        <!-- <div class="header-action">
          <el-button class="back-btn" @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
          </el-button>
        </div> -->
        <h2>sCare 社区养老服务</h2>
        <p class="subtitle">用户登录</p>
      </div>

      <div class="login-tabs">
        <button
          type="button"
          class="tab-chip"
          :class="{ active: loginType === 'code' }"
          @click="switchLoginType('code')"
        >
          验证码登录
        </button>
        <button
          type="button"
          class="tab-chip"
          :class="{ active: loginType === 'password' }"
          @click="switchLoginType('password')"
        >
          密码登录
        </button>
      </div>

      <el-form ref="formRef" :model="form" :rules="currentRules" size="large" label-width="0">
        <el-form-item prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            maxlength="11"
            :prefix-icon="Phone"
            clearable
          />
        </el-form-item>

        <el-form-item v-if="loginType === 'code'" prop="code">
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

        <el-form-item v-else prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
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
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="auth-footer">
        <p>首次使用？<router-link :to="registerLink">立即注册</router-link></p>
        <router-link v-if="loginType === 'password'" class="extra-link" :to="resetPasswordLink">忘记密码</router-link>
        <router-link class="extra-link" to="/quick">先去申请服务</router-link>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Phone, Message, Lock } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { authAPI } from '@/api'
import { useTokenStore, useUserStore } from '@/store'

const router = useRouter()
const route = useRoute()
const tokenStore = useTokenStore()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const countdown = ref(0)
const loginType = ref<'code' | 'password'>('code') // 默认验证码登录

const form = reactive({
  phone: '',
  code: '',
  password: ''
})

// 验证码登录规则
const codeRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

// 密码登录规则
const passwordRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 当前验证规则
const currentRules = computed(() => loginType.value === 'code' ? codeRules : passwordRules)

const registerLink = computed(() => {
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined
  return redirect ? { path: '/register', query: { redirect } } : { path: '/register' }
})

const resetPasswordLink = computed(() => {
  const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : undefined
  return redirect ? { path: '/reset-password', query: { redirect } } : { path: '/reset-password' }
})

const switchLoginType = (type: 'code' | 'password') => {
  if (loginType.value === type) return
  loginType.value = type
  formRef.value?.clearValidate()
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

const handleLogin = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const result = await authAPI.login({
        phone: form.phone,
        password: loginType.value === 'password' ? form.password : undefined,
        code: loginType.value === 'code' ? form.code : undefined,
        type: loginType.value
      })

      tokenStore.setToken(result.token)
      tokenStore.setRefreshToken(result.refresh_token)

      // Minimal user info from login response
      userStore.setUser({
        id: result.user_id,
        phone: result.phone,
        role: 'c_end',
        has_password: result.has_password
      })

      // Fetch profile for prefill (id_number/address/user_type)
      try {
        const check = await authAPI.checkToken()
        userStore.setUser(check.user)
        if (check.profile) {
          userStore.setProfile(check.profile)
        }
      } catch (e) {
        // ignore: profile may be missing
      }

      ElMessage.success('登录成功')

      // 跳转到原目标页面或首页
      const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/home'
      router.replace(redirect)
    } catch (error) {
      console.error('登录失败:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 100%;
  max-width: 460px;
}

.login-card :deep(.el-card__body) {
  padding: 28px 22px 22px;
}

.login-card :deep(.el-form-item) {
  margin-bottom: 20px;
}

.login-card :deep(.el-input__wrapper) {
  border-radius: 10px;
}

.login-card :deep(.el-button) {
  border-radius: 10px;
}

.header-action {
  position: absolute;
  left: 0;
  top: 2px;
}

.login-header {
  position: relative;
  text-align: center;
  margin-bottom: 24px;
}

.login-header h2 {
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

.login-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 18px;
}

.tab-chip {
  flex: 1;
  border: none;
  border-radius: 999px;
  padding: 12px 0;
  background: #f4f7fb;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-chip.active {
  background: rgba(64, 158, 255, 0.12);
  color: var(--color-primary, #409EFF);
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
  display: flex;
  flex-direction: column;
  gap: 10px;
  color: var(--text-color-secondary, #909399);
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
