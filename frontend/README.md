# sCare 前端执行与开发说明

## 1. 依赖说明

前端本地开发依赖：

| 组件 | 建议版本 | 用途 |
|------|----------|------|
| Node.js | 18+ | Vite、TypeScript、构建工具链 |
| npm / pnpm | 最新稳定版 | 安装依赖 |
| 后端 API | `http://localhost:8080` | 两个前端都依赖后端接口 |

默认访问地址：

| 应用 | 地址 | 说明 |
|------|------|------|
| 管理门户 | `http://localhost:3001` | B端管理后台 |
| C端 | `https://localhost:5174` | 本地开发启用 HTTPS |

启动顺序建议：

1. 先启动 `backend/docker-compose.yml` 中的 MySQL / Redis。
2. 启动后端 `go run . serve`。
3. 启动管理门户与 C 端前端。

## 2. 配置文件说明

### 2.1 管理门户：`frontend/management-portal/.env`

首次使用请复制：

```bash
cd frontend/management-portal
cp .env.development.example .env
```

关键配置项：

| 配置项 | 是否必须 | 默认值 | 说明 |
|--------|----------|--------|------|
| `VITE_API_BASE_URL` | 必须 | `http://localhost:8080/api/v1` | 管理门户请求的后端地址 |
| `VITE_APP_TITLE` | 否 | `sCare 管理后台` | 页面标题 |
| `VITE_AMAP_KEY` | 推荐 | 空 | 地图选点、围栏编辑需要 |
| `VITE_AMAP_SECURITY_JS_CODE` | 推荐 | 空 | 高德地图安全密钥 |
| `VITE_C_END_BASE_URL` | 推荐 | `https://localhost:5174` | 管理端生成二维码或跳转 C 端时使用 |

说明：

- 若只做登录、列表、基础 CRUD 联调，`VITE_API_BASE_URL` 足够
- 若涉及站点管理、地图选点、围栏编辑，必须补齐高德相关配置

### 2.2 C端

C端本地开发默认不依赖额外环境文件：

- 接口通过 Vite proxy 转发到 `http://localhost:8080`
- 默认 `baseURL` 为 `/api/v1`
- 可选环境变量 `VITE_MOCK_USER=true` 用于本地调试 Mock 用户逻辑

## 3. 启动方式

### 3.1 管理门户

```bash
cd frontend/management-portal
cp .env.development.example .env
npm install
npm run dev
```

启动后访问：`http://localhost:3001`

### 3.2 C端

```bash
cd frontend/c-end
npm install
npm run dev
```

启动后访问：`https://localhost:5174`

说明：

- C端使用 `@vitejs/plugin-basic-ssl`
- 浏览器首次访问可能提示本地证书风险，这是开发环境预期行为

## 4. 构建命令

### 管理门户

```bash
cd frontend/management-portal
npm run build
npm run preview
npm run lint
```

### C端

```bash
cd frontend/c-end
npm run build
npm run preview
```

## 5. 与后端联调

### 5.1 后端依赖

两个前端都要求后端已运行在 `http://localhost:8080`，至少保证以下接口可访问：

- `GET /api/v1/health`
- `POST /api/v1/b/auth/login`
- `POST /api/v1/c/auth/login`

### 5.2 管理门户代理与接口地址

管理门户默认通过 `.env` 中的 `VITE_API_BASE_URL` 直接请求后端。

推荐配置：

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_C_END_BASE_URL=https://localhost:5174
```

### 5.3 C端代理

C端 `vite.config.ts` 中默认代理：

- `/api` -> `http://localhost:8080`
- `/static` -> `http://localhost:8080`

因此本地开发通常无需单独配置 API 地址。

## 6. 常用调试场景

### 6.1 管理门户登录

使用种子管理员账号：

- 手机号：`13800000001`
- 密码：`Test@123`

### 6.2 C端登录

使用种子老年人账号：

- 手机号：`13800000008`
- 密码：`Test@123`

### 6.3 快速开通联调

当前快速开通接口需要的核心字段示例：

```ts
await authAPI.quickStart({
  phone: '13800138000',
  code: '123456',
  name: '张三',
  address: '北京市昌平区霍营街道示例地址',
  source_station_id: 1,
  service_type: 'meal',
  description: '需要送午餐',
  contact_name: '张三',
  contact_phone: '13800138000',
})
```

若需要更精确派单，可补充：

- `submit_lat` / `submit_lng`
- `service_lat` / `service_lng`
- `images`

## 7. 常见问题

### Q1: 管理门户可以打开，但接口全是 401 / 404

先检查：

1. 后端是否已启动在 `8080`
2. `frontend/management-portal/.env` 中 `VITE_API_BASE_URL` 是否正确
3. MySQL / Redis 是否正常启动

### Q2: C端页面打不开

通常是以下原因：

1. 浏览器未接受本地 HTTPS 证书
2. `5174` 端口被占用
3. 后端未启动，导致页面初始化请求失败

### Q3: 地图组件空白或报高德 Key 错误

请补齐：

- `VITE_AMAP_KEY`
- `VITE_AMAP_SECURITY_JS_CODE`

### Q4: 管理门户二维码跳转地址不对

请检查：

- `VITE_C_END_BASE_URL`

建议本地开发固定为：

```env
VITE_C_END_BASE_URL=https://localhost:5174
```

---

前端开发文档已按当前仓库实际启动方式、端口和配置项更新，可直接用于本地联调。
