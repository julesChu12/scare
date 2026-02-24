# 昌平区霍营街道社区养老信息分发平台

**版本**：v1.0.0
**状态**：开发中
**技术栈**：Go + Vue 3 + MySQL + Redis

---

## 项目简介

本项目是一个基于地理围栏的社区养老服务信息分发平台，采用 B/S 架构，提供统一的养老服务需求入口，通过地理围栏自动匹配站点，实现需求的智能分发与任务管理。

**核心特性**：
- 🎯 **地理围栏匹配**：基于点在多边形算法（射线法）实现需求自动分发
- 🔐 **权限管理**：基于自定义 RBAC 三表模型的细粒度权限控制
- 📱 **多端支持**：C端（用户端）+ 管理门户（工作人员+管理后台统一）
- 🚀 **PWA 支持**：C端支持 PWA，可离线使用
- 📊 **实时通知**：关键节点自动触发通知

---

## 技术架构

### 后端
- **语言**：Go 1.25
- **框架**：Gin（Web 框架）、GORM（ORM）
- **数据库**：MySQL 8.0+
- **缓存**：Redis 7.0+
- **权限**：自定义 RBAC（permissions/roles/role_permissions 三表）
- **认证**：JWT

### 前端
- **框架**：Vue 3 + TypeScript
- **构建**：Vite
- **UI 库**：Element Plus（C端 + 管理门户）
- **状态管理**：Pinia
- **路由**：Vue Router 4
- **图表**：ECharts + vue-echarts（管理门户）
- **权限**：自定义 RBAC + v-permission 指令（管理门户）
- **PWA**：Vite PWA（C端）

### 部署
- **容器化**：Docker + Docker Compose
- **反向代理**：Nginx
- **编排（可选）**：Kubernetes

---

## 快速开始

### 环境要求
- Go 1.25+
- Node.js 18+
- MySQL 8.0+
- Redis 7.0+
- npm 或 pnpm

### 1. 克隆项目
```bash
git clone https://github.com/your-org/community-elderly-care-platform.git
cd community-elderly-care-platform
```

### 2. 初始化数据库
```bash
cd database
./scripts/init.sh
```

### 3. 启动后端
```bash
cd backend
cp .env.example .env  # 修改配置
go mod download
go run cmd/server/main.go
```

### 4. 启动前端

#### C端（用户端）
```bash
cd frontend/c-end
npm install
npm run dev  # http://localhost:5174
```

#### 管理门户
```bash
cd frontend/management-portal
npm install
npm run dev  # http://localhost:3001
```

---

## 项目结构

### 整体结构

```
sCare/
├── backend/                    # 后端服务（Go 1.25）
│   ├── cmd/                    # 程序入口
│   │   └── server/             # API 服务启动入口
│   ├── internal/               # 私有应用代码
│   │   ├── config/             # 配置管理
│   │   ├── domain/             # 领域模型（Entity、VO）
│   │   ├── repository/         # 数据访问层
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # 中间件（认证、权限、CORS等）
│   │   └── algorithm/          # 算法实现（点在多边形）
│   ├── pkg/                    # 公共库
│   └── docs/                   # 后端文档
│       ├── 02-系统架构设计.md   # 架构设计
│       └── 03-数据库设计.md     # 数据库设计
│
├── frontend/                   # 前端项目
│   ├── c-end/                  # C端（用户端）
│   │   ├── docs/               # 前端文档
│   │   ├── public/             # PWA 静态资源
│   │   └── src/
│   │       ├── api/            # 接口层
│   │       ├── assets/         # 静态资源
│   │       ├── components/     # 组件
│   │       ├── composables/    # 组合式函数
│   │       ├── router/         # 路由
│   │       ├── store/          # Pinia
│   │       ├── utils/          # 工具
│   │       └── views/          # 页面
│   └── management-portal/      # 管理门户
│       ├── docs/
│       ├── public/
│       ├── package.json
│       ├── tsconfig.json
│       ├── vite.config.ts
│       └── src/
│       │       ├── api/
│       │       ├── assets/
│       │       ├── components/
│       │       ├── composables/
│       │       ├── config/
│       │       ├── directives/
│       │       ├── layouts/
│       │       ├── router/
│       │       │   ├── guards/
│       │       │   └── routes/
│       │       ├── store/
│       │       │   └── modules/
│       │       ├── types/
│       │       ├── utils/
│       │       └── views/
│       └── shared/
│           ├── api/
│           ├── constants/
│           ├── types/
│           └── utils/
│
├── docs/                       # 项目文档
│   ├── README.md              # 文档导航
│   ├── SPEC.md                # 项目规格文档
│   ├── PRD.md                 # 产品需求文档
│   ├── PROJECT_STATUS.md      # 项目开发状态
│   └── planning/              # 历史规划文档
│
├── database/                  # 数据库
│   ├── schema/                # 表结构
│   ├── seeds/                 # 测试数据
│   └── docs/                  # 数据库文档
│
├── deployment/                # 部署配置
│   └── docs/                  # 部署文档
│
└── docker-compose.yml         # Docker Compose配置
```

### 后端结构详解

