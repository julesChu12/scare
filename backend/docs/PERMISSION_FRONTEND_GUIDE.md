# 权限系统重构 - 前端接入指南

## 一、变更概述

本次重构将权限系统从 Casbin 迁移到自实现的数据库方案，统一了权限码格式，支持前端按钮级权限控制。

### 核心变更

| 变更项 | 旧方案 | 新方案 |
|--------|--------|--------|
| 权限存储 | Casbin policy.csv | MySQL permissions 表 |
| 权限码格式 | `perm_b_requests_list_get` | `service:request:list` |
| 菜单字段 | `permission` | `permission_code` |
| 中间件 | CasbinMiddleware | PermissionMiddleware |

### 权限码格式规范

```
module:resource:action

示例：
- dashboard              # 顶级菜单
- service:request        # 二级菜单
- service:request:list   # 按钮/API 权限
- service:task:claim     # 按钮/API 权限
```

---

## 二、API 变更

### 1. `/api/v1/b/auth/me` - 获取当前用户信息

**响应变更**：`permissions` 字段现在返回权限码列表（而非旧的 `perm_xxx` 格式）

```json
{
  "msg": "ok",
  "data": {
    "user": {
      "id": 4,
      "name": "王小红",
      "phone": "13800000004",
      "roles": ["staff"],
      "station_id": 1,
      "status": "active"
    },
    "permissions": [
      "dashboard",
      "service",
      "service:request",
      "service:request:list",
      "service:request:detail",
      "service:task",
      "service:task:pool",
      "service:task:my",
      "service:task:claim",
      "service:task:complete",
      "station",
      "station:list",
      "station:list:view",
      "station:zone",
      "station:zone:list"
    ]
  }
}
```

### 2. `/api/v1/b/menus/user` - 获取用户菜单

**响应变更**：菜单对象的 `permission` 字段改名为 `permission_code`

```json
{
  "msg": "success",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "name": "工作台",
      "path": "/dashboard",
      "component": "Dashboard",
      "icon": "Odometer",
      "permission_code": "dashboard",
      "sort": 1,
      "hidden": false,
      "status": "active",
      "children": []
    },
    {
      "id": 2,
      "parent_id": 0,
      "name": "服务管理",
      "path": "/services",
      "permission_code": "service",
      "children": [
        {
          "id": 6,
          "parent_id": 2,
          "name": "服务请求",
          "path": "/services/requests",
          "permission_code": "service:request:list"
        }
      ]
    }
  ]
}
```

### 3. `/api/v1/b/permissions/tree` - 获取权限树

**响应变更**：权限节点结构调整

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "tree": [
      {
        "id": "dashboard",
        "code": "dashboard",
        "label": "工作台",
        "type": "menu",
        "api_path": "",
        "method": "",
        "is_public": false,
        "disabled": false,
        "children": []
      },
      {
        "id": "service",
        "code": "service",
        "label": "服务管理",
        "type": "menu",
        "children": [
          {
            "id": "service:request",
            "code": "service:request",
            "label": "服务请求",
            "type": "menu",
            "children": [
              {
                "id": "service:request:list",
                "code": "service:request:list",
                "label": "请求列表",
                "type": "button",
                "api_path": "/api/v1/b/requests",
                "method": "GET"
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### 4. `/api/v1/b/roles/:role/permissions` - 获取/更新角色权限

**GET 响应**：返回权限码数组

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "role": "staff",
    "permissions": [
      "dashboard",
      "service:request:list",
      "service:task:pool"
    ]
  }
}
```

**PUT 请求**：提交权限码数组

```json
{
  "permissions": [
    "dashboard",
    "service:request:list",
    "service:request:detail",
    "service:task:pool",
    "service:task:my",
    "service:task:claim"
  ]
}
```

---

## 三、前端适配指南

### 1. 更新类型定义

```typescript
// 旧的菜单类型
interface OldMenu {
  permission: string;  // 删除
}

// 新的菜单类型
interface Menu {
  id: number;
  parent_id: number;
  name: string;
  path: string;
  component: string;
  icon: string;
  permission_code: string;  // 新字段名
  sort: number;
  hidden: boolean;
  status: string;
  children?: Menu[];
}

// 权限节点类型
interface PermissionNode {
  id: string;           // 权限码
  code: string;         // 权限码（同 id）
  label: string;        // 显示名称
  type: 'menu' | 'button' | 'resource';
  api_path?: string;
  method?: string;
  is_public?: boolean;
  disabled?: boolean;
  children?: PermissionNode[];
}

// 用户信息类型
interface UserInfo {
  user: {
    id: number;
    name: string;
    phone: string;
    roles: string[];
    station_id: number;
    status: string;
  };
  permissions: string[];  // 权限码数组
}
```

### 2. 权限判断工具函数

```typescript
// stores/permission.ts
import { defineStore } from 'pinia';

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    permissions: [] as string[],
  }),

  actions: {
    setPermissions(perms: string[]) {
      this.permissions = perms;
    },

    // 检查是否有某个权限
    hasPermission(code: string): boolean {
      // Admin 拥有所有权限
      const userStore = useUserStore();
      if (userStore.roles.includes('admin')) {
        return true;
      }
      return this.permissions.includes(code);
    },

    // 检查是否有任一权限
    hasAnyPermission(codes: string[]): boolean {
      return codes.some(code => this.hasPermission(code));
    },

    // 检查是否有所有权限
    hasAllPermissions(codes: string[]): boolean {
      return codes.every(code => this.hasPermission(code));
    },
  },
});
```

### 3. 按钮级权限指令

```typescript
// directives/permission.ts
import { usePermissionStore } from '@/stores/permission';
import type { Directive } from 'vue';

