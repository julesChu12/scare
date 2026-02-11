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
</style>
