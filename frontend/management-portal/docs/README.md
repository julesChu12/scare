# 管理后台前端文档

统一管理门户（B端工作人员 + 系统管理员）前端文档索引。

## 核心文档

| 文档 | 说明 |
|------|------|
| [06-B端前端设计.md](./06-B端前端设计.md) | 管理后台界面设计和交互流程 |

## 技术栈

- Vue 3 + TypeScript + Vite + Pinia
- Element Plus UI 组件库
- 自定义 RBAC 权限（v-permission 指令 + usePermission composable）
- ECharts + vue-echarts 数据可视化
- 高德地图 JS API（围栏编辑）

## 已完成功能（18个页面）

1. **认证**: 登录页、JWT Token 管理
2. **Dashboard**: 工作台统计概览
3. **站点管理**: 站点 CRUD、围栏管理（地图编辑）
4. **用户管理**: 用户列表、角色权限配置
5. **业务管理**: 服务需求、任务管理、服务类型
6. **内容管理**: 新闻公告、轮播图管理
7. **系统管理**: 菜单管理、日志管理、通知管理
8. **其他**: 工单统计、个人中心、老年人档案

## 权限体系（四层）

1. 路由守卫（`meta.permission_code`）
2. v-permission 指令（DOM 级别）
3. usePermission composable
4. 侧边栏菜单动态获取（`/b/menus/user`）

## 相关文档

- **后端 API**: `backend/docs/04-API接口设计.md`
- **数据库**: `backend/docs/03-数据库设计.md`
- **项目总览**: `README.md`

---

**最后更新**: 2026-02-24
