# sCare C端快速开通系统设计文档

**版本**: v1.0
**创建时间**: 2026-01-20
**作者**: SuperClaude Framework
**架构**: 扫码即用 PWA + 验证码登录 + 自动注册

---

## 📋 目录

1. [核心目标](#核心目标)
2. [用户决策清单](#用户决策清单)
3. [完整API设计](#完整api设计)
4. [数据库Schema](#数据库schema)
5. [用户流程图](#用户流程图)
6. [开发清单](#开发清单)
7. [工作量估算](#工作量估算)

---

## 🎯 核心目标

### 设计原则
- **扫码即用**：二维码粘贴在单元楼门口，居民扫码直接使用
- **零门槛**：无需预先注册，填写信息即注册
- **老年友好**：验证码登录，无需记密码
- **PWA应用**：离线可用，添加到主屏幕

### 核心用户场景
```
老人需要助浴服务
  ↓
扫社区服务站二维码
  ↓
填写：姓名 + 手机号 + 验证码 + 服务需求 + 地址
  ↓
自动注册 + 自动登录 + 自动创建服务请求
  ↓
工作人员接单 → 上门服务 → 结单 → 用户评价
```

---

## ✅ 用户决策清单

### 决策1：用户模型
**选择**: A. 手机号 = 用户 = 服务对象（最简单）
- 一个手机号对应一个账号
- 不支持一个手机号管理多个服务对象
- customer_type 等信息在个人中心补充

### 决策2：首次流程
**选择**: A. 极简一步（一个表单提交所有）
- 必填：姓名、手机号、验证码、服务类型、地址
- 可选：补充说明
- 跳过：customer_type、照片、紧急联系人（后续补充）

### 决策3：再次使用
**选择**: B. 二次确认身份（更安全）
- 检测到token，显示"欢迎回来，XXX"
- 提供[继续使用] [切换账号]选项
- 预填充姓名、手机号、上次地址

### 决策4：验证码失败
**选择**: B. 强制联系工作人员
- 提示："短信服务暂未开通，请联系工作人员：010-12345678"
- 开发环境使用Mock（控制台打印验证码）

### 决策5：定位失败
**选择**: B. 允许手动填地址
- GPS定位失败 → 允许手动输入地址
- 调用高德地理编码API解析地址 → 获取经纬度
- 根据经纬度匹配最近服务站

### 决策6：无可用服务站
**处理方式**:
1. 提示："您所在区域暂无服务站，已记录您的需求"
2. 将请求标记为"待扩展区域"（status: no_coverage）
3. 后台通知运营团队

### 决策7：用户评价
**选择**: C. 暂不做评价
- 先上线核心功能
- 评价功能后续迭代

### 决策8：二维码URL
**选择**: A. https://app.scare.com/quick?station=1
- 简短易读
- 适合打印在海报/宣传单

---

## 📐 完整API设计

### 1. 发送验证码
```http
POST /api/v1/c/auth/send-code
Content-Type: application/json

Request:
{
  "phone": "13800000001"
}

Response 200:
{
  "msg": "ok",
  "data": {
    "expires_in": 300
  }
}

Response 429 (频率限制):
{
  "msg": "发送过于频繁，请1分钟后再试",
  "data": null
}

Response 429 (每日限制):
{
  "msg": "今日发送次数已达上限（10次）",
  "data": null
}

业务规则:
- 同一手机号：60秒内只能发送1次
- 同一手机号：每天最多发送10次
- 验证码有效期：5分钟
- 验证码长度：6位数字
- Redis存储：
  - Key: sms_code:{phone}, Value: {code}, TTL: 300s
  - Key: sms_rate:{phone}:minute, Value: {timestamp}, TTL: 60s
  - Key: sms_rate:{phone}:daily, Value: {count}, TTL: 86400s
```

### 2. 快速开通（注册+登录+创建请求）
```http
POST /api/v1/c/auth/quick-start
Content-Type: application/json

Request:
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

Response 200 (新用户):
{
  "msg": "ok",
  "data": {
    "is_new_user": true,
    "token": {
      "access_token": "eyJhbGc...",
      "refresh_token": "eyJhbGc...",
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

Response 200 (老用户):
{
  "msg": "ok",
  "data": {
    "is_new_user": false,
    "token": {...},
    "user": {
      "id": 8,
      "phone": "13800000008",
      "name": "张大爷",
      "customer_type": "elderly"
    },
    "request": {...}
  }
}

Response 400 (验证码错误):
{
  "msg": "验证码错误或已过期",
  "data": null
}

Response 409 (手机号已注册):
{
  "msg": "该手机号已注册，请使用验证码登录",
  "data": {
    "existing_user": true
  }
}

业务逻辑:
1. 验证验证码（Redis: GET sms_code:{phone}）
2. 检查手机号是否已注册
   - 已注册 → 验证码登录 + 创建请求
   - 未注册 → 创建用户 + 创建档案 + 登录 + 创建请求
3. 生成JWT Token (type: c_end, roles: nil)
4. 创建服务请求（status: pending）
5. 删除已使用的验证码（Redis: DEL sms_code:{phone}）
6. 返回token + user + request

事务保证:
- 使用数据库事务
- 失败自动回滚（用户、档案、请求）
```

### 3. 地址解析（高德API代理）
```http
POST /api/v1/c/geocode
Content-Type: application/json

Request:
{
  "address": "北京市朝阳区幸福小区1号楼101"
}

Response 200:
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

Response 404 (无服务站覆盖):
{
  "msg": "您所在区域暂无服务站，已记录您的需求",
  "data": {
    "lat": 39.9042,
    "lng": 116.4074,
    "formatted_address": "北京市朝阳区幸福街道幸福小区1号楼101室",
    "no_coverage": true
  }
}

Response 400 (地址解析失败):
{
  "msg": "地址解析失败，请检查地址是否正确",
  "data": null
}

业务逻辑:
1. 调用高德地理编码API
   - URL: https://restapi.amap.com/v3/geocode/geo
   - Params: address, key
2. 解析返回的经纬度
3. 根据经纬度匹配最近服务站
   - 遍历所有活跃服务站
   - 计算Haversine距离
   - 返回最近的服务站
4. 如果所有服务站距离 > 10km → 无覆盖
```

### 4. 检查用户状态（供前端预填充）
```http
GET /api/v1/c/auth/check
Headers:
  Authorization: Bearer <token>

Response 200:
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

Response 401 (token无效):
{
  "msg": "unauthorized",
  "data": null
}

业务逻辑:
1. 从token解析user_id
2. 查询用户信息
3. 查询最近一次服务请求（用于预填充地址）
```

### 5. 更新个人资料
```http
PUT /api/v1/c/profile
Headers:
  Authorization: Bearer <token>
Content-Type: application/json

Request:
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

Response 200:
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

业务逻辑:
1. 从token获取user_id
2. 查询customer_profile
3. 更新字段
4. emergency_contact自动序列化为JSON
```

### 6. 已有功能（无需修改）
```http
# 创建服务请求（C端系统内详细创建）
POST /api/v1/c/requests
Headers: Authorization: Bearer <token>

# 查看我的请求列表
GET /api/v1/c/requests
Headers: Authorization: Bearer <token>

# 查看请求详情
GET /api/v1/c/requests/:id
Headers: Authorization: Bearer <token>

# 取消请求
POST /api/v1/c/requests/:id/cancel
Headers: Authorization: Bearer <token>

# 获取通知列表
GET /api/v1/c/notifications
Headers: Authorization: Bearer <token>

# 标记通知已读
POST /api/v1/c/notifications/:id/read
Headers: Authorization: Bearer <token>

# 上传文件
POST /api/v1/c/upload
Headers: Authorization: Bearer <token>
```

---

## 🗄️ 数据库Schema调整

### 1. 验证码存储（Redis）
```
# 验证码
Key: sms_code:{phone}
Value: {6位验证码}
TTL: 300秒

# 分钟级频率限制
Key: sms_rate:{phone}:minute
Value: {timestamp}
TTL: 60秒

# 每日频率限制
Key: sms_rate:{phone}:daily
Value: {发送次数}
TTL: 86400秒

示例:
SET sms_code:13800000001 "123456" EX 300
SETEX sms_rate:13800000001:minute 1 60
INCR sms_rate:13800000001:daily
EXPIRE sms_rate:13800000001:daily 86400
```

### 2. service_requests 表调整
```sql
ALTER TABLE `service_requests`
ADD COLUMN `station_id` bigint unsigned DEFAULT NULL COMMENT '扫码的服务站ID（来源追踪）',
ADD COLUMN `note` text COMMENT '用户补充说明',
ADD COLUMN `source` varchar(20) DEFAULT 'qrcode' COMMENT '请求来源：qrcode/app/admin',
ADD KEY `idx_station_status` (`station_id`, `status`),
ADD KEY `idx_source` (`source`);

COMMENT ON COLUMN `station_id`: '扫码携带的服务站ID，用于统计二维码效果';
COMMENT ON COLUMN `source`: 'qrcode-扫码快速通道, app-C端系统创建, admin-B端后台创建';
```

### 3. users 表调整
```sql
ALTER TABLE `users`
ADD COLUMN `source` varchar(20) DEFAULT 'quick_start' COMMENT '注册来源：quick_start/admin/import';

COMMENT ON COLUMN `source`:
  'quick_start-扫码快速注册, admin-后台创建, import-批量导入';
```

### 4. 数据示例
```sql
-- 用户通过扫码注册
INSERT INTO `users` (`phone`, `name`, `status`, `source`)
VALUES ('13800000014', '王奶奶', 'active', 'quick_start');

-- 创建客户档案（字段暂时为空，后续补充）
INSERT INTO `customer_profiles` (`user_id`)
VALUES (14);

-- 创建服务请求
INSERT INTO `service_requests` (
  `user_id`, `request_no`, `service_type`, `station_id`,
  `lat`, `lng`, `address`, `note`, `status`, `source`
) VALUES (
  14, 'R202601200001', 'bath_assistance', 1,
  39.9042, 116.4074, '北京市朝阳区幸福小区1号楼101',
  '腿脚不便，需要轮椅', 'pending', 'qrcode'
);
```

---

## 🔄 完整用户流程图

### 流程1：首次扫码注册+提交请求

```
┌─────────────────────────────────────────┐
│ 1. 用户扫码（单元楼门口海报）            │
│ URL: https://app.scare.com/quick?station=1│
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 2. 前端路由解析                          │
│ route: /quick                            │
│ query: { station: 1 }                    │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 3. 检查localStorage                      │
│ token = localStorage.getItem('token')    │
├─────────────────────────────────────────┤
│ if (token存在) {                         │
│   GET /api/v1/c/auth/check              │
│   if (有效) → 跳到流程2（预填充）        │
│   if (无效) → 清除token，继续            │
│ } else {                                │
│   显示空白表单                           │
│ }                                       │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 4. 快速通道表单（一屏内）                │
├─────────────────────────────────────────┤
│ 姓名：     [____________]               │
│ 手机号：   [____________]               │
│ 验证码：   [____] [获取验证码] 60s      │
│                                         │
│ 服务类型：  [助浴服务 ▼]                │
│           助浴/助餐/护理/康复/其他       │
│                                         │
│ 服务地址：  [📍 使用当前位置]           │
│           or [✍️ 手动输入]              │
│                                         │
│ 详细地址：  [____________]              │
│ 补充说明：  [____________] (可选)       │
│                                         │
│         [提交服务请求]                   │
│                                         │
│ ──────────────────                      │
│ 已有账号？[进入个人中心] →              │
└──────────────────┬──────────────────────┘
                   ↓
       用户点击[获取验证码]
                   ↓
┌─────────────────────────────────────────┐
│ 5. 发送验证码                            │
│ POST /api/v1/c/auth/send-code           │
│ { "phone": "13800000001" }              │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 6. 后端验证码逻辑                        │
├─────────────────────────────────────────┤
│ • 检查分钟级限制                         │
│   EXISTS sms_rate:{phone}:minute        │
│   → 429 "请1分钟后再试"                  │
│                                         │
│ • 检查每日限制                           │
│   GET sms_rate:{phone}:daily ≥ 10       │
│   → 429 "今日次数已达上限"               │
│                                         │
│ • 生成6位验证码                          │
│   code = Random(100000, 999999)         │
│                                         │
│ • 存储到Redis                            │
│   SET sms_code:{phone} {code} EX 300    │
│   SETEX sms_rate:{phone}:minute 1 60    │
│   INCR sms_rate:{phone}:daily           │
│                                         │
│ • 发送短信（Mock）                       │
│   if (env == "dev")                     │
│     log.Printf("验证码：%s", code)       │
│   else                                  │
│     return "短信服务未开通"              │
└──────────────────┬──────────────────────┘
                   ↓
      用户填写完表单，点击[提交]
                   ↓
┌─────────────────────────────────────────┐
│ 7. 提交快速开通                          │
│ POST /api/v1/c/auth/quick-start         │
│ {                                       │
│   "phone": "13800000001",               │
│   "code": "123456",                     │
│   "name": "张大爷",                      │
│   "service_request": {                  │
│     "service_type": "bath_assistance",  │
│     "station_id": 1,                    │
│     "address": "...",                   │
│     "lat": 39.9042,                     │
│     "lng": 116.4074,                    │
│     "note": "腿脚不便"                   │
│   }                                     │
│ }                                       │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 8. 后端业务逻辑（事务）                  │
├─────────────────────────────────────────┤
│ BEGIN TRANSACTION                       │
│                                         │
│ • 验证验证码                             │
│   code = GET sms_code:{phone}           │
│   if (code != input) → 400错误          │
│                                         │
│ • 检查手机号                             │
│   user = SELECT * FROM users            │
│         WHERE phone = ?                 │
│                                         │
│ • IF (user不存在) {                     │
│     // 新用户注册                        │
│     INSERT INTO users                   │
│     (phone, name, status, source)       │
│     VALUES (?, ?, 'active', 'quick_start')│
│                                         │
│     INSERT INTO customer_profiles       │
│     (user_id) VALUES (?)                │
│   }                                     │
│                                         │
│ • 生成JWT Token                         │
│   token = GenerateToken(                │
│     userID, "c_end", 0, nil             │
│   )                                     │
│                                         │
│ • 创建服务请求                           │
│   INSERT INTO service_requests          │
│   (user_id, request_no, service_type,   │
│    station_id, lat, lng, address,       │
│    note, status, source)                │
│   VALUES (?, ?, ?, ?, ?, ?, ?, ?,       │
│           'pending', 'qrcode')          │
│                                         │
│ • 删除验证码                             │
│   DEL sms_code:{phone}                  │
│                                         │
│ COMMIT                                  │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 9. 前端处理响应                          │
├─────────────────────────────────────────┤
│ • 保存token                              │
│   localStorage.setItem('token', access) │
│   localStorage.setItem('refresh', refresh)│
│                                         │
│ • 保存用户信息                           │
│   localStorage.setItem('user', JSON)    │
│                                         │
│ • 跳转到请求详情                         │
│   router.push(`/requests/${id}`)        │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 10. 请求详情页                           │
├─────────────────────────────────────────┤
│ ✅ 提交成功！                            │
│ 工作人员将在30分钟内联系您               │
│                                         │
│ 服务单号：R202601200001                  │
│ 服务类型：助浴服务                       │
│ 服务地址：朝阳区幸福小区1号楼101         │
│ 当前状态：🟡 待分派                      │
│ 提交时间：2026-01-20 10:30              │
│                                         │
│ [查看更多请求]  [完善个人资料]          │
└─────────────────────────────────────────┘
```

### 流程2：老用户再次扫码（预填充）

```
┌─────────────────────────────────────────┐
│ 1. 用户扫码                              │
│ URL: https://app.scare.com/quick?station=1│
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 2. 检测到localStorage有token             │
│ GET /api/v1/c/auth/check                │
└──────────────────┬──────────────────────┘
                   ↓
┌─────────────────────────────────────────┐
│ 3. Token有效，预填充表单                 │
├─────────────────────────────────────────┤
│ 欢迎回来，张大爷！                       │
│                                         │
│ 姓名：     张大爷 (不可改)               │
│ 手机号：   138****0001 (不可改)         │
│                                         │
│ 服务类型：  [助浴服务 ▼]                │
│ 服务地址：  [朝阳区幸福小区1号楼101]     │
│            ← 上次地址自动填充            │
│            [📍 重新定位]                 │
│                                         │
│ 补充说明：  [____________]              │
│                                         │
│         [提交服务请求]                   │
│                                         │
│ [切换账号]  [进入个人中心]              │
└──────────────────┬──────────────────────┘
                   ↓
       用户点击[提交]（无需验证码）
                   ↓
┌─────────────────────────────────────────┐
│ 4. 直接创建请求                          │
│ POST /api/v1/c/requests                 │
│ Headers: Authorization: Bearer <token>  │
│ {                                       │
│   "service_type": "bath_assistance",    │
│   "station_id": 1,                      │
│   "address": "...",                     │
│   "lat": 39.9042,                       │
│   "lng": 116.4074,                      │
│   "note": "..."                         │
│ }                                       │
└──────────────────┬──────────────────────┘
                   ↓
         跳转到请求详情页
```

---

## 📝 后端开发清单

### Phase 1: 验证码系统 ⏱️ 2小时

#### 文件创建
- [ ] `internal/service/sms_service.go`
- [ ] `internal/handler/sms_handler.go`

#### 功能实现
```go
// internal/service/sms_service.go
type SMSService struct {
    rdb *redis.Client
}

func (s *SMSService) SendCode(phone string) error {
    // 1. 检查分钟级限制
    // 2. 检查每日限制
    // 3. 生成6位验证码
    // 4. 存储到Redis
    // 5. 发送短信（Mock）
}

func (s *SMSService) VerifyCode(phone, code string) error {
    // 1. 从Redis获取验证码
    // 2. 对比验证
    // 3. 验证成功后删除
}

func (s *SMSService) CheckRateLimit(phone string) error {
    // 检查频率限制
}
```

#### 路由配置
```go
// internal/handler/routes.go
cAuthGroup := api.Group("/c/auth")
cAuthGroup.POST("/send-code", smsHandler.SendCode)
```

#### 单元测试
- [ ] `internal/service/sms_service_test.go`
  - [ ] TestSendCode_Success
  - [ ] TestSendCode_MinuteLimit
  - [ ] TestSendCode_DailyLimit
  - [ ] TestVerifyCode_Success
  - [ ] TestVerifyCode_Expired
  - [ ] TestVerifyCode_Invalid

---

### Phase 2: 快速开通API ⏱️ 3小时

#### 文件修改/创建
- [ ] `internal/service/auth_service.go` (新增QuickStart方法)
- [ ] `internal/handler/c_auth_handler.go` (新增QuickStart处理)

#### 功能实现
```go
// internal/service/auth_service.go
type QuickStartInput struct {
    Phone          string
    Code           string
    Name           string
    ServiceRequest RequestInput
}

type QuickStartOutput struct {
    IsNewUser bool
    Token     *Tokens
    User      *domain.User
    Request   *domain.ServiceRequest
}

func (s *AuthService) QuickStart(input QuickStartInput) (*QuickStartOutput, error) {
    // 开启事务
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 1. 验证验证码
    if err := s.smsService.VerifyCode(input.Phone, input.Code); err != nil {
        tx.Rollback()
        return nil, err
    }

    // 2. 检查手机号是否已注册
    user, err := s.userRepo.GetByPhone(input.Phone)
    isNewUser := false

    if errors.Is(err, gorm.ErrRecordNotFound) {
        // 新用户：创建用户 + 档案
        user = &domain.User{
            Phone:  input.Phone,
            Name:   input.Name,
            Status: "active",
            Source: "quick_start",
        }
        if err := tx.Create(user).Error; err != nil {
            tx.Rollback()
            return nil, err
        }

        profile := &domain.CustomerProfile{
            UserID: user.ID,
        }
        if err := tx.Create(profile).Error; err != nil {
            tx.Rollback()
            return nil, err
        }

        isNewUser = true
    } else if err != nil {
        tx.Rollback()
        return nil, err
    }

    // 3. 生成Token
    token, err := s.jwtManager.GenerateToken(user.ID, "c_end", 0, nil)
    refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, "c_end")

    // 4. 创建服务请求
    request := &domain.ServiceRequest{
        UserID:      user.ID,
        RequestNo:   generateRequestNo(),
        ServiceType: input.ServiceRequest.ServiceType,
        StationID:   &input.ServiceRequest.StationID,
        Lat:         &input.ServiceRequest.Lat,
        Lng:         &input.ServiceRequest.Lng,
        Address:     &input.ServiceRequest.Address,
        Note:        &input.ServiceRequest.Note,
        Status:      "pending",
        Source:      "qrcode",
    }
    if err := tx.Create(request).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return &QuickStartOutput{
        IsNewUser: isNewUser,
        Token:     &Tokens{AccessToken: token, RefreshToken: refreshToken},
        User:      user,
        Request:   request,
    }, nil
}
```

#### 单元测试
- [ ] `internal/service/auth_service_test.go`
  - [ ] TestQuickStart_NewUser
  - [ ] TestQuickStart_ExistingUser
  - [ ] TestQuickStart_InvalidCode
  - [ ] TestQuickStart_TransactionRollback

---

### Phase 3: 地址解析 ⏱️ 1.5小时

#### 文件创建
- [ ] `internal/service/geocode_service.go`
- [ ] `internal/handler/geocode_handler.go`

#### 高德API集成
```go
// internal/service/geocode_service.go
type GeocodeService struct {
    apiKey      string
    stationRepo *repository.StationRepository
}

type GeocodeResult struct {
    Lat              float64
    Lng              float64
    FormattedAddress string
    MatchedStation   *MatchedStation
    NoCoverage       bool
}

type MatchedStation struct {
    ID       int64
    Name     string
    Distance float64 // 米
}

func (s *GeocodeService) Geocode(address string) (*GeocodeResult, error) {
    // 1. 调用高德地理编码API
    url := fmt.Sprintf(
        "https://restapi.amap.com/v3/geocode/geo?address=%s&key=%s",
        url.QueryEscape(address),
        s.apiKey,
    )

    resp, err := http.Get(url)
    // 解析响应...

    // 2. 匹配最近服务站
    station, distance := s.matchNearestStation(lat, lng)

    // 3. 判断是否有覆盖（距离 < 10km）
    noCoverage := distance > 10000

    return &GeocodeResult{...}, nil
}

func (s *GeocodeService) matchNearestStation(lat, lng float64) (*domain.Station, float64) {
    stations, _ := s.stationRepo.ListActive()

    var nearest *domain.Station
    minDistance := math.MaxFloat64

    for _, station := range stations {
        // 使用Haversine公式计算距离
        distance := haversineDistance(lat, lng, station.Lat, station.Lng)
        if distance < minDistance {
            minDistance = distance
            nearest = station
        }
    }

    return nearest, minDistance
}
```

#### 配置管理
- [ ] `internal/config/config.go` 添加高德API Key
```go
type Config struct {
    // ...
    Amap struct {
        APIKey string `yaml:"api_key"`
    } `yaml:"amap"`
}
```

#### 单元测试
- [ ] `internal/service/geocode_service_test.go`
  - [ ] TestGeocode_Success
  - [ ] TestGeocode_InvalidAddress
  - [ ] TestMatchStation_HasCoverage
  - [ ] TestMatchStation_NoCoverage

---

### Phase 4: 个人资料API ⏱️ 1小时

#### 文件修改
- [ ] `internal/handler/c_auth_handler.go` (新增Check方法)
- [ ] `internal/handler/profile_handler.go` (新建)

#### 功能实现
```go
// internal/handler/c_auth_handler.go
func (h *CAuthHandler) Check(c *gin.Context) {
    userID, _ := GetUserID(c)

    // 获取用户信息
    user, err := h.userRepo.GetByID(userID)

    // 获取最近一次请求（用于预填充）
    var lastRequest *domain.ServiceRequest
    h.db.Where("user_id = ?", userID).
        Order("created_at DESC").
        First(&lastRequest)

    Respond(c, http.StatusOK, "ok", gin.H{
        "user": user,
        "last_request": lastRequest,
    })
}

// internal/handler/profile_handler.go
type ProfileHandler struct {
    customerRepo *repository.CustomerRepository
}

func (h *ProfileHandler) Update(c *gin.Context) {
    userID, _ := GetUserID(c)

    var req struct {
        CustomerType     *string              `json:"customer_type"`
        Gender           *string              `json:"gender"`
        BirthDate        *string              `json:"birth_date"`
        HealthStatus     *string              `json:"health_status"`
        DisabilityLevel  *string              `json:"disability_level"`
        EmergencyContact *domain.EmergencyContact `json:"emergency_contact"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        RespondError(c, http.StatusBadRequest, "invalid payload")
        return
    }

    // 获取档案
    profile, err := h.customerRepo.GetByUserID(userID)

    // 更新字段
    if req.CustomerType != nil {
        profile.CustomerType = req.CustomerType
    }
    // ... 更新其他字段

    if err := h.customerRepo.Update(profile); err != nil {
        RespondError(c, http.StatusInternalServerError, "update failed")
        return
    }

    Respond(c, http.StatusOK, "ok", profile)
}
```

---

### Phase 5: 路由和权限配置 ⏱️ 30分钟

#### 修改文件
- [ ] `internal/handler/routes.go`

#### 路由注册
```go
// 白名单路由（无需认证）
cAuthGroup := api.Group("/c/auth")
cAuthGroup.POST("/send-code", smsHandler.SendCode)          // 发送验证码
cAuthGroup.POST("/quick-start", cAuthHandler.QuickStart)    // 快速开通

// 地址解析（无需认证）
api.POST("/c/geocode", geocodeHandler.Geocode)

// 需要认证的路由
secured := api.Group("")
secured.Use(middleware.AuthMiddleware(jwtManager, blacklistService))
secured.GET("/c/auth/check", cAuthHandler.Check)            // 检查状态
secured.PUT("/c/profile", profileHandler.Update)            // 更新资料
```

#### Casbin权限策略
- [ ] 更新 `database/seeds/seed.sql`
```sql
-- C端公共权限（已有，确认包含）
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES
('p', 'role:authenticated', '/api/v1/c/auth/check', 'GET'),
('p', 'role:authenticated', '/api/v1/c/profile', 'PUT');
```

---

### Phase 6: 数据库迁移 ⏱️ 30分钟

#### 创建迁移文件
- [ ] `database/migrations/006_quick_start.sql`

```sql
-- =====================================================
-- sCare C端快速开通系统 - 数据库迁移
-- 版本：v1.0
-- 日期：2026-01-20
-- =====================================================

-- 1. 修改 service_requests 表
ALTER TABLE `service_requests`
ADD COLUMN `station_id` bigint unsigned DEFAULT NULL COMMENT '扫码的服务站ID（来源追踪）' AFTER `user_id`,
ADD COLUMN `note` text COMMENT '用户补充说明' AFTER `images`,
ADD COLUMN `source` varchar(20) DEFAULT 'app' COMMENT '请求来源：qrcode/app/admin' AFTER `note`;

ALTER TABLE `service_requests`
ADD KEY `idx_station_status` (`station_id`, `status`),
ADD KEY `idx_source` (`source`);

-- 2. 修改 users 表
ALTER TABLE `users`
ADD COLUMN `source` varchar(20) DEFAULT 'admin' COMMENT '注册来源：quick_start/admin/import' AFTER `status`;

-- 3. 更新已有数据
UPDATE `service_requests` SET `source` = 'app' WHERE `source` IS NULL;
UPDATE `users` SET `source` = 'admin' WHERE `source` IS NULL;

-- 4. 验证
SELECT
    COUNT(*) as total_requests,
    COUNT(DISTINCT source) as source_types,
    source
FROM `service_requests`
GROUP BY source;

SELECT
    COUNT(*) as total_users,
    COUNT(DISTINCT source) as source_types,
    source
FROM `users`
GROUP BY source;
```

#### 执行迁移
```bash
mysql -u root -p scare < database/migrations/006_quick_start.sql
```

---

### Phase 7: 集成测试 ⏱️ 1.5小时

#### 测试文件
- [ ] `internal/handler/quick_integration_test.go`

#### 测试场景
```go
func TestQuickStartIntegration(t *testing.T) {
    // 1. 发送验证码
    resp1 := POST("/api/v1/c/auth/send-code", {
        "phone": "13900000100"
    })
    assert.Equal(t, 200, resp1.StatusCode)

    // 2. 快速开通（新用户）
    resp2 := POST("/api/v1/c/auth/quick-start", {
        "phone": "13900000100",
        "code": "123456", // Mock环境固定验证码
        "name": "测试用户",
        "service_request": {
            "service_type": "bath_assistance",
            "station_id": 1,
            "address": "测试地址",
            "lat": 39.9042,
            "lng": 116.4074
        }
    })
    assert.Equal(t, 200, resp2.StatusCode)
    assert.True(t, resp2.Data.IsNewUser)

    token := resp2.Data.Token.AccessToken

    // 3. 检查用户状态
    resp3 := GET("/api/v1/c/auth/check", token)
    assert.Equal(t, 200, resp3.StatusCode)

    // 4. 更新个人资料
    resp4 := PUT("/api/v1/c/profile", token, {
        "customer_type": "elderly",
        "emergency_contact": {
            "name": "家属",
            "phone": "13900000101",
            "relation": "子女"
        }
    })
    assert.Equal(t, 200, resp4.StatusCode)
}

func TestRateLimit(t *testing.T) {
    // 分钟级限制
    for i := 0; i < 2; i++ {
        resp := POST("/api/v1/c/auth/send-code", {
            "phone": "13900000200"
        })
        if i == 0 {
            assert.Equal(t, 200, resp.StatusCode)
        } else {
            assert.Equal(t, 429, resp.StatusCode)
        }
    }
}
```

---

## ⏰ 工作量估算

### 后端开发：10小时
| 任务 | 预估时间 | 优先级 |
|------|---------|--------|
| 验证码系统 | 2h | P0 |
| 快速开通API | 3h | P0 |
| 地址解析 | 1.5h | P1 |
| 个人资料API | 1h | P1 |
| 路由配置 | 0.5h | P0 |
| 数据库迁移 | 0.5h | P0 |
| 集成测试 | 1.5h | P1 |

### 前端开发：17小时
| 任务 | 预估时间 | 优先级 |
|------|---------|--------|
| 项目初始化 | 1h | P0 |
| 快速通道页 | 4h | P0 |
| 请求详情页 | 2h | P0 |
| C端系统页面 | 6h | P1 |
| PWA优化 | 2h | P1 |
| 地图集成 | 2h | P1 |

### **总计：27小时（3-4个工作日）**

---

## 📅 开发排期建议

### Day 1（后端为主）⏱️ 7小时
- ✅ 验证码系统（2h）
- ✅ 快速开通API（3h）
- ✅ 地址解析（1.5h）
- ✅ 数据库迁移（0.5h）

### Day 2（前后端并行）⏱️ 8小时
- **后端**：个人资料API + 路由配置 + 测试（3h）
- **前端**：项目初始化 + 快速通道页（5h）

### Day 3（前端为主）⏱️ 7小时
- ✅ 请求详情页（2h）
- ✅ 登录页+首页（3h）
- ✅ 地图集成（2h）

### Day 4（完善+测试）⏱️ 5小时
- ✅ 请求列表+个人中心（3h）
- ✅ PWA优化（2h）
- ✅ 联调测试（前后端）

---

## 🔐 安全措施

### 1. 验证码安全
```yaml
频率限制:
  分钟级: 1次/分钟/手机号
  每日级: 10次/天/手机号

有效期: 5分钟

存储: Redis（自动过期）

验证:
  - 验证后立即删除
  - 最多验证3次（暂未实现）
```

### 2. Token安全
```yaml
JWT配置:
  access_token: 1小时有效期
  refresh_token: 7天有效期

黑名单:
  - 登出时加入黑名单
  - Redis存储，TTL = token剩余时间
```

### 3. 数据脱敏
```yaml
展示规则:
  手机号: 138****0001
  地址: XX小区X号楼（不显示门牌号）

完整信息仅对:
  - 用户本人
  - 接单工作人员
  - 站长/管理员
```

---

## 📊 监控指标

### 业务指标
```yaml
转化率:
  - 扫码UV → 提交请求
  - 新用户注册转化率
  - 服务完成率

活跃度:
  - DAU/MAU
  - 人均请求次数
  - 二维码扫码热力图

质量:
  - 请求响应时长
  - 服务评价（后续）
```

### 技术指标
```yaml
性能:
  - API响应时间（P95 < 200ms）
  - 验证码发送成功率
  - 地址解析准确率

稳定性:
  - 服务可用性 > 99.9%
  - 错误率 < 0.1%

资源:
  - Redis内存使用
  - 高德API调用量
```

---

## 🚀 未来优化方向

### Phase 2（迭代功能）
1. **用户评价系统**
   - 服务完成后评分
   - 评价工作人员
   - 评价内容审核

2. **智能推荐**
   - 根据历史服务推荐
   - 根据时间推荐（定期服务提醒）

3. **家庭账户**
   - 一个手机号管理多个服务对象
   - 子女帮父母申请服务

4. **消息推送**
   - 服务状态变更推送
   - 工作人员到达提醒
   - 定期关怀提醒

### Phase 3（高级功能）
1. **AI客服**
   - 智能问答
   - 服务咨询

2. **语音输入**
   - 老年人语音填单
   - 语音导航

3. **视频通话**
   - 远程评估
   - 服务监督

---

**文档结束** | 版本: v1.0 | 更新: 2026-01-20
