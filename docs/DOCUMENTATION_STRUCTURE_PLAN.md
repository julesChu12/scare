# 文档结构重组方案

**日期**：2026-01-31  
**目标**：按照子模块分类组织文档，提高可维护性

---

## 🎯 重组原则

1. **模块化**：按照代码模块对应文档位置
2. **就近原则**：文档放在对应模块目录下
3. **清晰层次**：根目录 → 模块目录 → 文档文件

---

## 📁 新文档结构

```
sCare/
├── README.md                              # ✅ 项目主文档（保留）
├── .gitignore                             # ✅ Git配置（保留）
│
├── docs/                                  # 📁 项目级通用文档
│   ├── README.md                          # 📄 文档索引（新建）
│   ├── PROJECT_STATUS.md                  # 📄 项目状态报告（保留）
│   ├── SPEC.md                            # 📄 项目规格文档（保留）
│   └── planning/                          # 📁 历史规划文档（归档）
│       ├── API_PLANNING.md
│       ├── MULTI_ROLE_REFACTORING.md
│       ├── ROUTES_AND_PERMISSIONS.md
│       ├── SPEC_MULTI_ROLE.md
│       └── 02-前端进度对照表.md
│
├── backend/                               # 📁 后端模块
│   └── docs/
│       ├── README.md                      # 📄 后端文档索引（新建）
│       ├── 01-开发指南.md                 # ➡️ 从 docs/DEVELOPMENT.md 移动
│       ├── 02-系统架构设计.md             # ✅ 已存在
│       ├── 03-数据库设计.md               # ✅ 已存在
│       ├── 04-API接口设计.md              # ✅ 已存在
│       ├── 05-配置说明.md                 # ✅ 已存在
│       ├── 07-部署方案.md                 # ✅ 已存在
│       ├── MVP_API_SPEC.md                # ✅ 已存在
│       └── testing/                       # 📁 测试文档
│           ├── TEST_REPORT.md             # ➡️ 从 docs/TEST_REPORT.md 移动
│           ├── 功能测试用例表.md          # ✅ 已存在
│           └── 围栏匹配测试记录.md        # ✅ 已存在
│
├── frontend/                              # 📁 前端模块
│   ├── c-end/
│   │   └── docs/
│   │       ├── README.md                  # 📄 C端文档索引（已存在）
│   │       ├── 05-C端前端设计.md          # ✅ 已存在
│   │       ├── c-end-design-prompt.md     # ✅ 已存在
│   │       ├── c-end-feature-specs.md     # ✅ 已存在
│   │       ├── c-end-ui-wireframes.md     # ✅ 已存在
│   │       └── SPEC.md                    # ✅ 已存在
│   │
│   └── management-portal/
│       └── docs/
│           ├── README.md                  # 📄 B端文档索引（已存在）
│           ├── 03-前端架构决策-权限集成方案.md  # ✅ 已存在
│           ├── 06-B端前端设计.md          # ✅ 已存在
│           └── SPEC.md                    # ✅ 已存在
│
├── database/                              # 📁 数据库模块
│   └── docs/                              # 📁 数据库文档（新建目录）
│       └── README.md                      # 📄 数据库文档索引（新建）
│
└── deployment/                            # 📁 部署模块
    └── docs/                              # 📁 部署文档（新建目录）
        └── README.md                      # 📄 部署文档索引（新建）
```

---

## 🔄 文档迁移清单

### 1. 移动到后端文档

| 源文件 | 目标位置 | 原因 |
|--------|----------|------|
| `docs/DEVELOPMENT.md` | `backend/docs/01-开发指南.md` | 主要是后端开发环境配置 |
| `docs/TEST_REPORT.md` | `backend/docs/testing/TEST_REPORT.md` | 后端API测试报告 |

### 2. 保留在项目根文档

| 文件 | 位置 | 原因 |
|------|------|------|
| `SPEC.md` | `docs/SPEC.md` | 项目级规格文档 |
| `PROJECT_STATUS.md` | `docs/PROJECT_STATUS.md` | 项目整体状态 |

