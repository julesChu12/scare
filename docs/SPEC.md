# SPEC: 基于地理围栏的社区养老信息分发平台

**Status**: Locked
**Owner**: Jules
**Baseline**: Fly Personal Fullstack (Go 1.25.1)
**日期**: 2025-01-16

---

## 1. Requirement Summary

### 问题陈述

随着我国人口老龄化进程加快，社区养老服务需求日益增长。根据国家统计局数据，60岁及以上人口占比持续上升。然而，传统社区养老服务模式存在以下问题：

1. **渠道分散**：需求上报依赖电话、微信群、线下登记等多种方式，缺乏统一入口
2. **人工派单**：依靠人工整理和转发完成任务派发，响应不及时，效率低下
3. **责任不清**：缺乏明确的分发规则和责任界定，处理过程不透明
4. **空间脱节**：未结合社区站点和服务人员的空间分布，无法实现精准派单

### 项目目标

设计并实现一个基于地理围栏的社区养老信息分发平台，实现：

1. **统一信息入口**：通过二维码为老年人及家属提供便捷的需求提交方式
2. **智能自动分发**：基于地理围栏规则（点在多边形算法）自动将需求分发至对应社区站点任务池
3. **处理过程透明**：全流程状态跟踪，责任界定清晰
4. **降低使用门槛**：老年人友好的界面设计，支持PWA离线访问

### 为什么现在要做

- 毕业设计项目要求
- 社区养老服务数字化转型需求迫切
- 地理围栏技术在社区服务场景的应用探索价值

---

## 2. Success Criteria

### Functional

- **F1: 用户认证** - 支持4种角色（elderly/family/staff/admin）登录和权限控制
- **F2: 需求提交** - 老年人/家属可在3分钟内完成需求提交流程
- **F3: 地理围栏匹配** - 基于点在多边形算法实现需求自动分发，匹配时间 < 50ms
- **F4: 任务管理** - 支持任务认领、状态更新、完成反馈，全流程可追溯
- **F5: 多渠道通知** - 需求创建、认领、完成等关键节点自动触发通知
- **F6: 围栏管理** - 管理员可通过地图绘制多边形围栏，实时预览覆盖范围

### Technical

- **Performance**:
  - 页面首屏加载时间 < 2秒
  - API 平均响应时间 < 500ms
  - 地理围栏匹配算法 < 50ms（已测试验证）

- **Reliability**:
  - 系统可用性 ≥ 99%
  - 支持 100 并发用户
  - 数据持久化保证（MySQL + Redis）

- **Observability**:
  - 完整的操作日志记录
  - 关键业务指标统计
  - 错误监控和告警机制

---

## 3. Baseline Used / Exceptions

### Baseline Used

- **架构模式**: Clean Architecture（分层架构）
- **API 风格**: RESTful 资源风格
- **数据存储**: MySQL 8.0 (空间数据支持) + Redis 缓存
- **可观测性**: 日志记录 + 操作审计
- **部署方式**: Docker 容器化 + Nginx 反向代理
- **权限系统**: 自定义 RBAC 三表模型
- **认证方式**: JWT Token

### Exceptions (if any)

无例外，完全采用标准 Baseline

---

## 4. Plan A (Chosen)

### Summary

采用 **前后端分离架构** + **地理围栏智能分发** 的技术方案：

**后端**：Go 1.25 + Gin + GORM + MySQL 8.0 + Redis
- 分层架构：Handler → Service → Repository
- 地理围栏匹配：射线法算法 + MySQL 空间索引
- 权限控制：自定义 RBAC（permissions/roles/role_permissions 三表）

**前端**：Vue 3 + TypeScript + Pinia
- C端：Naive UI + PWA（老年人友好）
- 管理门户：Element Plus + 三端合一（staff/station_manager/admin）

### Why this fits current task

- **毕业设计规模适中**：小型试点（<1000用户），技术栈学习曲线平缓
- **地理围栏核心技术**：MySQL 8.0 POLYGON 类型 + 空间索引满足需求
- **权限系统完整**：自定义 RBAC 三表模型提供细粒度权限控制
- **前后端独立开发**：符合毕业设计时间安排，可并行推进

### Risks

- **R1**: 地理围栏匹配性能（多围栏场景）
  - **Mitigation**: 空间索引优化 + 兜底机制
- **R2**: 前端权限复杂度（三端合一）
  - **Mitigation**: 自定义权限服务 + 动态路由守卫 + v-permission 指令
- **R3**: PWA 兼容性
  - **Mitigation**: 渐进增强 + 降级方案

### Mitigations

- 地理围栏：优先级排序 + 默认站点兜底
- 权限系统：自定义 RBAC 三表模型（permissions/roles/role_permissions）
- PWA：Service Worker 离线缓存 + 在线同步

---

## 5. Plan B (Fallback)

### When to switch

- 地理围栏匹配性能无法满足要求（> 100ms）
- MySQL 空间查询出现兼容性问题

### Trade-offs

- **PostGIS** 更强大但运维复杂度增加
- **外包匹配服务**（高德地图 API）增加依赖和成本

---

## 6. API & Contract

### Resources

- `/api/requests` - 服务需求
- `/api/tasks` - 任务管理
- `/api/stations` - 服务站点
- `/api/zones` - 地理围栏
- `/api/users` - 用户管理
- `/api/auth` - 认证授权

### Endpoints (核心接口)

#### 需求管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/requests | 创建需求（触发自动分发） |
| GET | /api/requests/:id | 查询需求详情 |
| GET | /api/requests | 需求列表（支持筛选） |
| PUT | /api/requests/:id | 更新需求 |
| DELETE | /api/requests/:id | 取消需求 |

