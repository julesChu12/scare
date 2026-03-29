<template>
  <div class="page-container">
    <div class="page-header">
      <h2>站点管理</h2>
      <el-button v-if="canCreateStation" type="primary" @click="handleAdd">新增站点</el-button>
    </div>

    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="站点名称">
          <el-input v-model="searchForm.name" placeholder="请输入站点名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 列表 -->
    <el-table v-loading="loading" :data="tableData" border style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="站点名称" min-width="150" />
      <el-table-column prop="code" label="编号" width="120" />
      <el-table-column prop="address" label="地址" min-width="200" show-overflow-tooltip />
      <el-table-column prop="phone" label="联系电话" width="150" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="280" fixed="right">
        <template #default="{ row }">
          <el-button link type="success" @click="handleShowQrCode(row)">站点二维码</el-button>
          <el-button v-if="canUpdateStation" link type="primary" @click="handleEdit(row)">编辑</el-button>
          <el-button v-if="canDeleteStation" link type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="fetchData"
      />
    </div>

    <!-- 弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增站点' : '编辑站点'"
      width="700px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="站点名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入站点名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="站点编号" prop="code">
              <el-input v-model="formData.code" placeholder="请输入站点编号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入联系电话" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-switch
                v-model="formData.status"
                active-value="active"
                inactive-value="inactive"
                active-text="启用"
                inactive-text="停用"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="地址" prop="address">
              <el-input v-model="formData.address" placeholder="请输入地址" type="textarea" :rows="2" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="qrDialogVisible"
      title="站点二维码"
      width="420px"
      destroy-on-close
    >
      <div class="qr-dialog-content">
        <div class="qr-station-name">{{ qrStationName }}</div>
        <div class="qr-hint">扫码后进入 C 端快速开通页，来源站点会以参数形式附带。</div>

        <div v-loading="qrLoading" class="qr-preview-wrapper">
          <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="站点二维码" class="qr-preview" />
        </div>

        <el-input
          :model-value="qrCodeUrl"
          type="textarea"
          :rows="3"
          readonly
        />
      </div>

      <template #footer>
        <el-button @click="qrDialogVisible = false">关闭</el-button>
        <el-button :disabled="!qrCodeUrl" @click="handleCopyQrLink">复制链接</el-button>
        <el-button type="primary" :disabled="!qrCodeDataUrl" @click="handleDownloadQrCode">下载二维码</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import QRCode from 'qrcode'
import { stationApi } from '@/api'
import { PERM_STATION_CREATE, PERM_STATION_DELETE, PERM_STATION_UPDATE } from '@/constants/permissions'
import { useAuthStore } from '@/store/modules/auth'
import type { Station } from '@/types/api'

// 数据
const authStore = useAuthStore()
const loading = ref(false)
const tableData = ref<Station[]>([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const searchForm = reactive({
  name: ''
})

// 弹窗
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const submitting = ref(false)
const formRef = ref<FormInstance>()
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  address: '',
  phone: '',
  status: 'active'
})

const qrDialogVisible = ref(false)
const qrLoading = ref(false)
const qrStationName = ref('')
const qrCodeUrl = ref('')
const qrCodeDataUrl = ref('')

const rules: FormRules = {
  name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入站点编号', trigger: 'blur' }]
}

const canCreateStation = computed(() => authStore.hasPermission(PERM_STATION_CREATE))
const canUpdateStation = computed(() => authStore.hasPermission(PERM_STATION_UPDATE))
const canDeleteStation = computed(() => authStore.hasPermission(PERM_STATION_DELETE))

const resolveCEndBaseUrl = () => {
  const configuredBaseUrl = import.meta.env.VITE_C_END_BASE_URL?.trim()
  if (configuredBaseUrl) {
    return configuredBaseUrl
  }

  if (import.meta.env.DEV) {
    return `${window.location.protocol}//${window.location.hostname}:5174`
  }

  return window.location.origin
}

const buildQuickStartUrl = (stationId: number) => {
  const url = new URL('/quick', resolveCEndBaseUrl())
  url.searchParams.set('source_station_id', String(stationId))
  return url.toString()
}

const copyText = async (text: string) => {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }

  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'true')
  textarea.style.position = 'fixed'
  textarea.style.top = '-9999px'
  document.body.appendChild(textarea)
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
}

// 方法
const fetchData = async () => {
  loading.value = true
  try {
    const res = await stationApi.getStations({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: searchForm.name || undefined
    })
    tableData.value = res.data.items
    pagination.total = res.data.total
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const resetSearch = () => {
  searchForm.name = ''
  handleSearch()
}

const handleAdd = () => {
  dialogType.value = 'add'
  dialogVisible.value = true
  // Reset form
  formData.id = 0
  formData.name = ''
  formData.code = ''
  formData.address = ''
  formData.phone = ''
  formData.status = 'active'
}

const handleEdit = (row: Station) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  Object.assign(formData, {
    id: row.id,
    name: row.name,
    code: row.code,
    address: row.address,
    phone: row.phone,
    status: row.status
  })
}

const handleDelete = (row: Station) => {
  ElMessageBox.confirm('确定删除该站点吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await stationApi.deleteStation(row.id)
      ElMessage.success('删除成功')
      fetchData()
    } catch (error) {
      console.error(error)
    }
  })
}

const handleShowQrCode = async (row: Station) => {
  qrDialogVisible.value = true
  qrLoading.value = true
  qrStationName.value = row.name
  qrCodeDataUrl.value = ''
  qrCodeUrl.value = buildQuickStartUrl(row.id)

  try {
    qrCodeDataUrl.value = await QRCode.toDataURL(qrCodeUrl.value, {
      width: 320,
      margin: 2,
      errorCorrectionLevel: 'M',
      color: {
        dark: '#111827',
        light: '#ffffff',
      },
    })
  } catch (error) {
    console.error(error)
    ElMessage.error('二维码生成失败')
  } finally {
    qrLoading.value = false
  }
}

const handleCopyQrLink = async () => {
  if (!qrCodeUrl.value) return

  try {
    await copyText(qrCodeUrl.value)
    ElMessage.success('链接已复制')
  } catch (error) {
    console.error(error)
    ElMessage.error('复制失败，请手动复制')
  }
}

const handleDownloadQrCode = () => {
  if (!qrCodeDataUrl.value) return

  const link = document.createElement('a')
  link.href = qrCodeDataUrl.value
  link.download = `station-${qrStationName.value || 'qrcode'}.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (dialogType.value === 'add') {
          await stationApi.createStation({
            name: formData.name,
            code: formData.code,
            address: formData.address,
            phone: formData.phone,
            status: formData.status
          })
        } else {
          await stationApi.updateStation(formData.id, {
            name: formData.name,
            code: formData.code,
            address: formData.address,
            phone: formData.phone,
            status: formData.status
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

const resetForm = () => {
  if (formRef.value) formRef.value.resetFields()
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
.search-bar {
  margin-bottom: 20px;
}
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.qr-dialog-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.qr-station-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.qr-hint {
  color: #6b7280;
  line-height: 1.6;
}

.qr-preview-wrapper {
  min-height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
}

.qr-preview {
  width: 320px;
  height: 320px;
  object-fit: contain;
}
</style>
