<template>
  <el-dialog
    v-model="visible"
    title="服务评价"
    width="90%"
    :max-width="400"
    :close-on-click-modal="false"
  >
    <div class="rating-content">
      <p class="rating-hint">请对本次服务进行评价</p>

      <div class="rating-stars">
        <el-rate
          v-model="rating"
          :texts="['很差', '较差', '一般', '满意', '非常满意']"
          show-text
          size="large"
        />
      </div>

      <el-input
        v-model="comment"
        type="textarea"
        placeholder="请输入您的评价和建议（可选）"
        :rows="4"
        maxlength="500"
        show-word-limit
      />
    </div>

    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="loading" :disabled="rating === 0">
        提交评价
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { requestsAPI } from '@/api'

const props = defineProps<{
  modelValue: boolean
  requestId: number
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success', rating: number, comment: string): void
}>()

const visible = ref(props.modelValue)
const rating = ref(0)
const comment = ref('')
const loading = ref(false)

// 同步 visible 状态
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    // 重置表单
    rating.value = 0
    comment.value = ''
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleCancel = () => {
  visible.value = false
}

const handleSubmit = async () => {
  if (rating.value === 0) {
    ElMessage.warning('请选择评分')
    return
  }

  loading.value = true
  try {
    await requestsAPI.rateRequest(props.requestId, {
      rating: rating.value,
      comment: comment.value || undefined
    })

    ElMessage.success('评价提交成功')
    emit('success', rating.value, comment.value)
    visible.value = false
  } catch (error) {
    console.error('评价提交失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.rating-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rating-hint {
  text-align: center;
  color: var(--text-color-regular, #606266);
  font-size: var(--font-size-base, 16px);
}

.rating-stars {
  display: flex;
  justify-content: center;
}

.rating-stars :deep(.el-rate) {
  height: auto;
}

.rating-stars :deep(.el-rate__icon) {
  font-size: 32px;
}
</style>
