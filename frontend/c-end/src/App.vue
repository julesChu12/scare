<template>
  <NetworkStatus />
  <router-view />
</template>

<script setup lang="ts">
import NetworkStatus from '@/components/NetworkStatus.vue'
</script>

<style>
:root {
  --primary-color: #409EFF;
  --primary-gradient: linear-gradient(135deg, #409EFF 0%, #66b1ff 100%);
  --bg-color: #f5f7fa;
  --card-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  --card-radius: 12px;
  --safe-area-bottom: env(safe-area-inset-bottom, 20px);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body,
#app {
  width: 100%;
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  font-size: var(--font-size-base, 16px);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
