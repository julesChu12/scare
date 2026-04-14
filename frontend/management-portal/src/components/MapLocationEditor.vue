<template>
  <div class="map-location-editor">
    <div ref="mapContainer" class="map-container"></div>
    <div class="map-overlay-hint">
      点击地图设置位置，或拖拽标记微调坐标
    </div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import AMapLoader from '@amap/amap-jsapi-loader'
import { ElMessage } from 'element-plus'

const securityCode = import.meta.env.VITE_AMAP_SECURITY_JS_CODE
if (securityCode) {
  ;(window as any)._AMapSecurityConfig = { securityJsCode: securityCode }
}

interface Props {
  latitude?: number
  longitude?: number
  zoom?: number
}

const props = withDefaults(defineProps<Props>(), {
  zoom: 15,
})

const emit = defineEmits<{
  'update:address': [value: string]
  'update:latitude': [value: number | undefined]
  'update:longitude': [value: number | undefined]
}>()

const defaultCenter: [number, number] = [116.397428, 39.90923]

const mapContainer = ref<HTMLDivElement | null>(null)
let AMapObj: any = null
let map: any = null
let marker: any = null
let geocoder: any = null

const hasValidLocation = (latitude?: number, longitude?: number) => (
  typeof latitude === 'number'
  && Number.isFinite(latitude)
  && typeof longitude === 'number'
  && Number.isFinite(longitude)
)

const buildPosition = (longitude: number, latitude: number): [number, number] => [longitude, latitude]

const clearPosition = () => {
  if (!map || !marker) return

  map.remove(marker)
  marker.off?.('dragend')
  marker = null
}

const syncPosition = (longitude: number, latitude: number, shouldCenter = true) => {
  if (!map || !AMapObj) return

  const position = buildPosition(longitude, latitude)
  if (!marker) {
    marker = new AMapObj.Marker({
      position,
      draggable: true,
      title: '站点位置',
    })
    marker.on('dragend', (event: any) => {
      const lng = event.lnglat.getLng()
      const lat = event.lnglat.getLat()
      syncPosition(lng, lat, false)
      emit('update:longitude', lng)
      emit('update:latitude', lat)
      void reverseGeocodePosition(lng, lat).catch((error) => {
        console.warn('拖拽标记反查地址失败:', error)
      })
    })
    map.add(marker)
  } else {
    marker.setPosition(position)
  }

  if (shouldCenter) {
    map.setCenter(position)
  }
}

const initMap = async () => {
  if (!mapContainer.value) return

  const amapKey = import.meta.env.VITE_AMAP_KEY
  if (!amapKey) {
    ElMessage.error('未配置高德地图 Key，请在 .env 中设置 VITE_AMAP_KEY')
    return
  }

  AMapObj = await AMapLoader.load({
    key: amapKey,
    version: '2.0',
    plugins: ['AMap.Geocoder', 'AMap.ToolBar', 'AMap.Scale'],
  })

  if (!mapContainer.value) return

  const center = hasValidLocation(props.latitude, props.longitude)
    ? buildPosition(props.longitude!, props.latitude!)
    : defaultCenter

  map = new AMapObj.Map(mapContainer.value, {
    zoom: props.zoom,
    center,
    viewMode: '2D',
  })

  map.addControl(new AMapObj.ToolBar())
  map.addControl(new AMapObj.Scale())
  geocoder = new AMapObj.Geocoder({
    city: '全国',
  })

  map.on('click', (event: any) => {
    const lng = event.lnglat.getLng()
    const lat = event.lnglat.getLat()
    syncPosition(lng, lat, false)
    emit('update:longitude', lng)
    emit('update:latitude', lat)
    void reverseGeocodePosition(lng, lat).catch((error) => {
      console.warn('地图选点反查地址失败:', error)
    })
  })

  if (hasValidLocation(props.latitude, props.longitude)) {
    syncPosition(props.longitude!, props.latitude!)
  }
}

const tryInitMap = async (attempts = 0) => {
  await nextTick()
  if (!mapContainer.value) return

  if ((mapContainer.value.clientWidth === 0 || mapContainer.value.clientHeight === 0) && attempts < 10) {
    window.setTimeout(() => {
      void tryInitMap(attempts + 1)
    }, 200)
    return
  }

  if (!map) {
    try {
      await initMap()
    } catch (error) {
      console.error('地图加载失败:', error)
      ElMessage.error('地图加载失败，请检查网络或 Key 配置')
    }
  }
}