export const vPermission: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const permissionStore = usePermissionStore();
    const value = binding.value;

    const hasPermission = Array.isArray(value)
      ? permissionStore.hasAnyPermission(value)
      : permissionStore.hasPermission(value);

    if (!hasPermission) {
      el.parentNode?.removeChild(el);
    }
  },
};

// main.ts
app.directive('permission', vPermission);
```

### 4. 在组件中使用

```vue
<template>
  <div>
    <!-- 按钮级权限控制 -->
    <el-button
      v-permission="'service:request:list'"
      @click="handleView"
    >
      查看列表
    </el-button>

    <el-button
      v-permission="'station:list:create'"
      type="primary"
      @click="handleCreate"
    >
      新建站点
    </el-button>

    <!-- 多权限（任一满足） -->
    <el-button
      v-permission="['service:task:claim', 'service:task:complete']"
      @click="handleTask"
    >
      处理任务
    </el-button>

    <!-- 编程式判断 -->
    <el-button
      v-if="hasPermission('system:user:update')"
      @click="handleEdit"
    >
      编辑
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { usePermissionStore } from '@/stores/permission';

const permissionStore = usePermissionStore();
const hasPermission = (code: string) => permissionStore.hasPermission(code);
</script>
```

### 5. 登录后初始化权限

```typescript
// 登录成功后
async function handleLogin() {
  const { data } = await login(form);

  // 保存 token
  setToken(data.token);

  // 获取用户信息和权限
  const { data: userInfo } = await getUserInfo();

  // 设置权限
  const permissionStore = usePermissionStore();
  permissionStore.setPermissions(userInfo.permissions);

  // 设置用户信息
  const userStore = useUserStore();
  userStore.setUser(userInfo.user);
}
```

---

## 四、权限码对照表

### 服务管理模块

| 权限码 | 说明 | API |
|--------|------|-----|
| `service` | 服务管理菜单 | - |
| `service:request` | 服务请求菜单 | - |
| `service:request:list` | 查看请求列表 | GET /api/v1/b/requests |
| `service:request:detail` | 查看请求详情 | GET /api/v1/b/requests/:id |
| `service:task` | 任务管理菜单 | - |
| `service:task:pool` | 查看任务池 | GET /api/v1/b/tasks/pool |
| `service:task:my` | 查看我的任务 | GET /api/v1/b/tasks/my |
| `service:task:claim` | 认领任务 | POST /api/v1/b/tasks/:id/claim |
| `service:task:complete` | 完成任务 | POST /api/v1/b/tasks/:id/complete |

### 站点管理模块

| 权限码 | 说明 | API |
|--------|------|-----|
| `station` | 站点管理菜单 | - |
| `station:list` | 站点列表菜单 | - |
| `station:list:view` | 查看站点列表 | GET /api/v1/b/stations |
| `station:list:detail` | 查看站点详情 | GET /api/v1/b/stations/:id |
| `station:list:create` | 创建站点 | POST /api/v1/b/stations |
| `station:list:update` | 编辑站点 | PUT /api/v1/b/stations/:id |
| `station:list:delete` | 删除站点 | DELETE /api/v1/b/stations/:id |
| `station:zone` | 服务围栏菜单 | - |
| `station:zone:list` | 查看围栏列表 | GET /api/v1/b/zones |
| `station:zone:create` | 创建围栏 | POST /api/v1/b/zones |
| `station:zone:update` | 编辑围栏 | PUT /api/v1/b/zones/:id |
| `station:zone:delete` | 删除围栏 | DELETE /api/v1/b/zones/:id |

### 系统管理模块

| 权限码 | 说明 | API |
|--------|------|-----|
| `system` | 系统管理菜单 | - |
| `system:user` | 用户管理菜单 | - |
| `system:user:list` | 查看用户列表 | GET /api/v1/b/users |
| `system:user:detail` | 查看用户详情 | GET /api/v1/b/users/:id |
| `system:user:create` | 创建用户 | POST /api/v1/b/users |
| `system:user:update` | 编辑用户 | PUT /api/v1/b/users/:id |
| `system:user:roles` | 分配角色 | PUT /api/v1/b/users/:id/roles |
| `system:role` | 角色管理菜单 | - |
| `system:role:list` | 查看角色列表 | GET /api/v1/b/roles |
| `system:role:permissions` | 查看角色权限 | GET /api/v1/b/roles/:role/permissions |
| `system:role:update` | 更新角色权限 | PUT /api/v1/b/roles/:role/permissions |
| `system:permission:tree` | 查看权限树 | GET /api/v1/b/permissions/tree |
| `system:menu` | 菜单管理菜单 | - |
| `system:menu:list` | 查看菜单列表 | GET /api/v1/b/menus |
| `system:menu:create` | 创建菜单 | POST /api/v1/b/menus |
| `system:menu:update` | 编辑菜单 | PUT /api/v1/b/menus/:id |
| `system:menu:delete` | 删除菜单 | DELETE /api/v1/b/menus/:id |

### 内容管理模块

| 权限码 | 说明 | API |
|--------|------|-----|
| `content` | 内容管理菜单 | - |
| `content:banner` | 轮播图管理菜单 | - |
| `content:banner:list` | 查看轮播图列表 | GET /api/v1/b/banners |
| `content:banner:create` | 创建轮播图 | POST /api/v1/b/banners |
| `content:banner:update` | 编辑轮播图 | PUT /api/v1/b/banners/:id |
| `content:banner:delete` | 删除轮播图 | DELETE /api/v1/b/banners/:id |
| `content:news` | 新闻管理菜单 | - |
| `content:news:list` | 查看新闻列表 | GET /api/v1/b/news |
| `content:news:detail` | 查看新闻详情 | GET /api/v1/b/news/:id |
| `content:news:create` | 创建新闻 | POST /api/v1/b/news |
| `content:news:update` | 编辑新闻 | PUT /api/v1/b/news/:id |
| `content:news:delete` | 删除新闻 | DELETE /api/v1/b/news/:id |

---

## 五、角色默认权限

### Admin（系统管理员）
- 拥有所有权限（代码中特殊处理，跳过权限检查）

### Station Manager（站点管理员）
- dashboard
- service:* (服务管理全部)
- station:list:view, station:list:detail
- station:zone:* (围栏管理全部)
- system:user:* (用户管理全部)

### Staff（工作人员）
- dashboard
- service:request:list, service:request:detail
- service:task:pool, service:task:my, service:task:claim, service:task:complete
- station:list:view
- station:zone:list

---

## 六、迁移检查清单

- [ ] 更新菜单类型定义：`permission` → `permission_code`
- [ ] 更新权限存储：使用新的权限码格式
- [ ] 更新权限判断逻辑：适配新的权限码
- [ ] 更新权限树组件：适配新的节点结构
- [ ] 测试 Admin 登录：应获取全部菜单
- [ ] 测试 Staff 登录：应获取正确过滤的菜单
- [ ] 测试按钮级权限：v-permission 指令正常工作
- [ ] 测试无权限访问：返回 403 Forbidden
