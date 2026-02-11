# SPEC: sCare MVP后端API实现

**Status**: Locked
**Owner**: sCare Development Team
**Baseline**: Go 1.25 + Gin + GORM
**日期**: 2026-01-18

---

## 1. Requirement Summary

### 问题陈述

数据库模型已完成(8表)，后端框架已搭建(54个Go文件)，现需要实现MVP版本的核心API接口，支持完整的需求提交→地理匹配→任务派发→认领→完成流程。

### 项目目标

实现**最小可演示路径**，打通完整业务流程：

1. **C端用户提交养老服务需求** (含地理围栏自动匹配)
2. **系统自动派单到服务站点** (基于地理位置)
3. **工作人员认领并处理任务**
4. **任务完成后用户查看结果**
5. **支持站长手动转派任务** (P1功能)

### 为什么现在要做

- **数据库已就绪**: GORM模型完成，AutoMigrate可用
- **框架已搭建**: 部分handler/service/repository已存在
- **前端等待联调**: 需要尽快提供可用接口
- **演示需求**: 需要完整流程验证业务逻辑

---

## 2. Success Criteria

### Functional

- **F1: 认证授权** - 5种角色登录，JWT认证，Casbin权限控制
- **F2: 需求提交** - C端提交需求，地理围栏匹配<50ms，拒绝围栏外请求
- **F3: 任务派发** - 需求自动分配到站点，写入任务池
- **F4: 任务认领** - 工作人员认领任务，并发控制防止重复认领
- **F5: 任务完成** - 提交完成信息(备注+图片)，状态流转正确
- **F6: 任务转派** - 站长可手动转派任务，记录历史
- **F7: 通知系统** - 关键节点写入站内信通知

### Technical

- **Performance**:
  - API平均响应时间 < 500ms
  - 地理围栏匹配 < 50ms
  - 支持100并发请求

- **Reliability**:
  - 任务认领并发安全(乐观锁)
  - 事务保证数据一致性
  - 错误处理完整(4xx/5xx)

- **Observability**:
  - 关键操作记录日志
  - API访问日志
  - 错误堆栈记录

---

## 3. Baseline Used / Exceptions

### Baseline Used

- **语言框架**: Go 1.25 + Gin + GORM
- **数据库**: MySQL 8.0 + Redis 7.0
- **认证**: JWT Token
- **权限**: Casbin RBAC
- **日志**: Zap日志库
- **配置**: Viper配置管理

### Exceptions

- **地理围栏拒绝策略**: 不在围栏内直接拒绝(不做兜底分配)
- **多围栏匹配**: 随机选择(不按优先级，MVP简化)
- **图片存储**: 本地文件系统(MVP后切OSS)
- **通知推送**: 仅写表(不发邮件/短信)

---

## 4. Plan A (Chosen)

### Summary

实现**26个API接口**(P0: 10个, P1: 8个, P2: 8个)，优先完成P0核心流程，支持：

- ✅ 认证登录 + 用户信息获取
- ✅ 需求提交 + 地理围栏匹配
- ✅ 任务派发 + 认领 + 完成
- ✅ 任务转派(站长权限)
- ✅ 站内信通知
- ✅ 图片上传(本地存储)

### Why this fits current task

- **MVP路径最短**: 完整流程演示只需P0接口
- **前后端可并行**: 前端可Mock P0接口先开发UI
- **风险可控**: 核心逻辑已有基础(geofence/engine.go已实现)
- **扩展性保留**: P1/P2接口设计完整，随时可补充

### Risks

- **R1**: 地理围栏匹配性能(多围栏场景)
  - **Mitigation**: 限制围栏数量(<100)，后期引入空间索引

- **R2**: 任务认领并发冲突
  - **Mitigation**: 数据库乐观锁 + 返回409冲突提示

- **R3**: 图片上传容量问题
  - **Mitigation**: 限制单文件5MB，MVP后切OSS

### Mitigations

