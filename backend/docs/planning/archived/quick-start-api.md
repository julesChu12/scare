# sCare C端快速开通 API 文档

**版本**: v1.0
**Base URL**: `https://api.scare.com/api/v1`
**协议**: HTTPS
**认证方式**: Bearer Token (JWT)

---

## 📚 目录

1. [认证接口](#认证接口)
2. [快速通道接口](#快速通道接口)
3. [C端系统接口](#c端系统接口)
4. [错误码说明](#错误码说明)

---

## 🔐 认证接口

### 1.1 发送验证码

**端点**: `POST /c/auth/send-code`
**认证**: 无需认证
**描述**: 向指定手机号发送6位数字验证码，有效期5分钟

**请求体**:
```json
{
  "phone": "13800000001"
}
```

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "expires_in": 300
  }
}
```

**错误响应**:
```json
// 429 - 分钟级限制
{
  "msg": "发送过于频繁，请1分钟后再试",
  "data": null
}

// 429 - 每日限制
{
  "msg": "今日发送次数已达上限（10次）",
  "data": null
}

// 500 - 生产环境
{
  "msg": "短信服务暂未开通，请联系工作人员",
  "data": null
}
```

**频率限制**:
- 同一手机号：60秒内最多发送1次
- 同一手机号：每天最多发送10次

---

### 1.2 快速开通（注册+登录+创建请求）

**端点**: `POST /c/auth/quick-start`
**认证**: 无需认证
**描述**: 一键完成注册、登录、创建服务请求

**请求体**:
```json
{
  "phone": "13800000001",
  "code": "123456",
  "name": "张大爷",
  "service_request": {
    "service_type": "bath_assistance",
    "station_id": 1,
    "address": "北京市朝阳区幸福小区1号楼101",
    "lat": 39.9042,
    "lng": 116.4074,
    "note": "腿脚不便，需要轮椅"
  }
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 是 | 手机号，11位数字 |
| code | string | 是 | 验证码，6位数字 |
| name | string | 是 | 用户姓名 |
| service_request | object | 是 | 服务请求信息 |
| service_request.service_type | string | 是 | 服务类型：bath_assistance(助浴)/meal_assistance(助餐)/nursing(护理)/rehabilitation(康复)/other(其他) |
| service_request.station_id | number | 是 | 服务站点ID（二维码携带） |
| service_request.address | string | 是 | 服务地址 |
| service_request.lat | number | 否 | 纬度（GPS定位成功时提供） |
| service_request.lng | number | 否 | 经度（GPS定位成功时提供） |
| service_request.note | string | 否 | 补充说明 |

**成功响应** (200 OK - 新用户):
```json
{
  "msg": "ok",
  "data": {
    "is_new_user": true,
    "token": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 3600
    },
    "user": {
      "id": 14,
      "phone": "13800000001",
      "name": "张大爷",
      "status": "active",
      "created_at": "2026-01-20T10:30:00Z"
    },
    "request": {
      "id": 1,
      "request_no": "R202601200001",
      "service_type": "bath_assistance",
      "status": "pending",
      "station_id": 1,
      "address": "北京市朝阳区幸福小区1号楼101",
      "note": "腿脚不便，需要轮椅",
      "created_at": "2026-01-20T10:30:00Z"
    }
  }
}
```

**成功响应** (200 OK - 老用户):
```json
{
  "msg": "ok",
  "data": {
    "is_new_user": false,
    "token": {...},
    "user": {
      "id": 8,
      "phone": "13800000008",
      "name": "张大爷",
      "customer_type": "elderly",
      "status": "active"
    },
    "request": {...}
  }
}
```

**错误响应**:
```json
// 400 - 验证码错误
{
  "msg": "验证码错误或已过期",
  "data": null
}

// 400 - 参数错误
{
  "msg": "invalid payload",
  "data": null
}
```

---

### 1.3 检查用户状态

**端点**: `GET /c/auth/check`
**认证**: 需要Token
**描述**: 检查当前token状态，获取用户信息（用于预填充表单）

**请求头**:
```
Authorization: Bearer {access_token}
```

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "user": {
      "id": 8,
      "phone": "13800000008",
      "name": "张大爷",
      "customer_type": "elderly",
      "status": "active"
    },
    "last_request": {
      "id": 10,
      "address": "北京市朝阳区幸福小区1号楼101",
      "service_type": "bath_assistance",
      "lat": 39.9042,
      "lng": 116.4074
    }
  }
}
```

**错误响应**:
```json
// 401 - Token无效
{
  "msg": "unauthorized",
  "data": null
}
```

---

## 🚀 快速通道接口

### 2.1 地址解析（高德API代理）

**端点**: `POST /c/geocode`
**认证**: 无需认证
**描述**: 将地址字符串解析为经纬度，并匹配最近服务站

**请求体**:
```json
{
  "address": "北京市朝阳区幸福小区1号楼101"
}
```

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "lat": 39.9042,
    "lng": 116.4074,
    "formatted_address": "北京市朝阳区幸福街道幸福小区1号楼101室",
    "matched_station": {
      "id": 1,
      "name": "朝阳区幸福街道服务站",
      "distance": 500
    }
  }
}
```

**无覆盖响应** (404):
```json
{
  "msg": "您所在区域暂无服务站，已记录您的需求",
  "data": {
    "lat": 39.9042,
    "lng": 116.4074,
    "formatted_address": "北京市朝阳区幸福街道幸福小区1号楼101室",
    "no_coverage": true
  }
}
```

**错误响应**:
```json
// 400 - 地址解析失败
{
  "msg": "地址解析失败，请检查地址是否正确",
  "data": null
}
```

---

## 👤 C端系统接口

### 3.1 更新个人资料

**端点**: `PUT /c/profile`
**认证**: 需要Token
**描述**: 更新用户的客户档案信息

**请求头**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**请求体**:
```json
{
  "customer_type": "elderly",
  "gender": "male",
  "birth_date": "1950-05-15",
  "health_status": "良好",
  "disability_level": "自理",
  "emergency_contact": {
    "name": "张小明",
    "phone": "13900000001",
    "relation": "子女"
  }
}
```

**字段说明**:
| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_type | string | 否 | 客户类型：elderly(老年人)/disabled(残障人士)/pregnant(孕妇)/child(儿童)/other(其他) |
| gender | string | 否 | 性别：male/female |
| birth_date | string | 否 | 出生日期，格式: YYYY-MM-DD |
| health_status | string | 否 | 健康状况 |
| disability_level | string | 否 | 失能等级：自理/轻度失能/中度失能/重度失能 |
| emergency_contact | object | 否 | 紧急联系人 |
| emergency_contact.name | string | 是 | 联系人姓名 |
| emergency_contact.phone | string | 是 | 联系人电话 |
| emergency_contact.relation | string | 是 | 关系 |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "id": 8,
    "user_id": 8,
    "customer_type": "elderly",
    "gender": "male",
    "birth_date": "1950-05-15",
    "health_status": "良好",
    "disability_level": "自理",
    "emergency_contact": {
      "name": "张小明",
      "phone": "13900000001",
      "relation": "子女"
    },
    "updated_at": "2026-01-20T11:00:00Z"
  }
}
```

---

### 3.2 创建服务请求

**端点**: `POST /c/requests`
**认证**: 需要Token
**描述**: C端系统内创建服务请求（带照片上传等完整功能）

**请求头**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**请求体**:
```json
{
  "service_type": "bath_assistance",
  "lat": 39.9042,
  "lng": 116.4074,
  "address": "北京市朝阳区幸福小区1号楼101",
  "contact_name": "张大爷",
  "contact_phone": "13800000001",
  "images": [
    "https://storage.scare.com/uploads/xxx.jpg"
  ],
  "note": "腿脚不便，需要轮椅"
}
```

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "id": 2,
    "request_no": "R202601200002",
    "user_id": 8,
    "service_type": "bath_assistance",
    "status": "pending",
    "lat": 39.9042,
    "lng": 116.4074,
    "address": "北京市朝阳区幸福小区1号楼101",
    "contact_name": "张大爷",
    "contact_phone": "13800000001",
    "images": ["https://storage.scare.com/uploads/xxx.jpg"],
    "note": "腿脚不便，需要轮椅",
    "created_at": "2026-01-20T11:30:00Z"
  }
}
```

