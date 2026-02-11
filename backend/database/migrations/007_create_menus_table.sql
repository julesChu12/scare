-- 007_create_menus_table.sql
-- 创建菜单表，用于动态菜单系统

CREATE TABLE IF NOT EXISTS `menus` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示顶级菜单',
    `name` VARCHAR(50) NOT NULL COMMENT '菜单名称',
    `path` VARCHAR(200) DEFAULT '' COMMENT '路由路径',
    `component` VARCHAR(200) DEFAULT '' COMMENT '前端组件路径',
    `icon` VARCHAR(50) DEFAULT '' COMMENT '图标名称',
    `permission` VARCHAR(100) DEFAULT '' COMMENT '所需权限标识',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值，数字越小越靠前',
    `hidden` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏 0-显示 1-隐藏',
    `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_permission` (`permission`),
    INDEX `idx_status` (`status`),
    INDEX `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';
