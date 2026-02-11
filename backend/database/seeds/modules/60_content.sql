-- 60_content.sql
-- 内容模块初始化（Banner + News）

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `banners` (`id`, `station_id`, `title`, `image_url`, `link_type`, `link_value`, `sort`, `status`, `created_at`, `updated_at`) VALUES
(1, 0, '社区养老服务平台', 'https://via.placeholder.com/800x400/4A90E2/FFFFFF?text=社区养老服务', 'none', '', 100, 'active', NOW(3), NOW(3)),
(2, 0, '专业护理团队', 'https://via.placeholder.com/800x400/50C878/FFFFFF?text=专业护理', 'none', '', 90, 'active', NOW(3), NOW(3)),
(3, 0, '24小时服务热线', 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=24小时服务', 'none', '', 80, 'active', NOW(3), NOW(3));

INSERT INTO `news` (`id`, `title`, `summary`, `content`, `cover_url`, `type`, `status`, `station_id`, `author_id`, `publish_at`, `view_count`, `created_at`, `updated_at`) VALUES
(1, '社区养老服务中心正式启用', '首批养老服务站点完成部署并投入使用。', '社区养老服务中心正式启用，为居民提供送餐、陪诊、照护等服务。', 'https://via.placeholder.com/600x300/4A90E2/FFFFFF?text=News+1', 'news', 'published', 0, 1, NOW(), 120, NOW(3), NOW(3)),
(2, '春季义诊活动公告', '本周六开展春季社区义诊，请有需要的居民提前预约。', '本周六上午9点到12点，在朝阳区幸福街道服务站开展义诊活动。', 'https://via.placeholder.com/600x300/50C878/FFFFFF?text=News+2', 'notice', 'published', 0, 1, NOW(), 80, NOW(3), NOW(3));
