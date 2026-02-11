# 多角色支持改造计划

## 改造范围

### 已完成 ✅
- [x] 数据库表结构改造
- [x] 现有数据迁移
- [x] 测试数据准备

### 待完成

## 一、后端改造（Backend）

### 1.1 数据模型层
**文件**：`backend/internal/model/user.go`

```go
// 修改 User 结构体
type User struct {
    ID          uint      `gorm:"primaryKey"`
    Phone       string    `gorm:"uniqueIndex"`
    Name        string
    Role        string    `gorm:"column:role;comment:已废弃,仅保留兼容性"` // 标记为废弃
    PrimaryRole string    `gorm:"column:primary_role"` // 新增
    Roles       []UserRole `gorm:"foreignKey:UserID"` // 关联
    // ... 其他字段
}

// 新增 UserRole 模型
type UserRole struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    Role      string    `gorm:"index"`
    IsPrimary bool      `gorm:"default:false"`
    Status    string    `gorm:"default:active"`
    GrantedAt time.Time
    // ...
}

// 辅助方法
func (u *User) GetActiveRoles() []string {
    var roles []string
    for _, ur := range u.Roles {
        if ur.Status == "active" {
            roles = append(roles, ur.Role)
        }
    }
    return roles
}

func (u *User) HasRole(role string) bool {
    for _, ur := range u.Roles {
        if ur.Role == role && ur.Status == "active" {
            return true
        }
    }
    return false
}

func (u *User) HasAnyRole(roles []string) bool {
    for _, r := range roles {
        if u.HasRole(r) {
            return true
        }
    }
    return false
}
```

### 1.2 认证层
**文件**：`backend/internal/service/auth_service.go`

```go
// 登录响应改造
type LoginResponse struct {
    Token        string   `json:"token"`
    RefreshToken string   `json:"refresh_token"`
    UserID       uint     `json:"user_id"`
    Name         string   `json:"name"`
    Phone        string   `json:"phone"`
    PrimaryRole  string   `json:"primary_role"`  // 主角色
    Roles        []string `json:"roles"`         // 所有角色（新增）
    StationID    *uint    `json:"station_id"`
}

// 登录时预加载角色
func (s *AuthService) Login(phone, password string) (*LoginResponse, error) {
    var user model.User
    err := s.db.Preload("Roles", "status = ?", "active").
        Where("phone = ?", phone).First(&user).Error
    // ...
    return &LoginResponse{
        // ...
        PrimaryRole: user.PrimaryRole,
        Roles:       user.GetActiveRoles(), // 返回所有角色
    }, nil
}
```

**新增**：`backend/internal/api/v1/user.go`
```go
// 角色切换接口
type SwitchRoleRequest struct {
    Role string `json:"role" binding:"required"`
}

func (h *UserHandler) SwitchRole(c *gin.Context) {
    var req SwitchRoleRequest
    // 验证用户是否拥有该角色
    userID := c.GetUint("user_id")
    var user model.User
    h.db.Preload("Roles").First(&user, userID)

    if !user.HasRole(req.Role) {
        c.JSON(403, gin.H{"msg": "您没有该角色权限"})
        return
    }

    // 更新 JWT Token 中的 current_role
    newToken := generateTokenWithRole(userID, req.Role)
    c.JSON(200, gin.H{
        "msg": "ok",
        "data": gin.H{
            "token": newToken,
            "current_role": req.Role,
        },
    })
}
```

### 1.3 权限中间件
**文件**：`backend/internal/middleware/casbin.go`

```go
// 修改权限检查逻辑
func (m *CasbinMiddleware) CheckPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")

        // 获取用户所有激活的角色
        var user model.User
        m.db.Preload("Roles", "status = ?", "active").First(&user, userID)

        // 检查是否有任一角色拥有权限
        hasPermission := false
        for _, ur := range user.Roles {
            ok, _ := m.enforcer.Enforce(
                fmt.Sprintf("role:%s", ur.Role),
                c.Request.URL.Path,
                c.Request.Method,
            )
            if ok {
                hasPermission = true
                // 设置当前使用的角色到 context
                c.Set("current_role", ur.Role)
                break
            }
        }

        if !hasPermission {
            c.JSON(403, gin.H{"msg": "forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 1.4 JWT Token 改造
**文件**：`backend/pkg/jwt/jwt.go`

```go
// JWT Claims 添加 current_role 字段
type Claims struct {
    UserID      uint     `json:"uid"`
    PrimaryRole string   `json:"primary_role"` // 主角色
    CurrentRole string   `json:"current_role"` // 当前使用的角色（新增）
    Roles       []string `json:"roles"`        // 所有角色（新增）
    StationID   uint     `json:"station_id"`
    jwt.RegisteredClaims
}

