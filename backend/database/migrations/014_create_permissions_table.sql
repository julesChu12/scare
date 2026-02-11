-- 014_create_permissions_table.sql
-- 创建权限表

CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `code` VARCHAR(100) NOT NULL COMMENT '权限码，格式: module:resource:action',
    `name` VARCHAR(100) NOT NULL COMMENT '权限名称',
    `description` VARCHAR(500) DEFAULT '' COMMENT '权限描述',
    `module` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '所属模块',
    `type` ENUM('menu', 'button', 'resource') NOT NULL DEFAULT 'resource' COMMENT '权限类型: menu-菜单 button-按钮 resource-资源',
    `parent_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '父权限ID，0表示顶级',
    `api_path` VARCHAR(200) DEFAULT '' COMMENT 'API路径，支持通配符*',
    `api_method` VARCHAR(20) DEFAULT '' COMMENT 'HTTP方法: GET/POST/PUT/DELETE',
    `is_public` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否公共权限，1-是 0-否',
    `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_code` (`code`, `deleted_at`),
    INDEX `idx_parent_id` (`parent_id`),
    INDEX `idx_type` (`type`),
    INDEX `idx_module` (`module`),
    INDEX `idx_api` (`api_path`, `api_method`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';
