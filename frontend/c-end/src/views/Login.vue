<template>
  <div class="login-container">
    <div class="login-card">
      <h1>霍营街道服务站</h1>

      <!-- 登录方式切换 -->
      <div class="login-tabs">
        <div
          class="tab-item"
          :class="{ active: loginType === 'code' }"
          @click="loginType = 'code'"
        >
          验证码登录
        </div>
        <div
          class="tab-item"
          :class="{ active: loginType === 'password' }"
          @click="loginType = 'password'"
        >
          密码登录
        </div>
      </div>

      <el-form :model="form" :rules="currentRules" ref="formRef" label-width="0">
        <el-form-item prop="phone">
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

        <!-- 验证码登录 -->
        <el-form-item prop="code" v-if="loginType === 'code'">
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

        <!-- 密码登录 -->
        <el-form-item prop="password" v-else>
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            @click="handleLogin"
            :loading="loading"
            size="large"
            style="width: 100%"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="tips">
        <p>首次使用？<router-link to="/home">返回首页选择服务</router-link></p>
      </div>
    </div>
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
const loginType = ref<'code' | 'password'>('password') // 默认密码登录

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
        role: 'c_end'
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
      const redirect = route.query.redirect as string || '/home'
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
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-card h1 {
  text-align: center;
  font-size: var(--font-size-title, 28px);
  font-weight: bold;
  margin-bottom: 24px;
  color: var(--text-color-primary, #303133);
}

.login-tabs {
  display: flex;
  margin-bottom: 24px;
  border-bottom: 1px solid #eee;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 12px 0;
  cursor: pointer;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-base, 16px);
  transition: color 0.3s, border-color 0.3s;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
}

.tab-item.active {
  color: var(--color-primary, #409EFF);
  border-bottom-color: var(--color-primary, #409EFF);
}

.tab-item:hover {
  color: var(--color-primary, #409EFF);
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
  font-size: var(--font-size-base, 16px);
}

.tips a {
  color: var(--color-primary, #409EFF);
  text-decoration: none;
}
</style>
