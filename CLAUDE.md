# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with this repository.

## 语言偏好

- **回复语言**：中文
- **代码注释**：中文
- **Git 提交信息**：中文

## 项目概述

昌平区霍营街道社区养老信息分发平台（毕业设计项目）。基于地理围栏的社区养老服务信息分发系统，包含 B 端管理门户和 C 端用户端（PWA）。Go 模块名：`community-elderly-care-platform`。

## 项目结构

### 整体目录结构

```
sCare/
├── backend/                    # 后端服务（Go 1.25）
│   ├── cmd/                    # Cobra CLI 程序入口
│   │   ├── server/             # API 服务启动入口
│   │   └── tools/              # 工具集（GORM Gen 等）
│   ├── internal/               # 私有应用代码（分层架构）
│   │   ├── config/             # 配置管理
│   │   ├── consts/             # 常量定义
│   │   ├── dao/                # GORM Gen 模型与查询
│   │   │   ├── model/          # 数据模型（.gen.go 文件禁止手动修改）
│   │   │   └── query/          # 查询构建器
│   │   ├── dto/                # 数据传输对象
│   │   ├── handler/            # HTTP 处理器
│   │   ├── service/            # 业务逻辑层
│   │   ├── repository/         # 数据访问层
│   │   ├── middleware/         # 中间件（JWT、权限、端隔离）
│   │   ├── router/             # 路由注册与依赖注入
│   │   │   ├── deps.go         # 依赖注入容器
│   │   │   ├── router.go       # 路由统一注册
│   │   │   ├── b_end.go        # B 端路由
│   │   │   └── c_end.go        # C 端路由
│   │   ├── notify/             # 通知服务（邮件）
│   │   └── storage/            # 文件存储（本地/OSS）
│   ├── pkg/                    # 公共库
│   │   ├── geo/                # 地理围栏引擎（射线法）
│   │   ├── jwt/                # JWT 工具
│   │   ├── crypto/             # 加密工具
│   │   ├── logger/             # Zap 日志工具
│   │   ├── redis/              # Redis 工具
│   │   └── email/              # 邮件工具
│   ├── database/               # 数据库
│   │   ├── schema/             # 表结构定义
│   │   ├── seeds/              # 种子数据
│   │   ├── migrations/         # 迁移脚本
│   │   └── docs/               # 数据库文档
│   ├── docs/                   # 后端文档
│   │   ├── testing/            # 测试文档
│   │   ├── 01-开发指南.md
│   │   ├── 02-系统架构设计.md
│   │   ├── 03-数据库设计.md
│   │   ├── 04-API接口设计.md
│   │   └── 05-配置说明.md
│   ├── scripts/                # 工具脚本（API 测试等）
│   └── storage/                # 文件存储目录
│
├── frontend/                   # 前端项目
│   ├── c-end/                  # C端用户端（PWA）
│   │   ├── src/
│   │   │   ├── api/            # API 层（模块化拆分）
│   │   │   ├── views/          # 页面组件
│   │   │   ├── components/     # 组件
│   │   │   ├── composables/    # 组合式函数
│   │   │   ├── router/         # 路由配置
│   │   │   ├── store/          # Pinia 状态管理
│   │   │   ├── utils/          # 工具函数
│   │   │   ├── styles/         # 样式文件（含 CSS 变量系统）
│   │   │   └── config/         # 配置文件
│   │   ├── docs/               # C端文档
│   │   └── public/             # PWA 静态资源
│   └── management-portal/      # B端管理门户
│       ├── src/
│       │   ├── pages/          # 页面组件
│       │   ├── layouts/        # 布局组件
│       │   ├── components/     # 组件
│       │   ├── composables/    # 组合式函数
│       │   ├── router/         # 路由配置
│       │   │   └── guards/     # 路由守卫（权限控制）
│       │   ├── store/
│       │   │   └── modules/    # Pinia 模块
│       │   ├── api/            # API 层（集中管理）
│       │   ├── directives/     # 自定义指令（权限）
│       │   ├── types/          # TypeScript 类型定义
│       │   ├── utils/          # 工具函数
│       │   └── config/         # 配置文件
│       ├── docs/               # B端文档
│       └── tests/              # Playwright E2E 测试
│
├── docs/                       # 项目级文档
│   ├── README.md              # 文档导航索引
│   ├── SPEC.md                # 项目规格文档（毕业设计）
│   ├── PRD.md                 # 产品需求文档
│   ├── PROJECT_STATUS.md      # 项目开发状态
│   └── BUSINESS_RULES.md      # 业务规则文档
│
├── deployment/                 # 部署配置
│   ├── docker/                 # Docker 配置
│   ├── nginx/                  # Nginx 配置
│   ├── kubernetes/             # K8s 配置
│   ├── scripts/                # 部署脚本
│   └── docs/                   # 部署文档
│
├── paper/                      # 毕业论文相关
│   ├── thesis/                 # 论文草稿
│   ├── thesis_new/             # 论文新版本
│   └── 附件/                   # 论文附件
│
├── CLAUDE.md                   # Claude Code AI 开发指南（本文件）
├── README.md                   # 项目说明文档
└── docker-compose.yml          # Docker Compose 配置
```

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