```
backend/
├── cmd/server/
│   └── main.go                # 服务启动入口
│
├── internal/
│   ├── config/                 # 配置管理
│   │   └── config.go          # 配置结构定义
│   ├── domain/                # 领域模型
│   │   ├── models.go          # 数据库模型（GORM）
│   │   ├── roles.go           # 角色定义
│   │   ├── status.go          # 状态定义
│   │   └── service_types.go   # 服务类型
│   ├── repository/            # 数据访问层
│   │   ├── user_repo.go
│   │   ├── request_repo.go
│   │   ├── station_repo.go
│   │   ├── zone_repo.go
│   │   └── task_repo.go
│   ├── service/               # 业务逻辑层
│   │   ├── auth_service.go    # 认证服务
│   │   ├── request_service.go # 需求处理服务
│   │   ├── geofence_service.go # 围栏匹配服务
│   │   ├── task_service.go    # 任务管理服务
│   │   └── notification_service.go
│   ├── handler/               # HTTP 处理器
│   │   ├── b_auth_handler.go  # B端认证
│   │   ├── c_auth_handler.go  # C端认证
│   │   ├── request_handler.go
│   │   ├── zone_handler.go
│   │   └── task_handler.go
│   ├── middleware/            # 中间件
│   │   ├── auth.go            # JWT 认证
│   │   ├── permission.go       # RBAC 权限检查
│   │   └── end_type.go        # 端隔离
│   └── geofence/              # 地理围栏引擎
│       └── engine.go          # 射线法实现
│
├── pkg/                       # 公共库
│   ├── jwt/                   # JWT 工具
│   ├── logger/                # 日志工具
│   ├── response/              # 统一响应格式
│   └── errors/                # 错误处理
│
└── configs/                   # 配置文件
    └── config.yaml            # 应用配置
```

### 前端结构详解

```
frontend/
├── c-end/                     # C端（用户端）
│   ├── docs/
│   ├── public/
│   └── src/
│       ├── api/
│       ├── assets/
│       ├── components/
│       ├── composables/
│       ├── router/
│       ├── store/
│       ├── utils/
│       └── views/
└── management-portal/         # 管理门户
    ├── docs/
    ├── public/
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── pages/             # 页面组件
        ├── layouts/           # 布局组件
        ├── components/
        ├── router/
        │   └── guards/
        ├── store/
        │   └── modules/
        ├── config/
        ├── directives/
        ├── composables/
        ├── types/
        ├── utils/
        └── api/
```

---

## 核心功能

### 1. 需求提交（C端）
- 扫码进入，填写养老服务需求
- 自动获取地理位置
- 表单验证与提交

### 2. 地理围栏匹配
- 点在多边形算法（射线法）
- 多围栏优先级排序
- 兜底规则（未命中时分配到最近站点）

### 3. 任务管理（管理门户）
- 工作人员：任务池、认领、完成
- 站点负责人：任务分配、转派、人员管理
- 系统管理员：围栏管理、站点管理、用户管理、角色权限配置

### 4. 权限控制（自定义 RBAC）
- 基于 permissions/roles/role_permissions 三表实现
- 3种角色：工作人员（staff）、站点负责人（station_manager）、系统管理员（admin）
- Admin 角色跳过所有权限检查
- 路由级权限 + v-permission 指令按钮级权限

### 5. 通知推送
- 需求创建、任务认领、任务完成等关键节点触发通知
- 支持站内信、邮件（异步队列）

---

## 权限角色说明

| 角色 | 权限范围 | 典型操作 |
|------|---------|---------|
| **工作人员（staff）** | 本人任务 | 查看任务池、认领任务、完成任务、查看个人统计 |
| **站点负责人（station_manager）** | 本站点 + 继承 staff | 任务分配、转派、查看站点统计、管理站点人员 |
| **系统管理员（admin）** | 全局 + 继承 station_manager | 围栏管理、站点管理、用户管理、角色权限配置 |

---

## 开发指南

### 后端开发
```bash
cd backend

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 构建
go build -o bin/server cmd/server/main.go
```

### 前端开发

#### C端用户端
```bash
cd frontend/c-end
npm install
npm run dev      # http://localhost:5174
npm run build    # 生产构建
```

#### B端管理门户
```bash
cd frontend/management-portal
npm install
npm run dev      # http://localhost:3001
npm run build    # 生产构建
npm run lint     # 代码检查
```

---

## Docker 部署

### 开发环境
```bash
docker-compose up -d
```

### 生产环境
```bash
docker-compose -f deployment/docker-compose.prod.yml up -d
```

---

## API 文档

启动后端后访问：
- Swagger UI: http://localhost:8080/swagger/index.html
- API 文档: `docs/api/openapi.yaml`

---

## 测试

### 后端测试
```bash
cd backend
go test -v ./internal/algorithm  # 点在多边形算法测试
go test -v ./internal/service     # 业务逻辑测试
```

### 前端测试
```bash
# 目前暂未实现自动化测试
```

**测试文档**：
- 功能测试用例：`backend/docs/testing/功能测试用例表.md`
- 围栏匹配测试：`backend/docs/testing/围栏匹配测试记录.md`
- MVP测试报告：`backend/docs/testing/TEST_REPORT.md`

---

## 参考文档

- **📋 文档索引**：`docs/README.md`（文档导航）
- **📄 产品需求**：`docs/PRD.md`
- **📄 项目规格**：`docs/SPEC.md`（毕业设计规格文档）
- **📊 项目状态**：`docs/PROJECT_STATUS.md`
- **🔧 开发指南**：`backend/docs/01-开发指南.md`
- **🏗️ 系统架构**：`backend/docs/02-系统架构设计.md`
- **🗄️ 数据库设计**：`backend/docs/03-数据库设计.md`
- **📡 API接口**：`backend/docs/04-API接口设计.md`
- **🧪 测试报告**：`backend/docs/testing/TEST_REPORT.md`

---

## 许可证

MIT License

---

## 联系方式

- **作者**：[Your Name]
- **学校**：北京邮电大学
- **邮箱**：your-email@example.com

---

**备注**：本项目为毕业设计项目，用于展示基于地理围栏的需求分发系统设计与实现。
