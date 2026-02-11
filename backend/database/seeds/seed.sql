-- =====================================================
-- sCare 平台测试数据种子入口（兼容单文件执行）
-- 说明：本文件由 modules/*.sql 合并生成，便于 docker exec 直接导入
-- =====================================================

-- >>> BEGIN: database/seeds/modules/00_reset_all.sql
-- 00_reset_all.sql
-- 清空业务表（用于完整重建测试数据）

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE role_permissions;
TRUNCATE TABLE permissions;
TRUNCATE TABLE roles;
TRUNCATE TABLE menus;
TRUNCATE TABLE task_histories;
TRUNCATE TABLE task_assignments;
TRUNCATE TABLE service_requests;
TRUNCATE TABLE service_zones;
TRUNCATE TABLE service_stations;
TRUNCATE TABLE notifications;
TRUNCATE TABLE customer_profiles;
TRUNCATE TABLE user_identities;
TRUNCATE TABLE news;
TRUNCATE TABLE banners;
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;

-- <<< END: database/seeds/modules/00_reset_all.sql

-- >>> BEGIN: database/seeds/modules/10_roles_permissions.sql
-- 10_roles_permissions.sql
-- 角色与权限模块初始化

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- =====================================================
-- 1. 初始化角色（B端）
-- =====================================================
INSERT INTO `roles` (`id`, `code`, `name`, `description`, `is_system`, `status`, `sort`, `created_at`, `updated_at`) VALUES
(1, 'admin', '系统管理员', '拥有系统全部权限', 1, 'active', 1, NOW(), NOW()),
(2, 'station_manager', '站点管理员', '管理所属站点的服务和人员', 1, 'active', 2, NOW(), NOW()),
(3, 'staff', '工作人员', '执行服务任务的一线人员', 1, 'active', 3, NOW(), NOW());

-- =====================================================
-- 2. 初始化权限
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
-- 公共模块
('public', '公共权限', '所有登录用户都有的权限', 'public', 'menu', 0, '', '', 1, 0),
('public:auth', '认证相关', '认证相关公共接口', 'public', 'menu', 1, '', '', 1, 1),
('public:auth:me', '获取当前用户', '获取当前登录用户信息', 'public', 'resource', 2, '/api/v1/*/auth/me', 'GET', 1, 1),
('public:auth:logout', '登出', '用户登出', 'public', 'resource', 2, '/api/v1/*/auth/logout', 'POST', 1, 2),
('public:upload', '文件上传', '文件上传相关', 'public', 'menu', 1, '', '', 1, 2),
('public:upload:file', '上传文件', '上传文件接口', 'public', 'resource', 5, '/api/v1/*/upload', 'POST', 1, 1),
('public:notification', '通知管理', '通知相关公共接口', 'public', 'menu', 1, '', '', 1, 3),
('public:notification:list', '查看通知', '查看通知列表', 'public', 'resource', 7, '/api/v1/*/notifications', 'GET', 1, 1),
('public:notification:read', '标记已读', '标记通知已读', 'public', 'resource', 7, '/api/v1/*/notifications/*/read', 'POST', 1, 2);

-- =====================================================
-- 仪表盘模块
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
('dashboard', '工作台', '工作台菜单', 'dashboard', 'menu', 0, '', '', 0, 10);

-- =====================================================
-- 服务管理模块
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
-- 服务管理父菜单
('service', '服务管理', '服务管理模块', 'service', 'menu', 0, '', '', 0, 20),

-- 服务请求
('service:request', '服务请求', '服务请求管理', 'service', 'menu', 12, '', '', 0, 1),
('service:request:list', '请求列表', '查看服务请求列表', 'service', 'button', 13, '/api/v1/b/requests', 'GET', 0, 1),
('service:request:detail', '请求详情', '查看服务请求详情', 'service', 'resource', 13, '/api/v1/b/requests/*', 'GET', 0, 2),

