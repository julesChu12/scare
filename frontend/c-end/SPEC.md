# sCare C端用户小程序/H5 技术规格说明书

## 📋 项目概述

### 项目信息
- **项目名称**: sCare 社区养老平台 - C端用户端
- **目标用户**: 老年人(elderly) + 家属(family)
- **平台**: 微信小程序 / H5（优先微信小程序）
- **技术栈**: 待定（推荐 uni-app / Taro / 原生微信小程序）
- **后端**: Go + Gin（已有API）
- **认证方式**: JWT Token
- **状态**: MVP 规划
- **日期**: 2026-01-19

### 核心定位
- ✅ **无复杂权限系统**：C端用户统一身份，无角色区分
- ✅ **用户属性标签**：elderly/family 只是用户类型标签，不影响权限
- ✅ **简单易用**：针对老年人群体，界面简洁、操作简单
- ✅ **核心功能**：发起需求 → 查看进度 → 接收通知

---

## 🎯 MVP 功能范围

### 用户类型说明
| 用户类型 | 说明 | 核心需求 |
|---------|------|---------|
| **老年人(elderly)** | 主要服务对象 | 发起服务需求、查看需求状态、接收服务通知 |
| **家属(family)** | 老年人的监护人 | 查看关联老年人的需求、接收进度通知 |

**重要**：用户类型只是属性标签，**不控制权限**。所有C端用户的功能权限完全相同。

---

## 📱 页面结构（MVP）

```
/
├── /pages/auth/
│   └── login.vue           # 登录页（手机号 + 验证码）
│
├── /pages/index/
│   └── index.vue           # 首页（需求发起入口 + 我的需求）
│
├── /pages/request/
│   ├── create.vue          # 发起需求（表单）
│   ├── list.vue            # 我的需求列表
│   └── detail.vue          # 需求详情（状态、进度、联系方式）
│
└── /pages/profile/
    └── index.vue           # 个人中心（基本信息、退出登录）
```

---

## 🎨 核心页面设计

### 1. 登录页 (`/pages/auth/login`)

**设计原则**：简单、大字体、易操作

**功能**：
- 手机号输入（大号输入框）
- 验证码登录（发送验证码按钮）
- 自动登录（记住登录态）

**交互流程**：
```
1. 用户输入手机号
2. 点击"获取验证码" → 调用 POST /api/v1/c/auth/sms
3. 输入验证码
4. 点击"登录" → 调用 POST /api/v1/c/auth/login
5. 接收 { access_token, user: { id, name, phone, type } }
6. 存储 token 到本地
7. 跳转到首页
```

**UI要点**：
- ✅ 字体大（≥18px）
- ✅ 按钮大（高度≥50px）
- ✅ 高对比度配色
- ✅ 错误提示清晰明确

---

### 2. 首页 (`/pages/index/index`)

**布局**：
```
┌─────────────────────────┐
│  欢迎，张奶奶 👋          │
│  （用户头像 + 姓名）       │
├─────────────────────────┤
│                         │
│   🏠 我要服务            │
│   [大按钮：发起服务需求]   │
│                         │
├─────────────────────────┤
│  📋 我的需求             │
│  ┌─────────────────┐   │
│  │ 陪同就医         │   │
│  │ 进行中 → 查看    │   │
│  └─────────────────┘   │
│  ┌─────────────────┐   │
│  │ 上门助洁         │   │
│  │ 已完成           │   │
│  └─────────────────┘   │
└─────────────────────────┘
```

**功能**：
- 用户信息展示（头像、姓名、用户类型标签）
- **发起服务需求**（大按钮，跳转到创建页面）
- 我的需求列表（最近5条，点击查看详情）
- 底部导航（首页、我的）

---

### 3. 发起需求页 (`/pages/request/create`)

**表单字段**（简化版）：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| 服务类型 | 单选 | ✅ | 陪同就医、上门助洁、日常照料、紧急呼叫 |
| 服务描述 | 文本框 | ✅ | 详细说明需求（大文本框） |
| 服务地址 | 地图选点 | ✅ | 调用微信地图API选择位置 |
| 联系电话 | 输入框 | ✅ | 默认填充用户手机号 |
| 预约时间 | 时间选择器 | ❌ | 可选，默认"尽快" |

**交互流程**：
```
1. 用户选择服务类型（图标+文字）
2. 填写服务描述（语音输入支持）
3. 选择服务地址（调用微信选择位置API）
4. 确认联系电话（默认已填）
5. 点击"提交需求"
6. 调用 POST /api/v1/c/requests
7. 成功后跳转到需求详情页
```

