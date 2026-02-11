# C 端前端接入指南

## 一、概述

C 端面向老年人（elderly）和家属（family）用户，提供服务请求、个人资料管理等功能。与 B 端不同，C 端不需要复杂的权限系统，仅需基本的身份认证。

### 用户类型

| 类型 | 说明 |
|------|------|
| `elderly` | 老年人，服务的直接接受者 |
| `family` | 家属，可代老年人发起服务请求 |

### 认证方式

- **验证码登录**（推荐）：手机号 + 短信验证码
- **密码登录**：手机号 + 密码

---

## 二、API 接口

### 基础路径

```
/api/v1/c/
```

### 认证相关

#### 1. 发送验证码

```http
POST /api/v1/c/auth/send-code
Content-Type: application/json

{
  "phone": "13800138000"
}
```

**响应**：
```json
{
  "msg": "验证码已发送",
  "data": null
}
```

**错误码**：
- `429` - 发送过于频繁（1分钟内重复发送）或今日次数已达上限

---

#### 2. 登录

```http
POST /api/v1/c/auth/login
Content-Type: application/json

// 验证码登录（默认）
{
  "phone": "13800138000",
  "code": "123456",
  "type": "code"
}

// 密码登录
{
  "phone": "13800138000",
  "password": "your_password",
  "type": "password"
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user_id": 123,
    "type": "c_end",
    "customer_type": "elderly",
    "name": "张三",
    "phone": "13800138000",
    "status": "active"
  }
}
```

---

#### 3. 刷新 Token

```http
POST /api/v1/c/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "token": "new_access_token",
    "refresh_token": "new_refresh_token"
  }
}
```

---

#### 4. 获取当前用户信息

```http
GET /api/v1/c/auth/me
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "user_id": 123,
    "type": "c_end",
    "customer_type": "elderly",
    "name": "张三",
    "phone": "13800138000",
    "status": "active"
  }
}
```

---

#### 5. 检查 Token 状态（用于预填充）

```http
GET /api/v1/c/auth/check
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "user": {
      "id": 123,
      "phone": "13800138000",
      "role": "c_end"
    },
    "profile": {
      "name": "张三",
      "id_number": "310***********1234",
      "address": "上海市浦东新区xxx路xxx号",
      "user_type": "elderly"
    }
  }
}
```

---

#### 6. 登出

```http
POST /api/v1/c/auth/logout
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": null
}
```

---

#### 7. 快速开通（注册+登录+创建服务请求）

一次性完成注册、登录和创建服务需求的全流程，适用于首次使用的用户。

```http
POST /api/v1/c/auth/quick-start
Content-Type: application/json

{
  "phone": "13800138000",
  "code": "123456",
  "name": "张三",
  "address": "上海市浦东新区xxx路xxx号",
  "latitude": 31.2304,
  "longitude": 121.4737,
  "service_type": "daily_care",
  "description": "需要日常照护服务"
}
```

**参数说明**：
- `phone` - 手机号（必填）
- `code` - 验证码（必填）
- `name` - 姓名（必填）
- `address` - 地址（与经纬度二选一）
- `latitude/longitude` - 经纬度（与地址二选一）
- `service_type` - 服务类型（必填）
- `description` - 服务描述（选填）

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 123,
      "phone": "13800138000",
      "role": "c_end"
    },
    "profile": {
      "name": "张三",
      "address": "上海市浦东新区xxx路xxx号"
    },
    "request": {
      "id": 456,
      "request_no": "REQ20240101001",
      "service_type": "daily_care",
      "status": "pending",
      "created_at": "2024-01-01T10:00:00Z"
    }
  }
}
```

---

### 个人资料

#### 1. 更新个人资料

```http
PUT /api/v1/c/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "张三",
  "id_number": "310101199001011234",
  "address": "上海市浦东新区xxx路xxx号",
  "user_type": "elderly"
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "user_id": 123,
    "customer_type": "elderly",
    "id_card": "310101199001011234",
    "address": "上海市浦东新区xxx路xxx号"
  }
}
```

---

### 地理编码

#### 1. 地址解析（地址 → 经纬度）

```http
POST /api/v1/c/geocode
Content-Type: application/json

