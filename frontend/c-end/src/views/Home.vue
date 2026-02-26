<template>
  <div class="home-container">
    <!-- 顶部位置栏 -->
    <div class="location-bar">
      <div class="location-left" @click="refreshLocation">
        <el-icon><Location /></el-icon>
        <span class="location-text">{{ currentDistrict }}</span>
        <span v-if="isDefaultLocation" class="location-hint">(默认)</span>
        <el-icon class="refresh-icon"><RefreshRight /></el-icon>
      </div>
      <div class="station-right" @click="showStationInfo">
        <el-icon><OfficeBuilding /></el-icon>
        <span class="station-name">{{ stationName }}</span>
        <el-icon><ArrowRight /></el-icon>
      </div>
    </div>

    <!-- Banner 轮播 -->
    <div class="banner-section" v-if="banners.length > 0">
      <el-carousel height="160px" :interval="4000" indicator-position="outside">
        <el-carousel-item v-for="banner in banners" :key="banner.id">
          <div class="banner-item" @click="viewBanner(banner)">
            <el-image :src="banner.image_url" fit="cover" class="banner-image">
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="banner-overlay" v-if="banner.title">
               <h3>{{ banner.title }}</h3>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>
    <!-- 缺省状态或者骨架屏 -->
    <div class="banner-section" v-else>
       <div class="banner-skeleton">
          <el-skeleton-item variant="image" style="width: 100%; height: 160px; border-radius: 12px;" />
       </div>
    </div>

    <!-- 服务网格 -->
    <div class="service-section">
      <div class="section-header">
        <h2 class="section-title">服务项目</h2>
        <span class="more-link" @click="goToAllServices">
          全部服务 <el-icon><ArrowRight /></el-icon>
        </span>
      </div>
      <div class="service-grid">
        <div
          v-for="service in homeServices"
          :key="service.key"
          class="service-item"
          @click="service.key === 'more' ? goToAllServices() : selectService(service.key)"
        >
          <span class="service-icon">{{ service.icon }}</span>
          <span class="service-name">{{ service.name }}</span>
        </div>
      </div>
    </div>

    <!-- 站点动态 -->
    <div class="news-section">
      <div class="section-header">
        <h2 class="section-title">站点动态</h2>
        <span class="more-link" @click="goToNews">
          更多 <el-icon><ArrowRight /></el-icon>
        </span>
      </div>
      <div class="news-list">
        <div v-for="news in newsList" :key="news.id" class="news-item" @click="viewNews(news)">
          <div class="news-tag" :class="news.type">{{ getTagText(news.type) }}</div>
          <div class="news-content">
            <h4 class="news-title">{{ news.title }}</h4>
            <p class="news-date">{{ formatDate(news.publish_at) }}</p>
          </div>
          <el-icon class="news-arrow"><ArrowRight /></el-icon>
        </div>
      </div>
    </div>

    <!-- 底部占位 -->
    <div class="bottom-placeholder"></div>

    <!-- 底部导航 -->
    <div class="bottom-nav">
      <div class="nav-item active">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </div>
      <div class="nav-item" @click="goToServices">
        <el-icon><List /></el-icon>
        <span>服务</span>
      </div>
      <div class="nav-item" @click="goToMine">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>

    <!-- 站点信息弹窗 -->
    <el-dialog
      v-model="stationDialogVisible"
      title="服务站点信息"
      width="85%"
      center
    >
      <div v-if="currentStation" class="station-dialog-content">
        <div class="station-info-item">
          <span class="info-label">站点名称</span>
          <span class="info-value">{{ currentStation.name }}</span>
        </div>
        <div class="station-info-item" v-if="currentStation.address">
          <span class="info-label">站点地址</span>
          <span class="info-value">{{ currentStation.address }}</span>
        </div>
        <div class="station-info-item" v-if="currentStation.phone">
          <span class="info-label">联系电话</span>
          <span class="info-value phone" @click="callStation">
            {{ currentStation.phone }}
            <el-icon><Phone /></el-icon>
          </span>
        </div>
        <div class="station-info-item" v-if="currentStation.service_area">
          <span class="info-label">服务范围</span>
          <span class="info-value">{{ currentStation.service_area }}</span>
        </div>
        <div class="station-info-item" v-if="currentStation.work_hours">
          <span class="info-label">工作时间</span>
          <span class="info-value">{{ currentStation.work_hours }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="stationDialogVisible = false">关闭</el-button>
        <el-button type="primary" v-if="currentStation?.phone" @click="callStation">
          <el-icon><Phone /></el-icon>
          拨打电话
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Location,
  RefreshRight,
  OfficeBuilding,
  ArrowRight,
  HomeFilled,
  List,
  User,
  Phone,
  Picture
} from '@element-plus/icons-vue'
import { geocodeAPI, stationAPI, bannerAPI, newsAPI } from '@/api'
import type { StationInfo } from '@/api'
import type { NewsItem } from '@/api/news'
import type { Banner } from '@/api/banner'
import { getCurrentPosition } from '@/utils/coordTransform'

