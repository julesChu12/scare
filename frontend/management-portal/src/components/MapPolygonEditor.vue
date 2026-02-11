<template>
  <div class="map-polygon-editor">
    <!-- 地图容器 -->
    <div ref="mapContainer" class="map-container"></div>
    
    <!-- 坐标列表面板 -->
    <div class="coord-panel">
      <div class="panel-header">
        <h4>围栏顶点 ({{ points.length }})</h4>
        <div class="actions">
          <el-button size="small" type="primary" link @click="clearMap">清除</el-button>
        </div>
      </div>
      <div class="point-list">
        <div 
          v-for="(point, index) in points" 
          :key="index"
          class="point-item"
        >
          <span class="index">{{ index + 1 }}</span>
          <span class="coord">{{ point.lng.toFixed(6) }}, {{ point.lat.toFixed(6) }}</span>
        </div>
        <div v-if="points.length === 0" class="empty-text">
          点击地图或使用多边形工具绘制围栏
        </div>
      </div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar">
      <el-button-group>
        <el-button type="primary" @click="startDraw">
          <el-icon><EditPen /></el-icon> 开始绘制
        </el-button>
        <el-button type="success" @click="endDraw">
          <el-icon><Check /></el-icon> 完成编辑
        </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, shallowRef, nextTick } from 'vue'