{
  "address": "上海市浦东新区世纪大道100号"
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "latitude": 31.2304,
    "longitude": 121.4737,
    "formatted_address": "上海市浦东新区世纪大道100号"
  }
}
```

---

#### 2. 逆地理编码（经纬度 → 地址）

```http
GET /api/v1/c/geocode/reverse?lat=31.2304&lng=121.4737
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "province": "上海市",
    "city": "上海市",
    "district": "浦东新区",
    "address": "上海市浦东新区世纪大道100号"
  }
}
```

---

### 服务请求

#### 1. 获取服务请求列表

```http
GET /api/v1/c/requests?page=1&page_size=10&status=pending
Authorization: Bearer <token>
```

**查询参数**：
- `page` - 页码，默认 1
- `page_size` - 每页数量，默认 10
- `status` - 状态筛选（可选）：`pending`, `processing`, `completed`, `cancelled`

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "request_no": "REQ20240101001",
        "service_type": "daily_care",
        "status": "pending",
        "description": "需要日常照护服务",
        "address": "上海市浦东新区xxx路xxx号",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10
  }
}
```

---

#### 2. 创建服务请求

```http
POST /api/v1/c/requests
Authorization: Bearer <token>
Content-Type: application/json

{
  "service_type": "daily_care",
  "description": "需要日常照护服务",
  "address": "上海市浦东新区xxx路xxx号",
  "latitude": 31.2304,
  "longitude": 121.4737,
  "contact_name": "张三",
  "contact_phone": "13800138000"
}
```

**服务类型**：
- `daily_care` - 日常照护
- `medical_care` - 医疗护理
- `rehabilitation` - 康复训练
- `mental_care` - 心理关怀
- `housekeeping` - 家政服务
- `meal_service` - 助餐服务
- `other` - 其他

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "request_no": "REQ20240101001",
    "service_type": "daily_care",
    "status": "pending",
    "created_at": "2024-01-01T10:00:00Z"
  }
}
```

---

#### 3. 获取服务请求详情

```http
GET /api/v1/c/requests/:id
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "request_no": "REQ20240101001",
    "service_type": "daily_care",
    "status": "processing",
    "description": "需要日常照护服务",
    "address": "上海市浦东新区xxx路xxx号",
    "contact_name": "张三",
    "contact_phone": "13800138000",
    "assigned_staff": {
      "id": 10,
      "name": "李四",
      "phone": "13900139000"
    },
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T11:00:00Z"
  }
}
```

---

#### 4. 取消服务请求

```http
POST /api/v1/c/requests/:id/cancel
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "临时有事，需要取消"
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": null
}
```

---

#### 5. 评价服务

```http
POST /api/v1/c/requests/:id/rate
Authorization: Bearer <token>
Content-Type: application/json

{
  "rating": 5,
  "comment": "服务很好，非常满意"
}
```

**参数说明**：
- `rating` - 评分，1-5 分
- `comment` - 评价内容（选填）

**响应**：
```json
{
  "msg": "ok",
  "data": null
}
```

---

### 站点匹配

#### 匹配最近的服务站点

```http
POST /api/v1/c/stations/match
Content-Type: application/json