**UI要点**：
- ✅ 表单分步骤，一页只显示一个问题
- ✅ 支持语音输入（微信API）
- ✅ 地图大图展示
- ✅ 提交按钮醒目

---

### 4. 需求列表页 (`/pages/request/list`)

**列表项设计**：
```
┌─────────────────────────┐
│ 🏥 陪同就医               │
│ 2026-01-19 10:30       │
│ 【进行中】工作人员已接单   │
│ 联系电话: 138****1234    │
│         [查看详情 →]     │
└─────────────────────────┘
```

**状态标签**：
- 🟡 待分配（黄色）
- 🔵 进行中（蓝色）
- 🟢 已完成（绿色）
- 🔴 已取消（灰色）

**功能**：
- 按时间倒序展示
- 状态筛选（全部、进行中、已完成）
- 下拉刷新
- 点击卡片查看详情

---

### 5. 需求详情页 (`/pages/request/detail`)

**页面布局**：
```
┌─────────────────────────┐
│  服务类型: 陪同就医        │
│  状态: 🔵 进行中          │
├─────────────────────────┤
│  📍 服务地址              │
│  北京市朝阳区xxx路xx号     │
│  [查看地图]              │
├─────────────────────────┤
│  📝 需求描述              │
│  需要陪同去医院看病，行动  │
│  不便需要轮椅...          │
├─────────────────────────┤
│  👤 联系方式              │
│  张奶奶 138****1234      │
├─────────────────────────┤
│  🕐 提交时间              │
│  2026-01-19 10:30       │
├─────────────────────────┤
│  👷 服务人员（如果已分配）  │
│  李师傅 156****5678       │
│  [一键拨号]              │
├─────────────────────────┤
│  📸 服务照片（如果已完成）  │
│  [图片展示]              │
└─────────────────────────┘

[取消需求] （仅待分配状态可点击）
```

**功能**：
- 展示完整需求信息
- 地图展示服务地址
- 一键拨号（联系服务人员）
- 查看服务照片（完成后）
- 取消需求（仅待分配状态）

---

### 6. 个人中心 (`/pages/profile/index`)

**功能**：
- 用户信息展示（头像、姓名、手机号、用户类型）
- 修改个人信息（姓名、头像）
- 退出登录

**简化版布局**：
```
┌─────────────────────────┐
│   [头像]                 │
│   张奶奶                  │
│   138****1234            │
│   用户类型: 老年人         │
├─────────────────────────┤
│ > 我的需求 (12)          │
│ > 个人信息               │
│ > 关于我们               │
│ > 退出登录               │
└─────────────────────────┘
```

---

## 🔌 后端接口契约

### 1. 认证接口

#### 1.1 发送验证码
```
POST /api/v1/c/auth/sms

Request:
{
  "phone": "13812341234"
}

Response:
{
  "code": 0,
  "msg": "验证码已发送",
  "data": {
    "expires_in": 300  // 5分钟有效期
  }
}
```

#### 1.2 登录
```
POST /api/v1/c/auth/login

Request:
{
  "phone": "13812341234",
  "sms_code": "123456"
}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "access_token": "eyJhbGc...",
    "user": {
      "id": 1,
      "name": "张奶奶",
      "phone": "13812341234",
      "type": "elderly",  // elderly / family
      "avatar": "https://..."
    }
  }
}
```

#### 1.3 获取用户信息
```
GET /api/v1/c/auth/me
Headers: Authorization: Bearer {token}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "user": {
      "id": 1,
      "name": "张奶奶",
      "phone": "13812341234",
      "type": "elderly",
      "avatar": "https://...",
      "created_at": "2026-01-01T00:00:00Z"
    }
  }
}
```

---

### 2. 需求管理接口

#### 2.1 创建需求
```
POST /api/v1/c/requests
Headers: Authorization: Bearer {token}

Request:
{
  "service_type": "accompany_medical",  // 服务类型枚举
  "description": "需要陪同去医院看病...",
  "address": "北京市朝阳区xxx路xx号",
  "latitude": 39.9042,
  "longitude": 116.4074,
  "contact_phone": "13812341234",
  "scheduled_time": "2026-01-20T10:00:00Z"  // 可选
}

Response:
{
  "code": 0,
  "msg": "需求提交成功",
  "data": {
    "request_id": 123,
    "status": "pending"
  }
}
```

