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
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
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
// 使用普通变量而非 ref/shallowRef 存储 AMap 实例，避免 Vue 代理导致的类型检查错误
let map: any = null
let AMapObj: any = null
let polygon: any = null
let polyEditor: any = null

// 数据
const points = ref<{ lng: number, lat: number }[]>([])

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

    // 尝试初始化地图，如果容器不可见（例如在弹窗动画中），则轮询直到可见
    tryInitMap()
  } catch (e) {
    console.error('AMap load failed:', e)
    ElMessage.error('地图加载失败，请检查网络或 Key 配置')
  }
})

function tryInitMap(attempts = 0) {
  if (!mapContainer.value) return
  
  // 如果容器没有宽度/高度，说明可能还在弹窗动画中，延迟重试
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
    center: [116.397428, 39.90923], // 默认北京
    viewMode: '2D'
  })

  // 解析初始数据
  if (props.modelValue && props.modelValue.length > 0) {
    points.value = props.modelValue
    const path = props.modelValue.map(p => [p.lng, p.lat])
    drawPolygon(path)
  }

  // 点击地图开始绘制（如果当前没有多边形）
  map.on('click', () => {
    if (!polygon && !polyEditor) {
      startDraw()
    }
  })
}

// 绘制/回显多边形
function drawPolygon(path: number[][]) {
  // 清理旧的
  if (polygon) {
    map.remove(polygon)
    if (polyEditor) {
      polyEditor.close()
      polyEditor = null
    }
  }

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
  if (!map || !AMapObj) return

  // 如果已有编辑实例，先关闭
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }

  // 如果没有多边形对象，创建一个空的
  if (!polygon) {
    polygon = new AMapObj.Polygon({
      path: [],
      strokeColor: "#409eff", 
      strokeWeight: 2,
      strokeOpacity: 0.8,
      fillOpacity: 0.2,
      fillColor: '#409eff',
      zIndex: 50,
    })
    map.add(polygon)
  }

  polyEditor = new AMapObj.PolygonEditor(map, polygon)
  
  // 监听编辑事件更新坐标
  polyEditor.on('addnode', updateModel)
  polyEditor.on('adjust', updateModel)
  polyEditor.on('removenode', updateModel)
  polyEditor.on('end', updateModel)

  polyEditor.open()
  
  if (!points.value.length) {
    ElMessage.info('请在地图上点击添加顶点绘制围栏')
  }
}

// 结束编辑
function endDraw() {
  if (polyEditor) {
    polyEditor.close()
    updateModel() // 确保最后状态同步
    polyEditor = null
  }
}

// 清除
function clearMap() {
  if (polyEditor) {
    polyEditor.close()
    polyEditor = null
  }
  if (polygon) {
    map.remove(polygon)
    polygon = null
  }
  points.value = []
  emit('update:modelValue', [])
}

// 更新数据模型
function updateModel() {
  if (!polygon) return
  
  const path = polygon.getPath()
  // path 是 AMap.LngLat 对象的数组
  const coords = path.map((p: any) => [p.lng, p.lat])
  const newPoints = coords.map((p: any) => ({ lng: p[0], lat: p[1] }))
  
  points.value = newPoints
  emit('update:modelValue', newPoints)
}

// 监听外部数据变化
watch(() => props.modelValue, (newVal) => {
  // 如果正在编辑中，不响应外部变化以免冲突，除非外部清空了
  if (!newVal || newVal.length === 0) {
    if (points.value.length > 0) {
      clearMap()
    }
    return
  }

  // 简单比对是否需要更新（避免死循环）
  if (JSON.stringify(newVal) !== JSON.stringify(points.value)) {
    points.value = newVal
    if (map && AMapObj) {
      const path = newVal.map(p => [p.lng, p.lat])
      drawPolygon(path)
    }
  }
}, { deep: true })

onUnmounted(() => {
  if (polyEditor) {
    polyEditor.close()
  }
  if (map) {
    map.destroy()
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