## 架构设计

### 后端分层架构

```
Handler (参数校验/响应格式化)
  → Service (业务逻辑/事务处理)
    → Repository (GORM 数据库操作)
      → Model (GORM Gen 生成，.gen.go 文件禁止手动修改)
```

**核心设计原则**：
- **依赖注入**：`internal/router/deps.go` 中的 `Deps` 容器，`NewDeps()` 按 Repo → Service → Handler 顺序初始化
- **路由注册**：`internal/router/router.go` 统一注册，B 端路由在 `b_end.go`，C 端在 `c_end.go`
- **常量定义**：`internal/consts/` 下按职责分文件（`roles.go`、`status.go`、`service_types.go`、`menu.go`）
- **通知系统**：`internal/notify/` 提供邮件发送（SMTP），通过 `MailSender` 接口解耦
- **文件存储**：`internal/storage/` 支持本地存储和阿里云 OSS，通过工厂模式切换

### 后端技术栈

**核心依赖**（基于 `go.mod`）：
- **Web 框架**：Gin 1.11
- **ORM**：GORM 1.31 + GORM Gen 0.3.27（代码生成）
- **数据库驱动**：MySQL 1.6、SQLite 1.6（开发）
- **缓存**：Redis 9.5（go-redis/v9）+ miniredis（测试）
- **认证**：JWT 5.2（golang-jwt/jwt/v5）
- **日志**：Zap 1.27 + lumberjack（日志轮转）
- **CLI**：Cobra 1.10 + Viper 1.18
- **API 文档**：Swagger 1.16（swaggo/swag）
- **测试**：testify 1.11
- **其他**：UUID 生成、Excel 处理、加密库

### 前端技术栈

**B 端管理门户**（`frontend/management-portal/package.json`）：
- Vue 3.4 + TypeScript 5.3
- Vite 6.0（构建工具）
- Element Plus 2.5（UI 组件库）
- Pinia 2.1（状态管理）
- Vue Router 4.2（路由）
- ECharts 5.4 + vue-echarts 6.6（图表）
- Axios 1.6（HTTP 客户端）
- 高德地图（地理围栏编辑）
- 富文本编辑器（@wangeditor/editor）
- Playwright 1.58（E2E 测试）

**C 端用户端**（`frontend/c-end/package.json`）：
- Vue 3.4 + TypeScript 5.4
- Vite 5.0（构建工具）
- Element Plus 2.5（UI 组件库）
- Pinia 2.1（状态管理）
- Vue Router 4.2（路由）
- Axios 1.6（HTTP 客户端）
- @vueuse/core 10.7（组合式工具库）
- vite-plugin-pwa 0.17（PWA 支持）

### 路由与中间件设计

**中间件链**：

| 路径前缀 | 说明 | 中间件链 |
|----------|------|---------|
| `/api/v1/b/auth/login` | B 端公开接口 | 无 |
| `/api/v1/b/*` | B 端受保护接口 | AuthMiddleware → RequireEndType("b_end") → PermissionMiddleware |
| `/api/v1/c/auth/*` | C 端公开接口 | 无 |
| `/api/v1/c/*` 部分公开 | C 端部分无需登录（新闻、轮播、站点匹配、地理编码） | 无 |
| `/api/v1/c/*` | C 端受保护接口 | AuthMiddleware → RequireEndType("c_end") |
| `/swagger/*` | Swagger 文档 | 无 |

**中间件职责**：
- `AuthMiddleware`：JWT token 验证，解析用户信息存入 context
- `RequireEndType`：端隔离，防止 B/C 端 token 跨端使用
- `PermissionMiddleware`：自定义 RBAC 权限检查，Admin 跳过

### 认证与权限系统

**JWT 双端认证**：
- Token 中 `type` 字段区分 `b_end` / `c_end`
- Claims 包含：`user_id`、`identities`、`primary_role`、`station_id`
- 支持访问令牌 + 刷新令牌机制

**Token 黑名单**：
- Redis 实现（`TokenBlacklistService`）
- 支持单 token 和用户级别撤销
- 登出时将 token 加入黑名单

**自定义 RBAC 权限**：
- 已从 Casbin 迁移为自定义 `PermissionService`
- 基于三表模型：`permissions`（权限）、`roles`（角色）、`role_permissions`（关联）
- Admin 角色跳过所有权限检查

**角色体系**：
- B 端：`admin` > `station_manager` > `staff`
- C 端：`elderly`（老年人）、`family`（家属）

**多身份支持**：
- `user_identities` 表支持一个用户绑定多个角色身份
- 用户可切换不同身份（不同站点角色）

### 地理围栏引擎（`pkg/geo/`）

**核心算法**（点在多边形 - 射线法）：

服务请求创建时根据用户坐标自动匹配服务站点：
1. 内存加载所有围栏，按优先级降序排列
2. BoundingBox 快速排除不相关围栏
3. 射线法（Ray Casting）精确判断点在多边形内
4. 返回首个匹配的 StationID

