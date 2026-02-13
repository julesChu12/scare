# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

社区养老信息分发平台后端服务，基于 Go + Gin + GORM + MySQL 构建。
Go 模块名：`community-elderly-care-platform`，Go 1.25+。

## 常用命令

```bash
air                          # 热重载启动（推荐，端口 8080）
go run . serve               # 标准启动（Cobra CLI）
go run . serve --port 8081   # 指定端口
go build ./...               # 编译检查
go test ./...                # 运行所有测试
go test ./internal/service/... -run TestAuthService  # 运行单个测试
go test ./pkg/geo/... -v     # 运行单个包测试（详细输出）

# GORM 模型重新生成（禁止手动修改 .gen.go 文件）
go run cmd/tools/gen/gorm_gen.go

# Swagger 文档更新（接口变更后必须执行）
swag init -g main.go -o docs --parseDependency --parseInternal

# API 回归测试
./scripts/test_api.sh

# 数据库迁移/种子数据（CLI 命令为 stub，需手动执行 SQL）
# database/migrations/*.sql  database/seeds/*.sql
```

## 分层架构

```
Handler (参数校验/响应格式化)
  → Service (业务逻辑/事务处理)
    → Repository (GORM 数据库操作)
      → Model (GORM Gen 生成，.gen.go 禁止手动修改)
```

依赖注入通过 `internal/router/deps.go` 的 `Deps` 容器统一管理，`NewDeps()` 按 Repo → Service → Handler 顺序初始化。新增模块时需在此文件注册。

扩展 GORM 模型：在 `internal/dao/model/` 中创建非 `.gen.go` 文件（如 `user_ext.go`）。

## 认证与权限体系

JWT 双端认证，token 中 `type` 字段区分 `b_end` / `c_end`。

中间件链（B 端受保护接口）：
```
AuthMiddleware → RequireEndType("b_end") → PermissionMiddleware
```

AuthMiddleware 解析 JWT 后向 Gin Context 注入：
- `user_id` (int64)、`user_type` (string)、`station_id` (int64)
- `user_identities` ([]string)、`user_primary` (string)
- `token_id` (string)、`token_expires_at` (time.Time)

Handler 层通过 `handler.GetUserID(c)`、`GetUserIdentities(c)` 等辅助函数读取。

权限系统已从 Casbin 迁移为自定义 `PermissionService`（内存缓存 + sync.RWMutex），Admin 角色跳过所有权限检查。Token 黑名单通过 Redis 实现，支持单 token 和用户级别撤销。

## API 约定

统一响应格式，使用 `handler.Respond()` / `handler.RespondPage()` / `handler.RespondError()`：
```json
{ "msg": "ok", "data": {...} }
{ "msg": "ok", "data": { "items": [], "page": 1, "page_size": 10, "total": 100 } }
{ "msg": "错误描述", "data": null }
```

分页参数：`page`（从 1 开始）、`page_size`（默认 10）。

## 路由结构

| 路径前缀 | 说明 | 认证 |
|----------|------|------|
| `/api/v1/b/auth/login` | B 端登录 | 无 |
| `/api/v1/b/auth/me`、`logout`、`menus/user` | B 端基础（无权限检查） | AuthMiddleware + RequireEndType |
| `/api/v1/b/*` | B 端业务接口 | + PermissionMiddleware |
| `/api/v1/c/auth/*` | C 端登录/注册 | 无 |
| `/api/v1/c/*` | C 端业务接口 | AuthMiddleware + RequireEndType("c_end") |
| `/swagger/*` | API 文档 | 无 |

## 角色体系

B 端：`admin` > `station_manager` > `staff`
C 端：`elderly`、`family`

多角色通过 `user_identities` 表支持，`primary` 字段标记主身份。获取角色的标准方式：
```go
identities := handler.GetUserIdentities(c) // 从 JWT Claims 中获取
```

## 地理围栏引擎（pkg/geo/）

服务请求创建时根据用户坐标自动匹配服务站点：BoundingBox 快速排除 → 射线法精确判断 → 按优先级返回首个匹配 StationID。围栏数据在 `GeofenceService.Reload()` 时加载到内存。

## 注意事项

1. 数据库使用逻辑外键，无数据库外键约束；软删除使用 `deleted_at`
2. MySQL 字符集统一 `utf8mb4`，DSN 需包含 `charset=utf8mb4&collation=utf8mb4_unicode_ci`
3. MySQL/Redis 运行在 Docker 容器中（`scare_mysql`、`scare_redis`）
4. 接口变更后必须更新 Swagger 文档
5. 本地开发用 `air` 启动；若服务已在运行，优先复用已有进程
6. `scripts/test_api.sh` 是标准 API 回归测试脚本，提交前建议执行
