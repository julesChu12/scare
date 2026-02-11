<template>
  <div class="news-list-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>站点动态</h1>
      <div style="width: 40px"></div>
    </div>

    <!-- 分类筛选 -->
    <div class="filter-tabs">
      <div
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentType === tab.value }"
        @click="selectType(tab.value)"
      >
        {{ tab.label }}
      </div>
    </div>

    <div class="content">
      <!-- 骨架屏 -->
      <div v-if="loading && newsList.length === 0" class="skeleton-list">
        <div v-for="i in 5" :key="i" class="skeleton-item">
          <div class="skeleton-tag"></div>
          <div class="skeleton-content">
            <div class="skeleton-title"></div>
            <div class="skeleton-date"></div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!loading && newsList.length === 0" class="empty-state">
        <el-icon class="empty-icon"><Document /></el-icon>
        <p class="empty-text">暂无动态</p>
      </div>

      <!-- 新闻列表 -->
      <div v-else class="news-list">
        <div
          v-for="news in newsList"
          :key="news.id"
          class="news-item"
          @click="viewNews(news)"
        >
          <div class="news-tag" :class="news.type">{{ getTagText(news.type) }}</div>
          <div class="news-content">
            <h4 class="news-title">{{ news.title }}</h4>
            <p class="news-summary" v-if="news.summary">{{ news.summary }}</p>
            <p class="news-date">{{ formatDate(news.publish_at) }}</p>
          </div>
          <el-icon class="news-arrow"><ArrowRight /></el-icon>
        </div>
      </div>

      <!-- 加载更多 -->
      <div v-if="newsList.length > 0" class="load-more">
        <span v-if="loading" class="loading-text">
          <el-icon class="is-loading"><Loading /></el-icon>
          加载中...
        </span>
        <span v-else-if="noMore" class="no-more-text">没有更多了</span>
        <span v-else class="load-more-text" @click="loadMore">加载更多</span>
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
      <div class="nav-item" @click="goToMine">
        <el-icon><User /></el-icon>
        <span>我的</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft,
  ArrowRight,
  HomeFilled,
  List,
  User,
  Loading,
  Document
} from '@element-plus/icons-vue'
import { newsAPI, type NewsItem } from '@/api/news'

const router = useRouter()

const tabs = [
  { label: '全部', value: '' },
  { label: '公告', value: 'notice' },
  { label: '活动', value: 'activity' },
  { label: '动态', value: 'news' }
]

const currentType = ref('')
const newsList = ref<NewsItem[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const noMore = computed(() => newsList.value.length >= total.value && total.value > 0)

const goBack = () => {
  router.back()
}

const selectType = (type: string) => {
  if (currentType.value === type) return
  currentType.value = type
  page.value = 1
  newsList.value = []
  fetchNews()
}

const fetchNews = async (isLoadMore = false) => {
  if (loading.value) return
  loading.value = true

  try {
    const result = await newsAPI.getList({
      page: page.value,
      page_size: pageSize,
      type: currentType.value || undefined
    })

    if (isLoadMore) {
      newsList.value = [...newsList.value, ...result.items]
    } else {
      newsList.value = result.items
    }
    total.value = result.total
  } catch (error) {
    console.error('获取新闻列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  if (loading.value || noMore.value) return
  page.value++
  fetchNews(true)
}

const viewNews = (news: NewsItem) => {
  router.push(`/news/${news.id}`)
}

const getTagText = (type: string) => {
  const map: Record<string, string> = {
    notice: '公告',
    activity: '活动',
    news: '动态'
  }
  return map[type] || '动态'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const goToHome = () => router.push('/home')
const goToServices = () => router.push('/services')
const goToMine = () => router.push('/mine')

onMounted(() => {
  fetchNews()
})
</script>

<style scoped>
.news-list-container {
  min-height: 100vh;
  background: #f5f7fa;
  display: flex;
  flex-direction: column;
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
  font-size: 18px;
  font-weight: bold;
}

/* 分类筛选 */
.filter-tabs {
  display: flex;
  background: white;
  padding: 12px 16px;
  gap: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.tab-item {
  padding: 6px 16px;
  border-radius: 16px;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s;
}

.tab-item.active {
  background: #409EFF;
  color: white;
}

/* 内容区域 */
.content {
  flex: 1;
  padding: 16px;
  padding-bottom: 80px;
}

/* 骨架屏 */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
}

.skeleton-tag {
  width: 40px;
  height: 24px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
  flex-shrink: 0;
}

.skeleton-content {
  flex: 1;
}

.skeleton-title {
  width: 80%;
  height: 18px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
  margin-bottom: 8px;
}

.skeleton-date {
  width: 40%;
  height: 14px;
  background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 4px;
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}

.empty-icon {
  font-size: 64px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 16px;
  color: #909399;
}

/* 新闻列表 */
.news-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.news-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.news-item:active {
  transform: scale(0.98);
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
  font-size: 15px;
  color: #303133;
  margin-bottom: 6px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-summary {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.news-date {
  font-size: 12px;
  color: #c0c4cc;
}

.news-arrow {
  color: #c0c4cc;
  flex-shrink: 0;
  margin-top: 4px;
}

/* 加载更多 */
.load-more {
  text-align: center;
  padding: 16px;
  font-size: 14px;
  color: #909399;
}

.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.load-more-text {
  color: #409EFF;
  cursor: pointer;
}

.no-more-text {
  color: #c0c4cc;
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
  padding: 10px 0 20px;
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
</style>
