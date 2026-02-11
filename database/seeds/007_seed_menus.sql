-- =====================================================
-- 菜单数据种子文件
-- 版本: v1.0.0
-- 说明: B端管理后台菜单配置
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有菜单数据
DELETE FROM `menus` WHERE `deleted_at` IS NULL;

-- =====================================================
-- 顶级菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(1, 0, '工作台', '/dashboard', 'Dashboard', 'Odometer', 'dashboard', 1, 0, 'active'),
(2, 0, '服务管理', '/services', '', 'Service', 'service', 2, 0, 'active'),
(3, 0, '站点管理', '/stations', '', 'OfficeBuilding', 'station', 3, 0, 'active'),
(4, 0, '内容管理', '/content', '', 'Document', 'content', 4, 0, 'active'),
(5, 0, '系统管理', '/system', '', 'Setting', 'system', 5, 0, 'active');

-- =====================================================
-- 服务管理子菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(6, 2, '任务管理', '/services/tasks', 'TaskManagement', '', 'service:task:pool', 1, 0, 'active'),
(7, 2, '服务请求', '/services/requests', 'ServiceRequests', '', 'service:request:list', 2, 0, 'active');

-- =====================================================
-- 站点管理子菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(8, 3, '站点列表', '/stations/list', 'StationList', '', 'station:list:view', 1, 0, 'active'),
(9, 3, '服务围栏', '/stations/zones', 'ServiceZones', '', 'station:zone:list', 2, 0, 'active');

-- =====================================================
-- 内容管理子菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(10, 4, '轮播图管理', '/content/banners', 'BannerManagement', '', 'content:banner:list', 1, 0, 'active'),
(11, 4, '新闻管理', '/content/news', 'NewsManagement', '', 'content:news:list', 2, 0, 'active');

-- =====================================================
-- 系统管理子菜单
-- =====================================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
(12, 5, '用户管理', '/system/users', 'UserManagement', '', 'system:user:list', 1, 0, 'active'),
(13, 5, '角色管理', '/system/roles', 'RoleManagement', '', 'system:role:list', 2, 0, 'active'),
(14, 5, '菜单管理', '/system/menus', 'MenuManagement', '', 'system:menu:list', 3, 0, 'active');

-- 验证数据
SELECT '菜单数据' AS 'Table', COUNT(*) AS 'Count' FROM `menus` WHERE `deleted_at` IS NULL;
