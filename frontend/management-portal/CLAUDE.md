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
npx playwright test  # 运行 E2E 测试（需先启动后端和前端服务）
npx playwright test --ui  # 以 UI 模式运行测试
```

E2E 测试使用 Playwright，测试文件在 `tests/` 目录，无独立配置文件。测试前需确保后端（localhost:8080）和前端（localhost:3001）服务已启动。

## 架构概览

Vue 3 + TypeScript 管理后台（B端），采用 Composition API (`<script setup lang="ts">`) + Element Plus。

### 核心分层

- **API 层** (`src/api/index.ts`): 所有后端接口集中定义，按模块导出（authApi、taskApi、stationApi 等）。所有端点以 `/b/` 为前缀。
- **HTTP 客户端** (`src/utils/request.ts`): Axios 实例，自动附加 JWT token，处理 401 自动登出。
- **状态管理** (`src/store/modules/auth.ts`): 唯一的 Pinia store，管理认证状态、token 刷新、权限检查。localStorage key 前缀为 `b_`。
- **路由** (`src/router/index.ts`): 静态路由定义，每个路由通过 `meta.permission_code` 控制访问权限。
- **布局** (`src/layouts/BasicLayout.vue`): 侧边栏 + 顶栏 + 内容区，侧边栏菜单从后端 API (`/b/menus/user`) 动态获取。

### 权限体系（四层）

自定义权限管理系统，权限数据从后端获取，前端负责 UI 控制：

1. **路由守卫** (`src/router/guards/permission.guard.ts`): 检查登录状态和 `meta.permission_code`
2. **v-permission 指令** (`src/directives/permission.ts`): DOM 级别权限控制，无权限则移除元素
3. **usePermission composable** (`src/composables/usePermission.ts`): 提供 `hasPermission`、`hasAnyPermission` 等方法
4. **后端最终验证**: 前端只做 UI 控制，后端返回 403 作为最终权限验证

**权限特性**:
- 权限列表从 `/b/auth/current` API 获取，存储在 `user.permissions` 数组
- Admin 角色拥有所有权限，跳过权限检查
- 三级角色：staff → station_manager → admin

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

## 测试账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| Admin | 13800000001 | Test@123 |
| Station Manager | 13800000002 | Test@123 |
| Staff | 13800000004 | Test@123 |
