<template>
  <div class="new-request-container">
    <!-- Header -->
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>申请服务</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <!-- 服务类型显示 -->
      <div class="service-type-display" v-if="selectedService">
        <span class="service-icon">{{ selectedService.icon }}</span>
        <div class="service-info">
          <span class="service-name">{{ selectedService.name }}</span>
          <span class="service-desc">{{ selectedService.description }}</span>
        </div>
        <el-button text @click="changeService">更换</el-button>
      </div>

      <!-- 表单 -->
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top" class="request-form">
        <el-form-item label="联系人姓名" prop="contact_name">
          <el-input
            v-model="form.contact_name"
            placeholder="请输入联系人姓名"
            size="large"
          />
        </el-form-item>

        <el-form-item label="联系电话" prop="contact_phone">
          <el-input
            v-model="form.contact_phone"
            placeholder="请输入联系电话"
            maxlength="11"
            size="large"
          />
        </el-form-item>

        <el-form-item label="服务地址" prop="address">
          <el-input
            v-model="form.address"
            type="textarea"
            placeholder="请输入详细地址"
            :rows="3"
          />
        </el-form-item>

        <el-form-item label="服务描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请简单描述您的需求（可选）"
            :rows="3"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            @click="handleSubmit"
            :loading="loading"
            class="submit-btn"
          >
            提交申请
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import { requestsAPI, geocodeAPI } from '@/api'
import { useUserStore } from '@/store'
import { allServiceTypes } from '@/config/serviceTypes'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const loading = ref(false)

// 获取服务类型
const serviceType = ref(route.query.type as string || 'other')

const selectedService = computed(() => {
  return allServiceTypes.find(s => s.key === serviceType.value)
})

// 表单数据
const form = reactive({
  contact_name: '',
  contact_phone: '',
  address: '',
  description: ''
})

// 表单验证规则
const rules = {
  contact_name: [
    { required: true, message: '请输入联系人姓名', trigger: 'blur' }
  ],
  contact_phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  address: [
    { required: true, message: '请输入服务地址', trigger: 'blur' }
  ]
}

const goBack = () => {
  router.back()
}

const changeService = () => {
  router.push('/services')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      // 尝试解析地址获取经纬度
      let lat: number | undefined
      let lng: number | undefined

      try {
        const geoResult = await geocodeAPI.geocode({ address: form.address })
        lat = geoResult.latitude
        lng = geoResult.longitude
      } catch (error) {
        console.warn('地址解析失败，将使用用户输入的地址:', error)
      }

      // 创建服务请求
      const result = await requestsAPI.createRequest({
        service_type: serviceType.value,
        contact_name: form.contact_name,
        contact_phone: form.contact_phone,
        address: form.address,
        description: form.description,
        latitude: lat,
        longitude: lng
      })

      ElMessage.success('服务申请已提交')
      router.push(`/requests/${result.id}`)
    } catch (error) {
      console.error('提交失败:', error)
    } finally {
      loading.value = false
    }
  })
}

// 预填充用户信息
onMounted(() => {
  if (userStore.profile) {
    form.contact_name = userStore.profile.name || ''
    form.address = userStore.profile.address || ''
  }
  if (userStore.user) {
    form.contact_phone = userStore.user.phone || ''
  }
})
</script>

<style scoped>
.new-request-container {
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
  font-size: var(--font-size-subtitle, 18px);
  font-weight: bold;
}

.content {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.service-type-display {
  background: white;
  border-radius: 12px;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.service-icon {
  font-size: 40px;
}

.service-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.service-name {
  font-size: var(--font-size-base, 16px);
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.service-desc {
  font-size: var(--font-size-sm, 14px);
  color: var(--text-color-secondary, #909399);
}

.request-form {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.submit-btn {
  width: 100%;
  font-size: var(--font-size-base, 16px);
}
</style>
