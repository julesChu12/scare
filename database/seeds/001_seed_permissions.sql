-- =====================================================
-- 权限系统数据种子文件
-- 版本: v1.0.0
-- 说明: 权限、角色、角色权限关联
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `role_permissions`;
TRUNCATE TABLE `roles`;
TRUNCATE TABLE `permissions`;
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 权限数据 (62条)
-- =====================================================

-- 公共权限
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(1, 'public', '公共权限', '所有登录用户都有的权限', 'public', 'menu', 0, '', '', 1, 'active', 0),
(2, 'public:auth', '认证相关', '认证相关公共接口', 'public', 'menu', 1, '', '', 1, 'active', 1),
(3, 'public:auth:me', '获取当前用户', '获取当前登录用户信息', 'public', 'resource', 2, '/api/v1/*/auth/me', 'GET', 1, 'active', 1),
(4, 'public:auth:logout', '登出', '用户登出', 'public', 'resource', 2, '/api/v1/*/auth/logout', 'POST', 1, 'active', 2),
(5, 'public:upload', '文件上传', '文件上传相关', 'public', 'menu', 1, '', '', 1, 'active', 2),
(6, 'public:upload:file', '上传文件', '上传文件接口', 'public', 'resource', 5, '/api/v1/*/upload', 'POST', 1, 'active', 1),
(7, 'public:notification', '通知管理', '通知相关公共接口', 'public', 'menu', 1, '', '', 1, 'active', 3),
(8, 'public:notification:list', '查看通知', '查看通知列表', 'public', 'resource', 7, '/api/v1/*/notifications', 'GET', 1, 'active', 1),
(9, 'public:notification:read', '标记已读', '标记通知已读', 'public', 'resource', 7, '/api/v1/*/notifications/*/read', 'POST', 1, 'active', 2);

-- 工作台
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(10, 'dashboard', '工作台', '工作台菜单', 'dashboard', 'menu', 0, '', '', 0, 'active', 10);

-- 服务管理
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(11, 'service', '服务管理', '服务管理模块', 'service', 'menu', 0, '', '', 0, 'active', 20),
(12, 'service:request', '服务请求', '服务请求管理', 'service', 'menu', 11, '', '', 0, 'active', 1),
(13, 'service:request:list', '请求列表', '查看服务请求列表', 'service', 'button', 12, '/api/v1/b/requests', 'GET', 0, 'active', 1),
(14, 'service:request:detail', '请求详情', '查看服务请求详情', 'service', 'resource', 12, '/api/v1/b/requests/*', 'GET', 0, 'active', 2),
(15, 'service:task', '任务管理', '任务管理', 'service', 'menu', 11, '', '', 0, 'active', 2),
(16, 'service:task:pool', '任务池', '查看任务池', 'service', 'button', 15, '/api/v1/b/tasks/pool', 'GET', 0, 'active', 1),
(17, 'service:task:my', '我的任务', '查看我的任务', 'service', 'button', 15, '/api/v1/b/tasks/my', 'GET', 0, 'active', 2),
(18, 'service:task:claim', '认领任务', '认领任务', 'service', 'button', 15, '/api/v1/b/tasks/*/claim', 'POST', 0, 'active', 3),
(19, 'service:task:complete', '完成任务', '完成任务', 'service', 'button', 15, '/api/v1/b/tasks/*/complete', 'POST', 0, 'active', 4);

