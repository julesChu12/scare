-- =====================================================
-- 种子数据主入口文件
-- 版本: v1.0.0
-- 说明: 按顺序执行所有种子文件
-- =====================================================
--
-- 使用方法:
-- docker exec -i scare_mysql mysql -u scare_user -pscare_pass scare_db < seed_all.sql
--
-- 或者单独执行某个模块:
-- docker exec -i scare_mysql mysql -u scare_user -pscare_pass scare_db < 001_seed_permissions.sql
--
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- =====================================================
-- 执行顺序（有依赖关系）
-- =====================================================
-- 1. 权限系统（基础数据，无依赖）
-- 2. 用户数据（依赖权限系统的角色）
-- 3. 站点和围栏（依赖用户）
-- 4. 服务需求和任务（依赖用户、站点）
-- 5. 通知（依赖用户、任务）
-- 6. 新闻（依赖用户）
-- 7. 菜单（独立）
-- =====================================================

SOURCE 001_seed_permissions.sql;
SOURCE 002_seed_users.sql;
SOURCE 003_seed_stations.sql;
SOURCE 004_seed_requests.sql;
SOURCE 005_seed_notifications.sql;
SOURCE 006_seed_news.sql;
SOURCE 007_seed_menus.sql;

-- =====================================================
-- 数据统计
-- =====================================================
SELECT '========== 种子数据统计 ==========' AS '';
SELECT '权限' AS 'Table', COUNT(*) AS 'Count' FROM `permissions`
UNION ALL SELECT '角色', COUNT(*) FROM `roles`
UNION ALL SELECT '角色权限', COUNT(*) FROM `role_permissions`
UNION ALL SELECT '用户', COUNT(*) FROM `users`
UNION ALL SELECT '老年人档案', COUNT(*) FROM `elderly_profiles`
UNION ALL SELECT '用户角色', COUNT(*) FROM `user_roles`
UNION ALL SELECT '服务站点', COUNT(*) FROM `service_stations`
UNION ALL SELECT '服务围栏', COUNT(*) FROM `service_zones`
UNION ALL SELECT '服务需求', COUNT(*) FROM `service_requests`
UNION ALL SELECT '任务分配', COUNT(*) FROM `task_assignments`
UNION ALL SELECT '任务历史', COUNT(*) FROM `task_histories`
UNION ALL SELECT '通知记录', COUNT(*) FROM `notifications`
UNION ALL SELECT '新闻', COUNT(*) FROM `news` WHERE `deleted_at` IS NULL
UNION ALL SELECT '菜单', COUNT(*) FROM `menus` WHERE `deleted_at` IS NULL;

-- =====================================================
-- 测试账号说明
-- =====================================================
-- 管理员: 13800000001 / Test@123
-- 霍营站站长: 13800000002 / Test@123
-- 龙泽站站长: 13800000003 / Test@123
-- 霍营站工作人员: 13800000004, 13800000005 / Test@123
-- 龙泽站工作人员: 13800000006, 13800000007 / Test@123
-- 老年人: 13800000008, 13800000009, 13800000010 / Test@123
-- 家属: 13800000011, 13800000012 / Test@123
-- =====================================================
