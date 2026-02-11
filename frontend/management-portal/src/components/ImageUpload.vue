<template>
  <div class="image-upload">
    <el-upload
      v-model:file-list="fileList"
      :action="uploadUrl"
      :headers="uploadHeaders"
      list-type="picture-card"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-remove="handleRemove"
      :before-upload="beforeUpload"
      :limit="limit"
      accept="image/*"
    >
      <el-icon><Plus /></el-icon>
    </el-upload>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { UploadUserFile, UploadFile } from 'element-plus'
import { useAuthStore } from '@/store/modules/auth'

interface Props {
  modelValue?: string[]
  limit?: number
  maxSize?: number // MB
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => [],
  limit: 9,
  maxSize: 5,
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const authStore = useAuthStore()

// 上传地址
const uploadUrl = computed(() =>
  `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api/v1'}/b/upload`
)

// 上传请求头（包含 Token）
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`,
}))

// 文件列表
const fileList = ref<UploadUserFile[]>([])

// 监听 modelValue 变化，初始化文件列表
watch(
  () => props.modelValue,
  (urls) => {
    if (urls && urls.length > 0) {
      fileList.value = urls.map((url, index) => ({
        name: `image-${index}`,
        url: url.startsWith('http') ? url : `${import.meta.env.VITE_API_BASE_URL?.replace('/api/v1', '') || 'http://localhost:3001'}${url}`,
        uid: Date.now() + index,
      }))
    }
  },
  { immediate: true }
)

/**
 * 上传前校验
 */
function beforeUpload(file: File): boolean {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }

  const isLtMaxSize = file.size / 1024 / 1024 < props.maxSize
  if (!isLtMaxSize) {
    ElMessage.error(`图片大小不能超过 ${props.maxSize}MB！`)
    return false
  }

  return true
}

/**
 * 上传成功
 */
function handleSuccess(response: any, file: UploadFile) {
  if (response.msg === 'ok' && response.data?.url) {
    // 更新文件列表中的 URL
    const index = fileList.value.findIndex(f => f.uid === file.uid)
    if (index !== -1) {
      fileList.value[index].url = response.data.url
    }

    // 通知父组件
    emitUrls()
    ElMessage.success('上传成功')
  } else {
    ElMessage.error('上传失败')
  }
}

/**
 * 上传失败
 */
function handleError(error: Error) {
  console.error('Upload error:', error)
  ElMessage.error('上传失败，请重试')
}

/**
 * 删除图片
 */
function handleRemove() {
  emitUrls()
}

/**
 * 发送 URL 列表到父组件
 */
function emitUrls() {
  const urls = fileList.value
    .filter(file => file.url)
    .map(file => {
      const url = file.url!
      // 移除域名部分，保留相对路径
      return url.replace(/^https?:\/\/[^\/]+/, '')
    })
  emit('update:modelValue', urls)
}
</script>

<style scoped lang="scss">
.image-upload {
  :deep(.el-upload-list) {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  :deep(.el-upload-list__item) {
    margin: 0;
  }
}
</style>