{
  "latitude": 31.2304,
  "longitude": 121.4737
}
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "station": {
      "id": 1,
      "name": "浦东新区养老服务中心",
      "address": "上海市浦东新区xxx路xxx号",
      "phone": "021-12345678",
      "distance": 1.5
    }
  }
}
```

---

### 内容展示

#### 1. 获取轮播图

```http
GET /api/v1/c/banners
```

**响应**：
```json
{
  "msg": "ok",
  "data": [
    {
      "id": 1,
      "title": "欢迎使用养老服务平台",
      "image_url": "https://example.com/banner1.jpg",
      "link_url": "https://example.com/detail",
      "sort": 1
    }
  ]
}
```

---

### 通知消息

#### 1. 获取通知列表

```http
GET /api/v1/c/notifications?page=1&page_size=10
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "服务请求已受理",
        "content": "您的服务请求 REQ20240101001 已被受理",
        "type": "request",
        "is_read": false,
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 1,
    "unread_count": 1
  }
}
```

---

#### 2. 标记通知已读

```http
POST /api/v1/c/notifications/:id/read
Authorization: Bearer <token>
```

**响应**：
```json
{
  "msg": "ok",
  "data": null
}
```

---

### 文件上传

#### 上传文件

```http
POST /api/v1/c/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <binary>
```

**响应**：
```json
{
  "msg": "ok",
  "data": {
    "url": "https://example.com/uploads/xxx.jpg",
    "filename": "xxx.jpg"
  }
}
```

---

## 三、前端实现指南

### 1. 类型定义

```typescript
// types/user.ts
export interface CEndUser {
  user_id: number;
  type: 'c_end';
  customer_type: 'elderly' | 'family';
  name: string;
  phone: string;
  status: string;
}

export interface CustomerProfile {
  name: string;
  id_number?: string;
  address?: string;
  user_type: 'elderly' | 'family';
}

// types/request.ts
export type ServiceType =
  | 'daily_care'
  | 'medical_care'
  | 'rehabilitation'
  | 'mental_care'
  | 'housekeeping'
  | 'meal_service'
  | 'other';

export type RequestStatus =
  | 'pending'
  | 'processing'
  | 'completed'
  | 'cancelled';

export interface ServiceRequest {
  id: number;
  request_no: string;
  service_type: ServiceType;
  status: RequestStatus;
  description?: string;
  address: string;
  contact_name: string;
  contact_phone: string;
  assigned_staff?: {
    id: number;
    name: string;
    phone: string;
  };
  created_at: string;
  updated_at: string;
}
```

---

### 2. API 封装

```typescript
// api/auth.ts
import request from '@/utils/request';

export function sendCode(phone: string) {
  return request.post('/c/auth/send-code', { phone });
}

export function login(data: { phone: string; code?: string; password?: string; type?: string }) {
  return request.post('/c/auth/login', data);
}

export function refresh(refreshToken: string) {
  return request.post('/c/auth/refresh', { refresh_token: refreshToken });
}

export function getUserInfo() {
  return request.get('/c/auth/me');
}

export function checkToken() {
  return request.get('/c/auth/check');
}

export function logout() {
  return request.post('/c/auth/logout');
}

export function quickStart(data: {
  phone: string;
  code: string;
  name: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  service_type: string;
  description?: string;
}) {
  return request.post('/c/auth/quick-start', data);
}

// api/request.ts
export function getRequests(params: { page?: number; page_size?: number; status?: string }) {
  return request.get('/c/requests', { params });
}

export function createRequest(data: {
  service_type: string;
  description?: string;
  address: string;
  latitude?: number;
  longitude?: number;
  contact_name: string;
  contact_phone: string;
}) {
  return request.post('/c/requests', data);
}

export function getRequestDetail(id: number) {
  return request.get(`/c/requests/${id}`);
}

export function cancelRequest(id: number, reason?: string) {
  return request.post(`/c/requests/${id}/cancel`, { reason });
}

export function rateRequest(id: number, rating: number, comment?: string) {
  return request.post(`/c/requests/${id}/rate`, { rating, comment });
}
```

---

### 3. 用户状态管理

```typescript
// stores/user.ts
import { defineStore } from 'pinia';
import { getUserInfo, logout } from '@/api/auth';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('c_token') || '',
    refreshToken: localStorage.getItem('c_refresh_token') || '',
    userInfo: null as CEndUser | null,
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    customerType: (state) => state.userInfo?.customer_type,
  },

  actions: {
    setTokens(token: string, refreshToken: string) {
      this.token = token;
      this.refreshToken = refreshToken;
      localStorage.setItem('c_token', token);
      localStorage.setItem('c_refresh_token', refreshToken);
    },

    async fetchUserInfo() {
      try {
        const { data } = await getUserInfo();
        this.userInfo = data;
        return data;
      } catch (error) {
        this.clearAuth();
        throw error;
      }
    },

    async logout() {
      try {
        await logout();
      } finally {
        this.clearAuth();
      }
    },

    clearAuth() {
      this.token = '';
      this.refreshToken = '';
      this.userInfo = null;
      localStorage.removeItem('c_token');
      localStorage.removeItem('c_refresh_token');
    },
  },
});
```

---

### 4. 请求拦截器

```typescript
// utils/request.ts
import axios from 'axios';
import { useUserStore } from '@/stores/user';
import { refresh } from '@/api/auth';

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
});

