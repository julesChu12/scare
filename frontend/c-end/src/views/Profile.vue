<template>
  <div class="profile-container">
    <div class="header">
      <h1>我的资料</h1>
    </div>

    <div class="content">
      <!-- 用户卡片 -->
      <div class="user-card info-card">
        <div class="avatar">
          <el-icon><User /></el-icon>
        </div>
        <div class="user-info">
          <h2>{{ userStore.profile?.name || '用户' }}</h2>
          <p>{{ maskPhone(userStore.user?.phone || '') }}</p>
        </div>
      </div>

      <!-- 站点信息卡片 -->
      <div v-if="currentStation" class="station-card info-card">
        <div class="card-header">
          <span class="card-icon">🏢</span>
          <span class="card-title">所属服务站点</span>
        </div>
        <div class="station-info">
          <div class="info-row">
            <span class="label">站点名称</span>
            <span class="value">{{ currentStation.name }}</span>
          </div>
          <div class="info-row">
            <span class="label">📍 地址</span>
            <span class="value">{{ currentStation.address }}</span>
          </div>
          <div class="info-row">
            <span class="label">📞 电话</span>
            <span class="value">{{ currentStation.phone }}</span>
          </div>
          <div v-if="currentStation.work_hours" class="info-row">
            <span class="label">🕐 工作时间</span>
            <span class="value">{{ currentStation.work_hours }}</span>
          </div>
        </div>
        <a :href="`tel:${currentStation.phone}`" class="call-button">
          <el-icon><Phone /></el-icon>
          拨打站点电话
        </a>
      </div>

      <!-- 个人紧急联系人卡片 -->
      <div v-if="userStore.profile?.emergency_contact" class="emergency-card info-card">
        <div class="card-header">
          <span class="card-icon">🚨</span>
          <span class="card-title">紧急联系人</span>
        </div>
        <div class="emergency-info">
          <div class="contact-main">
            <span class="contact-name">{{ userStore.profile.emergency_contact.name }}</span>
            <span class="contact-relation">（{{ userStore.profile.emergency_contact.relation }}）</span>
          </div>
          <div class="contact-phone">{{ userStore.profile.emergency_contact.phone }}</div>
        </div>
        <a :href="`tel:${userStore.profile.emergency_contact.phone}`" class="call-button emergency">
          <el-icon><Phone /></el-icon>
          拨打
        </a>
      </div>

      <!-- 信息模块卡片 -->
      <div class="module-card info-card" @click="goToBasicInfo">
        <div class="module-header">
          <span class="module-title">基本信息</span>
          <el-icon class="arrow"><ArrowRight /></el-icon>
        </div>
        <div class="module-summary">
          {{ getBasicInfoSummary() }}
        </div>
      </div>

      <div class="module-card info-card" @click="goToContactInfo">
        <div class="module-header">
          <span class="module-title">联系信息</span>
          <el-icon class="arrow"><ArrowRight /></el-icon>
        </div>
        <div class="module-summary">
          {{ getContactInfoSummary() }}
        </div>
      </div>

      <div class="module-card info-card" @click="goToAddressInfo">
        <div class="module-header">
          <span class="module-title">服务地址</span>
          <el-icon class="arrow"><ArrowRight /></el-icon>
        </div>
        <div class="module-summary">
          {{ getAddressSummary() }}
        </div>
      </div>

      <div class="module-card info-card" @click="goToHealthInfo">
        <div class="module-header">
          <span class="module-title">健康档案</span>
          <el-icon class="arrow"><ArrowRight /></el-icon>
        </div>
        <div class="module-summary">
          {{ getHealthSummary() }}
        </div>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <div class="nav-item" @click="goToHome">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </div>
      <div class="nav-item" @click="goToServices">
        <el-icon><List /></el-icon>
        <span>服务</span>
      </div>
      <div class="nav-item active">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowRight, User, Phone, HomeFilled, List } from '@element-plus/icons-vue'
