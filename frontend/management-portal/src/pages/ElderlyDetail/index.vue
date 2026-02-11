<template>
  <div class="elderly-detail">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
      <h2>老年人档案详情</h2>
      <el-button type="primary" @click="handleEdit">编辑</el-button>
    </div>

    <!-- 基本信息 -->
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
        </div>
      </template>
      <div class="basic-info">
        <el-avatar :size="100" :src="elderlyInfo.avatar">
          <el-icon :size="50"><User /></el-icon>
        </el-avatar>
        <el-descriptions :column="2" border class="info-descriptions">
          <el-descriptions-item label="姓名">{{ elderlyInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="性别">{{ elderlyInfo.gender === 'male' ? '男' : '女' }}</el-descriptions-item>
          <el-descriptions-item label="年龄">{{ calculateAge(elderlyInfo.birth_date) }}岁</el-descriptions-item>
          <el-descriptions-item label="民族">{{ elderlyInfo.ethnicity }}</el-descriptions-item>
          <el-descriptions-item label="身份证号">{{ maskIdCard(elderlyInfo.id_card) }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ elderlyInfo.phone }}</el-descriptions-item>
          <el-descriptions-item label="地址" :span="2">{{ elderlyInfo.address }}</el-descriptions-item>
          <el-descriptions-item label="所属站点">{{ elderlyInfo.station_name }}</el-descriptions-item>
          <el-descriptions-item label="建档时间">{{ formatDate(elderlyInfo.created_at) }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>

    <!-- 健康信息 -->
    <el-card class="health-card">
      <template #header>
        <div class="card-header">
          <span>健康信息</span>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="健康状况">
          <el-tag :type="getHealthStatusType(elderlyInfo.health_status)" size="small">
            {{ getHealthStatusText(elderlyInfo.health_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="慢性病">{{ elderlyInfo.chronic_diseases || '无' }}</el-descriptions-item>
        <el-descriptions-item label="过敏史">{{ elderlyInfo.allergies || '无' }}</el-descriptions-item>
        <el-descriptions-item label="特殊需求">{{ elderlyInfo.special_needs || '无' }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 紧急联系人 -->
    <el-card class="contacts-card">
      <template #header>
        <div class="card-header">
          <span>紧急联系人</span>
          <el-button type="primary" link @click="handleAddContact">
            <el-icon><Plus /></el-icon>
            添加
          </el-button>
        </div>
      </template>
      <el-table :data="emergencyContacts" stripe style="width: 100%" empty-text="暂无紧急联系人">
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="relationship" label="关系" width="100" />
        <el-table-column prop="phone" label="电话" width="150" />
        <el-table-column prop="is_primary" label="主要联系人" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.is_primary" type="success" size="small">是</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEditContact(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDeleteContact(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 服务记录 -->
    <el-card class="records-card">
      <template #header>
        <div class="card-header">
          <span>服务记录</span>
          <el-button type="primary" link @click="handleViewAllRecords">
            更多 >
          </el-button>
        </div>
      </template>
      <el-table :data="serviceRecords" stripe style="width: 100%" empty-text="暂无服务记录">
        <el-table-column prop="service_date" label="日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.service_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="service_type" label="服务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getServiceTypeTag(row.service_type)" size="small">
              {{ getServiceTypeText(row.service_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="staff_name" label="服务人员" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <span v-if="row.rating" class="rating">
              <el-icon><Star /></el-icon>
              {{ row.rating }}
            </span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
      </el-table>
    </el-card>

    <!-- 添加/编辑联系人弹窗 -->
    <el-dialog v-model="showContactDialog" :title="contactDialogTitle" width="500px">
      <el-form ref="contactFormRef" :model="contactForm" :rules="contactRules" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="contactForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="关系" prop="relationship">
          <el-select v-model="contactForm.relationship" placeholder="请选择关系" style="width: 100%">
            <el-option label="儿子" value="儿子" />
            <el-option label="女儿" value="女儿" />
            <el-option label="配偶" value="配偶" />
            <el-option label="兄弟" value="兄弟" />
            <el-option label="姐妹" value="姐妹" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="contactForm.phone" placeholder="请输入电话" />
        </el-form-item>
        <el-form-item label="主要联系人">
          <el-switch v-model="contactForm.is_primary" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showContactDialog = false">取消</el-button>
        <el-button type="primary" :loading="savingContact" @click="handleSaveContact">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowLeft, User, Plus, Star } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

// 老年人ID
const elderlyId = ref<number>(Number(route.params.id))

// 老年人信息
const elderlyInfo = ref({
  id: 0,
  name: '',
  gender: 'male',
  birth_date: '',
  ethnicity: '',
  id_card: '',
  phone: '',
  address: '',
  station_id: 0,
  station_name: '',
  health_status: 'good',
  chronic_diseases: '',
  allergies: '',
  special_needs: '',
  avatar: '',
  created_at: '',
})

// 紧急联系人
const emergencyContacts = ref<Array<{
  id: number
  name: string
  relationship: string
  phone: string
  is_primary: boolean
}>>([])

// 服务记录
const serviceRecords = ref<Array<{
  id: number
  service_date: string
  service_type: string
  staff_name: string
  status: string
  rating: number | null
  remark: string
}>>([])

// 联系人弹窗
const showContactDialog = ref(false)
const contactDialogTitle = ref('添加联系人')
const savingContact = ref(false)
const isEditContact = ref(false)

const contactFormRef = ref<FormInstance>()
const contactForm = reactive({
  id: null as number | null,
  name: '',
  relationship: '',
  phone: '',
  is_primary: false,
})

const contactRules: FormRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  relationship: [{ required: true, message: '请选择关系', trigger: 'change' }],
  phone: [{ required: true, message: '请输入电话', trigger: 'blur' }],
}

/**
 * 加载老年人详情
 */
async function loadElderlyDetail() {
  try {
    // 模拟数据
    elderlyInfo.value = {
      id: elderlyId.value,
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
      chronic_diseases: '高血压、糖尿病',
      allergies: '无',
      special_needs: '低盐低糖饮食',
      avatar: '',
      created_at: '2025-06-01',
    }
  } catch (error) {
    console.error('Failed to load elderly detail:', error)
  }
}

/**
 * 加载紧急联系人
 */
async function loadEmergencyContacts() {
  try {
    // 模拟数据
    emergencyContacts.value = [
      { id: 1, name: '张小明', relationship: '儿子', phone: '13900000001', is_primary: true },
      { id: 2, name: '张小红', relationship: '女儿', phone: '13900000002', is_primary: false },
    ]
  } catch (error) {
    console.error('Failed to load emergency contacts:', error)
  }
}

/**
 * 加载服务记录
 */
async function loadServiceRecords() {
  try {
    // 模拟数据
    serviceRecords.value = [
      { id: 1, service_date: '2026-02-01', service_type: 'meal', staff_name: '李师傅', status: 'completed', rating: 4.8, remark: '按时送达' },
      { id: 2, service_date: '2026-01-28', service_type: 'cleaning', staff_name: '王师傅', status: 'completed', rating: 5.0, remark: '打扫干净' },
      { id: 3, service_date: '2026-01-25', service_type: 'meal', staff_name: '李师傅', status: 'completed', rating: 4.9, remark: '' },
      { id: 4, service_date: '2026-01-20', service_type: 'medical', staff_name: '赵师傅', status: 'completed', rating: 4.7, remark: '测量血压正常' },
      { id: 5, service_date: '2026-01-15', service_type: 'accompany', staff_name: '刘师傅', status: 'completed', rating: 5.0, remark: '陪同就医' },
    ]
  } catch (error) {
    console.error('Failed to load service records:', error)
  }
}

/**
 * 返回列表
 */
function goBack() {
  router.push('/residents/elderly')
}

/**
 * 编辑档案
 */
function handleEdit() {
  // 跳转到编辑页面或打开编辑弹窗
  ElMessage.info('编辑功能开发中')
}

/**
 * 添加联系人
 */
function handleAddContact() {
  isEditContact.value = false
  contactDialogTitle.value = '添加联系人'
  resetContactForm()
  showContactDialog.value = true
}

/**
 * 编辑联系人
 */
function handleEditContact(contact: typeof emergencyContacts.value[0]) {
  isEditContact.value = true
  contactDialogTitle.value = '编辑联系人'
  Object.assign(contactForm, {
    id: contact.id,
    name: contact.name,
    relationship: contact.relationship,
    phone: contact.phone,
    is_primary: contact.is_primary,
  })
  showContactDialog.value = true
}

/**
 * 删除联系人
 */
async function handleDeleteContact(contact: typeof emergencyContacts.value[0]) {
  try {
    await ElMessageBox.confirm(`确定要删除联系人"${contact.name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    // 模拟删除
    emergencyContacts.value = emergencyContacts.value.filter(c => c.id !== contact.id)
    ElMessage.success('删除成功')
  } catch {
    // 取消删除
  }
}

/**
 * 保存联系人
 */
async function handleSaveContact() {
  if (!contactFormRef.value) return
  await contactFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      savingContact.value = true
      // 模拟保存
      await new Promise(resolve => setTimeout(resolve, 500))
      ElMessage.success(isEditContact.value ? '修改成功' : '添加成功')
      showContactDialog.value = false
      loadEmergencyContacts()
    } catch (error) {
      console.error('Failed to save contact:', error)
    } finally {
      savingContact.value = false
    }
  })
}

/**
 * 重置联系人表单
 */
function resetContactForm() {
  contactForm.id = null
  contactForm.name = ''
  contactForm.relationship = ''
  contactForm.phone = ''
  contactForm.is_primary = false
}

/**
 * 查看全部服务记录
 */
function handleViewAllRecords() {
  ElMessage.info('查看全部服务记录功能开发中')
}

/**
 * 计算年龄
 */
function calculateAge(birthDate: string): number {
  if (!birthDate) return 0
  return dayjs().diff(dayjs(birthDate), 'year')
}

/**
 * 身份证号脱敏
 */
function maskIdCard(idCard: string): string {
  if (!idCard || idCard.length < 15) return idCard
  return idCard.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2')
}

/**
 * 格式化日期
 */
function formatDate(date: string): string {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD')
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

/**
 * 获取服务类型标签颜色
 */
function getServiceTypeTag(type: string): string {
  const tagMap: Record<string, string> = {
    meal: 'success',
    medical: 'danger',
    cleaning: 'warning',
    shopping: 'primary',
    accompany: 'info',
  }
  return tagMap[type] || ''
}

/**
 * 获取服务类型文本
 */
function getServiceTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    meal: '助餐服务',
    medical: '医疗服务',
    cleaning: '清洁服务',
    shopping: '代购服务',
    accompany: '陪护服务',
  }
  return typeMap[type] || type
}

/**
 * 获取状态标签类型
 */
function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    processing: 'primary',
    completed: 'success',
    cancelled: 'info',
  }
  return typeMap[status] || ''
}

/**
 * 获取状态文本
 */
function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    pending: '待处理',
    processing: '处理中',
    completed: '已完成',
    cancelled: '已取消',
  }
  return textMap[status] || status
}

onMounted(() => {
  loadElderlyDetail()
  loadEmergencyContacts()
  loadServiceRecords()
})
</script>

<style scoped lang="scss">
.elderly-detail {
  padding: 20px;

  .page-header {
    display: flex;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      flex: 1;
      margin: 0 20px;
      font-size: 20px;
      font-weight: 500;
      color: #303133;
    }
  }

  .el-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .info-card {
    .basic-info {
      display: flex;
      align-items: flex-start;
      gap: 24px;

      .el-avatar {
        flex-shrink: 0;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      .info-descriptions {
        flex: 1;
      }
    }
  }

  .records-card {
    .rating {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      color: #e6a23c;
      font-weight: 500;

      .el-icon {
        font-size: 14px;
      }
    }
  }
}
</style>
