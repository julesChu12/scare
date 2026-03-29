<template>
  <transition name="slide-down">
    <div v-if="showOffline" class="network-status-bar">
      <el-icon><WarningFilled /></el-icon>
      <span>当前处于离线状态，部分功能可能不可用</span>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { WarningFilled } from '@element-plus/icons-vue'

const showOffline = ref(!navigator.onLine)

function handleOnline() {
  showOffline.value = false
}

function handleOffline() {
  showOffline.value = true
}

onMounted(() => {
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
})
</script>

<style scoped>
.network-status-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 16px;
  background: #f56c6c;
  color: #fff;
  font-size: calc(13px * var(--font-scale));
  font-weight: 500;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-100%);
  opacity: 0;
}
</style>
