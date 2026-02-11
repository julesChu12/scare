# sCare B端管理后台前端技术规格说明书

## 📋 项目概述

### 项目信息
- **项目名称**: sCare 社区养老平台 - B端管理后台
- **技术栈**: Vue 3.4 + TypeScript + Vite + Element Plus + Pinia
- **后端**: Go + Gin + GORM + Casbin + Redis
- **认证方式**: JWT (localStorage) + Redis 黑名单
- **权限模型**: RBAC（基于角色的访问控制）
- **状态**: Locked
- **日期**: 2026-01-19

### 核心特性
- ✅ 多角色并集权限（用户可拥有多个角色，权限取并集）
- ✅ 动态权限树（后端接口获取，支持扩展）
- ✅ 公共权限管理（role:authenticated 特殊角色）
- ✅ 树形权限分配（父子联动，只保存叶子节点）
- ✅ 权限实时生效（Token 黑名单机制）

---

## 🏗️ 架构设计

### 1. 权限模型

#### 1.1 角色定义
| 角色 | 角色代码 | 权限范围 |
|------|---------|---------|
| 超级管理员 | admin | 全部 API 权限 |
| 站长 | station_manager | 任务、需求、站点、围栏、用户（查看）、通知 |
| 工作人员 | staff | 任务、需求（查看）、站点（查看）、围栏（查看）、通知 |
| **公共角色** | **authenticated** | 个人信息、登出、文件上传 |

#### 1.2 权限计算规则
```typescript
// 用户权限 = 所有角色权限的并集
user.roles = ['staff', 'station_manager']
user.permissions = Set(
  permissions(role:staff) ∪
  permissions(role:station_manager) ∪
  permissions(role:authenticated)
)
```

#### 1.3 权限树结构（3层）
```
Level 1: 业务模块 (module) - 分组节点，仅用于展示
  Level 2: 功能资源 (resource) - 可勾选节点
    Level 3: 操作动作 (action) - 叶子节点，实际保存到后端
```

**示例**：
```
需求管理 (module, id: module_requests)
  ├── 需求列表 (resource, id: res_requests_list)
  │   ├── 查看需求列表 (action, id: perm_requests_list_get)
  │   └── 导出需求列表 (action, id: perm_requests_list_export)
  ├── 需求详情 (resource, id: res_requests_detail)
  │   ├── 查看需求详情 (action, id: perm_requests_detail_get)
  │   └── 修改需求 (action, id: perm_requests_detail_put)
  └── 需求创建 (action, id: perm_requests_create)
```

#### 1.4 权限树节点定义
```typescript
interface PermissionNode {
  id: string;                    // 唯一标识（用于勾选状态）
  label: string;                 // 显示名称
  type: 'module' | 'resource' | 'action'; // 节点类型
  resource?: string;             // RESTful 资源路径（仅 action 类型有值）
  method?: string;               // HTTP 方法（仅 action 类型有值）
  children?: PermissionNode[];   // 子节点
  disabled?: boolean;            // 是否禁用（公共权限专用）
  isPublic?: boolean;            // 是否公共权限
}
```

---

## 🎨 前端页面设计

### 2.1 页面结构

```
/
├── /login                    # 登录页
├── /                         # 首页布局（需要登录）
│   ├── /dashboard            # 工作台（默认首页）
│   ├── /system               # 系统管理
│   │   ├── /users            # 用户管理
│   │   ├── /roles            # 角色管理
│   │   └── /permissions      # 权限管理
│   ├── /tasks                # 任务管理
│   ├── /requests             # 需求管理
│   ├── /stations             # 站点管理
│   └── /zones                # 围栏管理
└── /403                      # 无权限页面
```

---

### 2.2 核心页面详细设计

#### 2.2.1 登录页 (`/login`)

**功能**：
- 用户名/密码登录
- 记住我（localStorage 存储用户名）
- 登录后跳转到首页或原访问页面

**交互流程**：
```
1. 用户输入用户名/密码 → 点击登录
2. 调用 POST /api/v1/b/auth/login
3. 接收 { access_token, refresh_token, user: { id, username, roles: [] } }
4. 存储到 localStorage 和 Pinia Store
5. 调用 GET /api/v1/b/auth/me 获取完整用户信息和权限列表
6. 跳转到 /dashboard 或 redirect 参数指定的页面
```

