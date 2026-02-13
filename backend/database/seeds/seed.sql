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

-- 密码哈希: $2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm (Test@123)
-- id_card_masked: 脱敏值直接存储，id_card/id_card_hmac 由 encrypt_seed 工具生成后 UPDATE
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `gender`, `birth_date`, `id_card_masked`, `station_id`, `status`, `created_at`, `updated_at`) VALUES
-- B端用户（9人）
(1, '13800000001', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '系统管理员', 'admin@scare.com', 'male', '1980-01-01', '1101**********0011', NULL, 'active', NOW(), NOW()),
(2, '13800000002', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李站长', 'lizhang@scare.com', 'male', '1975-05-15', '1101**********0022', 1, 'active', NOW(), NOW()),
(3, '13800000003', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王站长', 'wangzhang@scare.com', 'female', '1978-08-20', '1101**********0033', 2, 'active', NOW(), NOW()),
(4, '13800000004', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王小红', 'xiaohong@scare.com', 'female', '1990-03-10', '1101**********0044', 1, 'active', NOW(), NOW()),
(5, '13800000005', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '刘师傅', 'liushifu@scare.com', 'male', '1985-07-22', '1101**********0055', 1, 'active', NOW(), NOW()),
(6, '13800000006', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈护士', 'chenhushi@scare.com', 'female', '1992-11-05', '1101**********0066', 2, 'active', NOW(), NOW()),
(7, '13800000007', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵大哥', 'zhaodage@scare.com', 'male', '1988-09-18', '1101**********0077', 2, 'active', NOW(), NOW()),
(24, '13800000024', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙小明', 'sunxm@scare.com', 'male', '1995-08-12', '1101**********0244', 1, 'active', NOW(), NOW()),
(25, '13800000025', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '周护工', 'zhouhg@scare.com', 'female', '1991-12-05', '1101**********0255', 2, 'active', NOW(), NOW()),

-- C端用户（16人）
(8, '13800000008', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '张大爷', NULL, 'male', '1950-05-15', '1101**********0088', NULL, 'active', NOW(), NOW()),
(9, '13800000009', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '李奶奶', NULL, 'female', '1955-03-20', '1101**********0099', 1, 'active', NOW(), NOW()),
(10, '13800000010', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '王爷爷', NULL, 'male', '1948-11-10', '1101**********0100', NULL, 'active', NOW(), NOW()),
(11, '13800000011', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '孙女士', NULL, 'female', '1990-06-25', '1101**********0111', NULL, 'active', NOW(), NOW()),
(12, '13800000012', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '赵先生', NULL, 'male', '1965-02-14', '1101**********0122', NULL, 'active', NOW(), NOW()),
(13, '13800000013', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '小明', NULL, 'male', '2018-03-15', '1101**********0133', NULL, 'active', NOW(), NOW()),
(14, '13800000014', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '周阿姨', NULL, 'female', '1949-03-08', '1101**********0144', NULL, 'active', NOW(), NOW()),
(15, '13800000015', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '吴大爷', NULL, 'male', '1951-07-12', '1101**********0155', NULL, 'active', NOW(), NOW()),
(16, '13800000016', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '郑奶奶', NULL, 'female', '1952-09-25', '1101**********0166', NULL, 'active', NOW(), NOW()),
(17, '13800000017', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '冯爷爷', NULL, 'male', '1946-05-18', '1101**********0177', NULL, 'active', NOW(), NOW()),
(18, '13800000018', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '陈阿姨', NULL, 'female', '1953-12-02', '1101**********0188', NULL, 'active', NOW(), NOW()),
(19, '13800000019', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '杨大爷', NULL, 'male', '1950-08-15', '1101**********0199', NULL, 'active', NOW(), NOW()),
(20, '13800000020', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '黄女士', NULL, 'female', '1992-05-10', '1101**********0200', NULL, 'active', NOW(), NOW()),
(21, '13800000021', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '林先生', NULL, 'male', '1970-03-20', '1101**********0211', NULL, 'active', NOW(), NOW()),
(22, '13800000022', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '何大爷', NULL, 'male', '1947-12-08', '1101**********0222', NULL, 'active', NOW(), NOW()),
(23, '13800000023', '$2a$10$xNeRb.BO6iQVzx440yJPPuc6g.OMrFNIpsNjQVjsT8740S33ocALm', '马奶奶', NULL, 'female', '1956-01-15', '1101**********0233', NULL, 'active', NOW(), NOW());

-- =====================================================
-- 4. 初始化用户身份（user_identities）
-- =====================================================
INSERT INTO `user_identities` (`user_id`, `identity_type`, `is_primary`, `station_id`, `status`, `granted_at`, `created_at`, `updated_at`) VALUES
-- B端身份（10条）
(1, 'admin', 1, NULL, 'active', NOW(), NOW(), NOW()),
(2, 'station_manager', 1, 1, 'active', NOW(), NOW(), NOW()),
(3, 'station_manager', 1, 2, 'active', NOW(), NOW(), NOW()),
(4, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(5, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(6, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),
(7, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),
(9, 'staff', 0, 1, 'active', NOW(), NOW(), NOW()),
(24, 'staff', 1, 1, 'active', NOW(), NOW(), NOW()),
(25, 'staff', 1, 2, 'active', NOW(), NOW(), NOW()),

-- C端身份（16条）
(8, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(9, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(10, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(11, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),
(12, 'disabled', 1, NULL, 'active', NOW(), NOW(), NOW()),
(13, 'child', 1, NULL, 'active', NOW(), NOW(), NOW()),
(14, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(15, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(16, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(17, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(18, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(19, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(20, 'pregnant', 1, NULL, 'active', NOW(), NOW(), NOW()),
(21, 'family', 1, NULL, 'active', NOW(), NOW(), NOW()),
(22, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW()),
(23, 'elderly', 1, NULL, 'active', NOW(), NOW(), NOW());

-- =====================================================
-- 5. 初始化客户档案（C端服务对象）
-- =====================================================
INSERT INTO `customer_profiles` (
    `user_id`, `customer_type`, `gender`, `birth_date`, `address`,
    `health_status`, `disability_level`, `medical_history`, `special_needs`,
    `emergency_contact`, `created_at`, `updated_at`
) VALUES
(8, 'elderly', 'male', '1950-05-15', '北京市朝阳区幸福小区1号楼101', '良好', '自理',
 '高血压，需要定期测量血压', '每周需要社区医生上门测血压',
 '{"name":"张小明","phone":"13900000001","relation":"子女"}', NOW(), NOW()),

(9, 'elderly', 'female', '1955-03-20', '北京市朝阳区幸福小区2号楼202', '一般', '轻度失能',
 '糖尿病，行动不便', '需要助行器，定期血糖监测',
 '{"name":"李华","phone":"13900000002","relation":"子女"}', NOW(), NOW()),

(10, 'elderly', 'male', '1948-11-10', '北京市朝阳区幸福小区3号楼303', '较差', '中度失能',
 '心脏病，中风后遗症', '需要轮椅，24小时护理',
 '{"name":"王芳","phone":"13900000003","relation":"子女"}', NOW(), NOW()),

(11, 'pregnant', 'female', '1990-06-25', '北京市朝阳区康乐小区5号楼501', '良好', NULL,
 '孕27周，定期产检', '需要产前护理指导',
 '{"name":"孙先生","phone":"13900000004","relation":"配偶"}', NOW(), NOW()),

(12, 'disabled', 'male', '1965-02-14', '北京市朝阳区康乐小区6号楼602', '较差', '重度失能',
 '脊髓损伤，下肢瘫痪', '需要专业康复护理，定期更换导尿管',
 '{"name":"赵女士","phone":"13900000005","relation":"配偶"}', NOW(), NOW()),

(13, 'child', 'male', '2018-03-15', '北京市朝阳区幸福小区4号楼404', '良好', '自理',
 '无重大病史', '课后托管服务',
 '{"name":"小明妈妈","phone":"13900000006","relation":"母亲"}', NOW(), NOW()),

(14, 'elderly', 'female', '1949-03-08', '北京市朝阳区幸福小区5号楼505', '一般', '轻度失能',
 '关节炎，腰椎间盘突出', '需要定期理疗和康复训练',
 '{"name":"周明","phone":"13900000007","relation":"子女"}', NOW(), NOW()),

(15, 'elderly', 'male', '1951-07-12', '北京市朝阳区幸福小区6号楼606', '良好', '自理',
 '轻度白内障', '需要定期眼科检查',
 '{"name":"吴丽","phone":"13900000008","relation":"子女"}', NOW(), NOW()),

(16, 'elderly', 'female', '1952-09-25', '北京市朝阳区康乐小区1号楼101', '一般', '轻度失能',
 '骨质疏松，曾骨折', '需要防跌倒辅助和钙质补充',
 '{"name":"郑强","phone":"13900000009","relation":"子女"}', NOW(), NOW()),

(17, 'elderly', 'male', '1946-05-18', '北京市朝阳区康乐小区2号楼202', '较差', '重度失能',
 '阿尔茨海默症早期，高血压', '需要24小时看护，防走失',
 '{"name":"冯丽华","phone":"13900000010","relation":"子女"}', NOW(), NOW()),

(18, 'elderly', 'female', '1953-12-02', '北京市朝阳区幸福小区7号楼707', '一般', '自理',
 '慢性支气管炎', '需要定期呼吸功能检查',
 '{"name":"陈刚","phone":"13900000011","relation":"子女"}', NOW(), NOW()),

(19, 'elderly', 'male', '1950-08-15', '北京市朝阳区康乐小区3号楼303', '良好', '自理',
 '前列腺增生', '需要定期泌尿科复查',
 '{"name":"杨芳","phone":"13900000012","relation":"子女"}', NOW(), NOW()),

(20, 'pregnant', 'female', '1992-05-10', '北京市朝阳区幸福小区8号楼808', '良好', NULL,
 '孕32周，妊娠期糖尿病', '需要饮食指导和血糖监测',
 '{"name":"黄先生","phone":"13900000013","relation":"配偶"}', NOW(), NOW()),

(22, 'elderly', 'male', '1947-12-08', '北京市朝阳区康乐小区4号楼404', '较差', '中度失能',
 '帕金森病，行动迟缓', '需要康复训练和日常照护',
 '{"name":"何美","phone":"13900000014","relation":"子女"}', NOW(), NOW()),

(23, 'elderly', 'female', '1956-01-15', '北京市朝阳区幸福小区9号楼909', '一般', '自理',
 '高血脂，轻度脂肪肝', '需要饮食管理和定期体检',
 '{"name":"马军","phone":"13900000015","relation":"子女"}', NOW(), NOW());

-- <<< END: database/seeds/modules/40_users_profiles.sql

-- >>> BEGIN: database/seeds/modules/50_requests_tasks.sql
-- 50_requests_tasks.sql
-- 服务请求、任务、任务历史模块初始化
-- 覆盖全状态(dispatched/claimed/completed/cancelled) × 全类型(meal/medical/care/other) × 30天时间跨度

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `service_requests` (
    `id`, `request_no`, `user_id`, `service_type`, `status`,
    `description`, `submit_location_lat`, `submit_location_lng`,
    `contact_name`, `contact_phone`, `address`, `appointment_time`,
    `urgency`, `station_id`, `rating`, `feedback`,
    `created_at`, `updated_at`
) VALUES
-- ===== 已完成 completed（8条）=====
(1, 'REQ-2026011501', 8, 'meal', 'completed',
 '需要午餐送餐服务，清淡饮食', 39.9200, 116.4500,
 '张大爷', '13800000008', '北京市朝阳区幸福小区1号楼101',
 '2026-01-15 12:00:00', 'normal', 1, 5, '服务很好，送餐准时',
 '2026-01-15 08:00:00', '2026-01-15 12:30:00'),

(2, 'REQ-2026011801', 9, 'medical', 'completed',
 '需要陪同去医院复查糖尿病', 39.9200, 116.4500,
 '李奶奶', '13800000009', '北京市朝阳区幸福小区2号楼202',
 '2026-01-18 09:00:00', 'urgent', 1, 4, '陪诊很耐心，但等候时间较长',
 '2026-01-17 10:00:00', '2026-01-18 14:00:00'),

(3, 'REQ-2026012001', 10, 'care', 'completed',
 '需要日常照护服务，协助洗漱穿衣', 39.9200, 116.4500,
 '王爷爷', '13800000010', '北京市朝阳区幸福小区3号楼303',
 '2026-01-20 08:00:00', 'normal', 1, 5, '护理人员非常专业',
 '2026-01-19 14:00:00', '2026-01-20 10:00:00'),

(4, 'REQ-2026012201', 14, 'medical', 'completed',
 '需要上门测量血压和血糖', 39.9200, 116.4500,
 '周阿姨', '13800000014', '北京市朝阳区幸福小区5号楼505',
 '2026-01-22 10:00:00', 'normal', 1, 5, '测量仔细，记录详细',
 '2026-01-21 09:00:00', '2026-01-22 11:00:00'),

(5, 'REQ-2026012501', 15, 'meal', 'completed',
 '需要晚餐送餐，少盐少油', 39.9200, 116.4500,
 '吴大爷', '13800000015', '北京市朝阳区幸福小区6号楼606',
 '2026-01-25 17:30:00', 'normal', 1, 4, '饭菜可口',
 '2026-01-25 10:00:00', '2026-01-25 18:00:00'),

(6, 'REQ-2026012801', 16, 'care', 'completed',
 '需要协助康复训练', 39.9100, 116.4600,
 '郑奶奶', '13800000016', '北京市朝阳区康乐小区1号楼101',
 '2026-01-28 14:00:00', 'normal', 2, 3, '服务态度好，但时间有点短',
 '2026-01-27 09:00:00', '2026-01-28 15:30:00'),

(7, 'REQ-2026020101', 17, 'care', 'completed',
 '需要24小时看护服务', 39.9100, 116.4600,
 '冯爷爷', '13800000017', '北京市朝阳区康乐小区2号楼202',
 '2026-02-01 08:00:00', 'urgent', 2, 5, '护工非常负责',
 '2026-01-31 10:00:00', '2026-02-01 20:00:00'),

(8, 'REQ-2026020301', 8, 'meal', 'completed',
 '需要午餐送餐，糖尿病餐', 39.9200, 116.4500,
 '张大爷', '13800000008', '北京市朝阳区幸福小区1号楼101',
 '2026-02-03 12:00:00', 'normal', 1, 5, '配餐合理',
 '2026-02-03 08:00:00', '2026-02-03 12:20:00'),

-- ===== 已认领 claimed（6条）=====
(9, 'REQ-2026020801', 9, 'medical', 'claimed',
 '需要陪同去医院做CT检查', 39.9200, 116.4500,
 '李奶奶', '13800000009', '北京市朝阳区幸福小区2号楼202',
 '2026-02-10 09:00:00', 'urgent', 1, 0, NULL,
 '2026-02-08 10:00:00', '2026-02-08 10:30:00'),

(10, 'REQ-2026020802', 18, 'medical', 'claimed',
 '需要上门做呼吸功能检查', 39.9200, 116.4500,
 '陈阿姨', '13800000018', '北京市朝阳区幸福小区7号楼707',
 '2026-02-10 14:00:00', 'normal', 1, 0, NULL,
 '2026-02-08 11:00:00', '2026-02-08 11:30:00'),

(11, 'REQ-2026020901', 22, 'care', 'claimed',
 '需要康复训练指导', 39.9100, 116.4600,
 '何大爷', '13800000022', '北京市朝阳区康乐小区4号楼404',
 '2026-02-11 09:00:00', 'normal', 2, 0, NULL,
 '2026-02-09 08:00:00', '2026-02-09 09:00:00'),

(12, 'REQ-2026021001', 19, 'meal', 'claimed',
 '需要午餐和晚餐送餐', 39.9100, 116.4600,
 '杨大爷', '13800000019', '北京市朝阳区康乐小区3号楼303',
 '2026-02-10 11:30:00', 'normal', 2, 0, NULL,
 '2026-02-10 07:00:00', '2026-02-10 08:00:00'),

(13, 'REQ-2026021101', 10, 'care', 'claimed',
 '需要协助翻身和按摩', 39.9200, 116.4500,
 '王爷爷', '13800000010', '北京市朝阳区幸福小区3号楼303',
 '2026-02-12 08:00:00', 'urgent', 1, 0, NULL,
 '2026-02-11 15:00:00', '2026-02-11 16:00:00'),

(14, 'REQ-2026021102', 23, 'medical', 'claimed',
 '需要上门抽血体检', 39.9200, 116.4500,
 '马奶奶', '13800000023', '北京市朝阳区幸福小区9号楼909',
 '2026-02-13 08:00:00', 'normal', 1, 0, NULL,
 '2026-02-11 09:00:00', '2026-02-11 10:00:00'),

-- ===== 待认领 dispatched（7条）=====
(15, 'REQ-2026021201', 14, 'meal', 'dispatched',
 '需要午餐送餐，软食', 39.9200, 116.4500,
 '周阿姨', '13800000014', '北京市朝阳区幸福小区5号楼505',
 '2026-02-13 12:00:00', 'normal', 1, 0, NULL,
 '2026-02-12 08:00:00', '2026-02-12 08:00:00'),

(16, 'REQ-2026021202', 15, 'medical', 'dispatched',
 '需要陪同去眼科复查', 39.9200, 116.4500,
 '吴大爷', '13800000015', '北京市朝阳区幸福小区6号楼606',
 '2026-02-14 09:00:00', 'normal', 1, 0, NULL,
 '2026-02-12 10:00:00', '2026-02-12 10:00:00'),

(17, 'REQ-2026021203', 17, 'care', 'dispatched',
 '需要全天看护，家属外出', 39.9100, 116.4600,
 '冯爷爷', '13800000017', '北京市朝阳区康乐小区2号楼202',
 '2026-02-14 08:00:00', 'urgent', 2, 0, NULL,
 '2026-02-12 14:00:00', '2026-02-12 14:00:00'),

(18, 'REQ-2026021204', 16, 'other', 'dispatched',
 '需要代购日用品和药品', 39.9100, 116.4600,
 '郑奶奶', '13800000016', '北京市朝阳区康乐小区1号楼101',
 '2026-02-13 15:00:00', 'normal', 2, 0, NULL,
 '2026-02-12 16:00:00', '2026-02-12 16:00:00'),

(19, 'REQ-2026021301', 8, 'meal', 'dispatched',
 '需要午餐送餐，低盐饮食', 39.9200, 116.4500,
 '张大爷', '13800000008', '北京市朝阳区幸福小区1号楼101',
 '2026-02-13 12:00:00', 'normal', 1, 0, NULL,
 '2026-02-13 07:00:00', '2026-02-13 07:00:00'),

(20, 'REQ-2026021302', 20, 'other', 'dispatched',
 '需要产前护理指导和营养咨询', 39.9200, 116.4500,
 '黄女士', '13800000020', '北京市朝阳区幸福小区8号楼808',
 '2026-02-14 10:00:00', 'normal', 1, 0, NULL,
 '2026-02-13 09:00:00', '2026-02-13 09:00:00'),

(21, 'REQ-2026021303', 12, 'care', 'dispatched',
 '需要专业康复护理', 39.9100, 116.4600,
 '赵先生', '13800000012', '北京市朝阳区康乐小区6号楼602',
 '2026-02-14 09:00:00', 'urgent', 2, 0, NULL,
 '2026-02-13 10:00:00', '2026-02-13 10:00:00'),

-- ===== 已取消 cancelled（4条）=====
(22, 'REQ-2026012301', 11, 'other', 'cancelled',
 '需要产前护理指导', 39.9100, 116.4600,
 '孙女士', '13800000011', '北京市朝阳区康乐小区5号楼501',
 '2026-01-24 10:00:00', 'normal', 2, 0, NULL,
 '2026-01-23 15:00:00', '2026-01-23 15:30:00'),

(23, 'REQ-2026013001', 19, 'meal', 'cancelled',
 '需要午餐送餐（已自行解决）', 39.9100, 116.4600,
 '杨大爷', '13800000019', '北京市朝阳区康乐小区3号楼303',
 '2026-01-30 12:00:00', 'normal', 2, 0, NULL,
 '2026-01-30 08:00:00', '2026-01-30 09:00:00'),

(24, 'REQ-2026020501', 15, 'medical', 'cancelled',
 '需要陪诊（改期）', 39.9200, 116.4500,
 '吴大爷', '13800000015', '北京市朝阳区幸福小区6号楼606',
 '2026-02-06 09:00:00', 'normal', 1, 0, NULL,
 '2026-02-05 10:00:00', '2026-02-05 14:00:00'),

(25, 'REQ-2026021103', 18, 'other', 'cancelled',
 '需要代购药品（家属已购买）', 39.9200, 116.4500,
 '陈阿姨', '13800000018', '北京市朝阳区幸福小区7号楼707',
 '2026-02-12 10:00:00', 'normal', 1, 0, NULL,
 '2026-02-11 08:00:00', '2026-02-11 12:00:00');

-- =====================================================
-- 任务分配（task_assignments）
-- =====================================================
INSERT INTO `task_assignments` (
    `id`, `request_id`, `station_id`, `staff_id`, `status`,
    `claimed_at`, `completed_at`, `rating`, `feedback`, `staff_notes`,
    `created_at`, `updated_at`
) VALUES
-- ===== 已完成 completed（8条）=====
(1, 1, 1, 4, 'completed',
 '2026-01-15 08:30:00', '2026-01-15 12:30:00', 5, '服务很好，送餐准时', '老人饮食偏清淡，已备注',
 '2026-01-15 08:05:00', '2026-01-15 12:30:00'),

(2, 2, 1, 5, 'completed',
 '2026-01-17 11:00:00', '2026-01-18 14:00:00', 4, '陪诊很耐心，但等候时间较长', '糖尿病复查，需空腹',
 '2026-01-17 10:05:00', '2026-01-18 14:00:00'),

(3, 3, 1, 4, 'completed',
 '2026-01-19 15:00:00', '2026-01-20 10:00:00', 5, '护理人员非常专业', '协助洗漱穿衣，老人行动不便',
 '2026-01-19 14:05:00', '2026-01-20 10:00:00'),

(4, 4, 1, 24, 'completed',
 '2026-01-21 10:00:00', '2026-01-22 11:00:00', 5, '测量仔细，记录详细', '血压偏高，建议复查',
 '2026-01-21 09:05:00', '2026-01-22 11:00:00'),

(5, 5, 1, 5, 'completed',
 '2026-01-25 11:00:00', '2026-01-25 18:00:00', 4, '饭菜可口', '少盐少油，已与厨房确认',
 '2026-01-25 10:05:00', '2026-01-25 18:00:00'),

(6, 6, 2, 6, 'completed',
 '2026-01-27 10:00:00', '2026-01-28 15:30:00', 3, '服务态度好，但时间有点短', '康复训练45分钟',
 '2026-01-27 09:05:00', '2026-01-28 15:30:00'),

(7, 7, 2, 7, 'completed',
 '2026-01-31 11:00:00', '2026-02-01 20:00:00', 5, '护工非常负责', '24小时看护，夜间巡查3次',
 '2026-01-31 10:05:00', '2026-02-01 20:00:00'),

(8, 8, 1, 4, 'completed',
 '2026-02-03 08:30:00', '2026-02-03 12:20:00', 5, '配餐合理', '糖尿病餐，低GI食材',
 '2026-02-03 08:05:00', '2026-02-03 12:20:00'),

-- ===== 已认领 claimed（6条）=====
(9, 9, 1, 5, 'claimed',
 '2026-02-08 10:30:00', NULL, 0, NULL, NULL,
 '2026-02-08 10:05:00', '2026-02-08 10:30:00'),

(10, 10, 1, 24, 'claimed',
 '2026-02-08 11:30:00', NULL, 0, NULL, NULL,
 '2026-02-08 11:05:00', '2026-02-08 11:30:00'),

(11, 11, 2, 6, 'claimed',
 '2026-02-09 09:00:00', NULL, 0, NULL, NULL,
 '2026-02-09 08:05:00', '2026-02-09 09:00:00'),

(12, 12, 2, 7, 'claimed',
 '2026-02-10 08:00:00', NULL, 0, NULL, NULL,
 '2026-02-10 07:05:00', '2026-02-10 08:00:00'),

(13, 13, 1, 4, 'claimed',
 '2026-02-11 16:00:00', NULL, 0, NULL, NULL,
 '2026-02-11 15:05:00', '2026-02-11 16:00:00'),

(14, 14, 1, 24, 'claimed',
 '2026-02-11 10:00:00', NULL, 0, NULL, NULL,
 '2026-02-11 09:05:00', '2026-02-11 10:00:00'),

-- ===== 待认领 dispatched（7条）=====
(15, 15, 1, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-12 08:05:00', '2026-02-12 08:05:00'),

(16, 16, 1, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-12 10:05:00', '2026-02-12 10:05:00'),

(17, 17, 2, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-12 14:05:00', '2026-02-12 14:05:00'),

(18, 18, 2, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-12 16:05:00', '2026-02-12 16:05:00'),

(19, 19, 1, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-13 07:05:00', '2026-02-13 07:05:00'),

(20, 20, 1, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-13 09:05:00', '2026-02-13 09:05:00'),

(21, 21, 2, NULL, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-13 10:05:00', '2026-02-13 10:05:00'),

-- ===== 已取消 cancelled（4条）=====
(22, 22, 2, NULL, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-01-23 15:05:00', '2026-01-23 15:30:00'),

(23, 23, 2, NULL, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-01-30 08:05:00', '2026-01-30 09:00:00'),

(24, 24, 1, NULL, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-05 10:05:00', '2026-02-05 14:00:00'),

(25, 25, 1, NULL, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-11 08:05:00', '2026-02-11 12:00:00');

-- =====================================================
-- 任务历史（task_histories）
-- 完整生命周期：dispatched → claimed → completed / cancelled
-- =====================================================
INSERT INTO `task_histories` (
    `id`, `task_id`, `request_id`, `action`, `operator_id`,
    `from_staff_id`, `to_staff_id`, `from_station_id`, `to_station_id`,
    `status_before`, `status_after`, `remark`, `created_at`, `updated_at`
) VALUES
-- ===== 已完成任务历史（每条3步：dispatch → claim → complete）=====
-- 任务1
(1, 1, 1, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-01-15 08:05:00', '2026-01-15 08:05:00'),
(2, 1, 1, 'claim', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-15 08:30:00', '2026-01-15 08:30:00'),
(3, 1, 1, 'complete', 4, NULL, NULL, NULL, NULL, 'claimed', 'completed', '送餐完成，老人确认签收', '2026-01-15 12:30:00', '2026-01-15 12:30:00'),
-- 任务2
(4, 2, 2, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '紧急派单', '2026-01-17 10:05:00', '2026-01-17 10:05:00'),
(5, 2, 2, 'claim', 5, NULL, 5, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-17 11:00:00', '2026-01-17 11:00:00'),
(6, 2, 2, 'complete', 5, NULL, NULL, NULL, NULL, 'claimed', 'completed', '陪诊完成，已取药', '2026-01-18 14:00:00', '2026-01-18 14:00:00'),
-- 任务3
(7, 3, 3, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-01-19 14:05:00', '2026-01-19 14:05:00'),
(8, 3, 3, 'claim', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-19 15:00:00', '2026-01-19 15:00:00'),
(9, 3, 3, 'complete', 4, NULL, NULL, NULL, NULL, 'claimed', 'completed', '照护服务完成', '2026-01-20 10:00:00', '2026-01-20 10:00:00'),
-- 任务4
(10, 4, 4, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-01-21 09:05:00', '2026-01-21 09:05:00'),
(11, 4, 4, 'claim', 24, NULL, 24, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-21 10:00:00', '2026-01-21 10:00:00'),
(12, 4, 4, 'complete', 24, NULL, NULL, NULL, NULL, 'claimed', 'completed', '血压血糖测量完成，已记录', '2026-01-22 11:00:00', '2026-01-22 11:00:00'),
-- 任务5
(13, 5, 5, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-01-25 10:05:00', '2026-01-25 10:05:00'),
(14, 5, 5, 'claim', 5, NULL, 5, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-25 11:00:00', '2026-01-25 11:00:00'),
(15, 5, 5, 'complete', 5, NULL, NULL, NULL, NULL, 'claimed', 'completed', '晚餐送达，老人满意', '2026-01-25 18:00:00', '2026-01-25 18:00:00'),
-- 任务6
(16, 6, 6, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-01-27 09:05:00', '2026-01-27 09:05:00'),
(17, 6, 6, 'claim', 6, NULL, 6, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-27 10:00:00', '2026-01-27 10:00:00'),
(18, 6, 6, 'complete', 6, NULL, NULL, NULL, NULL, 'claimed', 'completed', '康复训练完成', '2026-01-28 15:30:00', '2026-01-28 15:30:00'),
-- 任务7
(19, 7, 7, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '紧急派单', '2026-01-31 10:05:00', '2026-01-31 10:05:00'),
(20, 7, 7, 'claim', 7, NULL, 7, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-01-31 11:00:00', '2026-01-31 11:00:00'),
(21, 7, 7, 'complete', 7, NULL, NULL, NULL, NULL, 'claimed', 'completed', '24小时看护完成，一切正常', '2026-02-01 20:00:00', '2026-02-01 20:00:00'),
-- 任务8
(22, 8, 8, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-03 08:05:00', '2026-02-03 08:05:00'),
(23, 8, 8, 'claim', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-03 08:30:00', '2026-02-03 08:30:00'),
(24, 8, 8, 'complete', 4, NULL, NULL, NULL, NULL, 'claimed', 'completed', '糖尿病餐送达', '2026-02-03 12:20:00', '2026-02-03 12:20:00'),

-- ===== 已认领任务历史（每条2步：dispatch → claim）=====
-- 任务9
(25, 9, 9, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '紧急派单', '2026-02-08 10:05:00', '2026-02-08 10:05:00'),
(26, 9, 9, 'claim', 5, NULL, 5, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-08 10:30:00', '2026-02-08 10:30:00'),
-- 任务10
(27, 10, 10, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-08 11:05:00', '2026-02-08 11:05:00'),
(28, 10, 10, 'claim', 24, NULL, 24, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-08 11:30:00', '2026-02-08 11:30:00'),
-- 任务11
(29, 11, 11, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-09 08:05:00', '2026-02-09 08:05:00'),
(30, 11, 11, 'claim', 6, NULL, 6, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-09 09:00:00', '2026-02-09 09:00:00'),
-- 任务12
(31, 12, 12, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-10 07:05:00', '2026-02-10 07:05:00'),
(32, 12, 12, 'claim', 7, NULL, 7, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-10 08:00:00', '2026-02-10 08:00:00'),
-- 任务13
(33, 13, 13, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '紧急派单', '2026-02-11 15:05:00', '2026-02-11 15:05:00'),
(34, 13, 13, 'claim', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-11 16:00:00', '2026-02-11 16:00:00'),
-- 任务14
(35, 14, 14, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-11 09:05:00', '2026-02-11 09:05:00'),
(36, 14, 14, 'claim', 24, NULL, 24, NULL, NULL, 'dispatched', 'claimed', '员工主动认领', '2026-02-11 10:00:00', '2026-02-11 10:00:00'),

-- ===== 待认领任务历史（每条1步：dispatch）=====
(37, 15, 15, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-12 08:05:00', '2026-02-12 08:05:00'),
(38, 16, 16, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-12 10:05:00', '2026-02-12 10:05:00'),
(39, 17, 17, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '紧急派单', '2026-02-12 14:05:00', '2026-02-12 14:05:00'),
(40, 18, 18, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-12 16:05:00', '2026-02-12 16:05:00'),
(41, 19, 19, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-13 07:05:00', '2026-02-13 07:05:00'),
(42, 20, 20, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-13 09:05:00', '2026-02-13 09:05:00'),
(43, 21, 21, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-13 10:05:00', '2026-02-13 10:05:00'),

-- ===== 已取消任务历史（每条2步：dispatch → cancel）=====
-- 任务22
(44, 22, 22, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-01-23 15:05:00', '2026-01-23 15:05:00'),
(45, 22, 22, 'cancel', 11, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '用户主动取消', '2026-01-23 15:30:00', '2026-01-23 15:30:00'),
-- 任务23
(46, 23, 23, 'dispatch', 3, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-01-30 08:05:00', '2026-01-30 08:05:00'),
(47, 23, 23, 'cancel', 19, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '用户自行解决', '2026-01-30 09:00:00', '2026-01-30 09:00:00'),
-- 任务24
(48, 24, 24, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-05 10:05:00', '2026-02-05 10:05:00'),
(49, 24, 24, 'cancel', 15, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '用户改期', '2026-02-05 14:00:00', '2026-02-05 14:00:00'),
-- 任务25
(50, 25, 25, 'dispatch', 2, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-11 08:05:00', '2026-02-11 08:05:00'),
(51, 25, 25, 'cancel', 18, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '家属已购买药品', '2026-02-11 12:00:00', '2026-02-11 12:00:00');

-- <<< END: database/seeds/modules/50_requests_tasks.sql

-- >>> BEGIN: database/seeds/modules/60_content.sql
-- 60_content.sql
-- 内容模块初始化（Banner + News）
-- 覆盖多链接类型、多状态、多站点

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `banners` (`id`, `station_id`, `title`, `image_url`, `link_type`, `link_value`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '社区养老服务平台', 'https://via.placeholder.com/800x400/4A90E2/FFFFFF?text=社区养老服务', 'none', '', 100, 'active', NOW(3), NOW(3)),
(2, 0, '专业护理团队', 'https://via.placeholder.com/800x400/50C878/FFFFFF?text=专业护理', 'none', '', 90, 'active', NOW(3), NOW(3)),
(3, 0, '24小时服务热线', 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=24小时服务', 'none', '', 80, 'active', NOW(3), NOW(3)),
(4, 1, '幸福街道春季义诊', 'https://via.placeholder.com/800x400/9B59B6/FFFFFF?text=春季义诊', 'news', '3', 70, 'active', NOW(3), NOW(3)),
(5, 2, '康乐小区康复训练班', 'https://via.placeholder.com/800x400/E67E22/FFFFFF?text=康复训练', 'news', '6', 60, 'active', NOW(3), NOW(3)),
(6, 0, '元宵节送温暖活动', 'https://via.placeholder.com/800x400/1ABC9C/FFFFFF?text=元宵送温暖', 'url', 'https://example.com/lantern-festival', 50, 'active', NOW(3), NOW(3)),
(7, 0, '已下线的测试轮播图', 'https://via.placeholder.com/800x400/95A5A6/FFFFFF?text=已下线', 'none', '', 10, 'inactive', NOW(3), NOW(3));

INSERT INTO `news` (`id`, `title`, `summary`, `content`, `cover_url`, `type`, `status`, `station_id`, `author_id`, `publish_at`, `view_count`, `created_at`, `updated_at`) VALUES
-- 新闻（news）
(1, '社区养老服务中心正式启用',
 '首批养老服务站点完成部署并投入使用。',
 '经过数月筹备，社区养老服务中心于2026年1月正式启用。首批两个服务站点——幸福街道站和康乐小区站已完成部署，可为辖区内老年人、残障人士、孕产妇及儿童提供送餐、陪诊、照护等多项服务。',
 'https://via.placeholder.com/600x300/4A90E2/FFFFFF?text=News+1', 'news', 'published', 0, 1,
 '2026-01-10 09:00:00', 320, '2026-01-10 08:00:00', '2026-01-10 09:00:00'),

(2, '幸福街道站首月服务报告出炉',
 '幸福街道站开站首月累计服务超过50人次。',
 '幸福街道服务站自1月15日正式运营以来，已累计完成送餐服务20次、陪诊服务12次、日常照护服务18次，服务满意度达4.6分（满分5分）。',
 'https://via.placeholder.com/600x300/2ECC71/FFFFFF?text=News+2', 'news', 'published', 1, 2,
 '2026-02-05 10:00:00', 185, '2026-02-05 09:00:00', '2026-02-05 10:00:00'),

(3, '智慧养老平台升级：新增在线预约功能',
 '平台新增在线预约功能，居民可通过手机端直接预约服务。',
 '为提升服务便捷性，平台于近日完成功能升级，新增在线预约模块。居民可通过微信小程序或PWA应用直接提交服务请求，系统将自动匹配最近的服务站点并派单。',
 'https://via.placeholder.com/600x300/3498DB/FFFFFF?text=News+3', 'news', 'published', 0, 1,
 '2026-02-10 14:00:00', 95, '2026-02-10 13:00:00', '2026-02-10 14:00:00'),

(4, '康乐小区站志愿者招募中',
 '康乐小区服务站面向社会招募志愿者，欢迎有爱心的居民加入。',
 '为扩大服务覆盖面，康乐小区服务站现面向社会公开招募志愿者。志愿者将参与送餐配送、陪同就医、日常探访等服务。有意者请联系站点负责人王站长。',
 'https://via.placeholder.com/600x300/E74C3C/FFFFFF?text=News+4', 'news', 'published', 2, 3,
 '2026-02-12 09:00:00', 42, '2026-02-12 08:00:00', '2026-02-12 09:00:00'),

-- 公告（notice）
(5, '春季义诊活动公告',
 '本周六开展春季社区义诊，请有需要的居民提前预约。',
 '定于2026年2月15日（周六）上午9:00-12:00，在幸福街道服务站开展春季社区义诊活动。届时将有内科、骨科、眼科等专家坐诊，请有需要的居民提前到站点登记预约。',
 'https://via.placeholder.com/600x300/50C878/FFFFFF?text=Notice+1', 'notice', 'published', 1, 2,
 '2026-02-11 10:00:00', 210, '2026-02-11 09:00:00', '2026-02-11 10:00:00'),

(6, '春节期间服务时间调整通知',
 '春节期间（1月28日-2月4日）服务时间有所调整。',
 '各位居民：春节期间（1月28日至2月4日），各服务站点将调整服务时间为每日9:00-16:00。紧急服务请拨打24小时热线。祝大家新春快乐！',
 'https://via.placeholder.com/600x300/F39C12/FFFFFF?text=Notice+2', 'notice', 'published', 0, 1,
 '2026-01-25 15:00:00', 450, '2026-01-25 14:00:00', '2026-01-25 15:00:00'),

(7, '系统维护公告',
 '2月16日凌晨2:00-5:00进行系统维护，届时服务将暂停。',
 '为提升系统稳定性和安全性，计划于2026年2月16日凌晨2:00至5:00进行系统维护升级。维护期间平台将暂停服务，请提前做好安排。',
 'https://via.placeholder.com/600x300/95A5A6/FFFFFF?text=Notice+3', 'notice', 'published', 0, 1,
 '2026-02-13 10:00:00', 15, '2026-02-13 09:00:00', '2026-02-13 10:00:00'),

-- 活动（activity）
(8, '元宵节送温暖——独居老人慰问活动',
 '元宵佳节，为辖区独居老人送去汤圆和祝福。',
 '2026年2月12日元宵节当天，两个服务站点联合开展"元宵送温暖"活动，为辖区内30余位独居老人送去汤圆、水果和节日祝福。活动受到老人们的热烈欢迎。',
 'https://via.placeholder.com/600x300/1ABC9C/FFFFFF?text=Activity+1', 'activity', 'published', 0, 1,
 '2026-02-12 16:00:00', 128, '2026-02-12 15:00:00', '2026-02-12 16:00:00'),

(9, '老年人智能手机使用培训班',
 '免费教老年人使用智能手机，包括微信、健康码、在线挂号等。',
 '为帮助老年人跨越"数字鸿沟"，幸福街道服务站将于每周三下午14:00-16:00开设智能手机使用培训班。课程内容包括微信基本操作、健康码使用、在线挂号预约等。',
 'https://via.placeholder.com/600x300/9B59B6/FFFFFF?text=Activity+2', 'activity', 'published', 1, 2,
 '2026-02-08 09:00:00', 76, '2026-02-08 08:00:00', '2026-02-08 09:00:00'),

-- 草稿（draft）
(10, '三月份健康讲座预告',
 '三月份将开展系列健康讲座，敬请期待。',
 '三月份计划开展以下健康讲座：3月5日"春季养生与饮食调理"、3月12日"高血压日常管理"、3月19日"糖尿病饮食指南"、3月26日"老年人防跌倒指南"。',
 'https://via.placeholder.com/600x300/BDC3C7/FFFFFF?text=Draft', 'notice', 'draft', 0, 1,
 NULL, 0, '2026-02-13 11:00:00', '2026-02-13 11:00:00'),

-- 已下线（offline）
(11, '2025年度工作总结',
 '2025年度社区养老服务工作总结报告。',
 '2025年度，社区养老服务项目完成前期调研、系统开发、站点建设等工作，为2026年正式运营奠定了基础。',
 'https://via.placeholder.com/600x300/7F8C8D/FFFFFF?text=Offline', 'news', 'offline', 0, 1,
 '2025-12-30 10:00:00', 560, '2025-12-30 09:00:00', '2026-01-15 10:00:00');

-- <<< END: database/seeds/modules/60_content.sql

-- >>> BEGIN: database/seeds/modules/70_notifications.sql
-- 70_notifications.sql
-- 通知模块初始化
-- 覆盖 system/task/request 类型 × 已读/未读 × 多渠道 × 多用户

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `notifications` (
    `id`, `user_id`, `title`, `body`, `type`, `related_id`, `related_type`,
    `channel`, `send_status`, `is_read`, `retry_count`, `created_at`, `updated_at`
) VALUES
-- ===== 系统通知 system（10条，7已读3未读）=====
(1, 2, '系统通知', '欢迎使用社区养老服务平台管理端。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-10 10:00:00'),
(2, 3, '系统通知', '欢迎使用社区养老服务平台管理端。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-10 11:00:00'),
(3, 4, '系统通知', '欢迎使用社区养老服务平台，您已被分配为幸福街道站员工。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-10 14:00:00'),
(4, 5, '系统通知', '欢迎使用社区养老服务平台，您已被分配为幸福街道站员工。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-11 08:00:00'),
(5, 6, '系统通知', '欢迎使用社区养老服务平台，您已被分配为康乐小区站员工。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-10 15:00:00'),
(6, 7, '系统通知', '欢迎使用社区养老服务平台，您已被分配为康乐小区站员工。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-01-10 09:00:00', '2026-01-10 16:00:00'),
(7, 24, '系统通知', '欢迎加入幸福街道服务站团队。', 'system', 0, 'system', 'in_app', 'sent', 1, 0, '2026-02-01 09:00:00', '2026-02-01 10:00:00'),
(8, 25, '系统通知', '欢迎加入康乐小区服务站团队。', 'system', 0, 'system', 'in_app', 'sent', 0, 0, '2026-02-01 09:00:00', '2026-02-01 09:00:00'),
(9, 2, '系统维护通知', '2月16日凌晨2:00-5:00将进行系统维护，届时服务暂停。', 'system', 0, 'system', 'in_app', 'sent', 0, 0, '2026-02-13 10:00:00', '2026-02-13 10:00:00'),
(10, 3, '系统维护通知', '2月16日凌晨2:00-5:00将进行系统维护，届时服务暂停。', 'system', 0, 'system', 'in_app', 'sent', 0, 0, '2026-02-13 10:00:00', '2026-02-13 10:00:00'),

-- ===== 任务通知 task（12条，8已读4未读）=====
-- 已完成任务的通知
(11, 4, '新任务分配', '您有一个新的送餐任务（REQ-2026011501），请及时处理。', 'task', 1, 'task', 'in_app', 'sent', 1, 0, '2026-01-15 08:05:00', '2026-01-15 08:10:00'),
(12, 5, '新任务分配', '您有一个新的陪诊任务（REQ-2026011801），请及时处理。', 'task', 2, 'task', 'in_app', 'sent', 1, 0, '2026-01-17 10:05:00', '2026-01-17 10:30:00'),
(13, 4, '新任务分配', '您有一个新的照护任务（REQ-2026012001），请及时处理。', 'task', 3, 'task', 'in_app', 'sent', 1, 0, '2026-01-19 14:05:00', '2026-01-19 14:30:00'),
(14, 24, '新任务分配', '您有一个新的医疗任务（REQ-2026012201），请及时处理。', 'task', 4, 'task', 'in_app', 'sent', 1, 0, '2026-01-21 09:05:00', '2026-01-21 09:30:00'),
(15, 6, '新任务分配', '您有一个新的康复训练任务（REQ-2026012801），请及时处理。', 'task', 6, 'task', 'in_app', 'sent', 1, 0, '2026-01-27 09:05:00', '2026-01-27 09:30:00'),
(16, 7, '新任务分配', '您有一个紧急看护任务（REQ-2026020101），请尽快处理。', 'task', 7, 'task', 'in_app', 'sent', 1, 0, '2026-01-31 10:05:00', '2026-01-31 10:20:00'),
(17, 4, '新任务分配', '您有一个新的送餐任务（REQ-2026020301），请及时处理。', 'task', 8, 'task', 'in_app', 'sent', 1, 0, '2026-02-03 08:05:00', '2026-02-03 08:15:00'),
(18, 5, '任务提醒', '您认领的陪诊任务（REQ-2026020801）预约时间为明天上午9点，请做好准备。', 'task', 9, 'task', 'in_app', 'sent', 1, 0, '2026-02-09 18:00:00', '2026-02-09 19:00:00'),
-- 当前进行中任务的通知（未读）
(19, 6, '新任务分配', '您有一个新的康复训练任务（REQ-2026020901），请及时处理。', 'task', 11, 'task', 'in_app', 'sent', 0, 0, '2026-02-09 08:05:00', '2026-02-09 08:05:00'),
(20, 7, '新任务分配', '您有一个新的送餐任务（REQ-2026021001），请及时处理。', 'task', 12, 'task', 'in_app', 'sent', 0, 0, '2026-02-10 07:05:00', '2026-02-10 07:05:00'),
(21, 4, '新任务分配', '您有一个紧急照护任务（REQ-2026021101），请尽快处理。', 'task', 13, 'task', 'in_app', 'sent', 0, 0, '2026-02-11 15:05:00', '2026-02-11 15:05:00'),
(22, 24, '新任务分配', '您有一个新的体检任务（REQ-2026021102），请及时处理。', 'task', 14, 'task', 'in_app', 'sent', 0, 0, '2026-02-11 09:05:00', '2026-02-11 09:05:00'),

-- ===== 请求通知 request（10条，7已读3未读）=====
-- 已完成请求的用户通知
(23, 8, '服务完成', '您的送餐服务请求（REQ-2026011501）已完成，请对服务进行评价。', 'request', 1, 'request', 'in_app', 'sent', 1, 0, '2026-01-15 12:30:00', '2026-01-15 13:00:00'),
(24, 9, '服务完成', '您的陪诊服务请求（REQ-2026011801）已完成，请对服务进行评价。', 'request', 2, 'request', 'in_app', 'sent', 1, 0, '2026-01-18 14:00:00', '2026-01-18 15:00:00'),
(25, 10, '服务完成', '您的照护服务请求（REQ-2026012001）已完成，请对服务进行评价。', 'request', 3, 'request', 'in_app', 'sent', 1, 0, '2026-01-20 10:00:00', '2026-01-20 11:00:00'),
(26, 14, '服务完成', '您的医疗服务请求（REQ-2026012201）已完成，请对服务进行评价。', 'request', 4, 'request', 'in_app', 'sent', 1, 0, '2026-01-22 11:00:00', '2026-01-22 12:00:00'),
(27, 15, '服务完成', '您的送餐服务请求（REQ-2026012501）已完成，请对服务进行评价。', 'request', 5, 'request', 'in_app', 'sent', 1, 0, '2026-01-25 18:00:00', '2026-01-25 19:00:00'),
(28, 8, '服务完成', '您的送餐服务请求（REQ-2026020301）已完成，请对服务进行评价。', 'request', 8, 'request', 'in_app', 'sent', 1, 0, '2026-02-03 12:20:00', '2026-02-03 13:00:00'),
(29, 16, '服务完成', '您的康复训练请求（REQ-2026012801）已完成，请对服务进行评价。', 'request', 6, 'request', 'in_app', 'sent', 1, 0, '2026-01-28 15:30:00', '2026-01-28 16:00:00'),
-- 进行中请求的用户通知（未读）
(30, 9, '服务已认领', '您的陪诊请求（REQ-2026020801）已被工作人员认领，请保持电话畅通。', 'request', 9, 'request', 'in_app', 'sent', 0, 0, '2026-02-08 10:30:00', '2026-02-08 10:30:00'),
(31, 10, '服务已认领', '您的照护请求（REQ-2026021101）已被工作人员认领，请保持电话畅通。', 'request', 13, 'request', 'in_app', 'sent', 0, 0, '2026-02-11 16:00:00', '2026-02-11 16:00:00'),
(32, 14, '服务已派单', '您的送餐请求（REQ-2026021201）已派单至幸福街道站，请耐心等待。', 'request', 15, 'request', 'in_app', 'sent', 0, 0, '2026-02-12 08:05:00', '2026-02-12 08:05:00');

-- <<< END: database/seeds/modules/70_notifications.sql

