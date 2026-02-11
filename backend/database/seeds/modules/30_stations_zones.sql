-- 30_stations_zones.sql
-- 站点与围栏模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `service_stations` (`id`, `name`, `code`, `address`, `phone`, `latitude`, `longitude`, `status`, `created_at`, `updated_at`) VALUES
(1, '朝阳区幸福街道服务站', 'CYXF001', '北京市朝阳区幸福大街100号', '010-12345678', 39.9200, 116.4500, 'active', NOW(), NOW()),
(2, '朝阳区康乐街道服务站', 'CYKL002', '北京市朝阳区康乐大街200号', '010-87654321', 39.9100, 116.4600, 'active', NOW(), NOW());

-- =====================================================
-- 2.1 初始化服务区域/地理围栏（service_zones）
-- =====================================================
INSERT INTO `service_zones` (`id`, `station_id`, `name`, `points`, `priority`, `status`, `created_at`, `updated_at`) VALUES
-- 站点1服务区域（幸福街道）
(1, 1, '幸福街道服务区',
 '[[116.44,39.91],[116.46,39.91],[116.46,39.93],[116.44,39.93],[116.44,39.91]]',
 1, 'active', NOW(), NOW()),
-- 站点2服务区域（康乐街道）
(2, 2, '康乐街道服务区',
 '[[116.45,39.90],[116.47,39.90],[116.47,39.92],[116.45,39.92],[116.45,39.90]]',
 1, 'active', NOW(), NOW());