**状态存储**：
```typescript
// localStorage
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "remember_username": "admin" // 如果勾选"记住我"
}
```

---

#### 2.2.2 首页布局 (`/`)

**组件结构**：
```
<Layout>
  ├── <Header>
  │   ├── Logo
  │   ├── 面包屑导航
  │   └── 用户信息下拉菜单（个人中心、修改密码、登出）
  ├── <Sidebar>
  │   └── 菜单（根据权限过滤）
  └── <Main>
      └── <router-view />
</Layout>
```

**菜单权限过滤**：
```typescript
// 菜单配置
const menuConfig = [
  {
    path: '/system',
    name: '系统管理',
    icon: 'Setting',
    permission: 'menu:system', // 需要的权限
    children: [
      { path: '/system/users', name: '用户管理', permission: 'perm_users_list_get' },
      { path: '/system/roles', name: '角色管理', permission: 'perm_roles_list_get' },
    ]
  },
  // ...
];

// 过滤逻辑
const filteredMenu = computed(() => {
  return filterMenuByPermissions(menuConfig, userStore.permissions);
});
```

---

#### 2.2.3 角色管理页 (`/system/roles`)

**功能**：
- 角色列表展示（表格）
- 新增角色（弹窗）
- 编辑角色（弹窗）
- 删除角色（二次确认）
- **权限分配**（树形勾选弹窗）

**权限分配交互**：
```
1. 点击"分配权限"按钮 → 打开弹窗
2. 调用 GET /api/v1/b/permissions/tree 获取完整权限树
3. 调用 GET /api/v1/b/roles/:role/permissions 获取当前角色已分配权限
4. 渲染树形勾选组件（Element Plus Tree）
   - 公共权限节点灰色显示，disabled=true
   - 根据已分配权限设置 checked 状态
5. 用户勾选/取消勾选节点
   - 勾选父节点 → 所有子节点自动勾选
   - 勾选子节点 → 父节点显示半选态（indeterminate）
   - 取消勾选父节点 → 所有子节点自动取消
6. 点击"保存"
   - 提取所有叶子节点（type='action'）的 id
   - 调用 PUT /api/v1/b/roles/:role/permissions
   - 后端返回影响的用户数和 Token 撤销信息
7. 刷新列表
```

**树形勾选组件关键配置**：
```vue
<el-tree
  ref="permissionTreeRef"
  :data="permissionTree"
  :props="{ children: 'children', label: 'label', disabled: 'disabled' }"
  node-key="id"
  show-checkbox
  :default-checked-keys="checkedKeys"
  :check-strictly="false"
  @check="handleCheckChange"
/>
```

```typescript
// 提取叶子节点逻辑
function getLeafPermissions(tree: PermissionNode[]): string[] {
  const leafIds: string[] = [];

  function traverse(nodes: PermissionNode[]) {
    nodes.forEach(node => {
      if (node.type === 'action') {
        leafIds.push(node.id);
      }
      if (node.children) {
        traverse(node.children);
      }
    });
  }

  traverse(tree);
  return leafIds;
}

// 保存时只传递叶子节点
function handleSave() {
  const checkedNodes = permissionTreeRef.value!.getCheckedNodes();
  const leafPermissions = checkedNodes
    .filter(node => node.type === 'action')
    .map(node => node.id);

  await roleApi.updatePermissions(currentRole.value, leafPermissions);
}
```

---

#### 2.2.4 权限管理页 (`/system/permissions`)

**功能**：
- 权限树完整展示（只读）
- 展示每个权限的关联角色
- 搜索权限（按名称/路径）

**布局**：
```
左侧：权限树（展开所有）
右侧：选中权限的详情
  - 权限名称
  - 资源路径 + HTTP 方法
  - 已分配的角色列表
  - 权限说明
```

**注意**：此页面不提供编辑功能，权限编辑统一在"角色管理"中进行。

---

#### 2.2.5 用户管理页 (`/system/users`)