-- 站点管理
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(20, 'station', '站点管理', '站点管理模块', 'station', 'menu', 0, '', '', 0, 'active', 30),
(21, 'station:list', '站点列表', '站点列表管理', 'station', 'menu', 20, '', '', 0, 'active', 1),
(22, 'station:list:view', '查看站点', '查看站点列表', 'station', 'button', 21, '/api/v1/b/stations', 'GET', 0, 'active', 1),
(23, 'station:list:detail', '站点详情', '查看站点详情', 'station', 'resource', 21, '/api/v1/b/stations/*', 'GET', 0, 'active', 2),
(24, 'station:list:create', '创建站点', '创建新站点', 'station', 'button', 21, '/api/v1/b/stations', 'POST', 0, 'active', 3),
(25, 'station:list:update', '编辑站点', '编辑站点信息', 'station', 'button', 21, '/api/v1/b/stations/*', 'PUT', 0, 'active', 4),
(26, 'station:list:delete', '删除站点', '删除站点', 'station', 'button', 21, '/api/v1/b/stations/*', 'DELETE', 0, 'active', 5),
(27, 'station:zone', '服务围栏', '服务围栏管理', 'station', 'menu', 20, '', '', 0, 'active', 2),
(28, 'station:zone:list', '围栏列表', '查看围栏列表', 'station', 'button', 27, '/api/v1/b/zones', 'GET', 0, 'active', 1),
(29, 'station:zone:create', '创建围栏', '创建新围栏', 'station', 'button', 27, '/api/v1/b/zones', 'POST', 0, 'active', 2),
(30, 'station:zone:update', '编辑围栏', '编辑围栏信息', 'station', 'button', 27, '/api/v1/b/zones/*', 'PUT', 0, 'active', 3),
(31, 'station:zone:delete', '删除围栏', '删除围栏', 'station', 'button', 27, '/api/v1/b/zones/*', 'DELETE', 0, 'active', 4);

-- 系统管理
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(32, 'system', '系统管理', '系统管理模块', 'system', 'menu', 0, '', '', 0, 'active', 40),
(33, 'system:user', '用户管理', '用户管理', 'system', 'menu', 32, '', '', 0, 'active', 1),
(34, 'system:user:list', '用户列表', '查看用户列表', 'system', 'button', 33, '/api/v1/b/users', 'GET', 0, 'active', 1),
(35, 'system:user:detail', '用户详情', '查看用户详情', 'system', 'resource', 33, '/api/v1/b/users/*', 'GET', 0, 'active', 2),
(36, 'system:user:create', '创建用户', '创建新用户', 'system', 'button', 33, '/api/v1/b/users', 'POST', 0, 'active', 3),
(37, 'system:user:update', '编辑用户', '编辑用户信息', 'system', 'button', 33, '/api/v1/b/users/*', 'PUT', 0, 'active', 4),
(38, 'system:user:roles', '分配角色', '分配用户角色', 'system', 'button', 33, '/api/v1/b/users/*/roles', 'PUT', 0, 'active', 5),
(39, 'system:role', '角色管理', '角色权限管理', 'system', 'menu', 32, '', '', 0, 'active', 2),
(40, 'system:role:list', '角色列表', '查看角色列表', 'system', 'button', 39, '/api/v1/b/roles', 'GET', 0, 'active', 1),
(41, 'system:role:permissions', '角色权限', '查看角色权限', 'system', 'button', 39, '/api/v1/b/roles/*/permissions', 'GET', 0, 'active', 2),
(42, 'system:role:update', '更新权限', '更新角色权限', 'system', 'button', 39, '/api/v1/b/roles/*/permissions', 'PUT', 0, 'active', 3),
(43, 'system:permission:tree', '权限树', '查看权限树', 'system', 'button', 39, '/api/v1/b/permissions/tree', 'GET', 0, 'active', 4),
(44, 'system:menu', '菜单管理', '菜单管理', 'system', 'menu', 32, '', '', 0, 'active', 3),
(45, 'system:menu:list', '菜单列表', '查看菜单列表', 'system', 'button', 44, '/api/v1/b/menus', 'GET', 0, 'active', 1),
(46, 'system:menu:detail', '菜单详情', '查看菜单详情', 'system', 'resource', 44, '/api/v1/b/menus/*', 'GET', 0, 'active', 2),
(47, 'system:menu:create', '创建菜单', '创建新菜单', 'system', 'button', 44, '/api/v1/b/menus', 'POST', 0, 'active', 3),
(48, 'system:menu:update', '编辑菜单', '编辑菜单信息', 'system', 'button', 44, '/api/v1/b/menus/*', 'PUT', 0, 'active', 4),
(49, 'system:menu:delete', '删除菜单', '删除菜单', 'system', 'button', 44, '/api/v1/b/menus/*', 'DELETE', 0, 'active', 5),
(50, 'system:menu:sort', '菜单排序', '批量更新菜单排序', 'system', 'button', 44, '/api/v1/b/menus/sort', 'PUT', 0, 'active', 6);

