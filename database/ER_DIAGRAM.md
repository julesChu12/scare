# sCare 数据库 ER 图与关系说明

**项目**: 昌平区霍营街道社区养老信息分发平台
**版本**: v1.0.0
**日期**: 2026-01-18

---

## 📊 ER 图 (Mermaid)

```mermaid
erDiagram
    %% ========================================
    %% 核心实体定义
    %% ========================================

    USERS ||--o| ELDERLY_PROFILES : "一对一(user_id)"
    USERS ||--o{ SERVICE_REQUESTS : "一对多(user_id)"
    USERS ||--o{ NOTIFICATIONS : "一对多(user_id)"
    USERS }o--|| SERVICE_STATIONS : "多对一(station_id)"

    SERVICE_STATIONS ||--o{ SERVICE_ZONES : "一对多(station_id)"
    SERVICE_STATIONS ||--o{ SERVICE_REQUESTS : "一对多(station_id)"
    SERVICE_STATIONS ||--o{ TASK_ASSIGNMENTS : "一对多(station_id)"

    SERVICE_REQUESTS ||--o{ TASK_ASSIGNMENTS : "一对多(request_id)"
    SERVICE_REQUESTS ||--o{ TASK_HISTORIES : "一对多(request_id)"

    TASK_ASSIGNMENTS ||--o{ TASK_HISTORIES : "一对多(task_id)"
    TASK_ASSIGNMENTS }o--|| USERS : "多对一(staff_id)"

    ELDERLY_PROFILES }o--o| USERS : "多对一(emergency_contact_user_id)-未来扩展"

    %% ========================================
    %% 用户表
    %% ========================================
    USERS {
        bigint id PK "主键ID"
        varchar phone UK "手机号(登录账号)"
        varchar password_hash "密码哈希"
        varchar name "用户姓名"
        varchar email "邮箱地址"
        varchar avatar "头像URL"
        varchar gender "性别(male/female/other)"
        date birth_date "出生日期"
        varchar id_card "身份证号"
        varchar role "角色(elderly/family/staff/station_manager/admin)"
        bigint station_id FK "关联站点ID(逻辑外键)"
        varchar status "状态(active/inactive/suspended)"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间(软删除)"
    }

    %% ========================================
    %% 老年人档案表
    %% ========================================
    ELDERLY_PROFILES {
        bigint id PK "主键ID"
        bigint user_id FK-UK "关联用户ID(逻辑外键,唯一)"
        varchar address "居住地址"
        decimal latitude "纬度"
        decimal longitude "经度"
        text health_info "健康信息"
        varchar emergency_contact_name "紧急联系人姓名"
        varchar emergency_contact_phone "紧急联系人电话"
        varchar emergency_contact_relation "紧急联系人关系"
        bigint emergency_contact_user_id FK "紧急联系人用户ID(未来扩展)"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 服务站点表
    %% ========================================
    SERVICE_STATIONS {
        bigint id PK "主键ID"
        varchar name "站点名称"
        varchar code UK "站点编码"
        varchar address "详细地址"
        varchar phone "联系电话"
        decimal latitude "纬度"
        decimal longitude "经度"
        varchar service_area "服务范围描述"
        int capacity "服务容量"
        varchar work_hours "工作时间"
        varchar status "状态(active/inactive)"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 服务围栏表
    %% ========================================
    SERVICE_ZONES {
        bigint id PK "主键ID"
        bigint station_id FK "关联站点ID(逻辑外键)"
        varchar name "围栏名称"
        json points "多边形顶点[[lng,lat],...]"
        int priority "优先级(数字越大越高)"
        varchar status "状态(active/inactive)"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 服务需求表
    %% ========================================
    SERVICE_REQUESTS {
        bigint id PK "主键ID"
        varchar request_no UK "需求编号(唯一)"
        bigint user_id FK "提交用户ID(逻辑外键)"
        varchar service_type "服务类型(常量)"
        varchar status "需求状态"
        text description "需求详细描述"
        decimal submit_location_lat "提交位置纬度"
        decimal submit_location_lng "提交位置经度"
        varchar contact_name "联系人姓名"
        varchar contact_phone "联系电话"
        varchar address "详细地址"
        datetime appointment_time "预约时间"
        varchar urgency "紧急程度(normal/urgent)"
        bigint station_id FK "分配站点ID(逻辑外键)"
        text reject_reason "拒绝原因"
        json images "图片URL列表"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 任务分配表
    %% ========================================
    TASK_ASSIGNMENTS {
        bigint id PK "主键ID"
        bigint request_id FK "关联需求ID(逻辑外键)"
        bigint station_id FK "分配站点ID(逻辑外键)"
        bigint staff_id FK "认领工作人员ID(逻辑外键)"
        varchar status "任务状态"
        datetime claimed_at "认领时间"
        datetime completed_at "完成时间"
        int rating "服务评分(0-5)"
        text feedback "用户评价"
        text staff_notes "工作人员备注"
        json images "完成照片URL列表"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 任务历史表
    %% ========================================
    TASK_HISTORIES {
        bigint id PK "主键ID"
        bigint task_id FK "关联任务ID(逻辑外键)"
        bigint request_id FK "关联需求ID(逻辑外键)"
        varchar action "操作类型"
        bigint operator_id FK "操作人ID(逻辑外键)"
        bigint from_staff_id FK "转派来源工作人员ID"
        bigint to_staff_id FK "转派目标工作人员ID"
        bigint from_station_id FK "转派来源站点ID"
        bigint to_station_id FK "转派目标站点ID"
        varchar status_before "操作前状态"
        varchar status_after "操作后状态"
        text remark "备注说明"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }

    %% ========================================
    %% 通知表
    %% ========================================
    NOTIFICATIONS {
        bigint id PK "主键ID"
        bigint user_id FK "接收用户ID(逻辑外键)"
        varchar title "通知标题"
        text body "通知内容"
        varchar type "通知类型"
        bigint related_id "关联业务ID"
        varchar related_type "关联业务类型(request/task)"
        varchar channel "渠道(in_app/email/sms)"
        varchar send_status "发送状态(pending/sent/failed)"
        datetime sent_at "发送时间"
        boolean is_read "是否已读"
        datetime read_at "阅读时间"
        int retry_count "重试次数"
        datetime created_at "创建时间"
        datetime updated_at "更新时间"
        datetime deleted_at "删除时间"
    }
```

