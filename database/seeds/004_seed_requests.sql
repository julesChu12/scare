-- =====================================================
-- 服务需求和任务数据种子文件
-- 版本: v1.0.0
-- 说明: 服务需求、任务分配、任务历史
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `task_histories`;
TRUNCATE TABLE `task_assignments`;
TRUNCATE TABLE `service_requests`;
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 服务需求数据
-- =====================================================
INSERT INTO `service_requests` (`request_no`, `user_id`, `service_type`, `status`, `description`, `submit_location_lat`, `submit_location_lng`, `contact_name`, `contact_phone`, `address`, `appointment_time`, `urgency`, `station_id`) VALUES
-- 已完成的需求
('REQ-2026011801', 8, 'meal', 'completed', '需要午餐送餐服务，饮食清淡，少油少盐', 40.0500000, 116.3800000, '张大爷', '13800000008', '北京市昌平区霍营街道华龙苑北里小区1号楼3单元501', '2026-01-18 11:30:00', 'normal', 1),
-- 处理中的需求
('REQ-2026011802', 9, 'medical', 'processing', '需要陪同去医院复查，预约明天上午9点', 40.0550000, 116.3850000, '李奶奶', '13800000009', '北京市昌平区霍营街道龙锦苑东一区2号楼2单元302', '2026-01-19 09:00:00', 'urgent', 1),
-- 待认领的需求
('REQ-2026011803', 10, 'care', 'dispatched', '需要日常照护服务，协助洗漱', 40.0620000, 116.3920000, '王大爷', '13800000010', '北京市昌平区龙泽园街道龙泽苑西区5号楼1单元102', NULL, 'normal', 2),
-- 待处理的需求
('REQ-2026011804', 8, 'cleaning', 'pending', '家里需要打扫卫生，周末方便', 40.0500000, 116.3800000, '张大爷', '13800000008', '北京市昌平区霍营街道华龙苑北里小区1号楼3单元501', '2026-01-20 10:00:00', 'normal', NULL);

-- =====================================================
-- 任务分配数据
-- =====================================================
INSERT INTO `task_assignments` (`request_id`, `station_id`, `staff_id`, `status`, `claimed_at`, `completed_at`, `rating`, `feedback`) VALUES
-- 已完成任务
(1, 1, 4, 'completed', '2026-01-18 10:30:00', '2026-01-18 12:00:00', 5, '服务态度好，送餐及时'),
-- 处理中任务
(2, 1, 5, 'processing', '2026-01-18 14:00:00', NULL, 0, NULL),
-- 待认领任务
(3, 2, NULL, 'dispatched', NULL, NULL, 0, NULL);

-- =====================================================
-- 任务历史数据
-- =====================================================
INSERT INTO `task_histories` (`task_id`, `request_id`, `action`, `operator_id`, `status_before`, `status_after`, `remark`) VALUES
-- 任务1的历史
(1, 1, 'dispatched', 2, NULL, 'dispatched', '系统自动派单到霍营站'),
(1, 1, 'claimed', 4, 'dispatched', 'claimed', '王小红认领任务'),
(1, 1, 'completed', 4, 'claimed', 'completed', '任务顺利完成'),
-- 任务2的历史
(2, 2, 'dispatched', 2, NULL, 'dispatched', '系统自动派单到霍营站'),
(2, 2, 'claimed', 5, 'dispatched', 'claimed', '刘小明认领任务'),
-- 任务3的历史
(3, 3, 'dispatched', 3, NULL, 'dispatched', '系统自动派单到龙泽站');

-- 验证数据
SELECT '服务需求' AS 'Table', COUNT(*) AS 'Count' FROM `service_requests`
UNION ALL SELECT '任务分配', COUNT(*) FROM `task_assignments`
UNION ALL SELECT '任务历史', COUNT(*) FROM `task_histories`;
