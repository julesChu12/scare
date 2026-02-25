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
          <span class="coord">{{ point.lng?.toFixed(6) ?? '—' }}, {{ point.lat?.toFixed(6) ?? '—' }}</span>
        </div>
        <div v-if="points.length === 0" class="empty-text">
          点击"开始绘制"后在地图上点击添加顶点
        </div>
      </div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar">
      <el-button-group>
        <el-button type="primary" @click="startDraw" :disabled="isDrawing">
          <el-icon><EditPen /></el-icon> 开始绘制
        </el-button>
        <el-button type="success" @click="endDraw" :disabled="!isDrawing">
          <el-icon><Check /></el-icon> 完成编辑
        </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import AMapLoader from '@amap/amap-jsapi-loader'
import { EditPen, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// 配置安全密钥
const securityCode = import.meta.env.VITE_AMAP_SECURITY_JS_CODE
if (securityCode) {
  ;(window as any)._AMapSecurityConfig = { securityJsCode: securityCode }
}

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
let map: any = null
let AMapObj: any = null
let polygon: any = null
let polyEditor: any = null

// 数据
const points = ref<{ lng: number, lat: number }[]>([])
const isDrawing = ref(false)

// 初始化地图
onMounted(async () => {
  try {
    const amapKey = import.meta.env.VITE_AMAP_KEY
    if (!amapKey) {
      ElMessage.error('未配置高德地图 Key，请在 .env 中设置 VITE_AMAP_KEY')
      return
    }

    AMapObj = await AMapLoader.load({
      key: amapKey,
      version: '2.0',
      plugins: ['AMap.PolygonEditor']
    })

    // 尝试初始化地图
    tryInitMap()
  } catch (e) {
    console.error('AMap load failed:', e)
    ElMessage.error('地图加载失败，请检查网络或 Key 配置')
  }
})

function tryInitMap(attempts = 0) {
  if (!mapContainer.value) return
  
  // 如果容器没有宽度/高度，延迟重试
  if (mapContainer.value.clientWidth === 0 || mapContainer.value.clientHeight === 0) {
    if (attempts < 10) {
      setTimeout(() => tryInitMap(attempts + 1), 200)
    }
    return
  }

  initMap()
}

function initMap() {
  if (!mapContainer.value || !AMapObj || map) return

  map = new AMapObj.Map(mapContainer.value, {
    zoom: 14,
    center: [116.397428, 39.90923],
    viewMode: '2D'
  })

  // 解析初始数据
  if (props.modelValue && props.modelValue.length > 0) {
    points.value = [...props.modelValue]
    const path = props.modelValue.map(p => [p.lng, p.lat])
    drawPolygon(path)
  }
}

// 绘制多边形
function drawPolygon(path: number[][]) {
  if (!map || !AMapObj) return
  
  // 清理旧的
  if (polygon) {
    map.remove(polygon)
    polygon = null
  }
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }

  if (path.length < 3) return

  polygon = new AMapObj.Polygon({
    path: path,
    strokeColor: "#409eff", 
    strokeWeight: 2,
    strokeOpacity: 0.8,
    fillOpacity: 0.2,
    fillColor: '#409eff',
    zIndex: 50,
  })

  map.add(polygon)
  map.setFitView([polygon])
}

// 开始编辑/绘制
function startDraw() {
  if (!map || !AMapObj) {
    ElMessage.warning('地图尚未加载完成')
    return
  }

  // 如果已有编辑实例，先关闭
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }

  // 清除旧的多边形
  if (polygon) {
    map.remove(polygon)
    polygon = null
  }

  points.value = []
  isDrawing.value = true

  // 创建编辑器（不传 target polygon，进入绘制模式）
  polyEditor = new AMapObj.PolygonEditor(map)

  // 监听 add 事件：用户完成绘制（双击闭合）时触发，获取新创建的多边形
  polyEditor.on('add', (event: any) => {
    polygon = event.target
    // 设置多边形样式
    polygon.setOptions({
      strokeColor: '#409eff',
      strokeWeight: 2,
      strokeOpacity: 0.8,
      fillOpacity: 0.2,
      fillColor: '#409eff',
      zIndex: 50,
    })
    updateModel()
  })
  polyEditor.on('addnode', updateModel)
  polyEditor.on('adjust', updateModel)
  polyEditor.on('removenode', updateModel)
  polyEditor.on('end', (event: any) => {
    if (event?.target) {
      polygon = event.target
    }
    updateModel()
  })

  // 开启绘制模式
  polyEditor.open()
  ElMessage.info('请在地图上点击添加顶点，双击结束绘制')
}

// 结束编辑
function endDraw() {
  if (polyEditor) {
    polyEditor.close()
    updateModel()
    polyEditor = null
  }
  isDrawing.value = false

  if (points.value.length < 3) {
    ElMessage.warning('围栏至少需要3个顶点')
  } else {
    ElMessage.success('围栏绘制完成')
  }
}

// 清除
function clearMap() {
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }
  if (polygon && map) {
    map.remove(polygon)
    polygon = null
  }
  points.value = []
  isDrawing.value = false
  emit('update:modelValue', [])
}

// 更新数据模型
function updateModel() {
  // 优先从 polygon 变量获取，其次从编辑器获取
  const target = polygon || polyEditor?.getTarget?.()
  if (!target) return

  const path = target.getPath()
  if (!path || path.length === 0) return

  const newPoints = path.map((p: any) => ({
    lng: typeof p.lng === 'function' ? p.lng() : p.lng,
    lat: typeof p.lat === 'function' ? p.lat() : p.lat
  }))

  points.value = newPoints
  emit('update:modelValue', newPoints)
}

// 监听外部数据变化
watch(() => props.modelValue, (newVal) => {
  // 如果正在编辑中，不响应外部变化
  if (isDrawing.value) return
  
  if (!newVal || newVal.length === 0) {
    if (points.value.length > 0) {
      clearMap()
    }
    return
  }

  // 避免死循环
  if (JSON.stringify(newVal) !== JSON.stringify(points.value)) {
    points.value = [...newVal]
    if (map && AMapObj) {
      const path = newVal.map(p => [p.lng, p.lat])
      drawPolygon(path)
    }
  }
}, { deep: true })

onUnmounted(() => {
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }
  if (map) {
    map.destroy()
    map = null
  }
  AMapObj = null
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