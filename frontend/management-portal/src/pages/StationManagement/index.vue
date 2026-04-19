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
      width="800px"
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
          <el-col :span="12">
            <el-form-item label="纬度">
              <el-input :model-value="locationLatitudeText" readonly placeholder="暂无定位坐标" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="经度">
              <el-input :model-value="locationLongitudeText" readonly placeholder="暂无定位坐标" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <div class="location-panel">
              <div class="location-panel__header">
                <div>
                  <div class="location-panel__title">当前位置</div>
                  <div class="location-panel__meta">
                    {{ hasStationLocation ? '点击地图或拖拽标记可微调当前站点坐标' : '请输入地址后定位，或直接在地图上选择站点位置' }}
                  </div>
                </div>
                <el-button
                  type="primary"
                  plain
                  :loading="locatingByAddress"
                  :disabled="!formData.address.trim()"
                  @click="handleLocateByAddress"
                >
                  解析地址定位
                </el-button>
              </div>
              <div class="location-panel__address">
                {{ formData.address || '未填写站点地址' }}
              </div>
              <div class="location-panel__map">
                <map-location-editor
                  ref="locationEditorRef"
                  :longitude="formData.longitude"
                  :latitude="formData.latitude"
                  :zoom="16"
                  @update:address="formData.address = $event"
                  @update:longitude="formData.longitude = $event"
                  @update:latitude="formData.latitude = $event"
                />
              </div>
            </div>
          </el-col>
        </el-row>

        <!-- 站点工作人员（仅编辑时显示） -->
        <el-row v-if="dialogType === 'edit'">
          <el-col :span="24">
            <div class="staff-section">
              <div class="staff-title">站点工作人员</div>
              <el-table :data="stationStaffs" v-loading="staffLoading" stripe size="small">
                <el-table-column prop="name" label="姓名" width="100" />
                <el-table-column prop="phone" label="手机号" width="140" />
                <el-table-column label="B端角色" min-width="120">
                  <template #default="{ row }">
                    <template v-if="row.b_end_identities?.length">
                      <el-tag
                        v-for="identity in row.b_end_identities"
                        :key="identity"
                        size="small"
                        style="margin-right: 4px"
                      >
                        {{ getIdentityName(identity) }}
                      </el-tag>
                    </template>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
                <el-table-column label="状态" width="80">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
                      {{ row.status === 'active' ? '正常' : '停用' }}
                    </el-tag>
                  </template>
                </el-table-column>
              </el-table>
              <div v-if="!staffLoading && stationStaffs.length === 0" class="staff-empty">
                暂无工作人员
              </div>
            </div>
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
import { stationApi, userApi } from '@/api'
import MapLocationEditor from '@/components/MapLocationEditor.vue'
import { PERM_STATION_CREATE, PERM_STATION_DELETE, PERM_STATION_UPDATE } from '@/constants/permissions'
import { useAuthStore } from '@/store/modules/auth'
import type { Station, User } from '@/types/api'

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
const locationEditorRef = ref<{
  geocodeAddress: (address: string) => Promise<{ formattedAddress?: string }>
} | null>(null)
const locatingByAddress = ref(false)
const formData = reactive({
  id: 0,
  name: '',
  code: '',
  address: '',
  phone: '',
  latitude: undefined as number | undefined,
  longitude: undefined as number | undefined,
  status: 'active'
})

const qrDialogVisible = ref(false)
const qrLoading = ref(false)
const qrStationName = ref('')
const qrCodeUrl = ref('')
const qrCodeDataUrl = ref('')

// 站点工作人员
const stationStaffs = ref<User[]>([])
const staffLoading = ref(false)

// 身份名称映射
const identityMap: Record<string, string> = {
  admin: '系统管理员',
  station_manager: '站点管理员',
  staff: '工作人员',
}

function getIdentityName(identity: string) {
  return identityMap[identity] || identity
}

