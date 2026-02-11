# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 语言偏好

回复、代码注释、Git 提交信息均使用中文。

## 常用命令

```bash
npm run dev          # 启动开发服务器 (port 3001, 代理 /api 和 /static 到 localhost:8080)
npm run build        # 类型检查 (vue-tsc) + Vite 构建
npm run lint         # ESLint 检查并自动修复
npm run format       # Prettier 格式化 src/
```

E2E 测试使用 Playwright，测试文件在 `tests/` 目录，无独立配置文件。

## 架构概览

Vue 3 + TypeScript 管理后台（B端），采用 Composition API (`<script setup lang="ts">`) + Element Plus。

### 核心分层

- **API 层** (`src/api/index.ts`): 所有后端接口集中定义，按模块导出（authApi、taskApi、stationApi 等）。所有端点以 `/b/` 为前缀。
- **HTTP 客户端** (`src/utils/request.ts`): Axios 实例，自动附加 JWT token，处理 401 自动登出。
- **状态管理** (`src/store/modules/auth.ts`): 唯一的 Pinia store，管理认证状态、token 刷新、权限检查。localStorage key 前缀为 `b_`。
- **路由** (`src/router/index.ts`): 静态路由定义，每个路由通过 `meta.permission_code` 控制访问权限。
- **布局** (`src/layouts/BasicLayout.vue`): 侧边栏 + 顶栏 + 内容区，侧边栏菜单从后端 API (`/b/menus/user`) 动态获取。

### 权限体系（多层）

1. **路由守卫** (`src/router/guards/permission.guard.ts`): 检查登录状态和 `meta.permission_code`
2. **v-permission 指令** (`src/directives/permission.ts`): DOM 级别权限控制，无权限则移除元素
3. **usePermission composable** (`src/composables/usePermission.ts`): 提供 `hasPermission`、`hasAnyPermission` 等方法
4. **Casbin 配置** (`src/config/casbin/`): 前端 RBAC 模型，三级角色：staff → station_manager → admin

### 页面组织

页面在 `src/pages/` 下按功能目录组织，每个目录包含 `index.vue`。共享组件在 `src/components/`（地图相关组件、图片上传）。

### 外部服务依赖

- 高德地图 (AMap): 地理围栏编辑和展示，key 通过 `VITE_AMAP_KEY` 配置
- ECharts: 统计图表（通过 vue-echarts）

## 环境配置

复制 `.env.development.example` 为 `.env.development`，需配置：
- `VITE_API_BASE_URL`: 后端 API 地址（默认 `http://localhost:8080/api/v1`）
- `VITE_AMAP_KEY`: 高德地图 API Key

## 类型定义

所有 API 相关的 TypeScript 接口集中在 `src/types/api.ts`。

## Vite 配置要点

- `@` 别名指向 `src/`
- 代码分割：element-plus、vue-vendor（vue + router + pinia）、echarts 各自独立 chunk
