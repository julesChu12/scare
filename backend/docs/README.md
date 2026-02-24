# 后端文档

**模块**：sCare Backend  
**语言**：Go 1.25  
**框架**：Gin + GORM  

---

## 📖 开发文档

| 文档 | 说明 |
|------|------|
| [01-开发指南.md](./01-开发指南.md) | 开发环境配置与快速开始 |
| [02-系统架构设计.md](./02-系统架构设计.md) | 系统分层架构与模块设计 |
| [03-数据库设计.md](./03-数据库设计.md) | 数据库表结构设计（Code First） |
| [04-API接口设计.md](./04-API接口设计.md) | RESTful API 设计规范 |
| [05-配置说明.md](./05-配置说明.md) | 配置文件说明 |
| [07-部署方案.md](./07-部署方案.md) | 生产环境部署方案 |
| [08-快速开通功能实现.md](./08-快速开通功能实现.md) | 快速开通系统完整实现 |
| [MVP_API_SPEC.md](./MVP_API_SPEC.md) | MVP API 实现规格 |
| [banner_api.md](./banner_api.md) | Banner 轮播图 API 文档 |

---

## 🔗 前端接入指南

| 文档 | 说明 |
|------|------|
| [C_END_FRONTEND_GUIDE.md](./C_END_FRONTEND_GUIDE.md) | C 端前端接入指南 |
| [PERMISSION_FRONTEND_GUIDE.md](./PERMISSION_FRONTEND_GUIDE.md) | 权限系统前端接入指南 |

---

## 🧪 测试文档

| 文档 | 说明 |
|------|------|
| [testing/TEST_REPORT.md](./testing/TEST_REPORT.md) | MVP回归测试报告 |
| [testing/测试覆盖率报告.md](./testing/测试覆盖率报告.md) | 详细测试覆盖率分析 |
| [testing/功能测试用例表.md](./testing/功能测试用例表.md) | 功能测试用例 |
| [testing/围栏匹配测试记录.md](./testing/围栏匹配测试记录.md) | 地理围栏匹配测试 |

---

## 🎯 快速开始

```bash
# 1. 启动依赖服务
docker-compose up -d

# 2. 安装依赖
cd backend
go mod download

# 3. 配置环境变量
cp .env.example .env

# 4. 启动后端（热加载）
air
```

详细说明：[01-开发指南.md](./01-开发指南.md)

---

## 📊 项目结构

```
backend/
├── cmd/                  # CLI 命令（Cobra）
├── internal/
│   ├── config/          # 配置管理
│   ├── consts/          # 常量定义
│   ├── dao/             # GORM Gen 生成的模型与查询
│   ├── dto/             # 数据传输对象
│   ├── handler/         # HTTP 处理器
│   ├── service/         # 业务逻辑
│   ├── repository/      # 数据访问层
│   ├── middleware/      # 中间件（JWT、权限等）
│   ├── router/          # 路由注册与依赖注入
│   ├── notify/          # 通知服务（邮件）
│   └── storage/         # 文件存储（本地/OSS）
├── pkg/                 # 公共库（geo、crypto、jwt、logger）
├── database/            # 数据库 Schema、迁移、种子数据
├── scripts/             # 工具脚本
└── docs/                # 项目文档
```

---

**最后更新**：2026-02-24
