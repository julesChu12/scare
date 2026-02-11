<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div>
            <span class="title">菜单管理</span>
            <span class="description">管理系统菜单配置，支持树形结构和排序</span>
          </div>
          <el-button type="primary" :icon="Plus" @click="handleCreate">
            新建菜单
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="menuTree"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="name" label="菜单名称" min-width="180">
          <template #default="{ row }">
            <el-icon v-if="row.icon" class="menu-icon">
              <component :is="row.icon" />
            </el-icon>
            <span>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" min-width="150" />
        <el-table-column prop="component" label="组件路径" min-width="200" />
        <el-table-column prop="permission_code" label="权限标识" min-width="150">
          <template #default="{ row }">
            <el-tag v-if="row.permission_code" type="info" size="small">
              {{ row.permission_code }}
            </el-tag>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="hidden" label="隐藏" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.hidden ? 'warning' : 'success'" size="small">
              {{ row.hidden ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAddChild(row)">
              添加子菜单
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-popconfirm
              title="确定要删除该菜单吗？"
              confirm-button-text="确定"
              cancel-button-text="取消"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑菜单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑菜单' : '新建菜单'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="上级菜单" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="parentMenuOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="选择上级菜单（不选则为顶级菜单）"
            clearable
            check-strictly
            :render-after-expand="false"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入菜单名称" />
        </el-form-item>
        <el-form-item label="路由路径" prop="path">
          <el-input v-model="formData.path" placeholder="请输入路由路径，如 /admin/menus" />
        </el-form-item>
        <el-form-item label="组件路径" prop="component">
          <el-input v-model="formData.component" placeholder="请输入组件路径，如 @/pages/MenuManagement/index.vue" />
        </el-form-item>
        <el-form-item label="图标" prop="icon">
          <el-select
            v-model="formData.icon"
            placeholder="选择图标"
            clearable
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="icon in iconOptions"
              :key="icon"
              :label="icon"
              :value="icon"
            >
              <el-icon class="icon-option">
                <component :is="icon" />
              </el-icon>
              <span>{{ icon }}</span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="权限标识" prop="permission_code">
          <el-input v-model="formData.permission_code" placeholder="请输入权限标识，如 menu:list" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="是否隐藏" prop="hidden">
          <el-switch v-model="formData.hidden" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">启用</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { menuApi } from '@/api'
import type { Menu, MenuRequest } from '@/types/api'

// 常用图标列表
const iconOptions = [
  'List',
  'User',
  'UserFilled',
  'Lock',
  'Setting',
  'Menu',
  'House',
  'Document',
  'Folder',
  'Files',
  'Edit',
  'Delete',
  'Plus',
  'Search',
  'View',
  'Hide',
  'Bell',
  'Message',
  'Calendar',
  'Clock',
  'Location',
  'Phone',
  'Star',
  'Flag',
  'Warning',
  'CircleCheck',
  'CircleClose',
  'Tools',
  'Operation',
  'DataAnalysis',
  'PieChart',
  'TrendCharts',
  'Grid',
  'Tickets',
  'OfficeBuilding',
  'School',
  'Shop',
]

// 状态
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editingId = ref<number | null>(null)

// 数据
const menuTree = ref<Menu[]>([])
const parentMenuOptions = ref<Menu[]>([])

// 表单
const formRef = ref<FormInstance>()
const formData = reactive<MenuRequest & { parent_id: number }>({
  parent_id: 0,
  name: '',
  path: '',
  component: '',
  icon: '',
  permission_code: '',
  sort: 0,
  hidden: false,
  status: 'active',
})

const formRules: FormRules = {
  name: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' },
    { min: 1, max: 50, message: '菜单名称长度为 1-50 个字符', trigger: 'blur' },
  ],
}

// 加载菜单树
const loadMenuTree = async () => {
  loading.value = true
  try {
    const res = await menuApi.getMenuTree()
    if (res.msg === 'success') {
      menuTree.value = res.data || []
      // 构建上级菜单选项（添加根节点选项）
      parentMenuOptions.value = [
        { id: 0, name: '顶级菜单', parent_id: 0, path: '', component: '', icon: '', permission_code: '', sort: 0, hidden: false, status: 'active' as const, children: menuTree.value },
      ]
    } else {
      ElMessage.error(res.msg || '获取菜单失败')
    }
  } catch (error) {
    ElMessage.error('获取菜单失败')
  } finally {
    loading.value = false
  }
}

// 重置表单
const resetForm = () => {
  formData.parent_id = 0
  formData.name = ''
  formData.path = ''
  formData.component = ''
  formData.icon = ''
  formData.permission_code = ''
  formData.sort = 0
  formData.hidden = false
  formData.status = 'active'
  editingId.value = null
}

// 新建菜单
const handleCreate = () => {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

// 添加子菜单
const handleAddChild = (parent: Menu) => {
  isEdit.value = false
  resetForm()
  formData.parent_id = parent.id
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = (menu: Menu) => {
  isEdit.value = true
  editingId.value = menu.id
  formData.parent_id = menu.parent_id
  formData.name = menu.name
  formData.path = menu.path
  formData.component = menu.component
  formData.icon = menu.icon
  formData.permission_code = menu.permission_code
  formData.sort = menu.sort
  formData.hidden = menu.hidden
  formData.status = menu.status
  dialogVisible.value = true
}

// 删除菜单
const handleDelete = async (menu: Menu) => {
  try {
    const res = await menuApi.deleteMenu(menu.id)
    if (res.msg === 'success') {
      ElMessage.success('删除成功')
      loadMenuTree()
    } else {
      ElMessage.error(res.msg || '删除失败')
    }
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const data: MenuRequest = {
      parent_id: formData.parent_id || undefined,
      name: formData.name,
      path: formData.path || undefined,
      component: formData.component || undefined,
      icon: formData.icon || undefined,
      permission_code: formData.permission_code || undefined,
      sort: formData.sort,
      hidden: formData.hidden,
      status: formData.status,
    }

    let res
    if (isEdit.value && editingId.value) {
      res = await menuApi.updateMenu(editingId.value, data)
    } else {
      res = await menuApi.createMenu(data)
    }

    if (res.msg === 'success' || res.msg === '创建成功' || res.msg === '更新成功') {
      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      dialogVisible.value = false
      loadMenuTree()
    } else {
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadMenuTree()
})
</script>

<style scoped>
.menu-management {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header .title {
  font-size: 18px;
  font-weight: 600;
}

.card-header .description {
  margin-left: 12px;
  font-size: 14px;
  color: #909399;
}

.menu-icon {
  margin-right: 8px;
  vertical-align: middle;
}

.text-muted {
  color: #909399;
}

.icon-option {
  margin-right: 8px;
  vertical-align: middle;
}

:deep(.el-table .cell) {
  display: flex;
  align-items: center;
}
</style>