#### 2.2 查询我的需求列表
```
GET /api/v1/c/requests?status=all&page=1&page_size=10
Headers: Authorization: Bearer {token}

Query Parameters:
- status: all | pending | in_progress | completed | cancelled
- page: 页码（默认1）
- page_size: 每页数量（默认10）

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 20,
    "page": 1,
    "page_size": 10,
    "requests": [
      {
        "id": 123,
        "service_type": "accompany_medical",
        "service_type_name": "陪同就医",
        "description": "需要陪同去医院...",
        "status": "in_progress",
        "status_name": "进行中",
        "address": "北京市朝阳区...",
        "contact_phone": "13812341234",
        "assigned_staff": {  // 如果已分配
          "name": "李师傅",
          "phone": "15612345678"
        },
        "created_at": "2026-01-19T10:30:00Z",
        "updated_at": "2026-01-19T11:00:00Z"
      }
    ]
  }
}
```

#### 2.3 查询需求详情
```
GET /api/v1/c/requests/:id
Headers: Authorization: Bearer {token}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "id": 123,
    "service_type": "accompany_medical",
    "service_type_name": "陪同就医",
    "description": "需要陪同去医院看病...",
    "status": "completed",
    "status_name": "已完成",
    "address": "北京市朝阳区xxx路xx号",
    "latitude": 39.9042,
    "longitude": 116.4074,
    "contact_phone": "13812341234",
    "scheduled_time": "2026-01-20T10:00:00Z",
    "assigned_staff": {
      "name": "李师傅",
      "phone": "15612345678",
      "avatar": "https://..."
    },
    "completion_photos": [  // 完成照片
      "https://...",
      "https://..."
    ],
    "completion_note": "服务已完成，老人安全到家",
    "created_at": "2026-01-19T10:30:00Z",
    "updated_at": "2026-01-19T15:00:00Z",
    "completed_at": "2026-01-19T15:00:00Z"
  }
}
```

#### 2.4 取消需求
```
POST /api/v1/c/requests/:id/cancel
Headers: Authorization: Bearer {token}

Request:
{
  "reason": "时间冲突，需要改期"  // 可选
}

Response:
{
  "code": 0,
  "msg": "需求已取消",
  "data": {
    "request_id": 123,
    "status": "cancelled"
  }
}
```

---

### 3. 通知接口（可选，MVP可暂缓）

#### 3.1 获取通知列表
```
GET /api/v1/c/notifications?page=1&page_size=20
Headers: Authorization: Bearer {token}

Response:
{
  "code": 0,
  "msg": "ok",
  "data": {
    "total": 5,
    "unread_count": 2,
    "notifications": [
      {
        "id": 1,
        "type": "request_assigned",
        "title": "需求已分配",
        "content": "您的陪同就医需求已由李师傅接单",
        "is_read": false,
        "created_at": "2026-01-19T11:00:00Z"
      }
    ]
  }
}
```

#### 3.2 标记已读
```
POST /api/v1/c/notifications/:id/read
Headers: Authorization: Bearer {token}

Response:
{
  "code": 0,
  "msg": "ok"
}
```

---

## 🎨 UI/UX 设计原则

### 针对老年人用户的设计
1. **大字体**：正文 ≥16px，标题 ≥20px
2. **大按钮**：高度 ≥50px，宽度尽量占满
3. **高对比度**：黑底白字或白底黑字，避免灰色
4. **简化操作**：每页最多3个操作，步骤清晰
5. **语音支持**：关键输入框支持语音输入
6. **一键拨号**：显示电话号码时提供直接拨号按钮

### 配色建议
- **主色**：#1890FF（蓝色，信任感）
- **成功**：#52C41A（绿色）
- **警告**：#FAAD14（橙色）
- **错误**：#F5222D（红色）
- **文字**：#262626（深灰）
- **背景**：#F5F5F5（浅灰）

---

## 📂 推荐技术选型

### 方案对比

| 技术方案 | 优势 | 劣势 | 推荐度 |
|---------|------|------|--------|
| **uni-app** | 一套代码多端、开发效率高、生态成熟 | 性能略逊于原生 | ⭐⭐⭐⭐⭐ |
| **Taro** | React语法、跨端支持、京东团队维护 | 学习成本较高 | ⭐⭐⭐⭐ |
| **原生微信小程序** | 性能最优、API完整 | 只能在微信生态 | ⭐⭐⭐ |
| **Vue3 H5** | 开发简单、调试方便 | 微信内体验差 | ⭐⭐ |

**推荐**：**uni-app** (Vue 3 + TypeScript)

---

## 🔐 安全与权限

