<template>
  <div class="quick-start-container">
    <div class="header">
      <el-button class="back-btn" circle @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>申请服务</h1>
      <div class="placeholder"></div>
    </div>

    <div class="service-type-display" v-if="selectedServiceName">
      <span class="service-icon">{{ selectedServiceIcon }}</span>
      <span class="service-name">{{ selectedServiceName }}</span>
    </div>

    <div class="form-container">
      <el-form :model="form" :rules="rules" ref="formRef" label-position="top">
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
            maxlength="11"
            size="large"
            :disabled="isLoggedIn"
          />
        </el-form-item>

        <el-form-item label="验证码" prop="code" v-if="!isLoggedIn">
          <div class="code-input-group">
            <el-input v-model="form.code" placeholder="请输入验证码" maxlength="6" size="large" />
            <el-button :disabled="countdown > 0" @click="sendCode" size="large">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" size="large" />
        </el-form-item>

        <el-form-item label="详细地址" prop="address">
          <el-input
            v-model="form.address"
            type="textarea"
            placeholder="请输入详细地址"
            :rows="2"
            size="large"
          />
          <div class="location-hint" v-if="hasLocation">
            <el-icon><Location /></el-icon>
            <span>已获取您的位置</span>
          </div>
        </el-form-item>

        <el-form-item label="服务类型" prop="service_type" v-if="!selectedServiceType">
          <el-select v-model="form.service_type" placeholder="请选择服务类型" size="large" style="width: 100%">
            <el-option
              v-for="service in allServiceTypes"
              :key="service.key"
              :label="service.name"
              :value="service.key"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="服务描述">
          <el-input
            v-model="form.description"
            type="textarea"
            placeholder="请简单描述您的需求（可选）"
            :rows="3"
            size="large"
          />
        </el-form-item>

        <el-form-item label="上传照片">
          <el-upload
            v-model:file-list="fileList"
            :http-request="customUpload"
            :before-upload="beforeUpload"
            list-type="picture-card"
            accept="image/jpeg,image/png,image/webp"
            :limit="5"
            :on-exceed="handleExceed"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="upload-hint">最多 5 张，支持 JPG/PNG/WebP，单张不超过 5MB</div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            @click="handleSubmit"
            :loading="loading"
            size="large"
            style="width: 100%"
          >
            {{ isLoggedIn ? '提交申请' : '提交并开通服务' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="tips" v-if="!isLoggedIn">
        <p>已有账号？<router-link to="/login">使用密码登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Location, Plus } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile, UploadRequestOptions } from 'element-plus'
import { authAPI, requestsAPI, uploadAPI } from '@/api'
import { useTokenStore, useUserStore } from '@/store'
import { allServiceTypes, getServiceTypeName, getServiceTypeIcon } from '@/config/serviceTypes'
import { getCurrentPosition } from '@/utils/coordTransform'

const route = useRoute()
const router = useRouter()
const tokenStore = useTokenStore()
const userStore = useUserStore()

// 表单引用
const formRef = ref<FormInstance>()

// 状态
const countdown = ref(0)
const loading = ref(false)
const hasLocation = ref(false)
const coordinates = ref<{ lng: number; lat: number } | null>(null)

// 上传相关
const fileList = ref<UploadFile[]>([])
const uploadedUrls = ref<string[]>([])

