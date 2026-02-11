# sCare Backend 项目规则

## 项目概述

社区养老信息分发平台后端服务，基于 Go + Gin + GORM + MySQL 构建。

## 技术栈

| 技术 | 版本/说明 |
|------|----------|
| Go | 1.25+ |
| Gin | Web Framework |
| GORM | ORM + Gen 代码生成 |
| MySQL | 8.0+ (utf8mb4) |
| Redis | 缓存/会话 |
| Casbin | RBAC 权限控制 |
| Air | 热重载开发 |

## 开发环境

### 数据库连接配置

```env
# MySQL
DB_HOST=localhost
DB_PORT=3306
DB_NAME=scare_db
DB_USER=scare_user
DB_PASSWORD=scare_pass
DB_CHARSET=utf8mb4
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Docker 服务

MySQL 和 Redis 部署在 Docker 中，容器名称：`scare_mysql`、`scare_redis`

```bash
# 连接 MySQL
docker exec -it scare_mysql mysql -u scare_user -pscare_pass scare_db

# 连接 Redis
docker exec -it scare_redis redis-cli

# 查看容器状态
docker ps | grep scare
```

### 测试账号

| 角色 | 手机号 | 密码 |
|------|--------|------|
| Admin | 13800000001 | Test@123 |
| Station Manager | 13800000002 | Test@123 |
| Staff | 13800000004 | Test@123 |

## 开发规范

### 分层架构

```
Handler (API层) → Service (业务层) → Repository (数据层) → Model (模型层)
```

| 层级 | 职责 |
|------|------|
| Handler | 参数校验、响应格式化，不含业务逻辑 |
| Service | 业务逻辑、数据组装、事务处理 |
| Repository | 数据库操作，使用 GORM Gen 查询 |
| Model | 数据模型定义（由 Gen 生成） |

### GORM Gen 模型

- 模型文件：`internal/dao/model/*.gen.go`
- **禁止手动修改 `.gen.go` 文件**
- 扩展模型：在同目录创建非 gen 文件（如 `user_ext.go`）
- 重新生成：`go run cmd/tools/gen/gorm_gen.go`

### 数据库连接
**重要**: 统一使用 UTF-8 字符集（MySQL 侧使用 `utf8mb4`）进行连接与写入：

```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&collation=utf8mb4_unicode_ci",
    cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
```

### API 响应格式

```json
// 成功
{ "msg": "ok", "data": { ... } }

// 错误
{ "msg": "错误描述", "data": null }
```

### 分页与筛选参数

| 参数 | 说明 |
|------|------|
| `page` | 页码，从 1 开始 |
| `page_size` | 每页数量，默认 10 |
| `role` | 角色筛选 |
| `status` | 状态筛选 |
| `station_id` | 站点筛选 |
| `keyword` | 关键词搜索 |

### 用户角色系统

| 角色代码 | 说明 |
|----------|------|
| `admin` | 系统管理员（全部权限） |
| `station_manager` | 站点管理员 |
| `staff` | 工作人员 |
| `elderly` | 老年人（C端） |
| `family` | 家属（C端） |

```go
// 获取用户角色的标准方式（Service 层）
roles, _ := userRoleRepo.GetActiveByUserID(userID)
// 兼容旧数据
if len(roles) == 0 && user.Role != "" {
    roles = []string{user.Role}
}
```

## 常用命令

```bash
# 启动开发服务器（热重载，推荐）
air

# 启动后端服务（不使用热重载）
go run . serve

# 编译检查
go build ./...

# 运行测试
go test ./...

# API 回归测试脚本
./scripts/test_api.sh

# 重新生成 GORM 模型
go run cmd/tools/gen/gorm_gen.go

# 重新生成 Swagger 文档
swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
```

## 目录结构

```
backend/
├── cmd/
│   ├── server/main.go      # 主入口
│   └── tools/gen/          # 代码生成工具
├── configs/                # 配置文件
├── internal/
│   ├── config/             # 配置加载
│   ├── consts/             # 常量定义
│   ├── dao/
│   │   ├── model/          # GORM Gen 生成的模型
│   │   └── query/          # GORM Gen 生成的查询
│   ├── handler/            # API 处理器
│   ├── middleware/         # 中间件
│   ├── repository/         # 数据访问层
│   ├── router/             # 路由配置
│   └── service/            # 业务逻辑层
├── pkg/                    # 公共包
├── scripts/                # 脚本
└── docs/                   # Swagger 文档
```

## 注意事项

1. 所有数据库表使用逻辑外键，不使用数据库外键约束
2. 软删除使用 `deleted_at` 字段
3. Admin 角色拥有所有权限，在代码中特殊处理
4. B端和C端使用不同的 JWT 类型（`b_end` / `c_end`）
5. **接口变更后必须更新 Swagger 文档**：`swag init -g main.go -o docs --parseDependency --parseInternal`
6. **本地开发可使用 `air` 启动服务；若服务已在运行，优先复用已有进程**
7. **MySQL 运行在 Docker 容器中（默认容器名：`scare_mysql`）**
8. **`scripts/test_api.sh` 是标准 API 回归测试脚本，提交前建议至少执行一轮**
