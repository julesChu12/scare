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

    <div v-if="sourceStation" class="source-station-card">
      <div class="source-station-title">扫码来源站点</div>
      <div class="source-station-name">{{ sourceStation.name }}</div>
      <div class="source-station-meta">
        <span v-if="sourceStation.address">{{ sourceStation.address }}</span>
        <span v-if="sourceStation.phone">联系电话：{{ sourceStation.phone }}</span>
      </div>
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
        <div v-if="showAccountHint" class="phone-status-hint neutral">
          {{ accountHint }}
        </div>

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
            placeholder="请输入老人实际需要上门服务的详细地址"
            :rows="1"
            size="large"
          />
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
            :rows="2"
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
            style="width: 100%; margin-top: 8px;"
          >
            提交申请
          </el-button>
        </el-form-item>
      </el-form>

      <div class="tips" v-if="!isLoggedIn">
        <p>已有账号？<router-link :to="loginLink">使用密码登录</router-link></p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import type { FormInstance, UploadFile, UploadRequestOptions } from 'element-plus'
import { authAPI, geocodeAPI, requestsAPI, stationAPI, uploadAPI } from '@/api'
import type { Station } from '@/api'
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
const sourceStation = ref<Station | null>(null)
const sourceStationId = ref<number | null>(null)

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
const loginLink = computed(() => ({
  path: '/login',
  query: { redirect: route.fullPath }
}))

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

const isValidPhone = (phone: string) => /^1[3-9]\d{9}$/.test(phone)
const showAccountHint = computed(() => !isLoggedIn.value && isValidPhone(form.phone))
const accountHint = computed(() => '提交后系统会自动校验账号状态；如该手机号未注册，会自动创建账号并继续提交服务申请。')

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

const parseSourceStationId = () => {
  const raw = route.query.source_station_id
  const parsed = parseInt(Array.isArray(raw) ? (raw[0] ?? '') : (raw ?? ''), 10)
  sourceStationId.value = Number.isFinite(parsed) && parsed > 0 ? parsed : null
}

const loadSourceStation = async () => {
  parseSourceStationId()
  if (!sourceStationId.value) {
    sourceStation.value = null
    return
  }

  try {
    sourceStation.value = await stationAPI.getStationById(sourceStationId.value)
  } catch (error) {
    sourceStation.value = null
    sourceStationId.value = null
    ElMessage.warning('二维码来源站点无效，已按普通入口处理')
  }
}

const matchStation = async (lat: number, lng: number) => {
  return stationAPI.matchStation({ latitude: lat, longitude: lng })
}

const confirmDispatchIfNeeded = async (preview: { assignedStation: Station | null; manualReview: boolean }) => {
  if (preview.manualReview) {
    const sourceLabel = sourceStation.value ? `来源站点“${sourceStation.value.name}”将仅作为入口记录，` : ''
    await ElMessageBox.confirm(
      `当前无法准确解析服务地址，${sourceLabel}本次申请将进入人工复核，暂不自动分派站点。是否继续提交？`,
      '确认提交方式',
      {
        confirmButtonText: '继续提交',
        cancelButtonText: '返回修改',
        type: 'warning'
      }
    )
    return
  }

  if (preview.assignedStation && sourceStation.value && sourceStation.value.id !== preview.assignedStation.id) {
    ElMessage.info(`您从“${sourceStation.value.name}”入口进入，本次服务将按服务地址由“${preview.assignedStation.name}”受理`)
  }
}