#### 任务管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tasks?station=xxx | 站点任务池 |
| POST | /api/tasks/:id/claim | 认领任务 |
| POST | /api/tasks/:id/complete | 完成任务 |
| GET | /api/tasks/my | 我的任务列表 |

#### 地理围栏

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/zones | 围栏列表 |
| POST | /api/zones | 创建围栏（POLYGON） |
| PUT | /api/zones/:id | 更新围栏 |
| DELETE | /api/zones/:id | 删除围栏 |
| POST | /api/zones/:id/match | 测试点是否在围栏内 |

#### 认证授权

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/logout | 用户登出 |
| GET | /api/auth/profile | 获取当前用户信息 |
| POST | /api/auth/refresh | 刷新Token |

### HTTP status semantics

- **200**: 成功
- **400**: 参数错误
- **401**: 未认证
- **403**: 无权限
- **404**: 资源不存在
- **500**: 服务器内部错误

### Error cases

```json
{
  "code": 400,
  "message": "参数错误",
  "error": "address字段不能为空"
}
```

---

## 7. Data & Consistency

### Data model

见 `backend/docs/03-数据库设计.md`：

- **用户表** (users): 4种角色（elderly/family/staff/admin）
- **站点表** (service_stations): 服务站点信息
- **围栏表** (service_zones): POLYGON 空间数据类型
- **需求表** (service_requests): 状态机流转
- **任务表** (task_assignments): 认领和完成记录
- **通知表** (notifications): 多渠道通知记录
- **日志表** (operation_logs): 操作审计日志

### Consistency level

- **强一致性**: 用户认证、权限控制、任务状态流转
- **最终一致性**: 通知发送、统计报表

### Idempotency

- 需求提交：基于 `request_no` 去重
- 任务认领：状态检查 + 乐观锁

### Retry / compensation

- 通知发送：重试 3 次，指数退避
- 围栏匹配失败：使用默认站点兜底

### DTM usage

**No** - 小型项目，不需要分布式事务

---

## 8. Observability

### Metrics

- **业务指标**：需求数量、响应时间、完成率、认领率
- **系统指标**：API 调用量、错误率、活跃用户数

### Logs

- **操作日志** (operation_logs): 关键操作审计
- **错误日志**: Gin 框架错误中间件
- **访问日志**: Nginx access.log

### Traces

无分布式追踪（单体架构）

### Alerts (work-hours oriented)

- 需求分发失败（兜底机制触发）
- 通知发送失败率 > 10%
- API 错误率 > 5%
- 数据库连接异常

---

## 9. Deployment & Migration

### Migration steps (pre-deploy)

1. 执行数据库初始化脚本（`03-数据库设计.md` §5.1）
2. 导入测试数据（§5.2）
3. 配置环境变量（.env 文件）
4. 构建前端静态资源

### Deploy strategy

**Docker Compose 部署**：

```bash
docker-compose up -d
```

服务清单：
- frontend (Nginx 静态服务)
- backend (Go API 服务)
- mysql (MySQL 8.0)
- redis (Redis 7.0)
- nginx (反向代理)

### Rollback plan

1. 停止服务：`docker-compose down`
2. 回滚代码：`git checkout <previous-commit>`
3. 重新部署：`docker-compose up -d`
4. 数据库备份恢复（如有 schema 变更）

---

## 10. Risks & Edge Cases

### R1: 地理围栏匹配失败

**场景**: 用户位置不在任何围栏内
**处理**: 使用默认站点兜底

### R2: 并发认领冲突

**场景**: 多个服务人员同时认领同一任务
**处理**: 数据库唯一索引 + 乐观锁

### R3: 通知发送失败

**场景**: 短信/邮件服务不可用
**处理**: 异步队列 + 重试机制 + 失败记录

### R4: PWA 离线数据冲突

**场景**: 离线提交的需求与服务器数据冲突
**处理**: 时间戳比较 + 人工确认

### R5: 权限升级攻击

**场景**: 恶意用户尝试越权操作
**处理**: 自定义 RBAC + 路由守卫 + 接口级权限检查

### R6: 地理围栏边界情况

**场景**: 点恰好在多边形边界上
**处理**: 射线法算法明确定义（边界外）

---

## 11. Next Steps for `claude -p`

### 阶段1: 后端核心框架（1-2天）

- [ ] 创建 Domain 层实体定义
- [ ] 实现 Repository 层数据访问
- [ ] 开发用户认证和权限 API（JWT + 自定义 RBAC）
- [ ] 实现地理围栏匹配算法（射线法）

### 阶段2: 前端基础功能（2-3天）

- [ ] 管理门户：登录页面 + 权限系统
- [ ] 管理门户：布局组件 + 动态菜单
- [ ] 管理门户：任务池界面 + 认领功能
- [ ] C端：需求提交表单 + PWA 配置

### 阶段3: 核心业务功能（2-3天）

- [ ] 后端：需求分发引擎 + 任务管理 API
- [ ] 后端：通知服务（邮件/短信）
- [ ] 管理门户：围栏管理（地图绘制）
- [ ] 管理门户：统计分析报表

### 阶段4: 测试和部署（1-2天）

- [ ] 单元测试（地理围栏算法）
- [ ] 集成测试（API 接口）
- [ ] Docker 部署配置
- [ ] 性能测试和优化

---

**参考文档**:
- `/Users/yt/Documents/project/sCare/docs/general/01-PRD-产品需求文档.md`
- `/Users/yt/Documents/project/sCare/backend/docs/02-系统架构设计.md`
- `/Users/yt/Documents/project/sCare/backend/docs/03-数据库设计.md`
