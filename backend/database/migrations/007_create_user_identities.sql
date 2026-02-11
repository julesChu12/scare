-- 007_create_user_identities.sql
-- 创建 user_identities 表，统一管理用户的多身份（B端角色 + C端用户类型）
-- 执行时间: 2026-02-05
-- 说明：将 B端角色和 C端用户类型统一到一张表，支持一人多身份

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- =====================================================
-- 1. 创建 user_identities 表
-- =====================================================
CREATE TABLE IF NOT EXISTS `user_identities` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `identity_type` VARCHAR(20) NOT NULL COMMENT '身份类型: B端(admin/station_manager/staff) C端(elderly/family)',
  `is_primary` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否主身份(每用户有且仅有一个)',
  `station_id` BIGINT DEFAULT NULL COMMENT '所属服务站ID(仅B端身份)',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '身份状态(active/inactive)',
  `granted_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '授予时间',
  `granted_by` BIGINT DEFAULT NULL COMMENT '授予人ID',
  `revoked_at` DATETIME DEFAULT NULL COMMENT '撤销时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  -- 唯一约束：同一用户同一身份类型只能有一条激活记录
  UNIQUE KEY `uk_user_identity` (`user_id`, `identity_type`, `deleted_at`),
  -- 索引：按用户ID查询
  KEY `idx_ui_user_id` (`user_id`),
  -- 索引：按身份类型查询
  KEY `idx_ui_identity_type` (`identity_type`),
  -- 索引：按状态查询
  KEY `idx_ui_status` (`status`),
  -- 索引：按站点查询
  KEY `idx_ui_station_id` (`station_id`),
  -- 索引：软删除
  KEY `idx_ui_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户身份表';

-- =====================================================
-- 2. 删除 users 表的 role 和 primary_role 字段
-- =====================================================
-- 检查并删除 role 字段
SET @db_name = DATABASE();
SET @table_name = 'users';

-- 删除 role 字段的索引（如果存在）
SET @index_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND INDEX_NAME = 'idx_users_role'
);
SET @drop_index_sql = IF(@index_exists > 0, 'ALTER TABLE `users` DROP INDEX `idx_users_role`', 'SELECT 1');
PREPARE stmt FROM @drop_index_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 删除 primary_role 字段的索引（如果存在）
SET @index_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND INDEX_NAME = 'idx_users_primary_role'
);
SET @drop_index_sql = IF(@index_exists > 0, 'ALTER TABLE `users` DROP INDEX `idx_users_primary_role`', 'SELECT 1');
PREPARE stmt FROM @drop_index_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 删除 role 字段（如果存在）
SET @column_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND COLUMN_NAME = 'role'
);
SET @drop_column_sql = IF(@column_exists > 0, 'ALTER TABLE `users` DROP COLUMN `role`', 'SELECT 1');
PREPARE stmt FROM @drop_column_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 删除 primary_role 字段（如果存在）
SET @column_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND COLUMN_NAME = 'primary_role'
);
SET @drop_column_sql = IF(@column_exists > 0, 'ALTER TABLE `users` DROP COLUMN `primary_role`', 'SELECT 1');
PREPARE stmt FROM @drop_column_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- =====================================================
-- 3. 清理 roles 表中的 C端角色
-- =====================================================
-- 先删除 role_permissions 中的关联记录
DELETE FROM `role_permissions` WHERE role_id IN (SELECT id FROM `roles` WHERE code IN ('elderly', 'family'));
-- 删除 C端角色
DELETE FROM `roles` WHERE code IN ('elderly', 'family');

-- =====================================================
-- 4. 删除旧的 user_roles 表（数据将通过 SEED 重建）
-- =====================================================
DROP TABLE IF EXISTS `user_roles`;
