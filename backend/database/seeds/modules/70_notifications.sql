-- 70_notifications.sql
-- 通知模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `notifications` (`id`, `user_id`, `title`, `body`, `type`, `related_id`, `related_type`, `channel`, `send_status`, `is_read`, `retry_count`, `created_at`, `updated_at`) VALUES
(1, 2, '系统通知', '欢迎使用社区养老服务平台。', 'system', 0, 'system', 'in_app', 'sent', 0, 0, NOW(3), NOW(3)),
(2, 4, '任务提醒', '您有新的待处理任务，请及时查看。', 'task', 3, 'task', 'in_app', 'sent', 0, 0, NOW(3), NOW(3));
