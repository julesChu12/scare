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
              placeholder="姓名/手机号/身份证号"
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
            <el-button type="primary" @click="handleSearch">查询</el-button>
            <el-button @click="handleReset">重置</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 数据表格 -->
      <el-table v-loading="loading" :data="elderlyList" stripe style="width: 100%">
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="age" label="年龄" width="80">
          <template #default="{ row }">
            {{ calculateAge(row.birth_date) }}岁
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="80">
          <template #default="{ row }">
            {{ row.gender === 'male' ? '男' : '女' }}
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="手机号" width="130">
          <template #default="{ row }">
            {{ maskPhone(row.phone) }}
          </template>
        </el-table-column>
        <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
        <el-table-column prop="station_name" label="所属站点" width="180" />
        <el-table-column prop="health_status" label="健康状况" width="100">
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
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="民族" prop="ethnicity">
              <el-input v-model="elderlyForm.ethnicity" placeholder="请输入民族" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="身份证号" prop="id_card">
              <el-input v-model="elderlyForm.id_card" placeholder="请输入身份证号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="elderlyForm.phone" placeholder="请输入手机号" />
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
          <el-col :span="24">
            <el-form-item label="地址" prop="address">
              <el-input v-model="elderlyForm.address" placeholder="请输入详细地址" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="慢性病" prop="chronic_diseases">
              <el-input v-model="elderlyForm.chronic_diseases" type="textarea" :rows="2" placeholder="请输入慢性病情况，多个用逗号分隔" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="特殊需求" prop="special_needs">
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
  station_id: null as number | null,
})

// 站点列表
const stationList = ref<Array<{ id: number; name: string }>>([])

// 老年人列表
const elderlyList = ref<Array<{
  id: number
  name: string
  gender: string
  birth_date: string
  ethnicity: string
  id_card: string
  phone: string
  address: string
  station_id: number
  station_name: string
  health_status: string
  chronic_diseases: string
  special_needs: string
}>>([])

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
  birth_date: null as Date | null,
  ethnicity: '汉族',
  id_card: '',
  phone: '',
  address: '',
  station_id: null as number | null,
  health_status: 'good',
  chronic_diseases: '',
  special_needs: '',
})

// 表单验证规则
const formRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  birth_date: [{ required: true, message: '请选择出生日期', trigger: 'change' }],
  phone: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  address: [{ required: true, message: '请输入地址', trigger: 'blur' }],
  station_id: [{ required: true, message: '请选择所属站点', trigger: 'change' }],
}

/**
 * 加载站点列表
 */
async function loadStations() {
  try {
    // 模拟数据
    stationList.value = [
      { id: 1, name: '霍营街道第一服务站' },
      { id: 2, name: '霍营街道第二服务站' },
      { id: 3, name: '回龙观服务站' },
    ]
  } catch (error) {
    console.error('Failed to load stations:', error)
  }
}

/**
 * 加载老年人列表
 */
async function loadElderlyList() {
  try {
    loading.value = true
    // 模拟数据
    elderlyList.value = [
      {
        id: 1,
        name: '张老人',
        gender: 'male',
        birth_date: '1948-05-15',
        ethnicity: '汉族',
        id_card: '110114194805151234',
        phone: '13800000001',
        address: '霍营街道xx小区1号楼101室',
        station_id: 1,
        station_name: '霍营街道第一服务站',
        health_status: 'good',
        chronic_diseases: '高血压',
        special_needs: '低盐饮食',
      },
      {
        id: 2,
        name: '李老人',
        gender: 'female',
        birth_date: '1944-08-20',
        ethnicity: '汉族',
        id_card: '110114194408201234',
        phone: '13800000002',
        address: '霍营街道xx小区2号楼201室',
        station_id: 1,
        station_name: '霍营街道第一服务站',
        health_status: 'normal',
        chronic_diseases: '糖尿病、高血压',
        special_needs: '低糖低盐饮食',
      },
      {
        id: 3,
        name: '王老人',
        gender: 'male',
        birth_date: '1951-03-10',
        ethnicity: '汉族',
        id_card: '110114195103101234',
        phone: '13800000003',
        address: '霍营街道xx小区3号楼301室',
        station_id: 2,
        station_name: '霍营街道第二服务站',
        health_status: 'good',
        chronic_diseases: '',
        special_needs: '',
      },
      {
        id: 4,
        name: '赵老人',
        gender: 'female',
        birth_date: '1946-11-25',
        ethnicity: '汉族',
        id_card: '110114194611251234',
        phone: '13800000004',
        address: '霍营街道xx小区1号楼501室',
        station_id: 1,
        station_name: '霍营街道第一服务站',
        health_status: 'poor',
        chronic_diseases: '心脏病、高血压、糖尿病',
        special_needs: '需要定期上门护理',
      },
      {
        id: 5,
        name: '刘老人',
        gender: 'male',
        birth_date: '1941-07-08',
        ethnicity: '汉族',
        id_card: '110114194107081234',
        phone: '13800000005',
        address: '霍营街道xx小区2号楼102室',
        station_id: 2,
        station_name: '霍营街道第二服务站',
        health_status: 'normal',
        chronic_diseases: '关节炎',
        special_needs: '',
      },
    ]
    pagination.total = 5
  } catch (error) {
    console.error('Failed to load elderly list:', error)
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
  searchForm.station_id = null
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
function handleEdit(row: typeof elderlyList.value[0]) {
  isEdit.value = true
  dialogTitle.value = '编辑档案'
  Object.assign(elderlyForm, {
    id: row.id,
    name: row.name,
    gender: row.gender,
    birth_date: row.birth_date ? new Date(row.birth_date) : null,
    ethnicity: row.ethnicity,
    id_card: row.id_card,
    phone: row.phone,
    address: row.address,
    station_id: row.station_id,
    health_status: row.health_status,
    chronic_diseases: row.chronic_diseases,
    special_needs: row.special_needs,
  })
  showDialog.value = true
}

/**
 * 查看详情
 */
function handleView(row: typeof elderlyList.value[0]) {
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
      // 模拟保存
      await new Promise(resolve => setTimeout(resolve, 1000))
      ElMessage.success(isEdit.value ? '修改成功' : '创建成功')
      showDialog.value = false
      loadElderlyList()
    } catch (error) {
      console.error('Failed to save:', error)
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
  elderlyForm.birth_date = null
  elderlyForm.ethnicity = '汉族'
  elderlyForm.id_card = ''
  elderlyForm.phone = ''
  elderlyForm.address = ''
  elderlyForm.station_id = null
  elderlyForm.health_status = 'good'
  elderlyForm.chronic_diseases = ''
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
  return textMap[status] || status
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