---

### 3.3 查看服务请求列表

**端点**: `GET /c/requests`
**认证**: 需要Token
**描述**: 查看当前用户的所有服务请求

**请求头**:
```
Authorization: Bearer {access_token}
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| status | string | 否 | 过滤状态：pending/dispatched/claimed/in_progress/completed/cancelled |
| page | number | 否 | 页码，默认1 |
| page_size | number | 否 | 每页数量，默认20 |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "total": 10,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "request_no": "R202601200001",
        "service_type": "bath_assistance",
        "status": "completed",
        "address": "北京市朝阳区幸福小区1号楼101",
        "created_at": "2026-01-20T10:30:00Z",
        "updated_at": "2026-01-20T15:00:00Z"
      },
      {
        "id": 2,
        "request_no": "R202601200002",
        "service_type": "meal_assistance",
        "status": "in_progress",
        "address": "北京市朝阳区幸福小区1号楼101",
        "created_at": "2026-01-20T11:30:00Z",
        "updated_at": "2026-01-20T12:00:00Z"
      }
    ]
  }
}
```

---

### 3.4 查看请求详情

**端点**: `GET /c/requests/:id`
**认证**: 需要Token
**描述**: 查看指定服务请求的详细信息

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | number | 请求ID |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "request_no": "R202601200001",
    "user_id": 8,
    "service_type": "bath_assistance",
    "status": "completed",
    "station_id": 1,
    "station_name": "朝阳区幸福街道服务站",
    "address": "北京市朝阳区幸福小区1号楼101",
    "contact_name": "张大爷",
    "contact_phone": "13800000001",
    "note": "腿脚不便，需要轮椅",
    "images": [],
    "assigned_staff": {
      "id": 5,
      "name": "刘师傅",
      "phone": "13800000005"
    },
    "created_at": "2026-01-20T10:30:00Z",
    "updated_at": "2026-01-20T15:00:00Z",
    "completed_at": "2026-01-20T15:00:00Z"
  }
}
```

---

### 3.5 取消服务请求

**端点**: `POST /c/requests/:id/cancel`
**认证**: 需要Token
**描述**: 取消指定的服务请求（仅pending/dispatched状态可取消）

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | number | 请求ID |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "id": 1,
    "status": "cancelled",
    "updated_at": "2026-01-20T12:00:00Z"
  }
}
```

