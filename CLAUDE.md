# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言偏好

- **回复语言**：中文
- **代码注释**：中文
- **Git 提交信息**：中文

## 项目概述

昌平区霍营街道社区养老信息分发平台（毕业设计项目）。基于地理围栏的社区养老服务信息分发系统，包含 B 端管理门户和 C 端用户端（PWA）。Go 模块名：`community-elderly-care-platform`。

## 常用命令

### 基础设施

```bash
docker compose up -d                                        # 启动 MySQL + Redis
docker exec -it scare_mysql mysql -u scare_user -pscare_pass scare_db  # 连接 MySQL
docker exec -it scare_redis redis-cli                       # 连接 Redis
./backend/scripts/rebuild_db.sh                             # 重建数据库（危险操作）
```

### 后端（工作目录：`backend/`）

```bash
air                          # 热重载启动（推荐，端口 8080）
go run . serve               # 标准启动（Cobra CLI）
go build ./...               # 编译检查
go test ./...                # 运行所有测试
go test ./pkg/geo/...        # 运行单个包测试
go test -run TestXxx ./internal/service/...  # 运行单个测试函数
./scripts/test_api.sh        # API 回归测试（提交前建议执行）

# GORM 模型重新生成（修改数据库表结构后）
go run cmd/tools/gen/gorm_gen.go

# Swagger 文档更新（接口变更后必须执行）
swag init -g main.go -o docs --parseDependency --parseInternal
```

### 管理门户（工作目录：`frontend/management-portal/`）

```bash
npm run dev          # 开发服务器（端口 3001，代理 /api → localhost:8080）
npm run build        # vue-tsc 类型检查 + Vite 构建
npm run lint         # ESLint 检查并自动修复
npm run format       # Prettier 格式化
```

首次运行需复制 `.env.development.example` → `.env.development`，填入 `VITE_AMAP_KEY`（高德地图 Key）。

### C 端用户端（工作目录：`frontend/c-end/`）

```bash
npm run dev          # 开发服务器（端口 5174）
npm run build        # vue-tsc 类型检查 + Vite 构建
```

## 架构概览

### 后端分层架构

```
Handler (参数校验/响应格式化)
  → Service (业务逻辑/事务处理)
    → Repository (GORM 数据库操作)
      → Model (GORM Gen 生成，.gen.go 文件禁止手动修改)
```

- **依赖注入**：`internal/router/deps.go` 中的 `Deps` 容器，`NewDeps()` 按 Repo → Service → Handler 顺序初始化
- **路由注册**：`internal/router/router.go` 统一注册，B 端路由在 `b_end.go`，C 端在 `c_end.go`
- **常量定义**：`internal/consts/` 下按职责分文件（`roles.go`、`status.go`、`service_types.go`、`menu.go`）
- **通知系统**：`internal/notify/` 提供邮件发送（SMTP），通过 `MailSender` 接口解耦
- **文件存储**：`internal/storage/` 支持本地存储和阿里云 OSS，通过工厂模式切换

### 路由与中间件

| 路径前缀 | 说明 | 中间件链 |
|----------|------|---------|
| `/api/v1/b/auth/login` | B 端公开接口 | 无 |
| `/api/v1/b/*` | B 端受保护接口 | AuthMiddleware → RequireEndType("b_end") → PermissionMiddleware |
| `/api/v1/c/auth/*` | C 端公开接口 | 无 |
| `/api/v1/c/*` 部分公开 | C 端部分无需登录（新闻、轮播、站点匹配、地理编码） | 无 |
| `/api/v1/c/*` | C 端受保护接口 | AuthMiddleware → RequireEndType("c_end") |
| `/swagger/*` | Swagger 文档 | 无 |

### 认证与权限

- **JWT 双端认证**：token 中 `type` 字段区分 `b_end` / `c_end`，Claims 包含 `user_id`、`identities`、`primary_role`、`station_id`
- **Token 黑名单**：Redis 实现（`TokenBlacklistService`），支持单 token 和用户级别撤销
- **自定义 RBAC 权限**：已从 Casbin 迁移为自定义 `PermissionService`（基于 `permissions`/`roles`/`role_permissions` 三表），Admin 角色跳过所有权限检查
- **角色体系**：B 端 `admin` > `station_manager` > `staff`；C 端 `elderly`、`family`
- **多身份支持**：`user_identities` 表支持一个用户绑定多个角色身份

### 地理围栏引擎（`pkg/geo/`）

服务请求创建时根据用户坐标自动匹配服务站点：
1. 内存加载所有围栏，按优先级降序排列
2. BoundingBox 快速排除不相关围栏
3. 射线法（Ray Casting）精确判断点在多边形内
4. 返回首个匹配的 StationID

### 前端架构

**技术栈**：Vue 3 + TypeScript + Vite + Element Plus + Pinia

**管理门户权限体系**（四层）：
1. **路由守卫**：`meta.permission_code` 控制页面访问
2. **v-permission 指令**：DOM 级别，无权限则移除元素
3. **usePermission composable**：`hasPermission()`、`hasAnyPermission()`
4. **侧边栏菜单**：从 `/b/menus/user` 动态获取

**API 层差异**：
- 管理门户：所有接口集中在 `src/api/index.ts` 单文件，按模块导出
- C 端：接口按模块拆分为独立文件（`src/api/auth.ts`、`requests.ts` 等），通过 `index.ts` 统一导出
- HTTP 客户端（`src/utils/request.ts`）自动附加 JWT、处理 401 自动登出

**C 端 PWA**：通过 `vite-plugin-pwa` 实现，支持离线使用和安装到桌面。

## 开发规范

### 后端关键约定

- API 统一响应：`{ "msg": "ok", "data": {...} }` / `{ "msg": "错误描述", "data": null }`
- 分页参数：`page`（从 1 开始）、`page_size`（默认 10）
- 数据库：逻辑外键，无数据库外键约束；软删除使用 `deleted_at`；字符集 `utf8mb4`
- GORM Gen 模型在 `internal/dao/model/*.gen.go`，**禁止手动修改**；扩展模型创建非 `.gen.go` 文件
- GORM Gen 查询在 `internal/dao/query/`，同样由工具生成

### 前端关键约定

- Composition API（`<script setup lang="ts">`）
- 管理门户页面组件在 `src/pages/`，C 端在 `src/views/`
- 状态管理：Pinia store，管理门户 localStorage key 前缀 `b_`
- 路径别名：`@` → `src/`
- 类型定义集中在 `src/types/api.ts`

### 数据库管理

- Schema 定义：`database/schema/schema.sql`（Docker 首次启动自动导入）
- 种子数据：`database/seeds/seed.sql`（Docker 首次启动自动导入）
- 迁移脚本：`database/migrations/`

## 测试账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| Admin | 13800000001 | Test@123 |
| Station Manager | 13800000002 | Test@123 |
| Staff | 13800000004 | Test@123 |
