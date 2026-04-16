-- 为已经在 schema.sql 中前置包含这些字段的环境提供幂等保护。
-- 这样 fresh init 后再执行增量迁移时，不会因为重复 ADD COLUMN 失败。

SET @col_exists = (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND column_name = 'service_location_lat'
);
SET @sql = IF(
  @col_exists = 0,
  'ALTER TABLE `service_requests` ADD COLUMN `service_location_lat` DECIMAL(10,7) DEFAULT NULL AFTER `submit_location_lng`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @col_exists = (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND column_name = 'service_location_lng'
);
SET @sql = IF(
  @col_exists = 0,
  'ALTER TABLE `service_requests` ADD COLUMN `service_location_lng` DECIMAL(10,7) DEFAULT NULL AFTER `service_location_lat`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @col_exists = (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND column_name = 'source_station_id'
);
SET @sql = IF(
  @col_exists = 0,
  'ALTER TABLE `service_requests` ADD COLUMN `source_station_id` BIGINT DEFAULT NULL AFTER `urgency`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @col_exists = (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND column_name = 'dispatch_basis'
);
SET @sql = IF(
  @col_exists = 0,
  'ALTER TABLE `service_requests` ADD COLUMN `dispatch_basis` VARCHAR(50) DEFAULT NULL AFTER `station_id`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @col_exists = (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND column_name = 'needs_manual_verify'
);
SET @sql = IF(
  @col_exists = 0,
  'ALTER TABLE `service_requests` ADD COLUMN `needs_manual_verify` TINYINT(1) DEFAULT 0 AFTER `dispatch_basis`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @idx_exists = (
  SELECT COUNT(*)
  FROM information_schema.statistics
  WHERE table_schema = DATABASE()
    AND table_name = 'service_requests'
    AND index_name = 'idx_service_requests_source_station_id'
);
SET @sql = IF(
  @idx_exists = 0,
  'ALTER TABLE `service_requests` ADD KEY `idx_service_requests_source_station_id` (`source_station_id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
