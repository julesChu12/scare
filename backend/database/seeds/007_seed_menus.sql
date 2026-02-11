-- 007_seed_menus.sql
-- 初始化菜单数据，与当前管理后台路由结构一致
-- 注意：此文件使用 permission_code 字段（需要先执行 012_alter_menus_permission_code.sql）

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 清空现有菜单数据（如需重新初始化）
-- TRUNCATE TABLE `menus`;

-- ============================================
-- 插入顶级菜单
-- ============================================
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 工作台
(1, 0, '工作台', '/dashboard', 'Dashboard', 'Odometer', 'dashboard', 1, 0, 'active'),

-- 服务管理
(2, 0, '服务管理', '/services', '', 'Service', 'service', 2, 0, 'active'),

-- 数据中心
(3, 0, '数据中心', '/statistics', '', 'DataAnalysis', 'statistics', 3, 0, 'active'),

-- 居民管理
(4, 0, '居民管理', '/residents', '', 'User', 'resident', 4, 0, 'active'),

-- 站点管理
(5, 0, '站点管理', '/stations', '', 'OfficeBuilding', 'station', 5, 0, 'active'),

-- 内容管理
(7, 0, '内容管理', '/content', '', 'Document', 'content', 6, 0, 'active'),

-- 系统管理（放在最后）
(6, 0, '系统管理', '/system', '', 'Setting', 'system', 7, 0, 'active');

-- ============================================
-- 插入二级菜单 - 服务管理
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 需求管理（station_manager, admin）
(2, '需求管理', '/services/requests', 'RequestManagement', '', 'service:request:list', 1, 0, 'active'),
-- 任务池（staff, station_manager, admin）
(2, '任务池', '/services/tasks/pool', 'TaskPool', '', 'service:task:pool', 2, 0, 'active'),
-- 我的任务（staff）
(2, '我的任务', '/services/tasks/my', 'MyTasks', '', 'service:task:my', 3, 0, 'active'),
-- 任务详情（隐藏菜单，所有角色可访问）
(2, '任务详情', '/services/tasks/:id', 'TaskDetail', '', 'service:task:detail', 4, 1, 'active');

-- ============================================
-- 插入二级菜单 - 数据中心（新增）
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 统计概览（station_manager, admin）
(3, '统计概览', '/statistics/overview', 'StatisticsOverview', '', 'statistics:overview', 1, 0, 'active'),
-- 报表导出（station_manager, admin）
(3, '报表导出', '/statistics/reports', 'StatisticsReports', '', 'statistics:reports', 2, 0, 'active');

-- ============================================
-- 插入二级菜单 - 居民管理（新增）
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 老年人档案（station_manager, admin）
(4, '老年人档案', '/residents/elderly', 'ElderlyManagement', '', 'resident:elderly:list', 1, 0, 'active'),
-- 档案详情（隐藏菜单）
(4, '档案详情', '/residents/elderly/:id', 'ElderlyDetail', '', 'resident:elderly:detail', 2, 1, 'active');

-- ============================================
-- 插入二级菜单 - 站点管理
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 站点列表（admin）
(5, '站点列表', '/stations/list', 'StationManagement', '', 'station:list:view', 1, 0, 'active'),
-- 服务围栏（admin）
(5, '服务围栏', '/stations/zones', 'ZoneManagement', '', 'station:zone:list', 2, 0, 'active');

-- ============================================
-- 插入二级菜单 - 系统管理
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 用户管理（admin）
(6, '用户管理', '/system/users', 'UserManagement', '', 'system:user:list', 1, 0, 'active'),
-- 角色管理（admin）
(6, '角色权限', '/system/roles', 'RolePermission', '', 'system:role:list', 2, 0, 'active'),
-- 菜单管理（admin）
(6, '菜单管理', '/system/menus', 'MenuManagement', '', 'system:menu:list', 3, 0, 'active');

-- ============================================
-- 插入二级菜单 - 内容管理
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 轮播图管理（admin）
(7, '轮播图管理', '/content/banners', 'BannerManagement', '', 'content:banner:list', 1, 0, 'active'),
-- 通知管理（admin）
(7, '通知管理', '/content/notifications', 'NotificationManagement', '', 'content:notification:list', 2, 0, 'active');

-- ============================================
-- 插入独立菜单 - 个人信息（新增）
-- ============================================
INSERT INTO `menus` (`parent_id`, `name`, `path`, `component`, `icon`, `permission_code`, `sort`, `hidden`, `status`) VALUES
-- 个人信息（所有角色可访问，不在侧边栏显示，通过头像下拉菜单访问）
(0, '个人信息', '/profile', 'Profile', 'User', 'profile', 99, 1, 'active');