// 生成 Token 时包含所有角色
func GenerateToken(user *model.User) (string, error) {
    claims := Claims{
        UserID:      user.ID,
        PrimaryRole: user.PrimaryRole,
        CurrentRole: user.PrimaryRole, // 默认使用主角色
        Roles:       user.GetActiveRoles(),
        StationID:   user.StationID,
        // ...
    }
    // ...
}
```

---

## 二、前端改造（Frontend）

### 2.1 类型定义
**文件**：`frontend/management-portal/src/types/api.ts`

```typescript
export interface User {
  user_id: number
  primary_role: 'elderly' | 'family' | 'staff' | 'station_manager' | 'admin'  // 主角色
  roles: ('elderly' | 'family' | 'staff' | 'station_manager' | 'admin')[]     // 所有角色（新增）
  current_role?: string  // 当前使用的角色（新增）
  name: string
  phone: string
  station_id: number | null
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user_id: number
  primary_role: string  // 主角色
  roles: string[]       // 所有角色（新增）
  name: string
  phone: string
  station_id: number | null
}

// 角色切换请求
export interface SwitchRoleRequest {
  role: string
}

// 角色切换响应
export interface SwitchRoleResponse {
  token: string
  current_role: string
}
```

### 2.2 API 层
**文件**：`frontend/management-portal/src/api/index.ts`

```typescript
export const authApi = {
  login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return request.post('/auth/login', data)
  },

  // 新增：角色切换
  switchRole(role: string): Promise<ApiResponse<SwitchRoleResponse>> {
    return request.post('/auth/switch-role', { role })
  },
}
```

### 2.3 状态管理
**文件**：`frontend/management-portal/src/store/modules/auth.ts`

```typescript
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>('')
  const user = ref<User | null>(null)
  const currentRole = ref<string>('') // 当前角色（新增）

  const isLoggedIn = computed(() => !!token.value)

  // 当前角色的计算属性
  const activeRole = computed(() => currentRole.value || user.value?.primary_role || '')

  // 检查是否拥有某个角色
  function hasRole(role: string): boolean {
    return user.value?.roles?.includes(role) || false
  }

  // 检查是否拥有任一角色
  function hasAnyRole(roles: string[]): boolean {
    return roles.some(role => hasRole(role))
  }

  // 登录
  async function login(credentials: LoginRequest) {
    const response = await authApi.login(credentials)
    const { token: accessToken, refresh_token, roles, primary_role, ...userData } = response.data

    setToken(accessToken)
    setRefreshToken(refresh_token)

    const userInfo: User = {
      ...userData,
      primary_role: primary_role as User['primary_role'],
      roles: roles as User['roles'],
      current_role: primary_role, // 默认使用主角色
    }
    setUser(userInfo)
    currentRole.value = primary_role // 设置当前角色

    return userInfo
  }

  // 切换角色
  async function switchRole(role: string) {
    if (!hasRole(role)) {
      ElMessage.error('您没有该角色权限')
      return false
    }

    try {
      const response = await authApi.switchRole(role)
      const { token: newToken, current_role } = response.data

      // 更新 Token 和当前角色
      setToken(newToken)
      currentRole.value = current_role
      if (user.value) {
        user.value.current_role = current_role
        setUser(user.value)
      }

      ElMessage.success(`已切换到${getRoleText(current_role)}角色`)
      return true
    } catch (error) {
      ElMessage.error('角色切换失败')
      return false
    }
  }

  function getRoleText(role: string): string {
    const roleMap: Record<string, string> = {
      elderly: '老年人',
      family: '家属',
      staff: '工作人员',
      station_manager: '站长',
      admin: '管理员',
    }
    return roleMap[role] || role
  }

  return {
    token,
    user,
    currentRole,
    activeRole,
    isLoggedIn,
    hasRole,
    hasAnyRole,
    login,
    switchRole,
    getRoleText,
    // ... 其他方法
  }
})
```

### 2.4 路由守卫
**文件**：`frontend/management-portal/src/router/guards/permission.guard.ts`

```typescript
export function setupPermissionGuard(router: Router) {
  router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()

    // ... 登录检查

    // 角色权限检查（支持多角色）
    const requiredRoles = to.meta.roles as string[] | undefined
    if (requiredRoles && requiredRoles.length > 0) {
      // 检查用户是否拥有任一所需角色
      if (!authStore.hasAnyRole(requiredRoles)) {
        ElMessage.error('权限不足，无法访问此页面')
        return next(from.path || '/')
      }
    }

    next()
  })
}
```

### 2.5 角色切换组件
**新建文件**：`frontend/management-portal/src/components/RoleSwitcher.vue`

```vue
<template>
  <el-dropdown v-if="user && user.roles.length > 1" @command="handleSwitchRole">
    <span class="role-switcher">
      <el-tag :type="getRoleTagType(activeRole)">
        {{ getRoleText(activeRole) }}
      </el-tag>
      <el-icon><ArrowDown /></el-icon>
    </span>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="role in user.roles"
          :key="role"
          :command="role"
          :disabled="role === activeRole"
        >
          <el-icon v-if="role === activeRole"><Check /></el-icon>
          {{ getRoleText(role) }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowDown, Check } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()
const router = useRouter()

const user = computed(() => authStore.user)
const activeRole = computed(() => authStore.activeRole)

