<template>
  <div class="banner-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>轮播图管理</h3>
            <p class="subtitle">管理首页轮播图</p>
          </div>
          <div class="header-actions">
            <el-button type="primary" :icon="Plus" @click="handleAdd">
              新增轮播图
            </el-button>
            <el-button
              :icon="Refresh"
              :loading="loading"
              @click="loadBanners"
            >
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <!-- 轮播图列表 -->
      <el-table
        v-loading="loading"
        :data="bannerList"
        stripe
        style="width: 100%"
        empty-text="暂无轮播图"
      >
        <!-- 排序 -->
        <el-table-column prop="sort" label="排序" width="80" />

        <!-- 图片预览 -->
        <el-table-column label="图片" width="160">
          <template #default="{ row }">
            <el-image
              :src="getImageUrl(row.image_url)"
              :preview-src-list="[getImageUrl(row.image_url)]"
              fit="cover"
              style="width: 120px; height: 60px; border-radius: 4px"
            />
          </template>
        </el-table-column>

        <!-- 标题 -->
        <el-table-column prop="title" label="标题" min-width="150" />

        <!-- 链接类型 -->
        <el-table-column label="链接类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getLinkTypeText(row.link_type) }}</el-tag>
          </template>
        </el-table-column>

        <!-- 链接值 -->
        <el-table-column prop="link_value" label="链接值" min-width="150" show-overflow-tooltip />

        <!-- 状态 -->
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 创建时间 -->
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click="handleEdit(row)">
              编辑
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
      :title="isEdit ? '编辑轮播图' : '新增轮播图'"
      width="500px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入标题" />
        </el-form-item>

        <el-form-item label="图片" prop="image_url">
          <div class="image-upload">
            <el-input v-model="formData.image_url" placeholder="请输入图片URL或上传图片" style="flex: 1" />
            <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :show-file-list="false"
              :on-success="handleImageUploadSuccess"
              :on-error="handleUploadError"
              :before-upload="beforeUpload"
              accept="image/*"
            >
              <el-button type="primary" :icon="Upload" style="margin-left: 10px">
                上传
              </el-button>
            </el-upload>
          </div>
          <div v-if="formData.image_url" class="image-preview">
            <el-image
              :src="getImageUrl(formData.image_url)"
              fit="cover"
              style="width: 200px; height: 100px; border-radius: 4px; margin-top: 10px"
            />
          </div>
        </el-form-item>

        <el-form-item label="链接类型" prop="link_type">
          <el-select v-model="formData.link_type" placeholder="请选择链接类型" style="width: 100%">
            <el-option label="无链接" value="none" />
            <el-option label="内部页面" value="page" />
            <el-option label="外部链接" value="url" />
            <el-option label="新闻详情" value="news" />
          </el-select>
        </el-form-item>

        <el-form-item v-if="formData.link_type !== 'none'" label="链接值" prop="link_value">
          <el-input v-model="formData.link_value" placeholder="请输入链接值" />
        </el-form-item>

        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import { Plus, Refresh, Upload } from '@element-plus/icons-vue'
import { bannerApi } from '@/api'
import { useAuthStore } from '@/store/modules/auth'
import type { Banner, BannerRequest } from '@/types/api'
import dayjs from 'dayjs'

const authStore = useAuthStore()

// 上传配置
const uploadUrl = computed(() => `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api/v1'}/b/upload`)
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`,
}))

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 轮播图列表
const bannerList = ref<Banner[]>([])

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
const formData = reactive<BannerRequest>({
  title: '',
  image_url: '',
  link_type: 'none',
  link_value: '',
  sort: 0,
  status: 'active',
  station_id: 0,
})

// 表单验证规则
const formRules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  image_url: [{ required: true, message: '请输入图片URL', trigger: 'blur' }],
  link_type: [{ required: true, message: '请选择链接类型', trigger: 'change' }],
}

/**
 * 加载轮播图列表
 */
async function loadBanners() {
  try {
    loading.value = true
    const response = await bannerApi.getBanners({
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    const { items, total } = response.data
    pagination.total = total
    bannerList.value = items
  } catch (error) {
    console.error('Failed to load banners:', error)
  } finally {
    loading.value = false
  }
}

/**
 * 新增轮播图
 */
function handleAdd() {
  isEdit.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

/**
 * 编辑轮播图
 */
function handleEdit(banner: Banner) {
  isEdit.value = true
  editingId.value = banner.id
  Object.assign(formData, {
    title: banner.title,
    image_url: banner.image_url,
    link_type: banner.link_type,
    link_value: banner.link_value,
    sort: banner.sort,
    status: banner.status,
    station_id: banner.station_id,
  })
  dialogVisible.value = true
}

/**
 * 删除轮播图
 */
async function handleDelete(banner: Banner) {
  try {
    await ElMessageBox.confirm(
      `确定要删除轮播图"${banner.title}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    await bannerApi.deleteBanner(banner.id)
    ElMessage.success('删除成功')
    loadBanners()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to delete banner:', error)
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
      await bannerApi.updateBanner(editingId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await bannerApi.createBanner(formData)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    loadBanners()
  } catch (error: any) {
    if (error !== false) {
      console.error('Failed to submit banner:', error)
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
    image_url: '',
    link_type: 'none',
    link_value: '',
    sort: 0,
    status: 'active',
    station_id: 0,
  })
  formRef.value?.resetFields()
}

/**
 * 获取链接类型文本
 */
function getLinkTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    none: '无链接',
    page: '内部页面',
    url: '外部链接',
    news: '新闻详情',
  }
  return typeMap[type] || type
}

/**
 * 图片上传成功
 */
const handleImageUploadSuccess: UploadProps['onSuccess'] = (response) => {
  if (response.data?.url) {
    // 移除域名部分，保留相对路径
    const url = response.data.url.replace(/^https?:\/\/[^\/]+/, '')
    formData.image_url = url
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
  loadBanners()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadBanners()
}

onMounted(() => {
  loadBanners()
})
</script>

<style scoped lang="scss">
.banner-management {
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

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .image-upload {
    display: flex;
    align-items: center;
    width: 100%;
  }

  .image-preview {
    margin-top: 10px;
  }
}
</style>
