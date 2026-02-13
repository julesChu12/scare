<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>菜单管理</h3>
            <p class="subtitle">配置系统导航菜单、路由路径及权限标识</p>
          </div>
          <div class="header-actions">
            <el-button :icon="Refresh" @click="loadMenuTree">刷新</el-button>
            <el-button type="primary" :icon="Plus" @click="handleCreate">
              新建菜单
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索菜单名称/路径"
          clearable
          style="width: 250px"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" @click="handleSearch">查询</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="filteredMenuTree"
        row-key="id"
        border
        default-expand-all
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      >
        <el-table-column prop="name" label="菜单名称" min-width="180" />
        
        <el-table-column label="图标" width="70" align="center">
          <template #default="{ row }">
            <el-icon v-if="row.icon" class="menu-icon-display">
              <component :is="row.icon" />
            </el-icon>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="path" label="路由路径" min-width="150">
          <template #default="{ row }">
            <code class="path-code">{{ row.path }}</code>
          </template>
        </el-table-column>

        <el-table-column prop="component" label="组件路径" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <span class="component-text">{{ row.component || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="permission_code" label="权限标识" min-width="150">
          <template #default="{ row }">
            <el-tag v-if="row.permission_code" type="info" size="small" effect="plain">
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
      width="700px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="24">
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
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路由路径" prop="path">
              <el-input v-model="formData.path" placeholder="请输入路由路径" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件路径" prop="component">
              <el-input v-model="formData.component" placeholder="请输入组件路径" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
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
          </el-col>
          <el-col :span="12">
            <el-form-item label="权限标识" prop="permission_code">
              <el-input v-model="formData.permission_code" placeholder="如 menu:list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" :max="999" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="是否隐藏" prop="hidden">
              <el-switch v-model="formData.hidden" active-text="是" inactive-text="否" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio value="active">启用</el-radio>
                <el-radio value="inactive">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
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

// 搜索
const searchKeyword = ref('')

// 数据
const menuTree = ref<Menu[]>([])
const parentMenuOptions = ref<Menu[]>([])

// 树形过滤逻辑
const filteredMenuTree = computed(() => {
  if (!searchKeyword.value) return menuTree.value

  const keyword = searchKeyword.value.toLowerCase()
  
  function filterTree(tree: Menu[]): Menu[] {
    return tree
      .map(node => ({ ...node }))
      .filter(node => {
        const matches = 
          node.name.toLowerCase().includes(keyword) || 
          (node.path && node.path.toLowerCase().includes(keyword))
        
        if (node.children) {
          node.children = filterTree(node.children)
        }
        
        return matches || (node.children && node.children.length > 0)
      })
  }

  return filterTree(menuTree.value)
})

const handleSearch = () => {
  // Computed property handles search
}

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

<style scoped lang="scss">
.menu-management {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-left {
      h3 {
        margin: 0 0 4px 0;
        font-size: 18px;
        font-weight: 500;
      }
      .subtitle {
        margin: 0;
        font-size: 13px;
        color: #909399;
      }
    }
    
    .header-actions {
      display: flex;
      gap: 12px;
    }
  }

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid #f0f0f0;
  }

  .menu-icon-display {
    font-size: 18px;
    color: #409eff;
    vertical-align: middle;
  }

  .path-code {
    padding: 2px 6px;
    background-color: #f5f7fa;
    border-radius: 4px;
    color: #f56c6c;
    font-family: monospace;
    font-size: 12px;
  }

  .component-text {
    font-size: 12px;
    color: #606266;
  }

  .text-muted {
    color: #c0c4cc;
    font-size: 12px;
  }

  .icon-option {
    margin-right: 8px;
    vertical-align: middle;
  }
}
</style>