-- 内容管理
INSERT INTO `permissions` (`id`, `code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `status`, `sort`) VALUES
(51, 'content', '内容管理', '内容管理模块', 'content', 'menu', 0, '', '', 0, 'active', 50),
(52, 'content:banner', '轮播图管理', '轮播图管理', 'content', 'menu', 51, '', '', 0, 'active', 1),
(53, 'content:banner:list', '轮播图列表', '查看轮播图列表', 'content', 'button', 52, '/api/v1/b/banners', 'GET', 0, 'active', 1),
(54, 'content:banner:create', '创建轮播图', '创建新轮播图', 'content', 'button', 52, '/api/v1/b/banners', 'POST', 0, 'active', 2),
(55, 'content:banner:update', '编辑轮播图', '编辑轮播图信息', 'content', 'button', 52, '/api/v1/b/banners/*', 'PUT', 0, 'active', 3),
(56, 'content:banner:delete', '删除轮播图', '删除轮播图', 'content', 'button', 52, '/api/v1/b/banners/*', 'DELETE', 0, 'active', 4),
(57, 'content:news', '新闻管理', '新闻管理', 'content', 'menu', 51, '', '', 0, 'active', 2),
(58, 'content:news:list', '新闻列表', '查看新闻列表', 'content', 'button', 57, '/api/v1/b/news', 'GET', 0, 'active', 1),
(59, 'content:news:detail', '新闻详情', '查看新闻详情', 'content', 'resource', 57, '/api/v1/b/news/*', 'GET', 0, 'active', 2),
(60, 'content:news:create', '创建新闻', '创建新新闻', 'content', 'button', 57, '/api/v1/b/news', 'POST', 0, 'active', 3),
(61, 'content:news:update', '编辑新闻', '编辑新闻信息', 'content', 'button', 57, '/api/v1/b/news/*', 'PUT', 0, 'active', 4),
(62, 'content:news:delete', '删除新闻', '删除新闻', 'content', 'button', 57, '/api/v1/b/news/*', 'DELETE', 0, 'active', 5);

-- =====================================================
-- 角色数据 (5条)
-- =====================================================
INSERT INTO `roles` (`id`, `code`, `name`, `description`, `is_system`, `status`, `sort`) VALUES
(1, 'admin', '系统管理员', '拥有系统全部权限', 1, 'active', 1),
(2, 'station_manager', '站点管理员', '管理所属站点的服务和人员', 1, 'active', 2),
(3, 'staff', '工作人员', '执行服务任务的一线人员', 1, 'active', 3),
(4, 'elderly', '老年人', 'C端老年用户', 1, 'active', 4),
(5, 'family', '家属', 'C端家属用户', 1, 'active', 5);

-- =====================================================
-- 角色权限关联 (41条)
-- =====================================================

-- station_manager 角色权限 (26条)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(2, 10), (2, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 17), (2, 18), (2, 19),
(2, 20), (2, 21), (2, 22), (2, 23), (2, 27), (2, 28), (2, 29), (2, 30), (2, 31),
(2, 32), (2, 33), (2, 34), (2, 35), (2, 36), (2, 37), (2, 38);

-- staff 角色权限 (15条)
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES
(3, 10), (3, 11), (3, 12), (3, 13), (3, 14), (3, 15), (3, 16), (3, 17), (3, 18), (3, 19),
(3, 20), (3, 21), (3, 22), (3, 27), (3, 28);

-- 验证数据
SELECT '权限数据' AS 'Table', COUNT(*) AS 'Count' FROM `permissions`
UNION ALL SELECT '角色数据', COUNT(*) FROM `roles`
UNION ALL SELECT '角色权限', COUNT(*) FROM `role_permissions`;
