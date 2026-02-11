-- 008_add_users_id_card_hmac.sql
-- 为 users 表增加 id_card_hmac 字段，用于存储身份证号的 HMAC 摘要
-- 执行时间: 2026-02-10

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

SET @db_name = DATABASE();
SET @table_name = 'users';

-- 增加 id_card_hmac 字段（如果不存在）
SET @column_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.COLUMNS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND COLUMN_NAME = 'id_card_hmac'
);
SET @add_column_sql = IF(
    @column_exists = 0,
    'ALTER TABLE `users` ADD COLUMN `id_card_hmac` VARCHAR(64) DEFAULT NULL COMMENT ''身份证号HMAC摘要'' AFTER `id_card`',
    'SELECT 1'
);
PREPARE stmt FROM @add_column_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 增加索引（如果不存在）
SET @index_exists = (
    SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS
    WHERE TABLE_SCHEMA = @db_name AND TABLE_NAME = @table_name AND INDEX_NAME = 'idx_users_id_card_hmac'
);
SET @add_index_sql = IF(
    @index_exists = 0,
    'ALTER TABLE `users` ADD INDEX `idx_users_id_card_hmac` (`id_card_hmac`)',
    'SELECT 1'
);
PREPARE stmt FROM @add_index_sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
