<template>
  <router-view />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()

// 应用启动时从 localStorage 恢复登录状态
onMounted(() => {
  authStore.init()
})
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 图片预览弹窗强制置顶，不被任何父级 overflow/clip 裁剪 */
.el-image-viewer__wrapper {
  z-index: 9999 !important;
  position: fixed !important;
  top: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  left: 0 !important;
  overflow: visible !important;
}
.el-image-viewer__mask {
  z-index: 1 !important;
  position: absolute !important;
}
.el-image-viewer__canvas {
  z-index: 2 !important;
  position: relative !important;
}
.el-image-viewer__btn {
  z-index: 3 !important;
  position: absolute !important;
}
</style>
