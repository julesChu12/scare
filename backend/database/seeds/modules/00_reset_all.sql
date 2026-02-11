-- 00_reset_all.sql
-- 清空业务表（用于完整重建测试数据）

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE role_permissions;
TRUNCATE TABLE permissions;
TRUNCATE TABLE roles;
TRUNCATE TABLE menus;
TRUNCATE TABLE task_histories;
TRUNCATE TABLE task_assignments;
TRUNCATE TABLE service_requests;
TRUNCATE TABLE service_zones;
TRUNCATE TABLE service_stations;
TRUNCATE TABLE notifications;
TRUNCATE TABLE customer_profiles;
TRUNCATE TABLE user_identities;
TRUNCATE TABLE news;
TRUNCATE TABLE banners;
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;
