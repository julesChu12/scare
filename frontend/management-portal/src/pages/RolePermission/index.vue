<template>
  <div class="role-permission">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>角色权限管理</h3>
            <p>配置各角色的系统权限</p>
          </div>
        </div>
      </template>

      <div class="role-permission-content">
        <!-- 角色选择 -->
        <div class="role-selector">
          <h4>选择角色</h4>
          <el-radio-group v-model="selectedRole" @change="loadRolePermissions">
            <el-radio-button value="admin">系统管理员</el-radio-button>
            <el-radio-button value="station_manager">站点管理员</el-radio-button>
            <el-radio-button value="staff">工作人员</el-radio-button>
          </el-radio-group>
        </div>

        <!-- 权限树 -->
        <div class="permission-tree" v-loading="loading">
          <h4>权限配置</h4>
          <el-alert
            v-if="selectedRole === 'admin'"
            type="warning"
            :closable="false"
            show-icon
            style="margin-bottom: 16px"
          >
            系统管理员拥有所有权限，无需单独配置
          </el-alert>

          <el-tree
            v-else
            ref="treeRef"
            :data="permissionTree"
            :props="treeProps"
            show-checkbox
            node-key="id"
            :default-expand-all="true"
            :check-strictly="false"
          />

          <div class="action-buttons" v-if="selectedRole !== 'admin'">
            <el-button type="primary" :loading="submitting" @click="savePermissions">
              保存配置
            </el-button>
            <el-button @click="resetPermissions">重置</el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 权限说明 -->
    <el-card style="margin-top: 16px">
      <template #header>
        <h3 style="margin: 0; font-size: 16px">权限说明</h3>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="系统管理员 (admin)">
          拥有系统所有权限，包括用户管理、角色管理、站点管理等
        </el-descriptions-item>
        <el-descriptions-item label="站点管理员 (station_manager)">
          管理所属站点的任务分配、工作人员管理等
        </el-descriptions-item>
        <el-descriptions-item label="工作人员 (staff)">
          查看任务池、认领任务、完成任务等基本操作
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { permissionApi } from '@/api'
import type { PermissionNode } from '@/types/api'
import type { ElTree } from 'element-plus'

// 选中的角色
const selectedRole = ref('staff')

// 权限树数据
const permissionTree = ref<PermissionNode[]>([])
const loading = ref(false)
const submitting = ref(false)

// 当前角色的权限
const currentPermissions = ref<string[]>([])

// 树组件引用
const treeRef = ref<InstanceType<typeof ElTree>>()

// 树配置
const treeProps = {
  label: 'label',
  children: 'children',
}

// 加载权限树
async function loadPermissionTree() {
  loading.value = true
  try {
    const res = await permissionApi.getPermissionTree()
    if (res.msg === 'ok') {
      permissionTree.value = res.data?.tree || []
    }
  } catch (error) {
    ElMessage.error('加载权限树失败')
  } finally {
    loading.value = false
  }
}

// 加载角色权限
async function loadRolePermissions() {
  if (selectedRole.value === 'admin') return

  loading.value = true
  try {
    const res = await permissionApi.getRolePermissions(selectedRole.value)
    if (res.msg === 'ok') {
      currentPermissions.value = res.data.permissions || []
      // 设置选中状态
      await nextTick()
      if (treeRef.value) {
        treeRef.value.setCheckedKeys(currentPermissions.value)
      }
    }
  } catch (error) {
    ElMessage.error('加载角色权限失败')
  } finally {
    loading.value = false
  }
}

// 保存权限配置
async function savePermissions() {
  if (!treeRef.value) return

  const checkedKeys = treeRef.value.getCheckedKeys(false) as string[]
  const halfCheckedKeys = treeRef.value.getHalfCheckedKeys() as string[]
  const allKeys = [...checkedKeys, ...halfCheckedKeys]

  submitting.value = true
  try {
    await permissionApi.updateRolePermissions(selectedRole.value, {
      permissions: allKeys,
    })
    ElMessage.success('权限配置保存成功')
    currentPermissions.value = allKeys
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    submitting.value = false
  }
}

// 重置权限
function resetPermissions() {
  if (treeRef.value) {
    treeRef.value.setCheckedKeys(currentPermissions.value)
  }
}

onMounted(async () => {
  await loadPermissionTree()
  await loadRolePermissions()
})
</script>

<style scoped lang="scss">
.role-permission {
  .card-header {
    .header-left {
      h3 {
        margin: 0 0 4px 0;
        font-size: 16px;
      }
      p {
        margin: 0;
        font-size: 13px;
        color: #909399;
      }
    }
  }

  .role-permission-content {
    display: flex;
    gap: 24px;

    .role-selector {
      flex-shrink: 0;
      width: 200px;

      h4 {
        margin: 0 0 12px 0;
        font-size: 14px;
        color: #303133;
      }

      :deep(.el-radio-group) {
        display: flex;
        flex-direction: column;

        .el-radio-button {
          margin-left: 0;
          margin-bottom: 8px;

          .el-radio-button__inner {
            width: 100%;
            text-align: left;
            border-radius: 4px;
            border: 1px solid #dcdfe6;
          }
        }
      }
    }

    .permission-tree {
      flex: 1;
      border-left: 1px solid #ebeef5;
      padding-left: 24px;

      h4 {
        margin: 0 0 16px 0;
        font-size: 14px;
        color: #303133;
      }

      .action-buttons {
        margin-top: 24px;
        padding-top: 16px;
        border-top: 1px solid #ebeef5;
      }
    }
  }
}
</style>
