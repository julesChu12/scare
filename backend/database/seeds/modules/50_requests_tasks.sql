-- 50_requests_tasks.sql
-- 服务请求、任务、任务历史模块初始化

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

INSERT INTO `service_requests` (
    `id`, `request_no`, `user_id`, `service_type`, `status`,
    `description`, `submit_location_lat`, `submit_location_lng`,
    `contact_name`, `contact_phone`, `address`, `appointment_time`,
    `urgency`, `station_id`, `rating`, `feedback`,
    `created_at`, `updated_at`
) VALUES
-- 已完成的请求
(1, 'REQ-2026020701', 8, 'meal', 'completed',
 '需要午餐送餐服务，清淡饮食', 39.9042, 116.4074,
 '张大爷', '13800000008', '北京市朝阳区幸福小区1号楼101',
 '2026-02-07 12:00:00', 'normal', 1, 5, '服务很好，送餐准时',
 '2026-02-07 08:00:00', '2026-02-07 12:30:00'),

-- 已认领的请求
(2, 'REQ-2026020702', 9, 'medical', 'claimed',
 '需要陪同去医院复查', 39.9042, 116.4074,
 '李奶奶', '13800000009', '北京市朝阳区幸福小区2号楼202',
 '2026-02-08 09:00:00', 'urgent', 1, 0, NULL,
 '2026-02-07 10:00:00', '2026-02-07 10:30:00'),

-- 待认领的请求
(3, 'REQ-2026020703', 10, 'care', 'dispatched',
 '需要日常照护服务，协助洗漱', 39.9042, 116.4074,
 '王爷爷', '13800000010', '北京市朝阳区幸福小区3号楼303',
 '2026-02-08 08:00:00', 'normal', 2, 0, NULL,
 '2026-02-07 14:00:00', '2026-02-07 14:00:00'),

-- 已取消的请求
(4, 'REQ-2026020704', 11, 'other', 'cancelled',
 '需要产前护理指导', 39.9042, 116.4074,
 '孙女士', '13800000011', '北京市朝阳区康乐小区5号楼501',
 '2026-02-09 10:00:00', 'normal', 2, 0, NULL,
 '2026-02-07 15:00:00', '2026-02-07 15:30:00');

-- =====================================================
-- 7. 初始化任务分配（task_assignments）
-- 任务从服务请求派生，状态与请求保持一致
-- =====================================================
INSERT INTO `task_assignments` (
    `id`, `request_id`, `station_id`, `staff_id`, `status`,
    `claimed_at`, `completed_at`, `rating`, `feedback`, `images`,
    `created_at`, `updated_at`
) VALUES
-- 已完成的任务（对应请求1）
(1, 1, 1, 4, 'completed',
 '2026-02-07 08:30:00', '2026-02-07 12:30:00', 5, '服务很好，送餐准时', '[]',
 '2026-02-07 08:00:00', '2026-02-07 12:30:00'),

-- 已认领的任务（对应请求2）
(2, 2, 1, 5, 'claimed',
 '2026-02-07 10:30:00', NULL, 0, NULL, NULL,
 '2026-02-07 10:00:00', '2026-02-07 10:30:00'),

-- 待认领的任务（对应请求3）
(3, 3, 2, 0, 'dispatched',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-07 14:00:00', '2026-02-07 14:00:00'),

-- 已取消的任务（对应请求4，状态同步）
(4, 4, 2, 0, 'cancelled',
 NULL, NULL, 0, NULL, NULL,
 '2026-02-07 15:00:00', '2026-02-07 15:30:00');

-- =====================================================
-- 8. 初始化任务历史（task_histories）
-- 记录任务状态变更的审计日志
-- =====================================================
INSERT INTO `task_histories` (
    `id`, `task_id`, `request_id`, `action`, `operator_id`,
    `from_staff_id`, `to_staff_id`, `from_station_id`, `to_station_id`,
    `status_before`, `status_after`, `remark`, `created_at`
) VALUES
-- 任务1（已完成）的历史
(1, 1, 1, 'dispatched', 1, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-07 08:00:00'),
(2, 1, 1, 'claimed', 4, NULL, 4, NULL, NULL, 'dispatched', 'claimed', '员工认领任务', '2026-02-07 08:30:00'),
(3, 1, 1, 'completed', 4, 4, 4, NULL, NULL, 'claimed', 'completed', '服务完成', '2026-02-07 12:30:00'),

-- 任务2（已认领）的历史
(4, 2, 2, 'dispatched', 1, NULL, NULL, NULL, 1, NULL, 'dispatched', '系统自动派单', '2026-02-07 10:00:00'),
(5, 2, 2, 'claimed', 5, NULL, 5, NULL, NULL, 'dispatched', 'claimed', '员工认领任务', '2026-02-07 10:30:00'),

-- 任务3（待认领）的历史
(6, 3, 3, 'dispatched', 1, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-07 14:00:00'),

-- 任务4（已取消）的历史
(7, 4, 4, 'dispatched', 1, NULL, NULL, NULL, 2, NULL, 'dispatched', '系统自动派单', '2026-02-07 15:00:00'),
(8, 4, 4, 'cancelled', 11, NULL, NULL, NULL, NULL, 'dispatched', 'cancelled', '用户取消请求', '2026-02-07 15:30:00');
