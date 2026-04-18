<template>
  <div class="news-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>新闻管理</h3>
            <p class="subtitle">管理新闻、公告和活动信息</p>
          </div>
          <div class="header-actions">
            <el-button
              v-permission="'content:news:create'"
              type="primary"
              :icon="Plus"
              @click="handleAdd"
            >
              新增新闻
            </el-button>
            <el-button :icon="Refresh" :loading="loading" @click="loadNewsList">
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="类型">
            <el-select
              v-model="filterForm.type"
              placeholder="全部类型"
              clearable
              style="width: 120px"
              @change="handleFilter"
            >
              <el-option label="新闻" value="news" />
              <el-option label="公告" value="notice" />
              <el-option label="活动" value="activity" />
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select
              v-model="filterForm.status"
              placeholder="全部状态"
              clearable
              style="width: 120px"
              @change="handleFilter"
            >
              <el-option label="草稿" value="draft" />
              <el-option label="已发布" value="published" />
              <el-option label="已归档" value="archived" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- 新闻列表 -->
      <el-table
        v-loading="loading"
        :data="newsList"
        stripe
        style="width: 100%"
        empty-text="暂无新闻"
      >
        <!-- 封面图 -->
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.cover_url"
              :src="getImageUrl(row.cover_url)"
              :preview-src-list="[getImageUrl(row.cover_url)]"
              preview-teleported
              fit="cover"
              style="width: 60px; height: 40px; border-radius: 4px"
            />
            <span v-else class="no-cover">无封面</span>
          </template>
        </el-table-column>

        <!-- 标题 -->
        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />

        <!-- 类型 -->
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 状态 -->
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 浏览量 -->
        <el-table-column prop="view_count" label="浏览量" width="90" align="center" />

        <!-- 发布时间 -->
        <el-table-column label="发布时间" width="160">
          <template #default="{ row }">
            {{ row.status === 'published' ? formatDateTime(row.publish_at) : '-' }}
          </template>
        </el-table-column>

        <!-- 创建时间 -->
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="success"
              size="small"
              link
              @click="handlePublish(row)"
            >
              发布
            </el-button>
            <el-button
              v-if="row.status === 'published'"
              type="warning"
              size="small"
              link
              @click="handleArchive(row)"
            >
              归档
            </el-button>
            <el-button type="danger" size="small" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑新闻' : '新增新闻'"
      width="1100px"
      top="2vh"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-position="top"
        style="height: 600px; overflow: hidden;"
      >
        <el-row :gutter="24" style="height: 100%">
          <!-- 左侧：元数据 -->
          <el-col :span="6" style="height: 100%; overflow-y: auto; padding-right: 10px; border-right: 1px solid #f0f0f0;">
            <el-form-item label="标题" prop="title">
              <el-input v-model="formData.title" placeholder="请输入新闻标题" maxlength="200" show-word-limit />
            </el-form-item>

            <el-row :gutter="12">
              <el-col :span="24">
                <el-form-item label="类型" prop="type">
                  <el-select v-model="formData.type" placeholder="请选择类型" style="width: 100%">
                    <el-option label="新闻" value="news" />
                    <el-option label="公告" value="notice" />
                    <el-option label="活动" value="activity" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="24">
                <el-form-item label="状态" prop="status">
                  <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                    <el-option label="草稿" value="draft" />
                    <el-option label="发布" value="published" />
                    <el-option label="归档" value="archived" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="封面图" prop="cover_url">
              <div class="cover-upload-wrapper">
                <el-input
                  v-model="formData.cover_url"
                  placeholder="图片URL或上传"
                  size="small"
                  style="margin-bottom: 8px"
                />
                <el-upload
                  class="upload-box"
                  drag
                  :action="uploadUrl"
                  :headers="uploadHeaders"
                  :show-file-list="false"
                  :on-success="handleCoverUploadSuccess"
                  :on-error="handleUploadError"
                  :before-upload="beforeUpload"
                  accept="image/*"
                >
                  <div v-if="formData.cover_url" class="cover-preview-box">
                    <img :src="getImageUrl(formData.cover_url)" class="preview-img" />
                    <div class="hover-mask">
                      <el-icon><Upload /></el-icon>
                      <span>更换封面</span>
                    </div>
                  </div>
                  <div v-else class="upload-placeholder">
                    <el-icon class="el-icon--upload"><upload-filled /></el-icon>
                    <div class="el-upload__text">拖拽或点击上传</div>
                  </div>
                </el-upload>
              </div>
            </el-form-item>

            <el-form-item label="摘要" prop="summary">
              <el-input
                v-model="formData.summary"
                type="textarea"
                :rows="4"
                placeholder="请输入新闻摘要（建议100字以内）"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
          </el-col>

          <!-- 右侧：富文本编辑器 -->
          <el-col :span="18" style="height: 100%; display: flex; flex-direction: column;">
            <el-form-item label="正文内容" prop="content" style="flex: 1; margin-bottom: 0; display: flex; flex-direction: column;">
              <div style="border: 1px solid #dcdfe6; border-radius: 4px; display: flex; flex-direction: column; height: 100%; width: 100%;">
                <Toolbar
                  style="border-bottom: 1px solid #dcdfe6"
                  :editor="editorRef"
                  :defaultConfig="toolbarConfig"
                  :mode="mode"
                />
                <Editor
                  style="flex: 1; overflow-y: hidden;"
                  v-model="formData.content"
                  :defaultConfig="editorConfig"
                  :mode="mode"
                  @onCreated="handleCreated"
                />
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, shallowRef, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import { Plus, Refresh, Upload, UploadFilled } from '@element-plus/icons-vue'
import { newsApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import type { News, NewsRequest, NewsType, NewsStatus } from '@/types/api'
import dayjs from 'dayjs'
import '@wangeditor/editor/dist/css/style.css' // 引入 css
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

const authStore = useAuthStore()

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()
const mode = 'default' // 或 'simple'

// 工具栏配置
const toolbarConfig = {}

// 编辑器配置
const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      // 上传图片配置
      server: `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api/v1'}/b/upload`,
      fieldName: 'file',
      headers: {
        Authorization: `Bearer ${authStore.token}`,
      },
      // 自定义插入图片
      customInsert(res: any, insertFn: any) {
        // res 即服务端的返回结果
        if (res.code === 0 && res.data && res.data.url) {
          // 从 res 中找到 url 文件的服务端地址
          // 我们的 API 返回的是相对路径或完整 URL，需要处理
          let url = res.data.url
          if (!url.startsWith('http')) {
             // 拼接完整路径用于显示
             const baseUrl = import.meta.env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:3001'
             url = `${baseUrl}${url}`
          }
          const alt = ''
          const href = ''
          // 插入图片
          insertFn(url, alt, href)
        } else {
          ElMessage.error('图片上传失败: ' + (res.msg || '未知错误'))
        }
      },
    },
  },
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})

