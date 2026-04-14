-- sCare 数据库迁移脚本
-- 命名格式: YYYYMMDDHHMMSS_description.sql
-- 执行顺序: 按文件名 ASCII 排序

-- 使用方式:
-- docker exec scare_mysql mysql -u scare_user -pscare_pass scare_db --default-character-set=utf8mb4 < migration.sql

-- ============================================
-- 迁移记录表（首次运行需创建）
-- ============================================
CREATE TABLE IF NOT EXISTS `_migrations` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `applied_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `checksum` VARCHAR(64) DEFAULT NULL,
    INDEX `idx_migrations_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================
-- 迁移脚本模板
-- ============================================
-- -- migrate: 2025-01-01 add_user_status
-- -- 描述: 添加用户状态字段
--
-- INSERT INTO `_migrations` (`name`, `checksum`)
-- SELECT '20250101120000_add_user_status.sql', MD5(LOAD_FILE('20250101120000_add_user_status.sql'))
-- WHERE NOT EXISTS (SELECT 1 FROM `_migrations` WHERE name = '20250101120000_add_user_status.sql');
--
-- -- 你的 SQL 变更写在这里
-- ALTER TABLE `users` ADD COLUMN `status` VARCHAR(20) DEFAULT 'active' COMMENT '用户状态';
