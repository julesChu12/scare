-- menus_seed.sql (旧版，建议使用 007_seed_menus.sql)
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `permission`, `sort`, `hidden`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '任务管理', '/tasks', 'Layout', 'List', 'task:list', 1, 0, 'active', NOW(), NOW()),
(2, 1, '任务池', '/tasks/pool', 'TaskPool', 'List', 'task:pool', 1, 0, 'active', NOW(), NOW()),
(3, 1, '我的任务', '/tasks/my', 'MyTasks', 'User', 'task:my', 2, 0, 'active', NOW(), NOW()),
(4, 0, '系统管理', '/admin', 'Layout', 'Setting', 'system:manage', 9, 0, 'active', NOW(), NOW()),
(5, 4, '用户管理', '/admin/users', 'UserManagement', 'UserFilled', 'user:list', 1, 0, 'active', NOW(), NOW()),
(6, 4, '权限管理', '/admin/permissions', 'RolePermission', 'Lock', 'permission:manage', 2, 0, 'active', NOW(), NOW()),
(7, 4, '站点管理', '/admin/stations', 'StationManagement', 'OfficeBuilding', 'station:manage', 3, 0, 'active', NOW(), NOW()),
(8, 4, '围栏管理', '/admin/zones', 'ZoneManagement', 'MapLocation', 'zone:manage', 4, 0, 'active', NOW(), NOW());