const handleCreated = (editor: any) => {
  editorRef.value = editor // 记录 editor 实例，重要！
}

// 上传配置
const uploadUrl = computed(() => `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api/v1'}/b/upload`)
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`,
}))

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 新闻列表
const newsList = ref<News[]>([])

// 筛选表单
const filterForm = reactive({
  type: '' as NewsType | '',
  status: '' as NewsStatus | '',
})

// 分页参数
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

// 表单数据
const formData = reactive<NewsRequest>({
  title: '',
  summary: '',
  content: '',
  cover_url: '',
  type: 'news',
  status: 'draft',
  station_id: 0,
})

// 表单验证规则
const formRules: FormRules = {
  title: [
    { required: true, message: '请输入新闻标题', trigger: 'blur' },
    { min: 2, max: 200, message: '标题长度在 2 到 200 个字符', trigger: 'blur' },
  ],
  type: [{ required: true, message: '请选择新闻类型', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

/**
 * 加载新闻列表
 */
async function loadNewsList() {
  try {
    loading.value = true
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
    }
    if (filterForm.type) {
      params.type = filterForm.type
    }
    if (filterForm.status) {
      params.status = filterForm.status
    }

    const response = await newsApi.getNewsList(params)
    const { items, total } = response.data
    pagination.total = total
    newsList.value = items || []
  } catch (error) {
    console.error('Failed to load news:', error)
    ElMessage.error('加载新闻列表失败')
  } finally {
    loading.value = false
  }
}

/**
 * 筛选变化
 */
function handleFilter() {
  pagination.page = 1
  loadNewsList()
}

/**
 * 新增新闻
 */
function handleAdd() {
  isEdit.value = false
  editingId.value = null
  resetForm()
  refreshEditorToken()
  dialogVisible.value = true
}

/**
 * 编辑新闻
 */
function handleEdit(news: News) {
  isEdit.value = true
  editingId.value = news.id
  Object.assign(formData, {
    title: news.title,
    summary: news.summary,
    content: news.content,
    cover_url: news.cover_url,
    type: news.type,
    status: news.status,
    station_id: news.station_id,
  })
  refreshEditorToken()
  dialogVisible.value = true
}

/**
 * 刷新编辑器上传 Token
 */
function refreshEditorToken() {
  const editor = editorRef.value
  if (editor) {
    // 重新设置 headers
    // 注意：menu key 是 'uploadImage'
    const menuConfig = editor.getMenuConfig('uploadImage')
    menuConfig.headers = {
      Authorization: `Bearer ${authStore.token}`,
    }
  }
}

/**
 * 发布新闻
 */
async function handlePublish(news: News) {
  try {
    await ElMessageBox.confirm(
      `确定要发布新闻"${news.title}"吗？`,
      '发布确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info',
      }
    )

    await newsApi.updateNews(news.id, { ...news, status: 'published' })
    ElMessage.success('发布成功')
    loadNewsList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to publish news:', error)
      ElMessage.error('发布失败')
    }
  }
}

/**
 * 归档新闻
 */
async function handleArchive(news: News) {
  try {
    await ElMessageBox.confirm(
      `确定要归档新闻"${news.title}"吗？归档后将不再显示给用户。`,
      '归档确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await newsApi.updateNews(news.id, { ...news, status: 'archived' })
    ElMessage.success('归档成功')
    loadNewsList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to archive news:', error)
      ElMessage.error('归档失败')
    }
  }
}

/**
 * 删除新闻
 */
async function handleDelete(news: News) {
  try {
    await ElMessageBox.confirm(
      `确定要删除新闻"${news.title}"吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await newsApi.deleteNews(news.id)
    ElMessage.success('删除成功')
    loadNewsList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete news:', error)
      ElMessage.error('删除失败')
    }
  }
}