**功能**：
- 用户列表（表格 + 分页）
- 新增用户（弹窗）
- 编辑用户（弹窗）
- **分配角色**（多选下拉框）
- 重置密码
- 启用/禁用用户

**分配角色交互**：
```
1. 点击"分配角色"按钮 → 打开弹窗
2. 展示多选下拉框（Element Plus Select multiple）
3. 选项：admin、station_manager、staff（不含 authenticated）
4. 用户可同时选择多个角色
5. 点击"保存"
   - 调用 PUT /api/v1/b/users/:id/roles
   - 后端自动撤销该用户的所有 Token（强制重新登录）
6. 刷新列表
```

**列表字段**：
| 字段 | 说明 |
|------|------|
| 用户名 | username |
| 角色 | roles（标签展示，如：admin、staff） |
| 所属站点 | station_name |
| 状态 | active / inactive |
| 创建时间 | created_at |
| 操作 | 编辑、分配角色、重置密码、启用/禁用 |

---

## 📦 状态管理（Pinia）

### 3.1 UserStore

```typescript
// stores/user.ts
import { defineStore } from 'pinia';

interface UserInfo {
  id: number;
  username: string;
  roles: string[];              // ['admin', 'staff']
  station_id: number;
  station_name: string;
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('access_token') || '',
    refreshToken: localStorage.getItem('refresh_token') || '',
    userInfo: null as UserInfo | null,
    permissions: new Set<string>(), // 权限并集
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    hasPermission: (state) => (permission: string) => {
      return state.permissions.has(permission);
    },
    hasAnyPermission: (state) => (permissions: string[]) => {
      return permissions.some(p => state.permissions.has(p));
    },
    hasAllPermissions: (state) => (permissions: string[]) => {
      return permissions.every(p => state.permissions.has(p));
    },
  },

  actions: {
    async login(username: string, password: string) {
      const { access_token, refresh_token, user } = await authApi.login({ username, password });

      this.token = access_token;
      this.refreshToken = refresh_token;
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('refresh_token', refresh_token);

      // 获取完整用户信息和权限
      await this.fetchUserInfo();
    },

    async fetchUserInfo() {
      const data = await authApi.me();
      this.userInfo = data.user;

      // 权限并集（后端已计算好）
      this.permissions = new Set(data.permissions || []);
    },

    async logout() {
      await authApi.logout();

      this.token = '';
      this.refreshToken = '';
      this.userInfo = null;
      this.permissions.clear();

      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    },
  },
});
```

---

## 🔐 权限控制实现

### 4.1 路由守卫

```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/layouts/BasicLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: 'dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'system/roles',
          component: () => import('@/views/system/Roles.vue'),
          meta: { permission: 'perm_roles_list_get' },
        },
        // ...
      ],
    },
  ],
});

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore();

  // 1. 检查登录态
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    return next({ path: '/login', query: { redirect: to.fullPath } });
  }

  // 2. 检查权限
  if (to.meta.permission && !userStore.hasPermission(to.meta.permission)) {
    return next('/403');
  }

  next();
});

export default router;
```

---

### 4.2 按钮级权限指令

```typescript
// directives/permission.ts
import type { Directive } from 'vue';
import { useUserStore } from '@/stores/user';

export const permission: Directive = {
  mounted(el, binding) {
    const userStore = useUserStore();
    const { value } = binding;

    if (!value) return;

    // 支持字符串或数组
    const permissions = Array.isArray(value) ? value : [value];
    const hasPermission = permissions.some(p => userStore.hasPermission(p));

    if (!hasPermission) {
      el.parentNode?.removeChild(el);
    }
  },
};

// main.ts
import { permission } from './directives/permission';
app.directive('permission', permission);
```

**使用示例**：
```vue
<template>
  <!-- 单个权限 -->
  <el-button v-permission="'perm_requests_create'" type="primary">
    创建需求
  </el-button>

  <!-- 多个权限（任一满足） -->
  <el-button v-permission="['perm_requests_update', 'perm_requests_delete']">
    编辑
  </el-button>
</template>
```

---

## 🔌 后端接口契约

### 5.1 认证接口

