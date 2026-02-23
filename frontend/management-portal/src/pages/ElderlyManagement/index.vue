<template>
  <div class="elderly-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <h3>老年人档案</h3>
            <p class="subtitle">管理辖区内老年人基本信息和服务记录</p>
          </div>
          <el-button type="primary" :icon="Plus" @click="handleCreate">
            新建档案
          </el-button>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="search-bar">
        <el-form :inline="true" :model="searchForm">
          <el-form-item>
            <el-input
              v-model="searchForm.keyword"
              placeholder="姓名/手机号"
              :prefix-icon="Search"
              clearable
              style="width: 220px"
              @keyup.enter="handleSearch"
            />
          </el-form-item>
          <el-form-item>
            <el-select v-model="searchForm.station_id" placeholder="全部站点" clearable style="width: 180px">
              <el-option
                v-for="station in stationList"
                :key="station.id"
                :label="station.name"
                :value="station.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-select v-model="searchForm.health_status" placeholder="健康状况" clearable style="width: 140px">
              <el-option label="良好" value="good" />
              <el-option label="一般" value="normal" />
              <el-option label="较差" value="poor" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 数据表格 -->
      <el-table v-loading="loading" :data="elderlyList" stripe style="width: 100%">
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column label="年龄" width="80">
          <template #default="{ row }">
            {{ row.birth_date ? calculateAge(row.birth_date) + '岁' : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="性别" width="80">
          <template #default="{ row }">
            {{ row.gender === 'male' ? '男' : row.gender === 'female' ? '女' : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="手机号" width="130">
          <template #default="{ row }">
            {{ maskPhone(row.phone) }}
          </template>
        </el-table-column>
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="station_name" label="所属站点" width="180" />
        <el-table-column label="健康状况" width="100">
          <template #default="{ row }">
            <el-tag :type="getHealthStatusType(row.health_status)" size="small">
              {{ getHealthStatusText(row.health_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              查看
            </el-button>
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 新建/编辑弹窗 -->
    <el-dialog v-model="showDialog" :title="dialogTitle" width="800px">
      <el-form ref="formRef" :model="elderlyForm" :rules="formRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓名" prop="name">
              <el-input v-model="elderlyForm.name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="性别" prop="gender">
              <el-radio-group v-model="elderlyForm.gender">
                <el-radio value="male">男</el-radio>
                <el-radio value="female">女</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="出生日期" prop="birth_date">
              <el-date-picker
                v-model="elderlyForm.birth_date"
                type="date"
                placeholder="选择出生日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="elderlyForm.phone" placeholder="请输入手机号" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="身份证号" prop="id_card">
              <el-input v-model="elderlyForm.id_card" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属站点" prop="station_id">
              <el-select v-model="elderlyForm.station_id" placeholder="请选择站点" style="width: 100%">
                <el-option
                  v-for="station in stationList"
                  :key="station.id"
                  :label="station.name"
                  :value="station.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="健康状况" prop="health_status">
              <el-select v-model="elderlyForm.health_status" placeholder="请选择健康状况" style="width: 100%">
                <el-option label="良好" value="good" />
                <el-option label="一般" value="normal" />
                <el-option label="较差" value="poor" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="失能等级">
              <el-select v-model="elderlyForm.disability_level" placeholder="请选择" clearable style="width: 100%">
                <el-option label="自理" value="自理" />
                <el-option label="轻度" value="轻度" />
                <el-option label="中度" value="中度" />
                <el-option label="重度" value="重度" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="地址" prop="address">
              <el-input v-model="elderlyForm.address" placeholder="请输入详细地址" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="病史">
              <el-input v-model="elderlyForm.medical_history" type="textarea" :rows="2" placeholder="请输入病史/慢性病情况" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="特殊需求">
              <el-input v-model="elderlyForm.special_needs" type="textarea" :rows="2" placeholder="请输入特殊需求" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'
import { elderlyApi, stationApi } from '@/api'
import type { ElderlyProfile } from '@/types/api'

const router = useRouter()

// 加载状态
const loading = ref(false)
const saving = ref(false)

// 弹窗控制
const showDialog = ref(false)
const dialogTitle = ref('新建档案')
const isEdit = ref(false)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  station_id: undefined as number | undefined,
  health_status: undefined as string | undefined,
})

// 站点列表
const stationList = ref<Array<{ id: number; name: string }>>([])

// 老年人列表
const elderlyList = ref<ElderlyProfile[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
})

// 表单
const formRef = ref<FormInstance>()
const elderlyForm = reactive({
  id: null as number | null,
  name: '',
  gender: 'male',
  birth_date: '',
  id_card: '',
  phone: '',
  address: '',
  station_id: undefined as number | undefined,
  health_status: 'good',
  disability_level: '',
  medical_history: '',
  special_needs: '',
})

// 表单验证规则
const formRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  station_id: [{ required: true, message: '请选择所属站点', trigger: 'change' }],
}

/**
 * 加载站点列表
 */
async function loadStations() {
  try {
    const res = await stationApi.getStations({ page: 1, page_size: 100 })
    stationList.value = (res.data?.items || []).map((s) => ({ id: s.id, name: s.name }))
  } catch (error) {
    console.error('加载站点失败:', error)
  }
}

/**
 * 加载老年人列表
 */
async function loadElderlyList() {
  try {
    loading.value = true
    const res = await elderlyApi.getList({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      station_id: searchForm.station_id,
      health_status: searchForm.health_status,
    })
    elderlyList.value = res.data?.items || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error('加载老人列表失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

/**
 * 搜索
 */
function handleSearch() {
  pagination.page = 1
  loadElderlyList()
}

/**
 * 重置搜索
 */
function handleReset() {
  searchForm.keyword = ''
  searchForm.station_id = undefined
  searchForm.health_status = undefined
  pagination.page = 1
  loadElderlyList()
}

/**
 * 新建档案
 */
function handleCreate() {
  isEdit.value = false
  dialogTitle.value = '新建档案'
  resetForm()
  showDialog.value = true
}

/**
 * 编辑档案
 */
function handleEdit(row: ElderlyProfile) {
  isEdit.value = true
  dialogTitle.value = '编辑档案'
  Object.assign(elderlyForm, {
    id: row.id,
    name: row.name,
    gender: row.gender,
    birth_date: row.birth_date || '',
    id_card: row.id_card || '',
    phone: row.phone,
    address: row.address || '',
    station_id: row.station_id || undefined,
    health_status: row.health_status || 'good',
    disability_level: row.disability_level || '',
    medical_history: row.medical_history || '',
    special_needs: row.special_needs || '',
  })
  showDialog.value = true
}

/**
 * 查看详情
 */
function handleView(row: ElderlyProfile) {
  router.push(`/residents/elderly/${row.id}`)
}

/**
 * 保存
 */
async function handleSave() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      saving.value = true
      const data = {
        name: elderlyForm.name,
        gender: elderlyForm.gender,
        birth_date: elderlyForm.birth_date || undefined,
        id_card: elderlyForm.id_card || undefined,
        address: elderlyForm.address || undefined,
        station_id: elderlyForm.station_id,
        health_status: elderlyForm.health_status || undefined,
        disability_level: elderlyForm.disability_level || undefined,
        medical_history: elderlyForm.medical_history || undefined,
        special_needs: elderlyForm.special_needs || undefined,
      }

      if (isEdit.value && elderlyForm.id) {
        await elderlyApi.update(elderlyForm.id, data)
        ElMessage.success('修改成功')
      } else {
        await elderlyApi.create({ ...data, phone: elderlyForm.phone })
        ElMessage.success('创建成功')
      }
      showDialog.value = false
      loadElderlyList()
    } catch (error: any) {
      const msg = error?.response?.data?.msg || '操作失败'
      ElMessage.error(msg)
    } finally {
      saving.value = false
    }
  })
}

/**
 * 重置表单
 */
function resetForm() {
  elderlyForm.id = null
  elderlyForm.name = ''
  elderlyForm.gender = 'male'
  elderlyForm.birth_date = ''
  elderlyForm.id_card = ''
  elderlyForm.phone = ''
  elderlyForm.address = ''
  elderlyForm.station_id = undefined
  elderlyForm.health_status = 'good'
  elderlyForm.disability_level = ''
  elderlyForm.medical_history = ''
  elderlyForm.special_needs = ''
}

/**
 * 分页大小变化
 */
function handleSizeChange(size: number) {
  pagination.pageSize = size
  pagination.page = 1
  loadElderlyList()
}

/**
 * 页码变化
 */
function handlePageChange(page: number) {
  pagination.page = page
  loadElderlyList()
}

/**
 * 计算年龄
 */
function calculateAge(birthDate: string): number {
  if (!birthDate) return 0
  return dayjs().diff(dayjs(birthDate), 'year')
}

/**
 * 手机号脱敏
 */
function maskPhone(phone: string): string {
  if (!phone || phone.length < 7) return phone
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 获取健康状况标签类型
 */
function getHealthStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    good: 'success',
    normal: 'warning',
    poor: 'danger',
  }
  return typeMap[status] || ''
}

/**
 * 获取健康状况文本
 */
function getHealthStatusText(status: string): string {
  const textMap: Record<string, string> = {
    good: '良好',
    normal: '一般',
    poor: '较差',
  }
  return textMap[status] || status || '-'
}

onMounted(() => {
  loadStations()
  loadElderlyList()
})
</script>

<style scoped lang="scss">
.elderly-management {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0 0 5px 0;
      font-size: 18px;
      font-weight: 500;
    }

    .subtitle {
      margin: 0;
      font-size: 13px;
      color: #909399;
    }
  }

  .search-bar {
    margin-bottom: 16px;
  }

  .pagination-container {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
