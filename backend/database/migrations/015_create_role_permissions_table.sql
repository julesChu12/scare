-- 015_create_role_permissions_table.sql
-- 创建角色权限关联表

CREATE TABLE IF NOT EXISTS `role_permissions` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `role_id` BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    `permission_id` BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_perm` (`role_id`, `permission_id`),
    INDEX `idx_role_id` (`role_id`),
    INDEX `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';
