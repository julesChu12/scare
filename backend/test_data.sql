-- 测试数据脚本：创建服务需求和任务
-- 用于前端 AAA MVP 测试

-- 插入服务需求（由elderly用户创建）
INSERT INTO service_requests (
    request_no, user_id, service_type, status, priority,
    description, submit_location_lat, submit_location_lng,
    contact_name, contact_phone, address, station_id,
    created_at, updated_at
) VALUES
-- 需求1：待派发的送餐服务
('REQ20240118001', 8, 'meal', 'submitted', 'normal',
 '需要送午餐，清淡少油，无辣', 40.053906, 116.363035,
 '张阿姨', '13900000001', '北京市昌平区霍营街道龙锦苑小区1号楼101室', 1,
 NOW(), NOW()),

-- 需求2：待派发的医疗服务
('REQ20240118002', 8, 'medical', 'submitted', 'high',
 '需要上门量血压和测血糖', 40.054506, 116.364035,
 '李爷爷', '13900000002', '北京市昌平区霍营街道龙锦苑小区2号楼203室', 1,
 NOW(), NOW()),

-- 需求3：待派发的清洁服务
('REQ20240118003', 9, 'cleaning', 'submitted', 'normal',
 '打扫卫生，整理房间', 40.055106, 116.365035,
 '王阿姨', '13900000003', '北京市昌平区霍营街道龙锦苑小区3号楼305室', 1,
 NOW(), NOW());

-- 插入任务（由station_manager派发）
INSERT INTO tasks (
    request_id, station_id, status,
    created_at, updated_at
)
SELECT
    id, station_id, 'dispatched',
    NOW(), NOW()
FROM service_requests
WHERE request_no IN ('REQ20240118001', 'REQ20240118002', 'REQ20240118003');

-- 更新服务需求状态为已派发
UPDATE service_requests
SET status = 'dispatched', updated_at = NOW()
WHERE request_no IN ('REQ20240118001', 'REQ20240118002', 'REQ20240118003');