---

## 🔗 实体关系详解

### 1. **Users ↔ ElderlyProfiles** (一对一)

**关系类型**: 1:1
**逻辑外键**: `elderly_profiles.user_id` → `users.id`
**约束**: `user_id` 唯一索引

**业务含义**:
- 一个用户只能有一份老年人档案
- 只有角色为 `elderly` 的用户才创建档案
- 档案扩展了用户的详细信息(居住地址、健康信息等)

**示例**:
```
用户: 张大爷(id=8, role=elderly)
档案: elderly_profiles(user_id=8, address="华龙苑北里", ...)
```

---

### 2. **Users ↔ ServiceStations** (多对一)

**关系类型**: N:1
**逻辑外键**: `users.station_id` → `service_stations.id`
**约束**: 可为NULL(elderly/family角色不关联站点)

**业务含义**:
- **Staff/StationManager** 用户必须关联一个服务站点
- **Elderly/Family/Admin** 用户不关联站点(`station_id` 为 NULL)
- 一个站点可以有多个工作人员

**示例**:
```
站点: 霍营站(id=1)
工作人员:
  - 张站长(id=2, role=station_manager, station_id=1)
  - 王小红(id=4, role=staff, station_id=1)
  - 刘小明(id=5, role=staff, station_id=1)
```

---

### 3. **ServiceStations ↔ ServiceZones** (一对多)

**关系类型**: 1:N
**逻辑外键**: `service_zones.station_id` → `service_stations.id`
**约束**: NOT NULL

**业务含义**:
- 一个站点可以有**多个服务围栏**(覆盖不同区域)
- 多个围栏可以设置不同的**优先级**
- 地理围栏用于自动匹配需求到站点

**示例**:
```
站点: 霍营站(id=1)
围栏:
  - A区-华龙苑北里(priority=10)
  - B区-龙锦苑东一区(priority=10)
```

**匹配逻辑**:
1. 用户提交需求时提供经纬度
2. 系统遍历所有围栏，使用**射线法**判断点是否在多边形内
3. 如果命中多个围栏，选择**优先级最高**的
4. 如果未命中任何围栏，使用**兜底规则**(分配到最近站点)

---

### 4. **Users ↔ ServiceRequests** (一对多)

**关系类型**: 1:N
**逻辑外键**: `service_requests.user_id` → `users.id`
**约束**: NOT NULL

**业务含义**:
- 一个用户(老年人或家属)可以提交**多个服务需求**
- 需求包含服务类型、描述、地址、预约时间等信息
- 需求创建后通过地理围栏**自动分配到站点**

**示例**:
```
用户: 张大爷(id=8)
需求:
  - REQ-001: 送餐服务(status=completed)
  - REQ-004: 家政保洁(status=pending)
```

