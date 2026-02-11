# sCare 社区养老服务 - C端前端

基于 Vue 3 + TypeScript + PWA 的移动端应用，支持扫码即用的快速开通流程。

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全的JavaScript超集
- **Vite** - 快速的前端构建工具
- **Pinia** - Vue 3 状态管理
- **Vue Router** - 官方路由管理
- **Element Plus** - 基于 Vue 3 的组件库
- **Axios** - HTTP 客户端
- **PWA** - 渐进式Web应用支持

## 项目结构

```
frontend/
├── README.md                      # 前端总览文档
├── STRUCTURE_ANALYSIS.md          # 目录结构分析
│
├── c-end/                         # C端用户端（独立项目）
│   ├── src/
│   │   ├── api/                   # API 服务层
│   │   ├── views/                 # 页面组件（16个）
│   │   ├── store/                 # Pinia 状态管理
│   │   ├── router/                # 路由配置
│   │   ├── components/            # 组件
│   │   ├── composables/           # 组合式函数
│   │   └── ...
│   ├── public/                    # PWA 静态资源
│   ├── docs/                      # C端文档
│   ├── package.json               # 依赖配置
│   └── vite.config.ts             # Vite配置（含PWA）
│
└── management-portal/             # B端管理门户（独立项目）
    ├── src/
    │   ├── pages/                 # 页面组件（6个）
    │   ├── layouts/               # 布局组件
    │   ├── router/                # 路由配置
    │   ├── store/                 # Pinia 状态管理
    │   ├── config/                # Casbin 配置
    │   └── ...
    ├── docs/                      # B端文档
    ├── package.json               # 依赖配置
    └── vite.config.ts             # Vite配置
```

## 快速开始

### C端用户端

```bash
cd frontend/c-end
npm install
npm run dev      # http://localhost:5174
npm run build    # 生产构建
npm run preview  # 预览构建
```

### B端管理门户

```bash
cd frontend/management-portal
npm install
npm run dev      # http://localhost:3001
npm run build    # 生产构建
npm run preview  # 预览构建
```

## 核心功能

### 扫码即用流程

1. 用户扫描服务站二维码
2. 自动跳转到快速开通页面（带 station 参数）
3. 填写基本信息（手机号、验证码、姓名、身份证、地址、服务类型、服务描述）
4. 提交后自动完成：
   - 注册账号
   - 登录系统
   - 创建服务请求
5. 跳转到服务详情页

### 页面路由

- `/quick` - 快速开通（支持扫码携带 ?station=X 参数）
- `/login` - 密码登录（已有账号）
- `/home` - C端首页
- `/requests` - 我的服务列表
- `/requests/:id` - 服务详情
- `/profile` - 个人资料编辑

### API 集成

所有 API 调用都通过统一的 Axios 客户端，自动处理：

- Token 注入（从 localStorage）
- 401 自动跳转登录
- 统一错误提示（Element Plus Message）
- 请求/响应拦截

## 开发说明

### 状态管理

使用 Pinia 管理全局状态：

```typescript
// Token 管理
import { useTokenStore } from '@/stores'
const tokenStore = useTokenStore()
tokenStore.setToken(token)
tokenStore.clearToken()

// 用户信息管理
import { useUserStore } from '@/stores'
const userStore = useUserStore()
userStore.setUser(user)
userStore.setProfile(profile)
```

### API 调用示例

```typescript
import { authAPI, requestsAPI } from '@/api'

// 发送验证码
await authAPI.sendCode({ phone: '13800138000' })

// 快速开通
const result = await authAPI.quickStart({
  phone: '13800138000',
  code: '123456',
  station_id: 1,
  name: '张三',
  id_number: '110101199001011234',
  address: '北京市朝阳区XXX',
  service_type: 'meal'
})

// 获取我的服务列表
const { requests } = await requestsAPI.getMyRequests()
```

### 路由守卫

已配置全局路由守卫，自动检查登录状态：

- `meta.requiresAuth: true` 的路由需要登录
- 未登录时自动跳转到 `/login` 并保存原始目标路径
- 登录成功后自动跳转回原目标页面

## PWA 支持

本项目已配置 PWA，支持：

- 安装到主屏幕
- 离线访问（缓存静态资源）
- API 请求缓存（NetworkFirst 策略，1小时）

Manifest 配置在 `vite.config.ts` 中，Service Worker 自动生成。

## 环境变量

开发环境 API 代理已配置：

```typescript
// vite.config.ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true
    }
  }
}
```

生产环境需要配置实际的后端地址。

## 待完成功能

前端基础功能已完成，以下为待补齐或优化项：

1. **样式美化** - 根据设计稿优化移动端 UI 细节
2. **加载与空状态** - 统一 loading、骨架屏和空态
3. **错误边界** - 关键页面错误提示与恢复
4. **定位体验** - 位置授权/地图选点（当前仅地址解析）
5. **图片上传** - 服务请求支持上传图片
6. **消息推送** - 服务状态变更通知
7. **离线体验** - 完善 PWA 离线缓存与重试

## 注意事项

1. **后端接口依赖** - 需启动后端并保证 `/api/v1` 可用
2. **验证码服务** - 需后端提供短信/验证码接口
3. **地址解析** - 高德 API 需要申请 Key 并配置

## 下一步

1. 联调前后端接口
2. 优化 UI 样式和交互
3. 真机测试 PWA 功能
4. 性能优化和测试