#### 5.1.1 登录
```
POST /api/v1/b/auth/login

Request:
{
  "username": "admin",
  "password": "Test@123"
}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 86400,
    "user": {
      "id": 1,
      "username": "admin",
      "roles": ["admin"]
    }
  }
}
```

#### 5.1.2 获取用户信息（含权限）
```
GET /api/v1/b/auth/me
Headers: Authorization: Bearer {token}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "roles": ["admin", "staff"],
      "station_id": 1,
      "station_name": "中心站点"
    },
    "permissions": [
      "perm_requests_list_get",
      "perm_requests_create",
      "perm_tasks_claim",
      // ... 所有角色权限的并集
    ]
  }
}
```

---

### 5.2 权限管理接口

#### 5.2.1 获取权限树
```
GET /api/v1/b/permissions/tree

Response: {
  "code": 0,
  "msg": "ok",
  "data": {
    "tree": [
      {
        "id": "public",
        "label": "公共权限",
        "type": "module",
        "isPublic": true,
        "children": [...]
      },
      {
        "id": "module_requests",
        "label": "需求管理",
        "type": "module",
        "children": [...]
      },
      // ... 其他模块
    ]
  }
}
```

#### 5.2.2 获取角色权限
```
GET /api/v1/b/roles/:role/permissions

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "role": "station_manager",
    "permissions": ["perm_requests_list_get", ...]
  }
}
```

#### 5.2.3 更新角色权限
```
PUT /api/v1/b/roles/:role/permissions

Request:
{
  "permissions": ["perm_requests_list_get", "perm_requests_create", ...]
}

Response:
{
  "code": 0,
  "msg": "权限更新成功",
  "data": {
    "affected_users": 5,
    "tokens_revoked": true,
    "revoked_user_ids": [2, 5, 8, 12, 15]
  }
}
```

---

### 5.3 用户管理接口

#### 5.3.1 更新用户角色
```
PUT /api/v1/b/users/:id/roles

Request:
{
  "roles": ["admin", "staff"]
}

Response:
{
  "code": 0,
  "msg": "角色更新成功，用户需要重新登录",
  "data": {
    "user_id": 5,
    "roles": ["admin", "staff"],
    "tokens_revoked": true
  }
}
```

---

## 🛠️ 后端实现调整

### 6.1 JWT Payload 调整

**修改前**：
```go
type Claims struct {
    UserID    int64  `json:"uid"`
    Role      string `json:"role"`      // 单一角色
    Type      string `json:"type"`
    StationID int64  `json:"station_id"`
    jwtlib.RegisteredClaims
}
```

**修改后**：
```go
type Claims struct {
    UserID    int64    `json:"uid"`
    Roles     []string `json:"roles"`    // 多个角色
    Type      string   `json:"type"`
    StationID int64    `json:"station_id"`
    jwtlib.RegisteredClaims
}
```

---

### 6.2 中间件逻辑调整

**Casbin 中间件**（权限检查）：
```go
// backend/internal/middleware/casbin.go
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles, exists := c.Get("user_roles")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "missing roles"})
            return
        }

        userRoles := roles.([]string)
        obj := c.Request.URL.Path
        act := c.Request.Method

        // 自动添加 authenticated 角色（公共权限）
        allRoles := append(userRoles, "authenticated")

        // 遍历所有角色，任一匹配即通过
        for _, role := range allRoles {
            allowed, err := enforcer.Enforce("role:"+role, obj, act)
            if err == nil && allowed {
                c.Next()
                return
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "msg": "permission denied",
        })
    }
}
```

---

### 6.3 authenticated 角色的策略初始化

**在 seed.sql 中添加**：
```sql
-- =====================================================
-- 公共权限（authenticated 角色 - 所有登录用户）
-- =====================================================
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'role:authenticated', '/api/v1/b/auth/me', 'GET'),
('p', 'role:authenticated', '/api/v1/b/auth/logout', 'POST'),
('p', 'role:authenticated', '/api/v1/c/auth/me', 'GET'),
('p', 'role:authenticated', '/api/v1/c/auth/logout', 'POST'),
('p', 'role:authenticated', '/api/v1/b/upload', 'POST'),
('p', 'role:authenticated', '/api/v1/c/upload', 'POST'),
('p', 'role:authenticated', '/api/v1/b/notifications', 'GET'),
('p', 'role:authenticated', '/api/v1/b/notifications/*/read', 'POST'),
('p', 'role:authenticated', '/api/v1/c/notifications', 'GET'),
('p', 'role:authenticated', '/api/v1/c/notifications/*/read', 'POST');
```

