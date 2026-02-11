-- =====================================================
-- 用户数据种子文件
-- 版本: v1.0.0
-- 说明: 用户、老年人档案、用户角色
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `user_roles`;
TRUNCATE TABLE `elderly_profiles`;
TRUNCATE TABLE `users`;
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 用户数据（密码统一为: Test@123）
-- =====================================================
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `avatar`, `gender`, `birth_date`, `id_card`, `role`, `station_id`, `status`) VALUES
-- 系统管理员
(1, '13800000001', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '系统管理员', 'admin@scare.com', NULL, 'male', '1985-06-15', '110105198506151234', 'admin', NULL, 'active'),
-- 霍营站点负责人
(2, '13800000002', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张站长', 'zhang@scare.com', NULL, 'male', '1980-03-20', '110105198003201234', 'station_manager', 1, 'active'),
-- 龙泽站点负责人
(3, '13800000003', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李站长', 'li@scare.com', NULL, 'female', '1982-08-10', '110105198208101234', 'station_manager', 2, 'active'),
-- 霍营站工作人员
(4, '13800000004', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '王小红', 'wang@scare.com', NULL, 'female', '1990-05-12', '110105199005121234', 'staff', 1, 'active'),
(5, '13800000005', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '刘小明', 'liu@scare.com', NULL, 'male', '1988-11-25', '110105198811251234', 'staff', 1, 'active'),
-- 龙泽站工作人员
(6, '13800000006', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '陈小花', 'chen@scare.com', NULL, 'female', '1992-02-18', '110105199202181234', 'staff', 2, 'active'),
(7, '13800000007', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '赵小强', 'zhao@scare.com', NULL, 'male', '1987-09-08', '110105198709081234', 'staff', 2, 'active'),
-- 老年人用户
(8, '13800000008', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张大爷', 'zhangdy@example.com', NULL, 'male', '1945-03-10', '110105194503101234', 'elderly', NULL, 'active'),
(9, '13800000009', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李奶奶', 'linn@example.com', NULL, 'female', '1948-07-22', '110105194807221234', 'elderly', NULL, 'active'),
(10, '13800000010', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '王大爷', 'wangdy@example.com', NULL, 'male', '1950-12-05', '110105195012051234', 'elderly', NULL, 'active'),
-- 家属用户
(11, '13800000011', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张小华', 'zhangxh@example.com', NULL, 'female', '1975-08-15', '110105197508151234', 'family', NULL, 'active'),
(12, '13800000012', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李小勇', 'lixy@example.com', NULL, 'male', '1978-04-20', '110105197804201234', 'family', NULL, 'active');

-- =====================================================
-- 老年人档案数据
-- =====================================================
INSERT INTO `elderly_profiles` (`user_id`, `address`, `latitude`, `longitude`, `health_info`, `emergency_contact_name`, `emergency_contact_phone`, `emergency_contact_relation`, `emergency_contact_user_id`) VALUES
(8, '北京市昌平区霍营街道华龙苑北里小区1号楼3单元501', 40.0500000, 116.3800000, '高血压，需定期服药', '张小华', '13800000011', '子女', 11),
(9, '北京市昌平区霍营街道龙锦苑东一区2号楼2单元302', 40.0550000, 116.3850000, '糖尿病，行动不便', '李小勇', '13800000012', '子女', 12),
(10, '北京市昌平区龙泽园街道龙泽苑西区5号楼1单元102', 40.0620000, 116.3920000, '身体健康', NULL, NULL, NULL, NULL);

-- =====================================================
-- 用户角色关联（仅 B 端用户需要）
-- =====================================================
INSERT INTO `user_roles` (`user_id`, `role`, `is_primary`, `status`) VALUES
(1, 'admin', TRUE, 'active'),
(2, 'station_manager', TRUE, 'active'),
(3, 'station_manager', TRUE, 'active'),
(4, 'staff', TRUE, 'active'),
(5, 'staff', TRUE, 'active'),
(6, 'staff', TRUE, 'active'),
(7, 'staff', TRUE, 'active'),
-- 多角色测试：李奶奶也是志愿者
(9, 'staff', FALSE, 'active');

-- 更新 primary_role
UPDATE `users` SET `primary_role` = `role` WHERE `role` IS NOT NULL AND `role` != '';

-- 验证数据
SELECT '用户数据' AS 'Table', COUNT(*) AS 'Count' FROM `users`
UNION ALL SELECT '老年人档案', COUNT(*) FROM `elderly_profiles`
UNION ALL SELECT '用户角色', COUNT(*) FROM `user_roles`;
