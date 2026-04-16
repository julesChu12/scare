# 数据库初始化说明文档

**项目**: 昌平区霍营街道社区养老信息分发平台
**版本**: v1.0.0
**日期**: 2026-01-18

---

## 📋 目录结构

```
database/
├── schema/
│   └── schema.sql           # 数据库表结构定义
├── seeds/
│   └── seed.sql             # 测试环境种子数据
├── scripts/
│   └── init.sh              # 自动化初始化脚本
└── README.md                # 本文档
```

---

## 🚀 快速开始

### 方式一: 使用初始化脚本 (推荐)

```bash
cd database/scripts

# 完整初始化(创建表+导入数据)
./init.sh

# 指定数据库密码
DB_PASSWORD=your_password ./init.sh

# 只创建表结构
./init.sh --schema-only

# 只导入种子数据
./init.sh --seed-only

# 强制重新创建(会清空现有数据)
./init.sh --force
```

### 方式二: 手动执行SQL

```bash
# 1. 登录MySQL
mysql -u root -p

# 2. 执行表结构创建
source /path/to/backend/database/schema/schema.sql

# 3. 导入种子数据
source /path/to/backend/database/seeds/seed_all.sql
```

### 方式三: GORM AutoMigrate (开发环境)

后端服务启动时会自动执行表迁移:

```go
// backend/cmd/server/main.go
if err := domain.AutoMigrate(db.DB); err != nil {
    log.Fatalf("数据库迁移失败: %v", err)
}
```

---

## 📊 数据库设计详解

### 设计原则

1. ✅ **逻辑外键**: 所有关联字段使用逻辑外键，不使用数据库外键约束
2. ✅ **软删除**: 所有表支持软删除 (`deleted_at` 字段)
3. ✅ **审计字段**: 统一的 `created_at`, `updated_at` 时间戳
4. ✅ **索引优化**: 关键查询字段建立索引
5. ✅ **JSON存储**: 地理围栏、图片列表使用JSON类型

---

### 表结构说明

#### 1. users (用户表)

**用途**: 存储所有用户信息，支持5种角色

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| phone | VARCHAR(20) | 手机号(登录账号) | 国内手机号登录惯例，唯一索引 |
| password_hash | VARCHAR(255) | 密码哈希 | bcrypt哈希存储；C端快速开通用户在首次设置密码前允许为空 |
| avatar | VARCHAR(255) | 头像URL | 扩展字段，提升用户体验 |
| gender | VARCHAR(10) | 性别 | 可选字段，业务统计需要 |
| birth_date | DATE | 出生日期 | 年龄计算，可选 |
| id_card | VARCHAR(18) | 身份证号 | 实名认证，可选，有索引 |
| role | VARCHAR(20) | 角色 | elderly/family/staff/station_manager/admin |
| station_id | BIGINT | 关联站点ID | 逻辑外键，staff/station_manager有值 |
| status | VARCHAR(20) | 状态 | active/inactive/suspended |

**角色说明**:
- `elderly`: 老年人用户
- `family`: 家属用户
- `staff`: 工作人员
- `station_manager`: 站点负责人
- `admin`: 系统管理员

---

#### 2. elderly_profiles (老年人档案表)

**用途**: 一对一关联用户表，存储老年人详细信息

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| user_id | BIGINT | 关联用户ID | 唯一索引，一对一关系 |
| address | VARCHAR(200) | 居住地址 | 服务地址参考 |
| latitude/longitude | DECIMAL(10,7) | 经纬度 | 精度约1.11cm，满足定位需求 |
| health_info | TEXT | 健康信息 | 自由文本，记录健康状况 |
| emergency_contact_name | VARCHAR(50) | 紧急联系人姓名 | MVP阶段简单字段 |
| emergency_contact_phone | VARCHAR(20) | 紧急联系人电话 | 紧急情况联系 |
| emergency_contact_relation | VARCHAR(20) | 关系 | 子女/配偶/其他 |
| emergency_contact_user_id | BIGINT | 关联用户ID | 未来扩展，支持用户关联 |

**紧急联系人设计**:
- **MVP阶段**: 使用简单字段(姓名+电话)
- **未来扩展**: 支持关联到用户表(`emergency_contact_user_id`)，实现家属账号关联

---

#### 3. service_stations (服务站点表)

