<template>
  <div class="page-container">
    <div class="page-header">
      <h2>围栏管理</h2>
      <el-button type="primary" @click="handleAdd">新增围栏</el-button>
    </div>

    <!-- 列表 -->
    <el-table v-loading="loading" :data="tableData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="围栏名称" min-width="150" />
      <el-table-column prop="station_id" label="所属站点ID" width="120" />
      <el-table-column prop="priority" label="优先级" width="100" />
      <el-table-column label="顶点数" width="100">
        <template #default="{ row }">
          {{ parsePoints(row.points)?.length || 0 }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增围栏' : '编辑围栏'"
      width="1000px"
      top="5vh"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="围栏名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入围栏名称" />
        </el-form-item>
        <el-form-item label="所属站点" prop="station_id">
          <el-input-number v-model="formData.station_id" :min="1" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-input-number v-model="formData.priority" :min="0" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="formData.status"
            active-value="active"
            inactive-value="inactive"
            active-text="启用"
            inactive-text="停用"
          />
        </el-form-item>

        <el-form-item label="服务范围" prop="points">
          <map-polygon-editor v-model="formData.points" />
          <div class="form-tip">请在地图上绘制围栏区域</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { zoneApi } from '@/api'
import type { Zone, ZonePoint } from '@/types/api'
import MapPolygonEditor from '@/components/MapPolygonEditor.vue'

// 数据
const loading = ref(false)
const tableData = ref<Zone[]>([])

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const formData = reactive({
  id: 0,
  name: '',
  station_id: 1,
  priority: 0,
  status: 'active',
  points: [] as ZonePoint[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入围栏名称', trigger: 'blur' }]
}

// 方法
const fetchData = async () => {
  loading.value = true
  try {
    const res = await zoneApi.getZones()
    tableData.value = res.data.items
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  dialogType.value = 'add'
  dialogVisible.value = true
  // Reset form
  formData.id = 0
  formData.name = ''
  formData.priority = 0
  formData.status = 'active'
  formData.points = []
}

const handleEdit = (row: Zone) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  
  const points = parsePoints(row.points)

  Object.assign(formData, {
    id: row.id,
    name: row.name,
    station_id: row.station_id,
    priority: row.priority,
    status: row.status,
    points: points
  })
}

function parsePoints(pointsData: string | ZonePoint[]): ZonePoint[] {
  try {
    if (typeof pointsData === 'string') {
      return JSON.parse(pointsData)
    } else if (Array.isArray(pointsData)) {
      return pointsData
    }
  } catch (e) {
    console.error('Failed to parse zone points:', e)
  }
  return []
}

const handleDelete = (row: Zone) => {
  ElMessageBox.confirm('确定删除该围栏吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await zoneApi.deleteZone(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      if (!formData.points || formData.points.length < 3) {
        ElMessage.warning('请绘制至少包含3个点的围栏')
        return
      }

      submitting.value = true
      try {
        if (dialogType.value === 'add') {
          await zoneApi.createZone({
            name: formData.name,
            station_id: formData.station_id,
            priority: formData.priority,
            status: formData.status,
            points: formData.points
          })
        } else {
          await zoneApi.updateZone(formData.id, {
            name: formData.name,
            station_id: formData.station_id,
            priority: formData.priority,
            status: formData.status,
            points: formData.points
          })
        }
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '编辑成功')
        dialogVisible.value = false
        fetchData()
      } catch (error) {
        console.error(error)
      } finally {
        submitting.value = false
      }
    }
  })
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
  background: white;
  min-height: 100%;
}
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