const router = useRouter()

// 定位数据
const currentDistrict = ref('定位中...')
const stationName = ref('匹配中...')
const currentLat = ref<number | undefined>()
const currentLng = ref<number | undefined>()
const currentStation = ref<StationInfo | null>(null)
const locationLoading = ref(false)
const isDefaultLocation = ref(false)
const stationDialogVisible = ref(false)

const banners = ref<Banner[]>([])

// 获取轮播图
const fetchBanners = async () => {
  try {
    const result = await bannerAPI.getBanners()
    banners.value = result || []
  } catch (error) {
    console.warn('获取轮播图失败:', error)
  }
}

const homeServices = ref([
  { key: 'meal', name: '送餐服务', icon: '🍱' },
  { key: 'medical', name: '就医陪护', icon: '🏥' },
  { key: 'care', name: '日常照护', icon: '👴' },
  { key: 'repair', name: '居家维修', icon: '🔧' },
  { key: 'cleaning', name: '家政保洁', icon: '🧹' },
  { key: 'company', name: '陪伴聊天', icon: '💬' },
  { key: 'emergency', name: '紧急救助', icon: '🚨' },
  { key: 'more', name: '更多服务', icon: '📋' }
])

const newsList = ref<NewsItem[]>([])

// 获取新闻列表
const fetchNews = async () => {
  try {
    const result = await newsAPI.getList({ page_size: 3 })
    newsList.value = result.items
  } catch (error) {
    console.warn('获取新闻失败:', error)
  }
}

// 获取标签文本
const getTagText = (type: string) => {
  const map: Record<string, string> = {
    notice: '公告',
    activity: '活动',
    news: '动态'
  }
  return map[type] || '动态'
}

// 格式化日期
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// 默认位置：霍营街道（GCJ-02 坐标）
const DEFAULT_LOCATION = { lat: 40.0579, lng: 116.3698, district: '昌平区霍营街道' }

// 获取定位并解析地址
const fetchLocation = async () => {
  locationLoading.value = true
  isDefaultLocation.value = false
  try {
    const pos = await getCurrentPosition()

    currentLat.value = pos.lat
    currentLng.value = pos.lng

    // 逆地理编码获取区县名称
    const geoResult = await geocodeAPI.reverseGeocode({
      latitude: pos.lat,
      longitude: pos.lng
    })
    currentDistrict.value = geoResult.district || geoResult.city || '未知位置'
  } catch (error) {
    if (import.meta.env.DEV) {
      console.group('[Home] 定位失败，回退到默认位置')
      console.warn('原因:', error)
      console.info('默认坐标:', DEFAULT_LOCATION)
      console.groupEnd()
    }
    // 回退到默认位置
    currentLat.value = DEFAULT_LOCATION.lat
    currentLng.value = DEFAULT_LOCATION.lng
    currentDistrict.value = DEFAULT_LOCATION.district
    isDefaultLocation.value = true
  } finally {
    locationLoading.value = false
  }
}

// 匹配服务站点
const fetchStation = async () => {
  try {
    if (currentLat.value === undefined || currentLng.value === undefined) {
      return
    }
    const result = await stationAPI.matchStation({
      latitude: currentLat.value,
      longitude: currentLng.value
    })
    currentStation.value = result
    stationName.value = result.name
    // 每次获取站点后刷新 Banner，因为 Banner 可能跟站点有关
    fetchBanners()
  } catch (error) {
    console.warn('匹配站点失败:', error)
    stationName.value = '暂无站点'
    // 即使没有站点，也尝试获取通用 Banner
    fetchBanners()
  }
}

