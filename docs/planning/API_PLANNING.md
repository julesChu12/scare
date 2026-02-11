# sCare MVP 后端API开发规划

**版本**: v1.0.0-MVP
**锁定日期**: 2026-01-18
**目标**: 完整需求流程打通，支持演示和测试

---

## 🎯 MVP核心路径

```
C端用户提交需求
  ↓
地理围栏匹配（不在围栏内拒绝，多围栏随机选择）
  ↓
任务派单到站点任务池
  ↓
工作人员认领任务
  ↓
处理并完成任务
  ↓
用户查看结果和评价
```

---

## 📊 API优先级分级

### P0 - 核心流程 (必须实现，阻塞演示)

| API | 路径 | 方法 | 角色 | 说明 |
|-----|------|------|------|------|
| **认证登录** | `/api/auth/login` | POST | All | 手机号+密码登录，返回JWT |
| **获取当前用户** | `/api/auth/profile` | GET | All | 返回用户信息+角色 |
| **提交需求** | `/api/requests` | POST | elderly/family | 含地理匹配，返回匹配结果 |
| **需求详情** | `/api/requests/:id` | GET | elderly/family/staff/admin | 查询单个需求 |
| **站点任务池** | `/api/tasks` | GET | staff/station_manager | 查询待认领任务列表 |
| **认领任务** | `/api/tasks/:id/claim` | POST | staff | 工作人员认领 |
| **完成任务** | `/api/tasks/:id/complete` | POST | staff | 提交完成信息 |
| **我的任务** | `/api/tasks/my` | GET | staff | 我认领的任务列表 |
| **创建站点** | `/api/stations` | POST | admin | 管理员创建站点 |
| **创建围栏** | `/api/zones` | POST | admin | 管理员创建地理围栏 |

**P0总计**: 10个核心接口

---

### P1 - 完善体验 (提升可用性)

| API | 路径 | 方法 | 角色 | 说明 |
|-----|------|------|------|------|
| **需求列表** | `/api/requests` | GET | elderly/family | 我的需求列表 |
| **任务转派** | `/api/tasks/:id/transfer` | POST | station_manager | 站长转派任务 |
| **通知列表** | `/api/notifications` | GET | All | 站内信列表 |
| **标记已读** | `/api/notifications/:id/read` | PUT | All | 标记通知已读 |
| **图片上传** | `/api/upload` | POST | All | 本地存储，返回URL |
| **站点列表** | `/api/stations` | GET | admin | 站点管理列表 |
| **围栏列表** | `/api/zones` | GET | admin | 围栏管理列表 |
| **更新任务状态** | `/api/tasks/:id/status` | PUT | staff | 更新为processing |

**P1总计**: 8个接口

---

### P2 - 增强功能 (可延后)

| API | 路径 | 方法 | 角色 | 说明 |
|-----|------|------|------|------|
| **用户列表** | `/api/users` | GET | admin | 用户管理 |
| **创建用户** | `/api/users` | POST | admin | 批量创建工作人员 |
| **更新用户** | `/api/users/:id` | PUT | admin | 修改用户信息 |
| **禁用用户** | `/api/users/:id/disable` | PUT | admin | 禁用账号 |
| **任务历史** | `/api/tasks/:id/history` | GET | staff/admin | 任务操作历史 |
| **统计报表** | `/api/stats/overview` | GET | station_manager/admin | 概览统计 |
| **更新围栏** | `/api/zones/:id` | PUT | admin | 修改围栏 |
| **删除围栏** | `/api/zones/:id` | DELETE | admin | 删除围栏 |

**P2总计**: 8个接口

---

## 🔐 权限矩阵