**状态流转**:
```
pending → dispatched → claimed → processing → completed/cancelled
```

---

### 5. **ServiceStations ↔ ServiceRequests** (一对多)

**关系类型**: 1:N
**逻辑外键**: `service_requests.station_id` → `service_stations.id`
**约束**: 可为NULL(未匹配前)

**业务含义**:
- 需求通过地理围栏匹配后，`station_id` 被填充
- 一个站点可以接收**多个需求**
- 站点容量(`capacity`)控制同时处理的任务数

**匹配流程**:
```
1. 用户提交需求(station_id = NULL, status = pending)
2. 地理围栏引擎匹配站点
3. 更新 station_id，状态变为 dispatched
4. 站点任务池显示新任务
```

---

### 6. **ServiceRequests ↔ TaskAssignments** (一对多)

**关系类型**: 1:N
**逻辑外键**: `task_assignments.request_id` → `service_requests.id`
**约束**: NOT NULL (注意: **去掉了uniqueIndex，支持一对多**)

**业务含义**:
- **MVP阶段**: 一个需求对应一个任务(1:1)
- **未来扩展**: 支持任务转派、重新派单(1:N)
- 任务记录了认领人、完成时间、用户评价等

**示例** (未来扩展):
```
需求: REQ-002(就医陪护)
任务:
  - Task-001: 分配给王小红，后转派(status=transferred)
  - Task-002: 转派给刘小明，完成(status=completed)
```

**转派场景**:
- 工作人员请假/忙碌，需要转给其他人
- 跨站点转派(如站点人手不足)
- 任务升级(普通→紧急，需要更有经验的人)

---

### 7. **ServiceStations ↔ TaskAssignments** (一对多)

**关系类型**: 1:N
**逻辑外键**: `task_assignments.station_id` → `service_stations.id`
**约束**: NOT NULL

**业务含义**:
- 任务分配到站点后，显示在**站点任务池**
- 站点工作人员可以认领任务
- 站点负责人可以手动分配任务

---

### 8. **Users(Staff) ↔ TaskAssignments** (一对多)

**关系类型**: N:1
**逻辑外键**: `task_assignments.staff_id` → `users.id`
**约束**: 可为NULL(未认领前)

**业务含义**:
- 任务初始状态 `staff_id = NULL`(待认领)
- 工作人员认领后填充 `staff_id`
- 一个工作人员可以同时处理**多个任务**

**认领流程**:
```
1. 任务进入站点任务池(staff_id = NULL, status = dispatched)
2. 工作人员点击"认领"
3. 更新 staff_id, claimed_at, status = claimed
4. 开始处理任务(status = processing)
5. 完成任务(status = completed, completed_at)
```

---

### 9. **TaskAssignments ↔ TaskHistories** (一对多)

**关系类型**: 1:N
**逻辑外键**: `task_histories.task_id` → `task_assignments.id`
**约束**: NOT NULL

**业务含义**:
- 一个任务的**所有操作**都记录在历史表
- 支持追溯任务的完整生命周期
- 记录转派来源/目标、状态变更、操作人等

**历史记录示例**:
```
Task-001 的历史:
  1. [2026-01-18 10:00] 系统自动派单到霍营站(action=dispatched)
  2. [2026-01-18 10:30] 王小红认领任务(action=claimed)
  3. [2026-01-18 11:00] 王小红转派给刘小明(action=transferred)
  4. [2026-01-18 12:00] 刘小明完成任务(action=completed)
```

---

### 10. **ServiceRequests ↔ TaskHistories** (一对多)

**关系类型**: 1:N
**逻辑外键**: `task_histories.request_id` → `service_requests.id`
**约束**: NOT NULL

**业务含义**:
- 通过 `request_id` 可以查询**该需求的所有任务历史**
- 即使任务被删除，历史记录仍保留
- 支持需求维度的全流程追溯

---

### 11. **Users ↔ Notifications** (一对多)

**关系类型**: 1:N
**逻辑外键**: `notifications.user_id` → `users.id`
**约束**: NOT NULL

**业务含义**:
- 一个用户可以接收**多条通知**
- 支持多渠道(站内信、邮件、短信)
- 通过 `related_id` + `related_type` 关联业务实体

**通知类型**:
- `request_created`: 需求创建成功
- `task_dispatched`: 任务派单到站点
- `task_claimed`: 任务被认领
- `task_completed`: 任务完成
- `task_transferred`: 任务被转派

