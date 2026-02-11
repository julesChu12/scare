<template>
  <div ref="mapContainer" class="map-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import AMapLoader from '@amap/amap-jsapi-loader'

interface Props {
  longitude: number
  latitude: number
  zoom?: number
}

const props = withDefaults(defineProps<Props>(), {
  zoom: 15,
})

const mapContainer = ref<HTMLDivElement | null>(null)
let map: any = null
let marker: any = null

/**
 * 初始化地图
 */
async function initMap() {
  try {
    const AMap = await AMapLoader.load({
      key: import.meta.env.VITE_AMAP_KEY,
      version: '2.0',
      plugins: ['AMap.Marker'],
    })

    if (!mapContainer.value) return

    // 创建地图
    map = new AMap.Map(mapContainer.value, {
      viewMode: '2D',
      zoom: props.zoom,
      center: [props.longitude, props.latitude],
    })

    // 添加标记
    marker = new AMap.Marker({
      position: [props.longitude, props.latitude],
      title: '服务位置',
    })

    map.add(marker)
  } catch (error) {
    console.error('地图加载失败:', error)
  }
}

/**
 * 更新地图中心和标记位置
 */
function updateMapCenter() {
  if (map && marker) {
    const center = [props.longitude, props.latitude]
    map.setCenter(center)
    marker.setPosition(center)
  }
}

// 监听坐标变化
watch(
  () => [props.longitude, props.latitude],
  () => {
    updateMapCenter()
  }
)

onMounted(() => {
  initMap()
})
</script>

<style scoped lang="scss">
.map-container {
  width: 100%;
  height: 100%;
  min-height: 300px;
}
</style>