### C端权限模型（简化）
- ✅ **无角色权限**：所有C端用户权限完全相同
- ✅ **数据隔离**：用户只能查看/操作自己的数据
- ✅ **后端校验**：所有接口检查 `user_id` 是否匹配

### 数据安全
```go
// 后端接口示例（确保数据隔离）
func (h *RequestHandler) Get(c *gin.Context) {
    requestID := c.Param("id")
    userID := c.GetInt64("user_id")  // 从 JWT 获取

    request, err := h.requestRepo.GetByID(requestID)
    if err != nil {
        RespondError(c, http.StatusNotFound, "需求不存在")
        return
    }

    // 关键检查：确保用户只能访问自己的需求
    if request.UserID != userID {
        RespondError(c, http.StatusForbidden, "无权访问")
        return
    }

    Respond(c, http.StatusOK, "ok", request)
}
```

---

## 📱 微信小程序特殊配置

### 1. 必需权限
```json
// app.json
{
  "permission": {
    "scope.userLocation": {
      "desc": "您的位置信息将用于服务地址定位"
    }
  },
  "requiredPrivateInfos": [
    "getLocation",
    "chooseLocation"
  ]
}
```

### 2. 关键API使用
```javascript
// 选择位置
wx.chooseLocation({
  success: (res) => {
    const { name, address, latitude, longitude } = res;
    // 保存位置信息
  }
})

// 拨打电话
wx.makePhoneCall({
  phoneNumber: '13812341234'
})

// 语音输入（需要插件）
```

---

## ✅ MVP 实现 Checklist

### 第一阶段：基础功能（1-2天）
- [ ] 登录页（手机号 + 验证码）
- [ ] 首页（用户信息 + 发起需求入口）
- [ ] 发起需求（简化表单，必填字段）
- [ ] 需求列表（基础列表）

### 第二阶段：完善体验（1-2天）
- [ ] 需求详情页（完整信息展示）
- [ ] 地图选点（微信API集成）
- [ ] 一键拨号
- [ ] 取消需求

### 第三阶段：优化迭代（1-2天）
- [ ] 语音输入支持
- [ ] 图片上传（服务照片）
- [ ] 通知推送（可选）
- [ ] 个人中心完善

---

## 🎯 MVP 关键指标

### 用户体验
- ✅ 登录成功率 > 95%
- ✅ 需求提交成功率 > 90%
- ✅ 页面加载时间 < 2s
- ✅ 操作步骤 ≤ 3步（发起需求）

### 功能完整性
- ✅ 用户可独立完成需求发起
- ✅ 用户可查看需求实时状态
- ✅ 用户可联系服务人员

---

## 📝 数据字典

### 服务类型枚举
```typescript
enum ServiceType {
  ACCOMPANY_MEDICAL = 'accompany_medical',    // 陪同就医
  HOME_CLEANING = 'home_cleaning',            // 上门助洁
  DAILY_CARE = 'daily_care',                  // 日常照料
  EMERGENCY_CALL = 'emergency_call'           // 紧急呼叫
}
```

### 需求状态枚举
```typescript
enum RequestStatus {
  PENDING = 'pending',           // 待分配
  IN_PROGRESS = 'in_progress',   // 进行中
  COMPLETED = 'completed',       // 已完成
  CANCELLED = 'cancelled'        // 已取消
}
```

### 用户类型枚举
```typescript
enum UserType {
  ELDERLY = 'elderly',  // 老年人
  FAMILY = 'family'     // 家属
}
```

---

## 🚀 部署与发布

### 微信小程序发布流程
1. 申请小程序账号（企业主体）
2. 配置服务器域名（HTTPS必需）
3. 上传代码到微信平台
4. 提交审核
5. 审核通过后发布

### 服务器域名要求
```
request合法域名:
https://api.scare.example.com

uploadFile合法域名:
https://upload.scare.example.com

downloadFile合法域名:
https://static.scare.example.com
```

---

## 📌 重要提醒

### 与B端的关键区别
1. **无权限系统**：C端用户不需要角色、权限管理
2. **无后台管理**：C端是纯用户操作界面
3. **数据隔离**：用户只能看到自己的数据
4. **简化操作**：针对老年人优化，步骤尽量少

### 后端接口复用
- ✅ 认证接口：C端和B端分离（`/api/v1/c/auth/*` vs `/api/v1/b/auth/*`）
- ✅ 业务接口：部分复用（如查询需求详情）
- ✅ 数据库：共用同一套表结构

---

**Spec 已锁定，可交付给其他AI实现MVP。**

**建议技术栈**: uni-app (Vue 3 + TypeScript) + Vant Weapp (UI组件库)