**多态关联示例**:
```
通知: "您的需求已创建"
  - user_id: 8 (张大爷)
  - type: request_created
  - related_id: 1 (REQ-001)
  - related_type: request

通知: "任务已认领"
  - user_id: 8 (张大爷)
  - type: task_claimed
  - related_id: 1 (Task-001)
  - related_type: task
```

---

### 12. **ElderlyProfiles ↔ Users** (多对一, 未来扩展)

**关系类型**: N:1
**逻辑外键**: `elderly_profiles.emergency_contact_user_id` → `users.id`
**约束**: 可为NULL

**业务含义**:
- **MVP阶段**: 使用简单字段存储紧急联系人信息
- **未来扩展**: 紧急联系人可以关联到用户账号
- 实现家属账号与老年人档案的关联

**应用场景**:
```
MVP阶段:
  张大爷档案:
    emergency_contact_name: "张小华"
    emergency_contact_phone: "13800000011"
    emergency_contact_user_id: NULL

未来扩展:
  张大爷档案:
    emergency_contact_name: "张小华"
    emergency_contact_phone: "13800000011"
    emergency_contact_user_id: 11 (关联到张小华的用户账号)

  好处:
    - 张小华登录后可以查看张大爷的需求
    - 张小华可以代替张大爷提交需求
    - 张小华收到张大爷相关的通知
```

---

## 📈 核心业务流程

### 流程1: 需求提交 → 任务派单 → 任务认领 → 任务完成

```
┌─────────────┐
│  用户提交需求  │ (service_requests: status=pending)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 地理围栏匹配  │ (service_zones: 射线法算法)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  自动派单到站点 │ (service_requests: station_id填充, status=dispatched)
│               │ (task_assignments: 创建任务记录)
│               │ (task_histories: 记录派单操作)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 工作人员认领  │ (task_assignments: staff_id填充, status=claimed)
│               │ (task_histories: 记录认领操作)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  处理任务     │ (task_assignments: status=processing)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  完成任务     │ (task_assignments: status=completed, completed_at填充)
│               │ (service_requests: status=completed)
│               │ (task_histories: 记录完成操作)
│               │ (notifications: 通知用户任务完成)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  用户评价     │ (task_assignments: rating, feedback填充)
└─────────────┘
```

---

### 流程2: 任务转派

```
┌─────────────┐
│ 工作人员A认领 │ (task_assignments: staff_id=A, status=claimed)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 站长/A发起转派│ (创建新任务: staff_id=B)
│               │ (旧任务: status=transferred)
│               │ (task_histories: 记录转派, from_staff_id=A, to_staff_id=B)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 工作人员B处理 │ (新任务: status=processing)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ 完成任务      │ (新任务: status=completed)
└─────────────┘
```

---

## 🔍 关键设计决策说明

### 为什么 TaskAssignments 是一对多而非一对一？

**原因**:
1. **支持任务转派**: 一个需求可能先分给A，后转给B
2. **支持任务重派**: 任务失败后可以重新派单
3. **历史追溯**: 保留所有任务记录，便于审计

**MVP实现**:
- 虽然设计支持一对多，但MVP阶段只创建一个任务
- 业务逻辑层控制：一个需求同时只有一个 `active` 任务

---

### 为什么使用逻辑外键而非数据库外键？

**优势**:
1. **灵活性**: 避免级联删除带来的数据丢失风险
2. **性能**: 减少数据库锁竞争，提升并发性能
3. **软删除友好**: GORM软删除机制下，外键约束会导致问题
4. **分布式友好**: 未来可能拆分数据库，逻辑外键更灵活

**数据完整性保障**:
- 应用层校验外键有效性
- 索引优化关联查询性能
- 业务逻辑层保证数据一致性

---

### 为什么地理围栏使用JSON而非GEOMETRY？

**对比**:

| 方案 | 优势 | 劣势 |
|------|------|------|
| JSON + 射线法 | 灵活、易调试、跨数据库 | 无空间索引 |
| GEOMETRY + ST_Contains | 空间索引优化 | 复杂、调试困难 |

**决策**:
- 当前业务规模小(<100个围栏)，JSON方案性能足够(18ms)
- 应用层算法可控，便于优化和调试
- 未来如需优化，可迁移到GEOMETRY(表结构无需大改)

---

## 📚 参考文档

- [数据库表结构脚本](schema/schema.sql)
- [种子数据](seeds/seed.sql)
- [数据库初始化文档](README.md)
- [系统架构设计](../backend/docs/02-系统架构设计.md)

---

**维护者**: sCare开发团队
**最后更新**: 2026-01-18