**用途**: 社区服务站点信息

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| name | VARCHAR(100) | 站点名称 | 业务标识 |
| code | VARCHAR(50) | 站点编码 | 唯一编码，便于系统识别 |
| address | VARCHAR(200) | 详细地址 | 站点位置信息 |
| phone | VARCHAR(20) | 联系电话 | 站点联系方式 |
| latitude/longitude | DECIMAL(10,7) | 经纬度 | 兜底匹配时计算最近站点 |
| service_area | VARCHAR(200) | 服务范围描述 | 业务说明 |
| capacity | INT | 服务容量 | 同时处理任务数上限 |
| work_hours | VARCHAR(100) | 工作时间 | 业务参考 |
| status | VARCHAR(20) | 状态 | active/inactive |

---

#### 4. service_zones (服务围栏表)

**用途**: 地理围栏，存储多边形顶点坐标

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| station_id | BIGINT | 关联站点ID | 逻辑外键 |
| name | VARCHAR(100) | 围栏名称 | 业务标识 |
| points | JSON | 多边形顶点 | `[[lng,lat],[lng,lat],...]` 格式 |
| priority | INT | 优先级 | 多围栏重叠时选择高优先级 |
| status | VARCHAR(20) | 状态 | active/inactive |

**Points格式说明**:
```json
[
  [116.378, 40.048],  // [经度, 纬度]
  [116.382, 40.048],
  [116.382, 40.052],
  [116.378, 40.052],
  [116.378, 40.048]   // 闭合多边形
]
```

**为什么不用MySQL GEOMETRY类型?**
- 当前实现: JSON + 应用层射线法算法
- 优势: 灵活、易于调试、跨数据库兼容
- 性能: 实测匹配耗时18ms(目标<50ms)，满足需求

---

#### 5. service_requests (服务需求表)

**用途**: 用户提交的养老服务需求

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| request_no | VARCHAR(50) | 需求编号 | 唯一编号，支持幂等性 |
| user_id | BIGINT | 提交用户ID | 逻辑外键 |
| service_type | VARCHAR(20) | 服务类型 | 硬编码常量(见service_types.go) |
| status | VARCHAR(20) | 需求状态 | 状态机流转 |
| description | TEXT | 详细描述 | 需求说明 |
| submit_location_lat/lng | DECIMAL(10,7) | 提交位置 | 用于地理围栏匹配 |
| contact_name/phone | VARCHAR | 联系方式 | 服务联系信息 |
| address | VARCHAR(200) | 详细地址 | 服务地址 |
| appointment_time | DATETIME | 预约时间 | 预约型服务 |
| urgency | VARCHAR(20) | 紧急程度 | normal/urgent |
| station_id | BIGINT | 分配站点ID | 匹配后填充 |
| reject_reason | TEXT | 拒绝原因 | 需求被拒绝时记录 |
| images | JSON | 图片URL列表 | 支持多图上传 |

**状态流转**:
```
pending → dispatched → claimed → processing → completed/cancelled
```

---

#### 6. task_assignments (任务分配表)

**用途**: 任务派单与处理，支持一对多关系

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| request_id | BIGINT | 关联需求ID | 支持一对多(一个需求可多次派单) |
| station_id | BIGINT | 分配站点ID | 逻辑外键 |
| staff_id | BIGINT | 认领工作人员ID | 未认领前为NULL |
| status | VARCHAR(20) | 任务状态 | 状态机流转 |
| claimed_at | DATETIME | 认领时间 | 响应时长统计 |
| completed_at | DATETIME | 完成时间 | 完成时长统计 |
| rating | INT | 服务评分 | 1-5分 |
| feedback | TEXT | 用户评价 | 服务质量反馈 |
| staff_notes | TEXT | 工作人员备注 | 内部记录 |
| images | JSON | 完成照片 | 服务证明 |

**一对多设计**:
- 去掉 `uniqueIndex`，支持一个需求多次派单
- MVP阶段: 一个需求只派一次单
- 未来扩展: 支持任务转派、重新派单

---

#### 7. task_histories (任务历史表)

**用途**: 记录任务的所有操作历史

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| task_id | BIGINT | 关联任务ID | 逻辑外键 |
| request_id | BIGINT | 关联需求ID | 逻辑外键 |
| action | VARCHAR(20) | 操作类型 | dispatched/claimed/transferred/completed/cancelled |
| operator_id | BIGINT | 操作人ID | 操作记录 |
| from_staff_id/to_staff_id | BIGINT | 转派工作人员 | 转派记录 |
| from_station_id/to_station_id | BIGINT | 转派站点 | 跨站转派 |
| status_before/after | VARCHAR(20) | 状态变更 | 状态流转记录 |
| remark | TEXT | 备注说明 | 操作说明 |

**操作类型**:
- `dispatched`: 派单
- `claimed`: 认领
- `transferred`: 转派
- `completed`: 完成
- `cancelled`: 取消

---

#### 8. notifications (通知表)

**用途**: 多渠道通知记录

