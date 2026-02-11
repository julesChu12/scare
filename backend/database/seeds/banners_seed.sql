-- Banner 初始化数据
-- 用于演示和测试
-- 执行方式: docker exec -i scare_mysql mysql -u scare_user -pscare_pass scare_db --default-character-set=utf8mb4 < banners_seed.sql

-- 设置字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 全局 Banner（station_id = 0，所有站点可见）
INSERT INTO banners (station_id, title, image_url, link_type, link_value, sort, status, created_at, updated_at) VALUES
(0, '社区养老服务平台', 'https://via.placeholder.com/800x400/4A90E2/FFFFFF?text=社区养老服务', 'none', '', 100, 'active', NOW(3), NOW(3)),
(0, '专业护理团队', 'https://via.placeholder.com/800x400/50C878/FFFFFF?text=专业护理', 'none', '', 90, 'active', NOW(3), NOW(3)),
(0, '24小时服务热线', 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=24小时服务', 'none', '', 80, 'active', NOW(3), NOW(3));

-- 站点专属 Banner 示例（需要根据实际 station_id 调整）
-- 假设站点 ID 为 1
-- INSERT INTO banners (station_id, title, image_url, link_type, link_value, sort, status, created_at, updated_at) VALUES
-- (1, '本站点特色服务', 'https://via.placeholder.com/800x400/9B59B6/FFFFFF?text=特色服务', 'none', '', 100, 'active', NOW(), NOW());

-- 带链接的 Banner 示例
-- link_type 可选值: none, url, news, service
-- INSERT INTO banners (station_id, title, image_url, link_type, link_value, sort, status, created_at, updated_at) VALUES
-- (0, '最新资讯', 'https://via.placeholder.com/800x400/3498DB/FFFFFF?text=最新资讯', 'news', '1', 70, 'active', NOW(), NOW()),
-- (0, '热门服务', 'https://via.placeholder.com/800x400/E74C3C/FFFFFF?text=热门服务', 'service', '1', 60, 'active', NOW(), NOW()),
-- (0, '了解更多', 'https://via.placeholder.com/800x400/1ABC9C/FFFFFF?text=了解更多', 'url', 'https://example.com', 50, 'active', NOW(), NOW());

-- 非激活状态的 Banner 示例
-- INSERT INTO banners (station_id, title, image_url, link_type, link_value, sort, status, created_at, updated_at) VALUES
-- (0, '已下线的活动', 'https://via.placeholder.com/800x400/95A5A6/FFFFFF?text=已下线', 'none', '', 40, 'inactive', NOW(), NOW());