**错误响应**:
```json
// 400 - 不可取消
{
  "msg": "该请求当前状态不可取消",
  "data": null
}
```

---

### 3.6 获取通知列表

**端点**: `GET /c/notifications`
**认证**: 需要Token
**描述**: 获取当前用户的通知列表

**请求头**:
```
Authorization: Bearer {access_token}
```

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | number | 否 | 页码，默认1 |
| page_size | number | 否 | 每页数量，默认20 |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "total": 5,
    "unread_count": 2,
    "items": [
      {
        "id": 1,
        "type": "request_accepted",
        "title": "服务请求已接单",
        "content": "您的助浴服务请求已被刘师傅接单",
        "is_read": false,
        "created_at": "2026-01-20T11:00:00Z"
      },
      {
        "id": 2,
        "type": "request_completed",
        "title": "服务已完成",
        "content": "您的助浴服务已完成，请评价",
        "is_read": true,
        "created_at": "2026-01-20T15:00:00Z"
      }
    ]
  }
}
```

---

### 3.7 标记通知已读

**端点**: `POST /c/notifications/:id/read`
**认证**: 需要Token
**描述**: 标记指定通知为已读

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
| 参数 | 类型 | 说明 |
|------|------|------|
| id | number | 通知ID |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": null
}
```

---

### 3.8 文件上传

**端点**: `POST /c/upload`
**认证**: 需要Token
**描述**: 上传图片等文件

**请求头**:
```
Authorization: Bearer {access_token}
Content-Type: multipart/form-data
```

**表单数据**:
| 字段 | 类型 | 说明 |
|------|------|------|
| file | file | 文件，最大5MB，支持jpg/png/jpeg |

**成功响应** (200 OK):
```json
{
  "msg": "ok",
  "data": {
    "url": "https://storage.scare.com/uploads/2026/01/20/xxx.jpg",
    "filename": "xxx.jpg",
    "size": 102400
  }
}
```

---

## ❌ 错误码说明

### HTTP状态码
| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或token无效 |
| 403 | 无权限访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

### 业务错误码
| 错误码 | 说明 |
|--------|------|
| INVALID_PHONE | 手机号格式错误 |
| CODE_INVALID | 验证码错误或已过期 |
| RATE_LIMIT_MINUTE | 发送过于频繁（分钟级） |
| RATE_LIMIT_DAILY | 今日次数已达上限 |
| SMS_NOT_ENABLED | 短信服务未开通 |
| USER_EXISTS | 用户已存在 |
| USER_NOT_FOUND | 用户不存在 |
| REQUEST_NOT_FOUND | 请求不存在 |
| REQUEST_CANNOT_CANCEL | 请求当前状态不可取消 |
| NO_STATION_COVERAGE | 无服务站覆盖 |
| GEOCODE_FAILED | 地址解析失败 |

---

## 📝 开发环境说明

### 验证码Mock
开发环境下，验证码将打印到服务器控制台：
```
========== 短信验证码（Mock） ==========
手机号：13800000001
验证码：123456
有效期：5分钟
======================================
```

### 测试账号
| 手机号 | 类型 | 说明 |
|--------|------|------|
| 13800000008 | C端用户 | 张大爷（老年人） |
| 13800000009 | 跨端用户 | 李奶奶（既是老人也是工作人员） |
| 13800000010 | C端用户 | 王爷爷（老年人） |
| 13800000011 | C端用户 | 孙女士（孕妇） |
| 13800000012 | C端用户 | 赵先生（残障人士） |
| 13800000013 | C端用户 | 小明（儿童） |

所有用户密码（仅用于B端登录）：`Test@123`

---

## 🔗 相关链接

- [设计文档](./quick-start-design.md)
- [测试覆盖率报告](./test-coverage-report.md)
- [高德地图API文档](https://lbs.amap.com/api/webservice/guide/api/georegeo)

---

**文档版本**: v1.0
**更新时间**: 2026-01-20
**维护者**: SuperClaude Framework
