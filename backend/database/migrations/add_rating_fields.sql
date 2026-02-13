-- Add rating fields to service_requests table
-- 使用 PREPARE 语句实现幂等 ADD COLUMN（MySQL 兼容）

SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'service_requests' AND column_name = 'rating');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE service_requests ADD COLUMN rating INT DEFAULT 0', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @col_exists = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'service_requests' AND column_name = 'feedback');
SET @sql = IF(@col_exists = 0, 'ALTER TABLE service_requests ADD COLUMN feedback TEXT', 'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
