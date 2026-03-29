<template>
  <div class="address-info-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>服务地址</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="详细地址" prop="address">
          <el-input
            v-model="form.address"
            type="textarea"
            placeholder="请输入详细地址，包括省市区街道门牌号等"
            :rows="4"
            size="large"
            clearable
          />
          <div class="field-tip">请填写准确的服务地址，以便工作人员上门服务</div>
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
  address: ''
})

const rules = {
  address: [
    { required: true, message: '请输入详细地址', trigger: 'blur' },
    { min: 5, message: '地址至少需要5个字符', trigger: 'blur' }
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
        address: form.address
      })

      userStore.updateProfile({
        address: form.address
      })

      ElMessage.success('服务地址已更新')
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
    form.address = userStore.profile.address || ''
  }
})
</script>

<style scoped>
.address-info-container {
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
