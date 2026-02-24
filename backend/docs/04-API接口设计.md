# API接口设计

**版本**：v1.1
**日期**：2026年1月16日
**注意**：**本文档仅作为高层设计参考。详细字段定义、请求示例和响应结构请查阅 Swagger 文档 (由代码自动生成)。**

---

## 1. 核心设计原则

### 1.1 接口风格
*   严格遵循 RESTful 规范。
*   所有 API 均以 `/api` 开头。
*   版本控制通过 URL 路径实现（如 `/api/v1/...`）。

### 1.2 认证与鉴权
*   **认证**: HTTP Header `Authorization: Bearer <JWT_TOKEN>`。
*   **鉴权**: 基于自定义 RBAC 三表模型（permissions/roles/role_permissions）。

### 1.3 状态码规范
*   `200`: 成功。
*   `400`: 参数错误（应用层校验失败）。
*   `401`: 未认证。
*   `403`: 权限不足。
*   `500`: 服务器内部错误。

---

## 2. 核心业务流程与状态机

Swagger 无法很好地展示业务状态流转，此处重点说明核心业务对象的状态机。

### 2.1 需求 (Request) 状态流转

| 状态码 (const) | 说明 | 前置状态 | 触发动作 |
|---------------|------|----------|---------|
| `pending` | 待分发 | (初始) | 用户提交需求 |
| `dispatched` | 已分发 | pending | 系统自动/人工分发 |
| `claimed` | 已认领 | dispatched | 服务人员抢单 |
| `processing` | 处理中 | claimed | 服务人员开始服务 |
| `completed` | 已完成 | processing | 服务人员确认完成 |
| `cancelled` | 已取消 | any | 用户/管理员取消 |

### 2.2 任务 (Task) 状态流转

| 状态码 (const) | 说明 | 对应需求状态 |
|---------------|------|-------------|
| `claimed` | 已认领 | claimed |
| `processing` | 处理中 | processing |
| `completed` | 已完成 | completed |
| `cancelled` | 已取消 | cancelled |

---

## 3. 核心资源定义 (High Level)

### 3.1 需求 (Requests)
*   `POST /api/requests`: 提交需求（含 lat/lng）。
*   `GET /api/requests`: 查询列表。
*   `POST /api/requests/:id/cancel`: 取消。

### 3.2 任务 (Tasks)
*   `GET /api/tasks/pool`: 获取本站点的待认领任务池。
*   `POST /api/tasks/:id/claim`: 认领任务。
*   `POST /api/tasks/:id/complete`: 完成任务（上传照片）。

### 3.3 围栏 (Zones) - 管理员
*   `POST /api/zones`: 创建围栏。
    *   **输入**: 围栏名称、站点ID、**顶点坐标数组 (JSON)**。
*   `GET /api/zones`: 获取围栏列表。

---

## 4. Swagger 使用指南

在开发模式下，启动服务器后访问：
`http://localhost:8080/swagger/index.html`

**如何生成/更新 Swagger 文档**:
```bash
# 在项目根目录下执行
swag init -g cmd/server/main.go -o docs/swagger
```