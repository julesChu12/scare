<template>
  <div class="user-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>用户管理</h3>
            <p>管理系统用户账号和角色分配</p>
          </div>
          <el-button type="primary" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新建用户
          </el-button>
        </div>
      </template>

      <!-- 搜索筛选 -->
      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名/手机号"
          clearable
          style="width: 200px"
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-select
          v-model="filterRole"
          placeholder="角色筛选"
          clearable
          style="width: 150px"
          @change="handleSearch"
        >
          <el-option label="系统管理员" value="admin" />
          <el-option label="站点管理员" value="station_manager" />
          <el-option label="工作人员" value="staff" />
        </el-select>
        <el-button @click="handleSearch">
          <el-icon><Search /></el-icon>
          搜索
        </el-button>
      </div>

      <!-- 用户列表 -->
      <el-table :data="users" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column label="B端角色" min-width="180">
          <template #default="{ row }">
            <template v-if="row.b_end_identities?.length">
              <el-tag
                v-for="identity in row.b_end_identities"
                :key="identity"
                style="margin-right: 5px"
              >
                {{ getIdentityName(identity) }}
              </el-tag>
            </template>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column label="C端角色" min-width="180">
          <template #default="{ row }">
            <template v-if="row.c_end_identities?.length">
              <el-tag
                v-for="identity in row.c_end_identities"
                :key="identity"
                type="warning"
                style="margin-right: 5px"
              >
                {{ getIdentityName(identity) }}
              </el-tag>
            </template>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="station_name" label="所属站点" min-width="150">
          <template #default="{ row }">
            {{ row.station_name || "-" }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : 'danger'"
              size="small"
            >
              {{ row.status === "active" ? "正常" : "禁用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="showEditDialog(row)"
            >
              编辑
            </el-button>
            <el-button
              type="primary"
              link
              size="small"
              @click="showRoleDialog(row)"
            >
              角色
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="loadUsers"
          @current-change="loadUsers"
        />
      </div>
    </el-card>

    <!-- 新建用户对话框 -->
    <el-dialog v-model="createDialogVisible" title="新建用户" width="500px">
      <el-form
        ref="createFormRef"
        :model="createForm"
        :rules="createRules"
        label-width="80px"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="createForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="createForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="所属站点" prop="station_id">
          <el-select
            v-model="createForm.station_id"
            placeholder="请选择站点"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="station in stations"
              :key="station.id"
              :label="station.name"
              :value="station.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select
            v-model="createForm.role"
            placeholder="请选择角色"
            style="width: 100%"
          >
            <el-option label="工作人员" value="staff" />
            <el-option label="站点管理员" value="station_manager" />
            <el-option label="系统管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="createForm.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitCreate"
          >确定</el-button
        >
      </template>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑用户" width="500px">
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="editRules"
        label-width="80px"
      >
        <el-form-item label="姓名" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="头像">
          <image-upload v-model="editForm.avatar" :limit="1" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select
            v-model="editForm.gender"
            placeholder="请选择性别"
            style="width: 100%"
          >
            <el-option label="男" value="male" />
            <el-option label="女" value="female" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="出生日期" prop="birth_date">
          <el-date-picker
            v-model="editForm.birth_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card">
          <el-input v-model="editForm.id_card" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="editForm.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitEdit"
          >确定</el-button
        >
      </template>
    </el-dialog>

    <!-- 身份分配对话框 -->
    <el-dialog v-model="roleDialogVisible" title="身份分配" width="500px">
      <div class="role-dialog-content">
        <p class="user-info" style="margin-bottom: 20px">
          用户：{{ currentUser?.name }} ({{ currentUser?.phone }})
        </p>

        <!-- B端身份 -->
        <div class="identity-section">
          <h4 style="margin-bottom: 10px; color: #606266">
            B端身份 (管理职能)
          </h4>
          <el-checkbox-group v-model="selectedBIdentities">
            <el-checkbox value="admin">系统管理员</el-checkbox>
            <el-checkbox value="station_manager">站点管理员</el-checkbox>
            <el-checkbox value="staff">工作人员</el-checkbox>
          </el-checkbox-group>
        </div>

        <el-divider />

        <!-- C端身份 -->
        <div class="identity-section">
          <h4 style="margin-bottom: 10px; color: #606266">
            C端身份 (服务对象)
          </h4>
          <el-checkbox-group v-model="selectedCIdentities">
            <el-checkbox value="elderly">长者</el-checkbox>
            <el-checkbox value="child">幼儿</el-checkbox>
            <el-checkbox value="pregnant">孕妇</el-checkbox>
            <el-checkbox value="disabled">残障人士</el-checkbox>
          </el-checkbox-group>
        </div>
      </div>
      <template #footer>
        <el-button @click="roleDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitRoles"
          >确定</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { Plus, Search } from "@element-plus/icons-vue";
import ImageUpload from "@/components/ImageUpload.vue";
import { userApi, stationApi } from "@/api";
import type { User, Station } from "@/types/api";
import type { FormInstance, FormRules } from "element-plus";

// 用户列表
const users = ref<User[]>([]);
const loading = ref(false);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

// 筛选条件
const searchKeyword = ref("");
const filterRole = ref("");

// 站点列表
const stations = ref<Station[]>([]);

// 新建用户表单
const createDialogVisible = ref(false);
const createFormRef = ref<FormInstance>();

interface CreateForm {
  name: string;
  phone: string;
  password: string;
  station_id: number | null;
  role: string;
  status: "active" | "inactive";
}

const createForm = reactive<CreateForm>({
  name: "",
  phone: "",
  password: "",
  station_id: null,
  role: "staff",
  status: "active",
});

const createRules: FormRules = {
  name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
  phone: [
    { required: true, message: "请输入手机号", trigger: "blur" },
    { pattern: /^1[3-9]\d{9}$/, message: "手机号格式不正确", trigger: "blur" },
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur", min: 6 },
  ],
};

// 编辑用户表单
const editDialogVisible = ref(false);
const editFormRef = ref<FormInstance>();
const editUserId = ref(0);

interface EditForm {
  name: string;
  avatar: string[];
  gender: string;
  birth_date: string;
  id_card: string;
  status: "active" | "inactive";
}

const editForm = reactive<EditForm>({
  name: "",
  avatar: [],
  gender: "",
  birth_date: "",
  id_card: "",
  status: "active",
});

const editRules: FormRules = {
  name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
  id_card: [
    {
      pattern: /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/,
      message: "身份证号格式不正确",
      trigger: "blur",
    },
  ],
};

// 角色分配
const roleDialogVisible = ref(false);
const currentUser = ref<User | null>(null);
const selectedBIdentities = ref<string[]>([]);
const selectedCIdentities = ref<string[]>([]);

const submitting = ref(false);

// 身份名称映射
const identityMap: Record<string, string> = {
  admin: "系统管理员",
  station_manager: "站点管理员",
  staff: "工作人员",
  user: "普通用户",
  child: "幼儿",
  elderly: "长者",
  disabled: "残障人士",
  pregnant: "孕妇",
};

function getIdentityName(identity: string) {
  return identityMap[identity] || identity;
}

// 加载用户列表
async function loadUsers() {
  loading.value = true;
  try {
    const res = await userApi.getUsers({
      page: pagination.page,
      page_size: pagination.pageSize,
    });
    if (res.msg === "ok") {
      users.value = res.data.items || [];
      pagination.total = res.data.total;
    }
  } catch (error) {
    ElMessage.error("加载用户列表失败");
  } finally {
    loading.value = false;
  }
}

// 加载站点列表
async function loadStations() {
  try {
    const res = await stationApi.getStations({ page: 1, page_size: 100 });
    if (res.msg === "ok") {
      stations.value = res.data.items || [];
    }
  } catch (error) {
    console.error("加载站点列表失败", error);
  }
}

// 搜索
function handleSearch() {
  pagination.page = 1;
  loadUsers();
}

// 显示新建对话框
function showCreateDialog() {
  createForm.name = "";
  createForm.phone = "";
  createForm.password = "";
  createForm.role = "staff";
  createForm.station_id = null;
  createForm.status = "active";
  createDialogVisible.value = true;
}

// 显示编辑对话框
function showEditDialog(user: User) {
  editUserId.value = user.id;
  editForm.name = user.name;
  editForm.avatar = user.avatar ? [user.avatar] : [];
  editForm.gender = user.gender || "";
  editForm.birth_date = user.birth_date || "";
  editForm.id_card = user.id_card || "";
  editForm.status = user.status;
  editDialogVisible.value = true;
}

// 提交新建用户
async function submitCreate() {
  if (!createFormRef.value) return;

  await createFormRef.value.validate(async (valid) => {
    if (!valid) return;

    submitting.value = true;
    try {
      await userApi.createUser({
        name: createForm.name,
        phone: createForm.phone,
        password: createForm.password,
        role: createForm.role,
        station_id: createForm.station_id || undefined,
        status: createForm.status,
      });
      ElMessage.success("创建成功");
      createDialogVisible.value = false;
      loadUsers();
    } catch (error) {
      ElMessage.error("创建失败");
    } finally {
      submitting.value = false;
    }
  });
}

// 提交编辑用户
async function submitEdit() {
  if (!editFormRef.value) return;

  await editFormRef.value.validate(async (valid) => {
    if (!valid) return;

    submitting.value = true;
    try {
      await userApi.updateUser(editUserId.value, {
        name: editForm.name,
        avatar: editForm.avatar[0] || "",
        gender: editForm.gender,
        birth_date: editForm.birth_date,
        id_card: editForm.id_card,
        status: editForm.status,
      });
      ElMessage.success("更新成功");
      editDialogVisible.value = false;
      loadUsers();
    } catch (error) {
      ElMessage.error("更新失败");
    } finally {
      submitting.value = false;
    }
  });
}

// 显示角色分配对话框
function showRoleDialog(user: User) {
  currentUser.value = user;
  selectedBIdentities.value = [...(user.b_end_identities || [])];
  selectedCIdentities.value = [...(user.c_end_identities || [])];
  roleDialogVisible.value = true;
}

// 提交角色分配
async function submitRoles() {
  if (!currentUser.value) return;

  submitting.value = true;
  try {
    const identities = [
      ...selectedBIdentities.value,
      ...selectedCIdentities.value,
    ];
    await userApi.updateUserIdentities(currentUser.value.id, {
      identities,
    });
    ElMessage.success("身份更新成功");
    roleDialogVisible.value = false;
    loadUsers();
  } catch (error) {
    ElMessage.error("角色更新失败");
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  loadUsers();
  loadStations();
});
</script>

<style scoped lang="scss">
.user-management {
  width: 100%;

  :deep(.el-card__body) {
    padding: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

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

  .filter-bar {
    display: flex;
    gap: 12px;
    margin-bottom: 0; /* Remove bottom margin if table has top margin, or keep it */
    padding-bottom: 16px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }

  .text-muted {
    color: #909399;
    font-size: 13px;
  }

  /* 强制表格铺满 */
  :deep(.el-table) {
    width: 100%;
  }
}

.role-dialog-content {
  .user-info {
    margin-bottom: 16px;
    color: #606266;
  }

  :deep(.el-checkbox-group) {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
}
</style>
