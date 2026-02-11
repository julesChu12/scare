-- 016_seed_permissions.sql
-- 权限初始数据

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- =====================================================
-- 公共权限模块
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
('system:user:roles', '分配角色', '分配用户角色', 'system', 'button', 34, '/api/v1/b/users/*/roles', 'PUT', 0, 5),

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
SELECT 2, id FROM `permissions` WHERE `code` IN (
    'dashboard',
    'service', 'service:request', 'service:request:list', 'service:request:detail',
    'service:task', 'service:task:pool', 'service:task:my', 'service:task:claim', 'service:task:complete',
    'station', 'station:list', 'station:list:view', 'station:list:detail',
    'station:zone', 'station:zone:list', 'station:zone:create', 'station:zone:update', 'station:zone:delete',
    'system', 'system:user', 'system:user:list', 'system:user:detail', 'system:user:create', 'system:user:update', 'system:user:roles'
);

-- Staff 角色权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 3, id FROM `permissions` WHERE `code` IN (
    'dashboard',
    'service', 'service:request', 'service:request:list', 'service:request:detail',
    'service:task', 'service:task:pool', 'service:task:my', 'service:task:claim', 'service:task:complete',
    'station', 'station:list', 'station:list:view',
    'station:zone', 'station:zone:list'
);