const beforeUpload = (file: File) => {
  const isImage = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  if (!isImage) {
    ElMessage.error('仅支持 JPG/PNG/WebP 格式图片')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

const customUpload = async (options: UploadRequestOptions) => {
  try {
    const result = await uploadAPI.upload(options.file as File)
    uploadedUrls.value.push(result.url)
    options.onSuccess!(result)
  } catch (error) {
    options.onError!(error as any)
    ElMessage.error('图片上传失败')
  }
}

const handleExceed = () => {
  ElMessage.warning('最多上传 5 张图片')
}

// 是否已登录
const isLoggedIn = computed(() => tokenStore.isLoggedIn)

// 从 URL 获取的服务类型
const selectedServiceType = computed(() => route.query.type as string || '')
const selectedServiceName = computed(() => selectedServiceType.value ? getServiceTypeName(selectedServiceType.value) : '')
const selectedServiceIcon = computed(() => selectedServiceType.value ? getServiceTypeIcon(selectedServiceType.value) : '')

// 表单数据
const form = reactive({
  phone: '',
  code: '',
  name: '',
  address: '',
  service_type: '',
  description: ''
})

// 表单验证规则
const rules = computed(() => ({
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' }
  ],
  code: isLoggedIn.value ? [] : [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { pattern: /^\d{6}$/, message: '验证码为6位数字', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  address: [
    { required: true, message: '请输入详细地址', trigger: 'blur' }
  ],
  service_type: selectedServiceType.value ? [] : [
    { required: true, message: '请选择服务类型', trigger: 'change' }
  ]
}))

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

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const serviceType = selectedServiceType.value || form.service_type

      if (isLoggedIn.value) {
        // 已登录用户：直接创建服务请求
        const result = await requestsAPI.createRequest({
          service_type: serviceType,
          address: form.address,
          contact_name: form.name,
          contact_phone: form.phone,
          lat: coordinates.value?.lat,
          lng: coordinates.value?.lng,
          images: uploadedUrls.value.length > 0 ? uploadedUrls.value : undefined
        })

        ElMessage.success('服务申请已提交！')
        router.push(`/requests/${result.id}`)
      } else {
        // 未登录用户：调用快速开通API（自动注册+登录+创建请求）
        const result = await authAPI.quickStart({
          phone: form.phone,
          code: form.code,
          name: form.name,
          address: form.address,
          latitude: coordinates.value?.lat,
          longitude: coordinates.value?.lng,
          service_type: serviceType,
          description: form.description || undefined
        })

        // 保存Token和用户信息
        tokenStore.setToken(result.token)
        tokenStore.setRefreshToken(result.refresh_token)
        userStore.setUser(result.user)
        userStore.setProfile(result.profile)

        ElMessage.success('服务开通成功！')
        router.push(`/requests/${result.request.id}`)
      }
    } catch (error) {
      console.error('提交失败:', error)
    } finally {
      loading.value = false
    }
  })
}

// 获取地理位置
const fetchLocation = async () => {
  // 先尝试从 sessionStorage 获取
  const cached = sessionStorage.getItem('user_location')
  if (cached) {
    try {
      coordinates.value = JSON.parse(cached)
      hasLocation.value = true
      return
    } catch (e) {
      // ignore
    }
  }

  // 重新获取
  try {
    const coords = await getCurrentPosition()
    coordinates.value = coords
    hasLocation.value = true
    sessionStorage.setItem('user_location', JSON.stringify(coords))
  } catch (error) {
    console.warn('获取位置失败，使用默认位置:', error)
    // 回退到默认位置（霍营街道）
    const fallback = { lat: 40.0579, lng: 116.3698 }
    coordinates.value = fallback
    hasLocation.value = true
    sessionStorage.setItem('user_location', JSON.stringify(fallback))
  }
}

// 预填充已登录用户信息
const prefillUserInfo = () => {
  if (isLoggedIn.value) {
    if (userStore.user?.phone) {
      form.phone = userStore.user.phone
    }
    if (userStore.profile?.name) {
      form.name = userStore.profile.name
    }
    if (userStore.profile?.address) {
      form.address = userStore.profile.address
    }
  }
}

// 页面加载
onMounted(async () => {
  // 获取地理位置
  fetchLocation()

  // 如果已登录，预填充用户信息
  if (isLoggedIn.value) {
    // 先尝试获取最新用户信息
    if (!userStore.user || !userStore.profile) {
      try {
        const result = await authAPI.checkToken()
        userStore.setUser(result.user)
        if (result.profile) {
          userStore.setProfile(result.profile)
        }
      } catch (error) {
        console.error('获取用户信息失败:', error)
      }
    }
    prefillUserInfo()
  }
})
</script>

<style scoped>
.quick-start-container {
  min-height: 100vh;
  background: var(--bg-color, #f5f5f5);
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: var(--font-size-subtitle, 18px);
  font-weight: bold;
  margin: 0;
}

.back-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.placeholder {
  width: 32px;
}

.service-type-display {
  background: white;
  margin: 16px;
  padding: 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.service-type-display .service-icon {
  font-size: 32px;
}

.service-type-display .service-name {
  font-size: var(--font-size-subtitle, 18px);
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.form-container {
  padding: 0 16px 32px;
}

:deep(.el-form-item__label) {
  font-size: var(--font-size-base, 16px);
  color: var(--text-color-primary, #303133);
  font-weight: 500;
}

:deep(.el-input__inner) {
  font-size: var(--font-size-base, 16px);
}

:deep(.el-textarea__inner) {
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

.location-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  color: var(--color-success, #67C23A);
  font-size: var(--font-size-sm, 14px);
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

.upload-hint {
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
  margin-top: 8px;
}
</style>
