<template>
  <div class="news-detail-container">
    <div class="header">
      <el-button text @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h1>详情</h1>
      <div style="width: 40px"></div>
    </div>

    <div class="content">
      <div v-if="loading" class="loading">
        <el-icon class="is-loading"><Loading /></el-icon>
        <p>加载中...</p>
      </div>

      <div v-else-if="news" class="detail">
        <!-- 标题 -->
        <h2 class="news-title">{{ news.title }}</h2>

        <!-- 元信息 -->
        <div class="news-meta">
          <span class="publish-time">
            <el-icon><Clock /></el-icon>
            {{ formatTime(news.publish_at) }}
          </span>
          <el-tag v-if="news.type" size="small" :type="getTagType(news.type)">
            {{ getTagText(news.type) }}
          </el-tag>
          <span v-if="news.view_count" class="view-count">
            {{ news.view_count }} 次浏览
          </span>
        </div>

        <el-divider />

        <!-- 摘要 -->
        <p v-if="news.summary" class="news-summary">{{ news.summary }}</p>

        <!-- 正文内容 -->
        <div class="news-content" v-html="news.content"></div>

        <!-- 封面图 -->
        <div v-if="news.cover_url" class="news-cover">
          <el-image
            :src="news.cover_url"
            :preview-src-list="[news.cover_url]"
            fit="cover"
            class="cover-image"
          />
        </div>
      </div>

      <div v-else class="empty">
        <el-empty description="未找到相关内容" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Loading, Clock } from '@element-plus/icons-vue'
import { newsAPI, type NewsItem } from '@/api/news'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const news = ref<NewsItem | null>(null)

const goBack = () => {
  router.back()
}

const loadNewsDetail = async () => {
  const id = parseInt(route.params.id as string, 10)
  if (!id) return

  loading.value = true
  try {
    news.value = await newsAPI.getDetail(id)
  } catch (error) {
    console.error('加载新闻详情失败:', error)
    news.value = null
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const getTagType = (type: string) => {
  const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    'notice': 'warning',
    'activity': 'success',
    'news': 'info'
  }
  return typeMap[type] || 'info'
}

const getTagText = (type: string) => {
  const textMap: Record<string, string> = {
    'notice': '公告',
    'activity': '活动',
    'news': '动态'
  }
  return textMap[type] || '动态'
}

onMounted(() => {
  loadNewsDetail()
})
</script>

<style scoped>
.news-detail-container {
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
  max-width: 800px;
  margin: 0 auto;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-color-secondary, #909399);
}

.detail {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.news-title {
  font-size: var(--font-size-title, 22px);
  font-weight: bold;
  color: var(--text-color-primary, #303133);
  line-height: 1.4;
  margin-bottom: 16px;
}

.news-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  color: var(--text-color-secondary, #909399);
  font-size: var(--font-size-small, 14px);
}

.publish-time {
  display: flex;
  align-items: center;
  gap: 4px;
}

.source {
  color: var(--text-color-secondary, #909399);
}

.view-count {
  color: var(--text-color-secondary, #909399);
}

.news-summary {
  font-size: var(--font-size-base, 16px);
  color: var(--text-color-secondary, #606266);
  line-height: 1.6;
  margin-bottom: 16px;
  padding: 12px;
  background: #f9fafc;
  border-radius: 8px;
}

.news-content {
  font-size: var(--font-size-base, 16px);
  line-height: 1.8;
  color: var(--text-color-primary, #303133);
}

.news-content :deep(p) {
  margin-bottom: 16px;
}

.news-content :deep(ul),
.news-content :deep(ol) {
  margin-bottom: 16px;
  padding-left: 24px;
}

.news-content :deep(li) {
  margin-bottom: 8px;
}

.news-content :deep(strong) {
  color: var(--text-color-primary, #303133);
}

.news-cover {
  margin-top: 20px;
}

.cover-image {
  width: 100%;
  max-height: 300px;
  border-radius: 8px;
}

.news-images {
  margin-top: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.news-image {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  cursor: pointer;
}

.empty {
  padding: 60px 20px;
}

/* 响应式调整 */
@media (max-width: 480px) {
  .content {
    padding: 15px;
  }

  .detail {
    padding: 16px;
  }

  .news-title {
    font-size: 20px;
  }

  .news-images {
    grid-template-columns: 1fr;
  }

  .news-image {
    height: 180px;
  }
}
</style>