const reverseGeocodePosition = async (longitude: number, latitude: number) => {
  if (!geocoder) {
    throw new Error('geocoder not ready')
  }

  return new Promise<string>((resolve, reject) => {
    let settled = false
    const finishReject = (error: Error) => {
      if (settled) return
      settled = true
      window.clearTimeout(timeoutId)
      reject(error)
    }
    const finishResolve = (address: string) => {
      if (settled) return
      settled = true
      window.clearTimeout(timeoutId)
      resolve(address)
    }
    const timeoutId = window.setTimeout(() => {
      finishReject(new Error('reverse geocode timeout'))
    }, 10000)

    geocoder.getAddress(buildPosition(longitude, latitude), (status: string, result: any) => {
      if (settled) return

      const formattedAddress = result?.regeocode?.formattedAddress?.trim()
      if (status !== 'complete' || !formattedAddress) {
        finishReject(new Error('reverse geocode failed'))
        return
      }

      emit('update:address', formattedAddress)
      finishResolve(formattedAddress)
    })
  })
}

const geocodeAddress = async (address: string) => {
  if (!geocoder) {
    throw new Error('geocoder not ready')
  }
  const normalizedAddress = address.trim()
  if (!normalizedAddress) {
    throw new Error('address required')
  }

  return new Promise<{ latitude: number; longitude: number; formattedAddress?: string }>((resolve, reject) => {
    let settled = false
    const finishReject = (error: Error) => {
      if (settled) return
      settled = true
      window.clearTimeout(timeoutId)
      reject(error)
    }
    const finishResolve = (payload: { latitude: number; longitude: number; formattedAddress?: string }) => {
      if (settled) return
      settled = true
      window.clearTimeout(timeoutId)
      resolve(payload)
    }
    const timeoutId = window.setTimeout(() => {
      finishReject(new Error('geocode timeout'))
    }, 10000)

    geocoder.getLocation(normalizedAddress, (status: string, result: any) => {
      if (settled) return
      const geocodeResult = result?.geocodes?.[0]
      if (status !== 'complete' || !geocodeResult?.location) {
        finishReject(new Error('geocode failed'))
        return
      }

      const longitude = typeof geocodeResult.location.getLng === 'function'
        ? geocodeResult.location.getLng()
        : geocodeResult.location.lng
      const latitude = typeof geocodeResult.location.getLat === 'function'
        ? geocodeResult.location.getLat()
        : geocodeResult.location.lat

      if (!Number.isFinite(longitude) || !Number.isFinite(latitude)) {
        finishReject(new Error('invalid geocode location'))
        return
      }

      syncPosition(longitude, latitude)
      if (geocodeResult.formattedAddress) {
        emit('update:address', geocodeResult.formattedAddress)
      }
      emit('update:longitude', longitude)
      emit('update:latitude', latitude)
      finishResolve({
        latitude,
        longitude,
        formattedAddress: geocodeResult.formattedAddress,
      })
    })
  })
}

defineExpose({
  geocodeAddress,
})

watch(
  () => [props.longitude, props.latitude],
  ([longitude, latitude]) => {
    if (!map) return
    if (!hasValidLocation(latitude, longitude)) {
      clearPosition()
      return
    }
    syncPosition(longitude!, latitude!)
  }
)

onMounted(() => {
  void tryInitMap()
})

onUnmounted(() => {
  clearPosition()
  if (map) {
    map.destroy()
    map = null
  }
  geocoder = null
  AMapObj = null
})
</script>

<style scoped>
.map-location-editor {
  position: relative;
  width: 100%;
  height: 100%;
}

.map-container {
  width: 100%;
  height: 100%;
  min-height: 320px;
}

.map-overlay-hint {
  position: absolute;
  left: 12px;
  top: 12px;
  z-index: 10;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(17, 24, 39, 0.78);
  color: #fff;
  font-size: 12px;
  line-height: 1.4;
  pointer-events: none;
}
</style>