- 地理围栏: 遍历算法 + 性能测试验证
- 并发控制: `WHERE staff_id IS NULL` + `RowsAffected`检查
- 图片存储: 配置化路径，便于后期切换

---

## 5. API & Contract

### P0 核心接口 (10个)

#### 认证授权

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/auth/login` | 用户登录 | Public |
| GET | `/api/auth/profile` | 获取当前用户信息 | All |

#### 需求管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/requests` | 提交需求(含地理匹配) | elderly/family |
| GET | `/api/requests/:id` | 需求详情 | elderly/family/staff/admin |

#### 任务管理

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/tasks` | 站点任务池 | staff/station_manager |
| POST | `/api/tasks/:id/claim` | 认领任务 | staff |
| POST | `/api/tasks/:id/complete` | 完成任务 | staff |
| GET | `/api/tasks/my` | 我的任务列表 | staff |

#### 基础数据

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/stations` | 创建站点 | admin |
| POST | `/api/zones` | 创建围栏 | admin |

### P1 完善接口 (8个)

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/requests` | 需求列表 | elderly/family |
| POST | `/api/tasks/:id/transfer` | 任务转派 | station_manager |
| GET | `/api/notifications` | 通知列表 | All |
| PUT | `/api/notifications/:id/read` | 标记已读 | All |
| POST | `/api/upload` | 图片上传 | All |
| GET | `/api/stations` | 站点列表 | admin |
| GET | `/api/zones` | 围栏列表 | admin |
| PUT | `/api/tasks/:id/status` | 更新任务状态 | staff |

### HTTP status semantics

- **200**: 成功
- **400**: 参数错误 / 业务规则拒绝(如围栏外)
- **401**: 未认证
- **403**: 无权限
- **404**: 资源不存在
- **409**: 冲突(如任务已被认领)
- **500**: 服务器内部错误

### Error cases

**标准错误响应**:
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
    }
  }
}
```

**常见错误码**:
- `invalid_credentials`: 用户名或密码错误
- `no_geofence_matched`: 不在服务范围内
- `task_already_claimed`: 任务已被认领
- `permission_denied`: 权限不足
- `resource_not_found`: 资源不存在

---

## 6. Data & Consistency

### Data model

**核心表**:
- `users`: 用户表(5种角色)
- `elderly_profiles`: 老年人档案
- `service_stations`: 服务站点
- `service_zones`: 地理围栏(JSON存顶点)
- `service_requests`: 服务需求
- `task_assignments`: 任务分配
- `task_histories`: 任务历史
- `notifications`: 通知记录

详见: `/database/ER_DIAGRAM.md`

### Consistency level

- **强一致性**:
  - 任务认领(事务 + 乐观锁)
  - 需求提交(事务保证需求+任务+通知原子性)

- **最终一致性**:
  - 通知写入(异步，允许延迟)

### Idempotency

- **需求提交**: 基于`request_no`去重(前端生成唯一编号)
- **任务认领**: `WHERE staff_id IS NULL`防止重复认领
- **任务完成**: `WHERE status = 'claimed'`防止重复完成

### Retry / compensation

- **地理匹配失败**: 直接拒绝，返回400
- **通知写入失败**: 记录日志，不影响主流程
- **图片上传失败**: 返回错误，前端重试

### DTM usage

**No** - MVP阶段不使用分布式事务，单库事务足够

---

## 7. Observability

### Metrics

- **业务指标**:
  - 需求提交数量(按服务类型分组)
  - 地理匹配成功率
  - 任务认领率
  - 任务完成率
  - 平均响应时长

- **系统指标**:
  - API调用QPS
  - API响应时间(P50/P95/P99)
  - 错误率(按HTTP状态码)
  - 数据库连接池使用率

### Logs

- **访问日志**: Gin中间件记录所有HTTP请求
- **业务日志**: 关键操作(需求创建、任务认领、任务完成)
- **错误日志**: panic recover + stack trace
- **审计日志**: 写入`task_histories`表