-- 任务管理
('service:task', '任务管理', '任务管理', 'service', 'menu', 12, '', '', 0, 2),
('service:task:pool', '任务池', '查看任务池', 'service', 'button', 16, '/api/v1/b/tasks/pool', 'GET', 0, 1),
('service:task:my', '我的任务', '查看我的任务', 'service', 'button', 16, '/api/v1/b/tasks/my', 'GET', 0, 2),
('service:task:claim', '认领任务', '认领任务', 'service', 'button', 16, '/api/v1/b/tasks/*/claim', 'POST', 0, 3),
('service:task:complete', '完成任务', '完成任务', 'service', 'button', 16, '/api/v1/b/tasks/*/complete', 'POST', 0, 4);

-- =====================================================
-- 站点管理模块
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
-- 站点管理父菜单
('station', '站点管理', '站点管理模块', 'station', 'menu', 0, '', '', 0, 30),

-- 站点列表
('station:list', '站点列表', '站点列表管理', 'station', 'menu', 21, '', '', 0, 1),
('station:list:view', '查看站点', '查看站点列表', 'station', 'button', 22, '/api/v1/b/stations', 'GET', 0, 1),
('station:list:detail', '站点详情', '查看站点详情', 'station', 'resource', 22, '/api/v1/b/stations/*', 'GET', 0, 2),
('station:list:create', '创建站点', '创建新站点', 'station', 'button', 22, '/api/v1/b/stations', 'POST', 0, 3),
('station:list:update', '编辑站点', '编辑站点信息', 'station', 'button', 22, '/api/v1/b/stations/*', 'PUT', 0, 4),
('station:list:delete', '删除站点', '删除站点', 'station', 'button', 22, '/api/v1/b/stations/*', 'DELETE', 0, 5),

-- 服务围栏
('station:zone', '服务围栏', '服务围栏管理', 'station', 'menu', 21, '', '', 0, 2),
('station:zone:list', '围栏列表', '查看围栏列表', 'station', 'button', 28, '/api/v1/b/zones', 'GET', 0, 1),
('station:zone:create', '创建围栏', '创建新围栏', 'station', 'button', 28, '/api/v1/b/zones', 'POST', 0, 2),
('station:zone:update', '编辑围栏', '编辑围栏信息', 'station', 'button', 28, '/api/v1/b/zones/*', 'PUT', 0, 3),
('station:zone:delete', '删除围栏', '删除围栏', 'station', 'button', 28, '/api/v1/b/zones/*', 'DELETE', 0, 4);

-- =====================================================
-- 系统管理模块
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
-- 系统管理父菜单
('system', '系统管理', '系统管理模块', 'system', 'menu', 0, '', '', 0, 40),

-- 用户管理
('system:user', '用户管理', '用户管理', 'system', 'menu', 33, '', '', 0, 1),
('system:user:list', '用户列表', '查看用户列表', 'system', 'button', 34, '/api/v1/b/users', 'GET', 0, 1),
('system:user:detail', '用户详情', '查看用户详情', 'system', 'resource', 34, '/api/v1/b/users/*', 'GET', 0, 2),
('system:user:create', '创建用户', '创建新用户', 'system', 'button', 34, '/api/v1/b/users', 'POST', 0, 3),
('system:user:update', '编辑用户', '编辑用户信息', 'system', 'button', 34, '/api/v1/b/users/*', 'PUT', 0, 4),
('system:user:roles', '分配角色', '分配用户角色', 'system', 'button', 34, '/api/v1/b/users/*/identities', 'PUT', 0, 5),

-- 角色管理
('system:role', '角色管理', '角色权限管理', 'system', 'menu', 33, '', '', 0, 2),
('system:role:list', '角色列表', '查看角色列表', 'system', 'button', 40, '/api/v1/b/roles', 'GET', 0, 1),
('system:role:permissions', '角色权限', '查看角色权限', 'system', 'button', 40, '/api/v1/b/roles/*/permissions', 'GET', 0, 2),
('system:role:update', '更新权限', '更新角色权限', 'system', 'button', 40, '/api/v1/b/roles/*/permissions', 'PUT', 0, 3),
('system:permission:tree', '权限树', '查看权限树', 'system', 'button', 40, '/api/v1/b/permissions/tree', 'GET', 0, 4),

