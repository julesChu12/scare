# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

> 全局项目信息见根目录 `CLAUDE.md`，本文件仅补充后端专属规则。

## 技术栈

Go 1.25 · Gin · GORM (Gen) · MySQL 8.0 (utf8mb4) · Redis · JWT · Zap Logger · Cobra CLI · Swagger (swag)

## 常用命令

```bash
air                          # 热重载启动（推荐，端口 8080）
go run . serve               # 标准启动
go build ./...               # 编译检查
go test ./...                # 运行所有测试
go test ./pkg/geo/...        # 运行单个包测试
./scripts/test_api.sh        # API 回归测试

# GORM 模型重新生成
go run cmd/tools/gen/gorm_gen.go

# Swagger 文档更新（接口变更后必须执行）
swag init -g main.go -o docs --parseDependency --parseInternal
```

## 数据库连接

统一使用 UTF-8 字符集连接：

```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&collation=utf8mb4_unicode_ci",
    cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
```

## 用户角色获取模式

Service 层获取用户角色的标准方式（需兼容旧数据）：

```go
identities, _ := userIdentityRepo.GetActiveByUserID(userID)
```

## 新增功能开发流程

添加新的业务模块时，按以下顺序操作：

1. **Model**：如有新表，修改 `database/schema/` 后运行 `go run cmd/tools/gen/gorm_gen.go` 生成模型
2. **Repository**：在 `internal/repository/` 创建 `xxx_repo.go`
3. **Service**：在 `internal/service/` 创建 `xxx_service.go`
4. **Handler**：在 `internal/handler/` 创建 `xxx_handler.go`
5. **依赖注入**：在 `internal/router/deps.go` 的 `Deps` 结构体和 `NewDeps()` 中注册
6. **路由**：在 `internal/router/b_end.go` 或 `c_end.go` 中注册路由
7. **Swagger**：运行 `swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal`

## 注意事项

1. `.gen.go` 文件由 GORM Gen 生成，**禁止手动修改**
2. 接口变更后**必须**更新 Swagger 文档
3. 本地开发可用 `air` 热重载；若服务已在运行，优先复用已有进程
4. `scripts/test_api.sh` 是标准 API 回归测试脚本，提交前建议执行
5. Handler 层公共辅助：`context.go`（从 gin.Context 提取用户信息）、`response.go`（统一响应）、`pagination.go`（分页解析）、`params.go`（路径参数）