**性能优化**：
- 启动时将围栏数据加载到内存
- 使用 BoundingBox 预过滤
- 优先级排序减少不必要的计算

**兜底规则**：未命中任何围栏时，分配到最近的站点。

### 前端架构设计

**技术栈统一**：Vue 3 + TypeScript + Vite + Element Plus + Pinia

**管理门户权限体系**（四层防护）：
1. **路由守卫**：`meta.permission_code` 控制页面访问
2. **v-permission 指令**：DOM 级别，无权限则移除元素
3. **usePermission composable**：`hasPermission()`、`hasAnyPermission()`
4. **侧边栏菜单**：从 `/b/menus/user` 动态获取

**API 层组织差异**：
- **管理门户**：所有接口集中在 `src/api/index.ts` 单文件，按模块导出
- **C 端**：接口按模块拆分为独立文件（`src/api/auth.ts`、`requests.ts` 等），通过 `index.ts` 统一导出

**HTTP 客户端**：
- 自动附加 JWT token
- 处理 401 自动登出并跳转登录页
- 统一错误处理和提示

**C 端 PWA 特性**：
- `vite-plugin-pwa` 实现
- 支持离线使用和安装到桌面
- Workbox 对 API 请求使用 NetworkFirst 缓存策略（1 小时过期）

**C 端适老化特性**：
- CSS 变量系统支持三档字体大小调节
- `useFontSize` composable 实现字体大小切换
- `[data-font-size]` 属性选择器驱动样式变化

## 开发规范

### 后端关键约定

- **API 统一响应**：`{ "msg": "ok", "data": {...} }` / `{ "msg": "错误描述", "data": null }`
- **分页参数**：`page`（从 1 开始）、`page_size`（默认 10）
- **数据库约定**：
  - 逻辑外键，无数据库外键约束
  - 软删除使用 `deleted_at`
  - 字符集 `utf8mb4`
- **GORM Gen 约定**：
  - 模型在 `internal/dao/model/*.gen.go`，**禁止手动修改**
  - 扩展模型创建非 `.gen.go` 文件
  - 查询在 `internal/dao/query/`，同样由工具生成

### 前端关键约定

- **组件风格**：Composition API（`<script setup lang="ts">`）
- **目录组织**：
  - 管理门户页面组件在 `src/pages/`
  - C 端页面组件在 `src/views/`
- **状态管理**：Pinia store
  - 管理门户 localStorage key 前缀 `b_`
  - C 端 localStorage key 前缀 `c_`
- **路径别名**：`@` → `src/`
- **类型定义**：集中在 `src/types/api.ts`

### 数据库管理

- **Schema 定义**：`database/schema/schema.sql`（Docker 首次启动自动导入）
- **种子数据**：`database/seeds/seed.sql`（Docker 首次启动自动导入）
- **迁移脚本**：`database/migrations/`

### 新增功能开发流程

**后端**：
1. Model：如有新表，修改 `database/schema/` 后运行 `go run cmd/tools/gen/gorm_gen.go` 生成模型
2. Repository：在 `internal/repository/` 创建 `xxx_repo.go`
3. Service：在 `internal/service/` 创建 `xxx_service.go`
4. Handler：在 `internal/handler/` 创建 `xxx_handler.go`
5. 依赖注入：在 `internal/router/deps.go` 注册
6. 路由：在 `internal/router/b_end.go` 或 `c_end.go` 中注册
7. Swagger：运行 `swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal`

**前端**：
1. 定义 API 接口（`src/api/`）
2. 创建页面组件（`src/pages/` 或 `src/views/`）
3. 添加路由配置（如需权限控制，添加 `meta.permission_code`）
4. 更新状态管理（如需要）

## 测试账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| Admin | 13800000001 | Test@123 |
| Station Manager | 13800000002 | Test@123 |
| Staff | 13800000004 | Test@123 |

## 文档索引

### 项目级文档
- `docs/README.md` - 文档导航索引
- `docs/SPEC.md` - 项目规格文档（毕业设计）
- `docs/PRD.md` - 产品需求文档
- `docs/PROJECT_STATUS.md` - 项目开发状态
- `docs/BUSINESS_RULES.md` - 业务规则文档

### 后端文档
- `backend/docs/01-开发指南.md` - 开发环境配置
- `backend/docs/02-系统架构设计.md` - 系统架构设计
- `backend/docs/03-数据库设计.md` - 数据库设计
- `backend/docs/04-API接口设计.md` - API 接口设计
- `backend/docs/05-配置说明.md` - 配置文件说明
- `backend/docs/testing/TEST_REPORT.md` - 测试报告

### 前端文档
- `frontend/management-portal/docs/README.md` - B 端文档索引
- `frontend/management-portal/SPEC.md` - B 端规格文档
- `frontend/c-end/docs/README.md` - C 端文档索引
- `frontend/c-end/SPEC.md` - C 端规格文档

### 数据库文档
- `backend/database/docs/README.md` - 数据库文档索引
- `backend/database/ER_DIAGRAM.md` - ER 图文档