-- 菜单管理
('system:menu', '菜单管理', '菜单管理', 'system', 'menu', 33, '', '', 0, 3),
('system:menu:list', '菜单列表', '查看菜单列表', 'system', 'button', 45, '/api/v1/b/menus', 'GET', 0, 1),
('system:menu:detail', '菜单详情', '查看菜单详情', 'system', 'resource', 45, '/api/v1/b/menus/*', 'GET', 0, 2),
('system:menu:create', '创建菜单', '创建新菜单', 'system', 'button', 45, '/api/v1/b/menus', 'POST', 0, 3),
('system:menu:update', '编辑菜单', '编辑菜单信息', 'system', 'button', 45, '/api/v1/b/menus/*', 'PUT', 0, 4),
('system:menu:delete', '删除菜单', '删除菜单', 'system', 'button', 45, '/api/v1/b/menus/*', 'DELETE', 0, 5),
('system:menu:sort', '菜单排序', '批量更新菜单排序', 'system', 'button', 45, '/api/v1/b/menus/sort', 'PUT', 0, 6);

-- =====================================================
-- 内容管理模块
-- =====================================================
INSERT INTO `permissions` (`code`, `name`, `description`, `module`, `type`, `parent_id`, `api_path`, `api_method`, `is_public`, `sort`) VALUES
-- 内容管理父菜单
('content', '内容管理', '内容管理模块', 'content', 'menu', 0, '', '', 0, 50),

-- 轮播图管理
('content:banner', '轮播图管理', '轮播图管理', 'content', 'menu', 52, '', '', 0, 1),
('content:banner:list', '轮播图列表', '查看轮播图列表', 'content', 'button', 53, '/api/v1/b/banners', 'GET', 0, 1),
('content:banner:create', '创建轮播图', '创建新轮播图', 'content', 'button', 53, '/api/v1/b/banners', 'POST', 0, 2),
('content:banner:update', '编辑轮播图', '编辑轮播图信息', 'content', 'button', 53, '/api/v1/b/banners/*', 'PUT', 0, 3),
('content:banner:delete', '删除轮播图', '删除轮播图', 'content', 'button', 53, '/api/v1/b/banners/*', 'DELETE', 0, 4),

-- 新闻管理
('content:news', '新闻管理', '新闻管理', 'content', 'menu', 52, '', '', 0, 2),
('content:news:list', '新闻列表', '查看新闻列表', 'content', 'button', 58, '/api/v1/b/news', 'GET', 0, 1),
('content:news:detail', '新闻详情', '查看新闻详情', 'content', 'resource', 58, '/api/v1/b/news/*', 'GET', 0, 2),
('content:news:create', '创建新闻', '创建新新闻', 'content', 'button', 58, '/api/v1/b/news', 'POST', 0, 3),
('content:news:update', '编辑新闻', '编辑新闻信息', 'content', 'button', 58, '/api/v1/b/news/*', 'PUT', 0, 4),
('content:news:delete', '删除新闻', '删除新闻', 'content', 'button', 58, '/api/v1/b/news/*', 'DELETE', 0, 5);

-- =====================================================
-- 角色权限分配
-- =====================================================

-- Admin 角色拥有所有权限（在代码中特殊处理，这里不需要插入）

-- Station Manager 角色权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM `roles` r
JOIN `permissions` p ON 1=1
WHERE r.code = 'station_manager'
  AND p.`code` IN (
    'dashboard',
    'service', 'service:request', 'service:request:list', 'service:request:detail',
    'service:task', 'service:task:pool', 'service:task:my', 'service:task:claim', 'service:task:complete',
    'station', 'station:list', 'station:list:view', 'station:list:detail',
    'station:zone', 'station:zone:list', 'station:zone:create', 'station:zone:update', 'station:zone:delete',
    'system', 'system:user', 'system:user:list', 'system:user:detail', 'system:user:create', 'system:user:update', 'system:user:roles'
);