import AMapLoader from '@amap/amap-jsapi-loader'
import { EditPen, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

export interface ZonePoint {
  lat: number
  lng: number
}

interface Props {
  modelValue?: ZonePoint[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: () => []
})

const emit = defineEmits<{
  'update:modelValue': [value: ZonePoint[]]
}>()

// 地图相关
const mapContainer = ref<HTMLElement>()
const map = shallowRef<any>(null)
const AMap = shallowRef<any>(null)
const polygon = shallowRef<any>(null)
const polyEditor = shallowRef<any>(null)

// 数据
const points = ref<{ lng: number, lat: number }[]>([])

// 初始化地图
onMounted(async () => {
  try {
    const amapKey = import.meta.env.VITE_AMAP_KEY
    const securityCode = import.meta.env.VITE_AMAP_SECURITY_JS_CODE
    if (!amapKey) {
      ElMessage.error('未配置高德地图 Key，请在 .env 中设置 VITE_AMAP_KEY')
      return
    }
    if (securityCode) {
      ;(window as any)._AMapSecurityConfig = { securityJsCode: securityCode }
    }
    AMap.value = await AMapLoader.load({
      key: amapKey,
      version: '2.0',
      plugins: ['AMap.PolygonEditor']
    })

    // 延迟一点初始化，确保 DOM 布局完成（尤其是在弹窗中）
    nextTick(() => {
      initMap()
    })
  } catch (e) {
    console.error('AMap load failed:', e)
    ElMessage.error('地图加载失败，请检查网络或 Key 配置')
  }
})

function initMap() {
  if (!mapContainer.value || !AMap.value) return

  map.value = new AMap.value.Map(mapContainer.value, {
    zoom: 14,
    center: [116.397428, 39.90923], // 默认北京，后续可定位或根据数据居中
    viewMode: '2D'
  })

  // 解析初始数据
  if (props.modelValue && props.modelValue.length > 0) {
    try {
      const path = props.modelValue.map(p => [p.lng, p.lat])
      drawPolygon(path)
      points.value = props.modelValue
      // 自动缩放以适应多边形
      if (polygon.value) {
        map.value.setFitView([polygon.value])
      }
    } catch (e) {
      console.error('Invalid polygon data:', e)
    }
  }

  initEditor()
}

function initEditor() {
  if (!map.value || !AMap.value) return

  // 如果没有多边形，创建一个空的（但在编辑器里可能需要先有实例）
  // 实际上 AMap.PolygonEditor 可以绑定到一个多边形对象
  // 如果没有多边形，我们可以在 startDraw 时处理
}

// 绘制/回显多边形
function drawPolygon(path: number[][]) {
  if (polygon.value) {
    map.value.remove(polygon.value)
  }

  polygon.value = new AMap.value.Polygon({
    path: path,
    strokeColor: "#FF33FF", 
    strokeWeight: 6,
    strokeOpacity: 0.2,
    fillOpacity: 0.4,
    fillColor: '#1791fc',
    zIndex: 50,
  })

  map.value.add(polygon.value)
}

// 开始编辑/绘制
function startDraw() {
  if (!map.value || !AMap.value) return

  if (!polygon.value) {
    // 创建一个新的空多边形用于绘制
    polygon.value = new AMap.value.Polygon({
      path: [],
      strokeColor: "#FF33FF", 
      strokeWeight: 6,
      strokeOpacity: 0.2,
      fillOpacity: 0.4,
      fillColor: '#1791fc',
      zIndex: 50,
    })
    map.value.add(polygon.value)
  }

  if (!polyEditor.value) {
    polyEditor.value = new AMap.value.PolygonEditor(map.value, polygon.value)
    
    // 监听编辑事件更新坐标
    polyEditor.value.on('addnode', updateModel)
    polyEditor.value.on('adjust', updateModel)
    polyEditor.value.on('removenode', updateModel)
    polyEditor.value.on('end', updateModel)
  } else {
    polyEditor.value.setTarget(polygon.value)
  }

  polyEditor.value.open()
  ElMessage.info('请在地图上点击绘制或拖动节点编辑')
}

// 结束编辑
function endDraw() {
  if (polyEditor.value) {
    polyEditor.value.close()
    updateModel() // 确保最后状态同步
  }
}

// 清除
function clearMap() {
  if (polyEditor.value) {
    polyEditor.value.close()
    polyEditor.value.setTarget()
  }
  if (polygon.value) {
    map.value.remove(polygon.value)
    polygon.value = null
  }
  points.value = []
  emit('update:modelValue', [])
}

// 更新数据模型
function updateModel() {
  if (!polygon.value) return
  
  const path = polygon.value.getPath()
  // path 是 AMap.LngLat 对象的数组
  const coords = path.map((p: any) => [p.lng, p.lat])
  const newPoints = coords.map((p: any) => ({ lng: p[0], lat: p[1] }))
  
  points.value = newPoints
  
  // 发送 ZonePoint 对象数组
  emit('update:modelValue', newPoints)
}

// 监听外部数据变化（例如重置表单时）
watch(() => props.modelValue, (newVal) => {
  if ((!newVal || newVal.length === 0) && polygon.value) {
    // 外部清空
    clearMap()
  } else if (newVal && newVal.length > 0) {
    // 检查是否需要重绘（简单对比长度或第一个点）
    if (!polygon.value) {
      const path = newVal.map(p => [p.lng, p.lat])
      drawPolygon(path)
      points.value = newVal
      map.value.setFitView([polygon.value])
    }
  }
}, { deep: true })

onUnmounted(() => {
  if (polyEditor.value) {
    polyEditor.value.close()
  }
  if (map.value) {
    map.value.destroy()
  }
})
</script>

<style scoped lang="scss">
  .map-polygon-editor {
    position: relative;
    width: 100%;
    height: 550px;
    border: 1px solid #dcdfe6;
  border-radius: 4px;
  display: flex;
  overflow: hidden;

  .map-container {
    flex: 1;
    height: 100%;
  }

  .coord-panel {
    width: 200px;
    background: #f5f7fa;
    border-left: 1px solid #dcdfe6;
    display: flex;
    flex-direction: column;

    .panel-header {
      padding: 10px;
      border-bottom: 1px solid #ebeef5;
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      h4 {
        margin: 0;
        font-size: 14px;
        color: #606266;
      }
    }

    .point-list {
      flex: 1;
      overflow-y: auto;
      padding: 10px;

      .point-item {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
        font-size: 12px;
        color: #606266;

        .index {
          display: inline-block;
          width: 20px;
          height: 20px;
          line-height: 20px;
          text-align: center;
          background: #e4e7ed;
          border-radius: 50%;
          margin-right: 8px;
          flex-shrink: 0;
        }

        .coord {
          font-family: monospace;
        }
      }

      .empty-text {
        text-align: center;
        color: #909399;
        font-size: 12px;
        margin-top: 20px;
      }
    }
  }

  .toolbar {
    position: absolute;
    top: 10px;
    left: 10px;
    z-index: 100;
  }
}
</style>