// 刷新定位
const refreshLocation = async () => {
  await fetchLocation()
  await fetchStation()
}

const showStationInfo = () => {
  if (currentStation.value) {
    stationDialogVisible.value = true
  }
}

const callStation = () => {
  if (currentStation.value?.phone) {
    window.location.href = `tel:${currentStation.value.phone}`
  }
}

const goToAllServices = () => { router.push('/services') }
const selectService = (key: string) => { router.push(`/quick?type=${key}`) }
const goToNews = () => { router.push('/news') }
const viewNews = (news: any) => { router.push(`/news/${news.id}`) }

const viewBanner = (banner: Banner) => {
  if (!banner.link_url) return

  // 简单判断是否为外部链接
  if (banner.link_url.startsWith('http')) {
    window.location.href = banner.link_url
  } else {
    router.push(banner.link_url)
  }
}
const goToServices = () => { router.push('/services') }
const goToMine = () => { router.push('/mine') }

// 页面加载时获取定位、站点和新闻
onMounted(async () => {
  await fetchLocation()
  await fetchStation()
  fetchNews() // 不阻塞，异步加载新闻
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f7fa;
}

/* 顶部位置栏 */
.location-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: white;
}

.location-left,
.station-right {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #303133;
  cursor: pointer;
}

.location-left .el-icon,
.station-right .el-icon {
  color: #409EFF;
}

.location-hint {
  font-size: 12px;
  color: #E6A23C;
}

.refresh-icon {
  font-size: 12px;
  color: #909399 !important;
}

.station-name {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Banner */
.banner-section {
  padding: 12px 16px;
}

.banner-item {
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  position: relative;
  cursor: pointer;
}

.banner-image {
  width: 100%;
  height: 100%;
  display: block;
}

.banner-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.6));
  color: white;
}

.banner-overlay h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.image-placeholder {
  width: 100%;
  height: 100%;
  background: #f0f2f5;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #909399;
  font-size: 32px;
}

.banner-skeleton {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-carousel__indicators--outside) {
  margin-top: 8px;
}

:deep(.el-carousel__indicator--horizontal .el-carousel__button) {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #dcdfe6;
}

:deep(.el-carousel__indicator--horizontal.is-active .el-carousel__button) {
  background: #409EFF;
  width: 20px;
  border-radius: 4px;
}

/* 服务网格 */
.service-section {
  padding: 16px;
  background: white;
  margin: 0 0 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.more-link {
  display: flex;
  align-items: center;
  gap: 2px;
  font-size: 14px;
  color: #909399;
  cursor: pointer;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.service-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 12px 8px;
  background: #f8f9fc;
  border-radius: 12px;
  transition: transform 0.2s, background 0.2s;
}

.service-item:active {
  transform: scale(0.95);
  background: #eef3fb;
}

.service-icon {
  font-size: 36px;
}

.service-name {
  font-size: 13px;
  color: #303133;
  text-align: center;
}

/* 站点动态 */
.news-section {
  padding: 16px;
  background: white;
}

.news-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.news-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafc;
  border-radius: 8px;
  cursor: pointer;
}

.news-tag {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  flex-shrink: 0;
}

.news-tag.notice {
  background: #fef0f0;
  color: #f56c6c;
}

.news-tag.activity {
  background: #f0f9eb;
  color: #67c23a;
}

.news-tag.news {
  background: #ecf5ff;
  color: #409eff;
}

.news-content {
  flex: 1;
  min-width: 0;
}

.news-title {
  font-size: 14px;
  color: #303133;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.news-date {
  font-size: 12px;
  color: #909399;
}

.news-arrow {
  color: #c0c4cc;
  flex-shrink: 0;
}

/* 底部占位 */
.bottom-placeholder {
  height: calc(60px + env(safe-area-inset-bottom, 20px));
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
  padding: 10px 0 calc(10px + env(safe-area-inset-bottom, 20px));
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: #909399;
  font-size: 12px;
}

.nav-item .el-icon {
  font-size: 24px;
}

.nav-item.active {
  color: #409EFF;
}

/* 站点信息弹窗 */
.station-dialog-content {
  padding: 8px 0;
}

.station-info-item {
  display: flex;
  flex-direction: column;
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.station-info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
}

.info-value {
  font-size: 15px;
  color: #303133;
  line-height: 1.5;
}

.info-value.phone {
  color: #409EFF;
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
}
</style>