-- Staff 角色权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM `roles` r
JOIN `permissions` p ON 1=1
WHERE r.code = 'staff'
  AND p.`code` IN (
    'dashboard',
    'service', 'service:request', 'service:request:list', 'service:request:detail',
    'service:task', 'service:task:pool', 'service:task:my', 'service:task:claim', 'service:task:complete',
    'station', 'station:list', 'station:list:view',
    'station:zone', 'station:zone:list'
);

-- <<< END: database/seeds/modules/10_roles_permissions.sql

-- >>> BEGIN: database/seeds/modules/20_menus.sql
-- 20_menus.sql
-- 菜单模块初始化
-- 注意：此文件依赖 menus.permission_code 字段（对应迁移: 002_alter_menus_permission_code.sql）

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

-- <<< END: database/seeds/modules/20_menus.sql

-- >>> BEGIN: database/seeds/modules/30_stations_zones.sql
-- 30_stations_zones.sql
-- 站点与围栏模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `service_stations` (`id`, `name`, `code`, `address`, `phone`, `latitude`, `longitude`, `status`, `created_at`, `updated_at`) VALUES
(1, '朝阳区幸福街道服务站', 'CYXF001', '北京市朝阳区幸福大街100号', '010-12345678', 39.9200, 116.4500, 'active', NOW(), NOW()),
(2, '朝阳区康乐街道服务站', 'CYKL002', '北京市朝阳区康乐大街200号', '010-87654321', 39.9100, 116.4600, 'active', NOW(), NOW());

-- =====================================================
-- 2.1 初始化服务区域/地理围栏（service_zones）
-- =====================================================
INSERT INTO `service_zones` (`id`, `station_id`, `name`, `points`, `priority`, `status`, `created_at`, `updated_at`) VALUES
-- 站点1服务区域（幸福街道）
(1, 1, '幸福街道服务区',
 '[[116.44,39.91],[116.46,39.91],[116.46,39.93],[116.44,39.93],[116.44,39.91]]',
 1, 'active', NOW(), NOW()),
-- 站点2服务区域（康乐街道）
(2, 2, '康乐街道服务区',
 '[[116.45,39.90],[116.47,39.90],[116.47,39.92],[116.45,39.92],[116.45,39.90]]',
 1, 'active', NOW(), NOW());

-- <<< END: database/seeds/modules/30_stations_zones.sql