| 接口类别 | elderly | family | staff | station_manager | admin |
|---------|---------|--------|-------|-----------------|-------|
| **认证相关** | ✅ | ✅ | ✅ | ✅ | ✅ |
| **提交需求** | ✅ | ✅ | ❌ | ❌ | ❌ |
| **查看需求** | ✅(自己) | ✅(自己) | ✅(站点) | ✅(站点) | ✅(全部) |
| **任务认领** | ❌ | ❌ | ✅ | ✅ | ✅ |
| **任务转派** | ❌ | ❌ | ❌ | ✅ | ✅ |
| **站点管理** | ❌ | ❌ | ❌ | ❌ | ✅ |
| **围栏管理** | ❌ | ❌ | ❌ | ❌ | ✅ |
| **用户管理** | ❌ | ❌ | ❌ | ❌ | ✅ |
| **通知查看** | ✅ | ✅ | ✅ | ✅ | ✅ |

**权限实现方式**: Casbin RBAC + 路由中间件

---

## 📝 核心API详细设计

### 1. POST /api/auth/login (认证登录)

**请求**:
```json
{
  "phone": "13800000001",
  "password": "Test@123"
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2026-01-19T10:00:00Z",
    "user": {
      "id": 1,
      "phone": "13800000001",
      "name": "系统管理员",
      "role": "admin",
      "station_id": null
    }
  }
}
```

**失败响应** (401):
```json
{
  "code": 401,
  "message": "用户名或密码错误",
  "error": "invalid_credentials"
}
```

---

### 2. POST /api/requests (提交需求)

**请求**:
```json
{
  "service_type": "meal",
  "description": "需要午餐送餐服务，清淡少油",
  "submit_location_lat": 40.0500000,
  "submit_location_lng": 116.3800000,
  "contact_name": "张大爷",
  "contact_phone": "13800000008",
  "address": "北京市昌平区霍营街道华龙苑北里小区1号楼3单元501",
  "appointment_time": "2026-01-18T11:30:00Z",
  "urgency": "normal"
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "需求提交成功",
  "data": {
    "request_id": 1,
    "request_no": "REQ-2026011801",
    "status": "dispatched",
    "matched_station": {
      "id": 1,
      "name": "霍营街道养老服务中心",
      "code": "STATION-HY-001"
    },
    "matched_zone": {
      "id": 1,
      "name": "霍营站A区-华龙苑北里",
      "match_method": "geofence"  // geofence / random(多围栏随机)
    }
  }
}
```

**拒绝响应** (400):
```json
{
  "code": 400,
  "message": "您的位置不在任何服务范围内",
  "error": "no_geofence_matched",
  "data": {
    "submit_location_lat": 40.1000000,
    "submit_location_lng": 116.5000000,
    "nearest_station": {
      "id": 2,
      "name": "龙泽园养老服务站",
      "distance_km": 5.2
    },
    "suggestion": "请联系最近的服务站点"
  }
}
```

---

### 3. GET /api/tasks (站点任务池)

**查询参数**:
- `status`: dispatched (默认只显示待认领)
- `page`: 1
- `page_size`: 10

