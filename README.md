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
- Docker / Docker Compose
- npm 或 pnpm

### 配置与依赖说明

| 模块 | 配置文件 | 是否必须 | 说明 |
|------|----------|----------|------|
| 后端 | `backend/.env` | 必须 | MySQL、Redis、JWT、存储、安全加密、高德 Key 等配置 |
| 管理门户 | `frontend/management-portal/.env` | 推荐 | `VITE_API_BASE_URL`、地图 Key、C 端跳转地址等 |
| C端 | 无必须文件 | 否 | 本地开发默认走 Vite 代理到 `http://localhost:8080`，可选 `VITE_MOCK_USER` |

本地默认端口：
- 后端 API：`8080`
- MySQL：`3306`
- Redis：`6379`
- 管理门户：`3001`
- C端：`5174`（HTTPS）

### 1. 克隆项目
```bash
git clone https://github.com/your-org/community-elderly-care-platform.git
cd community-elderly-care-platform
```

### 2. 启动 MySQL / Redis 并初始化数据库
```bash
cd backend
cp .env.example .env
docker compose up -d
go run . migrate
```

说明：
- `backend/docker-compose.yml` 会启动本地 MySQL 8.0 和 Redis 7.0
- MySQL 首次启动会自动执行 `backend/database/schema/schema.sql` 和 `backend/database/seeds/001~007`
- 增量结构变更统一通过 `go run . migrate` 执行，迁移记录保存在 `schema_migrations`
- 如需重建本地数据库，可执行 `docker compose down -v && docker compose up -d && go run . migrate`
- 当前默认种子数据会初始化 `12` 个用户、`4` 个站点、`6` 个围栏、`4` 条服务请求和 `18` 个菜单

### 3. 启动后端
```bash
cd backend
go mod download
go run . serve
```

启动成功后可访问：
- 健康检查：`http://localhost:8080/api/v1/health`
- Swagger：`http://localhost:8080/swagger/index.html`

### 4. 启动前端

#### C端（用户端）
```bash
cd frontend/c-end
npm install
npm run dev  # https://localhost:5174
```

说明：
- C端本地开发启用了 `basicSsl`，浏览器首次访问会出现本地证书提示
- C端接口默认走 `/api` 代理到 `http://localhost:8080`

#### 管理门户
```bash
cd frontend/management-portal
cp .env.development.example .env
npm install
npm run dev  # http://localhost:3001
```

说明：
- 管理门户通过 `VITE_API_BASE_URL` 指向后端，默认值为 `http://localhost:8080/api/v1`
- 涉及地图选点/围栏编辑页面时，需要配置 `VITE_AMAP_KEY` 与 `VITE_AMAP_SECURITY_JS_CODE`

### 5. 验证默认测试数据

数据库自动初始化后，默认会写入测试账号和基础业务数据。

示例管理员账号：
- 手机号：`13800000001`
- 密码：`Test@123`

示例状态检查：

```bash
curl http://localhost:8080/api/v1/health
curl -I http://localhost:3001
curl -k -I https://localhost:5174
```

---

## 项目结构

### 整体结构

```
sCare/
├── backend/                    # 后端服务（Go 1.25）
│   ├── main.go                 # CLI 入口（如 serve / migrate）
│   ├── cmd/                    # 子命令实现
│   ├── internal/               # 私有应用代码
│   │   ├── config/             # 配置管理
│   │   ├── consts/             # 常量定义
│   │   ├── dao/                # GORM Gen 生成模型与查询
│   │   ├── dto/                # 数据传输对象
│   │   ├── repository/         # 数据访问层
│   │   ├── service/            # 业务逻辑层
│   │   ├── handler/            # HTTP 处理器
│   │   ├── middleware/         # 中间件（认证、权限、端隔离等）
│   │   ├── router/             # 路由注册
│   │   ├── notify/             # 通知能力
│   │   └── storage/            # 文件存储封装
│   ├── pkg/                    # 公共库
│   │   ├── database/           # MySQL 连接
│   │   ├── geo/                # 地理围栏算法
│   │   ├── jwt/                # JWT 工具
│   │   ├── logger/             # 日志封装
│   │   └── redis/              # Redis 连接
│   ├── database/               # schema / seed / migration
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
├── database/                  # 历史兼容目录（仅保留说明性文件）
│
├── deployment/                # 部署配置
│   └── docs/                  # 部署文档
│
└── backend/docker-compose.yml # 本地 MySQL / Redis 依赖
```

### 后端结构详解

```
backend/
├── main.go                    # CLI 入口
├── cmd/                       # serve / migrate 等命令
├── database/                  # schema / seeds / migrations（数据库唯一真相源）
├── internal/
│   ├── config/                # 配置管理
│   ├── consts/                # 业务常量
│   ├── dao/                   # 生成模型与查询
│   ├── dto/                   # DTO 定义
│   ├── repository/            # 数据访问层
│   ├── service/               # 业务逻辑层
│   ├── handler/               # HTTP 处理器
│   ├── middleware/            # 认证、权限、端隔离
│   ├── router/                # 路由装配
│   ├── notify/                # 邮件/通知
│   └── storage/               # 文件存储
├── pkg/                       # 公共库
│   ├── database/              # GORM/MySQL 连接
│   ├── geo/                   # 围栏匹配算法
│   ├── jwt/                   # JWT 工具
│   ├── logger/                # 日志工具
│   └── redis/                 # Redis 客户端
├── docs/                      # 后端文档
├── tests/                     # API 测试脚本
└── storage/                   # 本地上传文件
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
go build -o scare .

# 标准启动
go run . serve

# 热重载（已安装 air 时）
air -c .air.toml
```

### 前端开发

#### C端用户端
```bash
cd frontend/c-end
npm install
npm run dev      # https://localhost:5174
npm run build    # 生产构建
```

#### B端管理门户
```bash
cd frontend/management-portal
cp .env.development.example .env
npm install
npm run dev      # http://localhost:3001
npm run build    # 生产构建
npm run lint     # 代码检查
```

---

## Docker 部署

### 开发环境
```bash
cd backend
docker compose up -d
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
go test -v ./pkg/geo              # 点在多边形算法测试
go test -v ./internal/service/... # 业务逻辑测试
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