async function handleSwitchRole(role: string) {
  const success = await authStore.switchRole(role)
  if (success) {
    // 切换角色后重定向到对应首页
    const homeRoutes: Record<string, string> = {
      staff: '/tasks/pool',
      station_manager: '/dashboard',
      elderly: '/my-requests',
      family: '/family/requests',
      admin: '/admin/users',
    }
    router.push(homeRoutes[role] || '/')
  }
}

function getRoleText(role: string): string {
  return authStore.getRoleText(role)
}

function getRoleTagType(role: string): string {
  const typeMap: Record<string, string> = {
    staff: 'primary',
    station_manager: 'warning',
    elderly: 'success',
    family: 'info',
    admin: 'danger',
  }
  return typeMap[role] || ''
}
</script>

<style scoped lang="scss">
.role-switcher {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 12px;
  border-radius: 4px;

  &:hover {
    background-color: #f5f7fa;
  }
}
</style>
```

### 2.6 布局组件更新
**文件**：`frontend/management-portal/src/layouts/BasicLayout.vue`

```vue
<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '200px'">
      <!-- 侧边栏 -->
    </el-aside>

    <el-container>
      <el-header>
        <div class="header-left">
          <el-icon @click="isCollapse = !isCollapse">
            <Expand v-if="isCollapse" />
            <Fold v-else />
          </el-icon>
        </div>

        <div class="header-right">
          <!-- 角色切换组件（新增） -->
          <RoleSwitcher />

          <!-- 用户下拉菜单 -->
          <el-dropdown @command="handleCommand">
            <el-button text>
              <el-icon><User /></el-icon>
              <span>{{ authStore.user?.name }}</span>
              <el-icon><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Expand, Fold, User, ArrowDown } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'
import RoleSwitcher from '@/components/RoleSwitcher.vue'

const authStore = useAuthStore()
const router = useRouter()
const isCollapse = ref(false)

async function handleCommand(command: string) {
  if (command === 'logout') {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示')
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }
}
</script>
```

---

## 三、用户管理改造

### 3.1 用户列表增强
**文件**：`frontend/management-portal/src/pages/Admin/UserManagement.vue`

```vue
<el-table-column label="角色" width="300">
  <template #default="{ row }">
    <el-tag
      v-for="role in row.roles"
      :key="role"
      :type="getRoleTagType(role)"
      size="small"
      style="margin-right: 8px"
    >
      {{ getRoleText(role) }}
      <el-icon v-if="role === row.primary_role"><Star /></el-icon>
    </el-tag>
  </template>
</el-table-column>
```

### 3.2 用户编辑对话框
```vue
<el-form-item label="用户角色" required>
  <el-checkbox-group v-model="form.roles">
    <el-checkbox label="elderly">老年人</el-checkbox>
    <el-checkbox label="family">家属</el-checkbox>
    <el-checkbox label="staff">工作人员</el-checkbox>
    <el-checkbox label="station_manager">站长</el-checkbox>
    <el-checkbox label="admin">管理员</el-checkbox>
  </el-checkbox-group>
  <div class="form-tip">至少选择一个角色</div>
</el-form-item>

<el-form-item label="主角色" required>
  <el-radio-group v-model="form.primary_role">
    <el-radio
      v-for="role in form.roles"
      :key="role"
      :label="role"
    >
      {{ getRoleText(role) }}
    </el-radio>
  </el-radio-group>
  <div class="form-tip">用户登录后默认使用的角色</div>
</el-form-item>
```

---

## 四、测试计划

### 4.1 单元测试
- [ ] 用户多角色数据模型测试
- [ ] 角色检查方法测试
- [ ] 角色切换逻辑测试

### 4.2 集成测试
- [ ] 登录返回多角色
- [ ] 角色切换API测试
- [ ] 权限检查（多角色场景）

### 4.3 端到端测试
- [ ] 多角色用户登录
- [ ] 角色切换流程
- [ ] 不同角色访问不同页面
- [ ] 角色权限边界测试

---

## 五、实施顺序

### Phase 1: 后端核心（优先级最高）
1. ✅ 数据库迁移（已完成）
2. 数据模型改造
3. 登录API改造
4. 角色切换API
5. 权限中间件改造

### Phase 2: 前端核心
1. 类型定义更新
2. Store 改造
3. 路由守卫更新
4. 角色切换组件

### Phase 3: 用户体验
1. 布局组件更新
2. 用户管理界面
3. 首页路由逻辑

### Phase 4: 测试与优化
1. 功能测试
2. 权限测试
3. 用户体验优化

---

## 六、预计工作量

- 后端改造：4-6小时
- 前端改造：4-6小时
- 用户管理：2-3小时
- 测试调试：2-3小时

**总计**：12-18小时（约1.5-2天）

---

## 七、回滚方案

如果改造出现问题，可以：
1. 数据库回滚：删除 user_roles 表，users.primary_role 字段
2. 代码回滚：Git revert
3. 数据不会丢失：原 users.role 字段仍保留

---

**准备好开始了吗？我建议按以下顺序进行：**
1. 先完成后端改造（数据层 → API层 → 权限层）
2. 再完成前端改造（类型 → Store → 组件）
3. 最后集成测试

**是否立即开始？或者有其他考虑？**