**日志格式**:
```json
{
  "level": "info",
  "time": "2026-01-18T10:00:00Z",
  "msg": "task claimed",
  "task_id": 1,
  "staff_id": 4,
  "request_id": 1,
  "duration_ms": 15
}
```

### Traces

**No** - MVP阶段不引入分布式追踪，使用日志关联ID

### Alerts (work-hours oriented)

- 地理匹配失败率 > 10%
- 任务认领冲突率 > 5%
- API错误率 > 3%
- API P95响应时间 > 1s
- 数据库连接池耗尽

---

## 8. Deployment & Migration

### Migration steps (pre-deploy)

1. 执行数据库初始化
   ```bash
   cd database/scripts
   DB_PASSWORD=xxx ./init.sh
   ```

2. 导入种子数据(测试账号+站点+围栏)
   ```bash
   mysql -u root -p scare_db < database/seeds/seed.sql
   ```

3. 配置环境变量
   ```bash
   cp .env.example .env
   # 修改数据库密码、JWT密钥等
   ```

4. 构建后端服务
   ```bash
   cd backend
   go build -o bin/server cmd/server/main.go
   ```

### Deploy strategy

**Docker Compose部署**:
```bash
docker-compose up -d
```

**服务清单**:
- `backend`: Go API服务(端口8080)
- `mysql`: MySQL 8.0(端口3306)
- `redis`: Redis 7.0(端口6379)

### Rollback plan

1. 停止服务: `docker-compose down`
2. 回滚代码: `git checkout <previous-commit>`
3. 重新部署: `docker-compose up -d`
4. 数据库备份恢复(如有schema变更)

---

## 9. Risks & Edge Cases

### R1: 地理围栏匹配边界

**场景1**: 用户坐标不在任何围栏内
- **处理**: 返回400，提示最近站点信息
- **前端**: 显示提示"不在服务范围，请联系XX站点"

**场景2**: 用户坐标命中多个围栏
- **处理**: 随机选择一个，返回`match_method: "random"`
- **日志**: 记录所有命中围栏ID，便于分析

**场景3**: 围栏JSON数据异常
- **处理**: 跳过该围栏，继续匹配其他围栏
- **日志**: 记录错误围栏ID和错误详情

### R2: 任务认领并发冲突

**场景**: 两个工作人员同时认领同一任务
- **数据库**: `WHERE staff_id IS NULL` + 检查`RowsAffected`
- **返回**: 409 Conflict，提示"任务已被XX认领"
- **前端**: 刷新任务列表，显示最新状态

### R3: 图片上传失败

**场景1**: 文件过大(>5MB)
- **处理**: 返回400，提示"文件大小不能超过5MB"

**场景2**: 磁盘空间不足
- **处理**: 返回500，记录错误日志
- **运维**: 监控磁盘使用率，<20%告警

**场景3**: 文件类型不支持
- **处理**: 返回400，只允许jpg/png/jpeg

### R4: 任务转派权限问题

**场景**: staff尝试转派任务
- **处理**: Casbin拦截，返回403
- **前端**: 隐藏转派按钮(role判断)

### R5: 需求提交重复

**场景**: 网络抖动导致重复提交
- **处理**: 前端生成唯一`request_no`，后端去重
- **返回**: 200，返回已存在的需求信息

---

## 10. Next Steps for `claude -p`

### 阶段1: P0核心接口实现 (3-4天)

**Day 1: 认证基础**
- [ ] 实现 `POST /api/auth/login` (JWT生成)
- [ ] 实现 `GET /api/auth/profile` (从Token解析用户)
- [ ] 配置Casbin RBAC权限模型
- [ ] 编写认证中间件 (`middleware/auth.go`)
- [ ] 编写权限中间件 (`middleware/casbin.go`)

