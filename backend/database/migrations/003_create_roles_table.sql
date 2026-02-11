-- 003_create_roles_table.sql
-- 创建角色表

CREATE TABLE IF NOT EXISTS `roles` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `code` VARCHAR(50) NOT NULL COMMENT '角色编码，如 admin, staff',
    `name` VARCHAR(100) NOT NULL COMMENT '角色名称',
    `description` VARCHAR(500) DEFAULT '' COMMENT '角色描述',
    `is_system` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否系统内置角色，1-是 0-否',
    `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态 active-启用 inactive-禁用',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_code` (`code`, `deleted_at`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 插入系统内置角色
INSERT INTO `roles` (`code`, `name`, `description`, `is_system`, `sort`) VALUES
('admin', '系统管理员', '拥有系统全部权限', 1, 1),
('station_manager', '站点管理员', '管理所属站点的服务和人员', 1, 2),
('staff', '工作人员', '执行服务任务的一线人员', 1, 3),
('elderly', '老年人', 'C端老年用户', 1, 4),
('family', '家属', 'C端家属用户', 1, 5);
