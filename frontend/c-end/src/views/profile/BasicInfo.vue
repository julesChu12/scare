<template>
  <div class="basic-info-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>基本信息</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="姓名" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入姓名"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender" size="large">
            <el-radio label="男" border>男</el-radio>
            <el-radio label="女" border>女</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="出生日期" prop="birth_date">
          <el-date-picker
            v-model="form.birth_date"
            type="date"
            placeholder="请选择出生日期"
            size="large"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>

        <el-form-item label="身份证号（可选）" prop="id_card">
          <el-input
            v-model="form.id_card"
            placeholder="请输入身份证号"
            maxlength="18"
            size="large"
            clearable
            show-password
          />
          <div class="field-tip">身份证号将加密保存，仅用于身份验证</div>
        </el-form-item>

        <div class="form-actions">
          <el-button @click="goBack" size="large">取消</el-button>
          <el-button type="primary" @click="handleSave" :loading="loading" size="large">
            保存
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { profileAPI } from '@/api'
import { useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  name: '',
  gender: '',
  birth_date: '',
  id_card: ''
})

const rules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  id_card: [
    {
      pattern: /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]$/,
      message: '身份证号格式不正确',
      trigger: 'blur'
    }
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
      await profileAPI.updateProfile({
        name: form.name,
        gender: form.gender || undefined,
        birth_date: form.birth_date || undefined,
        id_card: form.id_card || undefined
      })

      userStore.updateProfile({
        name: form.name,
        gender: form.gender,
        birth_date: form.birth_date,
        id_card: form.id_card
      })

      ElMessage.success('基本信息已更新')
      router.back()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败，请重试')
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  if (userStore.profile) {
    form.name = userStore.profile.name || ''
    form.gender = userStore.profile.gender || ''
    form.birth_date = userStore.profile.birth_date || ''
    form.id_card = userStore.profile.id_card || ''
  }
})
</script>

<style scoped>
.basic-info-container {
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
  font-size: 24px;
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

:deep(.el-form-item__label) {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-color-primary, #303133);
  margin-bottom: 8px;
}

:deep(.el-input__inner) {
  font-size: 18px;
  height: 48px;
}

:deep(.el-radio) {
  margin-right: 16px;
}

:deep(.el-radio.is-bordered) {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 18px;
  min-width: 100px;
  text-align: center;
}

:deep(.el-date-editor) {
  height: 48px;
}

:deep(.el-date-editor .el-input__inner) {
  font-size: 18px;
}

.field-tip {
  font-size: 14px;
  color: var(--text-color-secondary, #909399);
  margin-top: 4px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 32px;
}

.form-actions .el-button {
  flex: 1;
  height: 48px;
  font-size: 18px;
}
</style>
