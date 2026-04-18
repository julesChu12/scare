-- =====================================================
-- 用户数据种子文件
-- 版本: v2.0.0
-- 说明: 用户、用户身份(user_identities)、客户档案、头像
-- =====================================================

SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

USE `scare_db`;

-- 清空现有数据
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `user_identities`;
TRUNCATE TABLE `customer_profiles`;
TRUNCATE TABLE `users`;
SET FOREIGN_KEY_CHECKS = 1;

-- =====================================================
-- 用户数据（密码统一为: Test@123）
-- 注意: id_card 为明文占位，id_card_hmac 留空，
--       后续可通过 cmd/tools/encrypt_seed 工具加密
-- =====================================================
INSERT INTO `users` (`id`, `phone`, `password_hash`, `name`, `email`, `avatar`, `gender`, `birth_date`, `id_card`, `id_card_hmac`, `id_card_masked`, `station_id`, `status`) VALUES
-- 系统管理员
(1, '13800000001', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '系统管理员', 'admin@scare.com', '/static/b_end/20260329/头像.jpg', 'male', '1985-06-15', '110105198506151234', '', '1101**********1234', NULL, 'active'),
-- 霍营站点负责人
(2, '13800000002', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张站长', 'zhang@scare.com', '/static/b_end/20260329/头像.jpg', 'male', '1980-03-20', '110105198003201234', '', '1101**********1234', 1, 'active'),
-- 龙泽站点负责人
(3, '13800000003', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李站长', 'li@scare.com', '/static/b_end/20260329/头像.jpg', 'female', '1982-08-10', '110105198208101234', '', '1101**********1234', 2, 'active'),
-- 霍营站工作人员
(4, '13800000004', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '王小红', 'wang@scare.com', '/static/b_end/20260329/头像.jpg', 'female', '1990-05-12', '110105199005121234', '', '1101**********1234', 1, 'active'),
(5, '13800000005', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '刘小明', 'liu@scare.com', '/static/b_end/20260329/头像.jpg', 'male', '1988-11-25', '110105198811251234', '', '1101**********1234', 1, 'active'),
-- 龙泽站工作人员
(6, '13800000006', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '陈小花', 'chen@scare.com', '/static/b_end/20260329/头像.jpg', 'female', '1992-02-18', '110105199202181234', '', '1101**********1234', 2, 'active'),
(7, '13800000007', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '赵小强', 'zhao@scare.com', '/static/b_end/20260329/头像.jpg', 'male', '1987-09-08', '110105198709081234', '', '1101**********1234', 2, 'active'),
-- 老年人用户
(8, '13800000008', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张大爷', 'zhangdy@example.com', '/static/b_end/20260329/头像.jpg', 'male', '1945-03-10', '110105194503101234', '', '1101**********1234', NULL, 'active'),
(9, '13800000009', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李奶奶', 'linn@example.com', '/static/b_end/20260329/头像.jpg', 'female', '1948-07-22', '110105194807221234', '', '1101**********1234', NULL, 'active'),
(10, '13800000010', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '王大爷', 'wangdy@example.com', '/static/b_end/20260329/头像.jpg', 'male', '1950-12-05', '110105195012051234', '', '1101**********1234', NULL, 'active'),
-- 家属用户
(11, '13800000011', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '张小华', 'zhangxh@example.com', '/static/b_end/20260329/头像.jpg', 'female', '1975-08-15', '110105197508151234', '', '1101**********1234', NULL, 'active'),
(12, '13800000012', '$2a$10$zwHpVzGXPcbQAV4Tb5KoAuDUKtMhaIEo/lfn1oqraECMW5XadrGKK', '李小勇', 'lixy@example.com', '/static/b_end/20260329/头像.jpg', 'male', '1978-04-20', '110105197804201234', '', '1101**********1234', NULL, 'active');

-- =====================================================
-- 用户身份数据（认证系统依赖此表）
-- identity_type: admin/station_manager/staff/elderly/family
-- =====================================================
INSERT INTO `user_identities` (`user_id`, `identity_type`, `is_primary`, `station_id`, `status`) VALUES
-- B端身份
(1, 'admin', TRUE, NULL, 'active'),
(2, 'station_manager', TRUE, 1, 'active'),
(3, 'station_manager', TRUE, 2, 'active'),
(4, 'staff', TRUE, 1, 'active'),
(5, 'staff', TRUE, 1, 'active'),
(6, 'staff', TRUE, 2, 'active'),
(7, 'staff', TRUE, 2, 'active'),
-- C端身份
(8, 'elderly', TRUE, NULL, 'active'),
(9, 'elderly', TRUE, NULL, 'active'),
(10, 'elderly', TRUE, NULL, 'active'),
(11, 'family', TRUE, NULL, 'active'),
(12, 'family', TRUE, NULL, 'active'),
-- 多身份测试：李奶奶也是 B端志愿者
(9, 'staff', FALSE, 1, 'active');

-- =====================================================
-- C端客户档案（C端登录依赖此表）
-- =====================================================
INSERT INTO `customer_profiles` (`user_id`, `id_card`, `address`, `latitude`, `longitude`, `customer_type`, `gender`, `birth_date`, `health_status`, `disability_level`, `medical_history`, `special_needs`, `emergency_contact`) VALUES
(8, '110105194503101234', '北京市昌平区霍营街道华龙苑北里小区1号楼3单元501', 40.0500, 116.3800, 'elderly', 'male', '1945-03-10', '高血压', '轻度', '高血压病史10年，服用降压药控制良好', '饮食需少油少盐，行动缓慢需注意防摔', '{"name":"张小华","phone":"13800000011","relation":"子女","user_id":11}'),
(9, '110105194807221234', '北京市昌平区霍营街道龙锦苑东一区2号楼2单元302', 40.0550, 116.3850, 'elderly', 'female', '1948-07-22', '糖尿病', '中度', '二型糖尿病8年，注射胰岛素，骨质疏松', '需协助洗浴，轮椅出行，定期复查骨密度', '{"name":"李小勇","phone":"13800000012","relation":"子女","user_id":12}'),
(10, '110105195012051234', '北京市昌平区龙泽园街道龙泽苑西区5号楼1单元102', 40.0620, 116.3920, 'elderly', 'male', '1950-12-05', '健康', '自理', '无重大疾病史，身体健康', '无特殊需求', NULL),
(11, '110105197508151234', '北京市昌平区霍营街道华龙苑北里小区2号楼', 40.0510, 116.3810, 'family', 'female', '1975-08-15', NULL, NULL, NULL, NULL, NULL),
(12, '110105197804201234', '北京市昌平区霍营街道龙锦苑东一区3号楼', 40.0540, 116.3840, 'family', 'male', '1978-04-20', NULL, NULL, NULL, NULL, NULL);

-- 验证数据
SELECT '用户数据' AS 'Table', COUNT(*) AS 'Count' FROM `users`
UNION ALL SELECT '用户身份', COUNT(*) FROM `user_identities`
UNION ALL SELECT '客户档案', COUNT(*) FROM `customer_profiles`;
