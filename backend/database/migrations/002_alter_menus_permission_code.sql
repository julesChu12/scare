-- 002_alter_menus_permission_code.sql
-- 将 menus 表的 permission 字段改名为 permission_code，统一权限码格式

-- 重命名字段
ALTER TABLE `menus` CHANGE COLUMN `permission` `permission_code` VARCHAR(100) DEFAULT '' COMMENT '权限码，格式: module:resource:action';

-- 更新索引名称（先删除旧索引，再创建新索引）
ALTER TABLE `menus` DROP INDEX `idx_permission`;
ALTER TABLE `menus` ADD INDEX `idx_permission_code` (`permission_code`);

-- 更新现有菜单的权限码为新格式
UPDATE `menus` SET `permission_code` = 'service:request:list' WHERE `permission_code` = 'service:request:list';
UPDATE `menus` SET `permission_code` = 'service:request:detail' WHERE `permission_code` = 'service:request:detail';
UPDATE `menus` SET `permission_code` = 'service:task:pool' WHERE `permission_code` = 'service:task:pool';
UPDATE `menus` SET `permission_code` = 'service:task:my' WHERE `permission_code` = 'service:task:my';
UPDATE `menus` SET `permission_code` = 'station:list' WHERE `permission_code` = 'station:list';
UPDATE `menus` SET `permission_code` = 'station:zone:list' WHERE `permission_code` = 'station:zone:list';
UPDATE `menus` SET `permission_code` = 'system:user:list' WHERE `permission_code` = 'system:user:list';
UPDATE `menus` SET `permission_code` = 'system:role:list' WHERE `permission_code` = 'system:role:list';
UPDATE `menus` SET `permission_code` = 'system:menu:list' WHERE `permission_code` = 'system:menu:list';
UPDATE `menus` SET `permission_code` = 'content:banner:list' WHERE `permission_code` = 'content:banner:list';
UPDATE `menus` SET `permission_code` = 'content:news:list' WHERE `permission_code` = 'content:news:list';