---

## ⚠️ 关键注意事项

### 7.1 安全性

1. **Token 刷新机制**：
   - Access Token 过期时间：24小时
   - Refresh Token 过期时间：7天
   - 使用 Refresh Token 自动刷新，减少重新登录

2. **权限变更即时生效**：
   - 更新角色权限后，立即撤销相关用户的所有 Token
   - 用户下次请求时 Token 被拒绝，前端自动跳转到登录页

3. **公共权限保护**：
   - 前端：禁用公共权限节点的勾选/删除操作
   - 后端：更新权限时拒绝删除 `role:authenticated` 的策略

---

### 7.2 性能优化

1. **权限缓存**：
   ```typescript
   // 前端：用户权限缓存到 Pinia Store，避免重复请求
   if (!userStore.permissions.size) {
       await userStore.fetchUserInfo();
   }
   ```

2. **权限树缓存**：
   ```typescript
   // 前端：权限树首次获取后缓存到 sessionStorage
   const cachedTree = sessionStorage.getItem('permission_tree');
   if (cachedTree) {
       return JSON.parse(cachedTree);
   }
   const tree = await permissionApi.getTree();
   sessionStorage.setItem('permission_tree', JSON.stringify(tree));
   ```

---

## 📂 项目目录结构

```
frontend/management-portal/
├── public/
├── src/
│   ├── api/                    # API 接口封装
│   │   ├── auth.ts             # 认证接口
│   │   ├── user.ts             # 用户管理接口
│   │   ├── role.ts             # 角色管理接口
│   │   ├── permission.ts       # 权限管理接口
│   │   └── request.ts          # Axios 实例配置
│   ├── assets/                 # 静态资源
│   ├── components/             # 公共组件
│   │   ├── PermissionTree.vue  # 权限树勾选组件
│   │   └── ...
│   ├── directives/             # 自定义指令
│   │   └── permission.ts       # 权限指令
│   ├── layouts/                # 布局组件
│   │   └── BasicLayout.vue     # 基础布局
│   ├── router/                 # 路由配置
│   │   └── index.ts
│   ├── stores/                 # Pinia 状态管理
│   │   └── user.ts             # 用户状态
│   ├── utils/                  # 工具函数
│   │   └── request.ts          # Axios 拦截器
│   ├── views/                  # 页面组件
│   │   ├── Login.vue
│   │   ├── Dashboard.vue
│   │   └── system/
│   │       ├── Users.vue
│   │       ├── Roles.vue
│   │       └── Permissions.vue
│   ├── App.vue
│   └── main.ts
├── .env.development
├── package.json
└── SPEC.md
```

---

## ✅ 实现 Checklist

### 后端调整（优先）
- [ ] JWT Payload 改为 `roles: []`
- [ ] 修改 `GenerateToken()` 方法接收多角色
- [ ] 修改登录逻辑（查询用户的所有角色）
- [ ] 调整认证中间件（设置 `user_roles`）
- [ ] 调整 Casbin 中间件（自动注入 `authenticated` 角色）
- [ ] 新增接口：`GET /api/v1/b/permissions/tree`
- [ ] 新增接口：`GET /api/v1/b/roles/:role/permissions`
- [ ] 新增接口：`PUT /api/v1/b/roles/:role/permissions`
- [ ] 新增接口：`PUT /api/v1/b/users/:id/roles`
- [ ] 修改 `/auth/me` 接口返回权限并集
- [ ] 在 seed.sql 添加 `role:authenticated` 策略

### 前端实现
- [ ] 实现 UserStore（Pinia）
- [ ] 实现权限指令 `v-permission`
- [ ] 实现路由守卫
- [ ] 实现登录页
- [ ] 实现首页布局
- [ ] 实现角色管理页
- [ ] 实现权限树勾选组件
- [ ] 实现用户管理页

---

**Spec 已锁定，可以进入执行模式。**