**Day 2: 需求提交**
- [ ] 优化地理围栏匹配引擎 (`geofence/engine.go`)
  - 添加多围栏随机选择逻辑
  - 添加围栏外拒绝逻辑
  - 添加性能测试(验证<50ms)
- [ ] 实现 `POST /api/requests` (需求提交+地理匹配)
  - 事务处理(需求+任务+历史+通知)
  - 返回匹配结果详情
- [ ] 实现 `GET /api/requests/:id` (需求详情)

**Day 3: 任务管理**
- [ ] 实现 `GET /api/tasks` (站点任务池查询)
  - 分页支持
  - 状态筛选
  - 按urgency排序
- [ ] 实现 `POST /api/tasks/:id/claim` (任务认领)
  - 乐观锁并发控制
  - 写任务历史
  - 发送通知
- [ ] 实现 `GET /api/tasks/my` (我的任务列表)

**Day 4: 任务完成 + 基础数据**
- [ ] 实现 `POST /api/tasks/:id/complete` (任务完成)
  - 更新任务状态
  - 更新需求状态
  - 写任务历史
- [ ] 实现 `POST /api/stations` (创建站点)
- [ ] 实现 `POST /api/zones` (创建围栏)
- [ ] 集成测试 + API文档生成

### 阶段2: P1完善接口 (2-3天)

**Day 5: 需求与转派**
- [ ] 实现 `GET /api/requests` (需求列表查询)
- [ ] 实现 `POST /api/tasks/:id/transfer` (任务转派)
  - 创建新任务
  - 更新旧任务状态为transferred
  - 写转派历史
  - 发送通知
- [ ] 实现 `POST /api/upload` (图片上传)
  - 本地文件存储
  - 文件类型验证
  - 大小限制

**Day 6: 通知与管理**
- [ ] 实现 `GET /api/notifications` (通知列表)
- [ ] 实现 `PUT /api/notifications/:id/read` (标记已读)
- [ ] 实现 `GET /api/stations` (站点列表)
- [ ] 实现 `GET /api/zones` (围栏列表)
- [ ] 实现 `PUT /api/tasks/:id/status` (更新任务状态)

**Day 7: 测试与优化**
- [ ] 单元测试补充
- [ ] 集成测试完善
- [ ] 性能测试(压测100并发)
- [ ] Bug修复

### 阶段3: P2增强功能 (按需实施)

根据测试反馈决定是否实现

---

## 11. 参考文档

- **API规划详情**: `/docs/API_PLANNING.md`
- **数据库设计**: `/database/README.md`
- **ER图关系**: `/database/ER_DIAGRAM.md`
- **系统架构**: `/backend/docs/02-系统架构设计.md`
- **原始SPEC**: `/SPEC.md`

---

## 12. 实施检查清单

### 代码实现
- [ ] P0接口(10个)全部实现
- [ ] P1接口(8个)全部实现
- [ ] Casbin权限配置完成
- [ ] JWT认证中间件就绪
- [ ] 地理围栏匹配逻辑优化
- [ ] 错误处理统一封装

### 测试验证
- [ ] 单元测试覆盖率>80%
- [ ] 地理匹配性能测试通过(<50ms)
- [ ] 任务认领并发测试通过
- [ ] 完整流程集成测试通过
- [ ] API文档Swagger生成

### 部署准备
- [ ] Docker镜像构建成功
- [ ] 数据库初始化脚本验证
- [ ] 环境变量配置模板
- [ ] 日志配置完成
- [ ] 监控指标定义

### 文档完善
- [ ] API接口文档(Swagger)
- [ ] 权限矩阵说明
- [ ] 部署操作手册
- [ ] 故障排查指南

---

**Spec状态**: ✅ Locked
**可执行性**: ✅ Ready for `claude -p`
**预计工期**: 7-10天
**风险等级**: 🟢 Low (核心逻辑已有基础)

---

**下一步执行命令**:
```bash
claude -p "根据SPEC.md实现P0核心接口"
```