| 字段 | 类型 | 说明 | 设计理由 |
|------|------|------|---------|
| user_id | BIGINT | 接收用户ID | 逻辑外键 |
| title | VARCHAR(100) | 通知标题 | 通知摘要 |
| body | TEXT | 通知内容 | 详细内容 |
| type | VARCHAR(20) | 通知类型 | 业务分类 |
| related_id/type | BIGINT/VARCHAR | 关联业务 | 多态关联 |
| channel | VARCHAR(20) | 渠道 | in_app/email/sms |
| send_status | VARCHAR(20) | 发送状态 | pending/sent/failed |
| sent_at | DATETIME | 发送时间 | 发送记录 |
| is_read | BOOLEAN | 是否已读 | 已读标记 |
| read_at | DATETIME | 阅读时间 | 阅读记录 |
| retry_count | INT | 重试次数 | 失败重试 |

**通知类型**:
- `request_created`: 需求创建
- `task_claimed`: 任务认领
- `task_completed`: 任务完成
- `task_dispatched`: 任务派单

---

## 🔑 索引设计

### 主键索引
所有表的 `id` 字段自动创建主键索引

### 唯一索引
- `users.phone`: 手机号唯一
- `users.request_no`: 需求编号唯一
- `service_stations.code`: 站点编码唯一
- `elderly_profiles.user_id`: 一对一关系

### 普通索引
- **外键字段**: 所有 `_id` 后缀字段
- **状态字段**: `status`, `send_status`
- **时间字段**: `deleted_at`(软删除查询优化)
- **业务字段**: `type`, `channel`, `is_read`

---

## 📝 种子数据说明

### 测试账号

| 角色 | 手机号 | 密码 | 姓名 | 说明 |
|------|--------|------|------|------|
| 管理员 | 13800000001 | Test@123 | 系统管理员 | 全局权限 |
| 站长 | 13800000002 | Test@123 | 张站长 | 霍营站负责人 |
| 站长 | 13800000003 | Test@123 | 李站长 | 龙泽站负责人 |
| 工作人员 | 13800000004 | Test@123 | 王小红 | 霍营站工作人员 |
| 工作人员 | 13800000005 | Test@123 | 刘小明 | 霍营站工作人员 |
| 老年人 | 13800000008 | Test@123 | 张大爷 | 测试老年人 |
| 家属 | 13800000011 | Test@123 | 张小华 | 张大爷子女 |

### 测试站点

1. **霍营街道养老服务中心** (station-001)
   - 覆盖华龙苑北里、龙锦苑东一区
   - 容量: 20个任务

2. **龙泽园养老服务站** (station-002)
   - 覆盖龙泽苑西区、东区
   - 容量: 15个任务

3. **回龙观养老服务站** (station-003)
   - 覆盖回龙观社区
   - 容量: 10个任务

### 测试需求

- 已完成: 张大爷的送餐服务
- 处理中: 李奶奶的就医陪护
- 待认领: 王大爷的日常照护
- 待派单: 张大爷的家政保洁

---

## ⚙️ 环境变量配置

初始化脚本支持以下环境变量:

```bash
export DB_HOST=localhost       # 数据库主机
export DB_PORT=3306            # 数据库端口
export DB_USER=root            # 数据库用户
export DB_PASSWORD=your_pass   # 数据库密码
export DB_NAME=scare_db        # 数据库名称
```

---

## 🔄 数据迁移策略

### 开发环境
使用 GORM AutoMigrate，服务启动时自动建表

### 生产环境
1. 执行 `schema.sql` 创建表结构
2. **不要**导入 `seed.sql`(测试数据)
3. 使用数据库迁移工具管理版本(如 Flyway, Liquibase)

---

## 🐛 常见问题

### Q1: 密码哈希如何生成?
```go
import "golang.org/x/crypto/bcrypt"

hash, _ := bcrypt.GenerateFromPassword([]byte("Test@123"), bcrypt.DefaultCost)
```

### Q2: 为什么不用外键约束?
- 灵活性: 避免级联删除问题
- 性能: 减少数据库锁竞争
- 分布式: 未来可能拆分数据库

### Q3: 如何验证数据导入成功?
```sql
SELECT '用户数据' AS 'Table', COUNT(*) AS 'Count' FROM users
UNION ALL
SELECT '服务站点', COUNT(*) FROM service_stations;
```

---

## 📚 参考文档

- [数据库设计文档](../backend/docs/03-数据库设计.md)
- [系统架构设计](../backend/docs/02-系统架构设计.md)
- [SPEC规格文档](../SPEC.md)

---

**维护者**: sCare开发团队
**最后更新**: 2026-01-18
