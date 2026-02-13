-- =====================================================
-- 菜单数据种子文件
-- 版本: v1.0.0
-- 说明: B端管理后台菜单配置
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有菜单数据
TRUNCATE TABLE `menus`;

-- =====================================================
-- 顶级菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(1, 0, '工作台', '/dashboard', 'Dashboard', 'Odometer', 'dashboard', 1, 0, 'active'),
(2, 0, '服务管理', '/services', '', 'Service', 'service', 2, 0, 'active'),
(3, 0, '站点管理', '/stations', '', 'OfficeBuilding', 'station', 3, 0, 'active'),
(4, 0, '数据中心', '/data', '', 'DataAnalysis', 'data', 4, 0, 'active'),
(5, 0, '内容管理', '/content', '', 'Document', 'content', 5, 0, 'active'),
(6, 0, '通知中心', '/notifications', 'NotificationCenter', 'Bell', 'public:notification', 6, 0, 'active'),
(7, 0, '系统管理', '/system', '', 'Setting', 'system', 7, 0, 'active');

-- =====================================================
-- 服务管理子菜单 (parent_id=2)
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(8, 2, '服务请求', '/services/requests', 'ServiceRequests', '', 'service:request:list', 1, 0, 'active'),
(9, 2, '任务管理', '/services/tasks', 'TaskManagement', '', 'service:task:pool', 2, 0, 'active');

-- =====================================================
-- 站点管理子菜单 (parent_id=3)
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(10, 3, '站点列表', '/stations/list', 'StationList', '', 'station:list:view', 1, 0, 'active'),
(11, 3, '服务围栏', '/stations/zones', 'ServiceZones', '', 'station:zone:list', 2, 0, 'active');

-- =====================================================
-- 数据中心子菜单 (parent_id=4)
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(12, 4, '统计分析', '/data/statistics', 'StatisticsView', '', 'data:statistics:view', 1, 0, 'active'),
(13, 4, '报表管理', '/data/reports', 'ReportManagement', '', 'data:report:list', 2, 0, 'active');

-- =====================================================
-- 内容管理子菜单 (parent_id=5)
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(14, 5, '轮播图管理', '/content/banners', 'BannerManagement', '', 'content:banner:list', 1, 0, 'active'),
(15, 5, '新闻管理', '/content/news', 'NewsManagement', '', 'content:news:list', 2, 0, 'active');

-- =====================================================
-- 系统管理子菜单 (parent_id=7)
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(16, 7, '用户管理', '/system/users', 'UserManagement', '', 'system:user:list', 1, 0, 'active'),
(17, 7, '角色管理', '/system/roles', 'RoleManagement', '', 'system:role:list', 2, 0, 'active'),
(18, 7, '菜单管理', '/system/menus', 'MenuManagement', '', 'system:menu:list', 3, 0, 'active');

-- 验证数据
SELECT '菜单数据' AS 'Table', COUNT(*) AS 'Count' FROM `menus` WHERE `deleted_at` IS NULL;