const resolveDispatchPreview = async () => {
  let geocodeResult
  try {
    geocodeResult = await geocodeAPI.geocode({ address: form.address })
  } catch (_error) {
    return {
      assignedStation: null,
      resolvedAddress: form.address,
      manualReview: true
    }
  }

  const assignedStation = await matchStation(geocodeResult.latitude, geocodeResult.longitude)
  return {
    assignedStation,
    resolvedAddress: geocodeResult.formatted_address || form.address,
    manualReview: false
  }
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
      const preview = await resolveDispatchPreview()
      await confirmDispatchIfNeeded(preview)
      const requestPayload = {
        address: preview.resolvedAddress,
        submit_lat: coordinates.value?.lat,
        submit_lng: coordinates.value?.lng,
        source_station_id: sourceStationId.value || undefined,
        contact_name: form.name,
        contact_phone: form.phone,
        images: uploadedUrls.value.length > 0 ? uploadedUrls.value : undefined
      }

      if (isLoggedIn.value) {
        // 已登录用户：直接创建服务请求
        const result = await requestsAPI.createRequest({
          service_type: serviceType,
          ...requestPayload
        })

        ElMessage.success('服务申请已提交！')
        router.push(`/requests/${result.id}`)
      } else {
        // 未登录用户：调用快速开通API（自动注册+登录+创建请求）
        const result = await authAPI.quickStart({
          phone: form.phone,
          code: form.code,
          name: form.name,
          address: preview.resolvedAddress,
          submit_lat: coordinates.value?.lat,
          submit_lng: coordinates.value?.lng,
          source_station_id: sourceStationId.value || undefined,
          service_type: serviceType,
          description: form.description || undefined,
          images: uploadedUrls.value.length > 0 ? uploadedUrls.value : undefined,
          contact_name: form.name,
          contact_phone: form.phone
        })

        // 保存Token和用户信息
        tokenStore.setToken(result.token)
        tokenStore.setRefreshToken(result.refresh_token)
        userStore.setUser(result.user)
        userStore.setProfile(result.profile)

        ElMessage.success('服务申请已提交！')
        router.push(`/requests/${result.request.id}`)
      }
    } catch (error) {
      if (error === 'cancel' || error === 'close') {
        return
      }
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
    console.warn('获取位置失败:', error)
    coordinates.value = null
    hasLocation.value = false
    sessionStorage.removeItem('user_location')
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
  await loadSourceStation()

  // 获取地理位置
  await fetchLocation()

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
  height: 100vh;
  background: var(--bg-color, #f5f5f5);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 14px 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-shrink: 0;
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
  margin: 12px 16px;
  padding: 10px 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
}

.service-type-display .service-icon {
  font-size: calc(28px * var(--font-scale));
}

.service-type-display .service-name {
  font-size: var(--font-size-base, 16px);
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.phone-status-hint {
  margin: -6px 0 16px;
  padding: 10px 12px;
  border-radius: 10px;
  font-size: 13px;
  line-height: 1.5;
}

.phone-status-hint.neutral {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.source-station-card {
  background: #fff7ed;
  border: 1px solid #fed7aa;
  margin: 0 16px 12px;
  padding: 12px 14px;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex-shrink: 0;
}

.source-station-title {
  font-size: var(--font-size-sm, 14px);
  color: #9a3412;
}

.source-station-name {
  font-size: var(--font-size-base, 16px);
  font-weight: 600;
  color: #7c2d12;
}

.source-station-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #9a3412;
  font-size: var(--font-size-sm, 14px);
}

.form-container {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px 24px;
  display: flex;
  flex-direction: column;
}

.form-container :deep(.el-form) {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  font-size: var(--font-size-sm, 14px);
  color: var(--text-color-primary, #303133);
  font-weight: 500;
  padding-bottom: 4px !important;
  line-height: 1.2;
}

:deep(.el-input__inner) {
  font-size: var(--font-size-base, 16px);
}

:deep(.el-textarea__inner) {
  font-size: var(--font-size-base, 16px);
  resize: none;
}

.code-input-group {
  display: flex;
  gap: 8px;
  width: 100%;
}

.code-input-group .el-input {
  flex: 1;
}

.location-hint {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 4px;
  color: var(--color-success, #67C23A);
  font-size: var(--font-size-sm, 14px);
}

.tips {
  text-align: center;
  margin-top: auto;
  padding-top: 16px;
  padding-bottom: 8px;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
}

.tips a {
  color: var(--color-primary, #409EFF);
  text-decoration: none;
}

.submit-hint {
  text-align: center;
  margin: 8px 0 0;
  color: #909399;
  font-size: var(--font-size-sm, 12px);
}

.upload-hint {
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-sm, 14px);
  margin-top: 4px;
}

:deep(.el-upload--picture-card) {
  width: 60px;
  height: 60px;
  line-height: 60px;
}

:deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 60px;
  height: 60px;
}
</style>