/**
 * 提交表单
 */
async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value && editingId.value) {
      await newsApi.updateNews(editingId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await newsApi.createNews(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadNewsList()
  } catch (error: any) {
    if (error !== false) {
      console.error('Failed to submit news:', error)
      ElMessage.error('操作失败')
    }
  } finally {
    submitting.value = false
  }
}

/**
 * 重置表单
 */
function resetForm() {
  Object.assign(formData, {
    title: '',
    summary: '',
    content: '',
    cover_url: '',
    type: 'news',
    status: 'draft',
    station_id: 0,
  })
  formRef.value?.resetFields()
}

/**
 * 封面上传成功
 */
const handleCoverUploadSuccess: UploadProps['onSuccess'] = (response) => {
  if (response.data?.url) {
    // 移除域名部分，保留相对路径
    const url = response.data.url.replace(/^https?:\/\/[^\/]+/, '')
    formData.cover_url = url
    ElMessage.success('上传成功')
  }
}

/**
 * 上传失败
 */
const handleUploadError: UploadProps['onError'] = () => {
  ElMessage.error('上传失败')
}

/**
 * 上传前校验
 */
const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

/**
 * 获取类型文本
 */
function getTypeText(type: NewsType): string {
  const typeMap: Record<NewsType, string> = {
    news: '新闻',
    notice: '公告',
    activity: '活动',
  }
  return typeMap[type] || type
}

/**
 * 获取类型标签样式
 */
function getTypeTagType(type: NewsType): string {
  const typeMap: Record<NewsType, string> = {
    news: 'primary',
    notice: 'warning',
    activity: 'success',
  }
  return typeMap[type] || 'info'
}

/**
 * 获取状态文本
 */
function getStatusText(status: NewsStatus): string {
  const statusMap: Record<NewsStatus, string> = {
    draft: '草稿',
    published: '已发布',
    archived: '已归档',
  }
  return statusMap[status] || status
}

/**
 * 获取状态标签样式
 */
function getStatusTagType(status: NewsStatus): string {
  const statusMap: Record<NewsStatus, string> = {
    draft: 'info',
    published: 'success',
    archived: 'warning',
  }
  return statusMap[status] || 'info'
}

/**
 * 格式化日期时间
 */
function formatDateTime(dateTime?: string): string {
  if (!dateTime) return '-'
  return dayjs(dateTime).format('YYYY-MM-DD HH:mm')
}

/**
 * 获取图片完整URL
 */
function getImageUrl(path: string): string {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${import.meta.env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:3001'}${path}`
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadNewsList()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadNewsList()
}

onMounted(() => {
  loadNewsList()
})
</script>

<style scoped lang="scss">
.news-management {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;

    h3 {
      margin: 0 0 5px 0;
      font-size: 18px;
      font-weight: 500;
    }

    .subtitle {
      margin: 0;
      font-size: 13px;
      color: #909399;
    }

    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .filter-container {
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid #ebeef5;

    .filter-form {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;

      .el-form-item {
        margin-bottom: 0;
      }
    }
  }

  .no-cover {
    color: #909399;
    font-size: 12px;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .cover-upload {
    display: flex;
    align-items: center;
    width: 100%;
  }

  .cover-preview {
    margin-top: 10px;
  }

  .cover-upload-wrapper {
    width: 100%;
    
    .upload-box {
      width: 100%;
      
      :deep(.el-upload) {
        width: 100%;
      }
      
      :deep(.el-upload-dragger) {
        width: 100%;
        height: 140px;
        padding: 0;
        display: flex;
        justify-content: center;
        align-items: center;
      }
    }

    .cover-preview-box {
      width: 100%;
      height: 100%;
      position: relative;
      overflow: hidden;
      
      .preview-img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      
      .hover-mask {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5);
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        color: white;
        opacity: 0;
        transition: opacity 0.3s;
        cursor: pointer;
        
        .el-icon {
          font-size: 24px;
          margin-bottom: 5px;
        }
      }
      
      &:hover .hover-mask {
        opacity: 1;
      }
    }

    .upload-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      
      .el-icon--upload {
        font-size: 32px;
        color: #909399;
        margin-bottom: 8px;
      }
      
      .el-upload__text {
        font-size: 12px;
        color: #606266;
      }
    }
  }
}
</style>