-- >>> BEGIN: database/seeds/modules/40_users_profiles.sql
-- 40_users_profiles.sql
-- 用户、身份、客户档案模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 密码哈希: $2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `gender`, `birth_date`, `station_id`, `status`, `created_at`, `updated_at`) VALUES
-- B端用户
(1, '13800000001', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '系统管理员', 'admin@scare.com', 'male', '1980-01-01', NULL, 'active', NOW(), NOW()),
(2, '13800000002', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李站长', 'lizhang@scare.com', 'male', '1975-05-15', 1, 'active', NOW(), NOW()),
(3, '13800000003', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王站长', 'wangzhang@scare.com', 'female', '1978-08-20', 2, 'active', NOW(), NOW()),
(4, '13800000004', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王小红', 'xiaohong@scare.com', 'female', '1990-03-10', 1, 'active', NOW(), NOW()),
(5, '13800000005', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '刘师傅', 'liushifu@scare.com', 'male', '1985-07-22', 1, 'active', NOW(), NOW()),
(6, '13800000006', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈护士', 'chenhushi@scare.com', 'female', '1992-11-05', 2, 'active', NOW(), NOW()),
(7, '13800000007', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵大哥', 'zhaodage@scare.com', 'male', '1988-09-18', 2, 'active', NOW(), NOW()),

-- C端用户（服务对象）
(8, '13800000008', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '张大爷', NULL, 'male', '1950-05-15', NULL, 'active', NOW(), NOW()),
(9, '13800000009', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李奶奶', NULL, 'female', '1955-03-20', 1, 'active', NOW(), NOW()),  -- 跨端用户（员工+客户）
(10, '13800000010', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王爷爷', NULL, 'male', '1948-11-10', NULL, 'active', NOW(), NOW()),
(11, '13800000011', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙女士', NULL, 'female', '1990-06-25', NULL, 'active', NOW(), NOW()),  -- 孕妇
(12, '13800000012', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵先生', NULL, 'male', '1965-02-14', NULL, 'active', NOW(), NOW()),  -- 失能人士
(13, '13800000013', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '小明', NULL, 'male', '2018-03-15', NULL, 'active', NOW(), NOW());  -- 儿童

-- =====================================================
-- 4. 初始化用户身份（user_identities）
-- B端身份：admin, station_manager, staff
-- C端身份：elderly, family, pregnant, disabled, child
-- =====================================================
INSERT INTO `user_identities` (`user_id`, `identity_type`, `is_primary`, `station_id`, `status`, `granted_at`, `created_at`, `updated_at`) VALUES
-- B端身份（8条）
(1, 'admin', 1, NULL, 'active', NOW(), NOW(), NOW()),           -- 系统管理员
(2, 'station_manager', 1, 1, 'active', NOW(), NOW(), NOW()),    -- 李站长 - 站点1
(3, 'station_manager', 1, 2, 'active', NOW(), NOW(), NOW()),    -- 王站长 - 站点2
(4, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),              -- 王小红 - 站点1
(5, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),              -- 刘师傅 - 站点1
(6, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),              -- 陈护士 - 站点2
(7, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),              -- 赵大哥 - 站点2
(9, 'staff', 0, 1, 'active', NOW(), NOW(), NOW()),              -- 李奶奶的副身份（志愿者）

-- C端身份（6条）
(8, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),         -- 张大爷 - 老年人
(9, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),         -- 李奶奶 - 老年人（主身份）
(10, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),        -- 王爷爷 - 老年人
(11, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),       -- 孙女士 - 孕妇
(12, 'disabled', 1, NULL, 'active', NOW(), NOW(), NOW()),       -- 赵先生 - 失能人士
(13, 'child', 1, NULL, 'active', NOW(), NOW(), NOW());          -- 小明 - 儿童

-- =====================================================
-- 5. 初始化客户档案（C端服务对象）
-- =====================================================
INSERT INTO `customer_profiles` (
    `user_id`,
    `customer_type`,
    `gender`,
    `birth_date`,
    `address`,
    `health_status`,
    `disability_level`,
    `medical_history`,
    `special_needs`,
    `emergency_contact`,
    `created_at`,
    `updated_at`
) VALUES
-- 张大爷：老年人
(8, 'elderly', 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '高血压，需要定期测量血压', '每周需要社区医生上门测血压',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}', NOW(), NOW()),

-- 李奶奶：老年人（跨端用户）
(9, 'elderly', 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '糖尿病，行动不便', '需要助行器，定期血糖监测',
 '{"name":"李华","phone":"13900000002","relation":"子女"}', NOW(), NOW()),

-- 王爷爷：老年人
(10, 'elderly', 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '心脏病，中风后遗症', '需要轮椅，24小时护理',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}', NOW(), NOW()),

-- 孙女士：孕妇
(11, 'pregnant', 'female', '1990-06-25', '北京市朝阳区康乐小区5号楼501', '良好', NULL,
 '孕27周，定期产检', '需要产前护理指导',
 '{"name":"孙先生","phone":"13900000004","relation":"配偶"}', NOW(), NOW()),

-- 赵先生：失能人士
(12, 'disabled', 'male', '1965-02-14', '北京市朝阳区康乐小区6号楼602', '较差', '重度失能',
 '脊髓损伤，下肢瘫痪', '需要专业康复护理，定期更换导尿管',
 '{"name":"赵女士","phone":"13900000005","relation":"配偶"}', NOW(), NOW()),

-- 小明：儿童
(13, 'child', 'male', '2018-03-15', '北京市朝阳区幸福小区4号楼404', '良好', '自理',
 '无重大病史', '课后托管服务',
 '{"name":"小明妈妈","phone":"13900000006","relation":"母亲"}', NOW(), NOW());

-- <<< END: database/seeds/modules/40_users_profiles.sql

-- >>> BEGIN: database/seeds/modules/50_requests_tasks.sql
-- 50_requests_tasks.sql
-- 服务请求、任务、任务历史模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `service_requests` (
    `id`, `request_no`, `user_id`, `service_type`, `status`,
    `description`, `submit_location_lat`, `submit_location_lng`,
    `contact_name`, `contact_phone`, `address`, `appointment_time`,
    `urgency`, `station_id`, `rating`, `feedback`,
    `created_at`, `updated_at`
) VALUES
-- 已完成的请求
(1, 'REQ-2026020701', 8, 'meal', 'completed',
 '需要午餐送餐服务，清淡饮食', 39.9042, 116.4074,
 '张大爷', '13800000008', '北京市朝阳区幸福小区1号楼101',
 '2026-02-07 12:00:00', 'normal', 1, 5, '服务很好，送餐准时',
 '2026-02-07 08:00:00', '2026-02-07 12:30:00'),

-- 已认领的请求
(2, 'REQ-2026020702', 9, 'medical', 'claimed',
 '需要陪同去医院复查', 39.9042, 116.4074,
 '李奶奶', '13800000009', '北京市朝阳区幸福小区2号楼202',
 '2026-02-08 09:00:00', 'urgent', 1, 0, NULL,
 '2026-02-07 10:00:00', '2026-02-07 10:30:00'),

-- 待认领的请求
(3, 'REQ-2026020703', 10, 'care', 'dispatched',
 '需要日常照护服务，协助洗漱', 39.9042, 116.4074,
 '王爷爷', '13800000010', '北京市朝阳区幸福小区3号楼303',
 '2026-02-08 08:00:00', 'normal', 2, 0, NULL,
 '2026-02-07 14:00:00', '2026-02-07 14:00:00'),

-- 已取消的请求
(4, 'REQ-2026020704', 11, 'other', 'cancelled',
 '需要产前护理指导', 39.9042, 116.4074,
 '孙女士', '13800000011', '北京市朝阳区康乐小区5号楼501',
 '2026-02-09 10:00:00', 'normal', 2, 0, NULL,
 '2026-02-07 15:00:00', '2026-02-07 15:30:00');

-- =====================================================
-- 7. 初始化任务分配（task_assignments）
-- 任务从服务请求派生，状态与请求保持一致
-- =====================================================
INSERT INTO `task_assignments` (
    `id`, `request_id`, `station_id`, `staff_id`, `status`,
    `claimed_at`, `completed_at`, `rating`, `feedback`, `images`,
    `created_at`, `updated_at`
) VALUES
-- 已完成的任务（对应请求1）
(1, 1, 1, 4, 'completed',
 '2026-02-07 08:30:00', '2026-02-07 12:30:00', 5, '服务很好，送餐准时', '[]',
 '2026-02-07 08:00:00', '2026-02-07 12:30:00'),

-- 已认领的任务（对应请求2）
(2, 2, 1, 5, 'claimed',
 '2026-02-07 10:30:00', NULL, 0, NULL, NULL,
 '2026-02-07 10:00:00', '2026-02-07 10:30:00'),

-- 待认领的任务（对应请求3）
(3, 3, 2, 0, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-07 14:00:00', '2026-02-07 14:00:00'),

-- 已取消的任务（对应请求4，状态同步）
(4, 4, 2, 0, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-07 15:00:00', '2026-02-07 15:30:00');

-- =====================================================
-- 8. 初始化任务历史（task_histories）
-- 记录任务状态变更的审计日志
-- =====================================================
INSERT INTO `task_histories` (
    `id`, `task_id`, `request_id`, `action`, `operator_id`,
    `from_staff_id`, `to_staff_id`, `from_station_id`, `to_station_id`,
    `status_before`, `status_after`, `remark`, `created_at`
) VALUES
-- 任务1（已完成）的历史
(1, 1, 1, 'dispatched', 1, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-07 08:00:00'),
(2, 1, 1, 'claimed', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工认领任务', '2026-02-07 08:30:00'),
(3, 1, 1, 'completed', 4, 4, 4, NULL, NULL, 'claimed', 'completed', '服务完成', '2026-02-07 12:30:00'),

-- 任务2（已认领）的历史
(4, 2, 2, 'dispatched', 1, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-07 10:00:00'),
(5, 2, 2, 'claimed', 5, NULL, 5, NULL, NULL, 'dispatched', 'claimed', '员工认领任务', '2026-02-07 10:30:00'),

-- 任务3（待认领）的历史
(6, 3, 3, 'dispatched', 1, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-07 14:00:00'),

-- 任务4（已取消）的历史
(7, 4, 4, 'dispatched', 1, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-07 15:00:00'),
(8, 4, 4, 'cancelled', 11, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '用户取消请求', '2026-02-07 15:30:00');

-- <<< END: database/seeds/modules/50_requests_tasks.sql

-- >>> BEGIN: database/seeds/modules/60_content.sql
-- 60_content.sql
-- 内容模块初始化（Banner + News）

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `banners` (`id`, `station_id`, `title`, `image_url`, `link_type`, `link_value`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '社区养老服务平台', 'https://via.placeholder.com/800x400/4A90E2/FFFFFF?text=社区养老服务', 'none', '', 100, 'active', NOW(3), NOW(3)),
(2, 0, '专业护理团队', 'https://via.placeholder.com/800x400/50C878/FFFFFF?text=专业护理', 'none', '', 90, 'active', NOW(3), NOW(3)),
(3, 0, '24小时服务热线', 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=24小时服务', 'none', '', 80, 'active', NOW(3), NOW(3));

INSERT INTO `news` (`id`, `title`, `summary`, `content`, `cover_url`, `type`, `status`, `station_id`, `author_id`, `publish_at`, `view_count`, `created_at`, `updated_at`) VALUES
(1, '社区养老服务中心正式启用', '首批养老服务站点完成部署并投入使用。', '社区养老服务中心正式启用，为居民提供送餐、陪诊、照护等服务。', 'https://via.placeholder.com/600x300/4A90E2/FFFFFF?text=News+1', 'news', 'published', 0, 1, NOW(), 120, NOW(3), NOW(3)),
(2, '春季义诊活动公告', '本周六开展春季社区义诊，请有需要的居民提前预约。', '本周六上午9点到12点，在朝阳区幸福街道服务站开展义诊活动。', 'https://via.placeholder.com/600x300/50C878/FFFFFF?text=News+2', 'notice', 'published', 0, 1, NOW(), 80, NOW(3), NOW(3));

-- <<< END: database/seeds/modules/60_content.sql

-- >>> BEGIN: database/seeds/modules/70_notifications.sql
-- 70_notifications.sql
-- 通知模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `notifications` (`id`, `user_id`, `title`, `body`, `type`, `related_id`, `related_type`, `channel`, `send_status`, `is_read`, `retry_count`, `created_at`, `updated_at`) VALUES
(1, 2, '系统通知', '欢迎使用社区养老服务平台。', 'system', 0, 'system', 'in_app', 'sent', 0, 0, NOW(3), NOW(3)),
(2, 4, '任务提醒', '您有新的待处理任务，请及时查看。', 'task', 3, 'task', 'in_app', 'sent', 0, 0, NOW(3), NOW(3));

-- <<< END: database/seeds/modules/70_notifications.sql