### 3. 归档历史规划文档

| 文件 | 目标位置 | 原因 |
|------|----------|------|
| `docs/API_PLANNING.md` | `docs/planning/API_PLANNING.md` | 历史规划 |
| `docs/MULTI_ROLE_REFACTORING.md` | `docs/planning/MULTI_ROLE_REFACTORING.md` | 历史规划 |
| `docs/ROUTES_AND_PERMISSIONS.md` | `docs/planning/ROUTES_AND_PERMISSIONS.md` | 历史规划 |
| `docs/SPEC_MULTI_ROLE.md` | `docs/planning/SPEC_MULTI_ROLE.md` | 历史规划 |
| `docs/general/02-前端进度对照表.md` | `docs/planning/02-前端进度对照表.md` | 历史规划 |

### 4. 合并到项目需求文档

| 文件 | 目标位置 | 原因 |
|------|------|------|
| `docs/general/01-PRD-产品需求文档.md` | `docs/PRD.md` | 产品需求文档应在项目根文档 |

### 5. 删除临时文档

| 文件 | 原因 |
|------|------|
| `docs/DOCUMENTATION_AUDIT.md` | 临时审核报告，已完成可删除 |
| `docs/DOCUMENTATION_REORGANIZATION_PLAN.md` | 临时整理计划，已完成可删除 |
| `docs/general/README.md` | general目录将被移除 |

---

## 📝 新建文档索引

### `docs/README.md`

```markdown
# sCare 项目文档

## 📋 核心文档

- [SPEC.md](./SPEC.md) - 项目规格文档
- [PRD.md](./PRD.md) - 产品需求文档
- [PROJECT_STATUS.md](./PROJECT_STATUS.md) - 项目开发状态

## 📁 模块文档

- [后端文档](../backend/docs/) - 后端开发、API、数据库
- [C端前端文档](../frontend/c-end/docs/) - C端用户界面
- [B端前端文档](../frontend/management-portal/docs/) - 管理门户
- [数据库文档](../database/docs/) - 数据库设计与脚本
- [部署文档](../deployment/docs/) - 部署配置与指南

## 📚 历史规划文档

- [planning/](./planning/) - 历史规划文档归档
```

### `backend/docs/README.md`

```markdown
# 后端文档

## 📖 开发文档

- [01-开发指南.md](./01-开发指南.md) - 开发环境配置
- [02-系统架构设计.md](./02-系统架构设计.md) - 系统架构设计
- [03-数据库设计.md](./03-数据库设计.md) - 数据库设计
- [04-API接口设计.md](./04-API接口设计.md) - API接口设计
- [05-配置说明.md](./05-配置说明.md) - 配置说明
- [07-部署方案.md](./07-部署方案.md) - 部署方案
- [MVP_API_SPEC.md](./MVP_API_SPEC.md) - MVP API实现规格

## 🧪 测试文档

- [testing/TEST_REPORT.md](./testing/TEST_REPORT.md) - 测试报告
- [testing/功能测试用例表.md](./testing/功能测试用例表.md) - 功能测试用例
- [testing/围栏匹配测试记录.md](./testing/围栏匹配测试记录.md) - 围栏匹配测试
```

### `database/docs/README.md`

```markdown
# 数据库文档

## 📊 数据库设计

参考：[backend/docs/03-数据库设计.md](../../backend/docs/03-数据库设计.md)

## 📁 SQL 脚本

- [schema/schema.sql](../schema/schema.sql) - 表结构定义
- [seeds/seed.sql](../seeds/seed.sql) - 测试数据
```

### `deployment/docs/README.md`

```markdown
# 部署文档

## 🚀 部署方案

参考：[backend/docs/07-部署方案.md](../../backend/docs/07-部署方案.md)

## 📦 Docker 配置

- [docker-compose.yml](../../docker-compose.yml) - Docker Compose配置
```

---

## ✅ 执行步骤

1. 创建新目录
2. 移动文档文件
3. 创建文档索引
4. 更新README.md引用
5. 删除临时文档
6. 验证所有链接

---

**创建时间**：2026-01-31