**成功响应** (200):
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "task_id": 1,
        "request_id": 1,
        "request_no": "REQ-2026011801",
        "service_type": "meal",
        "service_type_name": "送餐服务",
        "urgency": "normal",
        "contact_name": "张大爷",
        "contact_phone": "13800000008",
        "address": "华龙苑北里1号楼3单元501",
        "appointment_time": "2026-01-18T11:30:00Z",
        "status": "dispatched",
        "created_at": "2026-01-18T10:00:00Z"
      }
    ]
  }
}
```

---

### 4. POST /api/tasks/:id/claim (认领任务)

**请求**: 无body

**成功响应** (200):
```json
{
  "code": 200,
  "message": "任务认领成功",
  "data": {
    "task_id": 1,
    "status": "claimed",
    "claimed_at": "2026-01-18T10:30:00Z",
    "staff": {
      "id": 4,
      "name": "王小红"
    }
  }
}
```

**失败响应** (409):
```json
{
  "code": 409,
  "message": "任务已被其他人认领",
  "error": "task_already_claimed",
  "data": {
    "claimed_by": {
      "id": 5,
      "name": "刘小明"
    },
    "claimed_at": "2026-01-18T10:28:00Z"
  }
}
```

---

### 5. POST /api/tasks/:id/complete (完成任务)

**请求**:
```json
{
  "staff_notes": "已完成送餐，老人身体状况良好",
  "images": [
    "/uploads/2026/01/18/task_1_photo1.jpg",
    "/uploads/2026/01/18/task_1_photo2.jpg"
  ]
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "任务已完成",
  "data": {
    "task_id": 1,
    "status": "completed",
    "completed_at": "2026-01-18T12:00:00Z",
    "request_status": "completed"
  }
}
```

---

### 6. POST /api/tasks/:id/transfer (任务转派)

**请求**:
```json
{
  "to_staff_id": 5,
  "remark": "王小红临时有事，转派给刘小明处理"
}
```

**成功响应** (200):
```json
{
  "code": 200,
  "message": "任务转派成功",
  "data": {
    "old_task_id": 1,
    "old_task_status": "transferred",
    "new_task_id": 2,
    "new_task_status": "claimed",
    "from_staff": {
      "id": 4,
      "name": "王小红"
    },
    "to_staff": {
      "id": 5,
      "name": "刘小明"
    },
    "history_recorded": true
  }
}
```

---

### 7. POST /api/upload (图片上传)

**请求**: `multipart/form-data`
- `file`: 图片文件(最大5MB)

**成功响应** (200):
```json
{
  "code": 200,
  "message": "上传成功",
  "data": {
    "url": "/uploads/2026/01/18/abc123.jpg",
    "filename": "abc123.jpg",
    "size": 102400,
    "mime_type": "image/jpeg"
  }
}
```

**本地存储路径**: `./uploads/YYYY/MM/DD/filename`
**MVP后切换**: 配置OSS后返回OSS URL

---

## 🔄 地理围栏匹配逻辑

### 匹配流程

```go
// 伪代码
func MatchGeofence(lat, lng float64) (*MatchResult, error) {
    // 1. 查询所有active围栏
    zones := queryActiveZones()

    // 2. 遍历检查点是否在多边形内(射线法)
    matchedZones := []Zone{}
    for zone := range zones {
        if isPointInPolygon(lat, lng, zone.Points) {
            matchedZones = append(matchedZones, zone)
        }
    }

    // 3. 处理匹配结果
    if len(matchedZones) == 0 {
        // 不在任何围栏内，拒绝接单
        return nil, ErrNoGeofenceMatched
    }

    if len(matchedZones) == 1 {
        // 唯一匹配
        return &MatchResult{
            Zone: matchedZones[0],
            Method: "geofence"
        }, nil
    }

    // 4. 多围栏随机选择(先消费先得到)
    randomIndex := rand.Intn(len(matchedZones))
    return &MatchResult{
        Zone: matchedZones[randomIndex],
        Method: "random"
    }, nil
}
```

### 容错处理

| 场景 | 处理方式 | HTTP状态码 |
|------|---------|-----------|
| 不在任何围栏内 | 拒绝接单，提示最近站点 | 400 |
| 命中多个围栏 | 随机选择一个 | 200 |
| 围栏JSON解析失败 | 跳过该围栏，继续匹配 | - |
| 所有围栏都异常 | 拒绝接单 | 500 |

---

## 🗃️ 数据库操作要点

### 需求提交流程

```go
// 事务操作
tx := db.Begin()

// 1. 创建需求记录(status=pending)
request := CreateRequest(tx, requestData)

// 2. 地理围栏匹配
zone, err := MatchGeofence(lat, lng)
if err != nil {
    tx.Rollback()
    return err // 返回400，拒绝接单
}

// 3. 更新需求状态和站点
request.StationID = &zone.StationID
request.Status = "dispatched"
tx.Save(&request)

// 4. 创建任务记录
task := CreateTask(tx, request.ID, zone.StationID)

// 5. 写任务历史
WriteHistory(tx, task.ID, "dispatched", systemUserID)

// 6. 写通知(站内信)
CreateNotification(tx, request.UserID, "request_created", request.ID)