// 请求拦截器
request.interceptors.request.use((config) => {
  const userStore = useUserStore();
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`;
  }
  return config;
});

// 响应拦截器
request.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    const userStore = useUserStore();

    // Token 过期，尝试刷新
    if (error.response?.status === 401 && userStore.refreshToken) {
      try {
        const { data } = await refresh(userStore.refreshToken);
        userStore.setTokens(data.token, data.refresh_token);

        // 重试原请求
        error.config.headers.Authorization = `Bearer ${data.token}`;
        return request(error.config);
      } catch {
        userStore.clearAuth();
        // 跳转登录页
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

export default request;
```

---

### 5. 验证码倒计时组件

```vue
<!-- components/SmsCodeButton.vue -->
<template>
  <el-button
    :disabled="countdown > 0"
    :loading="loading"
    @click="handleSend"
  >
    {{ countdown > 0 ? `${countdown}s 后重试` : '获取验证码' }}
  </el-button>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { sendCode } from '@/api/auth';
import { ElMessage } from 'element-plus';

const props = defineProps<{
  phone: string;
}>();

const loading = ref(false);
const countdown = ref(0);
let timer: number | null = null;

async function handleSend() {
  if (!props.phone) {
    ElMessage.warning('请输入手机号');
    return;
  }

  loading.value = true;
  try {
    await sendCode(props.phone);
    ElMessage.success('验证码已发送');
    startCountdown();
  } catch (error: any) {
    if (error.response?.status === 429) {
      ElMessage.error(error.response.data.msg || '发送过于频繁');
    } else {
      ElMessage.error('发送失败，请重试');
    }
  } finally {
    loading.value = false;
  }
}

function startCountdown() {
  countdown.value = 60;
  timer = window.setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0 && timer) {
      clearInterval(timer);
      timer = null;
    }
  }, 1000);
}
</script>
```

---

### 6. 服务类型映射

```typescript
// constants/serviceType.ts
export const SERVICE_TYPE_MAP: Record<string, string> = {
  daily_care: '日常照护',
  medical_care: '医疗护理',
  rehabilitation: '康复训练',
  mental_care: '心理关怀',
  housekeeping: '家政服务',
  meal_service: '助餐服务',
  other: '其他',
};

export const REQUEST_STATUS_MAP: Record<string, { label: string; color: string }> = {
  pending: { label: '待处理', color: 'warning' },
  processing: { label: '处理中', color: 'primary' },
  completed: { label: '已完成', color: 'success' },
  cancelled: { label: '已取消', color: 'info' },
};
```

---

## 四、错误码说明

| HTTP 状态码 | 说明 |
|-------------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 无效 |
| 403 | 无权限（用户被禁用等） |
| 404 | 资源不存在 |
| 429 | 请求过于频繁（验证码发送限制） |
| 500 | 服务器内部错误 |

---

## 五、注意事项

1. **Token 存储**：建议使用 `localStorage` 存储 Token，注意 XSS 防护
2. **Token 刷新**：Access Token 有效期较短，需实现自动刷新机制
3. **地理位置**：创建服务请求时，优先使用经纬度，地址作为备选
4. **验证码限制**：
   - 同一手机号 1 分钟内只能发送 1 次
   - 同一手机号每天最多发送 10 次
5. **文件上传**：支持图片格式（jpg、png、gif），单文件最大 5MB

---

## 六、测试账号

| 用户类型 | 手机号 | 密码 |
|----------|--------|------|
| 老年人 | 13800000010 | Test@123 |
| 家属 | 13800000011 | Test@123 |

> 注：验证码登录时，开发环境可使用固定验证码 `123456`
