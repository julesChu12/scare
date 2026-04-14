<template>
  <div class="password-setting-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>登录密码</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <div class="intro-card">
          <h2>设置手机号密码登录</h2>
          <p>首次设置密码无需填写当前密码；如果您已设置过密码，修改时必须验证当前密码。</p>
        </div>

        <el-form-item label="当前密码（首次设置可留空）" prop="current_password">
          <el-input
            v-model="form.current_password"
            type="password"
            placeholder="如已设置密码，请输入当前密码"
            size="large"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="新密码" prop="new_password">
          <el-input
            v-model="form.new_password"
            type="password"
            placeholder="请输入新密码，至少 6 位"
            size="large"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="确认新密码" prop="confirm_password">
          <el-input
            v-model="form.confirm_password"
            type="password"
            placeholder="请再次输入新密码"
            size="large"
            show-password
            clearable
          />
        </el-form-item>

        <div class="form-actions">
          <el-button @click="goBack" size="large">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="loading" size="large">
            保存密码
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { authAPI } from '@/api'
import { useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  current_password: '',
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
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const goBack = () => {
  router.back()
}

const handleSave = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      await authAPI.setPassword({
        current_password: form.current_password || undefined,
        new_password: form.new_password
      })

      if (userStore.user) {
        userStore.setUser({
          ...userStore.user,
          has_password: true
        })
      }

      ElMessage.success('登录密码已更新')
      router.back()
    } catch (error) {
      console.error('设置密码失败:', error)
    } finally {
      loading.value = false
    }
  })
}
</script>

<style scoped>
.password-setting-container {
  min-height: 100vh;
  background: var(--bg-color, #f5f5f5);
}

.header {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header h1 {
  font-size: calc(24px * var(--font-scale));
  font-weight: bold;
}

.content {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.el-form {
  background: white;
  padding: 24px;
  border-radius: 12px;
}

.intro-card {
  margin-bottom: 24px;
  padding: 16px;
  border-radius: 12px;
  background: #f5f9ff;
  color: #303133;
}

.intro-card h2 {
  margin: 0 0 8px;
  font-size: calc(20px * var(--font-scale));
}

.intro-card p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  font-size: calc(14px * var(--font-scale));
}

:deep(.el-form-item__label) {
  font-size: calc(18px * var(--font-scale));
  font-weight: 500;
  color: var(--text-color-primary, #303133);
  margin-bottom: 8px;
}

:deep(.el-input__inner) {
  font-size: calc(18px * var(--font-scale));
  height: 48px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.form-actions .el-button {
  flex: 1;
  height: 48px;
  font-size: calc(18px * var(--font-scale));
}
</style>
