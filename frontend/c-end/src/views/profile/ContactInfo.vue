<template>
  <div class="contact-info-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>联系信息</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="手机号">
          <el-input
            :model-value="userStore.user?.phone"
            disabled
            size="large"
          >
            <template #prefix>
              <el-icon><Phone /></el-icon>
            </template>
          </el-input>
          <div class="field-tip">手机号为登录账号，不可修改</div>
        </el-form-item>

        <el-divider content-position="left">
          <span style="font-size: 18px; font-weight: 500;">🚨 紧急联系人</span>
        </el-divider>

        <el-form-item label="联系人姓名" prop="emergency_name">
          <el-input
            v-model="form.emergency_name"
            placeholder="请输入紧急联系人姓名"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item label="联系人电话" prop="emergency_phone">
          <el-input
            v-model="form.emergency_phone"
            placeholder="请输入紧急联系人电话"
            maxlength="11"
            size="large"
            clearable
          >
            <template #prefix>
              <el-icon><Phone /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="与您的关系" prop="emergency_relation">
          <el-select
            v-model="form.emergency_relation"
            placeholder="请选择关系"
            size="large"
            style="width: 100%"
          >
            <el-option label="配偶" value="配偶" />
            <el-option label="子女" value="子女" />
            <el-option label="父母" value="父母" />
            <el-option label="兄弟姐妹" value="兄弟姐妹" />
            <el-option label="亲戚" value="亲戚" />
            <el-option label="朋友" value="朋友" />
            <el-option label="其他" value="其他" />
          </el-select>
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
import { ArrowLeft, Phone } from '@element-plus/icons-vue'
import { profileAPI } from '@/api'
import { useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  emergency_name: '',
  emergency_phone: '',
  emergency_relation: ''
})

const rules = {
  emergency_name: [
    { required: true, message: '请输入紧急联系人姓名', trigger: 'blur' }
  ],
  emergency_phone: [
    { required: true, message: '请输入紧急联系人电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  emergency_relation: [
    { required: true, message: '请选择关系', trigger: 'change' }
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
        emergency_contact: {
          name: form.emergency_name,
          phone: form.emergency_phone,
          relation: form.emergency_relation
        }
      })

      userStore.updateProfile({
        emergency_contact: {
          name: form.emergency_name,
          phone: form.emergency_phone,
          relation: form.emergency_relation
        }
      })

      ElMessage.success('联系信息已更新')
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
  if (userStore.profile?.emergency_contact) {
    const ec = userStore.profile.emergency_contact
    form.emergency_name = ec.name || ''
    form.emergency_phone = ec.phone || ''
    form.emergency_relation = ec.relation || ''
  }
})
</script>

<style scoped>
.contact-info-container {
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

:deep(.el-select .el-input__inner) {
  font-size: 18px;
}

.field-tip {
  font-size: 14px;
  color: var(--text-color-secondary, #909399);
  margin-top: 4px;
}

.el-divider {
  margin: 24px 0;
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