import { stationAPI, authAPI, type StationInfo } from '@/api'
import { useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()

const currentStation = ref<StationInfo | null>(null)

// 手机号脱敏
const maskPhone = (phone: string) => {
  if (!phone || phone.length !== 11) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 计算年龄
const calculateAge = (birthDate: string) => {
  if (!birthDate) return ''
  const today = new Date()
  const birth = new Date(birthDate)
  let age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--
  }
  return age
}

// 获取基本信息摘要
const getBasicInfoSummary = () => {
  const profile = userStore.profile
  if (!profile) return '未填写'

  const parts = []
  if (profile.name) parts.push(profile.name)
  if (profile.gender) parts.push(profile.gender)
  if (profile.birth_date) {
    const age = calculateAge(profile.birth_date)
    if (age) parts.push(`${age}岁`)
  }

  return parts.length > 0 ? parts.join(' · ') : '未填写'
}

// 获取联系信息摘要
const getContactInfoSummary = () => {
  const profile = userStore.profile
  if (!profile) return '未填写'

  const parts = []
  if (userStore.user?.phone) parts.push('手机号')
  if (profile.emergency_contact?.name) parts.push('紧急联系人')

  return parts.length > 0 ? parts.join('、') : '未填写'
}

// 获取地址摘要
const getAddressSummary = () => {
  const address = userStore.profile?.address
  if (!address) return '未填写'
  return address.length > 30 ? address.substring(0, 30) + '...' : address
}

// 获取健康档案摘要
const getHealthSummary = () => {
  const profile = userStore.profile
  if (!profile) return '未填写'

  const parts = []
  if (profile.health_status) parts.push(profile.health_status)
  if (profile.medical_history) {
    const history = profile.medical_history.length > 10
      ? profile.medical_history.substring(0, 10) + '...'
      : profile.medical_history
    parts.push(history)
  }

  return parts.length > 0 ? parts.join(' · ') : '未填写'
}

// 导航方法
const goToHome = () => router.push('/home')
const goToServices = () => router.push('/services')
const goToBasicInfo = () => router.push('/profile/basic')
const goToContactInfo = () => router.push('/profile/contact')
const goToAddressInfo = () => router.push('/profile/address')
const goToHealthInfo = () => router.push('/profile/health')

// 获取站点信息
const fetchStation = async () => {
  try {
    // 尝试获取用户位置
    if (userStore.profile?.latitude && userStore.profile?.longitude) {
      const result = await stationAPI.matchStation({
        latitude: userStore.profile.latitude,
        longitude: userStore.profile.longitude
      })
      currentStation.value = result.station
    }
  } catch (error) {
    console.error('获取站点信息失败:', error)
  }
}

// 加载用户信息
onMounted(async () => {
  if (!userStore.profile) {
    try {
      const result = await authAPI.checkToken()
      userStore.setUser(result.user)
      if (result.profile) {
        userStore.setProfile(result.profile)
      }
    } catch (error) {
      console.error('加载用户信息失败:', error)
    }
  }

  await fetchStation()
})
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  background: var(--bg-color, #f5f5f5);
  padding-bottom: 80px;
}

.header {
  background: white;
  padding: 15px 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 10;
}

.header h1 {
  font-size: 20px;
  font-weight: bold;
}

.content {
  padding: 16px;
  max-width: 600px;
  margin: 0 auto;
}

/* 卡片基础样式 */
.info-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: transform 0.2s, box-shadow 0.2s;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

/* 用户卡片 */
.user-card {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.avatar .el-icon {
  font-size: 32px;
}

.user-info h2 {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 4px;
  color: var(--text-color-primary, #303133);
}

.user-info p {
  font-size: 16px;
  color: var(--text-color-secondary, #606266);
}

/* 站点信息卡片 */
.station-card {
  border: 2px solid #409EFF;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.card-icon {
  font-size: 24px;
}

.card-title {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.station-info {
  margin-bottom: 12px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  font-size: 16px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-row .label {
  color: var(--text-color-secondary, #606266);
  flex-shrink: 0;
  margin-right: 12px;
}

.info-row .value {
  color: var(--text-color-primary, #303133);
  text-align: right;
  flex: 1;
}

.call-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  height: 48px;
  background: #409EFF;
  color: white;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  text-decoration: none;
  transition: background 0.2s;
}

.call-button:hover {
  background: #66b1ff;
}

.call-button.emergency {
  background: #F56C6C;
}

.call-button.emergency:hover {
  background: #f78989;
}

/* 紧急联系人卡片 */
.emergency-card {
  border: 2px solid #F56C6C;
}

.emergency-info {
  margin-bottom: 12px;
}

.contact-main {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-color-primary, #303133);
  margin-bottom: 4px;
}

.contact-name {
  margin-right: 4px;
}

.contact-relation {
  color: var(--text-color-secondary, #606266);
  font-size: 16px;
}

.contact-phone {
  font-size: 16px;
  color: var(--text-color-secondary, #606266);
}

/* 信息模块卡片 */
.module-card {
  cursor: pointer;
}

.module-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.module-title {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-color-primary, #303133);
}

.arrow {
  color: var(--text-color-secondary, #909399);
  font-size: 18px;
}

.module-summary {
  font-size: 16px;
  color: var(--text-color-secondary, #606266);
}

/* 底部导航 */
.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  justify-content: space-around;
  padding: 12px 0;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: var(--text-color-secondary, #909399);
  font-size: 12px;
}

.nav-item .el-icon {
  font-size: 24px;
}

.nav-item.active {
  color: var(--color-primary, #409EFF);
}
</style>