tx.Commit()
```

### 任务认领并发控制

```go
// 使用乐观锁避免并发认领
func ClaimTask(taskID, staffID int64) error {
    result := db.Model(&TaskAssignment{}).
        Where("id = ? AND status = ? AND staff_id IS NULL", taskID, "dispatched").
        Updates(map[string]interface{}{
            "staff_id": staffID,
            "status": "claimed",
            "claimed_at": time.Now(),
        })

    if result.RowsAffected == 0 {
        return ErrTaskAlreadyClaimed
    }

    // 写历史和通知
    WriteHistory(taskID, "claimed", staffID)
    return nil
}
```

---

## 📦 前后端并行开发边界

### 前端可以先开始的

**基于P0接口Mock**:
- C端需求提交表单（mock POST /api/requests）
- 管理端任务池列表（mock GET /api/tasks）
- 任务认领按钮（mock POST /api/tasks/:id/claim）

**Mock数据准备**:
```json
// mock/requests.json
{
  "code": 200,
  "data": {
    "request_id": 1,
    "request_no": "REQ-MOCK-001",
    "status": "dispatched",
    "matched_station": {...}
  }
}
```

### 联调时机

- **第1次联调**: 认证登录 + 获取用户信息（前端鉴权基础）
- **第2次联调**: 需求提交 + 地理匹配（核心功能验证）
- **第3次联调**: 任务认领 + 完成（完整流程打通）

---

## 🚀 实施计划

### Phase 1: P0核心接口 (3-4天)

**Day 1**:
- ✅ 认证登录 (`/api/auth/login`)
- ✅ 获取用户信息 (`/api/auth/profile`)
- ✅ 地理围栏匹配逻辑优化

**Day 2**:
- ✅ 需求提交 (`POST /api/requests`)
- ✅ 需求详情 (`GET /api/requests/:id`)
- ✅ 站点/围栏创建 (管理端初始化数据)

**Day 3**:
- ✅ 任务池查询 (`GET /api/tasks`)
- ✅ 任务认领 (`POST /api/tasks/:id/claim`)
- ✅ 我的任务 (`GET /api/tasks/my`)

**Day 4**:
- ✅ 任务完成 (`POST /api/tasks/:id/complete`)
- ✅ 集成测试 + 联调准备

### Phase 2: P1完善接口 (2-3天)

**Day 5**:
- ✅ 需求列表 (`GET /api/requests`)
- ✅ 任务转派 (`POST /api/tasks/:id/transfer`)
- ✅ 图片上传 (`POST /api/upload`)

**Day 6**:
- ✅ 通知列表 (`GET /api/notifications`)
- ✅ 站点/围栏管理列表

**Day 7**: 缓冲时间 + Bug修复

### Phase 3: P2增强功能 (按需)

根据测试反馈决定是否实施

---

## 📋 交付物清单

### 代码层面
- [ ] 10个P0接口实现
- [ ] 8个P1接口实现
- [ ] Casbin权限配置
- [ ] 单元测试(核心逻辑覆盖率>80%)
- [ ] API文档(Swagger)

### 文档层面
- [ ] API接口文档(本文档)
- [ ] 权限矩阵说明
- [ ] 部署配置文档
- [ ] 测试用例文档

### 部署层面
- [ ] Docker镜像构建
- [ ] 数据库初始化脚本
- [ ] 环境变量配置模板

---

## ⚠️ 已知限制与后续优化

| 限制 | MVP方案 | 后续优化 |
|------|---------|---------|
| 地理匹配性能 | 遍历所有围栏 | 引入空间索引/GeoHash |
| 图片存储 | 本地文件系统 | 切换OSS |
| 通知推送 | 仅写表 | 接入邮件/短信服务 |
| 并发控制 | 数据库乐观锁 | 引入Redis分布式锁 |
| 日志 | 文件日志 | ELK日志聚合 |

---

**状态**: 规划已锁定
**下一步**: 生成详细SPEC.md → 开始实现P0接口

需要我现在生成完整的SPEC.md文档吗？
