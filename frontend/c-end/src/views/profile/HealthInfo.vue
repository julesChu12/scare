<template>
  <div class="health-info-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>健康档案</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="健康状况" prop="health_status">
          <el-select
            v-model="form.health_status"
            placeholder="请选择健康状况"
            size="large"
            style="width: 100%"
            clearable
          >
            <el-option label="良好" value="良好" />
            <el-option label="一般" value="一般" />
            <el-option label="较差" value="较差" />
          </el-select>
        </el-form-item>

        <el-form-item label="慢性病史（可选）" prop="medical_history">
          <el-input
            v-model="form.medical_history"
            type="textarea"
            placeholder="如有高血压、糖尿病等慢性疾病，请填写"
            :rows="3"
            size="large"
            clearable
          />
          <div class="field-tip">便于工作人员了解您的健康状况，提供更好的服务</div>
        </el-form-item>

        <el-form-item label="特殊需求（可选）" prop="special_needs">
          <el-input
            v-model="form.special_needs"
            type="textarea"
            placeholder="如需要轮椅、助听器等辅助设备，或其他特殊需求"
            :rows="3"
            size="large"
            clearable
          />
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
  health_status: '',
  medical_history: '',
  special_needs: ''
})

const rules = {
  health_status: [
    { required: true, message: '请选择健康状况', trigger: 'change' }
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
        health_status: form.health_status,
        medical_history: form.medical_history || undefined,
        special_needs: form.special_needs || undefined
      })

      userStore.updateProfile({
        health_status: form.health_status,
        medical_history: form.medical_history,
        special_needs: form.special_needs
      })

      ElMessage.success('健康档案已更新')
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
    form.health_status = userStore.profile.health_status || ''
    form.medical_history = userStore.profile.medical_history || ''
    form.special_needs = userStore.profile.special_needs || ''
  }
})
</script>

<style scoped>
.health-info-container {
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

:deep(.el-select .el-input__inner) {
  font-size: calc(18px * var(--font-scale));
}

:deep(.el-textarea__inner) {
  font-size: calc(16px * var(--font-scale));
  line-height: 1.6;
  padding: 12px;
}

.field-tip {
  font-size: calc(14px * var(--font-scale));
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
  font-size: calc(18px * var(--font-scale));
}
</style>
