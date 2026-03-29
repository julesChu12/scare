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

      <div
        ref="ratingTrackRef"
        class="rating-stars"
        @pointerdown="handlePointerDown"
        @pointermove="handlePointerMove"
        @pointerup="handlePointerEnd"
        @pointercancel="handlePointerEnd"
      >
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
  (e: 'success', rating: number, feedback: string): void
}>()

const visible = ref(props.modelValue)
const rating = ref(0)
const comment = ref('')
const loading = ref(false)
const ratingTrackRef = ref<HTMLElement | null>(null)
const activePointerId = ref<number | null>(null)

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

const resolveRatingFromPointer = (clientX: number) => {
  const items = Array.from(
    ratingTrackRef.value?.querySelectorAll('.el-rate__item') ?? []
  ) as HTMLElement[]

  if (items.length === 0) {
    return
  }

  const rects = items.map((item) => item.getBoundingClientRect())
  const firstCenter = rects[0].left + rects[0].width / 2
  const lastRect = rects[rects.length - 1]
  const lastCenter = lastRect.left + lastRect.width / 2
  const clampedX = Math.min(Math.max(clientX, firstCenter), lastCenter)

  for (let index = 0; index < rects.length; index += 1) {
    const currentCenter = rects[index].left + rects[index].width / 2
    const nextRect = rects[index + 1]
    const nextCenter = nextRect ? nextRect.left + nextRect.width / 2 : Infinity
    const boundary = (currentCenter + nextCenter) / 2

    if (clampedX <= boundary) {
      rating.value = index + 1
      return
    }
  }

  rating.value = rects.length
}

const handlePointerDown = (event: PointerEvent) => {
  if (event.pointerType === 'mouse' && event.buttons !== 1) {
    return
  }

  activePointerId.value = event.pointerId
  ratingTrackRef.value?.setPointerCapture(event.pointerId)
  resolveRatingFromPointer(event.clientX)
}

const handlePointerMove = (event: PointerEvent) => {
  if (activePointerId.value !== event.pointerId) {
    return
  }

  if (event.cancelable) {
    event.preventDefault()
  }

  resolveRatingFromPointer(event.clientX)
}

const handlePointerEnd = (event: PointerEvent) => {
  if (activePointerId.value !== event.pointerId) {
    return
  }

  ratingTrackRef.value?.releasePointerCapture(event.pointerId)
  activePointerId.value = null
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
      feedback: comment.value || undefined
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
  touch-action: none;
  user-select: none;
  -webkit-user-select: none;
}

.rating-stars :deep(.el-rate) {
  height: auto;
}

.rating-stars :deep(.el-rate__icon) {
  font-size: calc(32px * var(--font-scale));
}
</style>