// 加载站点工作人员
async function loadStationStaffs(stationId: number) {
  staffLoading.value = true
  try {
    const res = await userApi.getUsers({ station_id: stationId, page: 1, page_size: 100 })
    if (res.msg === 'ok') {
      stationStaffs.value = res.data.items || []
    }
  } catch (error) {
    console.error('加载站点工作人员失败', error)
    stationStaffs.value = []
  } finally {
    staffLoading.value = false
  }
}

const rules: FormRules = {
  name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入站点编号', trigger: 'blur' }]
}

const canCreateStation = computed(() => authStore.hasPermission(PERM_STATION_CREATE))
const canUpdateStation = computed(() => authStore.hasPermission(PERM_STATION_UPDATE))
const canDeleteStation = computed(() => authStore.hasPermission(PERM_STATION_DELETE))
const hasStationLocation = computed(() => (
  typeof formData.latitude === 'number'
  && Number.isFinite(formData.latitude)
  && typeof formData.longitude === 'number'
  && Number.isFinite(formData.longitude)
))
const locationLatitudeText = computed(() => (
  hasStationLocation.value ? formData.latitude!.toFixed(6) : ''
))
const locationLongitudeText = computed(() => (
  hasStationLocation.value ? formData.longitude!.toFixed(6) : ''
))

const resolveCEndBaseUrl = () => {
  const configuredBaseUrl = import.meta.env.VITE_C_END_BASE_URL?.trim()
  if (configuredBaseUrl) {
    return configuredBaseUrl
  }

  if (import.meta.env.DEV) {
    return `https://${window.location.hostname}:5174`
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
  formData.latitude = undefined
  formData.longitude = undefined
  formData.status = 'active'
  stationStaffs.value = []
}

const handleEdit = async (row: Station) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  try {
    const res = await stationApi.getStation(row.id)
    const station = res.data
    Object.assign(formData, {
      id: station.id,
      name: station.name,
      code: station.code,
      address: station.address,
      phone: station.phone,
      latitude: station.latitude,
      longitude: station.longitude,
      status: station.status
    })
  } catch (error) {
    console.error('加载站点详情失败', error)
    Object.assign(formData, {
      id: row.id,
      name: row.name,
      code: row.code,
      address: row.address,
      phone: row.phone,
      latitude: row.latitude,
      longitude: row.longitude,
      status: row.status
    })
  }
  // 加载该站点的工作人员
  await loadStationStaffs(row.id)
}

const handleLocateByAddress = async () => {
  if (!formData.address.trim()) {
    ElMessage.warning('请先输入站点地址')
    return
  }

  if (!locationEditorRef.value?.geocodeAddress) {
    ElMessage.warning('地图尚未加载完成，请稍后再试')
    return
  }

  locatingByAddress.value = true
  try {
    const result = await locationEditorRef.value.geocodeAddress(formData.address)
    if (typeof result?.formattedAddress === 'string' && result.formattedAddress.trim()) {
      formData.address = result.formattedAddress.trim()
    }
    ElMessage.success('已根据地址定位，请确认地图落点')
  } catch (error) {
    console.error('地址定位失败', error)
    ElMessage.error('地址解析失败，请检查地址或手动在地图上选择')
  } finally {
    locatingByAddress.value = false
  }
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
            latitude: formData.latitude,
            longitude: formData.longitude,
            status: formData.status
          })
        } else {
          await stationApi.updateStation(formData.id, {
            name: formData.name,
            code: formData.code,
            address: formData.address,
            phone: formData.phone,
            latitude: formData.latitude,
            longitude: formData.longitude,
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

.staff-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.staff-title {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  margin-bottom: 12px;
}

.staff-empty {
  text-align: center;
  color: #909399;
  padding: 20px 0;
  font-size: 14px;
}

.location-panel {
  margin-top: 4px;
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #f9fafb;
}

.location-panel__header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 8px;
}

.location-panel__title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.location-panel__meta {
  font-size: 12px;
  color: #6b7280;
}

.location-panel__address {
  margin-bottom: 12px;
  color: #374151;
  line-height: 1.6;
}

.location-panel__map {
  height: 320px;
  overflow: hidden;
  border-radius: 10px;
  border: 1px solid #dbeafe;
  background: #fff;
}
</style>
