# Seeds Layout

- `seed.sql`: 主入口（已合并模块内容，可直接 `mysql < seed.sql` 导入）
- `modules/00_reset_all.sql`: 全表清理
- `modules/10_roles_permissions.sql`: 角色与权限
- `modules/20_menus.sql`: 菜单
- `modules/30_stations_zones.sql`: 站点与围栏
- `modules/40_users_profiles.sql`: 用户、身份、客户档案
- `modules/50_requests_tasks.sql`: 服务请求、任务、历史
- `modules/60_content.sql`: 轮播图与新闻
- `modules/70_notifications.sql`: 通知

历史版本脚本已归档到 `legacy/`。
